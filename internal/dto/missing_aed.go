package dto

type MissingAED struct {
	ID          string  `json:"id"`
	Latitude    float64 `json:"lat"`
	Longitude   float64 `json:"lon"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	City        string  `json:"city"`
	Description string  `json:"description"`
}
