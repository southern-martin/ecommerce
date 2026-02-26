# AI-Ready B2C Marketplace — Go Microservices + Clean Architecture

## Context

Build a future-proof B2C marketplace ecommerce platform with deep AI integration. The platform connects buyers and sellers (like Amazon/Shopee) with AI-powered features as a core competitive advantage. Backend uses **Go microservices with Clean Architecture** (20 services), communicating via **gRPC** internally and exposing **REST** to the frontend via an API Gateway. Features include full product attribute/variant system, promotions engine, shipping/logistics, returns/disputes, loyalty program, affiliate/referral system, tax engine, CMS, multi-currency, and i18n support.

---

## Part 1: Business Plan

### Revenue Model

| Stream | Description | Phase |
|---|---|---|
| **Commission** | 5-15% per transaction by category | MVP |
| **Subscription** | Seller tiers: Free (10 listings), Pro ($29/mo), Enterprise ($99/mo) | Phase 2 |
| **Featured Listings** | Sellers pay to boost visibility | Phase 2 |
| **Advertising** | Sponsored product placements, banners | Phase 3 |
| **AI Premium** | AI descriptions, dynamic pricing insights as paid seller tools | Phase 4 |

### User Types & Journeys

- **Buyer**: Browse/Search → View Product → Add to Cart → Checkout → Pay → Track → Receive → Review
- **Seller**: Register → Verify → Create Store → List Products → Receive Orders → Ship → Get Paid → Analytics
- **Admin**: Monitor → Manage Users → Approve Sellers → Handle Disputes → Reports → Configure AI

### AI-First Competitive Advantages
- Semantic search understands intent, not just keywords
- AI shopping assistant guides buying decisions conversationally
- Image-based search ("find something like this photo")
- Auto-generated SEO product descriptions for sellers
- Dynamic pricing suggestions based on market data
- Fraud detection on transactions and fake reviews

---

## Part 2: Tech Stack

| Layer | Technology |
|---|---|
| **Backend** | Go (Golang) — Gin framework |
| **Architecture** | Clean Architecture (Uncle Bob) per microservice |
| **Inter-service** | gRPC (protobuf) |
| **API Gateway** | Kong Gateway (DB-less, declarative YAML) |
| **ORM** | GORM |
| **Frontend** | Vite + React 18 + TypeScript + Tailwind CSS + Shadcn/ui (SPA, module-based) |
| **Routing** | React Router v6 (module-based lazy loading) |
| **Server State** | TanStack Query v5 (API caching, mutations) |
| **Client State** | Zustand (cart, auth, UI state) |
| **Database** | PostgreSQL 16 + pgvector |
| **Cache** | Redis 7 |
| **Search** | Elasticsearch 8 |
| **Message Broker** | NATS JetStream (async events between services) |
| **AI/ML** | Python FastAPI microservice + OpenAI/Claude APIs |
| **Storage** | S3-compatible (MinIO for dev) |
| **Auth** | JWT + OAuth2 (Google, Facebook, Apple) |
| **Payments** | Stripe + Stripe Connect (marketplace) |
| **Mobile** | Flutter 3.x + Dart (BLoC/Cubit, get_it+injectable, Dio, go_router, freezed) |
| **Mobile Apps** | Buyer App + Seller/Admin App (shared packages: core, api_client, ui_kit, shared_models) |
| **Containers** | Docker + docker-compose |
| **Orchestration** | Kubernetes (Phase 5) |
| **Observability** | OpenTelemetry + Jaeger, Prometheus + Grafana |
| **Logging** | zerolog (structured JSON) |

---

## Part 3: Microservice Decomposition

### Service Map & Ports

```
Service              HTTP    gRPC    Primary Store            Notes
─────────────────    ─────   ─────   ─────────────────────    ─────────────────────
kong                 8000    —       — (DB-less YAML)         Edge gateway, all traffic enters here
                     8001                                     Kong Admin API
auth                 8090    9090    PostgreSQL + Redis       REST for Kong, gRPC for inter-service
user                 8091    9091    PostgreSQL               REST for Kong, gRPC for inter-service
product              8081    9081    PostgreSQL               REST for Kong, gRPC for inter-service
cart                 8082    9082    Redis (+ PG backup)      REST for Kong, gRPC for inter-service
order                8083    9083    PostgreSQL               REST for Kong, gRPC for inter-service
payment              8084    9084    PostgreSQL               REST for Kong + Stripe webhooks
search               8085    9085    Elasticsearch            REST for Kong, gRPC for inter-service
review               8086    9086    PostgreSQL               REST for Kong, gRPC for inter-service
notification         8087    9087    PostgreSQL + Redis       REST for Kong + WebSocket
chat                 8088    9088    PostgreSQL + Redis       REST for Kong + WebSocket
media                8089    9089    PostgreSQL + S3          REST for Kong, gRPC for inter-service
ai                   8092    9092    PostgreSQL (pgvector)    REST for Kong, gRPC for inter-service
promotion            8093    9093    PostgreSQL + Redis       Coupons, vouchers, flash sales, bundles
return               8094    9094    PostgreSQL               Returns, refunds, disputes
shipping             8095    9095    PostgreSQL               Carrier integration, labels, tracking
loyalty              8096    9096    PostgreSQL + Redis       Points, cashback, membership tiers
affiliate            8097    9097    PostgreSQL               Referral tracking, affiliate commissions
tax                  8098    9098    PostgreSQL + Redis       Tax rules engine, jurisdiction config
cms                  8099    9099    PostgreSQL + S3          Banners, landing pages, content scheduling
ai-services (Python) 8000    —       —                        Internal only, called by AI service
web (Next.js)        3000    —       —                        Frontend SPA/SSR
```

### Architecture Flow
```
                    ┌──────────────┐
  Client (Browser)  │   Kong GW    │  :8000 (proxy) / :8001 (admin)
  ──────────────────▶  DB-less     │
                    │  YAML config │
                    └──────┬───────┘
                           │
              ┌────────────┼────────────────────┐
              │            │                    │
         REST routes    JWT plugin         Rate Limiting
              │       (validates token)    CORS, Logging
              │            │                    │
              ▼            ▼                    ▼
    ┌─────────────────────────────────────────────────┐
    │          Go Microservices (REST + gRPC)          │
    │                                                  │
    │  auth:8090  user:8091  product:8081  cart:8082   │
    │  order:8083  payment:8084  search:8085          │
    │  promotion:8093  return:8094  shipping:8095     │
    │  loyalty:8096  affiliate:8097  tax:8098  ...    │
    │                                                  │
    │  Inter-service calls use gRPC (:9xxx ports)     │
    └─────────────────────────────────────────────────┘
```

**Key change**: Each Go service now exposes **both** HTTP (Gin REST) and gRPC:
- **HTTP (REST)**: Kong routes external traffic to these endpoints
- **gRPC**: Used for inter-service communication only (faster, typed)

### Database-per-Service Strategy
Each service owns its own PostgreSQL database (e.g., `ecommerce_auth`, `ecommerce_products`, `ecommerce_orders`, `ecommerce_promotions`, `ecommerce_returns`, `ecommerce_shipping`, `ecommerce_loyalty`, `ecommerce_affiliates`, `ecommerce_tax`, `ecommerce_cms`). Services **never** access each other's databases directly — they communicate via gRPC calls or NATS events.

---

## Part 4: Project Structure

```
/opt/openAi/ecommerce/
├── go.work                               # Go workspace file
├── go.work.sum
├── Makefile                              # Root: proto gen, build all, test all, docker
├── docker-compose.yml
├── .env.example
├── .gitignore
│
├── proto/                                # Shared protobuf definitions
│   ├── auth/auth.proto
│   ├── user/user.proto
│   ├── product/product.proto
│   ├── cart/cart.proto
│   ├── order/order.proto
│   ├── payment/payment.proto
│   ├── search/search.proto
│   ├── review/review.proto
│   ├── notification/notification.proto
│   ├── chat/chat.proto
│   ├── media/media.proto
│   ├── ai/ai.proto
│   ├── promotion/promotion.proto
│   ├── return/return.proto
│   ├── shipping/shipping.proto
│   ├── loyalty/loyalty.proto
│   ├── affiliate/affiliate.proto
│   ├── tax/tax.proto
│   └── cms/cms.proto
│
├── pkg/                                  # Shared Go packages (go module)
│   ├── go.mod
│   ├── logger/logger.go                  # zerolog structured logging
│   ├── errors/errors.go                  # Custom error types + gRPC code mapping
│   ├── auth/jwt.go                       # JWT generation, validation, claims
│   ├── middleware/                        # Shared Gin middleware
│   │   ├── auth.go                       # JWT auth middleware
│   │   ├── cors.go
│   │   ├── ratelimit.go
│   │   └── correlation.go               # Correlation ID propagation
│   ├── pagination/pagination.go          # Cursor/offset pagination helpers
│   ├── validator/validator.go            # Input validation helpers
│   ├── events/                           # NATS event schemas
│   │   ├── publisher.go
│   │   ├── subscriber.go
│   │   └── subjects.go                   # Event subject constants
│   ├── server/graceful.go                # Graceful shutdown (gRPC + HTTP)
│   ├── circuitbreaker/circuitbreaker.go  # gRPC circuit breaker (gobreaker)
│   ├── tracing/tracing.go               # OpenTelemetry init
│   ├── unitofwork/uow.go                # GORM transaction wrapper
│   ├── i18n/i18n.go                     # Multi-language message catalogs, locale detection
│   ├── money/money.go                   # Currency-safe amounts, multi-currency formatting
│   └── tax/client.go                    # gRPC client for tax-service
│
├── kong/                                 # Kong Gateway configuration
│   ├── kong.yml                          # Declarative config (services, routes, plugins)
│   ├── kong.dev.yml                      # Dev overrides
│   ├── kong.prod.yml                     # Production overrides
│   └── plugins/                          # Custom Kong plugins (if needed)
│       └── jwt-claims-to-header/         # Extract JWT claims → X-User-ID header
│
├── services/                             # All Go microservices (20 services)
│   ├── auth/
│   ├── user/
│   ├── product/                          # Reference implementation below
│   ├── cart/
│   ├── order/
│   ├── payment/
│   ├── search/
│   ├── review/
│   ├── notification/
│   ├── chat/
│   ├── media/
│   ├── ai/
│   ├── promotion/                        # Coupons, vouchers, flash sales, bundles
│   ├── return/                           # Returns, refunds, disputes
│   ├── shipping/                         # Carrier integration, labels, tracking
│   ├── loyalty/                          # Points, cashback, membership tiers
│   ├── affiliate/                        # Referral tracking, affiliate commissions
│   ├── tax/                              # Tax rules engine, jurisdiction config
│   └── cms/                              # Banners, landing pages, content scheduling
│
├── apps/
│   ├── web/                              # Vite + React SPA (module-based, see Part 17)
│   ├── mobile/
│   │   ├── buyer_app/                    # Flutter Buyer App (see Part 18)
│   │   ├── seller_app/                   # Flutter Seller/Admin App (see Part 18)
│   │   └── packages/                     # Shared Flutter packages
│   │       ├── core/                     # Shared domain entities, use cases, DI
│   │       ├── api_client/               # Dio HTTP client, interceptors
│   │       ├── ui_kit/                   # Shared widgets, theme, design tokens
│   │       └── shared_models/            # DTOs, enums shared between apps
│   └── ai-services/                      # Python FastAPI ML service
│
├── infra/
│   ├── docker/
│   │   └── postgres/
│   │       └── create-multiple-dbs.sh    # Creates per-service databases
│   ├── k8s/                              # Kubernetes manifests (Phase 5)
│   │   ├── namespace.yaml
│   │   ├── product/deployment.yaml
│   │   ├── gateway/ingress.yaml
│   │   └── ...
│   └── scripts/
│       └── setup-dev.sh
│
└── .github/
    └── workflows/
        └── ci.yml                        # Matrix CI: lint, test, build per service
```

---

## Part 5: Clean Architecture per Service (Reference: Product Service)

```
services/product/
├── cmd/
│   └── main.go                           # Entry point, DI wiring, startup
├── internal/
│   ├── domain/                           # INNER RING — zero dependencies
│   │   ├── entity/
│   │   │   ├── product.go                # Product struct, NewProduct(), business methods
│   │   │   ├── category.go              # Category + CategoryAttribute definitions
│   │   │   ├── attribute.go             # AttributeDefinition, AttributeType, ProductAttributeValue
│   │   │   ├── option.go                # ProductOption, ProductOptionValue (variant axes)
│   │   │   └── variant.go              # Variant with SKU, price override, stock, images
│   │   ├── repository/                   # Interface definitions (ports)
│   │   │   ├── product_repository.go     # ProductRepository interface
│   │   │   ├── category_repository.go
│   │   │   ├── attribute_repository.go  # AttributeDefinition + ProductAttributeValue CRUD
│   │   │   ├── option_repository.go     # ProductOption + OptionValue CRUD
│   │   │   └── variant_repository.go    # Variant CRUD, stock operations
│   │   ├── valueobject/
│   │   │   └── money.go
│   │   └── event/
│   │       └── publisher.go              # EventPublisher interface
│   │
│   ├── usecase/                          # MIDDLE RING — orchestration
│   │   ├── product_usecase.go            # CreateProduct, GetProduct, ListProducts, etc.
│   │   ├── category_usecase.go           # Category CRUD + attribute assignment
│   │   ├── variant_usecase.go            # Option management, variant generation, stock ops
│   │   └── interfaces.go                 # Input/Output DTOs
│   │
│   ├── adapter/                          # OUTER RING — implementations
│   │   ├── handler/
│   │   │   ├── http/                     # Gin HTTP handlers (health, debug)
│   │   │   │   └── product_handler.go
│   │   │   └── grpc/                     # gRPC server implementations
│   │   │       └── product_grpc.go
│   │   ├── repository/
│   │   │   └── postgres/                 # GORM repository implementations
│   │   │       ├── product_repository.go
│   │   │       ├── category_repository.go
│   │   │       ├── attribute_repository.go
│   │   │       ├── option_repository.go
│   │   │       ├── variant_repository.go
│   │   │       └── models.go             # GORM models (DB-specific, map to/from entities)
│   │   ├── presenter/
│   │   │   └── product_presenter.go      # Entity → proto response mapping
│   │   └── middleware/
│   │       └── auth.go
│   │
│   └── infrastructure/                   # OUTERMOST — frameworks & drivers
│       ├── config/config.go              # Viper config loading
│       ├── database/database.go          # GORM connection setup
│       ├── cache/redis.go                # Redis client
│       ├── messaging/nats.go             # NATS publisher/subscriber
│       └── grpc/clients.go              # gRPC client connections to other services
│
├── migrations/                           # golang-migrate SQL files
│   ├── 000001_create_products.up.sql
│   ├── 000001_create_products.down.sql
│   ├── 000002_create_attributes.up.sql
│   ├── 000002_create_attributes.down.sql
│   ├── 000003_create_options_variants.up.sql
│   └── 000003_create_options_variants.down.sql
├── Dockerfile
├── Makefile
└── go.mod
```

### Dependency Flow (Clean Architecture Rule)

```
    ┌─────────────────────────────────────────┐
    │         INFRASTRUCTURE                   │
    │  config, database, cache, messaging      │
    │         ↓ injected into ↓                │
    ├─────────────────────────────────────────┤
    │         ADAPTER (outer ring)             │
    │  handler/grpc, handler/http              │
    │  repository/postgres (IMPLEMENTS domain  │
    │  repository interfaces)                  │
    │         ↓ calls ↓                        │
    ├─────────────────────────────────────────┤
    │         USE CASE (middle ring)           │
    │  Orchestrates entities, calls repository │
    │  interfaces, publishes events            │
    │  NO framework code here                  │
    │         ↓ uses ↓                         │
    ├─────────────────────────────────────────┤
    │         DOMAIN (inner ring)              │
    │  entity/ — pure Go structs              │
    │  repository/ — interface definitions     │
    │  DEPENDS ON: NOTHING                     │
    └─────────────────────────────────────────┘

    Source code dependencies ALWAYS point inward.
```

---

## Part 6: Key Code Patterns

### Domain Entity (zero dependencies)
```go
// internal/domain/entity/product.go
type ProductStatus string
const (
    ProductStatusDraft    ProductStatus = "draft"
    ProductStatusActive   ProductStatus = "active"
    ProductStatusInactive ProductStatus = "inactive"
    ProductStatusArchived ProductStatus = "archived"
)

type Product struct {
    ID              string
    SellerID        string
    CategoryID      string
    Name            string
    Slug            string
    Description     string
    BasePriceCents  int64             // base price (variants can override)
    Currency        string            // ISO 4217 (USD, AUD, etc.)
    Status          ProductStatus
    HasVariants     bool              // true if product uses options/variants
    Tags            []string
    ImageURLs       []string
    Embedding       []float64         // pgvector AI embedding
    AttributeValues []ProductAttributeValue  // filled-in category attributes
    Options         []ProductOption          // variant axes (Color, Size, etc.)
    Variants        []Variant                // generated from option combinations
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

func NewProduct(sellerID, name, desc string, price int64, currency, catID string) *Product {
    return &Product{
        ID: uuid.New().String(), SellerID: sellerID,
        Name: name, Slug: slugify(name), Description: desc,
        BasePriceCents: price, Currency: currency, CategoryID: catID,
        HasVariants: false,
        Status: ProductStatusDraft, CreatedAt: time.Now(), UpdatedAt: time.Now(),
    }
}

// internal/domain/entity/attribute.go

// AttributeType defines what kind of values an attribute accepts.
type AttributeType string
const (
    AttributeTypeText        AttributeType = "text"         // free text
    AttributeTypeNumber      AttributeType = "number"       // numeric value
    AttributeTypeSelect      AttributeType = "select"       // single choice from predefined options
    AttributeTypeMultiSelect AttributeType = "multi_select" // multiple choices
    AttributeTypeColor       AttributeType = "color"        // color swatch (hex value)
    AttributeTypeBool        AttributeType = "bool"         // yes/no
)

// AttributeDefinition is an admin-defined attribute template assigned to categories.
// Example: Category "Clothing" has attributes [Brand(text,required), Material(select), Gender(select)]
type AttributeDefinition struct {
    ID            string
    Name          string          // "Brand", "Material", "Screen Size"
    Slug          string          // "brand", "material", "screen-size"
    Type          AttributeType
    Required      bool            // seller must fill this in
    Filterable    bool            // appears in search filter panel
    Options       []string        // predefined choices for select/multi_select types
    Unit          string          // optional unit label: "cm", "kg", "inches"
    SortOrder     int
    CreatedAt     time.Time
}

// CategoryAttribute links an AttributeDefinition to a Category.
type CategoryAttribute struct {
    CategoryID    string
    AttributeID   string
    SortOrder     int             // display order within this category
}

// ProductAttributeValue stores the seller's filled-in value for a category attribute.
type ProductAttributeValue struct {
    ID            string
    ProductID     string
    AttributeID   string
    AttributeName string          // denormalized for read performance
    Value         string          // single text/number/select value
    Values        []string        // for multi_select type
}

// internal/domain/entity/option.go

// ProductOption represents a variant axis defined by the seller (max 3 per product).
// Example: Product "T-Shirt" has options: [Color, Size]
type ProductOption struct {
    ID        string
    ProductID string
    Name      string              // "Color", "Size", "Material"
    SortOrder int                 // display order
    Values    []ProductOptionValue
}

// ProductOptionValue is one value within an option.
// Example: Option "Color" has values: [Red, Blue, Black]
type ProductOptionValue struct {
    ID         string
    OptionID   string
    Value      string             // "Red", "XL", "Cotton"
    ColorHex   string             // optional hex code for color swatches (#FF0000)
    SortOrder  int
}

// internal/domain/entity/variant.go

// Variant represents a specific purchasable combination of option values.
// Example: Color=Red + Size=XL → SKU "TSH-RED-XL", price 2500, stock 10
type Variant struct {
    ID             string
    ProductID      string
    SKU            string         // unique stock-keeping unit
    Name           string         // auto-generated: "Red / XL"
    PriceCents     int64          // 0 means use product base price
    CompareAtCents int64          // original/compare-at price (for showing discounts)
    CostCents      int64          // cost of goods (for profit calculation)
    Stock          int
    LowStockAlert  int            // notify seller when stock drops below this
    Weight         int            // grams (for shipping calculation)
    IsDefault      bool           // default selected variant on product page
    IsActive       bool           // can be deactivated without deleting
    ImageURLs      []string       // variant-specific images (e.g., different color photos)
    Barcode        string         // UPC/EAN/ISBN
    OptionValues   []VariantOptionValue  // which option values make up this variant
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

// VariantOptionValue links a variant to the specific option values it represents.
// Example: Variant "TSH-RED-XL" → [{OptionName:"Color", Value:"Red"}, {OptionName:"Size", Value:"XL"}]
type VariantOptionValue struct {
    VariantID     string
    OptionID      string
    OptionValueID string
    OptionName    string          // denormalized: "Color"
    Value         string          // denormalized: "Red"
}
```

