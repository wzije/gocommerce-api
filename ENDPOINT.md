## API Endpoints

### 1. User Registration

**Endpoint**: `POST /register`

**Description**: Register a new user.

**Request Headers**:
- Content-Type: `application/json`

**Request Body**:
```json
{
  "email": "testing@gmail.com",
  "password": "localhost"
}
```

**Response**:
- 201 Created: Returns the newly created user's details.

### 2. User Login

**Endpoint**: `POST /login`

**Description**: Log in a user and retrieve an access token.

**Request Headers**:
- Content-Type: `application/json`

**Request Body**:
```json
{
  "email": "{{email}}",
  "password": "{{password}}"
}
```

**Response**:
- 200 OK: Returns the access token and user details.

### 3. Get My Profile

**Endpoint**: `GET /me`

**Description**: Retrieve the authenticated user's profile information.

**Request Headers**:
- Content-Type: `application/json`
- Authorization: `Bearer {{token}}`

**Response**:
- 200 OK: Returns the user's profile details.

### 4. Get Products

**Endpoint**: `GET /products`

**Description**: Retrieve a list of products.

**Request Headers**:
- Content-Type: `application/json`
- Authorization: `Bearer {{token}}`

**Response**:
- 200 OK: Returns a list of products.

### 5. Get Product Details

**Endpoint**: `GET /products/{id}`

**Description**: Retrieve details of a specific product by its ID.

**Request Headers**:
- Content-Type: `application/json`
- Authorization: `Bearer {{token}}`

**Path Parameters**:
- `id`: The ID of the product.

**Response**:
- 200 OK: Returns the product details.

### 6. Get My Orders

**Endpoint**: `GET /orders`

**Description**: Retrieve a list of orders placed by the authenticated user.

**Request Headers**:
- Content-Type: `application/json`
- Authorization: `Bearer {{token}}`

**Response**:
- 200 OK: Returns a list of the user's orders.

### 7. Checkout / Create Order

**Endpoint**: `POST /orders/checkout`

**Description**: Create a new order (checkout process).

**Request Headers**:
- Content-Type: `application/json`
- Authorization: `Bearer {{token}}`

**Request Body**:
```json
{
  "shop_id": 1,
  "shipping_address": "Yogyakarta",
  "details": [
    {
      "product_id": 1,
      "quantity": 1
    }
  ]
}
```

**Response**:
- 201 Created: Returns the created order's details.

### 8. Create Payment for Order

**Endpoint**: `POST /orders/payment`

**Description**: Create a payment for an existing order.

**Request Headers**:
- Content-Type: `application/json`
- Authorization: `Bearer {{token}}`

**Request Body**:
```json
{
  "order_id": 5,
  "payment_method": "transfer",
  "amount": 10000
}
```

**Response**:
- 201 Created: Returns the payment details.

### 9. Get My Warehouses

**Endpoint**: `GET /warehouses`

**Description**: Retrieve a list of warehouses associated with the authenticated user.

**Request Headers**:
- Content-Type: `application/json`
- Authorization: `Bearer {{token}}`

**Response**:
- 200 OK: Returns a list of warehouses.

### 10. Increase Product Stock in a Specific Warehouse

**Endpoint**: `POST /warehouses/increase`

**Description**: Increase the stock of a specific product in a specific warehouse.

**Request Headers**:
- Content-Type: `application/json`
- Authorization: `Bearer {{token}}`

**Request Body**:
```json
{
  "product_id": 5,
  "warehouse_id": 1,
  "quantity": 10
}
```

**Response**:
- 200 OK: Returns a confirmation of the stock increase.

### 11. Decrease Product Stock in a Specific Warehouse

**Endpoint**: `POST /warehouses/reduce`

**Description**: Decrease the stock of a specific product in a specific warehouse.

**Request Headers**:
- Content-Type: `application/json`
- Authorization: `Bearer {{token}}`

**Request Body**:
```json
{
  "product_id": 31,
  "warehouse_id": 13,
  "quantity": 10
}
```

**Response**:
- 200 OK: Returns a confirmation of the stock decrease.

### 12. Transfer Product Stock Between Warehouses

**Endpoint**: `POST /warehouses/transfer`

**Description**: Transfer stock of a product from one warehouse to another.

**Request Headers**:
- Content-Type: `application/json`
- Authorization: `Bearer {{token}}`

**Request Body**:
```json
{
  "source_warehouse_id": 1,
  "target_warehouse_id": 2,
  "quantity": 30,
  "product_id": 3
}
```

**Response**:
- 200 OK: Returns a confirmation of the stock transfer.

### 13. Change Warehouse Status

**Endpoint**: `POST /warehouses/{id}/status`

**Description**: Change the active status of a warehouse.

**Request Headers**:
- Content-Type: `application/json`
- Authorization: `Bearer {{token}}`

**Path Parameters**:
- `id`: The ID of the warehouse.

**Query Parameters**:
- `is_active`: Set to `true` to activate the warehouse or `false` to deactivate it.

**Response**:
- 200 OK: Returns a confirmation of the status change.

---



