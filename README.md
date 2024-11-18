# Product Service for E-Commerce Application

[![GoDoc](https://pkg.go.dev/badge/github.com/tittuvarghese/ss-go-product-service)](https://pkg.go.dev/github.com/tittuvarghese/ss-go-product-service)
[![Build Status](https://travis-ci.org/tittuvarghese/ss-go-product-service.svg?branch=main)](https://travis-ci.org/tittuvarghese/ss-go-product-service)

The **Product Service** is part of the e-commerce microservices architecture. It manages products in the system and provides CRUD functionality such as creating, updating, and retrieving product details. The service is built using gRPC for efficient communication between microservices.

This service handles product-related operations, such as managing product information, including its ID, name, price, dimensions, and other attributes.

## API Overview

The **Product Service** exposes the following gRPC methods for managing products:

### 1. **Create Product**
- **RPC Method**: `CreateProduct`
- **Request Type**: `CreateProductRequest`
- **Response Type**: `CreateProductResponse`
- **Description**: Creates a new product in the system.

#### Request (CreateProductRequest)
```proto
message CreateProductRequest {
  Product product = 1; // The product to create
}
```

#### Response (CreateProductResponse)
```proto
message CreateProductResponse {
  string message = 1; // Success or failure message
}
```

### 2. **Get Product**
- **RPC Method**: `GetProduct`
- **Request Type**: `GetProductRequest`
- **Response Type**: `GetProductResponse`
- **Description**: Retrieves a single product by its product ID.

#### Request (GetProductRequest)
```proto
message GetProductRequest {
  string product_id = 1; // UUID of the product to retrieve
}
```

#### Response (GetProductResponse)
```proto
message GetProductResponse {
  string message = 1;   // Success or failure message
  Product product = 2;  // The product object
}
```

### 3. **Get Multiple Products**
- **RPC Method**: `GetProducts`
- **Request Type**: `GetProductsRequest`
- **Response Type**: `GetProductsResponse`
- **Description**: Retrieves multiple products by their product IDs.

#### Request (GetProductsRequest)
```proto
message GetProductsRequest {
  repeated string query = 1; // List of product IDs to retrieve
}
```

#### Response (GetProductsResponse)
```proto
message GetProductsResponse {
  string message = 1;    // Success or failure message
  repeated Product products = 2;  // List of products
}
```

### 4. **Update Product**
- **RPC Method**: `UpdateProduct`
- **Request Type**: `UpdateProductRequest`
- **Response Type**: `UpdateProductResponse`
- **Description**: Updates an existing product in the system.

#### Request (UpdateProductRequest)
```proto
message UpdateProductRequest {
  string product_id = 1; // UUID of the product to update
  Product product = 2;   // Updated product information
}
```

#### Response (UpdateProductResponse)
```proto
message UpdateProductResponse {
  string message = 1;  // Success or failure message
}
```

## Product Message Definition

The **Product** message structure contains the following fields:

```proto
message Product {
  string product_id = 1; // UUID
  string name = 2;
  int32 quantity = 3;
  string type = 4;
  string category = 5;
  repeated string image_urls = 6;
  double price = 7;
  message Size {
    double width = 1;
    double height = 2;
  }
  Size size = 8;
  double weight = 9;
  double shipping_base_price = 10;
  int32 base_delivery_timelines = 11; // in days
  string seller_id = 12; // Seller information (ID only for simplicity)
}
```

- **product_id**: The unique identifier for the product (UUID).
- **name**: The name of the product.
- **quantity**: The quantity of the product available.
- **type**: The type of the product (e.g., "electronics", "clothing").
- **category**: The category the product belongs to (e.g., "smartphones", "furniture").
- **image_urls**: A list of URLs for images associated with the product.
- **price**: The price of the product.
- **size**: The size of the product, including width and height.
- **weight**: The weight of the product.
- **shipping_base_price**: The base shipping price.
- **base_delivery_timelines**: The estimated delivery time (in days).
- **seller_id**: The identifier of the seller providing the product.

## Running the Service Locally

### Prerequisites

Before running the Product Service locally, ensure the following:

- Go 1.18 or higher
- Protocol Buffers (Protobuf) Compiler (`protoc`)
- gRPC Go Plugin for Protobuf (`protoc-gen-go` and `protoc-gen-go-grpc`)

### Steps to Run Locally

1. Clone the repository:
   ```bash
   git clone https://github.com/tittuvarghese/ss-go-product-service.git
   cd product-service
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Generate gRPC code from the `proto` file:
   ```bash
   protoc --go_out=. --go-grpc_out=. proto/product.proto
   ```

4. Start the Product Service:
   ```bash
   go run cmd/main.go
   ```

The service will start and listen for gRPC requests on the specified port (e.g., `50051`).


## Example Usage with Gateway Service

### Create a New Product
```bash
curl -X POST http://localhost:8080/product/create \
   -d '{"product": {"name": "Smartphone", "quantity": 10, "price": 299.99, "category": "electronics", "type": "mobile", "seller_id": "12345"}}' \
   -H "Content-Type: application/json"
```

### Get a Single Product by ID
```bash
curl -X GET http://localhost:8080/product/12345 \
   -H "Authorization: Bearer <your_jwt_token>"
```

### Update a Product
```bash
curl -X POST http://localhost:8080/product/update \
   -d '{"product_id": "12345", "product": {"name": "Smartphone Pro", "quantity": 15, "price": 349.99}}' \
   -H "Content-Type: application/json" \
   -H "Authorization: Bearer <your_jwt_token>"
```

### Get Multiple Products by IDs
```bash
curl -X GET http://localhost:8080/products \
   -d '{"query": ["12345", "67890"]}' \
   -H "Authorization: Bearer <your_jwt_token>"
```

## Architecture

The **Product Service** is a core component of the e-commerce microservices ecosystem, and it facilitates the management of products. Here’s how it fits into the overall architecture:

- **gRPC Communication**: The service communicates with other services through gRPC to ensure fast and efficient communication.
- **Database**: Stores product data such as product details, price, quantity, and seller information.
- **JWT Authentication**: Secures endpoints by using JWT for authenticated access.
