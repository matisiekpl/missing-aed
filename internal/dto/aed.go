package dto

type AED struct {
	ID        string  `json:"id"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Source    string  `json:"source"`
}
