package apperrors

import "errors"

var (
	ErrOrderNotFound                = errors.New("order not found")
	ErrOrderAlreadyCancelled        = errors.New("order already cancelled")
	ErrUserIdIsEmpty                = errors.New("userId is empty")
	ErrProductNotFound              = errors.New("product not found")
	ErrCartNotFound                 = errors.New("cart not found")
	ErrCartQuantityIsInvalid        = errors.New("cart quantity is invalid")
	ErrFailedToIncreaseProductQuota = errors.New("failed to increase product quota")
	ErrFailedToDecreaseProductQuota = errors.New("failed to decrease product quota")
	ErrUnauthorized                 = errors.New("insufficient permissions")
	ErrAuthenticationFailed         = errors.New("authentication failed")
)
