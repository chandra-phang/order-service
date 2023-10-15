package repositories

import (
	"database/sql"
	"log"
	"order-service/apperrors"
	"order-service/model"
	"time"

	"github.com/labstack/echo/v4"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepositoryInstance(db *sql.DB) model.IOrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r OrderRepository) CreateOrder(ctx echo.Context, order model.Order) error {
	sqlStatement := `
		INSERT INTO orders
			(id, user_id, status, created_at, updated_at)
		VALUES
			(?, ?, ?, ?, ?)
	`

	params := []interface{}{
		order.ID,
		order.UserID,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
	}

	_, err := r.db.Exec(sqlStatement, params...)
	if err != nil {
		return err
	}

	return nil
}

func (r OrderRepository) CancelOrder(ctx echo.Context, orderID string) error {
	sqlStatement := `
		UPDATE orders
		SET
			status = ?, updated_at = ?
		WHERE id = ?
	`

	params := []interface{}{
		model.OrderStatusCancelled,
		time.Now(),
		orderID,
	}

	_, err := r.db.Exec(sqlStatement, params...)
	if err != nil {
		return err
	}

	return nil
}

func (r OrderRepository) GetOrder(ctx echo.Context, orderID string) (*model.Order, error) {
	sqlStatement := `
		SELECT
			id,
			user_id,
			status,
			created_at,
			updated_at
		FROM orders
		WHERE id = ?
	`

	results, err := r.db.Query(sqlStatement, orderID)
	if err != nil {
		return nil, err
	}

	var order model.Order
	for results.Next() {
		err = results.Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			log.Println("failed to scan", err)
			return nil, err
		}
	}

	if order.ID == "" {
		return nil, apperrors.ErrOrderNotFound
	}

	return &order, nil
}

func (r OrderRepository) GetOrders(ctx echo.Context, userID string) ([]model.Order, error) {
	sqlStatement := `
		SELECT
			id,
			user_id,
			status,
			created_at,
			updated_at
		FROM orders
		WHERE user_id = ?
		ORDER BY updated_at DESC
	`

	results, err := r.db.Query(sqlStatement, userID)
	if err != nil {
		return nil, err
	}

	var orders = make([]model.Order, 0)
	for results.Next() {
		var order model.Order
		err = results.Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			log.Println("failed to scan", err)
			return nil, err
		}

		orders = append(orders, order)
	}
	return orders, nil
}
