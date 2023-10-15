package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/api/controllers"
	"order-service/config"
	v1req "order-service/dto/request/v1"
	v1resp "order-service/dto/response/v1"
	"order-service/request"

	"github.com/labstack/echo/v4"
)

// contextKey is a custom type for putting auth-related values into a context.
type contextKey string

// userContextKey is a key for saving a User object into a context.
const UserContextKey contextKey = "userID"

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		dto := v1req.AuthenticateDTO{
			SourceUri: c.Request().Host + c.Request().RequestURI,
		}

		data, err := json.Marshal(dto)
		if err != nil {
			return controllers.WriteError(c, http.StatusInternalServerError, err)
		}

		authHeader := c.Request().Header.Get("Authorization")

		// pass the Authorization header to auth service API
		authUrl := fmt.Sprintf("%s/v1/authenticate", config.GetConfig().AuthSvcHost)
		resp, statusCode, err := request.PostWithAuthorization(authUrl, data, authHeader)
		if err != nil {
			return controllers.WriteError(c, http.StatusInternalServerError, err)
		}
		if statusCode != http.StatusOK {
			return controllers.WriteErrorMsg(c, statusCode, "Authentication failed")
		}

		var response v1resp.AuthenticateDTO
		if err := json.Unmarshal(resp, &response); err != nil {
			return controllers.WriteError(c, statusCode, err)
		}

		c.Set(string(UserContextKey), response.Result.User.ID)

		return next(c)
	}
}