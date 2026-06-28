package product_image

import (
	"context"

	"ne-project/src/internal/models/entity"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type productImageRepository struct {
	db *sqlx.DB
}

func NewProductImageRepository(db *sqlx.DB) ProductImageRepositoryItf {
	return &productImageRepository{db: db}
}

func (r *productImageRepository) GetByProductID(ctx context.Context, productID string) ([]entity.ProductImage, error) {
	rows, err := r.db.QueryxContext(ctx, getProductImagesQuery, productID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("productID", productID).Msg("err find product images")
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var images []entity.ProductImage
	for rows.Next() {
		var img entity.ProductImage
		if err := rows.StructScan(&img); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

func (r *productImageRepository) AddImage(ctx context.Context, image *entity.ProductImage) error {
	err := r.db.QueryRowContext(ctx, insertProductImageQuery, image.ProductID, image.URL, image.AltText, image.SortOrder, image.IsPrimary).Scan(&image.ID)
	return err
}

func (r *productImageRepository) DeleteImage(ctx context.Context, imageID string) error {
	_, err := r.db.ExecContext(ctx, deleteProductImageQuery, imageID)
	return err
}

func (r *productImageRepository) SetPrimary(ctx context.Context, productID string, imageID string) error {
	_, err := r.db.ExecContext(ctx, setPrimaryImageQuery, imageID, productID)
	return err
}

func (r *productImageRepository) UnsetPrimary(ctx context.Context, productID string) error {
	_, err := r.db.ExecContext(ctx, unsetPrimaryImageQuery, productID)
	return err
}

func (r *productImageRepository) CountImages(ctx context.Context, productID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, countProductImagesQuery, productID).Scan(&count)
	return count, err
}
