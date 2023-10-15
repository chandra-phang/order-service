package middleware

import (
	"net/http"
	"order-service/api/controllers"
	"order-service/httpconnector"

	"github.com/labstack/echo/v4"
)

// userContextKey is a key for saving a UserID into a context.
const UserContextKey = "userID"

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authServiceCon := httpconnector.GetAuthServiceConnector()
		userID, err := authServiceCon.Authenticate(ctx)
		if err != nil {
			return controllers.WriteError(ctx, http.StatusUnauthorized, err)
		}

		ctx.Set(UserContextKey, userID)

		return next(ctx)
	}
}
