package product

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/repositories/order"
	"ne-project/src/internal/repositories/product"
	"ne-project/src/internal/repositories/product_image"
	"ne-project/src/internal/repositories/wishlist"
)

type ProductServiceItf interface {
	GetAllProducts(ctx context.Context, userID string, req *dto.GetProductsQuery) (*dto.ProductListResponse, error)
	GetProductByID(ctx context.Context, id string, userID string) (*dto.ProductDetailResponse, error)
	GetProductBySlug(ctx context.Context, slug string, userID string) (*dto.ProductDetailResponse, error)
	CreateProduct(ctx context.Context, storeID string, product *dto.CreateProductRequest) (*entity.Product, error)
	UpdateProduct(ctx context.Context, id string, storeID string, product *dto.UpdateProductRequest) (*entity.Product, error)
	DeleteProduct(ctx context.Context, id string, storeID string) error
	CreateMultipleProducts(ctx context.Context, storeID string, products []entity.Product) ([]entity.Product, error)
	GetSellerProductDetail(ctx context.Context, id string, storeID string) (*dto.SellerProductDetailResponse, error)
	AddProductImage(ctx context.Context, productID string, storeID string, req *dto.AddProductImageRequest) (*entity.ProductImage, error)
	DeleteProductImage(ctx context.Context, productID string, imageID string, storeID string) error
	SetPrimaryImage(ctx context.Context, productID string, imageID string, storeID string) error
}

type productService struct {
	productRepository      product.ProductRepositoryItf
	wishlistRepository     wishlist.WishlistRepositoryItf
	productImageRepository product_image.ProductImageRepositoryItf
	orderRepository        order.OrderRepositoryItf
}

func NewProductService(productRepository product.ProductRepositoryItf, wishlistRepository wishlist.WishlistRepositoryItf, productImageRepository product_image.ProductImageRepositoryItf, orderRepository order.OrderRepositoryItf) ProductServiceItf {
	return &productService{
		productRepository:      productRepository,
		wishlistRepository:     wishlistRepository,
		productImageRepository: productImageRepository,
		orderRepository:        orderRepository,
	}
}
