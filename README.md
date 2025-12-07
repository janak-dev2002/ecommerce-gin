# E-Commerce Backend - Gin Framework

A comprehensive RESTful API backend for an e-commerce platform built with Go, Gin web framework, and MySQL database. Features include user authentication, product management, shopping cart, order processing, payment integration, and cloud storage for product images.

## 🚀 Features

### Authentication & Authorization
- **User Registration & Login** with JWT-based authentication
- **Access & Refresh Token** mechanism for secure session management
- **Role-based Access Control** (Admin/User roles)
- **Password Hashing** using bcrypt

### Product Management
- Product CRUD operations (Create, Read, Update, Delete)
- Product listing with pagination
- Product search by slug
- Admin-only product management endpoints
- Image upload to Cloudflare R2 storage

### Shopping Cart
- Add products to cart
- Update cart item quantities
- Remove items from cart
- Clear entire cart
- View cart with calculated totals

### Order Management
- Checkout process with cart validation
- Order creation and tracking
- Order history for users
- Admin order management (view all orders, update status)
- Order status tracking (pending, processing, shipped, delivered, cancelled)

### Payment Integration
- Payment intent creation
- Simulated payment gateway
- Payment webhook handling
- Payment status tracking

### File Upload
- Product image upload to Cloudflare R2
- Support for multiple image formats (JPEG, PNG, WebP, GIF)
- Automatic content-type detection
- Public URL generation for uploaded images

## 🛠️ Tech Stack

- **Language**: Go 1.25.3
- **Web Framework**: Gin v1.11.0
- **Database**: MySQL with GORM v1.31.1
- **Authentication**: JWT (golang-jwt/jwt/v4)
- **Cloud Storage**: Cloudflare R2 (AWS S3 compatible)
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **Environment Config**: godotenv

## 📁 Project Structure

```
ecommerce-gin/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
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
│   │   └── upload_controller.go
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
│   │   └── admin.go             # Admin authorization middleware
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
├── .env                         # Environment variables
├── go.mod
└── README.md
```

## 🔧 Installation & Setup

### Prerequisites
- Go 1.25 or higher
- MySQL 8.0 or higher
- Cloudflare R2 account (for image uploads)

### 1. Clone the repository
```bash
git clone <repository-url>
cd ecommerce-gin
```

### 2. Install dependencies
```bash
go mod download
```

### 3. Generate JWT Secret

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

# JWT Configuration
# IMPORTANT: Generate a secure random secret using the command above!
JWT_SECRET=your-generated-base64-secret-here-minimum-32-characters
ACCESS_TOKEN_MINUTES=15
REFRESH_TOKEN_DAYS=7

# Cloudflare R2 / S3
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

### 6. Run the application
```bash
go run ./cmd/api/main.go
```

The server will start on `http://localhost:8080`

## 📡 API Endpoints

### Authentication
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
| GET | `/api/products` | List all products | No |
| GET | `/api/products/:slug` | Get product by slug | No |
| POST | `/api/admin/products` | Create product | Admin |
| PUT | `/api/admin/products/:id` | Update product | Admin |
| DELETE | `/api/admin/products/:id` | Delete product | Admin |

### Shopping Cart
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/cart/add` | Add item to cart | Yes |
| GET | `/api/cart` | Get user cart | Yes |
| PUT | `/api/cart/:id` | Update item quantity | Yes |
| DELETE | `/api/cart/:id` | Remove item from cart | Yes |
| DELETE | `/api/cart` | Clear cart | Yes |

### Orders
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/orders/checkout` | Create order from cart | Yes |
| GET | `/api/orders` | Get user orders | Yes |
| GET | `/api/orders/:id` | Get order details | Yes |
| GET | `/api/admin/orders` | Get all orders | Admin |
| PUT | `/api/admin/orders/:id/status` | Update order status | Admin |

### Payment
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/payment/create-intent` | Create payment intent | Yes |
| POST | `/webhook/payment` | Payment webhook | No |

### File Upload
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/upload/product` | Upload product image | Admin |

### Health Check
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/health` | Health check | No |

## 🔐 Authentication Flow

1. **Sign Up**: User registers with email and password
2. **Login**: Returns access token (JWT) and refresh token (HTTP-only cookie)
3. **Access Protected Routes**: Include access token in `Authorization: Bearer <token>` header
4. **Token Refresh**: Use `/auth/refresh` endpoint with refresh token cookie to get new access token
5. **Logout**: Clears refresh token cookie

## 💾 Database Models

### User
- ID, Email, Password (hashed), Role (admin/user)
- Timestamps (CreatedAt, UpdatedAt)

### Product
- ID, Name, Slug, Description, Price, Stock, ImageURL
- Timestamps (CreatedAt, UpdatedAt)

### CartItem
- ID, UserID, ProductID, Quantity
- Relationships: User, Product

### Order
- ID, UserID, TotalAmount, Status, PaymentIntentID
- Timestamps (CreatedAt, UpdatedAt)
- Relationships: User, OrderItems

### OrderItem
- ID, OrderID, ProductID, Quantity, PriceAtOrder

### PaymentIntent
- ID, OrderID, Amount, Status, ExternalID
- Timestamps (CreatedAt, UpdatedAt)

## 🚦 Status Codes

- `200 OK` - Successful request
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## 🧪 Example Requests

### Register User
```bash
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123"
  }'
```

### Add to Cart
```bash
curl -X POST http://localhost:8080/api/cart/add \
  -H "Authorization: Bearer <your-access-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity": 2
  }'
```

### Upload Product Image
```bash
curl -X POST http://localhost:8080/api/upload/product \
  -H "Authorization: Bearer <admin-access-token>" \
  -F "image=@/path/to/image.jpg"
```

## 🔒 Security Features

- Password hashing with bcrypt
- JWT-based authentication
- HTTP-only cookies for refresh tokens
- Role-based access control
- SQL injection prevention (GORM prepared statements)
- Input validation
- Secure file upload with type validation

## 🌐 Cloudflare R2 Setup

1. Create a Cloudflare account and enable R2
2. Create a new R2 bucket
3. Generate API tokens (Access Key ID and Secret Access Key)
4. Enable public access on the bucket to get the public URL
5. Add credentials to `.env` file

## 📝 License

This project is open source and available under the MIT License.

## 👨‍💻 Developer

Created as part of the Go learning journey.

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!

## 📧 Support

For support, email your-email@example.com or create an issue in the repository.

---

**Happy Coding! 🚀**
