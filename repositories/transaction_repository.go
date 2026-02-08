package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"strings"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	var res *models.Transaction

	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	productRepository := NewProductRepository(repo.db)

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)
	for _, item := range items {
		product, err := productRepository.GetByID(item.ProductID)
		if err != nil {
			return nil, err
		}

		subtotal := item.Quantity * product.Price
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, product.ID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING ID", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for i, detail := range details {
		details[i].TransactionID = transactionID

		var detailID int
		err := tx.QueryRow("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING ID", transactionID, detail.ProductID, detail.Quantity, detail.Subtotal).Scan(&detailID)
		if err != nil {
			return nil, err
		}
		details[i].ID = detailID
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	res = &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}

	return res, nil
}

func (repo *TransactionRepository) GetTotalTransactions(from *time.Time, to *time.Time) (int, error) {
	query := "SELECT COUNT(*) FROM transactions"
	args := []interface{}{}
	argPos := 1
	conditions := []string{}

	if from != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argPos))
		args = append(args, *from)
		argPos++
	}

	if to != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argPos))
		args = append(args, *to)
		argPos++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int
	err := repo.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

func (repo *TransactionRepository) GetTotalRevenue(from *time.Time, to *time.Time) (int, error) {
	query := "SELECT COALESCE(SUM(total_amount), 0) FROM transactions"
	args := []interface{}{}
	argPos := 1
	conditions := []string{}

	if from != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argPos))
		args = append(args, *from)
		argPos++
	}
	if to != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argPos))
		args = append(args, *to)
		argPos++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err := repo.db.QueryRow(query, args...).Scan(&total)
	return total, err
}

func (repo *TransactionRepository) GetBestSellingProduct(from *time.Time, to *time.Time) (*models.SoldProduct, error) {
	query := `
SELECT p.name, SUM(td.quantity) AS qty_sold
FROM transaction_details td
JOIN products p ON td.product_id = p.id
JOIN transactions t ON td.transaction_id = t.id
`
	args := []interface{}{}
	argPos := 1
	conditions := []string{}

	if from != nil {
		conditions = append(conditions, fmt.Sprintf("t.created_at >= $%d", argPos))
		args = append(args, *from)
		argPos++
	}
	if to != nil {
		conditions = append(conditions, fmt.Sprintf("t.created_at <= $%d", argPos))
		args = append(args, *to)
		argPos++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
GROUP BY p.id
ORDER BY qty_sold DESC
LIMIT 1
`

	var name string
	var qty int
	err := repo.db.QueryRow(query, args...).Scan(&name, &qty)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &models.SoldProduct{
		Name:         name,
		QuantitySold: qty,
	}, nil
}
