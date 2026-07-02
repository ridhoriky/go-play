package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepositoryItf {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, orders []entity.Order, items []entity.OrderItem) error {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err begin tx for order create")
		return err
	}
	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			zerolog.Ctx(ctx).Error().Err(rbErr).Msg("err rollback order create")
		}
	}()

	err = r.validateAndDeductStock(ctx, tx, items)
	if err != nil {
		return err
	}

	err = r.insertOrders(ctx, tx, orders)
	if err != nil {
		return err
	}

	if len(items) > 0 {
		query, args := buildBulkInsertOrderItemsQuery(items)
		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *orderRepository) validateAndDeductStock(ctx context.Context, tx *sqlx.Tx, items []entity.OrderItem) error {
	var productIDs []string
	reqQty := make(map[string]int)

	for i := range items {
		pid := items[i].ProductID
		productIDs = append(productIDs, pid)
		reqQty[pid] += items[i].Quantity
	}
	sort.Strings(productIDs)
	productIDs = removeDuplicateStrings(productIDs)

	for _, pid := range productIDs {
		var stock int
		var isActive bool
		err := tx.QueryRowContext(ctx, getProductStockForUpdateQuery, pid).Scan(&stock, &isActive)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return dto.NewError(http.StatusBadRequest, fmt.Sprintf("product %s not found", pid))
			}
			return err
		}

		qty := reqQty[pid]
		if !isActive {
			return dto.NewError(http.StatusBadRequest, "a product is not active")
		}
		if stock < qty {
			return dto.NewError(http.StatusBadRequest, "insufficient stock for a product")
		}
		_, err = tx.ExecContext(ctx, deductProductStockQuery, qty, pid)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *orderRepository) insertOrders(ctx context.Context, tx *sqlx.Tx, orders []entity.Order) error {
	for i := range orders {
		order := &orders[i]
		_, err := tx.ExecContext(ctx, insertOrderQuery,
			order.ID, order.BuyerID, order.StoreID, order.OrderNumber, order.TotalAmount, order.Status, order.ShippingAddress,
			order.ShippingCost, order.PaymentMethod, order.PaymentRef, order.Notes, order.CreatedAt, order.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func removeDuplicateStrings(s []string) []string {
	if len(s) == 0 {
		return s
	}
	res := []string{s[0]}
	for i := 1; i < len(s); i++ {
		if s[i] != s[i-1] {
			res = append(res, s[i])
		}
	}
	return res
}

func buildBulkInsertOrderItemsQuery(items []entity.OrderItem) (string, []any) {
	insertColumns := []string{"id", "order_id", "product_id", "product_name", "product_image", "quantity", "price", "subtotal", "created_at"}
	values := make([]string, 0, len(items))
	args := make([]any, 0, len(items)*9)
	argPos := 1

	for i := range items {
		item := &items[i]
		values = append(values, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)", argPos, argPos+1, argPos+2, argPos+3, argPos+4, argPos+5, argPos+6, argPos+7, argPos+8))
		args = append(args, item.ID, item.OrderID, item.ProductID, item.ProductName, item.ProductImage, item.Quantity, item.Price, item.Subtotal, item.CreatedAt)
		argPos += 9
	}
	return fmt.Sprintf(insertOrderItemBulkQuery, strings.Join(insertColumns, ","), strings.Join(values, ",")), args
}

func (r *orderRepository) GetByID(ctx context.Context, id string) (*entity.Order, []entity.OrderItem, error) {
	var o entity.Order
	err := r.db.QueryRowContext(ctx, getOrderByIDQuery, id).Scan(
		&o.ID, &o.BuyerID, &o.StoreID, &o.OrderNumber, &o.TotalAmount, &o.Status, &o.ShippingAddress,
		&o.ShippingCost, &o.PaymentMethod, &o.PaymentRef, &o.Notes, &o.CreatedAt, &o.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, dto.NewError(http.StatusNotFound, "order not found")
	}
	if err != nil {
		return nil, nil, err
	}

	items, err := r.getOrderItems(ctx, o.ID)
	if err != nil {
		return nil, nil, err
	}

	return &o, items, nil
}

func (r *orderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*entity.Order, []entity.OrderItem, error) {
	var o entity.Order
	err := r.db.QueryRowContext(ctx, getOrderByOrderNumberQuery, orderNumber).Scan(
		&o.ID, &o.BuyerID, &o.StoreID, &o.OrderNumber, &o.TotalAmount, &o.Status, &o.ShippingAddress,
		&o.ShippingCost, &o.PaymentMethod, &o.PaymentRef, &o.Notes, &o.CreatedAt, &o.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, dto.NewError(http.StatusNotFound, "order not found")
	}
	if err != nil {
		return nil, nil, err
	}

	items, err := r.getOrderItems(ctx, o.ID)
	if err != nil {
		return nil, nil, err
	}

	return &o, items, nil
}

func (r *orderRepository) GetByPaymentRef(ctx context.Context, paymentRef string) (*entity.Order, []entity.OrderItem, error) {
	var o entity.Order
	err := r.db.QueryRowContext(ctx, getOrderByPaymentRefQuery, paymentRef).Scan(
		&o.ID, &o.BuyerID, &o.StoreID, &o.OrderNumber, &o.TotalAmount, &o.Status, &o.ShippingAddress,
		&o.ShippingCost, &o.PaymentMethod, &o.PaymentRef, &o.Notes, &o.CreatedAt, &o.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, dto.NewError(http.StatusNotFound, "order not found for the given payment ref")
	}
	if err != nil {
		return nil, nil, err
	}

	items, err := r.getOrderItems(ctx, o.ID)
	if err != nil {
		return nil, nil, err
	}

	return &o, items, nil
}

func (r *orderRepository) getOrderItems(ctx context.Context, orderID string) ([]entity.OrderItem, error) {
	var items []entity.OrderItem
	rows, err := r.db.QueryContext(ctx, getOrderItemsByOrderIDQuery, orderID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var item entity.OrderItem
		if scanErr := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.ProductName, &item.ProductImage,
			&item.Quantity, &item.Price, &item.Subtotal, &item.CreatedAt,
		); scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *orderRepository) GetByBuyerID(ctx context.Context, buyerID string, params *dto.GetOrdersQuery) ([]entity.Order, int, error) {
	query := getOrdersByBuyerIDBaseQuery
	args := []any{buyerID}
	return r.listOrders(ctx, query, args, params)
}

func (r *orderRepository) GetByStoreID(ctx context.Context, storeID string, params *dto.GetOrdersQuery) ([]entity.Order, int, error) {
	query := getOrdersByStoreIDBaseQuery
	args := []any{storeID}
	return r.listOrders(ctx, query, args, params)
}

func (r *orderRepository) listOrders(ctx context.Context, baseQuery string, baseArgs []any, params *dto.GetOrdersQuery) ([]entity.Order, int, error) {
	query := baseQuery
	args := baseArgs
	argPos := len(args) + 1

	if params.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, params.Status)
		argPos++
	}

	sortBy := "created_at"
	allowedSortColumns := map[string]string{
		"created_at":   "created_at",
		"total_amount": "total_amount",
	}
	if val, ok := allowedSortColumns[params.SortBy]; ok {
		sortBy = val
	}

	sortDir := "DESC"
	if params.SortDir == "ASC" {
		sortDir = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortDir)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)

	offset := (params.Page - 1) * params.Limit
	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var orders []entity.Order
	var total int

	for rows.Next() {
		var o entity.Order
		var rowCount int
		if scanErr := rows.Scan(
			&o.ID, &o.BuyerID, &o.StoreID, &o.OrderNumber, &o.TotalAmount, &o.Status, &o.ShippingAddress,
			&o.ShippingCost, &o.PaymentMethod, &o.PaymentRef, &o.Notes, &o.CreatedAt, &o.UpdatedAt, &rowCount,
		); scanErr != nil {
			return nil, 0, scanErr
		}
		orders = append(orders, o)
		if total == 0 {
			total = rowCount
		}
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	_, err := r.db.ExecContext(ctx, updateOrderStatusQuery, status, id)
	return err
}

func (r *orderRepository) UpdatePayment(ctx context.Context, id string, paymentMethod string, paymentRef string) error {
	_, err := r.db.ExecContext(ctx, updateOrderPaymentQuery, paymentMethod, paymentRef, id)
	return err
}

func (r *orderRepository) CancelOrder(ctx context.Context, orderID string, items []entity.OrderItem) error {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		return err
	}
	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			zerolog.Ctx(ctx).Error().Err(rbErr).Msg("err rollback stock restore")
		}
	}()

	_, err = tx.ExecContext(ctx, updateOrderStatusQuery, "canceled", orderID)
	if err != nil {
		return err
	}

	for i := range items {
		item := &items[i]
		_, err = tx.ExecContext(ctx, restoreProductStockQuery, item.Quantity, item.ProductID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *orderRepository) HasActiveOrdersForProduct(ctx context.Context, productID string) (bool, error) {
	var hasActive bool
	err := r.db.QueryRowContext(ctx, hasActiveOrdersForProductQuery, productID).Scan(&hasActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		zerolog.Ctx(ctx).Error().Err(err).Str("productID", productID).Msg("err check active orders for product")
		return false, err
	}
	return hasActive, nil
}

func (r *orderRepository) GetRecentOrdersByProductID(ctx context.Context, productID string, limit int) ([]dto.OrderSummary, error) {
	rows, err := r.db.QueryContext(ctx, getRecentOrdersByProductIDQuery, productID, limit)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("productID", productID).Msg("err find recent orders for product")
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var orders []dto.OrderSummary
	for rows.Next() {
		var o dto.OrderSummary
		if scanErr := rows.Scan(
			&o.ID, &o.BuyerName, &o.TotalAmount, &o.Status, &o.CreatedAt,
		); scanErr != nil {
			return nil, scanErr
		}
		orders = append(orders, o)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
