package example

type Order struct {
	Id     string  `json:"id"`
	Price  float64 `json:"price"`
	UserId string  `json:"user_id"`
}
