package httpconnector

import (
	"encoding/json"
	"net/http"
	"order-service/apperrors"
	"order-service/config"
	v1request "order-service/dto/request/v1"
	v1response "order-service/dto/response/v1"
	"order-service/request"

	"github.com/labstack/echo/v4"
)

var authServiceCon *AuthServiceConnector

type AuthServiceConnector struct {
	Host            string
	AuthenticateUri string
}

func InitAuthServiceConnector(cfg config.Config) {
	authServiceCon = &AuthServiceConnector{
		Host:            cfg.AuthSvcHost,
		AuthenticateUri: cfg.AuthenticateUri,
	}
}

func GetAuthServiceConnector() *AuthServiceConnector {
	return authServiceCon
}

func (con AuthServiceConnector) Authenticate(ctx echo.Context) (string, error) {
	dto := v1request.AuthenticateDTO{
		SourceUri: ctx.Request().Host + ctx.Request().RequestURI,
	}

	data, err := json.Marshal(dto)
	if err != nil {
		return "", err
	}

	// pass the Authorization header to auth service API
	authorization := ctx.Request().Header.Get("Authorization")

	authUrl := con.Host + con.AuthenticateUri
	resp, statusCode, err := request.Post(authUrl, data, authorization)
	if err != nil {
		return "", err
	}
	if statusCode != http.StatusOK {
		return "", apperrors.ErrAuthenticationFailed
	}

	var response v1response.AuthenticateDTO
	if err := json.Unmarshal(resp, &response); err != nil {
		return "", err
	}

	return response.Result.User.ID, nil
}
