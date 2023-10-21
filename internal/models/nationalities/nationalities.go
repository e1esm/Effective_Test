package nationalities

type NationalityResponse struct {
	Count         int           `json:"count"`
	Name          string        `json:"name"`
	Nationalities []Nationality `json:"country"`
}

type Nationality struct {
	ID          string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
