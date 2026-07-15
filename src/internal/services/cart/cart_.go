package cart

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"ne-project/src/internal/config/appconfig"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/repositories/cart"
	"ne-project/src/internal/repositories/product"
	"ne-project/src/internal/repositories/store"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type cartService struct {
	cartRepo    cart.CartRepositoryItf
	productRepo product.ProductRepositoryItf
	storeRepo   store.StoreRepositoryItf
	cfg         *appconfig.Config
}

func NewCartService(
	cartRepo cart.CartRepositoryItf,
	productRepo product.ProductRepositoryItf,
	storeRepo store.StoreRepositoryItf,
	cfg *appconfig.Config,
) CartServiceItf {
	return &cartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
		storeRepo:   storeRepo,
		cfg:         cfg,
	}
}

func (s *cartService) AddToCart(ctx context.Context, buyerID string, req *dto.AddToCartRequest) (*dto.CartItemResponse, error) {
	// Validate Product
	prodDetail, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}
	if prodDetail == nil {
		return nil, dto.NewError(http.StatusNotFound, "product not found")
	}
	if !prodDetail.IsActive {
		return nil, dto.NewError(http.StatusBadRequest, "product is not active")
	}

	// Validate Store
	st, err := s.storeRepo.GetByID(ctx, prodDetail.StoreID)
	if err != nil {
		return nil, err
	}
	if st == nil || st.DeletedAt != nil {
		return nil, dto.NewError(http.StatusBadRequest, "store is not active")
	}

	if st.UserID == buyerID {
		return nil, dto.NewError(http.StatusForbidden, "cannot add your own product to cart")
	}

	// Check if product already in cart
	existingCart, err := s.cartRepo.GetByBuyerAndProduct(ctx, buyerID, req.ProductID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	var cartID string
	var finalQuantity int

	if existingCart != nil {
		cartID, finalQuantity, err = s.updateExistingCart(ctx, existingCart, &prodDetail.Product, req.Quantity)
	} else {
		cartID, finalQuantity, err = s.addNewCart(ctx, buyerID, &prodDetail.Product, req.Quantity)
	}
	if err != nil {
		return nil, err
	}

	subtotal := prodDetail.Price.Mul(decimal.NewFromInt(int64(finalQuantity)))

	return &dto.CartItemResponse{
		ID:       cartID,
		Product:  s.mapProductToResponse(prodDetail, prodDetail.CategoryName),
		Quantity: finalQuantity,
		Subtotal: subtotal,
		Store:    s.mapStoreToResponse(st),
	}, nil
}

func (s *cartService) updateExistingCart(ctx context.Context, existingCart *entity.Cart, prod *entity.Product, additionalQuantity int) (string, int, error) {
	finalQuantity := existingCart.Quantity + additionalQuantity
	if prod.Stock < finalQuantity {
		return "", 0, dto.NewError(http.StatusBadRequest, "insufficient product stock")
	}

	err := s.cartRepo.UpdateQuantity(ctx, existingCart.ID, finalQuantity)
	if err != nil {
		return "", 0, err
	}
	return existingCart.ID, finalQuantity, nil
}

