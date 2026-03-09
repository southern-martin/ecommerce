#!/usr/bin/env bash
#
# seed.sh — Populate the ecommerce dev environment with realistic demo data.
#
# Prerequisites: curl, jq, psql (for role updates)
# Usage:         bash scripts/seed.sh
#
# The script is idempotent: re-running it will skip already-created resources
# (the auth service returns 409 for duplicate emails, etc.).

set -euo pipefail

# ─── Colours ──────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Colour

# ─── Service base URLs (docker-compose ports) ────────────────────
AUTH_URL="${AUTH_URL:-http://localhost:28090}"
USER_URL="${USER_URL:-http://localhost:28091}"
PRODUCT_URL="${PRODUCT_URL:-http://localhost:28081}"
CART_URL="${CART_URL:-http://localhost:28082}"
ORDER_URL="${ORDER_URL:-http://localhost:28083}"
PROMOTION_URL="${PROMOTION_URL:-http://localhost:28093}"
REVIEW_URL="${REVIEW_URL:-http://localhost:28086}"
CMS_URL="${CMS_URL:-http://localhost:28099}"
TAX_URL="${TAX_URL:-http://localhost:28098}"
SHIPPING_URL="${SHIPPING_URL:-http://localhost:28095}"
KONG_URL="${KONG_URL:-http://localhost:28000}"

# Database (for role updates — register always creates "buyer")
POSTGRES_HOST="${POSTGRES_HOST:-localhost}"
POSTGRES_PORT="${POSTGRES_PORT:-15432}"
POSTGRES_USER="${POSTGRES_USER:-ecommerce}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-ecommerce_secret}"

# ─── Helpers ──────────────────────────────────────────────────────

info()    { echo -e "${CYAN}[INFO]${NC}  $*"; }
success() { echo -e "${GREEN}[OK]${NC}    $*"; }
warn()    { echo -e "${YELLOW}[WARN]${NC}  $*"; }
fail()    { echo -e "${RED}[FAIL]${NC}  $*"; }
section() { echo -e "\n${BOLD}━━━ $* ━━━${NC}"; }

# Curl wrapper that returns the HTTP body. Sets $HTTP_CODE as a side-effect.
api() {
  local method="$1" url="$2"
  shift 2
  local tmpfile
  tmpfile=$(mktemp)
  HTTP_CODE=$(curl -s -o "$tmpfile" -w '%{http_code}' -X "$method" "$url" \
    -H "Content-Type: application/json" "$@")
  cat "$tmpfile"
  rm -f "$tmpfile"
}

# ─── 1. Wait for services ────────────────────────────────────────
section "Waiting for services to be healthy"

wait_for_service() {
  local name="$1" url="$2" max_attempts="${3:-30}"
  local attempt=0
  while true; do
    attempt=$((attempt + 1))
    if curl -sf "${url}/health" >/dev/null 2>&1; then
      success "$name is healthy"
      return 0
    fi
    if [ "$attempt" -ge "$max_attempts" ]; then
      fail "$name did not become healthy after ${max_attempts}s"
      return 1
    fi
    sleep 1
  done
}

wait_for_service "auth"      "$AUTH_URL"
wait_for_service "product"   "$PRODUCT_URL"
wait_for_service "promotion" "$PROMOTION_URL"
wait_for_service "cms"       "$CMS_URL"
wait_for_service "tax"       "$TAX_URL"

# ─── 2. Create Users ─────────────────────────────────────────────
section "Creating users"

# Register a user; returns JSON with user_id and access_token.
# If the user already exists (409), log in instead.
register_or_login() {
  local email="$1" password="$2"
  local body
  body=$(api POST "${AUTH_URL}/api/v1/auth/register" \
    -d "{\"email\":\"${email}\",\"password\":\"${password}\"}")

  if [ "$HTTP_CODE" = "201" ]; then
    success "Registered ${email}"
    echo "$body"
    return
  fi

  if [ "$HTTP_CODE" = "409" ]; then
    warn "${email} already exists — logging in"
    body=$(api POST "${AUTH_URL}/api/v1/auth/login" \
      -d "{\"email\":\"${email}\",\"password\":\"${password}\"}")
    if [ "$HTTP_CODE" = "200" ]; then
      echo "$body"
      return
    fi
  fi

  fail "Could not register/login ${email} (HTTP ${HTTP_CODE})"
  echo "$body" >&2
  return 1
}

# Helper: update user role via psql (auth service always registers as "buyer")
update_role() {
  local user_id="$1" role="$2"
  PGPASSWORD="$POSTGRES_PASSWORD" psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" \
    -U "$POSTGRES_USER" -d ecommerce_auth -tAc \
    "UPDATE auth_users SET role = '${role}' WHERE id = '${user_id}';" >/dev/null 2>&1 || true
}

