# Quick Start Guide - Testing Swagger Documentation

## Step 1: Start the Server

Make sure your application is running:
```bash
go run ./cmd/api
```

You should see output similar to:
```
Server running on port 8080
[GIN-debug] Listening and serving HTTP on :8080
```

## Step 2: Access Swagger UI

Open your web browser and navigate to:
```
http://localhost:8080/swagger/index.html
```

## Step 3: Test the API

### A. Create a User Account

1. Find the **Auth** section in Swagger UI
2. Click on `POST /auth/signup`
3. Click "Try it out"
4. Modify the request body:
```json
{
  "name": "Test User",
  "email": "test@example.com",
  "password": "password123",
  "role": "customer"
}
```
5. Click "Execute"
6. You should see a 201 response with `{"message": "user created"}`

### B. Login

1. Click on `POST /auth/login`
2. Click "Try it out"
3. Enter credentials:
```json
{
  "email": "test@example.com",
  "password": "password123"
}
```
4. Click "Execute"
5. Copy the `access_token` from the response

### C. Authorize Swagger

1. Click the **"Authorize"** button (🔓) at the top right of Swagger UI
2. In the "Value" field, enter: `Bearer <paste-your-access-token-here>`
3. Click "Authorize"
4. Click "Close"

Now you're authenticated and can test protected endpoints!

### D. Test Protected Endpoints

#### Create a Product (Admin Only)

**Note:** First create an admin user with role "admin", then login and authorize.

1. Find `POST /api/admin/products/` under **Products**
2. Click "Try it out"
3. Enter product data:
```json
{
  "name": "Gaming Laptop",
  "description": "High-performance gaming laptop with RTX 4080",
  "price": 1999.99,
  "quantity": 5,
  "category": "Electronics",
  "image_url": "https://example.com/laptop.jpg"
}
```
4. Click "Execute"

#### List Products (Public)

1. Find `GET /products` under **Products**
2. Click "Try it out"
3. Optionally set query parameters:
   - page: 1
   - limit: 10
   - category: Electronics
4. Click "Execute"

#### Add to Cart

1. Find `POST /api/cart/add` under **Cart**
2. Click "Try it out"
3. Enter:
```json
{
  "product_id": 1,
  "quantity": 2
}
```
4. Click "Execute"

#### View Cart

1. Find `GET /api/cart/` under **Cart**
2. Click "Try it out"
3. Click "Execute"

#### Checkout

1. Find `POST /api/orders/checkout` under **Orders**
2. Click "Try it out"
3. Click "Execute"
4. Note the `order_id` in the response

#### View Orders

1. Find `GET /api/orders/` under **Orders**
2. Click "Try it out"
3. Click "Execute"

## Step 4: Testing Admin Features

### Create an Admin User

1. Use `POST /auth/signup` with:
```json
{
  "name": "Admin User",
  "email": "admin@example.com",
  "password": "admin123",
  "role": "admin"
}
```

2. Login with admin credentials
3. Authorize with the admin access token

### Test Admin Endpoints

#### List All Orders
- `GET /api/admin/orders` - View all orders with filters

#### Update Order Status
- `PUT /api/admin/orders/{id}/status`
```json
{
  "status": "confirmed"
}
```

#### View Order Statistics
- `GET /api/admin/orders/stats` - See order counts and revenue

## Available Endpoints Reference

### 🔓 Public Endpoints (No Authentication Required)
- POST `/auth/signup` - Register
- POST `/auth/login` - Login
- POST `/auth/refresh` - Refresh token
- GET `/products` - List products
- GET `/products/{slug}` - Get product by slug
- GET `/health` - Health check

### 🔒 Authenticated Endpoints (User)
- POST `/api/auth/logout` - Logout
- GET `/api/profile` - Get profile
- POST `/api/cart/add` - Add to cart
- GET `/api/cart/` - View cart
- PUT `/api/cart/{id}` - Update cart item
- DELETE `/api/cart/{id}` - Remove from cart
- DELETE `/api/cart/` - Clear cart
- POST `/api/orders/checkout` - Checkout
- GET `/api/orders/` - My orders
- GET `/api/orders/{id}` - Order details
- POST `/api/payment/start/{orderId}` - Start payment

### 👑 Admin Only Endpoints
- POST `/api/admin/products/` - Create product
- PUT `/api/admin/products/{id}` - Update product
- DELETE `/api/admin/products/{id}` - Delete product
- GET `/api/admin/orders` - List all orders
- GET `/api/admin/orders/{id}` - Get order
- PUT `/api/admin/orders/{id}/status` - Update order status
- GET `/api/admin/orders/stats` - Order statistics
- POST `/api/upload/product` - Upload image to S3

## Troubleshooting

### Issue: "401 Unauthorized"
**Solution:** Make sure you've authorized with a valid Bearer token

### Issue: "403 Forbidden" on admin endpoints
**Solution:** Make sure you're logged in with an admin account (role: "admin")

### Issue: Swagger UI not loading
**Solution:** 
1. Make sure the server is running
2. Check you're accessing `http://localhost:8080/swagger/index.html`
3. Regenerate docs: `swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal`

### Issue: "product not found or inactive"
**Solution:** Create a product first using the admin create product endpoint

### Issue: "cart is empty" when checking out
**Solution:** Add products to cart first using the add to cart endpoint

## Tips

1. **Save your tokens**: Copy access tokens to a text file for easy reuse
2. **Use realistic data**: Use valid email formats and reasonable prices
3. **Check responses**: Always review the response to understand what happened
4. **Follow the workflow**: 
   - Create products (admin) → Add to cart → Checkout → View orders
5. **Test error cases**: Try invalid inputs to see error responses

## Next Steps

- Explore all endpoints in the Swagger UI
- Test different query parameters for listing endpoints
- Try the payment workflow
- Upload product images
- Test order status transitions
- Review the full API documentation in `SWAGGER_DOCUMENTATION.md`

## Support

If you encounter issues, check:
1. Server logs in the terminal
2. Browser console for client-side errors
3. Response bodies for error messages
4. The comprehensive documentation in `SWAGGER_DOCUMENTATION.md`
