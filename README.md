# E-Commerce Backend - Gin Framework

A comprehensive RESTful API backend for an e-commerce platform built with Go, Gin web framework, MySQL database, and Redis cache. Features include user authentication, product management, shopping cart, order processing, payment integration, cloud storage for product images, **complete Swagger/OpenAPI documentation**, and **Docker containerization** for easy deployment.

## 🚀 Features

### 🐳 Docker Deployment
- **Fully containerized** application with Docker
- **Docker Compose** for multi-container orchestration
- **Development and production** configurations
- **One-command deployment** with all dependencies
- **Persistent storage** with Docker volumes
- **Health checks** for all services
- **Web-based management tools** (phpMyAdmin, Redis Commander)

### 📚 API Documentation
- **Interactive Swagger UI** for testing all endpoints
- **Complete OpenAPI 3.0 specification** with request/response schemas
- **Authentication support** in Swagger with JWT Bearer tokens
- **Example requests and responses** for all endpoints
- **Quick start guide** for API testing

### Authentication & Authorization
- **User Registration & Login** with JWT-based authentication
- **Access & Refresh Token** mechanism for secure session management
- **Role-based Access Control** (Admin/Customer roles)
- **Password Hashing** using bcrypt
- **Secure logout** with token invalidation
- **HTTP-only cookies** for refresh tokens

### Product Management
- Product CRUD operations (Create, Read, Update, Delete)
- Product listing with **pagination and filtering**
- Product search by name and category
- Product retrieval by slug with **Redis caching**
- **Rate limiting** on product endpoints
- Admin-only product management endpoints
- Image upload to AWS S3/Cloudflare R2 storage

### Shopping Cart
- Add products to cart with **stock validation**
- Update cart item quantities
- Remove items from cart
- Clear entire cart
- View cart with **calculated totals and product details**
- **Automatic quantity updates** for existing cart items

### Order Management
- Checkout process with **cart and stock validation**
- **Atomic order creation** with database transactions
- Order creation and tracking
- Order history for users
- Admin order management (view all orders, update status, statistics)
- Order status tracking with **allowed transitions** (pending → confirmed → shipped → delivered)
- **Stock management** (deduction on checkout, restoration on cancellation)
- **Order filtering** by status and date range

### Payment Integration
- Payment intent creation
- Simulated payment gateway integration
- Payment webhook handling
- Payment status tracking
- **Automatic order confirmation** on successful payment

### File Upload
- Product image upload to **AWS S3** or Cloudflare R2
- Support for multiple image formats (JPEG, PNG, WebP)
- Automatic content-type detection
- Public URL generation for uploaded images
- **Admin-only** upload access

### Performance & Caching
- **Redis caching** for product data (5-minute TTL)
- Cache invalidation on product updates
- Database query optimization with GORM
- **Efficient pagination** for large datasets

## 🛠️ Tech Stack

- **Language**: Go 1.25.3
- **Web Framework**: Gin v1.11.0
- **Database**: MySQL with GORM v1.31.1
- **Cache**: Redis v9.17.2
- **Authentication**: JWT (golang-jwt/jwt/v4 v4.5.2)
- **API Documentation**: Swagger/OpenAPI (swaggo/swag v1.16.6)
- **Cloud Storage**: AWS S3 / Cloudflare R2 (AWS SDK v2)
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **Environment Config**: godotenv v1.5.1

## 📁 Project Structure

