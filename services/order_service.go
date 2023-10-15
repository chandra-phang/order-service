package services

import (
	"database/sql"
	"encoding/json"
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

type IOrderService interface {
	// svc CRUD methods for domain objects
	CreateOrder(ctx echo.Context, dto v1request.CreateOrderDTO) error
	CancelOrder(ctx echo.Context, orderID string) error
	ListOrder(ctx echo.Context) ([]model.Order, error)
}

type orderSvc struct {
	dbCon         *sql.DB
	orderRepo     model.IOrderRepository
	orderItemRepo model.IOrderItemRepository
}

var orderSvcSingleton IOrderService

func InitOrderService(h handlers.Handler) {
	orderSvcSingleton = orderSvc{
		dbCon:         db.GetDB(),
		orderRepo:     repositories.NewOrderRepositoryInstance(h.DB),
		orderItemRepo: repositories.NewOrderItemRepositoryInstance(h.DB),
	}
}

func GetOrderService() IOrderService {
	return orderSvcSingleton
}

func (svc orderSvc) CreateOrder(ctx echo.Context, dto v1request.CreateOrderDTO) error {
	userID := ctx.Get(string(middleware.UserContextKey))
	if userID == nil {
		return apperrors.ErrUserIdIsEmpty
	}

	// add DB transaction
	tx, _ := svc.dbCon.Begin()
	defer tx.Rollback()

	// create order
	order := model.Order{
		ID:        lib.GenerateUUID(),
		UserID:    userID.(string),
		Status:    model.OrderStatusCompleted,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := svc.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	reqDTO := v1request.IncreaseProductBookedQuotaDTO{}

	// validate orderItems
	for _, orderItemDTO := range dto.OrderItems {
		// send request to product-service to get product by ID
		url := config.GetConfig().ProductSvcHost + fmt.Sprintf("/v1/products/%s", orderItemDTO.ProductID)
		_, statusCode, err := request.Get(url)
		if err != nil {
			return err
		}

		if statusCode != http.StatusOK {
			return apperrors.ErrProductNotFound
		}

		orderItem := model.OrderItem{
			ID:        lib.GenerateUUID(),
			OrderID:   order.ID,
			ProductID: orderItemDTO.ProductID,
			Quantity:  orderItemDTO.Quantity,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = svc.orderItemRepo.CreateOrderItem(ctx, orderItem)
		if err != nil {
			return err
		}

		productDTO := v1request.ProductDTO{
			ProductID: orderItem.ProductID,
			Quantity:  orderItem.Quantity,
		}
		reqDTO.Products = append(reqDTO.Products, productDTO)
	}

	// send request to product-service to increase-booked-quota
	url := config.GetConfig().ProductSvcHost + "/v1/products/increase-booked-quota"
	reqBody, err := json.Marshal(reqDTO)
	if err != nil {
		return err
	}

	_, statusCode, err := request.Put(url, reqBody)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return apperrors.ErrFailedToIncreaseProductQuota
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (svc orderSvc) CancelOrder(ctx echo.Context, orderID string) error {
	userID := ctx.Get(string(middleware.UserContextKey))
	if userID == nil {
		return apperrors.ErrUserIdIsEmpty
	}

	// add DB transaction
	tx, _ := svc.dbCon.Begin()
	defer tx.Rollback()

	order, err := svc.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		return apperrors.ErrUnauthorized
	}

	if order.Status == model.OrderStatusCancelled {
		return apperrors.ErrOrderAlreadyCancelled
	}

	err = svc.orderRepo.CancelOrder(ctx, orderID)
	if err != nil {
		return err
	}

	orderItems, err := svc.orderItemRepo.GetOrderItems(ctx, order.ID)
	if err != nil {
		return err
	}

	reqDTO := v1request.DecreaseProductBookedQuotaDTO{}
	for _, orderItem := range orderItems {
		productDTO := v1request.ProductDTO{
			ProductID: orderItem.ProductID,
			Quantity:  orderItem.Quantity,
		}

		reqDTO.Products = append(reqDTO.Products, productDTO)
	}

	// send request to product-service to decrease-booked-quota
	url := config.GetConfig().ProductSvcHost + "/v1/products/decrease-booked-quota"
	reqBody, err := json.Marshal(reqDTO)
	if err != nil {
		return err
	}

	_, statusCode, err := request.Put(url, reqBody)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return apperrors.ErrFailedToDecreaseProductQuota
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (svc orderSvc) ListOrder(ctx echo.Context) ([]model.Order, error) {
	userID := ctx.Get(string(middleware.UserContextKey))
	if userID == nil {
		return nil, apperrors.ErrUserIdIsEmpty
	}

	orders, err := svc.orderRepo.GetOrders(ctx, userID.(string))
	if err != nil {
		return nil, err
	}

	return orders, nil
}
