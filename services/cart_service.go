package services

import (
	"database/sql"
	"fmt"
	"net/http"
	"order-service/api/middleware"
	"order-service/apperrors"
	"order-service/config"
	"order-service/db"
	v1request "order-service/dto/request/v1"
	"order-service/handlers"
	"order-service/lib"
	"order-service/model"
	"order-service/repositories"
	"order-service/request"
	"time"

	"github.com/labstack/echo/v4"
)

type ICartService interface {
	// svc CRUD methods for domain objects
	AddToCart(ctx echo.Context, dto v1request.AddToCartDTO) error
}

type cartSvc struct {
	dbCon           *sql.DB
	cartProductRepo model.ICartProductRepository
}

var cartSvcSingleton ICartService

func InitCartService(h handlers.Handler) {
	cartSvcSingleton = cartSvc{
		dbCon:           db.GetDB(),
		cartProductRepo: repositories.NewCartProductRepositoryInstance(h.DB),
	}
}

func GetCartService() ICartService {
	return cartSvcSingleton
}

func (svc cartSvc) AddToCart(ctx echo.Context, dto v1request.AddToCartDTO) error {
	userID := ctx.Get(middleware.UserContextKey)
	if userID == "" {
		return apperrors.ErrUserIdIsEmpty
	}

	url := config.GetConfig().ProductSvcHost + fmt.Sprintf("/v1/products/%s", dto.ProductID)
	_, statusCode, err := request.Get(url)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return apperrors.ErrProductNotFound
	}

	// add DB transaction
	tx, _ := svc.dbCon.Begin()
	defer tx.Rollback()

	cartProduct, err := svc.cartProductRepo.GetCartProduct(ctx, userID.(string), dto.ProductID)
	if err != nil && err != apperrors.ErrCartNotFound {
		return err
	}

	if cartProduct == nil {
		cartProduct = &model.CartProduct{
			ID:        lib.GenerateUUID(),
			UserID:    userID.(string),
			ProductID: dto.ProductID,
			Quantity:  0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = svc.cartProductRepo.CreateCartProduct(ctx, *cartProduct)
		if err != nil {
			return err
		}
	}

	if cartProduct.Quantity+dto.Quantity < 0 {
		return apperrors.ErrCartQuantityIsInvalid
	}

	err = svc.cartProductRepo.UpdateQuantity(ctx, cartProduct.ID, dto.Quantity)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
