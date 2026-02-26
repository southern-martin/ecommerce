#!/bin/bash
# =============================================================================
# Deploy all ecommerce microservices to GCP Cloud Run
# =============================================================================
#
# Usage:
#   ./deploy.sh <project-id> <region>
#
# Example:
#   ./deploy.sh my-ecommerce-project us-central1
#
# Prerequisites:
#   - gcloud CLI installed and authenticated
#   - Docker images already built and pushed to GCR
#   - Cloud SQL, Memorystore (Redis), and other managed services provisioned
#
# =============================================================================

set -euo pipefail

PROJECT="${1:?Usage: $0 <project-id> <region>}"
REGION="${2:?Usage: $0 <project-id> <region>}"

IMAGE_BASE="gcr.io/${PROJECT}"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC}  $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC}  $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# =============================================================================
# Environment variables (override via environment or .env file)
# =============================================================================
POSTGRES_HOST="${POSTGRES_HOST:-/cloudsql/${PROJECT}:${REGION}:ecommerce-db}"
POSTGRES_USER="${POSTGRES_USER:-ecommerce}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-}"
REDIS_URL="${REDIS_URL:-redis://10.0.0.3:6379}"
NATS_URL="${NATS_URL:-nats://nats:4222}"
ELASTICSEARCH_URL="${ELASTICSEARCH_URL:-http://elasticsearch:9200}"
JWT_SECRET="${JWT_SECRET:-}"
STRIPE_SECRET_KEY="${STRIPE_SECRET_KEY:-}"
STRIPE_WEBHOOK_SECRET="${STRIPE_WEBHOOK_SECRET:-}"

# Common env vars for all services
COMMON_ENV="POSTGRES_HOST=${POSTGRES_HOST},POSTGRES_PORT=5432,POSTGRES_USER=${POSTGRES_USER},REDIS_URL=${REDIS_URL},NATS_URL=${NATS_URL},LOG_LEVEL=info"

if [ -n "${POSTGRES_PASSWORD}" ]; then
  COMMON_ENV="${COMMON_ENV},POSTGRES_PASSWORD=${POSTGRES_PASSWORD}"
fi

# =============================================================================
# Service definitions
# Format: name:http_port:db_name:access:extra_env
#   access: public = --allow-unauthenticated
#           internal = --no-allow-unauthenticated
# =============================================================================
declare -a SERVICES=(
  "auth:8090:ecommerce_auth:public:JWT_SECRET=${JWT_SECRET}"
  "user:8091:ecommerce_users:public:AUTH_GRPC_ADDR=auth:9090"
  "product:8081:ecommerce_products:public:"
  "cart:8082:ecommerce_cart:internal:"
  "order:8083:ecommerce_orders:internal:"
  "payment:8084:ecommerce_payments:internal:STRIPE_SECRET_KEY=${STRIPE_SECRET_KEY},STRIPE_WEBHOOK_SECRET=${STRIPE_WEBHOOK_SECRET},PLATFORM_COMMISSION_RATE=0.10"
  "tax:8098:ecommerce_tax:internal:"
  "shipping:8085:ecommerce_shipping:internal:"
  "return:8086:ecommerce_returns:internal:"
  "search:8087:ecommerce_search:public:ELASTICSEARCH_URL=${ELASTICSEARCH_URL}"
  "review:8088:ecommerce_reviews:internal:"
  "media:8089:ecommerce_media:internal:"
  "notification:8092:ecommerce_notifications:internal:"
  "cms:8099:ecommerce_cms:internal:"
  "promotion:8093:ecommerce_promotions:internal:"
  "loyalty:8096:ecommerce_loyalty:internal:"
  "affiliate:8097:ecommerce_affiliates:internal:"
  "chat:8094:ecommerce_chat:internal:"
  "ai:8095:ecommerce_ai:internal:"
)

# =============================================================================
# Deploy each service
# =============================================================================
FAILED=()
SUCCEEDED=()

for service_def in "${SERVICES[@]}"; do
  IFS=':' read -r NAME HTTP_PORT DB_NAME ACCESS EXTRA_ENV <<< "${service_def}"

  log_info "Deploying ${NAME}-service..."

  # Build env vars string
  ENV_VARS="${COMMON_ENV},HTTP_PORT=${HTTP_PORT},DB_NAME=${DB_NAME}"
  if [ -n "${EXTRA_ENV}" ]; then
    ENV_VARS="${ENV_VARS},${EXTRA_ENV}"
  fi

  # Set access flag
  if [ "${ACCESS}" = "public" ]; then
    ACCESS_FLAG="--allow-unauthenticated"
  else
    ACCESS_FLAG="--no-allow-unauthenticated"
  fi

  # Deploy to Cloud Run
  if gcloud run deploy "${NAME}-service" \
    --project="${PROJECT}" \
    --image="${IMAGE_BASE}/${NAME}:latest" \
    --platform=managed \
    --region="${REGION}" \
    --port="${HTTP_PORT}" \
    --set-env-vars="${ENV_VARS}" \
    --memory=512Mi \
    --cpu=1 \
    --min-instances=0 \
    --max-instances=10 \
    --concurrency=100 \
    --timeout=300 \
    ${ACCESS_FLAG} \
    --quiet; then
    SUCCEEDED+=("${NAME}")
    log_info "${NAME}-service deployed successfully."
  else
    FAILED+=("${NAME}")
    log_error "Failed to deploy ${NAME}-service."
  fi

  echo ""
done

# =============================================================================
# Summary
# =============================================================================
echo ""
echo "============================================="
echo "  Deployment Summary"
echo "============================================="
log_info "Succeeded: ${#SUCCEEDED[@]}/${#SERVICES[@]}"
for svc in "${SUCCEEDED[@]}"; do
  echo "  - ${svc}"
done

if [ ${#FAILED[@]} -gt 0 ]; then
  log_error "Failed: ${#FAILED[@]}/${#SERVICES[@]}"
  for svc in "${FAILED[@]}"; do
    echo "  - ${svc}"
  done
  exit 1
fi

echo ""
log_info "All services deployed successfully to ${REGION}!"
echo ""
echo "Service URLs:"
for service_def in "${SERVICES[@]}"; do
  IFS=':' read -r NAME _ _ _ _ <<< "${service_def}"
  URL=$(gcloud run services describe "${NAME}-service" \
    --project="${PROJECT}" \
    --region="${REGION}" \
    --format='value(status.url)' 2>/dev/null || echo "N/A")
  echo "  ${NAME}-service: ${URL}"
done