### Repository Interface (domain layer)
```go
// internal/domain/repository/product_repository.go
type ProductFilter struct {
    SellerID, CategoryID, Status, Query, SortBy string
    MinPrice, MaxPrice                          int64
    Attributes                                  map[string][]string // attribute slug → values for faceted filtering
    Page, PageSize                              int
}

type ProductRepository interface {
    Create(ctx context.Context, p *entity.Product) error
    GetByID(ctx context.Context, id string) (*entity.Product, error)
    GetBySlug(ctx context.Context, slug string) (*entity.Product, error)
    List(ctx context.Context, filter ProductFilter) ([]*entity.Product, int64, error)
    Update(ctx context.Context, p *entity.Product) error
    Delete(ctx context.Context, id string) error
}

// internal/domain/repository/attribute_repository.go
type AttributeRepository interface {
    // Attribute definitions (admin)
    CreateDefinition(ctx context.Context, attr *entity.AttributeDefinition) error
    GetDefinitionByID(ctx context.Context, id string) (*entity.AttributeDefinition, error)
    ListDefinitions(ctx context.Context) ([]*entity.AttributeDefinition, error)
    UpdateDefinition(ctx context.Context, attr *entity.AttributeDefinition) error
    DeleteDefinition(ctx context.Context, id string) error

    // Category ↔ attribute assignment (admin)
    AssignToCategory(ctx context.Context, categoryID, attributeID string, sortOrder int) error
    RemoveFromCategory(ctx context.Context, categoryID, attributeID string) error
    ListByCategory(ctx context.Context, categoryID string) ([]*entity.AttributeDefinition, error)

    // Product attribute values (seller)
    SetProductValues(ctx context.Context, productID string, values []entity.ProductAttributeValue) error
    GetProductValues(ctx context.Context, productID string) ([]entity.ProductAttributeValue, error)
}

// internal/domain/repository/option_repository.go
type OptionRepository interface {
    CreateOption(ctx context.Context, option *entity.ProductOption) error
    UpdateOption(ctx context.Context, option *entity.ProductOption) error
    DeleteOption(ctx context.Context, optionID string) error
    ListByProduct(ctx context.Context, productID string) ([]entity.ProductOption, error)

    CreateOptionValue(ctx context.Context, value *entity.ProductOptionValue) error
    UpdateOptionValue(ctx context.Context, value *entity.ProductOptionValue) error
    DeleteOptionValue(ctx context.Context, valueID string) error
}

// internal/domain/repository/variant_repository.go
type VariantRepository interface {
    Create(ctx context.Context, v *entity.Variant) error
    GetByID(ctx context.Context, id string) (*entity.Variant, error)
    GetBySKU(ctx context.Context, sku string) (*entity.Variant, error)
    ListByProduct(ctx context.Context, productID string) ([]entity.Variant, error)
    Update(ctx context.Context, v *entity.Variant) error
    Delete(ctx context.Context, id string) error
    BulkCreate(ctx context.Context, variants []entity.Variant) error
    UpdateStock(ctx context.Context, variantID string, delta int) error  // atomic stock adjustment
    SetOptionValues(ctx context.Context, variantID string, values []entity.VariantOptionValue) error
}
```

### Use Case (orchestration, no framework deps)
```go
// internal/usecase/product_usecase.go
type ProductUseCase struct {
    productRepo   repository.ProductRepository
    categoryRepo  repository.CategoryRepository
    attributeRepo repository.AttributeRepository
    publisher     event.EventPublisher
}

func NewProductUseCase(
    pr repository.ProductRepository,
    cr repository.CategoryRepository,
    ar repository.AttributeRepository,
    pub event.EventPublisher,
) *ProductUseCase {
    return &ProductUseCase{productRepo: pr, categoryRepo: cr, attributeRepo: ar, publisher: pub}
}

func (uc *ProductUseCase) CreateProduct(ctx context.Context, input CreateProductInput) (*entity.Product, error) {
    if _, err := uc.categoryRepo.GetByID(ctx, input.CategoryID); err != nil {
        return nil, fmt.Errorf("category not found: %w", err)
    }

    // Validate required category attributes are provided
    if err := uc.validateAttributes(ctx, input.CategoryID, input.AttributeValues); err != nil {
        return nil, err
    }

    product := entity.NewProduct(input.SellerID, input.Name, input.Description, input.PriceCents, input.Currency, input.CategoryID)
    product.AttributeValues = input.AttributeValues

    if err := uc.productRepo.Create(ctx, product); err != nil {
        return nil, err
    }

    // Save attribute values
    if len(input.AttributeValues) > 0 {
        if err := uc.attributeRepo.SetProductValues(ctx, product.ID, input.AttributeValues); err != nil {
            return nil, err
        }
    }

    _ = uc.publisher.Publish("product.created", product)
    return product, nil
}

// validateAttributes checks that all required category attributes have values.
func (uc *ProductUseCase) validateAttributes(ctx context.Context, categoryID string, values []entity.ProductAttributeValue) error {
    required, err := uc.attributeRepo.ListByCategory(ctx, categoryID)
    if err != nil {
        return err
    }
    provided := make(map[string]bool)
    for _, v := range values {
        provided[v.AttributeID] = true
    }
    for _, attr := range required {
        if attr.Required && !provided[attr.ID] {
            return fmt.Errorf("required attribute %q is missing", attr.Name)
        }
    }
    return nil
}

// internal/usecase/variant_usecase.go

type VariantUseCase struct {
    productRepo  repository.ProductRepository
    optionRepo   repository.OptionRepository
    variantRepo  repository.VariantRepository
    publisher    event.EventPublisher
}

// AddOption adds a variant axis to a product (e.g., "Color").
// Max 3 options per product. After adding options, call GenerateVariants.
func (uc *VariantUseCase) AddOption(ctx context.Context, productID string, name string, values []string) (*entity.ProductOption, error) {
    existing, _ := uc.optionRepo.ListByProduct(ctx, productID)
    if len(existing) >= 3 {
        return nil, errors.New("maximum 3 options per product")
    }
    option := &entity.ProductOption{
        ID: uuid.New().String(), ProductID: productID,
        Name: name, SortOrder: len(existing) + 1,
    }
    if err := uc.optionRepo.CreateOption(ctx, option); err != nil {
        return nil, err
    }
    for i, v := range values {
        val := &entity.ProductOptionValue{
            ID: uuid.New().String(), OptionID: option.ID,
            Value: v, SortOrder: i + 1,
        }
        if err := uc.optionRepo.CreateOptionValue(ctx, val); err != nil {
            return nil, err
        }
        option.Values = append(option.Values, *val)
    }
    return option, nil
}

// GenerateVariants creates all variant combinations from product options.
// Example: Color[Red,Blue] × Size[S,M] → 4 variants.
func (uc *VariantUseCase) GenerateVariants(ctx context.Context, productID string) ([]entity.Variant, error) {
    product, err := uc.productRepo.GetByID(ctx, productID)
    if err != nil {
        return nil, err
    }
    options, err := uc.optionRepo.ListByProduct(ctx, productID)
    if err != nil {
        return nil, err
    }
    if len(options) == 0 {
        return nil, errors.New("no options defined")
    }

    combinations := generateCombinations(options)
    var variants []entity.Variant
    for i, combo := range combinations {
        name := buildVariantName(combo) // "Red / S"
        v := entity.Variant{
            ID: uuid.New().String(), ProductID: productID,
            SKU:  fmt.Sprintf("%s-%d", product.Slug, i+1), // seller can override later
            Name: name, PriceCents: 0, // 0 = inherit base price
            Stock: 0, IsDefault: i == 0, IsActive: true,
            OptionValues: combo,
        }
        variants = append(variants, v)
    }
    if err := uc.variantRepo.BulkCreate(ctx, variants); err != nil {
        return nil, err
    }

    // Mark product as having variants
    product.HasVariants = true
    _ = uc.productRepo.Update(ctx, product)
    _ = uc.publisher.Publish("product.updated", product)
    return variants, nil
}

// UpdateVariantStock atomically adjusts variant stock (positive=restock, negative=purchase).
func (uc *VariantUseCase) UpdateVariantStock(ctx context.Context, variantID string, delta int) error {
    return uc.variantRepo.UpdateStock(ctx, variantID, delta)
}
```

### GORM Repository Implementation (adapter layer)
```go
// internal/adapter/repository/postgres/product_repository.go
type productRepository struct{ db *gorm.DB }

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
    return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, p *entity.Product) error {
    model := toProductModel(p) // entity → GORM model
    return r.db.WithContext(ctx).Create(&model).Error
}

func (r *productRepository) GetByID(ctx context.Context, id string) (*entity.Product, error) {
    var model ProductModel
    if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, entity.ErrProductNotFound
        }
        return nil, err
    }
    return model.ToEntity(), nil // GORM model → entity
}
```

### gRPC Handler (adapter layer)
```go
// internal/adapter/handler/grpc/product_grpc.go
type ProductGRPCServer struct {
    productpb.UnimplementedProductServiceServer
    useCase *usecase.ProductUseCase
}

func (s *ProductGRPCServer) CreateProduct(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.ProductResponse, error) {
    product, err := s.useCase.CreateProduct(ctx, usecase.CreateProductInput{
        SellerID: req.SellerId, Name: req.Name,
        Description: req.Description, PriceCents: req.Price.AmountCents,
        CategoryID: req.CategoryId, Tags: req.Tags,
    })
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
    }
    return presenter.ToProductResponse(product), nil
}
```

### Entry Point (DI wiring)
```go
// cmd/main.go
func main() {
    cfg := config.Load()
    log := logger.New(cfg.LogLevel)
    db := database.Connect(cfg.DB)
    natsConn := messaging.ConnectNATS(cfg.NatsURL)

    // Repository implementations
    productRepo := postgres.NewProductRepository(db)
    categoryRepo := postgres.NewCategoryRepository(db)
    attributeRepo := postgres.NewAttributeRepository(db)
    optionRepo := postgres.NewOptionRepository(db)
    variantRepo := postgres.NewVariantRepository(db)
    publisher := messaging.NewNATSPublisher(natsConn)

    // Use cases (injecting interfaces)
    productUC := usecase.NewProductUseCase(productRepo, categoryRepo, attributeRepo, publisher)
    variantUC := usecase.NewVariantUseCase(productRepo, optionRepo, variantRepo, publisher)

    // gRPC server
    grpcServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
    productpb.RegisterProductServiceServer(grpcServer, grpchandler.NewProductGRPCServer(productUC, variantUC))

    // HTTP server (health + debug)
    router := gin.New()
    router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
    httpServer := &http.Server{Handler: router}

    // Graceful startup + shutdown
    srv := server.New(log, grpcServer, httpServer)
    srv.Run(cfg.GRPCPort, cfg.HTTPPort)
}
```

---

## Part 7: Protobuf Definitions (Key Services)

### Product Service Proto
```protobuf
syntax = "proto3";
package product;
option go_package = "proto/product";

service ProductService {
    rpc CreateProduct(CreateProductRequest) returns (ProductResponse);
    rpc GetProduct(GetProductRequest) returns (ProductResponse);
    rpc ListProducts(ListProductsRequest) returns (ListProductsResponse);
    rpc UpdateProduct(UpdateProductRequest) returns (ProductResponse);
    rpc DeleteProduct(DeleteProductRequest) returns (Empty);

    // Attribute definitions (admin)
    rpc CreateAttributeDefinition(CreateAttributeDefinitionRequest) returns (AttributeDefinitionResponse);
    rpc ListAttributeDefinitions(ListAttributeDefinitionsRequest) returns (ListAttributeDefinitionsResponse);
    rpc AssignAttributeToCategory(AssignAttributeToCategoryRequest) returns (Empty);
    rpc ListCategoryAttributes(ListCategoryAttributesRequest) returns (ListAttributeDefinitionsResponse);

    // Product options & variants (seller)
    rpc AddProductOption(AddProductOptionRequest) returns (ProductOptionResponse);
    rpc RemoveProductOption(RemoveProductOptionRequest) returns (Empty);
    rpc GenerateVariants(GenerateVariantsRequest) returns (ListVariantsResponse);
    rpc UpdateVariant(UpdateVariantRequest) returns (VariantResponse);
    rpc UpdateVariantStock(UpdateVariantStockRequest) returns (Empty);
    rpc ListVariants(ListVariantsRequest) returns (ListVariantsResponse);

    // Categories
    rpc GetCategories(Empty) returns (CategoriesResponse);
}

message Product {
    string id = 1;
    string seller_id = 2;
    string category_id = 3;
    string name = 4;
    string slug = 5;
    string description = 6;
    Price base_price = 7;
    string status = 8;
    bool has_variants = 9;
    repeated string tags = 10;
    repeated string image_urls = 11;
    repeated ProductAttributeValue attribute_values = 12;
    repeated ProductOption options = 13;
    repeated Variant variants = 14;
    double rating_avg = 15;
    int32 rating_count = 16;
    string created_at = 17;
    string updated_at = 18;
}

message Price {
    int64 amount_cents = 1;
    string currency = 2;
}

// --- Attribute system ---

message AttributeDefinition {
    string id = 1;
    string name = 2;
    string slug = 3;
    string type = 4;             // text, number, select, multi_select, color, bool
    bool required = 5;
    bool filterable = 6;
    repeated string options = 7; // predefined choices for select types
    string unit = 8;             // optional: "cm", "kg"
    int32 sort_order = 9;
}

message ProductAttributeValue {
    string attribute_id = 1;
    string attribute_name = 2;
    string value = 3;            // single value
    repeated string values = 4;  // for multi_select
}

// --- Option / Variant system ---

message ProductOption {
    string id = 1;
    string name = 2;                          // "Color", "Size"
    int32 sort_order = 3;
    repeated ProductOptionValue values = 4;
}

message ProductOptionValue {
    string id = 1;
    string value = 2;            // "Red", "XL"
    string color_hex = 3;        // optional: "#FF0000"
    int32 sort_order = 4;
}

message Variant {
    string id = 1;
    string product_id = 2;
    string sku = 3;
    string name = 4;                          // "Red / XL"
    int64 price_cents = 5;                    // 0 = use base price
    int64 compare_at_cents = 6;               // original price for discounts
    int64 cost_cents = 7;                     // cost of goods
    int32 stock = 8;
    int32 low_stock_alert = 9;
    int32 weight_grams = 10;
    bool is_default = 11;
    bool is_active = 12;
    repeated string image_urls = 13;
    string barcode = 14;                      // UPC/EAN
    repeated VariantOptionValue option_values = 15;
}

message VariantOptionValue {
    string option_id = 1;
    string option_value_id = 2;
    string option_name = 3;       // "Color"
    string value = 4;             // "Red"
}
```

### Order Service Proto
```protobuf
service OrderService {
    rpc CreateOrder(CreateOrderRequest) returns (OrderResponse);
    rpc GetOrder(GetOrderRequest) returns (OrderResponse);
    rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);
    rpc UpdateOrderStatus(UpdateOrderStatusRequest) returns (OrderResponse);
    rpc CancelOrder(CancelOrderRequest) returns (OrderResponse);
}

message Order {
    string id = 1; string order_number = 2; string buyer_id = 3;
    string status = 4; int64 subtotal_cents = 5; int64 shipping_cents = 6;
    int64 tax_cents = 7; int64 discount_cents = 8; int64 total_cents = 9;
    string currency = 10; Address shipping_address = 11;
    repeated OrderItem items = 12; double fraud_score = 13;
    string created_at = 14; string updated_at = 15;
}

message OrderItem {
    string id = 1; string product_id = 2; string variant_id = 3;
    string product_name = 4; int32 quantity = 5; int64 unit_price_cents = 6;
    int64 total_cents = 7; string seller_id = 8;
}
```

### Auth Service Proto
```protobuf
service AuthService {
    rpc Register(RegisterRequest) returns (AuthResponse);
    rpc Login(LoginRequest) returns (AuthResponse);
    rpc RefreshToken(RefreshTokenRequest) returns (AuthResponse);
    rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
    rpc Logout(LogoutRequest) returns (Empty);
    rpc ForgotPassword(ForgotPasswordRequest) returns (Empty);
    rpc ResetPassword(ResetPasswordRequest) returns (Empty);
}

message AuthResponse {
    string access_token = 1; string refresh_token = 2;
    int64 expires_in = 3; UserInfo user = 4;
}
```

### Promotion Service Proto
```protobuf
service PromotionService {
    rpc CreateCoupon(CreateCouponRequest) returns (CouponResponse);
    rpc GetCoupon(GetCouponRequest) returns (CouponResponse);
    rpc ListCoupons(ListCouponsRequest) returns (ListCouponsResponse);
    rpc ValidateCoupon(ValidateCouponRequest) returns (ValidateCouponResponse);  // called by cart/order
    rpc RedeemCoupon(RedeemCouponRequest) returns (Empty);
    rpc CreateFlashSale(CreateFlashSaleRequest) returns (FlashSaleResponse);
    rpc ListFlashSales(ListFlashSalesRequest) returns (ListFlashSalesResponse);
    rpc CreateBundle(CreateBundleRequest) returns (BundleResponse);
    rpc ListBundles(ListBundlesRequest) returns (ListBundlesResponse);
}

message Coupon {
    string id = 1; string code = 2; string type = 3;           // percentage, fixed_amount, free_shipping
    int64 discount_value = 4;                                   // cents or percentage * 100
    int64 min_order_cents = 5; int64 max_discount_cents = 6;
    int32 usage_limit = 7; int32 usage_count = 8;
    int32 per_user_limit = 9;
    string scope = 10;                                          // all, category, product, seller
    repeated string scope_ids = 11;                             // category/product/seller IDs
    string starts_at = 12; string expires_at = 13;
    bool is_active = 14; string created_by = 15;                // seller_id or "platform"
}

message FlashSale {
    string id = 1; string name = 2;
    string starts_at = 3; string ends_at = 4;
    repeated FlashSaleItem items = 5;
    bool is_active = 6;
}

message FlashSaleItem {
    string product_id = 1; string variant_id = 2;
    int64 sale_price_cents = 3; int32 quantity_limit = 4; int32 sold_count = 5;
}

message Bundle {
    string id = 1; string name = 2; string seller_id = 3;
    repeated string product_ids = 4;
    int64 bundle_price_cents = 5;                               // price for buying all together
    int64 savings_cents = 6;                                    // how much buyer saves
    bool is_active = 7;
}
```

### Return Service Proto
```protobuf
service ReturnService {
    rpc CreateReturn(CreateReturnRequest) returns (ReturnResponse);
    rpc GetReturn(GetReturnRequest) returns (ReturnResponse);
    rpc ListReturns(ListReturnsRequest) returns (ListReturnsResponse);
    rpc ApproveReturn(ApproveReturnRequest) returns (ReturnResponse);
    rpc RejectReturn(RejectReturnRequest) returns (ReturnResponse);
    rpc ProcessRefund(ProcessRefundRequest) returns (RefundResponse);
    rpc CreateDispute(CreateDisputeRequest) returns (DisputeResponse);
    rpc ResolveDispute(ResolveDisputeRequest) returns (DisputeResponse);
    rpc ListDisputes(ListDisputesRequest) returns (ListDisputesResponse);
}

message Return {
    string id = 1; string order_id = 2; string buyer_id = 3; string seller_id = 4;
    string status = 5;                                          // requested, approved, rejected, shipped_back, received, refunded
    string reason = 6;                                          // defective, wrong_item, not_as_described, changed_mind, other
    string description = 7;
    repeated string image_urls = 8;                             // proof photos
    repeated ReturnItem items = 9;
    int64 refund_amount_cents = 10;
    string refund_method = 11;                                  // original_payment, wallet_credit
    string return_tracking_number = 12;
    string created_at = 13; string updated_at = 14;
}

message ReturnItem {
    string order_item_id = 1; string product_id = 2; string variant_id = 3;
    int32 quantity = 4; string reason = 5;
}

message Dispute {
    string id = 1; string order_id = 2; string return_id = 3;
    string buyer_id = 4; string seller_id = 5;
    string status = 6;                                          // open, under_review, resolved_buyer, resolved_seller, escalated
    string type = 7;                                            // item_not_received, not_as_described, unauthorized_charge
    string description = 8;
    repeated DisputeMessage messages = 9;
    string resolution = 10; string resolved_by = 11;
    string created_at = 12; string resolved_at = 13;
}

message DisputeMessage {
    string id = 1; string sender_id = 2; string sender_role = 3;   // buyer, seller, admin
    string message = 4; repeated string attachments = 5; string created_at = 6;
}
```

### Shipping Service Proto
```protobuf
service ShippingService {
    rpc GetShippingRates(GetShippingRatesRequest) returns (ShippingRatesResponse);    // rate shopping
    rpc CreateShipment(CreateShipmentRequest) returns (ShipmentResponse);
    rpc GetShipment(GetShipmentRequest) returns (ShipmentResponse);
    rpc GenerateLabel(GenerateLabelRequest) returns (LabelResponse);
    rpc GetTracking(GetTrackingRequest) returns (TrackingResponse);
    rpc ListCarriers(ListCarriersRequest) returns (ListCarriersResponse);
    rpc ConfigureCarrier(ConfigureCarrierRequest) returns (CarrierResponse);          // seller carrier setup
}

message ShippingRate {
    string carrier_code = 1; string service_name = 2;
    int64 rate_cents = 3; string currency = 4;
    int32 estimated_days_min = 5; int32 estimated_days_max = 6;
}

message Shipment {
    string id = 1; string order_id = 2; string seller_id = 3;
    string carrier_code = 4; string service_code = 5;
    string tracking_number = 6; string label_url = 7;
    string status = 8;                                          // pending, label_created, picked_up, in_transit, delivered, exception
    Address origin = 9; Address destination = 10;
    repeated ShipmentItem items = 11;
    int32 weight_grams = 12;
    string created_at = 13; string updated_at = 14;
}

message TrackingEvent {
    string status = 1; string description = 2;
    string location = 3; string timestamp = 4;
}

message Carrier {
    string code = 1; string name = 2;                           // fedex, ups, dhl, auspost, etc.
    bool is_active = 3;
    repeated string supported_countries = 4;
}
```