# --- Admin ---
ADMIN_JSON=$(register_or_login "admin@ecommerce.com" "Admin123!")
ADMIN_TOKEN=$(echo "$ADMIN_JSON" | jq -r '.access_token')
ADMIN_ID=$(echo "$ADMIN_JSON" | jq -r '.user_id')
update_role "$ADMIN_ID" "admin"
info "Admin ID: ${ADMIN_ID}"

# --- Sellers ---
declare -a SELLER_IDS SELLER_TOKENS
for i in 1 2 3; do
  SELLER_JSON=$(register_or_login "seller${i}@demo.com" "Seller123!")
  SELLER_TOKENS[$i]=$(echo "$SELLER_JSON" | jq -r '.access_token')
  SELLER_IDS[$i]=$(echo "$SELLER_JSON" | jq -r '.user_id')
  update_role "${SELLER_IDS[$i]}" "seller"
  info "Seller${i} ID: ${SELLER_IDS[$i]}"
done

# --- Buyers ---
declare -a BUYER_IDS BUYER_TOKENS
for i in 1 2 3; do
  BUYER_JSON=$(register_or_login "buyer${i}@demo.com" "Buyer123!")
  BUYER_TOKENS[$i]=$(echo "$BUYER_JSON" | jq -r '.access_token')
  BUYER_IDS[$i]=$(echo "$BUYER_JSON" | jq -r '.user_id')
  info "Buyer${i} ID: ${BUYER_IDS[$i]}"
done

# Re-login admin & sellers to pick up updated roles in JWT
info "Re-logging in admin & sellers to refresh JWT roles..."
ADMIN_JSON=$(api POST "${AUTH_URL}/api/v1/auth/login" \
  -d '{"email":"admin@ecommerce.com","password":"Admin123!"}')
ADMIN_TOKEN=$(echo "$ADMIN_JSON" | jq -r '.access_token')

for i in 1 2 3; do
  SELLER_JSON=$(api POST "${AUTH_URL}/api/v1/auth/login" \
    -d "{\"email\":\"seller${i}@demo.com\",\"password\":\"Seller123!\"}")
  SELLER_TOKENS[$i]=$(echo "$SELLER_JSON" | jq -r '.access_token')
done

# ─── 3. Create Categories (via Kong — admin) ─────────────────────
section "Creating product categories"

create_category() {
  local name="$1" sort_order="$2"
  local body
  body=$(api POST "${KONG_URL}/api/v1/admin/categories" \
    -H "Authorization: Bearer ${ADMIN_TOKEN}" \
    -d "{\"name\":\"${name}\",\"sort_order\":${sort_order}}")

  if [ "$HTTP_CODE" = "201" ]; then
    success "Category: ${name}"
    echo "$body" | jq -r '.id'
  elif echo "$body" | jq -e '.id' >/dev/null 2>&1; then
    warn "Category '${name}' may already exist"
    echo "$body" | jq -r '.id'
  else
    # Try to find existing category by listing
    warn "Category '${name}' creation returned HTTP ${HTTP_CODE}, looking up..."
    local list
    list=$(api GET "${PRODUCT_URL}/api/v1/categories")
    local cat_id
    cat_id=$(echo "$list" | jq -r --arg n "$name" '.categories[]? | select(.name==$n) | .id' 2>/dev/null | head -1)
    if [ -n "$cat_id" ] && [ "$cat_id" != "null" ]; then
      success "Found existing category '${name}': ${cat_id}"
      echo "$cat_id"
    else
      fail "Could not create or find category '${name}'"
      echo ""
    fi
  fi
}

CAT_ELECTRONICS=$(create_category "Electronics" 1)
CAT_CLOTHING=$(create_category "Clothing" 2)
CAT_HOME=$(create_category "Home & Kitchen" 3)
CAT_SPORTS=$(create_category "Sports" 4)
CAT_BOOKS=$(create_category "Books" 5)

# ─── 4. Create Products (via Kong — sellers) ─────────────────────
section "Creating products"

