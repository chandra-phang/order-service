package repositories

import (
	"database/sql"
	"log"
	"order-service/apperrors"
	"order-service/model"
	"time"

	"github.com/labstack/echo/v4"
)

type CartProductRepository struct {
	db *sql.DB
}

func NewCartProductRepositoryInstance(db *sql.DB) model.ICartProductRepository {
	return &CartProductRepository{
		db: db,
	}
}

func (r CartProductRepository) CreateCartProduct(ctx echo.Context, cartProduct model.CartProduct) error {
	sqlStatement := `
		INSERT INTO cart_products
			(id, user_id, product_id, quantity, created_at, updated_at)
		VALUES
			(?, ?, ?, ?, ?, ?)
	`

	params := []interface{}{
		cartProduct.ID,
		cartProduct.UserID,
		cartProduct.ProductID,
		cartProduct.Quantity,
		cartProduct.CreatedAt,
		cartProduct.UpdatedAt,
	}

	_, err := r.db.Exec(sqlStatement, params...)
	if err != nil {
		return err
	}

	return nil
}

func (r CartProductRepository) GetCartProduct(ctx echo.Context, userID string, productID string) (*model.CartProduct, error) {
	sqlStatement := `
		SELECT
			id,
			user_id,
			product_id,
			quantity,
			created_at,
			updated_at
		FROM cart_products
		WHERE user_id = ? AND product_id = ?
	`

	params := []interface{}{userID, productID}
	results, err := r.db.Query(sqlStatement, params...)
	if err != nil {
		return nil, err
	}

	var cartProduct model.CartProduct
	for results.Next() {
		err = results.Scan(&cartProduct.ID, &cartProduct.UserID, &cartProduct.ProductID, &cartProduct.Quantity, &cartProduct.CreatedAt, &cartProduct.UpdatedAt)
		if err != nil {
			log.Println("failed to scan", err)
			return nil, err
		}
	}

	if cartProduct.ID == "" {
		return nil, apperrors.ErrCartNotFound
	}

	return &cartProduct, nil
}

func (r CartProductRepository) UpdateQuantity(ctx echo.Context, cartProductID string, quantity int) error {
	sqlStatement := `
		UPDATE cart_products
		SET
			quantity = quantity + ?, updated_at = ?
		WHERE id = ?
	`

	params := []interface{}{
		quantity,
		time.Now(),
		cartProductID,
	}

	_, err := r.db.Exec(sqlStatement, params...)
	if err != nil {
		return err
	}

	return nil
}