### Loyalty Service Proto
```protobuf
service LoyaltyService {
    rpc GetMembership(GetMembershipRequest) returns (MembershipResponse);
    rpc GetPointsBalance(GetPointsBalanceRequest) returns (PointsBalanceResponse);
    rpc EarnPoints(EarnPointsRequest) returns (PointsTransactionResponse);            // called after order completion
    rpc RedeemPoints(RedeemPointsRequest) returns (PointsTransactionResponse);        // called during checkout
    rpc ListTransactions(ListTransactionsRequest) returns (ListTransactionsResponse);
    rpc GetTiers(Empty) returns (TiersResponse);
}

message Membership {
    string user_id = 1; string tier = 2;                        // bronze, silver, gold, platinum
    int64 points_balance = 3; int64 lifetime_points = 4;
    string tier_expires_at = 5; string joined_at = 6;
}

message PointsTransaction {
    string id = 1; string user_id = 2; string type = 3;        // earn, redeem, expire, adjust
    int64 points = 4; string source = 5;                        // order, review, referral, promotion
    string reference_id = 6;                                    // orderId, reviewId, etc.
    string description = 7; string created_at = 8;
}

message Tier {
    string name = 1; int64 min_points = 2;
    double cashback_rate = 3;                                   // e.g., 0.02 = 2%
    double points_multiplier = 4;                               // e.g., 1.5x for gold
    bool free_shipping = 5;
    int32 priority_support_hours = 6;                           // response SLA
}
```

### Affiliate Service Proto
```protobuf
service AffiliateService {
    rpc CreateAffiliateLink(CreateAffiliateLinkRequest) returns (AffiliateLinkResponse);
    rpc GetAffiliateStats(GetAffiliateStatsRequest) returns (AffiliateStatsResponse);
    rpc ListReferrals(ListReferralsRequest) returns (ListReferralsResponse);
    rpc TrackClick(TrackClickRequest) returns (Empty);                                 // called on link click
    rpc TrackConversion(TrackConversionRequest) returns (Empty);                       // called on order completion
    rpc RequestPayout(RequestPayoutRequest) returns (PayoutResponse);
    rpc GetProgram(Empty) returns (AffiliateProgramResponse);
}

message AffiliateLink {
    string id = 1; string user_id = 2;
    string code = 3;                                            // unique referral code
    string target_url = 4;                                      // product or category URL
    int64 click_count = 5; int64 conversion_count = 6;
    int64 total_earnings_cents = 7;
    string created_at = 8;
}

message Referral {
    string id = 1; string referrer_id = 2; string referred_id = 3;
    string order_id = 4; int64 order_total_cents = 5;
    int64 commission_cents = 6;                                 // calculated from program rate
    string status = 7;                                          // pending, confirmed, paid
    string created_at = 8;
}

message AffiliateProgram {
    double commission_rate = 1;                                 // e.g., 0.05 = 5%
    int64 min_payout_cents = 2;                                 // minimum $50 to request payout
    int32 cookie_days = 3;                                      // attribution window (30 days)
    int64 referrer_bonus_cents = 4;                             // bonus for referrer
    int64 referred_bonus_cents = 5;                             // discount for new user
}
```

### Tax Service Proto
```protobuf
service TaxService {
    rpc CalculateTax(CalculateTaxRequest) returns (TaxCalculationResponse);           // called by order service
    rpc GetTaxRules(GetTaxRulesRequest) returns (ListTaxRulesResponse);
    rpc CreateTaxRule(CreateTaxRuleRequest) returns (TaxRuleResponse);
    rpc UpdateTaxRule(UpdateTaxRuleRequest) returns (TaxRuleResponse);
    rpc ListTaxZones(Empty) returns (ListTaxZonesResponse);
}

message TaxCalculation {
    int64 subtotal_cents = 1;
    int64 tax_amount_cents = 2;
    repeated TaxBreakdown breakdown = 3;
}

message TaxBreakdown {
    string tax_name = 1;                                        // "GST", "VAT", "State Sales Tax"
    double rate = 2;                                            // 0.10 = 10%
    int64 amount_cents = 3;
    string jurisdiction = 4;                                    // "AU", "CA-ON", "US-CA"
}

message TaxRule {
    string id = 1; string country_code = 2;
    string state_code = 3;                                      // optional (US states, CA provinces)
    string tax_name = 4; double rate = 5;
    string category = 6;                                        // product category override
    bool inclusive = 7;                                          // true = tax included in price (AU/EU)
    string starts_at = 8; string expires_at = 9;
}
```

### CMS Service Proto
```protobuf
service CMSService {
    rpc CreateBanner(CreateBannerRequest) returns (BannerResponse);
    rpc ListBanners(ListBannersRequest) returns (ListBannersResponse);
    rpc UpdateBanner(UpdateBannerRequest) returns (BannerResponse);
    rpc DeleteBanner(DeleteBannerRequest) returns (Empty);
    rpc CreatePage(CreatePageRequest) returns (PageResponse);
    rpc GetPage(GetPageRequest) returns (PageResponse);
    rpc ListPages(ListPagesRequest) returns (ListPagesResponse);
    rpc UpdatePage(UpdatePageRequest) returns (PageResponse);
    rpc ScheduleContent(ScheduleContentRequest) returns (ScheduleResponse);
}

message Banner {
    string id = 1; string title = 2;
    string image_url = 3; string link_url = 4;
    string position = 5;                                        // hero, sidebar, category_top, checkout
    int32 sort_order = 6;
    string starts_at = 7; string ends_at = 8;
    bool is_active = 9;
    string target_audience = 10;                                // all, new_users, returning, specific_tier
}

message Page {
    string id = 1; string title = 2; string slug = 3;
    string content_html = 4;                                    // rich text content
    string meta_title = 5; string meta_description = 6;
    string status = 7;                                          // draft, published, scheduled
    string published_at = 8; string created_at = 9;
}
```

---

## Part 8: Event-Driven Communication (NATS)

### Event Map

| Event Subject | Publisher | Subscribers | Payload |
|---|---|---|---|
| `user.registered` | Auth | User, Notification, Loyalty, Affiliate | userId, email, name, referralCode? |
| `user.verified` | User | Notification | userId |
| `seller.approved` | User (admin action) | Notification, Product | sellerId, userId |
| `product.created` | Product | Search, AI | productId, name, description, categoryId |
| `product.updated` | Product | Search, AI | productId, changed fields |
| `product.deleted` | Product | Search, AI, Cart | productId |
| `cart.checkout` | Cart | Order | userId, cartItems[], couponCode?, pointsRedeem? |
| `order.created` | Order | Payment, Notification, Product (stock), Tax | orderId, items[], total, shippingAddress |
| `order.status_changed` | Order | Notification, User | orderId, oldStatus, newStatus |
| `order.completed` | Order | Loyalty (earn points), Affiliate (conversion), Notification | orderId, buyerId, total, sellerId |
| `payment.succeeded` | Payment | Order, Notification, User, Loyalty | paymentId, orderId, amount |
| `payment.failed` | Payment | Order, Notification | paymentId, orderId, reason |
| `review.created` | Review | Product (rating update), AI (sentiment), Loyalty (earn points) | reviewId, productId, rating |
| `shipment.created` | Shipping | Order, Notification | shipmentId, orderId, trackingNumber, carrier |
| `shipment.status_changed` | Shipping | Order, Notification | shipmentId, status, trackingEvents[] |
| `shipment.delivered` | Shipping | Order, Notification, Review (prompt), Loyalty | shipmentId, orderId |
| `return.requested` | Return | Order, Notification, Seller | returnId, orderId, reason, items[] |
| `return.approved` | Return | Notification, Shipping (return label) | returnId, orderId, refundAmount |
| `return.completed` | Return | Payment (refund), Product (restock), Notification | returnId, orderId, refundAmount |
| `dispute.opened` | Return | Notification (buyer+seller+admin) | disputeId, orderId, type |
| `dispute.resolved` | Return | Notification, Payment (if refund) | disputeId, resolution, refundAmount? |
| `coupon.redeemed` | Promotion | Order | couponId, orderId, discountAmount |
| `flash_sale.started` | Promotion | Product, Notification, Search | flashSaleId, items[] |
| `flash_sale.ended` | Promotion | Product, Search | flashSaleId |
| `points.earned` | Loyalty | Notification | userId, points, source, balance |
| `points.redeemed` | Loyalty | Notification | userId, points, orderId, balance |
| `tier.upgraded` | Loyalty | Notification, User | userId, oldTier, newTier |
| `referral.converted` | Affiliate | Notification, Loyalty | referrerId, referredId, commission |
| `payout.processed` | Payment | Notification (seller) | payoutId, sellerId, amount |
| `payout.requested` | Affiliate | Payment, Notification | affiliateId, amount |
| `ai.embedding.ready` | AI | Search | productId, embedding[] |

---

## Part 9: Kong Gateway Configuration (DB-less Declarative)

### Kong Declarative Config (`kong/kong.yml`)

```yaml
_format_version: "3.0"
_transform: true

# ============================================================
# GLOBAL PLUGINS
# ============================================================
plugins:
  - name: cors
    config:
      origins: ["http://localhost:3000", "https://yourdomain.com"]
      methods: [GET, POST, PUT, PATCH, DELETE, OPTIONS]
      headers: [Authorization, Content-Type, X-Request-ID]
      credentials: true
      max_age: 3600

  - name: rate-limiting
    config:
      minute: 120
      hour: 5000
      policy: redis
      redis:
        host: redis
        port: 6379

  - name: correlation-id
    config:
      header_name: X-Request-ID
      generator: uuid#counter
      echo_downstream: true

  - name: request-size-limiting
    config:
      allowed_payload_size: 10  # MB

  - name: prometheus    # Metrics for monitoring

  - name: file-log
    config:
      path: /dev/stdout
      reopen: false

# ============================================================
# JWT CONSUMER & CREDENTIALS (for token validation)
# ============================================================
consumers:
  - username: ecommerce-app
    jwt_secrets:
      - key: ecommerce-jwt
        algorithm: HS256
        secret: "${JWT_SECRET}"   # Injected via env

# ============================================================
# SERVICES & ROUTES
# ============================================================

# ---- AUTH SERVICE (public, no JWT required) ----
services:
  - name: auth-service
    url: http://auth:8090
    routes:
      - name: auth-register
        paths: ["/api/v1/auth/register"]
        methods: [POST]
        strip_path: false
      - name: auth-login
        paths: ["/api/v1/auth/login"]
        methods: [POST]
        strip_path: false
      - name: auth-refresh
        paths: ["/api/v1/auth/refresh"]
        methods: [POST]
        strip_path: false
      - name: auth-forgot-password
        paths: ["/api/v1/auth/forgot-password"]
        methods: [POST]
        strip_path: false
      - name: auth-reset-password
        paths: ["/api/v1/auth/reset-password"]
        methods: [POST]
        strip_path: false
      - name: auth-oauth
        paths: ["/api/v1/auth/oauth"]
        methods: [GET, POST]
        strip_path: false
      - name: auth-logout
        paths: ["/api/v1/auth/logout"]
        methods: [POST]
        strip_path: false
        plugins:
          - name: jwt    # Logout requires auth

  # ---- USER SERVICE (protected) ----
  - name: user-service
    url: http://user:8091
    routes:
      - name: user-profile
        paths: ["/api/v1/users"]
        strip_path: false
        plugins:
          - name: jwt
          - name: request-transformer
            config:
              add:
                headers:
                  - "X-User-ID:$(jwt.claims.sub)"

      - name: user-addresses
        paths: ["/api/v1/users/me/addresses"]
        strip_path: false
        plugins:
          - name: jwt

      - name: seller-profile
        paths: ["/api/v1/sellers"]
        strip_path: false
        plugins:
          - name: jwt

  # ---- PRODUCT SERVICE ----
  - name: product-service
    url: http://product:8081
    routes:
      - name: products-public
        paths: ["/api/v1/products"]
        methods: [GET]
        strip_path: false
        # No JWT — public browsing

      - name: categories-public
        paths: ["/api/v1/categories"]
        methods: [GET]
        strip_path: false

      - name: seller-products
        paths: ["/api/v1/seller/products"]
        methods: [POST, PATCH, DELETE, GET]
        strip_path: false
        plugins:
          - name: jwt
          - name: acl
            config:
              allow: [seller, admin]

      - name: seller-product-options
        paths: ["/api/v1/seller/products/~*/options"]
        methods: [POST, DELETE, GET]
        strip_path: false
        plugins:
          - name: jwt
          - name: acl
            config:
              allow: [seller, admin]

      - name: seller-product-variants
        paths: ["/api/v1/seller/products/~*/variants"]
        methods: [POST, PATCH, GET]
        strip_path: false
        plugins:
          - name: jwt
          - name: acl
            config:
              allow: [seller, admin]

      - name: admin-attributes
        paths: ["/api/v1/admin/attributes"]
        methods: [POST, PATCH, DELETE, GET]
        strip_path: false
        plugins:
          - name: jwt
          - name: acl
            config:
              allow: [admin]

      - name: category-attributes
        paths: ["/api/v1/categories/~*/attributes"]
        methods: [GET]
        strip_path: false

  # ---- CART SERVICE (protected) ----
  - name: cart-service
    url: http://cart:8082
    routes:
      - name: cart
        paths: ["/api/v1/cart"]
        strip_path: false
        plugins:
          - name: jwt

  # ---- ORDER SERVICE (protected) ----
  - name: order-service
    url: http://order:8083
    routes:
      - name: orders
        paths: ["/api/v1/orders"]
        strip_path: false
        plugins:
          - name: jwt
      - name: seller-orders
        paths: ["/api/v1/seller/orders"]
        strip_path: false
        plugins:
          - name: jwt
          - name: acl
            config:
              allow: [seller, admin]

  # ---- PAYMENT SERVICE ----
  - name: payment-service
    url: http://payment:8084
    routes:
      - name: payments
        paths: ["/api/v1/payments"]
        methods: [POST, GET]
        strip_path: false
        plugins:
          - name: jwt
      - name: stripe-webhooks
        paths: ["/api/v1/payments/webhooks/stripe"]
        methods: [POST]
        strip_path: false
        # No JWT — Stripe signs webhooks with its own secret
        plugins:
          - name: ip-restriction
            config:
              allow: ["3.18.12.63", "3.130.192.202"]   # Stripe webhook IPs

  # ---- SEARCH SERVICE (public) ----
  - name: search-service
    url: http://search:8085
    routes:
      - name: search
        paths: ["/api/v1/search"]
        methods: [GET, POST]
        strip_path: false

  # ---- REVIEW SERVICE ----
  - name: review-service
    url: http://review:8086
    routes:
      - name: reviews-read
        paths: ["/api/v1/products/~*/reviews"]
        methods: [GET]
        strip_path: false
      - name: reviews-write
        paths: ["/api/v1/products/~*/reviews"]
        methods: [POST]
        strip_path: false
        plugins:
          - name: jwt

  # ---- NOTIFICATION SERVICE (protected) ----
  - name: notification-service
    url: http://notification:8087
    routes:
      - name: notifications
        paths: ["/api/v1/notifications"]
        strip_path: false
        plugins:
          - name: jwt
      - name: ws-notifications
        paths: ["/ws/notifications"]
        protocols: [ws, wss]
        strip_path: false
        plugins:
          - name: jwt    # Token passed via query param for WebSocket

  # ---- CHAT SERVICE (protected) ----
  - name: chat-service
    url: http://chat:8088
    routes:
      - name: conversations
        paths: ["/api/v1/conversations"]
        strip_path: false
        plugins:
          - name: jwt
      - name: ws-chat
        paths: ["/ws/chat"]
        protocols: [ws, wss]
        strip_path: false
        plugins:
          - name: jwt

  # ---- MEDIA SERVICE (protected) ----
  - name: media-service
    url: http://media:8089
    routes:
      - name: media-upload
        paths: ["/api/v1/media"]
        methods: [POST, GET, DELETE]
        strip_path: false
        plugins:
          - name: jwt
          - name: request-size-limiting
            config:
              allowed_payload_size: 50  # 50MB for images/videos

  # ---- AI SERVICE ----
  - name: ai-service
    url: http://ai:8092
    routes:
      - name: ai-chat
        paths: ["/api/v1/ai/chat"]
        methods: [POST]
        strip_path: false
        plugins:
          - name: jwt
          - name: rate-limiting
            config:
              minute: 30   # Stricter rate limit for AI endpoints
      - name: ai-recommendations
        paths: ["/api/v1/ai/recommendations"]
        methods: [GET]
        strip_path: false
        plugins:
          - name: jwt
      - name: ai-generate
        paths: ["/api/v1/ai/generate-description"]
        methods: [POST]
        strip_path: false
        plugins:
          - name: jwt
          - name: acl
            config:
              allow: [seller, admin]
      - name: ai-image-search
        paths: ["/api/v1/search/image"]
        methods: [POST]
        strip_path: false

  # ---- PROMOTION SERVICE ----
  - name: promotion-service
    url: http://promotion:8093
    routes:
      - name: coupons-validate
        paths: ["/api/v1/coupons/validate"]
        methods: [POST]
        strip_path: false
        plugins:
          - name: jwt
      - name: coupons-public
        paths: ["/api/v1/coupons"]
        methods: [GET]
        strip_path: false
      - name: seller-coupons
        paths: ["/api/v1/seller/coupons"]
        methods: [POST, PATCH, DELETE, GET]
        strip_path: false
        plugins:
          - name: jwt
          - name: acl
            config:
              allow: [seller, admin]
      - name: flash-sales
        paths: ["/api/v1/flash-sales"]
        methods: [GET]
        strip_path: false
      - name: bundles
        paths: ["/api/v1/bundles"]
        methods: [GET]
        strip_path: false
      - name: admin-promotions
        paths: ["/api/v1/admin/promotions"]
        methods: [POST, PATCH, DELETE, GET]
        strip_path: false
        plugins:
          - name: jwt
          - name: acl
            config:
              allow: [admin]

  # ---- RETURN SERVICE ----
  - name: return-service
    url: http://return:8094
    routes:
      - name: buyer-returns
        paths: ["/api/v1/returns"]
        methods: [POST, GET]
        strip_path: false
        plugins:
          - name: jwt
      - name: buyer-return-detail
        paths: ["/api/v1/returns/~*"]
        methods: [GET]
        strip_path: false
        plugins:
          - name: jwt
      - name: seller-returns
        paths: ["/api/v1/seller/returns"]
        methods: [GET, PATCH]
        strip_path: false
        plugins:
          - name: jwt
          - name: acl
            config:
              allow: [seller, admin]
      - name: disputes
        paths: ["/api/v1/disputes"]
        methods: [POST, GET]
        strip_path: false
        plugins:
          - name: jwt
      - name: admin-disputes
        paths: ["/api/v1/admin/disputes"]
        methods: [GET, PATCH]
        strip_path: false
        plugins:
          - name: jwt
          - name: acl
            config:
              allow: [admin]

  # ---- SHIPPING SERVICE ----
  - name: shipping-service
    url: http://shipping:8095
    routes:
      - name: shipping-rates
        paths: ["/api/v1/shipping/rates"]
        methods: [POST]
        strip_path: false
        plugins:
          - name: jwt
      - name: shipments
        paths: ["/api/v1/shipments"]
        methods: [POST, GET]
        strip_path: false
        plugins:
          - name: jwt
      - name: shipment-tracking
        paths: ["/api/v1/tracking/~*"]
        methods: [GET]
        strip_path: false
      - name: seller-shipments
        paths: ["/api/v1/seller/shipments"]
        methods: [POST, GET]
        strip_path: false
        plugins:
          - name: jwt
          - name: acl
            config:
              allow: [seller, admin]
      - name: admin-carriers
        paths: ["/api/v1/admin/carriers"]
        methods: [POST, PATCH, GET]
        strip_path: false
        plugins:
          - name: jwt
          - name: acl
            config:
              allow: [admin]

  # ---- LOYALTY SERVICE ----
  - name: loyalty-service
    url: http://loyalty:8096
    routes:
      - name: loyalty-membership
        paths: ["/api/v1/loyalty/membership"]
        methods: [GET]
        strip_path: false
        plugins:
          - name: jwt
      - name: loyalty-points
        paths: ["/api/v1/loyalty/points"]
        methods: [GET]
        strip_path: false
        plugins:
          - name: jwt
      - name: loyalty-transactions
        paths: ["/api/v1/loyalty/transactions"]
        methods: [GET]
        strip_path: false
        plugins:
          - name: jwt
      - name: loyalty-redeem
        paths: ["/api/v1/loyalty/redeem"]
        methods: [POST]
        strip_path: false
        plugins:
          - name: jwt
      - name: loyalty-tiers
        paths: ["/api/v1/loyalty/tiers"]
        methods: [GET]
        strip_path: false

  # ---- AFFILIATE SERVICE ----
  - name: affiliate-service
    url: http://affiliate:8097
    routes:
      - name: affiliate-links
        paths: ["/api/v1/affiliate/links"]
        methods: [POST, GET]
        strip_path: false
        plugins:
          - name: jwt
      - name: affiliate-stats
        paths: ["/api/v1/affiliate/stats"]
        methods: [GET]
        strip_path: false
        plugins:
          - name: jwt
      - name: affiliate-referrals
        paths: ["/api/v1/affiliate/referrals"]
        methods: [GET]
        strip_path: false
        plugins:
          - name: jwt
      - name: affiliate-payout
        paths: ["/api/v1/affiliate/payout"]
        methods: [POST]
        strip_path: false
        plugins:
          - name: jwt
      - name: affiliate-track
        paths: ["/api/v1/r/~*"]
        methods: [GET]
        strip_path: false
      - name: admin-affiliates
        paths: ["/api/v1/admin/affiliates"]
        methods: [GET, PATCH]
        strip_path: false
        plugins:
          - name: jwt
          - name: acl
            config:
              allow: [admin]

  # ---- TAX SERVICE (internal + admin) ----
  - name: tax-service
    url: http://tax:8098
    routes:
      - name: admin-tax-rules
        paths: ["/api/v1/admin/tax"]
        methods: [POST, PATCH, DELETE, GET]
        strip_path: false
        plugins:
          - name: jwt
          - name: acl
            config:
              allow: [admin]
      - name: tax-zones
        paths: ["/api/v1/tax/zones"]
        methods: [GET]
        strip_path: false

  # ---- CMS SERVICE ----
  - name: cms-service
    url: http://cms:8099
    routes:
      - name: banners-public
        paths: ["/api/v1/banners"]
        methods: [GET]
        strip_path: false
      - name: pages-public
        paths: ["/api/v1/pages/~*"]
        methods: [GET]
        strip_path: false
      - name: admin-banners
        paths: ["/api/v1/admin/banners"]
        methods: [POST, PATCH, DELETE, GET]
        strip_path: false
        plugins:
          - name: jwt
          - name: acl
            config:
              allow: [admin]
      - name: admin-pages
        paths: ["/api/v1/admin/pages"]
        methods: [POST, PATCH, DELETE, GET]
        strip_path: false
        plugins:
          - name: jwt
          - name: acl
            config:
              allow: [admin]

  # ---- ADMIN SERVICE (admin only) ----
  - name: admin-service
    url: http://admin:8093
    routes:
      - name: admin
        paths: ["/api/v1/admin"]
        strip_path: false
        plugins:
          - name: jwt
          - name: acl
            config:
              allow: [admin]
          - name: rate-limiting
            config:
              minute: 300
```

