package v1

type IncreaseProductBookedQuotaDTO struct {
	Products []ProductDTO `json:"products"`
}

type DecreaseProductBookedQuotaDTO struct {
	Products []ProductDTO `json:"products"`
}

type ProductDTO struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}
