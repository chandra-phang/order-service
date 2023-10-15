package repositories

import (
	"database/sql"
	"order-service/model"

	"github.com/labstack/echo/v4"
)

type OrderItemRepository struct {
	db *sql.DB
}

func NewOrderItemRepositoryInstance(db *sql.DB) model.IOrderItemRepository {
	return &OrderItemRepository{
		db: db,
	}
}

func (r OrderItemRepository) CreateOrderItem(ctx echo.Context, orderItem model.OrderItem) error {
	sqlStatement := `
		INSERT INTO order_items
			(id, order_id, product_id, quantity, created_at, updated_at)
		VALUES
			(?, ?, ?, ?, ?, ?)
	`

	params := []interface{}{
		orderItem.ID,
		orderItem.OrderID,
		orderItem.ProductID,
		orderItem.Quantity,
		orderItem.CreatedAt,
		orderItem.UpdatedAt,
	}

	_, err := r.db.Exec(sqlStatement, params...)
	if err != nil {
		return err
	}

	return nil
}

func (r OrderItemRepository) GetOrderItems(ctx echo.Context, orderID string) ([]model.OrderItem, error) {
	sqlStatement := `
		SELECT
			id,
			order_id,
			product_id,
			quantity,
			created_at,
			updated_at
		FROM order_items
		WHERE order_id = ?
	`

	results, err := r.db.Query(sqlStatement, orderID)
	if err != nil {
		return nil, err
	}

	var orderItems = make([]model.OrderItem, 0)
	for results.Next() {
		var orderItem model.OrderItem
		err = results.Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.ProductID, &orderItem.Quantity, &orderItem.CreatedAt, &orderItem.UpdatedAt)
		if err != nil {
			return nil, err
		}

		orderItems = append(orderItems, orderItem)
	}
	return orderItems, nil
}