### Kong Plugins Used

| Plugin | Purpose | Scope |
|---|---|---|
| **jwt** | Validate JWT tokens, extract claims | Per-route (protected endpoints) |
| **acl** | Role-based access control (seller, admin) | Per-route |
| **rate-limiting** | Request throttling (Redis-backed) | Global + per-route overrides |
| **cors** | Cross-origin resource sharing | Global |
| **correlation-id** | Distributed tracing request ID | Global |
| **request-transformer** | Add X-User-ID header from JWT claims | Per-route |
| **ip-restriction** | Whitelist Stripe webhook IPs | Stripe webhook route |
| **request-size-limiting** | Max payload size | Global (10MB) + media (50MB) |
| **prometheus** | Metrics for Grafana dashboards | Global |
| **file-log** | Structured access logs | Global |

### How JWT Auth Works with Kong

```
1. User logs in → Auth Service returns JWT (access + refresh tokens)
2. Client sends requests with: Authorization: Bearer <token>
3. Kong JWT plugin validates token signature + expiry
4. Kong adds X-User-ID header (from jwt.claims.sub) to upstream request
5. Go service reads X-User-ID from request header — no JWT parsing needed
6. For role checks: Kong ACL plugin checks jwt.claims.role against allowed groups
```

### Each Go Service HTTP Handler Pattern (with Kong)
```go
// Services no longer need JWT middleware — Kong handles it
// Services just read the X-User-ID and X-User-Role headers set by Kong

func (h *ProductHandler) CreateProduct(c *gin.Context) {
    userID := c.GetHeader("X-User-ID")       // Set by Kong from JWT claims
    userRole := c.GetHeader("X-User-Role")   // Set by Kong from JWT claims

    if userID == "" {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    var input CreateProductRequest
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    product, err := h.useCase.CreateProduct(c.Request.Context(), usecase.CreateProductInput{
        SellerID: userID,
        Name:     input.Name,
        // ...
    })
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(201, product)
}
```

---

## Part 10: Database Schema (Key Tables)

### Product Database (`ecommerce_products`)

```sql
-- Core product table
CREATE TABLE products (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id       UUID NOT NULL,
    category_id     UUID REFERENCES categories(id),
    name            VARCHAR(500) NOT NULL,
    slug            VARCHAR(500) UNIQUE NOT NULL,
    description     TEXT,
    base_price_cents BIGINT NOT NULL,
    currency        VARCHAR(3) NOT NULL DEFAULT 'USD',
    status          VARCHAR(20) NOT NULL DEFAULT 'draft',  -- draft, active, inactive, archived
    has_variants    BOOLEAN NOT NULL DEFAULT false,
    tags            TEXT[],
    image_urls      TEXT[],
    embedding       vector(1536),                          -- pgvector AI embedding
    rating_avg      DECIMAL(3,2) DEFAULT 0,
    rating_count    INT DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_products_seller ON products(seller_id);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_slug ON products(slug);

-- Admin-defined attribute templates
CREATE TABLE attribute_definitions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,            -- "Brand", "Material"
    slug        VARCHAR(100) UNIQUE NOT NULL,      -- "brand", "material"
    type        VARCHAR(20) NOT NULL,              -- text, number, select, multi_select, color, bool
    required    BOOLEAN NOT NULL DEFAULT false,
    filterable  BOOLEAN NOT NULL DEFAULT false,    -- indexed for faceted search
    options     TEXT[],                             -- predefined choices for select types
    unit        VARCHAR(20),                       -- "cm", "kg", "inches"
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Which attributes belong to which category
CREATE TABLE category_attributes (
    category_id   UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    attribute_id  UUID NOT NULL REFERENCES attribute_definitions(id) ON DELETE CASCADE,
    sort_order    INT NOT NULL DEFAULT 0,
    PRIMARY KEY (category_id, attribute_id)
);

-- Seller-provided attribute values for a product
CREATE TABLE product_attribute_values (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id    UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    attribute_id  UUID NOT NULL REFERENCES attribute_definitions(id),
    attribute_name VARCHAR(100) NOT NULL,          -- denormalized for reads
    value         TEXT,                            -- single value
    values        TEXT[],                          -- for multi_select
    UNIQUE(product_id, attribute_id)
);
CREATE INDEX idx_pav_product ON product_attribute_values(product_id);
CREATE INDEX idx_pav_attribute_value ON product_attribute_values(attribute_id, value);  -- faceted search

-- Seller-defined variant axes (max 3 per product)
CREATE TABLE product_options (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id  UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    name        VARCHAR(50) NOT NULL,              -- "Color", "Size"
    sort_order  INT NOT NULL DEFAULT 0
);
CREATE INDEX idx_options_product ON product_options(product_id);

-- Values for each option
CREATE TABLE product_option_values (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    option_id   UUID NOT NULL REFERENCES product_options(id) ON DELETE CASCADE,
    value       VARCHAR(100) NOT NULL,             -- "Red", "XL"
    color_hex   VARCHAR(7),                        -- "#FF0000" for color swatches
    sort_order  INT NOT NULL DEFAULT 0
);
CREATE INDEX idx_option_values_option ON product_option_values(option_id);

-- Purchasable variants (one per option combination)
CREATE TABLE product_variants (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id      UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku             VARCHAR(100) UNIQUE NOT NULL,
    name            VARCHAR(200) NOT NULL,         -- "Red / XL"
    price_cents     BIGINT NOT NULL DEFAULT 0,     -- 0 = use base price
    compare_at_cents BIGINT NOT NULL DEFAULT 0,    -- strikethrough price
    cost_cents      BIGINT NOT NULL DEFAULT 0,     -- cost of goods
    stock           INT NOT NULL DEFAULT 0,
    low_stock_alert INT NOT NULL DEFAULT 5,
    weight_grams    INT NOT NULL DEFAULT 0,
    is_default      BOOLEAN NOT NULL DEFAULT false,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    image_urls      TEXT[],
    barcode         VARCHAR(50),                   -- UPC/EAN/ISBN
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_variants_product ON product_variants(product_id);
CREATE INDEX idx_variants_sku ON product_variants(sku);

-- Links variant to its option values
CREATE TABLE variant_option_values (
    variant_id      UUID NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    option_id       UUID NOT NULL REFERENCES product_options(id),
    option_value_id UUID NOT NULL REFERENCES product_option_values(id),
    option_name     VARCHAR(50) NOT NULL,          -- denormalized: "Color"
    value           VARCHAR(100) NOT NULL,         -- denormalized: "Red"
    PRIMARY KEY (variant_id, option_id)
);

-- Product images with AI embeddings
CREATE TABLE product_images (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id      UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    variant_id      UUID REFERENCES product_variants(id) ON DELETE SET NULL,
    url             TEXT NOT NULL,
    alt_text        VARCHAR(500),
    image_embedding vector(512),                   -- CLIP visual similarity
    sort_order      INT NOT NULL DEFAULT 0
);
```

### Other AI-Specific Columns

| Table | AI Column | Type | Purpose |
|---|---|---|---|
| `products` | `embedding` | `vector(1536)` | Semantic search via pgvector |
| `product_images` | `image_embedding` | `vector(512)` | Visual similarity (CLIP) |
| `orders` | `fraud_score` | `decimal(5,4)` | ML fraud detection score |
| `reviews` | `sentiment_score` | `decimal(3,2)` | NLP sentiment analysis |
| `user_interactions` | — | — | Behavioral data for recommendations |

### Promotion Database (`ecommerce_promotions`)

```sql
CREATE TABLE coupons (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code            VARCHAR(50) UNIQUE NOT NULL,
    type            VARCHAR(20) NOT NULL,               -- percentage, fixed_amount, free_shipping
    discount_value  BIGINT NOT NULL,                    -- cents or percentage * 100
    min_order_cents BIGINT DEFAULT 0,
    max_discount_cents BIGINT DEFAULT 0,                -- cap for percentage discounts
    usage_limit     INT DEFAULT 0,                      -- 0 = unlimited
    usage_count     INT DEFAULT 0,
    per_user_limit  INT DEFAULT 1,
    scope           VARCHAR(20) DEFAULT 'all',          -- all, category, product, seller
    scope_ids       TEXT[],                              -- IDs for scoped coupons
    created_by      VARCHAR(50) NOT NULL,                -- seller_id or "platform"
    starts_at       TIMESTAMPTZ NOT NULL,
    expires_at      TIMESTAMPTZ,
    is_active       BOOLEAN DEFAULT true,
    created_at      TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX idx_coupons_code ON coupons(code);

CREATE TABLE coupon_usages (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    coupon_id   UUID NOT NULL REFERENCES coupons(id),
    user_id     UUID NOT NULL,
    order_id    UUID NOT NULL,
    discount_cents BIGINT NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT now(),
    UNIQUE(coupon_id, user_id, order_id)
);

CREATE TABLE flash_sales (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(200) NOT NULL,
    starts_at   TIMESTAMPTZ NOT NULL,
    ends_at     TIMESTAMPTZ NOT NULL,
    is_active   BOOLEAN DEFAULT true,
    created_at  TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE flash_sale_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flash_sale_id   UUID NOT NULL REFERENCES flash_sales(id) ON DELETE CASCADE,
    product_id      UUID NOT NULL,
    variant_id      UUID,
    sale_price_cents BIGINT NOT NULL,
    quantity_limit  INT NOT NULL,
    sold_count      INT DEFAULT 0
);

CREATE TABLE bundles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(200) NOT NULL,
    seller_id       UUID NOT NULL,
    product_ids     UUID[] NOT NULL,
    bundle_price_cents BIGINT NOT NULL,
    savings_cents   BIGINT NOT NULL,
    is_active       BOOLEAN DEFAULT true,
    created_at      TIMESTAMPTZ DEFAULT now()
);
```

### Return Database (`ecommerce_returns`)

```sql
CREATE TABLE returns (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id            UUID NOT NULL,
    buyer_id            UUID NOT NULL,
    seller_id           UUID NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'requested',   -- requested, approved, rejected, shipped_back, received, refunded
    reason              VARCHAR(50) NOT NULL,                       -- defective, wrong_item, not_as_described, changed_mind
    description         TEXT,
    image_urls          TEXT[],
    refund_amount_cents BIGINT,
    refund_method       VARCHAR(20),                                -- original_payment, wallet_credit
    return_tracking     VARCHAR(100),
    created_at          TIMESTAMPTZ DEFAULT now(),
    updated_at          TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX idx_returns_order ON returns(order_id);
CREATE INDEX idx_returns_buyer ON returns(buyer_id);
CREATE INDEX idx_returns_seller ON returns(seller_id);

CREATE TABLE return_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    return_id       UUID NOT NULL REFERENCES returns(id) ON DELETE CASCADE,
    order_item_id   UUID NOT NULL,
    product_id      UUID NOT NULL,
    variant_id      UUID,
    quantity        INT NOT NULL,
    reason          VARCHAR(50)
);

CREATE TABLE disputes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id        UUID NOT NULL,
    return_id       UUID REFERENCES returns(id),
    buyer_id        UUID NOT NULL,
    seller_id       UUID NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'open',            -- open, under_review, resolved_buyer, resolved_seller, escalated
    type            VARCHAR(30) NOT NULL,                           -- item_not_received, not_as_described, unauthorized_charge
    description     TEXT NOT NULL,
    resolution      TEXT,
    resolved_by     UUID,
    created_at      TIMESTAMPTZ DEFAULT now(),
    resolved_at     TIMESTAMPTZ
);

CREATE TABLE dispute_messages (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dispute_id  UUID NOT NULL REFERENCES disputes(id) ON DELETE CASCADE,
    sender_id   UUID NOT NULL,
    sender_role VARCHAR(10) NOT NULL,                               -- buyer, seller, admin
    message     TEXT NOT NULL,
    attachments TEXT[],
    created_at  TIMESTAMPTZ DEFAULT now()
);
```

### Shipping Database (`ecommerce_shipping`)

```sql
CREATE TABLE carriers (
    code                VARCHAR(20) PRIMARY KEY,                    -- fedex, ups, dhl, auspost
    name                VARCHAR(100) NOT NULL,
    is_active           BOOLEAN DEFAULT true,
    supported_countries TEXT[] NOT NULL,
    api_base_url        VARCHAR(500),
    created_at          TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE carrier_credentials (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id   UUID NOT NULL,
    carrier_code VARCHAR(20) NOT NULL REFERENCES carriers(code),
    credentials JSONB NOT NULL,                                     -- encrypted API keys
    is_active   BOOLEAN DEFAULT true,
    UNIQUE(seller_id, carrier_code)
);

CREATE TABLE shipments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id        UUID NOT NULL,
    seller_id       UUID NOT NULL,
    carrier_code    VARCHAR(20) REFERENCES carriers(code),
    service_code    VARCHAR(50),
    tracking_number VARCHAR(100),
    label_url       TEXT,
    status          VARCHAR(20) DEFAULT 'pending',                  -- pending, label_created, picked_up, in_transit, delivered, exception
    origin          JSONB NOT NULL,                                 -- address JSON
    destination     JSONB NOT NULL,
    weight_grams    INT,
    rate_cents      BIGINT,
    currency        VARCHAR(3) DEFAULT 'USD',
    created_at      TIMESTAMPTZ DEFAULT now(),
    updated_at      TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX idx_shipments_order ON shipments(order_id);
CREATE INDEX idx_shipments_tracking ON shipments(tracking_number);

CREATE TABLE tracking_events (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipment_id UUID NOT NULL REFERENCES shipments(id) ON DELETE CASCADE,
    status      VARCHAR(50) NOT NULL,
    description TEXT,
    location    VARCHAR(200),
    event_at    TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT now()
);
```

### Loyalty Database (`ecommerce_loyalty`)

```sql
CREATE TABLE memberships (
    user_id         UUID PRIMARY KEY,
    tier            VARCHAR(20) NOT NULL DEFAULT 'bronze',          -- bronze, silver, gold, platinum
    points_balance  BIGINT NOT NULL DEFAULT 0,
    lifetime_points BIGINT NOT NULL DEFAULT 0,
    tier_expires_at TIMESTAMPTZ,
    joined_at       TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE points_transactions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES memberships(user_id),
    type            VARCHAR(10) NOT NULL,                           -- earn, redeem, expire, adjust
    points          BIGINT NOT NULL,
    source          VARCHAR(20) NOT NULL,                           -- order, review, referral, promotion, signup
    reference_id    VARCHAR(100),                                   -- orderId, reviewId
    description     VARCHAR(500),
    created_at      TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX idx_points_user ON points_transactions(user_id);

CREATE TABLE tiers (
    name                VARCHAR(20) PRIMARY KEY,
    min_points          BIGINT NOT NULL,
    cashback_rate       DECIMAL(5,4) NOT NULL,                      -- 0.02 = 2%
    points_multiplier   DECIMAL(3,1) NOT NULL DEFAULT 1.0,          -- 1.5x for gold
    free_shipping       BOOLEAN DEFAULT false,
    priority_support_hours INT DEFAULT 48
);

INSERT INTO tiers VALUES
    ('bronze',   0,     0.01, 1.0, false, 48),
    ('silver',   1000,  0.02, 1.2, false, 24),
    ('gold',     5000,  0.03, 1.5, true,  12),
    ('platinum', 15000, 0.05, 2.0, true,  4);
```

### Affiliate Database (`ecommerce_affiliates`)

```sql
CREATE TABLE affiliate_program (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    commission_rate     DECIMAL(5,4) NOT NULL DEFAULT 0.05,         -- 5%
    min_payout_cents    BIGINT NOT NULL DEFAULT 5000,               -- $50
    cookie_days         INT NOT NULL DEFAULT 30,
    referrer_bonus_cents BIGINT DEFAULT 0,
    referred_bonus_cents BIGINT DEFAULT 0,
    is_active           BOOLEAN DEFAULT true
);

CREATE TABLE affiliate_links (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL,
    code        VARCHAR(20) UNIQUE NOT NULL,
    target_url  TEXT,
    click_count BIGINT DEFAULT 0,
    conversion_count BIGINT DEFAULT 0,
    total_earnings_cents BIGINT DEFAULT 0,
    created_at  TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX idx_affiliate_links_user ON affiliate_links(user_id);
CREATE INDEX idx_affiliate_links_code ON affiliate_links(code);

CREATE TABLE referrals (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    referrer_id     UUID NOT NULL,
    referred_id     UUID NOT NULL,
    order_id        UUID,
    order_total_cents BIGINT,
    commission_cents BIGINT,
    status          VARCHAR(20) DEFAULT 'pending',                  -- pending, confirmed, paid
    created_at      TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE affiliate_payouts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    amount_cents    BIGINT NOT NULL,
    status          VARCHAR(20) DEFAULT 'requested',                -- requested, processing, completed, failed
    payout_method   VARCHAR(20) NOT NULL,                           -- bank_transfer, stripe
    created_at      TIMESTAMPTZ DEFAULT now(),
    completed_at    TIMESTAMPTZ
);
```

### Tax Database (`ecommerce_tax`)

```sql
CREATE TABLE tax_zones (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_code    VARCHAR(2) NOT NULL,
    state_code      VARCHAR(10),                                    -- US states, CA provinces
    name            VARCHAR(100) NOT NULL,                          -- "Australia", "California", "Ontario"
    UNIQUE(country_code, state_code)
);

CREATE TABLE tax_rules (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    zone_id         UUID NOT NULL REFERENCES tax_zones(id),
    tax_name        VARCHAR(50) NOT NULL,                           -- "GST", "VAT", "State Sales Tax"
    rate            DECIMAL(7,5) NOT NULL,                          -- 0.10000 = 10%
    category        VARCHAR(50),                                    -- product category override (null = all)
    inclusive       BOOLEAN DEFAULT false,                           -- true = price includes tax (AU/EU style)
    starts_at       TIMESTAMPTZ DEFAULT now(),
    expires_at      TIMESTAMPTZ,
    is_active       BOOLEAN DEFAULT true
);
CREATE INDEX idx_tax_rules_zone ON tax_rules(zone_id);
```

### CMS Database (`ecommerce_cms`)

```sql
CREATE TABLE banners (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           VARCHAR(200) NOT NULL,
    image_url       TEXT NOT NULL,
    link_url        TEXT,
    position        VARCHAR(30) NOT NULL,                           -- hero, sidebar, category_top, checkout
    sort_order      INT DEFAULT 0,
    target_audience VARCHAR(20) DEFAULT 'all',                      -- all, new_users, returning, specific_tier
    starts_at       TIMESTAMPTZ DEFAULT now(),
    ends_at         TIMESTAMPTZ,
    is_active       BOOLEAN DEFAULT true,
    created_at      TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE pages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           VARCHAR(200) NOT NULL,
    slug            VARCHAR(200) UNIQUE NOT NULL,
    content_html    TEXT NOT NULL,
    meta_title      VARCHAR(200),
    meta_description VARCHAR(500),
    status          VARCHAR(20) DEFAULT 'draft',                    -- draft, published, scheduled
    published_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT now(),
    updated_at      TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE content_schedules (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_type    VARCHAR(20) NOT NULL,                           -- banner, page
    content_id      UUID NOT NULL,
    action          VARCHAR(20) NOT NULL,                           -- publish, unpublish, activate, deactivate
    scheduled_at    TIMESTAMPTZ NOT NULL,
    executed        BOOLEAN DEFAULT false,
    created_at      TIMESTAMPTZ DEFAULT now()
);
```

### Payment Database — Seller Wallet & Settlement Additions

```sql
-- Added to ecommerce_payments database

CREATE TABLE seller_wallets (
    seller_id           UUID PRIMARY KEY,
    available_balance   BIGINT NOT NULL DEFAULT 0,                  -- available for withdrawal
    pending_balance     BIGINT NOT NULL DEFAULT 0,                  -- from orders not yet settled
    currency            VARCHAR(3) DEFAULT 'USD',
    updated_at          TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE wallet_transactions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id       UUID NOT NULL REFERENCES seller_wallets(seller_id),
    type            VARCHAR(20) NOT NULL,                           -- sale, commission_deducted, payout, refund_debit, adjustment
    amount_cents    BIGINT NOT NULL,                                -- positive = credit, negative = debit
    reference_type  VARCHAR(20),                                    -- order, payout, refund
    reference_id    UUID,
    description     VARCHAR(500),
    created_at      TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX idx_wallet_tx_seller ON wallet_transactions(seller_id);

CREATE TABLE settlement_periods (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id       UUID NOT NULL,
    period_start    TIMESTAMPTZ NOT NULL,
    period_end      TIMESTAMPTZ NOT NULL,
    total_sales     BIGINT NOT NULL DEFAULT 0,
    total_commission BIGINT NOT NULL DEFAULT 0,
    total_refunds   BIGINT NOT NULL DEFAULT 0,
    net_amount      BIGINT NOT NULL DEFAULT 0,
    status          VARCHAR(20) DEFAULT 'pending',                  -- pending, calculated, settled
    settled_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE payouts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id       UUID NOT NULL,
    amount_cents    BIGINT NOT NULL,
    currency        VARCHAR(3) DEFAULT 'USD',
    method          VARCHAR(20) NOT NULL,                           -- stripe_connect, bank_transfer
    stripe_transfer_id VARCHAR(100),
    status          VARCHAR(20) DEFAULT 'requested',                -- requested, processing, completed, failed
    requested_at    TIMESTAMPTZ DEFAULT now(),
    completed_at    TIMESTAMPTZ
);
CREATE INDEX idx_payouts_seller ON payouts(seller_id);
```

