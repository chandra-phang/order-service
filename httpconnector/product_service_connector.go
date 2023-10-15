package httpconnector

import (
	"encoding/json"
	"net/http"
	"order-service/apperrors"
	"order-service/config"
	v1request "order-service/dto/request/v1"
	"order-service/model"
	"order-service/request"

	"github.com/labstack/echo/v4"
)

var productServiceCon *ProductServiceConnector

type ProductServiceConnector struct {
	Host                   string
	GetProductUri          string
	IncreaseBookedQuotaUri string
	DecreaseBookedQuotaUri string
}

func InitProductServiceConnector(cfg config.Config) {
	productServiceCon = &ProductServiceConnector{
		Host:                   cfg.ProductSvcHost,
		GetProductUri:          cfg.GetProductUri,
		IncreaseBookedQuotaUri: cfg.IncreaseBookedQuotaUri,
		DecreaseBookedQuotaUri: cfg.DecreaseBookedQuotaUri,
	}
}

func GetProductServiceConnector() *ProductServiceConnector {
	return productServiceCon
}

func (con ProductServiceConnector) GetProduct(ctx echo.Context, productID string) error {
	url := con.Host + con.GetProductUri + productID
	_, statusCode, err := request.Get(url)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return apperrors.ErrProductNotFound
	}

	return nil
}

func (con ProductServiceConnector) IncreaseBookedQuota(ctx echo.Context, orderItems []model.OrderItem) error {
	reqDTO := v1request.IncreaseProductBookedQuotaDTO{}

	for _, orderItem := range orderItems {
		// prepare requestBody for increase-booked-quota API
		productDTO := v1request.ProductDTO{
			ProductID: orderItem.ProductID,
			Quantity:  orderItem.Quantity,
		}
		reqDTO.Products = append(reqDTO.Products, productDTO)
	}

	url := con.Host + con.IncreaseBookedQuotaUri
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

	return nil
}

func (con ProductServiceConnector) DecreaseBookedQuota(ctx echo.Context, orderItems []model.OrderItem) error {
	reqDTO := v1request.DecreaseProductBookedQuotaDTO{}
	for _, orderItem := range orderItems {
		productDTO := v1request.ProductDTO{
			ProductID: orderItem.ProductID,
			Quantity:  orderItem.Quantity,
		}

		reqDTO.Products = append(reqDTO.Products, productDTO)
	}

	url := con.Host + con.DecreaseBookedQuotaUri
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

	return nil
}
