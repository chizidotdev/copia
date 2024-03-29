// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package repository

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	// Create product images for a product:
	BulkCreateProductImages(ctx context.Context, arg BulkCreateProductImagesParams) ([]ProductImage, error)
	// Clear cart items
	ClearCartItems(ctx context.Context, userID uuid.UUID) error
	// Create a new customer
	CreateCustomer(ctx context.Context, arg CreateCustomerParams) (Customer, error)
	// Create a new order
	CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error)
	// Create a new order item
	CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (OrderItem, error)
	// Create a new product
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	// Create a product image for a product:
	CreateProductImage(ctx context.Context, arg CreateProductImageParams) (ProductImage, error)
	// Create a new store
	CreateStore(ctx context.Context, arg CreateStoreParams) (Store, error)
	// Create a new user
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	// Remove a cart item
	DeleteCartItem(ctx context.Context, arg DeleteCartItemParams) (CartItem, error)
	// Delete a customer by ID
	DeleteCustomer(ctx context.Context, arg DeleteCustomerParams) error
	// Delete an order by ID
	DeleteOrder(ctx context.Context, arg DeleteOrderParams) (Order, error)
	// Delete an order item by ID
	DeleteOrderItem(ctx context.Context, id uuid.UUID) (OrderItem, error)
	// Delete a product by ID
	DeleteProduct(ctx context.Context, arg DeleteProductParams) error
	// Delete a product image by ID:
	DeleteProductImage(ctx context.Context, id uuid.UUID) error
	// Delete a store by ID
	DeleteStore(ctx context.Context, id uuid.UUID) error
	// Delete a user by ID
	DeleteUser(ctx context.Context, id uuid.UUID) error
	// Get cart items by user id
	GetCartItems(ctx context.Context, userID uuid.UUID) ([]GetCartItemsRow, error)
	// Get a customer by ID
	GetCustomer(ctx context.Context, arg GetCustomerParams) (GetCustomerRow, error)
	// Get an order by ID
	GetOrder(ctx context.Context, id uuid.UUID) (Order, error)
	// Get an order item by ID
	GetOrderItem(ctx context.Context, id uuid.UUID) (OrderItem, error)
	// Get a product by ID
	GetProduct(ctx context.Context, id uuid.UUID) (Product, error)
	// Get a product image by ID:
	GetProductImage(ctx context.Context, id uuid.UUID) (ProductImage, error)
	// Get a store by ID
	GetStore(ctx context.Context, id uuid.UUID) (Store, error)
	// Get a store by user_id
	GetStoreByUserId(ctx context.Context, userID uuid.UUID) (Store, error)
	// Get a user by ID
	GetUser(ctx context.Context, id uuid.UUID) (GetUserRow, error)
	// Get a user by email
	GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error)
	// List all customers for a store
	ListCustomers(ctx context.Context, storeID uuid.UUID) ([]ListCustomersRow, error)
	// List all order items for an order
	ListOrderItems(ctx context.Context, orderID uuid.UUID) ([]ListOrderItemsRow, error)
	// List all product images for a product:
	ListProductImages(ctx context.Context, productID uuid.UUID) ([]ProductImage, error)
	// List all products by store ID
	ListProductsByStore(ctx context.Context, storeID uuid.UUID) ([]Product, error)
	// List all order items for a store
	ListStoreOrderItems(ctx context.Context, storeID uuid.UUID) ([]ListStoreOrderItemsRow, error)
	// List all stores
	ListStores(ctx context.Context, userID uuid.UUID) ([]Store, error)
	// List all orders for a user
	ListUserOrders(ctx context.Context, userID uuid.UUID) ([]Order, error)
	// List all users
	ListUsers(ctx context.Context) ([]User, error)
	// Search a stores product by title and description
	SearchProducts(ctx context.Context, query string) ([]Product, error)
	// Search a store by name
	SearchStores(ctx context.Context, query string) ([]SearchStoresRow, error)
	// Set a product image as primary and others as non-primary:
	SetPrimaryImage(ctx context.Context, arg SetPrimaryImageParams) error
	// Update cart item quantity
	UpdateCartItemQuantity(ctx context.Context, arg UpdateCartItemQuantityParams) (CartItem, error)
	// Update an order by ID
	UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error)
	// Update an order item status
	UpdateOrderItem(ctx context.Context, arg UpdateOrderItemParams) (OrderItem, error)
	// Update a product by ID
	UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error)
	// Update a product image by ID:
	UpdateProductImageURL(ctx context.Context, arg UpdateProductImageURLParams) error
	// Update a store by ID
	UpdateStore(ctx context.Context, arg UpdateStoreParams) ([]Store, error)
	// Update a user by ID
	UpdateUser(ctx context.Context, arg UpdateUserParams) error
	// Add or update a cart item
	UpsertCartItem(ctx context.Context, arg UpsertCartItemParams) (CartItem, error)
	// Upsert a user by email
	UpsertUser(ctx context.Context, arg UpsertUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