### User Database — Social Features Additions

```sql
-- Added to ecommerce_users database

CREATE TABLE user_follows (
    follower_id     UUID NOT NULL,
    seller_id       UUID NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (follower_id, seller_id)
);
CREATE INDEX idx_follows_seller ON user_follows(seller_id);

CREATE TABLE product_shares (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID,                                           -- null for anonymous shares
    product_id      UUID NOT NULL,
    platform        VARCHAR(20) NOT NULL,                           -- facebook, twitter, whatsapp, copy_link, email
    created_at      TIMESTAMPTZ DEFAULT now()
);
```

---

## Part 11: Shared Packages (`pkg/`)

| Package | Purpose |
|---|---|
| `pkg/logger` | zerolog structured JSON logging with service name, correlation ID |
| `pkg/errors` | Custom error types with gRPC status code mapping |
| `pkg/auth` | JWT generation/validation, claims extraction |
| `pkg/middleware` | Gin middleware: internal auth header parsing (X-User-ID from Kong), logging |
| `pkg/pagination` | Cursor + offset pagination helpers |
| `pkg/events` | NATS publisher/subscriber + event subject constants |
| `pkg/server` | Graceful shutdown for dual gRPC+HTTP servers |
| `pkg/circuitbreaker` | gobreaker-based gRPC client interceptor |
| `pkg/tracing` | OpenTelemetry init with OTLP exporter |
| `pkg/unitofwork` | GORM transaction wrapper (Unit of Work pattern) |
| `pkg/validator` | Input validation helpers |
| `pkg/i18n` | Multi-language support: message catalogs, locale detection from Accept-Language header |
| `pkg/money` | Currency-safe money type with multi-currency formatting and conversion |
| `pkg/tax` | Tax calculation client (gRPC call to tax-service) |

---

## Part 12: Key Patterns

- **Kong handles cross-cutting concerns**: JWT validation, CORS, rate limiting, ACL, correlation IDs, logging, IP restrictions — services stay focused on business logic
- **Dependency Injection**: Constructor injection in `cmd/main.go` — no DI framework needed
- **Repository Pattern**: Interface in `domain/repository/`, GORM implementation in `adapter/repository/postgres/`
- **Unit of Work**: GORM transaction wrapper for atomic multi-table operations (e.g., create order + items)
- **Domain Events**: NATS publish after successful DB commit, subscribers in other services
- **Circuit Breaker**: gobreaker wrapping gRPC client calls with configurable thresholds
- **Order State Machine**: Enum-based transitions with allowed-transitions map
- **Graceful Shutdown**: Signal handling, health status flip, drain in-flight requests
- **Health Checks**: HTTP `/health` + `/ready` endpoints, gRPC health protocol — Kong uses these for upstream health
- **Dual HTTP+gRPC per service**: HTTP endpoints for Kong routing, gRPC for inter-service communication
- **OpenTelemetry**: Distributed tracing across gRPC calls via otelgrpc interceptors + Kong OpenTelemetry plugin

---

## Part 13: Testing Strategy

| Layer | Tool | What to Test |
|---|---|---|
| Domain entities | Go stdlib + testify | Business logic, validations, state transitions |
| Use cases | testify mock, table-driven tests | Orchestration, event publishing, error paths |
| Repositories | testcontainers-go + real PostgreSQL | SQL queries, filters, pagination, GORM behavior |
| gRPC handlers | testify mock on use cases | Proto conversion, gRPC error code mapping |
| Kong config | decK validate (lint kong.yml) | Route definitions, plugin config correctness |
| E2E (backend) | docker-compose + HTTP via Kong | Full request flow through Kong to services |
| Mock generation | mockery | Auto-generate mocks from repository interfaces |
| Flutter domain | dart test | Entity logic, use case Either results |
| Flutter BLoC | bloc_test + mockito | State transitions, event→state mapping |
| Flutter widgets | flutter_test | Widget rendering, user interaction |
| Flutter integration | integration_test | Full user flows on device/emulator |
| Flutter golden | golden_toolkit | Visual regression for UI components |

---

## Part 14: Docker Compose (All Services)

Infrastructure: **Kong Gateway (DB-less)**, PostgreSQL (pgvector), Redis, NATS JetStream, Elasticsearch, MinIO, Mailhog
Go services: auth, user, product, cart, order, payment, search, review, notification, chat, media, ai
Apps: ai-services (Python FastAPI), web (Next.js)

Each Go service connects to its own database (`ecommerce_<name>`) and shared Redis/NATS.

### Kong in Docker Compose
```yaml
kong:
  image: kong/kong-gateway:3.6
  environment:
    KONG_DATABASE: "off"                          # DB-less mode
    KONG_DECLARATIVE_CONFIG: /etc/kong/kong.yml    # Declarative YAML
    KONG_PROXY_LISTEN: "0.0.0.0:8000"
    KONG_ADMIN_LISTEN: "0.0.0.0:8001"
    KONG_LOG_LEVEL: info
    KONG_PLUGINS: "bundled"                       # All bundled plugins available
    KONG_NGINX_PROXY_PROXY_BUFFER_SIZE: "128k"
    KONG_NGINX_PROXY_PROXY_BUFFERS: "4 256k"
  ports:
    - "8000:8000"    # Proxy (all API traffic enters here)
    - "8001:8001"    # Admin API (dev only, disable in prod)
  volumes:
    - ./kong/kong.yml:/etc/kong/kong.yml:ro
  healthcheck:
    test: ["CMD", "kong", "health"]
    interval: 10s
    timeout: 5s
    retries: 5
  depends_on:
    auth:
      condition: service_started
    product:
      condition: service_started
  restart: unless-stopped
```

### Frontend connects to Kong
```
NEXT_PUBLIC_API_URL=http://localhost:8000/api/v1    # Kong proxy port
NEXT_PUBLIC_WS_URL=ws://localhost:8000/ws           # WebSocket via Kong
```

---

## Part 15: CI/CD Pipeline

GitHub Actions with matrix strategy:
- **Lint**: golangci-lint per service
- **Test**: `go test ./...` per service with PostgreSQL + Redis services
- **Build**: Multi-stage Docker build per service
- **Deploy**: Push images to registry, update K8s manifests

---

## Part 16: AI Integration Architecture

```
Product Created/Updated
  → NATS event "product.created"
  → AI Service subscribes
  → Calls Python FastAPI /embeddings endpoint
  → Stores vector in pgvector (ecommerce_ai database)
  → Publishes "ai.embedding.ready"
  → Search Service updates Elasticsearch index with embedding

User searches "comfortable running shoes for flat feet"
  → Gateway → Search Service
  → Generate query embedding (via AI Service → Python)
  → Hybrid search: pgvector cosine similarity + Elasticsearch BM25
  → Return ranked results
```

---

## Part 17: Frontend — Vite + React + Tailwind (Module-Based Architecture)

### Tech Stack
- **Vite** — fast HMR, optimized production builds
- **React 18** + TypeScript
- **React Router v6** — file-based lazy loading per module
- **Tailwind CSS** — utility-first styling
- **Shadcn/ui** — Radix-based accessible components (copy-paste, fully owned)
- **TanStack Query v5** — server state, caching, optimistic updates
- **Zustand** — lightweight client state (cart, auth, UI)
- **React Hook Form + Zod** — form handling + validation
- **Axios** — HTTP client with interceptors (JWT refresh, error handling)
- **Socket.io Client** — real-time notifications & chat

### Module-Based Folder Structure

Each module is a **self-contained feature** with its own components, hooks, services, types, and routes. Modules only import from `shared/` or other modules via explicit barrel exports.

```
apps/web/
├── index.html
├── vite.config.ts
├── tailwind.config.ts
├── tsconfig.json
├── package.json
├── Dockerfile
│
├── public/
│   ├── favicon.ico
│   └── assets/
│
└── src/
    ├── main.tsx                          # App entry point
    ├── App.tsx                           # Root: providers, router
    ├── routes.tsx                        # Central route definitions (lazy imports)
    ├── vite-env.d.ts
    │
    ├── shared/                           # ═══ SHARED (cross-module) ═══
    │   ├── components/
    │   │   ├── ui/                       # Shadcn/ui primitives
    │   │   │   ├── button.tsx
    │   │   │   ├── input.tsx
    │   │   │   ├── dialog.tsx
    │   │   │   ├── dropdown-menu.tsx
    │   │   │   ├── select.tsx
    │   │   │   ├── table.tsx
    │   │   │   ├── card.tsx
    │   │   │   ├── badge.tsx
    │   │   │   ├── toast.tsx
    │   │   │   ├── skeleton.tsx
    │   │   │   ├── avatar.tsx
    │   │   │   ├── tabs.tsx
    │   │   │   └── ... (all shadcn components)
    │   │   ├── layout/
    │   │   │   ├── RootLayout.tsx        # Shell: header + main + footer
    │   │   │   ├── Header.tsx            # Nav, search bar, cart icon, user menu
    │   │   │   ├── Footer.tsx
    │   │   │   ├── Sidebar.tsx           # Dashboard sidebar (seller/admin)
    │   │   │   ├── MobileNav.tsx
    │   │   │   └── DashboardLayout.tsx   # Sidebar + content layout
    │   │   ├── feedback/
    │   │   │   ├── LoadingSpinner.tsx
    │   │   │   ├── ErrorBoundary.tsx
    │   │   │   ├── EmptyState.tsx
    │   │   │   └── PageNotFound.tsx
    │   │   ├── forms/
    │   │   │   ├── FormField.tsx         # Reusable form field wrapper
    │   │   │   ├── ImageUploader.tsx
    │   │   │   └── AddressForm.tsx
    │   │   └── data-display/
    │   │       ├── DataTable.tsx          # Generic sortable/filterable table
    │   │       ├── Pagination.tsx
    │   │       ├── PriceDisplay.tsx
    │   │       └── RatingStars.tsx
    │   │
    │   ├── hooks/
    │   │   ├── useAuth.ts                # Auth state + actions (wraps Zustand store)
    │   │   ├── useSocket.ts              # Socket.io connection management
    │   │   ├── useDebounce.ts
    │   │   ├── useMediaQuery.ts
    │   │   ├── useInfiniteScroll.ts
    │   │   └── useLocalStorage.ts
    │   │
    │   ├── stores/                       # Zustand stores (client state)
    │   │   ├── auth.store.ts             # User, tokens, login/logout actions
    │   │   ├── cart.store.ts             # Cart items, add/remove/update, totals
    │   │   └── ui.store.ts              # Sidebar open, theme, modals
    │   │
    │   ├── lib/
    │   │   ├── api-client.ts             # Axios instance + interceptors (JWT refresh, errors)
    │   │   ├── query-client.ts           # TanStack Query client config
    │   │   ├── utils.ts                  # cn(), formatPrice(), formatDate()
    │   │   ├── validators.ts             # Shared Zod schemas
    │   │   └── constants.ts              # API_URL, ORDER_STATUSES, ROLES
    │   │
    │   ├── types/
    │   │   ├── api.types.ts              # ApiResponse<T>, PaginatedResponse<T>
    │   │   ├── user.types.ts
    │   │   ├── product.types.ts
    │   │   ├── order.types.ts
    │   │   ├── cart.types.ts
    │   │   └── common.types.ts
    │   │
    │   ├── guards/
    │   │   ├── AuthGuard.tsx             # Redirect to /login if unauthenticated
    │   │   ├── GuestGuard.tsx            # Redirect to /shop if already authenticated
    │   │   ├── RoleGuard.tsx             # Check user role (seller, admin)
    │   │   └── SellerGuard.tsx
    │   │
    │   └── providers/
    │       ├── AppProviders.tsx           # Compose all providers
    │       ├── QueryProvider.tsx          # TanStack QueryClientProvider
    │       ├── ThemeProvider.tsx          # Dark/light mode
    │       ├── ToastProvider.tsx
    │       └── SocketProvider.tsx         # Socket.io context
    │
    ├── modules/                          # ═══ FEATURE MODULES ═══
    │   │
    │   ├── auth/                         # ── Auth Module ──
    │   │   ├── index.ts                  # Barrel export
    │   │   ├── routes.tsx                # /login, /register, /forgot-password
    │   │   ├── components/
    │   │   │   ├── LoginForm.tsx
    │   │   │   ├── RegisterForm.tsx
    │   │   │   ├── ForgotPasswordForm.tsx
    │   │   │   ├── ResetPasswordForm.tsx
    │   │   │   ├── OAuthButtons.tsx       # Google, Facebook, Apple login
    │   │   │   └── AuthLayout.tsx         # Centered card layout for auth pages
    │   │   ├── pages/
    │   │   │   ├── LoginPage.tsx
    │   │   │   ├── RegisterPage.tsx
    │   │   │   └── ForgotPasswordPage.tsx
    │   │   ├── services/
    │   │   │   └── auth.api.ts            # login(), register(), refreshToken(), etc.
    │   │   ├── hooks/
    │   │   │   ├── useLogin.ts            # useMutation wrapper
    │   │   │   └── useRegister.ts
    │   │   └── types/
    │   │       └── auth.types.ts
    │   │
    │   ├── shop/                         # ── Shop Module (buyer browsing) ──
    │   │   ├── index.ts
    │   │   ├── routes.tsx                # /products, /products/:slug, /categories/:slug
    │   │   ├── components/
    │   │   │   ├── ProductCard.tsx
    │   │   │   ├── ProductGrid.tsx
    │   │   │   ├── ProductGallery.tsx     # Image carousel/zoom
    │   │   │   ├── ProductInfo.tsx        # Price, rating, add to cart
    │   │   │   ├── ProductVariantSelector.tsx  # Option buttons/dropdowns (Color swatches, Size pills)
    │   │   │   ├── ProductAttributes.tsx      # Display filled-in attributes (Brand, Material, etc.)
    │   │   │   ├── ProductReviews.tsx
    │   │   │   ├── CategoryNav.tsx            # Category tree sidebar
    │   │   │   ├── FilterPanel.tsx            # Price range, rating, dynamic category attributes
    │   │   │   ├── AttributeFilter.tsx        # Renders filter controls per attribute type (checkbox, range, swatch)
    │   │   │   ├── SortDropdown.tsx
    │   │   │   └── BreadcrumbNav.tsx
    │   │   ├── pages/
    │   │   │   ├── ProductListPage.tsx    # Grid + filters + pagination
    │   │   │   ├── ProductDetailPage.tsx  # Gallery + info + reviews + recommendations
    │   │   │   ├── CategoryPage.tsx
    │   │   │   └── HomePage.tsx           # Hero, featured, trending, categories
    │   │   ├── services/
    │   │   │   ├── product.api.ts         # getProducts(), getProductBySlug(), etc.
    │   │   │   └── category.api.ts
    │   │   ├── hooks/
    │   │   │   ├── useProducts.ts         # useQuery: product list with filters
    │   │   │   ├── useProduct.ts          # useQuery: single product
    │   │   │   └── useCategories.ts
    │   │   └── types/
    │   │       └── shop.types.ts          # FilterState, SortOption
    │   │
    │   ├── search/                       # ── Search Module ──
    │   │   ├── index.ts
    │   │   ├── routes.tsx                # /search?q=
    │   │   ├── components/
    │   │   │   ├── SearchBar.tsx          # Autocomplete + debounce
    │   │   │   ├── SearchResults.tsx
    │   │   │   ├── SearchFilters.tsx
    │   │   │   ├── ImageSearchUpload.tsx  # AI visual search
    │   │   │   └── SearchSuggestions.tsx
    │   │   ├── pages/
    │   │   │   └── SearchResultsPage.tsx
    │   │   ├── services/
    │   │   │   └── search.api.ts
    │   │   └── hooks/
    │   │       ├── useSearch.ts
    │   │       └── useSearchSuggestions.ts
    │   │
    │   ├── cart/                          # ── Cart Module ──
    │   │   ├── index.ts
    │   │   ├── routes.tsx                # /cart
    │   │   ├── components/
    │   │   │   ├── CartDrawer.tsx         # Slide-over cart preview
    │   │   │   ├── CartItem.tsx
    │   │   │   ├── CartSummary.tsx        # Subtotal, shipping, tax, total
    │   │   │   ├── CartEmpty.tsx
    │   │   │   └── QuantitySelector.tsx
    │   │   ├── pages/
    │   │   │   └── CartPage.tsx           # Full cart view
    │   │   ├── services/
    │   │   │   └── cart.api.ts            # Syncs Zustand cart ↔ API
    │   │   └── hooks/
    │   │       ├── useCart.ts             # Wraps Zustand store + API sync
    │   │       └── useCartSync.ts         # Sync local cart to server on login
    │   │
    │   ├── checkout/                     # ── Checkout Module ──
    │   │   ├── index.ts
    │   │   ├── routes.tsx                # /checkout, /checkout/success
    │   │   ├── components/
    │   │   │   ├── CheckoutStepper.tsx    # Steps: address → payment → review
    │   │   │   ├── ShippingForm.tsx
    │   │   │   ├── PaymentForm.tsx        # Stripe Elements integration
    │   │   │   ├── OrderReview.tsx
    │   │   │   ├── CouponInput.tsx
    │   │   │   └── CheckoutSummary.tsx
    │   │   ├── pages/
    │   │   │   ├── CheckoutPage.tsx
    │   │   │   └── CheckoutSuccessPage.tsx
    │   │   ├── services/
    │   │   │   ├── order.api.ts           # createOrder()
    │   │   │   └── payment.api.ts         # createPaymentIntent()
    │   │   └── hooks/
    │   │       └── useCheckout.ts         # Multi-step form state
    │   │
    │   ├── account/                      # ── Customer Account Module ──
    │   │   ├── index.ts
    │   │   ├── routes.tsx                # /account/profile, /account/orders, etc.
    │   │   ├── components/
    │   │   │   ├── AccountSidebar.tsx
    │   │   │   ├── ProfileForm.tsx
    │   │   │   ├── AddressBook.tsx
    │   │   │   ├── AddressCard.tsx
    │   │   │   ├── OrderList.tsx
    │   │   │   ├── OrderDetail.tsx
    │   │   │   ├── OrderTimeline.tsx       # Status progression visual
    │   │   │   ├── WishlistGrid.tsx
    │   │   │   └── ChangePasswordForm.tsx
    │   │   ├── pages/
    │   │   │   ├── ProfilePage.tsx
    │   │   │   ├── OrdersPage.tsx
    │   │   │   ├── OrderDetailPage.tsx
    │   │   │   ├── AddressesPage.tsx
    │   │   │   └── WishlistPage.tsx
    │   │   ├── services/
    │   │   │   ├── profile.api.ts
    │   │   │   ├── address.api.ts
    │   │   │   └── order.api.ts
    │   │   └── hooks/
    │   │       ├── useProfile.ts
    │   │       ├── useOrders.ts
    │   │       └── useAddresses.ts
    │   │
    │   ├── seller/                       # ── Seller Dashboard Module ──
    │   │   ├── index.ts
    │   │   ├── routes.tsx                # /seller/dashboard, /seller/products, etc.
    │   │   ├── components/
    │   │   │   ├── SellerSidebar.tsx
    │   │   │   ├── StatsCard.tsx
    │   │   │   ├── RevenueChart.tsx       # Recharts/Chart.js
    │   │   │   ├── ProductForm.tsx        # Create/edit product form (basic info + category selection)
    │   │   │   ├── AttributeForm.tsx     # Dynamic form fields based on category attributes
    │   │   │   ├── OptionManager.tsx     # Define variant axes (Color, Size) + values
    │   │   │   ├── VariantManager.tsx    # Edit generated variants (SKU, price, stock, images per variant)
    │   │   │   ├── VariantTable.tsx      # Bulk edit variant grid (inline price/stock/SKU editing)
    │   │   │   ├── ImageManager.tsx      # Drag-drop image upload + reorder + assign to variants
    │   │   │   ├── SellerProductTable.tsx
    │   │   │   ├── SellerOrderTable.tsx
    │   │   │   ├── PayoutHistory.tsx
    │   │   │   └── StoreSettingsForm.tsx
    │   │   ├── pages/
    │   │   │   ├── SellerDashboardPage.tsx
    │   │   │   ├── SellerProductsPage.tsx
    │   │   │   ├── SellerProductCreatePage.tsx
    │   │   │   ├── SellerProductEditPage.tsx
    │   │   │   ├── SellerOrdersPage.tsx
    │   │   │   ├── SellerAnalyticsPage.tsx
    │   │   │   ├── SellerPayoutsPage.tsx
    │   │   │   └── SellerSettingsPage.tsx
    │   │   ├── services/
    │   │   │   ├── seller-product.api.ts
    │   │   │   ├── seller-order.api.ts
    │   │   │   ├── seller-analytics.api.ts
    │   │   │   └── seller-payout.api.ts
    │   │   └── hooks/
    │   │       ├── useSellerProducts.ts
    │   │       ├── useSellerOrders.ts
    │   │       └── useSellerStats.ts
    │   │
    │   ├── admin/                        # ── Admin Module ──
    │   │   ├── index.ts
    │   │   ├── routes.tsx                # /admin/dashboard, /admin/users, etc.
    │   │   ├── components/
    │   │   │   ├── AdminSidebar.tsx
    │   │   │   ├── DashboardStats.tsx
    │   │   │   ├── UserManagementTable.tsx
    │   │   │   ├── SellerApprovalTable.tsx
    │   │   │   ├── OrderManagementTable.tsx
    │   │   │   ├── ProductModerationTable.tsx
    │   │   │   ├── AttributeDefinitionForm.tsx  # Create/edit attribute definitions
    │   │   │   ├── AttributeDefinitionTable.tsx # List all attributes
    │   │   │   ├── CategoryAttributeManager.tsx # Assign/reorder attributes per category
    │   │   │   ├── RevenueReportChart.tsx
    │   │   │   └── SystemSettingsForm.tsx
    │   │   ├── pages/
    │   │   │   ├── AdminDashboardPage.tsx
    │   │   │   ├── AdminUsersPage.tsx
    │   │   │   ├── AdminSellersPage.tsx
    │   │   │   ├── AdminProductsPage.tsx
    │   │   │   ├── AdminAttributesPage.tsx      # Manage attribute definitions + category assignments
    │   │   │   ├── AdminCategoriesPage.tsx       # Category CRUD + attribute assignment
    │   │   │   ├── AdminOrdersPage.tsx
    │   │   │   ├── AdminReportsPage.tsx
    │   │   │   └── AdminSettingsPage.tsx
    │   │   ├── services/
    │   │   │   ├── admin-user.api.ts
    │   │   │   ├── admin-seller.api.ts
    │   │   │   ├── admin-attribute.api.ts       # CRUD attribute definitions + category assignment
    │   │   │   └── admin-report.api.ts
    │   │   └── hooks/
    │   │       ├── useAdminUsers.ts
    │   │       ├── useAdminAttributes.ts        # useQuery/useMutation for attribute definitions
    │   │       └── useAdminDashboard.ts
    │   │
    │   ├── reviews/                      # ── Reviews Module ──
    │   │   ├── index.ts
    │   │   ├── components/
    │   │   │   ├── ReviewCard.tsx
    │   │   │   ├── ReviewForm.tsx
    │   │   │   ├── ReviewList.tsx
    │   │   │   └── ReviewSummary.tsx       # Rating distribution bar chart
    │   │   ├── services/
    │   │   │   └── review.api.ts
    │   │   └── hooks/
    │   │       └── useReviews.ts
    │   │
    │   ├── notifications/                # ── Notifications Module ──
    │   │   ├── index.ts
    │   │   ├── components/
    │   │   │   ├── NotificationBell.tsx   # Header bell icon + badge count
    │   │   │   ├── NotificationDropdown.tsx
    │   │   │   ├── NotificationItem.tsx
    │   │   │   └── NotificationList.tsx
    │   │   ├── pages/
    │   │   │   └── NotificationsPage.tsx
    │   │   ├── services/
    │   │   │   └── notification.api.ts
    │   │   └── hooks/
    │   │       └── useNotifications.ts    # useQuery + WebSocket subscription
    │   │
    │   ├── chat/                         # ── Chat Module ──
    │   │   ├── index.ts
    │   │   ├── routes.tsx                # /messages, /messages/:conversationId
    │   │   ├── components/
    │   │   │   ├── ChatWidget.tsx         # Floating chat bubble
    │   │   │   ├── ConversationList.tsx
    │   │   │   ├── ChatWindow.tsx
    │   │   │   ├── MessageBubble.tsx
    │   │   │   └── ChatInput.tsx
    │   │   ├── pages/
    │   │   │   └── MessagesPage.tsx
    │   │   ├── services/
    │   │   │   └── chat.api.ts
    │   │   └── hooks/
    │   │       ├── useConversations.ts
    │   │       └── useMessages.ts         # useQuery + WebSocket for real-time
    │   │
    │   ├── ai/                           # ── AI Features Module ──
    │   │   ├── index.ts
    │   │   ├── components/
    │   │   │   ├── AIChatAssistant.tsx    # Floating AI shopping assistant
    │   │   │   ├── AIChatWindow.tsx
    │   │   │   ├── AIChatMessage.tsx
    │   │   │   ├── RecommendationCarousel.tsx  # "You might also like"
    │   │   │   ├── SimilarProducts.tsx
    │   │   │   └── AIDescriptionGenerator.tsx  # Seller tool: auto-generate desc
    │   │   ├── services/
    │   │   │   └── ai.api.ts              # chat(), getRecommendations(), generateDesc()
    │   │   └── hooks/
    │   │       ├── useAIChat.ts
    │   │       └── useRecommendations.ts
    │   │
    │   ├── promotions/                   # ── Promotions Module ──
    │   │   ├── index.ts
    │   │   ├── components/
    │   │   │   ├── CouponInput.tsx        # Apply coupon code in cart/checkout
    │   │   │   ├── CouponBadge.tsx        # Show applied discount
    │   │   │   ├── FlashSaleBanner.tsx    # Countdown timer + deals
    │   │   │   ├── FlashSaleGrid.tsx      # Flash sale product listing
    │   │   │   ├── BundleDealCard.tsx     # Bundle offer display
    │   │   │   ├── SellerCouponForm.tsx   # Seller: create/edit coupons
    │   │   │   ├── SellerCouponTable.tsx  # Seller: list coupons
    │   │   │   ├── AdminFlashSaleForm.tsx # Admin: create flash sales
    │   │   │   └── AdminPromotionTable.tsx
    │   │   ├── pages/
    │   │   │   ├── FlashSalesPage.tsx     # Public flash sale listing
    │   │   │   ├── SellerCouponsPage.tsx
    │   │   │   ├── AdminPromotionsPage.tsx
    │   │   │   └── AdminFlashSalesPage.tsx
    │   │   ├── services/
    │   │   │   ├── coupon.api.ts          # validate, list, redeem
    │   │   │   ├── flash-sale.api.ts
    │   │   │   └── bundle.api.ts
    │   │   └── hooks/
    │   │       ├── useCoupons.ts
    │   │       ├── useFlashSales.ts
    │   │       └── useBundles.ts
    │   │
    │   ├── returns/                      # ── Returns & Disputes Module ──
    │   │   ├── index.ts
    │   │   ├── components/
    │   │   │   ├── ReturnRequestForm.tsx  # Buyer: request return with reason + photos
    │   │   │   ├── ReturnStatusTimeline.tsx # Visual status progression
    │   │   │   ├── ReturnItemSelector.tsx # Select which items to return
    │   │   │   ├── DisputeForm.tsx        # Open a dispute
    │   │   │   ├── DisputeChat.tsx        # Dispute message thread
    │   │   │   ├── SellerReturnTable.tsx  # Seller: manage return requests
    │   │   │   ├── SellerReturnActions.tsx # Approve/reject with notes
    │   │   │   ├── AdminDisputeTable.tsx  # Admin: all disputes
    │   │   │   └── AdminDisputeDetail.tsx # Admin: review + resolve dispute
    │   │   ├── pages/
    │   │   │   ├── BuyerReturnsPage.tsx
    │   │   │   ├── BuyerReturnDetailPage.tsx
    │   │   │   ├── CreateReturnPage.tsx
    │   │   │   ├── SellerReturnsPage.tsx
    │   │   │   ├── AdminDisputesPage.tsx
    │   │   │   └── AdminDisputeDetailPage.tsx
    │   │   ├── services/
    │   │   │   ├── return.api.ts
    │   │   │   └── dispute.api.ts
    │   │   └── hooks/
    │   │       ├── useReturns.ts
    │   │       └── useDisputes.ts
    │   │
    │   ├── shipping/                     # ── Shipping Module ──
    │   │   ├── index.ts
    │   │   ├── components/
    │   │   │   ├── ShippingRateSelector.tsx  # Rate shopping at checkout
    │   │   │   ├── TrackingTimeline.tsx      # Visual tracking events
    │   │   │   ├── TrackingPage.tsx          # Public tracking lookup
    │   │   │   ├── SellerShipmentForm.tsx    # Create shipment + generate label
    │   │   │   ├── SellerShipmentTable.tsx
    │   │   │   ├── LabelPreview.tsx          # PDF label preview + print
    │   │   │   ├── CarrierSetupForm.tsx      # Seller: configure carrier API keys
    │   │   │   └── AdminCarrierTable.tsx     # Admin: manage carrier list
    │   │   ├── pages/
    │   │   │   ├── TrackingLookupPage.tsx
    │   │   │   ├── SellerShipmentsPage.tsx
    │   │   │   ├── SellerCarrierSetupPage.tsx
    │   │   │   └── AdminCarriersPage.tsx
    │   │   ├── services/
    │   │   │   └── shipping.api.ts
    │   │   └── hooks/
    │   │       ├── useShippingRates.ts
    │   │       ├── useShipments.ts
    │   │       └── useTracking.ts
    │   │
    │   ├── loyalty/                      # ── Loyalty & Rewards Module ──
    │   │   ├── index.ts
    │   │   ├── components/
    │   │   │   ├── PointsBalance.tsx      # Current points + tier display
    │   │   │   ├── TierProgressBar.tsx    # Progress to next tier
    │   │   │   ├── TierBenefitsCard.tsx   # What each tier offers
    │   │   │   ├── PointsHistory.tsx      # Transaction list
    │   │   │   ├── RedeemPointsInput.tsx  # Use points at checkout
    │   │   │   └── EarnPointsBanner.tsx   # Motivational: "Earn X points on this order"
    │   │   ├── pages/
    │   │   │   ├── LoyaltyDashboardPage.tsx
    │   │   │   └── TiersInfoPage.tsx
    │   │   ├── services/
    │   │   │   └── loyalty.api.ts
    │   │   └── hooks/
    │   │       ├── useMembership.ts
    │   │       ├── usePointsBalance.ts
    │   │       └── usePointsHistory.ts
    │   │
    │   ├── affiliate/                    # ── Affiliate & Referral Module ──
    │   │   ├── index.ts
    │   │   ├── components/
    │   │   │   ├── ReferralLinkGenerator.tsx  # Generate + copy referral link
    │   │   │   ├── ReferralStats.tsx          # Clicks, conversions, earnings
    │   │   │   ├── ReferralTable.tsx          # List all referrals
    │   │   │   ├── PayoutRequestForm.tsx      # Request payout
    │   │   │   ├── ShareButtons.tsx           # Social share (Facebook, Twitter, WhatsApp, Email)
    │   │   │   └── ReferralBanner.tsx         # "Invite friends, earn $X"
    │   │   ├── pages/
    │   │   │   ├── AffiliateDashboardPage.tsx
    │   │   │   ├── AffiliatePayoutsPage.tsx
    │   │   │   └── AdminAffiliatesPage.tsx
    │   │   ├── services/
    │   │   │   └── affiliate.api.ts
    │   │   └── hooks/
    │   │       ├── useAffiliateStats.ts
    │   │       └── useReferrals.ts
    │   │
    │   └── cms/                          # ── CMS / Content Module ──
    │       ├── index.ts
    │       ├── components/
    │       │   ├── BannerCarousel.tsx      # Hero banners on homepage
    │       │   ├── PromoBanner.tsx         # Sidebar / category page banners
    │       │   ├── StaticPage.tsx          # Render CMS page content
    │       │   ├── AdminBannerForm.tsx     # Admin: create/edit banners
    │       │   ├── AdminBannerTable.tsx
    │       │   ├── AdminPageEditor.tsx     # Rich text page editor
    │       │   ├── AdminPageTable.tsx
    │       │   └── ContentScheduler.tsx   # Schedule publish/unpublish
    │       ├── pages/
    │       │   ├── StaticPageView.tsx      # Public: /pages/:slug
    │       │   ├── AdminBannersPage.tsx
    │       │   └── AdminPagesPage.tsx
    │       ├── services/
    │       │   └── cms.api.ts
    │       └── hooks/
    │           ├── useBanners.ts
    │           └── usePages.ts
    │
    └── styles/
        └── globals.css                    # Tailwind @import + custom CSS variables
```

