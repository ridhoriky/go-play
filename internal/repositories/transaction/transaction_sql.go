package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ne-project/internal/database"
	"ne-project/internal/dto"
	"ne-project/internal/models"
	"time"
)
func (r *TransactionRepository) Checkout(
	ctx context.Context,
	req *dto.CheckoutRequest,
) (*models.Transaction, error) {

	var result *models.Transaction

	err := database.WithTransaction(ctx, r.db, func(tx *sql.Tx) error {

		totalAmount := 0
		details := make([]models.TransactionDetail, 0)

		for _, item := range req.Items {

			var price, stock int

			var productName string

			err := tx.QueryRowContext(
				ctx,
				lockProductQuery,
				item.ProductID,
			).Scan(&productName, &price, &stock)

			if err == sql.ErrNoRows {
				return fmt.Errorf("product id %d not found", item.ProductID)
			}

			if err != nil {
				return err
			}

			if stock < item.Quantity {
				return errors.New("product out of stock")
			}

			subtotal := price * item.Quantity
			totalAmount += subtotal

			_, err = tx.ExecContext(
				ctx,
				updateStockQuery,
				item.Quantity,
				item.ProductID,
			)

			if err != nil {
				return err
			}

			details = append(details, models.TransactionDetail{
				ProductID: item.ProductID,
				ProductName: productName,
				Quantity:  item.Quantity,
				Subtotal:  subtotal,
			})
		}
		var txID int
		var createdAt time.Time

		err := tx.QueryRowContext(
			ctx,
			insertTxQuery,
			totalAmount,
		).Scan(&txID)

		if err != nil {
			return err
		}

		rows := make([][]any, 0, len(details))

		for _, d := range details {
			rows = append(rows, []any{
				txID,
				d.ProductID,
				d.Quantity,
				d.Subtotal,
			})
		}

		query, args, err := database.BuildBulkInsert(
			"transaction_details",
			[]string{
				"transaction_id",
				"product_id",
				"quantity",
				"subtotal",
			},
			rows,
		)

		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}

		for i := range details {
			details[i].TransactionID = txID
		}

		result = &models.Transaction{
			ID:          txID,
			TotalAmount: totalAmount,
			CreatedAt:   createdAt,
			Details:     details,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TransactionRepository) GetToday(
	ctx context.Context,
) (*dto.TodaySummaryResponse, error) {

	var summary dto.TodaySummaryResponse

	err := r.db.QueryRowContext(ctx, querySummary).
		Scan(&summary.TotalRevenue, &summary.TotalTransactions)

	if err != nil {
		return nil, err
	}

	var best dto.BestSellingProduct

	err = r.db.QueryRowContext(ctx, queryBestProduct).
		Scan(&best.Name, &best.QuantitySold)

	if err == sql.ErrNoRows {

		summary.BestSelling = dto.BestSellingProduct{
			Name:         "-",
			QuantitySold: 0,
		}

		return &summary, nil
	}

	if err != nil {
		return nil, err
	}

	summary.BestSelling = best

	return &summary, nil
}
