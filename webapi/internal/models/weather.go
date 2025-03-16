package models

type WeatherInfo struct {
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
	Location    string  `json:"location"`
}

type WeatherResponse struct {
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
}