### Module Rules

1. **Self-contained**: Each module has its own `components/`, `pages/`, `services/`, `hooks/`, `types/`
2. **Barrel exports**: Each module exposes only what's needed via `index.ts`
3. **No cross-module imports** except through `shared/` — if two modules need the same component, move it to `shared/`
4. **Lazy loading**: Each module's `routes.tsx` uses `React.lazy()` for code splitting
5. **API isolation**: Each module's `services/*.api.ts` handles its own API calls via the shared Axios instance

### Route Configuration (`src/routes.tsx`)

```tsx
import { lazy, Suspense } from 'react';
import { createBrowserRouter, Navigate } from 'react-router-dom';
import { RootLayout } from '@/shared/components/layout/RootLayout';
import { DashboardLayout } from '@/shared/components/layout/DashboardLayout';
import { AuthGuard } from '@/shared/guards/AuthGuard';
import { GuestGuard } from '@/shared/guards/GuestGuard';
import { RoleGuard } from '@/shared/guards/RoleGuard';
import { LoadingSpinner } from '@/shared/components/feedback/LoadingSpinner';

// Lazy-loaded module pages (code-split per module)
const HomePage = lazy(() => import('@/modules/shop/pages/HomePage'));
const ProductListPage = lazy(() => import('@/modules/shop/pages/ProductListPage'));
const ProductDetailPage = lazy(() => import('@/modules/shop/pages/ProductDetailPage'));
const SearchResultsPage = lazy(() => import('@/modules/search/pages/SearchResultsPage'));
const CartPage = lazy(() => import('@/modules/cart/pages/CartPage'));
const CheckoutPage = lazy(() => import('@/modules/checkout/pages/CheckoutPage'));
const LoginPage = lazy(() => import('@/modules/auth/pages/LoginPage'));
const RegisterPage = lazy(() => import('@/modules/auth/pages/RegisterPage'));
// ... all other pages

const Lazy = ({ children }: { children: React.ReactNode }) => (
  <Suspense fallback={<LoadingSpinner />}>{children}</Suspense>
);

export const router = createBrowserRouter([
  {
    element: <RootLayout />,   // Header + Footer shell
    children: [
      // ── Public (Shop) ──
      { path: '/', element: <Lazy><HomePage /></Lazy> },
      { path: '/products', element: <Lazy><ProductListPage /></Lazy> },
      { path: '/products/:slug', element: <Lazy><ProductDetailPage /></Lazy> },
      { path: '/categories/:slug', element: <Lazy><CategoryPage /></Lazy> },
      { path: '/search', element: <Lazy><SearchResultsPage /></Lazy> },
      { path: '/cart', element: <Lazy><CartPage /></Lazy> },
      { path: '/flash-sales', element: <Lazy><FlashSalesPage /></Lazy> },
      { path: '/tracking/:number', element: <Lazy><TrackingLookupPage /></Lazy> },
      { path: '/pages/:slug', element: <Lazy><StaticPageView /></Lazy> },
      { path: '/r/:code', element: <Lazy><ReferralLandingPage /></Lazy> },

      // ── Auth (guest only) ──
      { element: <GuestGuard />, children: [
        { path: '/login', element: <Lazy><LoginPage /></Lazy> },
        { path: '/register', element: <Lazy><RegisterPage /></Lazy> },
      ]},

      // ── Protected (any authenticated user) ──
      { element: <AuthGuard />, children: [
        { path: '/checkout', element: <Lazy><CheckoutPage /></Lazy> },
        { path: '/checkout/success', element: <Lazy><CheckoutSuccessPage /></Lazy> },

        // ── Account (buyer) ──
        { path: '/account', children: [
          { index: true, element: <Navigate to="profile" /> },
          { path: 'profile', element: <Lazy><ProfilePage /></Lazy> },
          { path: 'orders', element: <Lazy><OrdersPage /></Lazy> },
          { path: 'orders/:id', element: <Lazy><OrderDetailPage /></Lazy> },
          { path: 'addresses', element: <Lazy><AddressesPage /></Lazy> },
          { path: 'wishlist', element: <Lazy><WishlistPage /></Lazy> },
          { path: 'returns', element: <Lazy><BuyerReturnsPage /></Lazy> },
          { path: 'returns/new/:orderId', element: <Lazy><CreateReturnPage /></Lazy> },
          { path: 'returns/:id', element: <Lazy><BuyerReturnDetailPage /></Lazy> },
        ]},

        // ── Loyalty ──
        { path: '/loyalty', element: <Lazy><LoyaltyDashboardPage /></Lazy> },
        { path: '/loyalty/tiers', element: <Lazy><TiersInfoPage /></Lazy> },

        // ── Affiliate ──
        { path: '/affiliate', element: <Lazy><AffiliateDashboardPage /></Lazy> },
        { path: '/affiliate/payouts', element: <Lazy><AffiliatePayoutsPage /></Lazy> },

        // ── Messages ──
        { path: '/messages', element: <Lazy><MessagesPage /></Lazy> },
      ]},

      // ── Seller (role: seller) ──
      { element: <RoleGuard roles={['seller', 'admin']} />, children: [
        { path: '/seller', element: <DashboardLayout sidebar="seller" />, children: [
          { index: true, element: <Navigate to="dashboard" /> },
          { path: 'dashboard', element: <Lazy><SellerDashboardPage /></Lazy> },
          { path: 'products', element: <Lazy><SellerProductsPage /></Lazy> },
          { path: 'products/new', element: <Lazy><SellerProductCreatePage /></Lazy> },
          { path: 'products/:id/edit', element: <Lazy><SellerProductEditPage /></Lazy> },
          { path: 'orders', element: <Lazy><SellerOrdersPage /></Lazy> },
          { path: 'returns', element: <Lazy><SellerReturnsPage /></Lazy> },
          { path: 'shipments', element: <Lazy><SellerShipmentsPage /></Lazy> },
          { path: 'carriers', element: <Lazy><SellerCarrierSetupPage /></Lazy> },
          { path: 'coupons', element: <Lazy><SellerCouponsPage /></Lazy> },
          { path: 'analytics', element: <Lazy><SellerAnalyticsPage /></Lazy> },
          { path: 'payouts', element: <Lazy><SellerPayoutsPage /></Lazy> },
          { path: 'settings', element: <Lazy><SellerSettingsPage /></Lazy> },
        ]},
      ]},

      // ── Admin (role: admin) ──
      { element: <RoleGuard roles={['admin']} />, children: [
        { path: '/admin', element: <DashboardLayout sidebar="admin" />, children: [
          { index: true, element: <Navigate to="dashboard" /> },
          { path: 'dashboard', element: <Lazy><AdminDashboardPage /></Lazy> },
          { path: 'users', element: <Lazy><AdminUsersPage /></Lazy> },
          { path: 'sellers', element: <Lazy><AdminSellersPage /></Lazy> },
          { path: 'products', element: <Lazy><AdminProductsPage /></Lazy> },
          { path: 'attributes', element: <Lazy><AdminAttributesPage /></Lazy> },
          { path: 'categories', element: <Lazy><AdminCategoriesPage /></Lazy> },
          { path: 'orders', element: <Lazy><AdminOrdersPage /></Lazy> },
          { path: 'returns', element: <Lazy><AdminDisputesPage /></Lazy> },
          { path: 'returns/:id', element: <Lazy><AdminDisputeDetailPage /></Lazy> },
          { path: 'promotions', element: <Lazy><AdminPromotionsPage /></Lazy> },
          { path: 'flash-sales', element: <Lazy><AdminFlashSalesPage /></Lazy> },
          { path: 'carriers', element: <Lazy><AdminCarriersPage /></Lazy> },
          { path: 'banners', element: <Lazy><AdminBannersPage /></Lazy> },
          { path: 'pages', element: <Lazy><AdminPagesPage /></Lazy> },
          { path: 'affiliates', element: <Lazy><AdminAffiliatesPage /></Lazy> },
          { path: 'tax', element: <Lazy><AdminTaxRulesPage /></Lazy> },
          { path: 'reports', element: <Lazy><AdminReportsPage /></Lazy> },
          { path: 'settings', element: <Lazy><AdminSettingsPage /></Lazy> },
        ]},
      ]},
    ],
  },
]);
```

### API Client Pattern (`shared/lib/api-client.ts`)

```typescript
import axios from 'axios';
import { useAuthStore } from '@/shared/stores/auth.store';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000/api/v1';

export const apiClient = axios.create({
  baseURL: API_URL,
  headers: { 'Content-Type': 'application/json' },
});

// Attach JWT to every request
apiClient.interceptors.request.use((config) => {
  const token = useAuthStore.getState().accessToken;
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

// Auto-refresh on 401
apiClient.interceptors.response.use(
  (res) => res,
  async (error) => {
    if (error.response?.status === 401 && !error.config._retry) {
      error.config._retry = true;
      try {
        const newToken = await useAuthStore.getState().refreshToken();
        error.config.headers.Authorization = `Bearer ${newToken}`;
        return apiClient(error.config);
      } catch {
        useAuthStore.getState().logout();
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);
```

### Module API Service Pattern (`modules/shop/services/product.api.ts`)

```typescript
import { apiClient } from '@/shared/lib/api-client';
import type { Product, ProductFilter, PaginatedResponse } from '@/shared/types';

export const productApi = {
  list: (params: ProductFilter) =>
    apiClient.get<PaginatedResponse<Product>>('/products', { params }),

  getBySlug: (slug: string) =>
    apiClient.get<Product>(`/products/${slug}`),

  getByCategory: (slug: string, params: ProductFilter) =>
    apiClient.get<PaginatedResponse<Product>>(`/categories/${slug}/products`, { params }),
};
```

### Module Hook Pattern (`modules/shop/hooks/useProducts.ts`)

```typescript
import { useQuery } from '@tanstack/react-query';
import { productApi } from '../services/product.api';
import type { ProductFilter } from '@/shared/types';

export function useProducts(filters: ProductFilter) {
  return useQuery({
    queryKey: ['products', filters],
    queryFn: () => productApi.list(filters).then(res => res.data),
    staleTime: 30_000,        // 30s before refetch
    placeholderData: (prev) => prev,  // Keep previous data while loading
  });
}

export function useProduct(slug: string) {
  return useQuery({
    queryKey: ['product', slug],
    queryFn: () => productApi.getBySlug(slug).then(res => res.data),
    enabled: !!slug,
  });
}
```

### Product TypeScript Types (`shared/types/product.types.ts`)

