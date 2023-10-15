package services

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"order-service/api/middleware"
	"order-service/apperrors"
	"order-service/config"
	"order-service/db"
	v1request "order-service/dto/request/v1"
	"order-service/handlers"
	"order-service/model"
	"order-service/repositories"
	"order-service/request"

	"github.com/labstack/echo/v4"
)

type IOrderService interface {
	// svc CRUD methods for domain objects
	CreateOrder(ctx echo.Context, dto v1request.CreateOrderDTO) error
	CancelOrder(ctx echo.Context, orderID string) error
	ListOrder(ctx echo.Context) ([]model.Order, error)
}

type orderSvc struct {
	config        *config.Config
	dbCon         *sql.DB
	orderRepo     model.IOrderRepository
	orderItemRepo model.IOrderItemRepository
}

var orderSvcSingleton IOrderService

func InitOrderService(h handlers.Handler) {
	orderSvcSingleton = orderSvc{
		config:        config.GetConfig(),
		dbCon:         db.GetDB(),
		orderRepo:     repositories.NewOrderRepositoryInstance(h.DB),
		orderItemRepo: repositories.NewOrderItemRepositoryInstance(h.DB),
	}
}

func GetOrderService() IOrderService {
	return orderSvcSingleton
}

func (svc orderSvc) CreateOrder(ctx echo.Context, dto v1request.CreateOrderDTO) error {
	userID := ctx.Get(middleware.UserContextKey)
	if userID == nil {
		return apperrors.ErrUserIdIsEmpty
	}

	userIdString := userID.(string)

	// add DB transaction
	tx, _ := svc.dbCon.Begin()
	defer tx.Rollback()

	// create order
	order := new(model.Order).Initialize(userIdString)
	err := svc.orderRepo.CreateOrder(ctx, *order)
	if err != nil {
		return err
	}

	reqDTO := v1request.IncreaseProductBookedQuotaDTO{}

	// validate orderItems
	for _, orderItemDTO := range dto.OrderItems {
		// validate productID with product-service - get product
		url := svc.config.ProductSvcHost + svc.config.GetProductUri + orderItemDTO.ProductID
		_, statusCode, err := request.Get(url)
		if err != nil {
			return err
		}

		if statusCode != http.StatusOK {
			return apperrors.ErrProductNotFound
		}

		// create orderItem
		orderItem := new(model.OrderItem).Initialize(order.ID, orderItemDTO.ProductID, orderItemDTO.Quantity)
		err = svc.orderItemRepo.CreateOrderItem(ctx, *orderItem)
		if err != nil {
			return err
		}

		// prepare requestBody for increase-booked-quota API
		productDTO := v1request.ProductDTO{
			ProductID: orderItem.ProductID,
			Quantity:  orderItem.Quantity,
		}
		reqDTO.Products = append(reqDTO.Products, productDTO)
	}

	// send request to product-service to increase-booked-quota
	url := svc.config.ProductSvcHost + svc.config.IncreaseBookedQuotaUri
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
	userID := ctx.Get(middleware.UserContextKey)
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

	// validate userID
	if order.UserID != userID {
		return apperrors.ErrUnauthorized
	}

	// validate order.Status
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

	// prepare requestBody to decrease-booked-quota API
	reqDTO := v1request.DecreaseProductBookedQuotaDTO{}
	for _, orderItem := range orderItems {
		productDTO := v1request.ProductDTO{
			ProductID: orderItem.ProductID,
			Quantity:  orderItem.Quantity,
		}

		reqDTO.Products = append(reqDTO.Products, productDTO)
	}

	// send request to product-service to decrease-booked-quota
	url := svc.config.ProductSvcHost + svc.config.DecreaseBookedQuotaUri
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
	userID := ctx.Get(middleware.UserContextKey)
	if userID == nil {
		return nil, apperrors.ErrUserIdIsEmpty
	}

	// retrieve orders by userID
	orders, err := svc.orderRepo.GetOrders(ctx, userID.(string))
	if err != nil {
		return nil, err
	}

	return orders, nil
}