```
ecommerce-gin/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point with Swagger setup
├── docs/
│   ├── docs.go                  # Generated Swagger documentation
│   ├── swagger.json             # OpenAPI JSON specification
│   └── swagger.yaml             # OpenAPI YAML specification
├── internal/
│   ├── cache/
│   │   ├── cache.go             # Redis cache interface
│   │   └── redis.go             # Redis client implementation
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── controllers/
│   │   ├── auth_controller.go   # Authentication endpoints
│   │   ├── product_controller.go
│   │   ├── cart_controller.go
│   │   ├── order_controller.go
│   │   ├── admin_order_controller.go
│   │   ├── payment_controller.go
│   │   ├── payment_webhook_controller.go
│   │   ├── upload_controller.go
│   │   └── models.go            # Swagger model definitions
│   ├── database/
│   │   ├── db.go                # Database connection
│   │   ├── user_repo.go         # User repository
│   │   ├── product_repo.go
│   │   ├── cart_repo.go
│   │   ├── order_repo.go
│   │   ├── order_admin_repo.go
│   │   └── payment_repo.go
│   ├── middleware/
│   │   ├── jwt.go               # JWT authentication middleware
│   │   ├── admin.go             # Admin authorization middleware
│   │   └── rate_limit.go        # Rate limiting middleware
│   ├── models/
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── cart_item.go
│   │   ├── order.go
│   │   ├── order_item.go
│   │   └── payment_intent.go
│   ├── routes/
│   │   ├── routes.go            # Main routes setup
│   │   ├── product_routes.go
│   │   ├── cart_routes.go
│   │   ├── order_routes.go
│   │   ├── admin_routes.go
│   │   ├── payment_routes.go
│   │   └── upload_routes.go
│   ├── services/
│   │   ├── payment_service.go   # Payment processing logic
│   │   └── s3_service.go        # S3/R2 upload service
│   └── utils/
│       ├── hash.go              # Hashing utilities
│       ├── password.go          # Password validation
│       └── token.go             # JWT token generation
├── uploads/                     # Local uploads directory
├── Dockerfile                   # Docker image build instructions
├── docker-compose.yml           # Production Docker setup
├── docker-compose.dev.yml       # Development Docker setup
├── .dockerignore                # Files to exclude from Docker build
├── .env                         # Environment variables
├── .env.docker                  # Docker environment template
├── start-docker.sh              # Quick start script (Linux/Mac)
├── start-docker.bat             # Quick start script (Windows)
├── go.mod
├── README.md
├── DOCKER_GUIDE.md              # Complete Docker guide for beginners
├── SWAGGER_DOCUMENTATION.md     # Comprehensive API documentation
├── SWAGGER_QUICK_START.md       # Quick start guide for testing
└── S3_STORAGE_GUIDE.md          # S3/R2 setup guide
```

## 🔧 Installation & Setup

You can run this application in two ways:

### 🐳 Option A: Docker (Recommended - Easiest!)

**Perfect for beginners! Everything runs in containers.**