```typescript
// --- Attribute system (admin-defined) ---
export type AttributeType = 'text' | 'number' | 'select' | 'multi_select' | 'color' | 'bool';

export interface AttributeDefinition {
  id: string;
  name: string;
  slug: string;
  type: AttributeType;
  required: boolean;
  filterable: boolean;
  options: string[];        // predefined choices for select types
  unit?: string;            // "cm", "kg"
  sortOrder: number;
}

export interface ProductAttributeValue {
  attributeId: string;
  attributeName: string;
  value: string;
  values?: string[];        // for multi_select
}

// --- Option / Variant system (seller-defined) ---
export interface ProductOption {
  id: string;
  name: string;             // "Color", "Size"
  sortOrder: number;
  values: ProductOptionValue[];
}

export interface ProductOptionValue {
  id: string;
  value: string;            // "Red", "XL"
  colorHex?: string;        // "#FF0000" for color swatches
  sortOrder: number;
}

export interface Variant {
  id: string;
  productId: string;
  sku: string;
  name: string;             // "Red / XL"
  priceCents: number;       // 0 = use base price
  compareAtCents: number;
  costCents: number;
  stock: number;
  lowStockAlert: number;
  weightGrams: number;
  isDefault: boolean;
  isActive: boolean;
  imageUrls: string[];
  barcode?: string;
  optionValues: VariantOptionValue[];
}

export interface VariantOptionValue {
  optionId: string;
  optionValueId: string;
  optionName: string;       // "Color"
  value: string;            // "Red"
}

// --- Product (with all nested data) ---
export interface Product {
  id: string;
  sellerId: string;
  categoryId: string;
  name: string;
  slug: string;
  description: string;
  basePriceCents: number;
  currency: string;
  status: 'draft' | 'active' | 'inactive' | 'archived';
  hasVariants: boolean;
  tags: string[];
  imageUrls: string[];
  attributeValues: ProductAttributeValue[];
  options: ProductOption[];
  variants: Variant[];
  ratingAvg: number;
  ratingCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface ProductFilter {
  categoryId?: string;
  sellerId?: string;
  status?: string;
  query?: string;
  minPrice?: number;
  maxPrice?: number;
  attributes?: Record<string, string[]>;  // attributeSlug → selected values
  sortBy?: string;
  page?: number;
  pageSize?: number;
}
```

### Zustand Cart Store Pattern (`shared/stores/cart.store.ts`)

```typescript
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface CartItem {
  productId: string;
  variantId?: string;       // required when product has variants
  variantName?: string;     // "Red / XL" — display label
  sku?: string;
  name: string;             // product name
  price: number;            // effective price (variant price or base price)
  image: string;            // variant image or product image
  quantity: number;
}

interface CartStore {
  items: CartItem[];
  addItem: (item: CartItem) => void;
  removeItem: (productId: string, variantId?: string) => void;
  updateQuantity: (productId: string, quantity: number, variantId?: string) => void;
  clearCart: () => void;
  totalItems: () => number;
  totalPrice: () => number;
}

export const useCartStore = create<CartStore>()(
  persist(
    (set, get) => ({
      items: [],
      addItem: (item) => set((state) => {
        const existing = state.items.find(
          i => i.productId === item.productId && i.variantId === item.variantId
        );
        if (existing) {
          return { items: state.items.map(i =>
            i.productId === item.productId && i.variantId === item.variantId
              ? { ...i, quantity: i.quantity + item.quantity }
              : i
          )};
        }
        return { items: [...state.items, item] };
      }),
      removeItem: (productId, variantId) => set((state) => ({
        items: state.items.filter(i => !(i.productId === productId && i.variantId === variantId))
      })),
      updateQuantity: (productId, quantity, variantId) => set((state) => ({
        items: state.items.map(i =>
          i.productId === productId && i.variantId === variantId ? { ...i, quantity } : i
        )
      })),
      clearCart: () => set({ items: [] }),
      totalItems: () => get().items.reduce((sum, i) => sum + i.quantity, 0),
      totalPrice: () => get().items.reduce((sum, i) => sum + i.price * i.quantity, 0),
    }),
    { name: 'cart-storage' }   // Persist to localStorage
  )
);
```

### SEO Consideration (SPA without SSR)

Since this is a Vite SPA (no SSR), SEO is handled via:
1. **Prerendering** — Use `vite-plugin-ssr` or `react-snap` to pre-render critical pages (homepage, product pages, categories) at build time
2. **Dynamic rendering** — Serve pre-rendered HTML to search engine bots via a service like Rendertron or prerender.io
3. **Meta tags** — React Helmet Async for dynamic `<title>`, `<meta description>`, Open Graph tags per page
4. **Sitemap** — Generate sitemap.xml at build time from product/category slugs

---

## Part 18: Mobile — Flutter Clean Architecture (Two Apps + Shared Packages)

### Tech Stack
- **Flutter 3.x** + Dart
- **State Management**: BLoC / Cubit (flutter_bloc)
- **Dependency Injection**: get_it + injectable (code generation)
- **HTTP Client**: Dio (interceptors, JWT refresh, error handling)
- **Local Storage**: flutter_secure_storage (tokens), Hive (cache)
- **Navigation**: go_router (declarative routing)
- **Forms**: reactive_forms or flutter_form_builder
- **Image**: cached_network_image, image_picker, image_cropper
- **Payments**: flutter_stripe
- **Real-time**: socket_io_client
- **Code Gen**: freezed (immutable models), json_serializable, injectable_generator, build_runner

### Two Apps + Shared Packages Structure

```
apps/mobile/
├── buyer_app/                            # ═══ BUYER APP ═══
│   ├── pubspec.yaml
│   ├── lib/
│   │   ├── main.dart                     # Entry point, DI init, runApp
│   │   ├── app.dart                      # MaterialApp, go_router, BlocProviders
│   │   ├── injection.dart                # get_it setup, @InjectableInit
│   │   ├── injection.config.dart         # Generated by injectable
│   │   │
│   │   ├── core/                         # App-level config
│   │   │   ├── config/
│   │   │   │   ├── app_config.dart       # API_URL, env-specific config
│   │   │   │   └── theme.dart            # ThemeData, colors, typography
│   │   │   ├── router/
│   │   │   │   └── app_router.dart       # go_router route definitions
│   │   │   └── constants/
│   │   │       └── app_constants.dart
│   │   │
│   │   └── features/                     # ═══ FEATURE MODULES ═══
│   │       │
│   │       ├── auth/                     # ── Auth Feature ──
│   │       │   ├── domain/
│   │       │   │   ├── entities/
│   │       │   │   │   └── user.dart                 # Pure domain entity
│   │       │   │   ├── repositories/
│   │       │   │   │   └── auth_repository.dart      # Abstract repository interface
│   │       │   │   └── usecases/
│   │       │   │       ├── login_usecase.dart
│   │       │   │       ├── register_usecase.dart
│   │       │   │       ├── logout_usecase.dart
│   │       │   │       └── refresh_token_usecase.dart
│   │       │   ├── data/
│   │       │   │   ├── datasources/
│   │       │   │   │   ├── auth_remote_datasource.dart   # Dio API calls
│   │       │   │   │   └── auth_local_datasource.dart    # Secure storage (tokens)
│   │       │   │   ├── models/
│   │       │   │   │   ├── user_model.dart               # DTO with fromJson/toJson
│   │       │   │   │   ├── user_model.freezed.dart       # Generated
│   │       │   │   │   └── login_request.dart
│   │       │   │   └── repositories/
│   │       │   │       └── auth_repository_impl.dart     # Implements domain interface
│   │       │   └── presentation/
│   │       │       ├── bloc/
│   │       │       │   ├── auth_bloc.dart
│   │       │       │   ├── auth_event.dart
│   │       │       │   └── auth_state.dart
│   │       │       ├── pages/
│   │       │       │   ├── login_page.dart
│   │       │       │   ├── register_page.dart
│   │       │       │   └── forgot_password_page.dart
│   │       │       └── widgets/
│   │       │           ├── login_form.dart
│   │       │           ├── register_form.dart
│   │       │           └── social_login_buttons.dart
│   │       │
│   │       ├── home/                     # ── Home Feature ──
│   │       │   ├── domain/
│   │       │   │   └── usecases/
│   │       │   │       ├── get_featured_products.dart
│   │       │   │       ├── get_categories.dart
│   │       │   │       └── get_banners.dart
│   │       │   ├── data/
│   │       │   │   ├── datasources/
│   │       │   │   │   └── home_remote_datasource.dart
│   │       │   │   └── repositories/
│   │       │   │       └── home_repository_impl.dart
│   │       │   └── presentation/
│   │       │       ├── cubit/
│   │       │       │   ├── home_cubit.dart
│   │       │       │   └── home_state.dart
│   │       │       ├── pages/
│   │       │       │   └── home_page.dart
│   │       │       └── widgets/
│   │       │           ├── banner_carousel.dart
│   │       │           ├── category_grid.dart
│   │       │           ├── featured_products.dart
│   │       │           └── trending_section.dart
│   │       │
│   │       ├── shop/                     # ── Shop/Products Feature ──
│   │       │   ├── domain/
│   │       │   │   ├── entities/
│   │       │   │   │   ├── product.dart
│   │       │   │   │   ├── category.dart
│   │       │   │   │   ├── variant.dart              # Variant with SKU, price, stock, optionValues
│   │       │   │   │   ├── product_option.dart       # ProductOption + ProductOptionValue
│   │       │   │   │   └── product_attribute.dart    # AttributeDefinition + ProductAttributeValue
│   │       │   │   ├── repositories/
│   │       │   │   │   └── product_repository.dart
│   │       │   │   └── usecases/
│   │       │   │       ├── get_products.dart
│   │       │   │       ├── get_product_detail.dart
│   │       │   │       └── get_products_by_category.dart
│   │       │   ├── data/
│   │       │   │   ├── datasources/
│   │       │   │   │   └── product_remote_datasource.dart
│   │       │   │   ├── models/
│   │       │   │   │   ├── product_model.dart
│   │       │   │   │   └── product_model.freezed.dart
│   │       │   │   └── repositories/
│   │       │   │       └── product_repository_impl.dart
│   │       │   └── presentation/
│   │       │       ├── bloc/
│   │       │       │   ├── product_list_bloc.dart
│   │       │       │   ├── product_detail_cubit.dart
│   │       │       │   └── filter_cubit.dart
│   │       │       ├── pages/
│   │       │       │   ├── product_list_page.dart
│   │       │       │   ├── product_detail_page.dart
│   │       │       │   └── category_page.dart
│   │       │       └── widgets/
│   │       │           ├── product_card.dart
│   │       │           ├── product_grid.dart
│   │       │           ├── product_image_gallery.dart
│   │       │           ├── variant_selector.dart        # Color swatches + size pills
│   │       │           ├── product_attributes.dart      # Display attributes (Brand, Material)
│   │       │           ├── attribute_filter_panel.dart   # Dynamic category-based filters
│   │       │           ├── filter_bottom_sheet.dart
│   │       │           └── sort_dropdown.dart
│   │       │
│   │       ├── search/                   # ── Search Feature ──
│   │       │   ├── domain/usecases/
│   │       │   │   ├── search_products.dart
│   │       │   │   └── search_by_image.dart
│   │       │   ├── data/...
│   │       │   └── presentation/
│   │       │       ├── cubit/search_cubit.dart
│   │       │       ├── pages/search_page.dart
│   │       │       └── widgets/
│   │       │           ├── search_bar_widget.dart
│   │       │           ├── search_suggestions.dart
│   │       │           ├── recent_searches.dart
│   │       │           └── image_search_button.dart
│   │       │
│   │       ├── cart/                      # ── Cart Feature ──
│   │       │   ├── domain/
│   │       │   │   ├── entities/cart.dart, cart_item.dart
│   │       │   │   ├── repositories/cart_repository.dart
│   │       │   │   └── usecases/
│   │       │   │       ├── add_to_cart.dart
│   │       │   │       ├── remove_from_cart.dart
│   │       │   │       ├── update_quantity.dart
│   │       │   │       └── get_cart.dart
│   │       │   ├── data/...
│   │       │   └── presentation/
│   │       │       ├── bloc/cart_bloc.dart
│   │       │       ├── pages/cart_page.dart
│   │       │       └── widgets/
│   │       │           ├── cart_item_card.dart
│   │       │           ├── cart_summary.dart
│   │       │           ├── quantity_stepper.dart
│   │       │           └── swipe_to_delete.dart
│   │       │
│   │       ├── checkout/                 # ── Checkout Feature ──
│   │       │   ├── domain/usecases/
│   │       │   │   ├── create_order.dart
│   │       │   │   ├── create_payment_intent.dart
│   │       │   │   └── apply_coupon.dart
│   │       │   ├── data/...
│   │       │   └── presentation/
│   │       │       ├── bloc/checkout_bloc.dart
│   │       │       ├── pages/
│   │       │       │   ├── checkout_page.dart
│   │       │       │   └── order_success_page.dart
│   │       │       └── widgets/
│   │       │           ├── address_selector.dart
│   │       │           ├── payment_method_selector.dart
│   │       │           ├── stripe_payment_sheet.dart
│   │       │           ├── coupon_input.dart
│   │       │           └── order_summary.dart
│   │       │
│   │       ├── orders/                   # ── Orders Feature ──
│   │       │   ├── domain/usecases/
│   │       │   │   ├── get_orders.dart
│   │       │   │   ├── get_order_detail.dart
│   │       │   │   └── cancel_order.dart
│   │       │   ├── data/...
│   │       │   └── presentation/
│   │       │       ├── cubit/orders_cubit.dart
│   │       │       ├── pages/
│   │       │       │   ├── orders_page.dart
│   │       │       │   └── order_detail_page.dart
│   │       │       └── widgets/
│   │       │           ├── order_card.dart
│   │       │           ├── order_timeline.dart
│   │       │           └── order_status_badge.dart
│   │       │
│   │       ├── profile/                  # ── Profile/Account Feature ──
│   │       │   ├── domain/usecases/
│   │       │   │   ├── get_profile.dart
│   │       │   │   ├── update_profile.dart
│   │       │   │   └── manage_addresses.dart
│   │       │   ├── data/...
│   │       │   └── presentation/
│   │       │       ├── cubit/profile_cubit.dart
│   │       │       ├── pages/
│   │       │       │   ├── profile_page.dart
│   │       │       │   ├── edit_profile_page.dart
│   │       │       │   ├── addresses_page.dart
│   │       │       │   └── settings_page.dart
│   │       │       └── widgets/
│   │       │           ├── profile_header.dart
│   │       │           ├── address_card.dart
│   │       │           └── address_form.dart
│   │       │
│   │       ├── wishlist/                 # ── Wishlist Feature ──
│   │       │   ├── domain/...
│   │       │   ├── data/...
│   │       │   └── presentation/
│   │       │       ├── cubit/wishlist_cubit.dart
│   │       │       └── pages/wishlist_page.dart
│   │       │
│   │       ├── reviews/                  # ── Reviews Feature ──
│   │       │   ├── domain/usecases/
│   │       │   │   ├── get_reviews.dart
│   │       │   │   └── create_review.dart
│   │       │   ├── data/...
│   │       │   └── presentation/
│   │       │       ├── cubit/reviews_cubit.dart
│   │       │       └── widgets/
│   │       │           ├── review_card.dart
│   │       │           ├── review_form.dart
│   │       │           ├── rating_bar.dart
│   │       │           └── rating_summary.dart
│   │       │
│   │       ├── chat/                     # ── Chat Feature ──
│   │       │   ├── domain/...
│   │       │   ├── data/...
│   │       │   └── presentation/
│   │       │       ├── bloc/chat_bloc.dart
│   │       │       ├── pages/
│   │       │       │   ├── conversations_page.dart
│   │       │       │   └── chat_page.dart
│   │       │       └── widgets/
│   │       │           ├── conversation_tile.dart
│   │       │           ├── message_bubble.dart
│   │       │           └── chat_input.dart
│   │       │
│   │       ├── notifications/            # ── Notifications Feature ──
│   │       │   ├── data/...
│   │       │   └── presentation/
│   │       │       ├── cubit/notification_cubit.dart
│   │       │       └── pages/notifications_page.dart
│   │       │
│   │       └── ai/                       # ── AI Features ──
│   │           └── presentation/
│   │               ├── cubit/ai_chat_cubit.dart
│   │               └── widgets/
│   │                   ├── ai_chat_fab.dart           # Floating action button
│   │                   ├── ai_chat_bottom_sheet.dart
│   │                   └── recommendation_carousel.dart
│   │
│   ├── test/                             # Unit + widget tests mirror lib/ structure
│   ├── integration_test/                 # Integration tests
│   ├── android/
│   ├── ios/
│   └── assets/
│
├── seller_app/                           # ═══ SELLER/ADMIN APP ═══
│   ├── pubspec.yaml
│   ├── lib/
│   │   ├── main.dart
│   │   ├── app.dart
│   │   ├── injection.dart
│   │   │
│   │   ├── core/
│   │   │   ├── config/
│   │   │   ├── router/app_router.dart
│   │   │   └── constants/
│   │   │
│   │   └── features/
│   │       ├── auth/                     # Seller/Admin login
│   │       │   └── (same clean arch structure)
│   │       │
│   │       ├── dashboard/                # ── Dashboard Feature ──
│   │       │   ├── domain/usecases/get_dashboard_stats.dart
│   │       │   └── presentation/
│   │       │       ├── cubit/dashboard_cubit.dart
│   │       │       ├── pages/dashboard_page.dart
│   │       │       └── widgets/
│   │       │           ├── stats_card.dart
│   │       │           ├── revenue_chart.dart
│   │       │           ├── recent_orders_list.dart
│   │       │           └── sales_overview.dart
│   │       │
│   │       ├── products/                 # ── Product Management ──
│   │       │   ├── domain/usecases/
│   │       │   │   ├── create_product.dart
│   │       │   │   ├── update_product.dart
│   │       │   │   ├── delete_product.dart
│   │       │   │   └── manage_variants.dart
│   │       │   └── presentation/
│   │       │       ├── bloc/product_management_bloc.dart
│   │       │       ├── pages/
│   │       │       │   ├── products_page.dart
│   │       │       │   ├── product_form_page.dart    # Create + Edit
│   │       │       │   └── variant_manager_page.dart
│   │       │       └── widgets/
│   │       │           ├── product_form.dart          # Basic info + category selection
│   │       │           ├── attribute_form.dart        # Dynamic fields from category attributes
│   │       │           ├── option_manager.dart        # Define variant axes (Color, Size) + values
│   │       │           ├── variant_table.dart         # Inline edit grid for all variants (SKU, price, stock)
│   │       │           ├── image_picker_grid.dart     # Upload + assign images to variants
│   │       │           ├── variant_form.dart          # Single variant edit dialog
│   │       │           ├── inventory_input.dart
│   │       │           └── ai_description_button.dart
│   │       │
│   │       ├── orders/                   # ── Order Management ──
│   │       │   ├── domain/usecases/
│   │       │   │   ├── get_seller_orders.dart
│   │       │   │   ├── update_order_status.dart
│   │       │   │   └── add_tracking.dart
│   │       │   └── presentation/
│   │       │       ├── bloc/seller_orders_bloc.dart
│   │       │       ├── pages/
│   │       │       │   ├── seller_orders_page.dart
│   │       │       │   └── order_detail_page.dart
│   │       │       └── widgets/
│   │       │           ├── order_card.dart
│   │       │           ├── status_update_dialog.dart
│   │       │           └── tracking_form.dart
│   │       │
│   │       ├── analytics/                # ── Analytics Feature ──
│   │       │   └── presentation/
│   │       │       ├── cubit/analytics_cubit.dart
│   │       │       ├── pages/analytics_page.dart
│   │       │       └── widgets/
│   │       │           ├── revenue_line_chart.dart
│   │       │           ├── orders_bar_chart.dart
│   │       │           ├── top_products_list.dart
│   │       │           └── date_range_picker.dart
│   │       │
│   │       ├── payouts/                  # ── Payouts Feature ──
│   │       │   └── presentation/
│   │       │       ├── cubit/payout_cubit.dart
│   │       │       └── pages/payouts_page.dart
│   │       │
│   │       ├── store_settings/           # ── Store Settings ──
│   │       │   └── presentation/
│   │       │       └── pages/store_settings_page.dart
│   │       │
│   │       ├── chat/                     # ── Chat with Buyers ──
│   │       │   └── (same as buyer_app chat feature)
│   │       │
│   │       └── notifications/
│   │           └── (same structure)
│   │
│   ├── test/
│   ├── integration_test/
│   ├── android/
│   └── ios/
│
└── packages/                             # ═══ SHARED FLUTTER PACKAGES ═══
    │
    ├── core/                             # Core domain + infrastructure
    │   ├── pubspec.yaml
    │   └── lib/
    │       ├── core.dart                 # Barrel export
    │       ├── error/
    │       │   ├── exceptions.dart       # ServerException, CacheException, etc.
    │       │   └── failures.dart         # Failure sealed class (ServerFailure, NetworkFailure)
    │       ├── usecase/
    │       │   └── usecase.dart          # abstract UseCase<Type, Params>
    │       ├── network/
    │       │   └── network_info.dart     # Connectivity check interface
    │       └── utils/
    │           ├── either.dart           # Either<Failure, Success> (or use dartz/fpdart)
    │           ├── input_converter.dart
    │           └── date_formatter.dart
    │
    ├── api_client/                       # Shared Dio HTTP client
    │   ├── pubspec.yaml
    │   └── lib/
    │       ├── api_client.dart           # Barrel export
    │       ├── dio_client.dart           # Dio instance + interceptors
    │       ├── interceptors/
    │       │   ├── auth_interceptor.dart  # Attach JWT, auto-refresh on 401
    │       │   ├── logging_interceptor.dart
    │       │   └── error_interceptor.dart # Map HTTP errors → domain exceptions
    │       ├── api_endpoints.dart        # All endpoint constants
    │       └── api_response.dart         # Generic ApiResponse<T>, PaginatedResponse<T>
    │
    ├── ui_kit/                           # Shared design system
    │   ├── pubspec.yaml
    │   └── lib/
    │       ├── ui_kit.dart               # Barrel export
    │       ├── theme/
    │       │   ├── app_theme.dart        # Light + dark ThemeData
    │       │   ├── app_colors.dart       # Brand color palette
    │       │   ├── app_typography.dart   # Text styles
    │       │   └── app_spacing.dart      # Spacing constants
    │       ├── widgets/
    │       │   ├── app_button.dart
    │       │   ├── app_text_field.dart
    │       │   ├── app_card.dart
    │       │   ├── app_bottom_sheet.dart
    │       │   ├── app_dialog.dart
    │       │   ├── loading_overlay.dart
    │       │   ├── error_widget.dart
    │       │   ├── empty_state.dart
    │       │   ├── rating_stars.dart
    │       │   ├── price_text.dart
    │       │   ├── cached_image.dart
    │       │   └── shimmer_loading.dart
    │       └── extensions/
    │           ├── context_extensions.dart  # context.theme, context.screenSize
    │           └── string_extensions.dart
    │
    └── shared_models/                    # Shared DTOs & enums
        ├── pubspec.yaml
        └── lib/
            ├── shared_models.dart        # Barrel export
            ├── enums/
            │   ├── order_status.dart
            │   ├── product_status.dart
            │   └── user_role.dart
            └── dtos/
                ├── user_dto.dart
                ├── product_dto.dart
                ├── order_dto.dart
                ├── pagination_dto.dart
                └── ... (all with freezed + json_serializable)
```