create_product() {
  local token="$1" category_id="$2" name="$3" description="$4" price_cents="$5" stock="$6"
  shift 6
  local tags_json="${1:-[]}"

  local payload
  payload=$(jq -n \
    --arg cat "$category_id" \
    --arg name "$name" \
    --arg desc "$description" \
    --argjson price "$price_cents" \
    --argjson stock "$stock" \
    --argjson tags "$tags_json" \
    '{
      category_id: $cat,
      name: $name,
      description: $desc,
      base_price_cents: $price,
      currency: "USD",
      product_type: "simple",
      stock_quantity: $stock,
      tags: $tags
    }')

  local body
  body=$(api POST "${KONG_URL}/api/v1/seller/products" \
    -H "Authorization: Bearer ${token}" \
    -d "$payload")

  if [ "$HTTP_CODE" = "201" ]; then
    success "Product: ${name}"
    echo "$body" | jq -r '.id // .product.id // empty' 2>/dev/null
  else
    warn "Product '${name}' — HTTP ${HTTP_CODE} (may already exist)"
    echo ""
  fi
}

# --- Seller 1: Electronics ---
info "Seller 1 → Electronics"
create_product "${SELLER_TOKENS[1]}" "$CAT_ELECTRONICS" \
  "ProBook Laptop 15\"" \
  "15.6-inch Full HD display, Intel Core i7, 16GB RAM, 512GB SSD. Perfect for professionals and creators." \
  129999 50 '["laptop","computer","electronics"]'

create_product "${SELLER_TOKENS[1]}" "$CAT_ELECTRONICS" \
  "SoundWave Wireless Headphones" \
  "Active noise cancellation, 30-hour battery life, Bluetooth 5.3. Premium sound quality for music lovers." \
  19999 200 '["headphones","audio","wireless"]'

create_product "${SELLER_TOKENS[1]}" "$CAT_ELECTRONICS" \
  "Galaxy Ultra Smartphone" \
  "6.7-inch AMOLED display, 108MP camera, 5G connectivity, 256GB storage. The ultimate mobile experience." \
  99999 150 '["smartphone","mobile","5g"]'

create_product "${SELLER_TOKENS[1]}" "$CAT_ELECTRONICS" \
  "SlimTab Pro 11" \
  "11-inch Liquid Retina display, M2 chip, 128GB storage. Lightweight and powerful tablet for on-the-go." \
  79999 80 '["tablet","portable","touchscreen"]'

create_product "${SELLER_TOKENS[1]}" "$CAT_ELECTRONICS" \
  "FitTrack Smartwatch" \
  "Heart rate monitor, GPS, 7-day battery, water resistant to 50m. Your personal health companion." \
  29999 300 '["smartwatch","wearable","fitness"]'

# --- Seller 2: Clothing ---
info "Seller 2 → Clothing"
create_product "${SELLER_TOKENS[2]}" "$CAT_CLOTHING" \
  "Classic Cotton T-Shirt" \
  "100% organic cotton, pre-shrunk, available in 12 colours. Comfortable everyday essential." \
  2499 500 '["tshirt","cotton","casual"]'

create_product "${SELLER_TOKENS[2]}" "$CAT_CLOTHING" \
  "Slim Fit Denim Jeans" \
  "Premium stretch denim, 5-pocket design, machine washable. Modern slim fit for everyday style." \
  5999 300 '["jeans","denim","casual"]'

create_product "${SELLER_TOKENS[2]}" "$CAT_CLOTHING" \
  "Urban Runner Sneakers" \
  "Lightweight mesh upper, memory foam insole, rubber outsole. Ideal for running and casual wear." \
  8999 200 '["sneakers","shoes","running"]'

create_product "${SELLER_TOKENS[2]}" "$CAT_CLOTHING" \
  "All-Weather Jacket" \
  "Water-resistant shell, breathable lining, adjustable hood. Stay dry and comfortable in any weather." \
  12999 100 '["jacket","outerwear","waterproof"]'

create_product "${SELLER_TOKENS[2]}" "$CAT_CLOTHING" \
  "Floral Summer Dress" \
  "Lightweight chiffon fabric, floral print, knee-length. Perfect for warm-weather occasions." \
  4999 250 '["dress","summer","floral"]'

# --- Seller 3: Home & Kitchen ---
info "Seller 3 → Home & Kitchen"
create_product "${SELLER_TOKENS[3]}" "$CAT_HOME" \
  "BrewMaster Coffee Maker" \
  "12-cup programmable drip coffee maker, built-in grinder, thermal carafe. Fresh coffee every morning." \
  8999 120 '["coffee","kitchen","appliance"]'

create_product "${SELLER_TOKENS[3]}" "$CAT_HOME" \
  "PowerBlend Pro Blender" \
  "1200W motor, 6 stainless steel blades, 64oz BPA-free jar. Crushes ice and blends smoothies in seconds." \
  6999 150 '["blender","kitchen","appliance"]'