#### Prerequisites
- [Docker Desktop](https://www.docker.com/products/docker-desktop) (includes Docker Compose)

#### Quick Start (30 seconds!)

**Windows:**
```bash
# Double-click start-docker.bat
# OR run in PowerShell:
.\start-docker.bat
```

**Linux/Mac:**
```bash
chmod +x start-docker.sh
./start-docker.sh
```

That's it! The script will:
1. Check if Docker is installed
2. Create `.env` file from template
3. Build and start all services (MySQL, Redis, Go app)
4. Show you the URLs to access

**Access your application:**
- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html
- phpMyAdmin: http://localhost:8081
- Redis Commander: http://localhost:8082

**For detailed Docker instructions, see [`DOCKER_GUIDE.md`](DOCKER_GUIDE.md)**

### 💻 Option B: Local Development (Manual Setup)

**For when you want full control and direct access to services.**

#### Prerequisites
- Go 1.25 or higher
- MySQL 8.0 or higher
- Redis 6.0 or higher
- AWS S3 account or Cloudflare R2 account (for image uploads)
- swaggo/swag CLI tool (for regenerating API docs)

#### 1. Clone the repository
```bash
git clone <repository-url>
cd ecommerce-gin
```

#### 2. Install dependencies
```bash
go mod download
```

#### 3. Install Swagger CLI (optional, for regenerating docs)
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

#### 4. Generate JWT Secret

**What is JWT Secret?**

The JWT (JSON Web Token) Secret is a cryptographic key used to sign and verify authentication tokens. It ensures that:
- Tokens cannot be forged or tampered with
- Only your server can generate valid tokens
- Token authenticity can be verified

**Why do you need it?**

JWT tokens contain user information (like user ID and role) and are used for authentication. The secret key:
- **Signs** the token when a user logs in
- **Verifies** the token on every authenticated request
- **Protects** against token manipulation by attackers

⚠️ **Security Note**: If someone gets your JWT secret, they can create fake authentication tokens and impersonate any user, including admins!

**Generate a secure random JWT secret:**

**Windows (PowerShell):**
```powershell
$bytes = New-Object byte[] 32; (New-Object Security.Cryptography.RNGCryptoServiceProvider).GetBytes($bytes); [Convert]::ToBase64String($bytes)
```

**Linux/Mac:**
```bash
openssl rand -base64 32
```

**Node.js:**
```bash
node -e "console.log(require('crypto').randomBytes(32).toString('base64'))"
```

Copy the generated string and use it as your `JWT_SECRET` in the `.env` file.

### 4. Configure environment variables

Copy `.env.example` to `.env` and update the values:

```bash
cp .env.example .env
```

Edit `.env` file:

```env
# Application
APP_ENV=development
PORT=8080

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=ecommerce_go

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
# IMPORTANT: Generate a secure random secret using the command above!
JWT_SECRET=your-generated-base64-secret-here-minimum-32-characters
ACCESS_TOKEN_MINUTES=15
REFRESH_TOKEN_DAYS=7

# AWS S3 / Cloudflare R2
S3_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com
S3_KEY=your-access-key-id
S3_SECRET=your-secret-access-key
S3_BUCKET=your-bucket-name
S3_REGION=auto
S3_PUBLIC_URL=https://pub-xxxxx.r2.dev
```

### 5. Create MySQL database
```sql
CREATE DATABASE ecommerce_go;
```

### 6. Start Redis server
```bash
redis-server
```

### 7. Run the application
```bash
go run ./cmd/api/main.go
```

The server will start on `http://localhost:8080`

### 8. Access API Documentation

Open your browser and navigate to:
```
http://localhost:8080/swagger/index.html
```

You'll see the interactive Swagger UI where you can:
- View all available endpoints
- Test API calls directly from the browser
- See request/response schemas
- Authenticate with JWT tokens

## 📚 API Documentation

This project includes comprehensive API documentation using Swagger/OpenAPI:

- **Swagger UI**: `http://localhost:8080/swagger/index.html` - Interactive API testing interface
- **OpenAPI JSON**: `http://localhost:8080/swagger/doc.json` - Machine-readable API specification
- **Quick Start Guide**: See `SWAGGER_QUICK_START.md` for step-by-step testing instructions
- **Full Documentation**: See `SWAGGER_DOCUMENTATION.md` for comprehensive API reference

### Using Swagger UI

1. **Navigate to** `http://localhost:8080/swagger/index.html`
2. **Create a user** using the `/auth/signup` endpoint
3. **Login** using the `/auth/login` endpoint and copy the `access_token`
4. **Click "Authorize"** button (🔓) at the top right
5. **Enter**: `Bearer <your-access-token>` and click "Authorize"
6. **Now you can test** all authenticated endpoints!

### Regenerating Swagger Documentation

If you modify API endpoints or add new ones:

```bash
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```

## 📡 API Endpoints

### Quick Reference

For detailed endpoint documentation with request/response schemas, visit the **Swagger UI** at:
```
http://localhost:8080/swagger/index.html
```

### Authentication (Public)
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/auth/signup` | Register new user | No |
| POST | `/auth/login` | User login | No |
| POST | `/auth/refresh` | Refresh access token | No |
| POST | `/api/auth/logout` | User logout | Yes |
| GET | `/api/profile` | Get user profile | Yes |

### Products
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/products` | List all products (pagination, filters) | No |
| GET | `/products/:slug` | Get product by slug (cached) | No |
| POST | `/api/admin/products/` | Create product | Admin |
| PUT | `/api/admin/products/:id` | Update product | Admin |
| DELETE | `/api/admin/products/:id` | Delete product (soft delete) | Admin |

### Shopping Cart
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/cart/add` | Add item to cart | Yes |
| GET | `/api/cart/` | Get user cart | Yes |
| PUT | `/api/cart/:id` | Update item quantity | Yes |
| DELETE | `/api/cart/:id` | Remove item from cart | Yes |
| DELETE | `/api/cart/` | Clear entire cart | Yes |

### Orders
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/orders/checkout` | Create order from cart | Yes |
| GET | `/api/orders/` | Get user orders | Yes |
| GET | `/api/orders/:id` | Get order details | Yes |

### Admin - Order Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/admin/orders` | Get all orders (filters, pagination) | Admin |
| GET | `/api/admin/orders/:id` | Get specific order | Admin |
| PUT | `/api/admin/orders/:id/status` | Update order status | Admin |
| GET | `/api/admin/orders/stats` | Get order statistics & revenue | Admin |

### Payment
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/payment/start/:orderId` | Start payment process | Yes |
| POST | `/api/payment/webhook` | Payment webhook handler | No |

### File Upload
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/upload/product` | Upload product image to S3 | Admin |

### Health & Documentation
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/health` | Health check | No |
| GET | `/swagger/*any` | Swagger UI & API docs | No |

## 🔐 Authentication Flow

1. **Sign Up**: User registers with name, email, password, and role (customer/admin)
   ```json
   POST /auth/signup
   {
     "name": "John Doe",
     "email": "john@example.com",
     "password": "password123",
     "role": "customer"
   }
   ```

2. **Login**: Returns access token (JWT) and refresh token (HTTP-only cookie)
   ```json
   POST /auth/login
   {
     "email": "john@example.com",
     "password": "password123"
   }
   
   Response:
   {
     "access_token": "eyJhbGc...",
     "expires_in": 900,
     "user": {
       "id": 1,
       "email": "john@example.com",
       "role": "customer",
       "name": "John Doe"
     }
   }
   ```

3. **Access Protected Routes**: Include access token in `Authorization: Bearer <token>` header

4. **Token Refresh**: Use `/auth/refresh` endpoint with refresh token cookie to get new access token
   - Refresh token is automatically rotated for security
   - Refresh tokens are stored hashed in the database

5. **Logout**: Clears refresh token from database and removes cookie
   ```json
   POST /api/auth/logout
   Authorization: Bearer <token>
   ```

## 💾 Database Models

### User
- **Fields**: ID, Name, Email, Password (hashed), Role (admin/customer), RefreshToken, TokenExpiry
- **Timestamps**: CreatedAt, UpdatedAt, DeletedAt
- **Relations**: Has many CartItems, Has many Orders

### Product
- **Fields**: ID, Name, Slug, Description, Price, Quantity (stock), Category, ImageURL, IsActive
- **Timestamps**: CreatedAt, UpdatedAt, DeletedAt (soft delete)
- **Relations**: Used in CartItems, Used in OrderItems
- **Features**: Auto-generated slug, Redis caching

### CartItem
- **Fields**: ID, UserID, ProductID, Quantity
- **Timestamps**: CreatedAt, UpdatedAt
- **Relations**: Belongs to User, Belongs to Product
- **Constraints**: Stock validation on add/update

### Order
- **Fields**: ID, UserID, TotalPrice, Status (pending/confirmed/shipped/delivered/cancelled)
- **Timestamps**: CreatedAt, UpdatedAt, DeletedAt
- **Relations**: Belongs to User, Has many OrderItems
- **Features**: Atomic creation with transaction, Stock deduction

### OrderItem
- **Fields**: ID, OrderID, ProductID, Quantity, Price (at order time), Subtotal
- **Timestamps**: CreatedAt, UpdatedAt
- **Relations**: Belongs to Order, References Product
- **Purpose**: Preserves product price at time of purchase

### PaymentIntent
- **Fields**: ID, OrderID, Amount, Status (pending/paid/failed), GatewayRef
- **Timestamps**: CreatedAt, UpdatedAt
- **Relations**: Links to Order
- **Features**: Simulated payment gateway integration

## 🚦 Status Codes

- `200 OK` - Successful request
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## 🧪 Example Requests

### 1. Register User
```bash
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "user@example.com",
    "password": "SecurePass123",
    "role": "customer"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123"
  }'
```

### 3. Create Product (Admin)
```bash
curl -X POST http://localhost:8080/api/admin/products/ \
  -H "Authorization: Bearer <your-access-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Gaming Laptop",
    "description": "High-performance gaming laptop",
    "price": 1999.99,
    "quantity": 10,
    "category": "Electronics",
    "image_url": "https://example.com/laptop.jpg"
  }'
```

### 4. List Products with Filters
```bash
# Get Electronics, page 1, 10 items
curl -X GET "http://localhost:8080/products?category=Electronics&page=1&limit=10"

# Search for "laptop"
curl -X GET "http://localhost:8080/products?search=laptop"
```

### 5. Add to Cart
```bash
curl -X POST http://localhost:8080/api/cart/add \
  -H "Authorization: Bearer <your-access-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity": 2
  }'
```

### 6. View Cart
```bash
curl -X GET http://localhost:8080/api/cart/ \
  -H "Authorization: Bearer <your-access-token>"
```

### 7. Checkout
```bash
curl -X POST http://localhost:8080/api/orders/checkout \
  -H "Authorization: Bearer <your-access-token>"
```

### 8. Upload Product Image
```bash
curl -X POST http://localhost:8080/api/upload/product \
  -H "Authorization: Bearer <admin-access-token>" \
  -F "image=@/path/to/image.jpg"
```

### 9. Update Order Status (Admin)
```bash
curl -X PUT http://localhost:8080/api/admin/orders/1/status \
  -H "Authorization: Bearer <admin-access-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "confirmed"
  }'
```

### 10. Get Order Statistics (Admin)
```bash
curl -X GET "http://localhost:8080/api/admin/orders/stats?from=2024-01-01&to=2024-12-31" \
  -H "Authorization: Bearer <admin-access-token>"
```

## 💡 Pro Tips for Testing

1. **Use Swagger UI**: The easiest way to test all endpoints interactively
   - Visit: `http://localhost:8080/swagger/index.html`
   - All examples and schemas are provided

2. **Postman Collection**: Import the OpenAPI spec from `/swagger/doc.json` into Postman

3. **Create Admin User**: Set `"role": "admin"` in signup request to test admin endpoints

4. **Test Complete Flow**:
   - Sign up → Login → Create products (admin) → Add to cart → Checkout → View orders

## 🔒 Security Features

- **Password hashing** with bcrypt (cost factor 10)
- **JWT-based authentication** with configurable expiration
- **HTTP-only cookies** for refresh tokens (prevents XSS attacks)
- **Role-based access control** (customer/admin roles)
- **Token rotation** on refresh (prevents token reuse)
- **SQL injection prevention** (GORM prepared statements)
- **Input validation** on all endpoints
- **Secure file upload** with type validation and size limits
- **Middleware-based authentication** for protected routes
- **Rate limiting** on sensitive endpoints
- **CORS configuration** for production deployment

## 🚀 Performance Optimizations

- **Redis caching** for frequently accessed products (5-minute TTL)
- **Cache invalidation** on product updates
- **Database indexing** on frequently queried fields (email, slug)
- **Pagination** for large datasets
- **Efficient SQL queries** with GORM eager loading
- **Database transactions** for data consistency
- **Connection pooling** for database connections

## 📊 Order Status Workflow

Orders follow a strict status transition workflow:

```
pending → confirmed → shipped → delivered
    ↓
cancelled (can cancel from pending or confirmed)
```

**Status Transitions:**
- `pending` → Can change to: `confirmed`, `cancelled`
- `confirmed` → Can change to: `shipped`, `cancelled`
- `shipped` → Can change to: `delivered`
- `delivered` → Final state (no changes)
- `cancelled` → Final state (no changes)

**Stock Management:**
- Stock is deducted when order is created (status: pending)
- Stock is restored if order is cancelled
- Admin can track order statistics and revenue by date range

## 🌐 AWS S3 / Cloudflare R2 Setup

### Option 1: AWS S3

1. Create an AWS account and navigate to S3
2. Create a new S3 bucket
3. Configure bucket permissions for public access (if needed)
4. Create IAM user with S3 access permissions
5. Generate Access Key ID and Secret Access Key
6. Update `.env` file with S3 credentials

### Option 2: Cloudflare R2

1. Create a Cloudflare account and enable R2
2. Create a new R2 bucket
3. Generate API tokens (Access Key ID and Secret Access Key)
4. Enable public access on the bucket to get the public URL
5. Add credentials to `.env` file

**Environment Variables:**
```env
S3_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com
S3_KEY=your-access-key-id
S3_SECRET=your-secret-access-key
S3_BUCKET=your-bucket-name
S3_REGION=auto
S3_PUBLIC_URL=https://pub-xxxxx.r2.dev
```

For detailed setup instructions, see `S3_STORAGE_GUIDE.md`

## 🧪 Testing the API

### Using Swagger UI (Recommended)

1. Start the server: `go run ./cmd/api/main.go`
2. Open browser: `http://localhost:8080/swagger/index.html`
3. Follow the testing guide in `SWAGGER_QUICK_START.md`

### Using cURL

See the **Example Requests** section above for cURL commands.

### Using Postman

1. Import OpenAPI specification from `http://localhost:8080/swagger/doc.json`
2. Set up environment variables for base URL and tokens
3. Test endpoints interactively

## 📖 Additional Documentation

- **`SWAGGER_DOCUMENTATION.md`** - Comprehensive API reference with all endpoints, request/response schemas, and examples
- **`SWAGGER_QUICK_START.md`** - Step-by-step guide for testing the API using Swagger UI
- **`S3_STORAGE_GUIDE.md`** - Detailed guide for setting up AWS S3 or Cloudflare R2
- **Swagger UI** - Interactive API documentation at `/swagger/index.html`

## 🛠️ Development Workflow

### Making Changes to the API

1. **Modify controller** - Update or add endpoints in `internal/controllers/`
2. **Add Swagger annotations** - Document the endpoint with proper Swagger comments
3. **Update routes** - Register the route in `internal/routes/`
4. **Regenerate docs** - Run: `swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal`
5. **Test** - Use Swagger UI to test your changes
6. **Commit** - Commit your changes including generated docs

### Project Best Practices

- **Always use transactions** for operations that modify multiple tables
- **Validate input** before database operations
- **Check stock availability** before adding to cart or creating orders
- **Invalidate cache** when updating products
- **Use proper HTTP status codes** for responses
- **Log errors** for debugging
- **Follow RESTful principles** for endpoint design

## 📝 License

This project is open source and available under the MIT License.

## 👨‍💻 Developer

Created as part of the Go learning journey - Building production-ready RESTful APIs.

## 🎯 Learning Outcomes

This project demonstrates:
- ✅ Building RESTful APIs with Gin framework
- ✅ JWT authentication and authorization
- ✅ Database design and GORM ORM
- ✅ Redis caching for performance
- ✅ File uploads to cloud storage (S3/R2)
- ✅ API documentation with Swagger/OpenAPI
- ✅ Middleware implementation (auth, admin, rate limiting)
- ✅ Transaction management for data integrity
- ✅ Error handling and validation
- ✅ Clean architecture and project structure

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 🐛 Known Issues & Future Enhancements

### Known Issues
- None at the moment

### Planned Features
- [ ] Email notifications for orders
- [ ] Product reviews and ratings
- [ ] Wishlist functionality
- [ ] Advanced search with Elasticsearch
- [ ] Real payment gateway integration (Stripe/PayPal)
- [ ] Order tracking with shipping providers
- [ ] Inventory management system
- [ ] Discount codes and promotions
- [ ] User address management
- [ ] Multiple image uploads per product

## 📧 Support

For support:
- 📖 Read the documentation in `SWAGGER_DOCUMENTATION.md`
- 🧪 Try the Swagger UI for interactive testing
- 📝 Check `SWAGGER_QUICK_START.md` for testing examples
- 🐛 Create an issue in the repository

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin) - Fast HTTP web framework
- [GORM](https://gorm.io/) - Fantastic ORM library for Golang
- [Swaggo](https://github.com/swaggo/swag) - Automatically generate RESTful API documentation
- [Redis](https://redis.io/) - In-memory data structure store
- [JWT-Go](https://github.com/golang-jwt/jwt) - JSON Web Token implementation

---

**Happy Coding! 🚀**

**Don't forget to ⭐ this repository if you found it helpful!**
