# E-Commerce API - Swagger Documentation

## Overview
This document provides information about the Swagger documentation for the E-Commerce Backend API built with Gin framework.

## Accessing Swagger UI

Once the application is running, you can access the Swagger documentation at:

```
http://localhost:8080/swagger/index.html
```

## API Endpoints Summary

### Authentication Endpoints (Public)
- **POST** `/auth/signup` - Register a new user
- **POST** `/auth/login` - Login and receive access token
- **POST** `/auth/refresh` - Refresh access token using refresh cookie
- **POST** `/api/auth/logout` - Logout user (requires authentication)

### Product Endpoints
**Public:**
- **GET** `/products` - List all products (with pagination and filtering)
- **GET** `/products/{slug}` - Get product by slug

**Admin Only:**
- **POST** `/api/admin/products/` - Create a new product
- **PUT** `/api/admin/products/{id}` - Update product details
- **DELETE** `/api/admin/products/{id}` - Delete (soft delete) a product

### Cart Endpoints (Authenticated)
- **POST** `/api/cart/add` - Add item to cart
- **GET** `/api/cart/` - Get user's cart
- **PUT** `/api/cart/{id}` - Update cart item quantity
- **DELETE** `/api/cart/{id}` - Remove item from cart
- **DELETE** `/api/cart/` - Clear entire cart

### Order Endpoints (Authenticated)
- **POST** `/api/orders/checkout` - Create order from cart
- **GET** `/api/orders/` - Get user's orders
- **GET** `/api/orders/{id}` - Get order details

### Admin Order Endpoints (Admin Only)
- **GET** `/api/admin/orders` - List all orders (with filters and pagination)
- **GET** `/api/admin/orders/{id}` - Get specific order by ID
- **PUT** `/api/admin/orders/{id}/status` - Update order status
- **GET** `/api/admin/orders/stats` - Get order statistics and revenue

### Payment Endpoints (Authenticated)
- **POST** `/api/payment/start/{orderId}` - Start payment process for an order
- **POST** `/api/payment/webhook` - Payment gateway webhook handler

### Upload Endpoints (Admin Only)
- **POST** `/api/upload/product` - Upload product image to AWS S3

### Health Check (Public)
- **GET** `/health` - Health check endpoint

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Most endpoints require authentication.

### Using Authentication in Swagger UI

1. First, call the `/auth/signup` endpoint to create a user account
2. Then call `/auth/login` with your credentials
3. Copy the `access_token` from the response
4. Click the **"Authorize"** button at the top of the Swagger UI
5. Enter: `Bearer <your-access-token>` (replace `<your-access-token>` with your actual token)
6. Click "Authorize" and then "Close"
7. Now you can test authenticated endpoints

**Example:**
```
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## User Roles

The API supports two user roles:
- **customer** (default) - Regular users who can shop and place orders
- **admin** - Administrators who can manage products and orders

### Creating an Admin User

To create an admin user, use the `/auth/signup` endpoint with the role field:

```json
{
  "name": "Admin User",
  "email": "admin@example.com",
  "password": "securepassword",
  "role": "admin"
}
```

## Request/Response Examples

### 1. User Signup
**Request:**
```json
POST /auth/signup
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "role": "customer"
}
```

**Response:**
```json
{
  "message": "user created"
}
```

### 2. User Login
**Request:**
```json
POST /auth/login
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "user": {
    "id": 1,
    "email": "john@example.com",
    "role": "customer",
    "name": "John Doe"
  }
}
```

### 3. Create Product (Admin)
**Request:**
```json
POST /api/admin/products/
Authorization: Bearer <token>

{
  "name": "Laptop",
  "description": "High-performance laptop",
  "price": 999.99,
  "quantity": 10,
  "category": "Electronics",
  "image_url": "https://example.com/laptop.jpg"
}
```

**Response:**
```json
{
  "id": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "name": "Laptop",
  "slug": "laptop",
  "description": "High-performance laptop",
  "price": 999.99,
  "quantity": 10,
  "category": "Electronics",
  "image_url": "https://example.com/laptop.jpg",
  "is_active": true
}
```

### 4. Add to Cart
**Request:**
```json
POST /api/cart/add
Authorization: Bearer <token>

{
  "product_id": 1,
  "quantity": 2
}
```

**Response:**
```json
{
  "message": "added to cart"
}
```

### 5. Checkout
**Request:**
```json
POST /api/orders/checkout
Authorization: Bearer <token>
```

**Response:**
```json
{
  "message": "order placed",
  "order_id": 1,
  "total": 1999.98
}
```

## Query Parameters

### Product Listing
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10)
- `search` - Search term for product name
- `category` - Filter by category

Example: `/products?page=1&limit=10&category=Electronics&search=laptop`

### Admin Order Listing
- `status` - Filter by order status (pending, confirmed, shipped, delivered, cancelled)
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 20)
- `from` - Start date (YYYY-MM-DD)
- `to` - End date (YYYY-MM-DD)

Example: `/api/admin/orders?status=pending&page=1&limit=20&from=2024-01-01&to=2024-12-31`

## Order Status Workflow

Orders follow this status workflow:
1. **pending** → Can transition to: confirmed, cancelled
2. **confirmed** → Can transition to: shipped, cancelled
3. **shipped** → Can transition to: delivered
4. **delivered** → Final state
5. **cancelled** → Final state

## Error Handling

All error responses follow this format:
```json
{
  "error": "error message description"
}
```

Common HTTP status codes:
- **200** - Success
- **201** - Created
- **400** - Bad Request (validation errors)
- **401** - Unauthorized (missing or invalid token)
- **403** - Forbidden (insufficient permissions)
- **404** - Not Found
- **422** - Unprocessable Entity
- **500** - Internal Server Error

## Rate Limiting

The `/products/{slug}` endpoint has rate limiting enabled (2 requests per time window) to demonstrate rate limiting functionality.

## File Upload

The image upload endpoint accepts multipart/form-data with the following specifications:
- **Field name**: `image`
- **Allowed formats**: .jpg, .jpeg, .png, .webp
- **Storage**: AWS S3

## Refresh Token

Refresh tokens are stored as HTTP-only cookies for security. The `/auth/refresh` endpoint:
- Reads the refresh token from the cookie
- Validates and rotates the refresh token
- Returns a new access token

## Development Setup

1. Install swaggo/swag:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Generate Swagger docs:
```bash
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```

3. Run the application:
```bash
go run ./cmd/api
```

4. Access Swagger UI:
```
http://localhost:8080/swagger/index.html
```

## Regenerating Documentation

If you make changes to the API annotations, regenerate the docs:

```bash
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```

## Additional Resources

- Swagger/OpenAPI Specification: https://swagger.io/specification/
- Swaggo Documentation: https://github.com/swaggo/swag
- Gin Framework: https://github.com/gin-gonic/gin

## Notes

- All timestamps are in RFC3339 format (ISO 8601)
- Prices are in decimal format (e.g., 999.99)
- The API uses soft deletes for products (sets deleted_at timestamp)
- Redis caching is implemented for product retrieval (5-minute TTL)
- Database transactions are used for checkout and order status updates
- Stock quantities are automatically updated during checkout
- Stock is restored when orders are cancelled

## Support

For API support, please contact: support@ecommerce.com