create_product "${SELLER_TOKENS[3]}" "$CAT_HOME" \
  "LED Architect Desk Lamp" \
  "Adjustable colour temperature, USB charging port, touch controls. Reduce eye strain while working." \
  4499 200 '["lamp","desk","lighting"]'

create_product "${SELLER_TOKENS[3]}" "$CAT_HOME" \
  "Luxury Velvet Throw Pillow" \
  "18x18 inch, removable zippered cover, hypoallergenic fill. Add a touch of elegance to any room." \
  2999 400 '["pillow","decor","living-room"]'

create_product "${SELLER_TOKENS[3]}" "$CAT_HOME" \
  "Modern Minimalist Wall Clock" \
  "12-inch silent quartz movement, brushed aluminium frame. Clean design that suits any interior." \
  3499 180 '["clock","decor","wall-art"]'

# ─── 5. Create Promotions ────────────────────────────────────────
section "Creating promotions"

create_coupon() {
  local code="$1" type="$2" discount_value="$3" expires_at="$4"
  local body
  body=$(api POST "${PROMOTION_URL}/api/v1/admin/promotions/coupons" \
    -H "Authorization: Bearer ${ADMIN_TOKEN}" \
    -H "X-User-ID: ${ADMIN_ID}" \
    -d "{
      \"code\": \"${code}\",
      \"type\": \"${type}\",
      \"discount_value\": ${discount_value},
      \"min_order_cents\": 1000,
      \"usage_limit\": 1000,
      \"per_user_limit\": 3,
      \"scope\": \"global\",
      \"expires_at\": \"${expires_at}\"
    }")

  if [ "$HTTP_CODE" = "201" ]; then
    success "Coupon: ${code}"
  else
    warn "Coupon '${code}' — HTTP ${HTTP_CODE} (may already exist)"
  fi
}

EXPIRES_AT=$(date -u -v+1y '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null || date -u -d '+1 year' '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null || echo "2027-12-31T23:59:59Z")

create_coupon "WELCOME10" "percentage" 10 "$EXPIRES_AT"
create_coupon "SUMMER20"  "percentage" 20 "$EXPIRES_AT"

# ─── 6. Create CMS Pages ─────────────────────────────────────────
section "Creating CMS pages"

create_page() {
  local title="$1" content_html="$2" meta_title="$3" meta_desc="$4"
  local payload
  payload=$(jq -n \
    --arg title "$title" \
    --arg html "$content_html" \
    --arg mt "$meta_title" \
    --arg md "$meta_desc" \
    '{title: $title, content_html: $html, meta_title: $mt, meta_description: $md}')

  local body
  body=$(api POST "${CMS_URL}/api/v1/admin/pages" \
    -H "Authorization: Bearer ${ADMIN_TOKEN}" \
    -H "X-User-ID: ${ADMIN_ID}" \
    -d "$payload")

  if [ "$HTTP_CODE" = "201" ]; then
    local page_id
    page_id=$(echo "$body" | jq -r '.page.id // .id // empty' 2>/dev/null)
    success "Page: ${title}"

    # Publish the page
    if [ -n "$page_id" ] && [ "$page_id" != "null" ]; then
      api PATCH "${CMS_URL}/api/v1/admin/pages/${page_id}/publish" \
        -H "Authorization: Bearer ${ADMIN_TOKEN}" \
        -H "X-User-ID: ${ADMIN_ID}" >/dev/null 2>&1
      info "  Published page ${page_id}"
    fi
  else
    warn "Page '${title}' — HTTP ${HTTP_CODE} (may already exist)"
  fi
}

create_page "About Us" \
  "<h1>About Us</h1><p>We are a modern e-commerce platform connecting buyers and sellers worldwide. Our mission is to provide a seamless shopping experience with the widest selection of quality products at competitive prices.</p><h2>Our Story</h2><p>Founded in 2024, we started with a simple idea: make online shopping better for everyone. Today we serve thousands of customers across the globe.</p>" \
  "About Us - Ecommerce" \
  "Learn about our e-commerce platform, our mission, and how we connect buyers and sellers worldwide."

create_page "Terms of Service" \
  "<h1>Terms of Service</h1><h2>1. Acceptance of Terms</h2><p>By accessing and using this platform, you accept and agree to be bound by the terms and provision of this agreement.</p><h2>2. Use of Service</h2><p>You agree to use the service only for purposes that are permitted by these Terms and any applicable law or regulation.</p><h2>3. User Accounts</h2><p>You are responsible for maintaining the confidentiality of your account credentials and for all activities that occur under your account.</p><h2>4. Limitation of Liability</h2><p>In no event shall the platform be liable for any indirect, incidental, special, consequential, or punitive damages.</p>" \
  "Terms of Service - Ecommerce" \
  "Read our terms of service to understand your rights and responsibilities when using our platform."

