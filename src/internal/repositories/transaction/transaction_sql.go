package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/preference"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

func (r *TransactionRepository) Checkout(
	ctx context.Context, req *dto.CreateTransactionRequest,
) (*entity.TransactionWithDetails, error) {

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err tx create users")
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
) ([]entity.TransactionDetail, decimal.Decimal, error) {
	totalAmount := decimal.NewFromInt(0)
	details := make([]entity.TransactionDetail, 0, len(items))

	for _, item := range items {
		detail, err := r.processSingleItem(ctx, tx, item)
		if err != nil {
			return nil, decimal.Zero, err
		}
		totalAmount = totalAmount.Add(detail.Subtotal)
		details = append(details, detail)
	}

	return details, totalAmount, nil
}

// processSingleItem
func (r *TransactionRepository) processSingleItem(
	ctx context.Context, tx *sql.Tx, item dto.CheckoutItem,
) (entity.TransactionDetail, error) {
	product, err := r.fetchProduct(ctx, tx, item.ProductID)
	if err != nil {
		return entity.TransactionDetail{}, err
	}

	if err = r.lockAndValidateStock(ctx, tx, product.Name, item); err != nil {
		return entity.TransactionDetail{}, err
	}

	if item.Quantity <= 0 {
		zerolog.Ctx(ctx).Error().Str("id", item.ProductID).Msg("err product out of stock")
		return entity.TransactionDetail{}, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Quantity for product with id '%s' must be greater than zero", item.ProductID),
		}
	}

	if _, err = tx.ExecContext(ctx, deductStockQuery, item.Quantity, item.ProductID); err != nil {
		zerolog.Ctx(ctx).Error().Str("id", item.ProductID).Msg("err update stock product")
		return entity.TransactionDetail{}, err
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
		zerolog.Ctx(ctx).Error().Err(err).Str("id", productID).Msg(preference.ErrProductNotFound)
		return nil, &dto.Error{
			Code:    http.StatusNotFound,
			Message: preference.ErrProductNotFound,
		}
	}
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", productID).Msg("err fetch product")
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
		zerolog.Ctx(ctx).Error().Err(err).Str("id", item.ProductID).Msg("err lock and validate product")
		return err
	}

	if lockedStock < item.Quantity {
		zerolog.Ctx(ctx).Error().Str("id", item.ProductID).Msg("err product out of stock")
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
	ctx context.Context, tx *sql.Tx, totalAmount decimal.Decimal,
) (entity.Transaction, error) {
	trx := entity.Transaction{}
	err := tx.QueryRowContext(ctx, insertTransactionQuery,
		uuid.New().String(), totalAmount, entity.TransactionStatusPaid,
	).Scan(&trx.ID, &trx.TotalAmount, &trx.Status, &trx.CreatedAt)
	if err != nil {
		zerolog.Ctx(ctx).Error().Msg("err insert transaction")
		return entity.Transaction{}, err
	}
	return trx, nil
}

// insertTransactionDetails
func (r *TransactionRepository) insertTransactionDetails(
	ctx context.Context, tx *sql.Tx, transactionID string, details []entity.TransactionDetail,
) error {

	const fieldPerRow = 7
	placeholders := make([]string, len(details))
	args := make([]any, 0, len(details)*fieldPerRow)

	for i, detail := range details {
		base := i * fieldPerRow
		placeholders[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d,$%d, $%d)",
			base+1, base+2, base+3, base+4, base+5, base+6, base+7,
		)

		detail.TransactionID = transactionID

		args = append(args,
			detail.ID,
			detail.TransactionID,
			detail.ProductID,
			detail.ProductName,
			detail.Quantity,
			detail.Price,
			detail.Subtotal,
		)
	}
	query := insertTransactionDetailQuery + strings.Join(placeholders, ",")
	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		zerolog.Ctx(ctx).Error().Msg("err insert transaction details")
		return err
	}

	return nil
}

// buildTransactionDetail
func buildTransactionDetail(p *dto.ProductSnapshot, item dto.CheckoutItem) entity.TransactionDetail {
	subtotal := p.Price.Mul(decimal.NewFromInt(int64(item.Quantity)))
	return entity.TransactionDetail{
		ID:          uuid.New().String(),
		ProductID:   p.ID,
		ProductName: p.Name,
		Quantity:    item.Quantity,
		Price:       p.Price,
		Subtotal:    subtotal,
	}
}

func (repo *TransactionRepository) GetTransactionByID(ctx context.Context, id string) (*entity.TransactionWithDetails, error) {
	query := getTransactionByIDQuery

	rows, err := repo.db.QueryContext(ctx, query, id)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("transaction_id", id).Msg("failed to query transaction with details")
		return nil, err
	}
	defer rows.Close()

	// Initialize result as pointer
	result := &entity.TransactionWithDetails{
		Items: make([]entity.TransactionDetail, 0),
	}

	// Track if we've scanned transaction header
	txScanned := false

	for rows.Next() {
		var (
			txID          string
			txTotalAmount decimal.Decimal
			txStatus      string
			txCreatedAt   time.Time
			txDetailID    sql.NullString
			transactionID sql.NullString
			productID     sql.NullString
			productName   sql.NullString
			quantity      sql.NullInt64
			price         sql.NullString
			subtotal      sql.NullString
		)

		err := rows.Scan(
			&txID,
			&txTotalAmount,
			&txStatus,
			&txCreatedAt,
			&txDetailID,
			&transactionID,
			&productID,
			&productName,
			&quantity,
			&price,
			&subtotal,
		)

		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("failed to scan transaction with details row")
			return nil, err
		}

		// Scan transaction header only once
		if !txScanned {
			result.Transaction = entity.Transaction{
				ID:          txID,
				TotalAmount: txTotalAmount,
				Status:      entity.TransactionStatus(txStatus),
				CreatedAt:   txCreatedAt,
			}
			txScanned = true
		}

		// Add transaction detail if it exists (LEFT JOIN may return NULL for details)
		if txDetailID.Valid && productID.Valid {
			priceDecimal, _ := decimal.NewFromString(price.String)
			subtotalDecimal, _ := decimal.NewFromString(subtotal.String)

			detail := entity.TransactionDetail{
				ID:            txDetailID.String,
				TransactionID: transactionID.String,
				ProductID:     productID.String,
				ProductName:   productName.String,
				Quantity:      int(quantity.Int64),
				Price:         priceDecimal,
				Subtotal:      subtotalDecimal,
			}
			result.Items = append(result.Items, detail)
		}
	}

	if err := rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error iterating transaction rows")
		return nil, err
	}

	// Check if transaction was found
	if !txScanned {
		zerolog.Ctx(ctx).Warn().Str("transaction_id", id).Msg("transaction not found")
		return nil, sql.ErrNoRows
	}

	return result, nil
}