### Clean Architecture Layers (per feature)

```
┌─────────────────────────────────────────────────┐
│                PRESENTATION                      │
│  BLoC / Cubit ← Pages ← Widgets                │
│  Depends on: Domain (use cases)                  │
│  Framework: Flutter, flutter_bloc                │
├─────────────────────────────────────────────────┤
│                DOMAIN                            │
│  Entities ← Use Cases ← Repository Interfaces   │
│  Depends on: NOTHING (pure Dart)                 │
│  No Flutter imports, no external packages        │
├─────────────────────────────────────────────────┤
│                DATA                              │
│  Models ← DataSources ← Repository Impls        │
│  Depends on: Domain (implements interfaces)      │
│  Framework: Dio, Hive, flutter_secure_storage    │
└─────────────────────────────────────────────────┘

Dependency rule: outer layers depend on inner layers, NEVER the reverse.
```

### Key Code Patterns

**UseCase Base Class** (`packages/core/lib/usecase/usecase.dart`)
```dart
import 'package:dartz/dartz.dart';
import '../error/failures.dart';

abstract class UseCase<Type, Params> {
  Future<Either<Failure, Type>> call(Params params);
}

class NoParams {}
```

**Domain Entity** (`features/shop/domain/entities/product.dart`)
```dart
class Product {
  final String id;
  final String name;
  final String slug;
  final String description;
  final int basePriceCents;
  final String currency;
  final String status;
  final bool hasVariants;
  final List<String> imageUrls;
  final List<ProductAttributeValue> attributeValues;
  final List<ProductOption> options;
  final List<Variant> variants;
  final double ratingAvg;
  final int ratingCount;
  final String categoryId;
  final String sellerId;

  const Product({
    required this.id, required this.name, required this.slug,
    required this.description, required this.basePriceCents,
    this.currency = 'USD', this.status = 'active',
    this.hasVariants = false, this.imageUrls = const [],
    this.attributeValues = const [], this.options = const [],
    this.variants = const [], this.ratingAvg = 0,
    this.ratingCount = 0, required this.categoryId, required this.sellerId,
  });

  /// Effective price: if has variants, use default variant price; otherwise base price.
  int get effectivePriceCents {
    if (!hasVariants) return basePriceCents;
    final defaultVariant = variants.where((v) => v.isDefault).firstOrNull;
    if (defaultVariant != null && defaultVariant.priceCents > 0) {
      return defaultVariant.priceCents;
    }
    return basePriceCents;
  }

  String get formattedPrice => '\$${(effectivePriceCents / 100).toStringAsFixed(2)}';

  /// Total stock across all active variants, or 0 for non-variant products.
  int get totalStock => variants.fold(0, (sum, v) => sum + (v.isActive ? v.stock : 0));
}

class ProductAttributeValue {
  final String attributeId;
  final String attributeName;
  final String value;
  final List<String>? values; // for multi_select

  const ProductAttributeValue({
    required this.attributeId, required this.attributeName,
    required this.value, this.values,
  });
}

class ProductOption {
  final String id;
  final String name;           // "Color", "Size"
  final int sortOrder;
  final List<ProductOptionValue> values;

  const ProductOption({
    required this.id, required this.name,
    this.sortOrder = 0, this.values = const [],
  });
}

class ProductOptionValue {
  final String id;
  final String value;          // "Red", "XL"
  final String? colorHex;      // "#FF0000"
  final int sortOrder;

  const ProductOptionValue({
    required this.id, required this.value,
    this.colorHex, this.sortOrder = 0,
  });
}

class Variant {
  final String id;
  final String sku;
  final String name;           // "Red / XL"
  final int priceCents;        // 0 = use base price
  final int compareAtCents;
  final int stock;
  final bool isDefault;
  final bool isActive;
  final List<String> imageUrls;
  final String? barcode;
  final List<VariantOptionValue> optionValues;

  const Variant({
    required this.id, required this.sku, required this.name,
    this.priceCents = 0, this.compareAtCents = 0,
    this.stock = 0, this.isDefault = false, this.isActive = true,
    this.imageUrls = const [], this.barcode,
    this.optionValues = const [],
  });
}

class VariantOptionValue {
  final String optionName;     // "Color"
  final String value;          // "Red"

  const VariantOptionValue({required this.optionName, required this.value});
}
```

**Data Model with Freezed** (`features/shop/data/models/product_model.dart`)
```dart
@freezed
class ProductModel with _$ProductModel {
  const factory ProductModel({
    required String id,
    required String name,
    required String slug,
    required String description,
    @JsonKey(name: 'base_price_cents') required int basePriceCents,
    @Default('USD') String currency,
    @Default('active') String status,
    @JsonKey(name: 'has_variants') @Default(false) bool hasVariants,
    @JsonKey(name: 'image_urls') @Default([]) List<String> imageUrls,
    @JsonKey(name: 'attribute_values') @Default([]) List<ProductAttributeValueModel> attributeValues,
    @Default([]) List<ProductOptionModel> options,
    @Default([]) List<VariantModel> variants,
    @JsonKey(name: 'rating_avg') @Default(0) double ratingAvg,
    @JsonKey(name: 'rating_count') @Default(0) int ratingCount,
    @JsonKey(name: 'category_id') required String categoryId,
    @JsonKey(name: 'seller_id') required String sellerId,
  }) = _ProductModel;

  factory ProductModel.fromJson(Map<String, dynamic> json) => _$ProductModelFromJson(json);
}

@freezed
class VariantModel with _$VariantModel {
  const factory VariantModel({
    required String id,
    required String sku,
    required String name,
    @JsonKey(name: 'price_cents') @Default(0) int priceCents,
    @JsonKey(name: 'compare_at_cents') @Default(0) int compareAtCents,
    @Default(0) int stock,
    @JsonKey(name: 'is_default') @Default(false) bool isDefault,
    @JsonKey(name: 'is_active') @Default(true) bool isActive,
    @JsonKey(name: 'image_urls') @Default([]) List<String> imageUrls,
    String? barcode,
    @JsonKey(name: 'option_values') @Default([]) List<VariantOptionValueModel> optionValues,
  }) = _VariantModel;

  factory VariantModel.fromJson(Map<String, dynamic> json) => _$VariantModelFromJson(json);
}

extension ProductModelX on ProductModel {
  Product toEntity() => Product(
    id: id, name: name, slug: slug, description: description,
    basePriceCents: basePriceCents, currency: currency, status: status,
    hasVariants: hasVariants, imageUrls: imageUrls,
    attributeValues: attributeValues.map((a) => a.toEntity()).toList(),
    options: options.map((o) => o.toEntity()).toList(),
    variants: variants.map((v) => v.toEntity()).toList(),
    ratingAvg: ratingAvg, ratingCount: ratingCount,
    categoryId: categoryId, sellerId: sellerId,
  );
}
```

**Repository Interface** (`features/shop/domain/repositories/product_repository.dart`)
```dart
abstract class ProductRepository {
  Future<Either<Failure, PaginatedResult<Product>>> getProducts(ProductFilter filter);
  Future<Either<Failure, Product>> getProductBySlug(String slug);
  Future<Either<Failure, List<Category>>> getCategories();
}
```

**Repository Implementation** (`features/shop/data/repositories/product_repository_impl.dart`)
```dart
@Injectable(as: ProductRepository)
class ProductRepositoryImpl implements ProductRepository {
  final ProductRemoteDataSource _remote;

  ProductRepositoryImpl(this._remote);

  @override
  Future<Either<Failure, PaginatedResult<Product>>> getProducts(ProductFilter filter) async {
    try {
      final response = await _remote.getProducts(filter);
      return Right(response.toEntity());
    } on ServerException catch (e) {
      return Left(ServerFailure(e.message));
    } on NetworkException {
      return Left(NetworkFailure());
    }
  }
}
```

**BLoC** (`features/shop/presentation/bloc/product_list_bloc.dart`)
```dart
@injectable
class ProductListBloc extends Bloc<ProductListEvent, ProductListState> {
  final GetProducts _getProducts;

  ProductListBloc(this._getProducts) : super(ProductListInitial()) {
    on<LoadProducts>(_onLoadProducts);
    on<LoadMoreProducts>(_onLoadMore);
    on<ApplyFilters>(_onApplyFilters);
  }

  Future<void> _onLoadProducts(LoadProducts event, Emitter<ProductListState> emit) async {
    emit(ProductListLoading());
    final result = await _getProducts(ProductFilter(page: 1, pageSize: 20));
    result.fold(
      (failure) => emit(ProductListError(failure.message)),
      (products) => emit(ProductListLoaded(products: products.items, hasMore: products.hasMore)),
    );
  }
}
```

**Dio Auth Interceptor** (`packages/api_client/lib/interceptors/auth_interceptor.dart`)
```dart
@injectable
class AuthInterceptor extends Interceptor {
  final TokenStorage _tokenStorage;
  final Dio _dio;

  AuthInterceptor(this._tokenStorage, this._dio);

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) async {
    final token = await _tokenStorage.getAccessToken();
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode == 401) {
      try {
        final refreshToken = await _tokenStorage.getRefreshToken();
        final response = await _dio.post('/auth/refresh', data: {'refresh_token': refreshToken});
        final newToken = response.data['access_token'];
        await _tokenStorage.saveAccessToken(newToken);

        // Retry original request
        err.requestOptions.headers['Authorization'] = 'Bearer $newToken';
        final retryResponse = await _dio.fetch(err.requestOptions);
        handler.resolve(retryResponse);
        return;
      } catch (_) {
        await _tokenStorage.clear();
        // Navigate to login (via event bus or global key)
      }
    }
    handler.next(err);
  }
}
```

### Dependency Injection Setup (`injection.dart`)
```dart
import 'package:get_it/get_it.dart';
import 'package:injectable/injectable.dart';
import 'injection.config.dart';

final getIt = GetIt.instance;

@InjectableInit()
void configureDependencies() => getIt.init();
```

All classes annotated with `@injectable`, `@singleton`, or `@lazySingleton` are auto-registered. Repository impls use `@Injectable(as: ProductRepository)` to bind interface → implementation.

### Navigation (go_router)
```dart
final appRouter = GoRouter(
  initialLocation: '/',
  redirect: (context, state) {
    final isLoggedIn = getIt<AuthBloc>().state is Authenticated;
    final isAuthRoute = state.matchedLocation.startsWith('/auth');
    if (!isLoggedIn && !isAuthRoute && _requiresAuth(state.matchedLocation)) return '/auth/login';
    if (isLoggedIn && isAuthRoute) return '/';
    return null;
  },
  routes: [
    GoRoute(path: '/', builder: (_, __) => const HomePage()),
    GoRoute(path: '/products', builder: (_, __) => const ProductListPage()),
    GoRoute(path: '/products/:slug', builder: (_, state) =>
      ProductDetailPage(slug: state.pathParameters['slug']!)),
    GoRoute(path: '/search', builder: (_, __) => const SearchPage()),
    GoRoute(path: '/cart', builder: (_, __) => const CartPage()),
    GoRoute(path: '/checkout', builder: (_, __) => const CheckoutPage()),
    GoRoute(path: '/orders', builder: (_, __) => const OrdersPage()),
    GoRoute(path: '/orders/:id', builder: (_, state) =>
      OrderDetailPage(orderId: state.pathParameters['id']!)),
    GoRoute(path: '/profile', builder: (_, __) => const ProfilePage()),
    GoRoute(path: '/wishlist', builder: (_, __) => const WishlistPage()),
    GoRoute(path: '/messages', builder: (_, __) => const ConversationsPage()),
    ShellRoute(builder: (_, __, child) => AuthShell(child: child), routes: [
      GoRoute(path: '/auth/login', builder: (_, __) => const LoginPage()),
      GoRoute(path: '/auth/register', builder: (_, __) => const RegisterPage()),
    ]),
  ],
);
```

### Testing Strategy (Flutter)

| Layer | Tool | What to Test |
|---|---|---|
| Domain entities | dart test | Business logic, value calculations |
| Use cases | dart test + mockito | Orchestration, Either<Failure, Success> |
| Repositories | dart test + mockito | Data source calls, error mapping, model→entity |
| BLoC/Cubit | bloc_test | State transitions, event handling |
| Widgets | flutter_test | Widget rendering, user interaction |
| Integration | integration_test | Full user flows on device/emulator |
| Golden | golden_toolkit | Visual regression tests for UI components |

```dart
// Example: BLoC test
blocTest<ProductListBloc, ProductListState>(
  'emits [Loading, Loaded] when LoadProducts succeeds',
  build: () {
    when(() => getProducts(any())).thenAnswer(
      (_) async => Right(PaginatedResult(items: [mockProduct], hasMore: false)),
    );
    return ProductListBloc(getProducts);
  },
  act: (bloc) => bloc.add(LoadProducts()),
  expect: () => [
    ProductListLoading(),
    ProductListLoaded(products: [mockProduct], hasMore: false),
  ],
);
```

---

## Part 19: Implementation Roadmap

### Phase 1 — Foundation (Week 1-2)
- Initialize Go workspace, `pkg/` shared packages (including `pkg/i18n`, `pkg/money`, `pkg/tax`), proto definitions, Makefile
- Generate protobuf Go code (`make proto`)
- Set up docker-compose with all infrastructure services (20 Go services)
- Implement **Auth Service** (register, login, JWT, refresh, OAuth, social login)
- Implement **User Service** (profiles, addresses, seller profiles, follow sellers)

### Phase 2 — Core Commerce (Week 3-5)
- Implement **Product Service** (CRUD, categories, attributes, options, variants, stock)
- Implement **Cart Service** (Redis-backed, merge on login, coupon application, points redemption)
- Implement **Order Service** (state machine, multi-seller order splitting)
- Implement **Payment Service** (Stripe intents, webhooks, Connect, seller wallets, settlements)
- Implement **Tax Service** (tax rules engine, jurisdiction config, tax calculation)
- Wire NATS events: order→payment→stock→tax

### Phase 3 — Shipping & Fulfillment (Week 6-7)
- Implement **Shipping Service** (carrier integration framework, rate shopping, label generation, tracking)
- Integrate first carriers (FedEx, UPS, DHL, USPS)
- Implement **Return Service** (return requests, approval workflow, refund processing)
- Implement **Dispute Service** (dispute creation, messaging, admin resolution)
- Wire NATS events: shipment→notification, return→payment→product(restock)

### Phase 4 — Supporting Services (Week 8-9)
- Implement **Search Service** (Elasticsearch index/query, autocomplete, faceted attribute filters)
- Implement **Review Service** (CRUD, rating aggregation)
- Implement **Media Service** (S3 presigned URLs, image processing)
- Implement **Notification Service** (email, push, WebSocket, i18n templates)
- Implement **CMS Service** (banners, static pages, content scheduling)
- Configure **Kong Gateway** (kong.yml routes for all 20 services, JWT, ACL, rate limiting, CORS)

### Phase 5 — Promotions & Engagement (Week 10-11)
- Implement **Promotion Service** (coupons, vouchers, flash sales, bundles)
- Implement **Loyalty Service** (points, tiers, cashback, earn/redeem)
- Implement **Affiliate Service** (referral links, click tracking, conversion tracking, payouts)
- Wire NATS events: order.completed→loyalty(earn), referral→affiliate(conversion)

### Phase 6 — Real-Time & AI (Week 12-13)
- Implement **Chat Service** (WebSocket messaging)
- Implement **AI Service** (embedding pipeline, recommendation proxy, chatbot)
- Deploy Python FastAPI AI service
- Connect frontend AI components (chat widget, search, recommendations)

### Phase 7 — React Frontend (Week 14-17)
- **Week 14**: Core modules — Auth, Shop (products, variants, attributes, filters), Search, Cart
- **Week 15**: Checkout (shipping rate selection, coupon, points, tax), Orders, Returns, Tracking
- **Week 16**: Seller Dashboard (products, orders, returns, shipments, coupons, carriers, analytics, payouts)
- **Week 17**: Admin Dashboard (users, sellers, products, attributes, categories, orders, disputes, promotions, flash sales, carriers, banners, pages, affiliates, tax, reports)
- **Week 17**: Loyalty dashboard, Affiliate dashboard, CMS pages, Social features (follow, share)

### Phase 8 — Flutter Mobile Apps (Week 18-21)
- **Week 18**: Set up Flutter monorepo, shared packages (core, api_client, ui_kit, shared_models)
- **Week 18**: Configure get_it + injectable DI, Dio client with auth interceptor
- **Week 19**: **Buyer App** — Auth (login, register, OAuth), Home, Shop (product list + detail + variants)
- **Week 19**: **Buyer App** — Search (text + image), Cart (coupons, points), Checkout (shipping rates, tax)
- **Week 20**: **Buyer App** — Orders, Returns, Tracking, Profile, Wishlist, Reviews, Loyalty, Affiliate
- **Week 20**: **Buyer App** — Notifications, Chat, AI assistant
- **Week 21**: **Seller App** — Auth, Dashboard, Product Management (CRUD + attributes + variants + images)
- **Week 21**: **Seller App** — Order management, Returns, Shipments, Coupons, Analytics, Payouts
- **Week 21**: Push notifications (Firebase Cloud Messaging), deep linking

### Phase 9 — Production (Week 22-24)
- OpenTelemetry tracing across all services
- Circuit breakers on all gRPC calls
- i18n: multi-language support (English, Spanish, French, Chinese, Arabic)
- Multi-currency: exchange rate service integration
- Integration tests (Go: testcontainers, Flutter: integration_test)
- E2E tests against Kong gateway
- GitHub Actions CI/CD pipeline (Go + React + Flutter)
- Kubernetes manifests: Deployments, Services, HPA, Ingress
- App Store / Play Store submission preparation
- Load testing with k6
- Security audit (OWASP, PCI compliance for payments)

---

## Verification Plan

### Infrastructure
1. `make run-infra` — PostgreSQL, Redis, NATS, Elasticsearch healthy
2. `make proto` — all protobuf code generated without errors (20 services)
3. `make build` — all 20 Go services compile
4. `make test` — all unit tests pass
5. `docker-compose up` — all services + Kong start and report healthy
6. `curl http://localhost:8001/status` — Kong Admin API confirms all routes loaded (60+ routes)

### Auth & Users
7. `curl POST /api/v1/auth/register` → user created, JWT returned, loyalty membership auto-created
8. `curl POST /api/v1/auth/register?ref=ABC123` → referral tracked in affiliate service
9. `curl POST /api/v1/users/:id/follow` → follow seller, verify follower count

### Products & Variants
10. Create product with category attributes → required attributes validated
11. Add options (Color: Red/Blue, Size: S/M/L) → generate 6 variants with SKUs
12. Update variant prices/stock individually → verify in product detail response

### Shopping & Checkout
13. Add product (with variant) to cart → verify variant-specific price/image
14. Apply coupon code → discount validated and applied
15. Get shipping rates → rate shopping returns multiple carrier options
16. Calculate tax → correct jurisdiction-based tax breakdown
17. Redeem loyalty points → points deducted, discount applied
18. Place order → payment intent created, stock reserved, tax calculated

### Fulfillment & Post-Order
19. Seller creates shipment → label generated, tracking number assigned
20. Tracking events update → buyer notified, status timeline updates
21. Shipment delivered → review prompt sent, loyalty points earned
22. Buyer requests return → seller notified, return workflow starts
23. Return approved → refund processed, stock restocked
24. Dispute opened → admin notified, messaging thread created

### Promotions & Engagement
25. Flash sale starts → products shown at sale price, countdown timer
26. Flash sale ends → prices revert, sold counts tracked
27. Bundle deal purchased → bundle discount applied correctly
28. Loyalty tier upgrade → notification sent, new benefits activated
29. Affiliate link clicked → click tracked with attribution cookie
30. Referred user places order → commission calculated, referrer notified

### CMS & Content
31. Admin creates banner → appears on homepage with correct scheduling
32. Admin publishes page → accessible at `/pages/:slug`
33. Content schedule executes → auto-publish/unpublish at scheduled time

### Security & Access Control
34. `curl /api/v1/admin/dashboard` without admin role → Kong returns 403
35. Seller cannot access another seller's products/orders → 403
36. Rate limiting works on auth endpoints (10/min) and AI endpoints (30/min)

### Mobile
37. **Flutter Buyer App**: login → browse → filter by attributes → select variant → cart → apply coupon → checkout → track → return → loyalty
38. **Flutter Seller App**: login → create product with variants → manage orders → create shipment → view returns → manage coupons → view payouts
39. `flutter test` — all unit + widget tests pass for both apps
40. `flutter test integration_test/` — critical user flows pass on device