create_page "Privacy Policy" \
  "<h1>Privacy Policy</h1><h2>Information We Collect</h2><p>We collect information you provide directly, such as your name, email address, and shipping details when you create an account or place an order.</p><h2>How We Use Your Information</h2><p>We use the information to process orders, communicate with you, and improve our services.</p><h2>Data Security</h2><p>We implement appropriate security measures to protect your personal information against unauthorized access, alteration, or destruction.</p><h2>Contact Us</h2><p>If you have questions about this privacy policy, please contact us at privacy@ecommerce.com.</p>" \
  "Privacy Policy - Ecommerce" \
  "Our privacy policy explains how we collect, use, and protect your personal information."

# ─── 7. Create Tax Rules ─────────────────────────────────────────
section "Creating tax rules"

# First, fetch existing zones
info "Fetching tax zones..."
ZONES_JSON=$(api GET "${TAX_URL}/api/v1/tax/zones")

get_zone_id() {
  local country="$1" state="${2:-}"
  local zone_id
  if [ -n "$state" ]; then
    zone_id=$(echo "$ZONES_JSON" | jq -r --arg cc "$country" --arg sc "$state" \
      '.zones[]? | select(.country_code==$cc and .state_code==$sc) | .id' 2>/dev/null | head -1)
  else
    zone_id=$(echo "$ZONES_JSON" | jq -r --arg cc "$country" \
      '.zones[]? | select(.country_code==$cc and (.state_code=="" or .state_code==null)) | .id' 2>/dev/null | head -1)
  fi
  echo "$zone_id"
}

create_tax_rule() {
  local zone_id="$1" tax_name="$2" rate="$3"
  if [ -z "$zone_id" ] || [ "$zone_id" = "null" ]; then
    warn "Skipping tax rule '${tax_name}' — zone not found"
    return
  fi

  local body
  body=$(api POST "${TAX_URL}/api/v1/admin/tax/rules" \
    -H "Authorization: Bearer ${ADMIN_TOKEN}" \
    -H "X-User-ID: ${ADMIN_ID}" \
    -d "{
      \"zone_id\": \"${zone_id}\",
      \"tax_name\": \"${tax_name}\",
      \"rate\": ${rate},
      \"category\": \"general\",
      \"inclusive\": false
    }")

  if [ "$HTTP_CODE" = "201" ]; then
    success "Tax rule: ${tax_name} (${rate}%)"
  else
    warn "Tax rule '${tax_name}' — HTTP ${HTTP_CODE} (may already exist)"
  fi
}

# Look up zone IDs
ZONE_US=$(get_zone_id "US" "")
ZONE_CA=$(get_zone_id "US" "CA")
ZONE_NY=$(get_zone_id "US" "NY")
ZONE_TX=$(get_zone_id "US" "TX")

# If zones don't exist yet, log a warning — the tax service may need its own
# seed data for zones first. We still attempt to create rules.
if [ -z "$ZONE_US" ] || [ "$ZONE_US" = "null" ]; then
  warn "US federal zone not found. Tax zones may need to be seeded first."
  warn "Attempting to create rules with placeholder zone IDs..."
  ZONE_US="zone-us"
  ZONE_CA="zone-us-ca"
  ZONE_NY="zone-us-ny"
  ZONE_TX="zone-us-tx"
fi

create_tax_rule "$ZONE_US" "US Federal Tax" 0
create_tax_rule "$ZONE_CA" "California State Tax" 7.25
create_tax_rule "$ZONE_NY" "New York State Tax" 8
create_tax_rule "$ZONE_TX" "Texas State Tax" 6.25

# ─── Done ─────────────────────────────────────────────────────────
section "Seed complete"
echo ""
echo -e "${GREEN}${BOLD}Summary:${NC}"
echo "  Users:       1 admin, 3 sellers, 3 buyers"
echo "  Categories:  5 (Electronics, Clothing, Home & Kitchen, Sports, Books)"
echo "  Products:    15 (5 per seller)"
echo "  Coupons:     2 (WELCOME10, SUMMER20)"
echo "  CMS pages:   3 (About Us, Terms of Service, Privacy Policy)"
echo "  Tax rules:   4 (US Federal, CA, NY, TX)"
echo ""
echo -e "${CYAN}Admin login:${NC}   admin@ecommerce.com / Admin123!"
echo -e "${CYAN}Seller login:${NC}  seller1@demo.com / Seller123!"
echo -e "${CYAN}Buyer login:${NC}   buyer1@demo.com / Buyer123!"
echo ""
