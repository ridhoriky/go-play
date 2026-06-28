package product_image

import (
	"context"

	"ne-project/src/internal/models/entity"
)

type ProductImageRepositoryItf interface {
	GetByProductID(ctx context.Context, productID string) ([]entity.ProductImage, error)
	AddImage(ctx context.Context, image *entity.ProductImage) error
	DeleteImage(ctx context.Context, imageID string) error
	SetPrimary(ctx context.Context, productID string, imageID string) error
	UnsetPrimary(ctx context.Context, productID string) error
	CountImages(ctx context.Context, productID string) (int, error)
}
