package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (r *TransactionRepository) Checkout(
	ctx context.Context, req *dto.CreateTransactionRequest,
) (*entity.TransactionWithDetails, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	details, totalAmount, err := r.processItems(ctx, tx, req.Items)
	if err != nil {
		return nil, err
	}

	trx, err := r.insertTransaction(ctx, tx, totalAmount)
	if err != nil {
		return nil, err
	}

	if err = r.insertTransactionDetails(ctx, tx, trx.ID, details); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &entity.TransactionWithDetails{Transaction: trx, Items: details}, nil
}

// processItem
func (r *TransactionRepository) processItems(
	ctx context.Context, tx *sql.Tx, items []dto.CheckoutItem,
) ([]entity.TransactionDetail, float64, error) {
	var totalAmount float64
	var details []entity.TransactionDetail

	for _, item := range items {
		detail, err := r.processSingleItem(ctx, tx, item)
		if err != nil {
			return nil, 0, err
		}
		totalAmount += detail.Subtotal.InexactFloat64()
		details = append(details, *detail)
	}

	return details, totalAmount, nil
}

// processSingleItem
func (r *TransactionRepository) processSingleItem(
	ctx context.Context, tx *sql.Tx, item dto.CheckoutItem,
) (*entity.TransactionDetail, error) {
	product, err := r.fetchProduct(ctx, tx, item.ProductID)
	if err != nil {
		return nil, err
	}

	if err = r.lockAndValidateStock(ctx, tx, product.Name, item); err != nil {
		return nil, err
	}

	if _, err = tx.ExecContext(ctx, deductStockQuery, item.Quantity, item.ProductID); err != nil {
		return nil, err
	}

	return buildTransactionDetail(product, item), nil
}

// fetchProduct
func (r *TransactionRepository) fetchProduct(
	ctx context.Context, tx *sql.Tx, productID string,
) (*dto.ProductSnapshot, error) {
	p := &dto.ProductSnapshot{}
	err := tx.QueryRowContext(ctx, getProductForCheckoutQuery, productID).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, &dto.Error{
			Code:    http.StatusNotFound,
			Message: "Product not found: " + productID,
		}
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

// lockAndValidateStock
func (r *TransactionRepository) lockAndValidateStock(
	ctx context.Context, tx *sql.Tx, productName string, item dto.CheckoutItem,
) error {
	var lockedStock int
	if err := tx.QueryRowContext(ctx, lockProductStockQuery, item.ProductID).Scan(&lockedStock); err != nil {
		return err
	}

	if lockedStock < item.Quantity {
		return &dto.Error{
			Code: http.StatusUnprocessableEntity,
			Message: fmt.Sprintf(
				"Insufficient stock for product '%s' (available: %d)",
				productName, lockedStock,
			),
		}
	}
	return nil
}

// insertTransaction
func (r *TransactionRepository) insertTransaction(
	ctx context.Context, tx *sql.Tx, totalAmount float64,
) (*entity.Transaction, error) {
	trx := &entity.Transaction{}
	err := tx.QueryRowContext(ctx, insertTransactionQuery,
		uuid.New().String(), totalAmount, entity.TransactionStatusPaid,
	).Scan(&trx.ID, &trx.TotalAmount, &trx.Status, &trx.CreatedAt)
	if err != nil {
		return nil, err
	}
	return trx, nil
}

// insertTransactionDetails
func (r *TransactionRepository) insertTransactionDetails(
	ctx context.Context, tx *sql.Tx, transactionID string, details []entity.TransactionDetail,
) error {
	for i := range details {
		details[i].TransactionID = transactionID
		_, err := tx.ExecContext(ctx, insertTransactionDetailQuery,
			details[i].ID,
			details[i].TransactionID,
			details[i].ProductID,
			details[i].ProductName,
			details[i].Quantity,
			details[i].Price,
			details[i].Subtotal,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// buildTransactionDetail
func buildTransactionDetail(p *dto.ProductSnapshot, item dto.CheckoutItem) *entity.TransactionDetail {
	subtotal := p.Price * float64(item.Quantity)
	return &entity.TransactionDetail{
		ID:          uuid.New().String(),
		ProductID:   p.ID,
		ProductName: p.Name,
		Quantity:    item.Quantity,
		Price:       decimal.NewFromFloat(p.Price),
		Subtotal:    decimal.NewFromFloat(subtotal),
	}
}