func (s *cartService) addNewCart(ctx context.Context, buyerID string, prod *entity.Product, quantity int) (string, int, error) {
	if prod.Stock < quantity {
		return "", 0, dto.NewError(http.StatusBadRequest, "insufficient product stock")
	}

	currentCarts, err := s.cartRepo.GetByBuyerID(ctx, buyerID)
	if err != nil {
		return "", 0, err
	}
	if len(currentCarts) >= s.cfg.Business.MaxCartItems {
		return "", 0, dto.NewError(http.StatusBadRequest, "max cart items limit reached")
	}

	cartID := uuid.New().String()
	newCart := &entity.Cart{
		ID:        cartID,
		BuyerID:   buyerID,
		ProductID: prod.ID,
		Quantity:  quantity,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = s.cartRepo.Add(ctx, newCart)
	if err != nil {
		return "", 0, err
	}
	return cartID, quantity, nil
}

func (s *cartService) GetCart(ctx context.Context, buyerID string) (*dto.CartResponse, error) {
	carts, err := s.cartRepo.GetByBuyerID(ctx, buyerID)
	if err != nil {
		return nil, err
	}

	var items []dto.CartItemResponse
	totalAmount := decimal.Zero
	totalItems := 0

	for _, c := range carts {
		prodDetail, err := s.productRepo.GetByID(ctx, c.ProductID)
		if err != nil {
			continue // skip or handle error
		}
		if prodDetail == nil {
			continue
		}

		st, err := s.storeRepo.GetByID(ctx, prodDetail.StoreID)
		if err != nil || st == nil {
			continue
		}

		subtotal := prodDetail.Price.Mul(decimal.NewFromInt(int64(c.Quantity)))
		totalAmount = totalAmount.Add(subtotal)
		totalItems += c.Quantity

		items = append(items, dto.CartItemResponse{
			ID:       c.ID,
			Product:  s.mapProductToResponse(prodDetail, prodDetail.CategoryName),
			Quantity: c.Quantity,
			Subtotal: subtotal,
			Store:    s.mapStoreToResponse(st),
		})
	}

	if items == nil {
		items = []dto.CartItemResponse{}
	}

	return &dto.CartResponse{
		Items:       items,
		TotalAmount: totalAmount,
		TotalItems:  totalItems,
	}, nil
}

func (s *cartService) UpdateQuantity(ctx context.Context, buyerID string, cartID string, req *dto.UpdateCartRequest) error {
	c, err := s.cartRepo.GetByID(ctx, cartID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.NewError(http.StatusNotFound, "cart item not found")
		}
		return err
	}
	if c.BuyerID != buyerID {
		return dto.NewError(http.StatusNotFound, "cart item not found")
	}

	prodDetail, err := s.productRepo.GetByID(ctx, c.ProductID)
	if err != nil {
		return err
	}
	if prodDetail == nil {
		return dto.NewError(http.StatusNotFound, "product not found")
	}

	if prodDetail.Stock < req.Quantity {
		return dto.NewError(http.StatusBadRequest, "insufficient product stock")
	}

	return s.cartRepo.UpdateQuantity(ctx, cartID, req.Quantity)
}

func (s *cartService) RemoveFromCart(ctx context.Context, buyerID string, cartID string) error {
	c, err := s.cartRepo.GetByID(ctx, cartID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.NewError(http.StatusNotFound, "cart item not found")
		}
		return err
	}
	if c.BuyerID != buyerID {
		return dto.NewError(http.StatusNotFound, "cart item not found")
	}

	return s.cartRepo.Delete(ctx, cartID)
}

func (s *cartService) ClearCart(ctx context.Context, buyerID string) error {
	return s.cartRepo.DeleteByBuyerID(ctx, buyerID)
}

func (s *cartService) mapProductToResponse(prod *entity.ProductDetail, categoryName string) dto.ProductResponse {
	return dto.ProductResponse{
		ID:           prod.ID,
		StoreID:      prod.StoreID,
		Name:         prod.Name,
		Slug:         prod.Slug,
		Description:  prod.Description,
		Price:        prod.Price,
		Stock:        prod.Stock,
		CategoryID:   prod.CategoryID,
		CategoryName: categoryName,
		RatingAvg:    prod.RatingAvg,
		TotalSold:    prod.TotalSold,
		IsActive:     prod.IsActive,
		PrimaryImage: prod.PrimaryImage,
		CreatedAt:    prod.CreatedAt,
		UpdatedAt:    prod.UpdatedAt,
	}
}

func (s *cartService) mapStoreToResponse(st *entity.Store) dto.StoreResponse {
	return dto.StoreResponse{
		ID:          st.ID,
		UserID:      st.UserID,
		StoreName:   st.StoreName,
		Slug:        st.Slug,
		Description: st.Description,
		LogoURL:     st.LogoURL,
		BannerURL:   st.BannerURL,
		IsVerified:  st.IsVerified,
		RatingAvg:   st.RatingAvg,
		TotalSales:  st.TotalSales,
		CreatedAt:   st.CreatedAt,
		UpdatedAt:   st.UpdatedAt,
	}
}
