package structs

type Product struct {
	Id          int     `json:"id,omitempty"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageUrl    string  `json:"image_url"`
	Price       float64 `json:"price"`
	Weight      float64 `json:"weight,omitempty"`
}
