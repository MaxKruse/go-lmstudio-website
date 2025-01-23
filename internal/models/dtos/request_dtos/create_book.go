package requestdtos

// @Description		Create book request
type CreateBookRequest struct {
	Title         string  `json:"title"`
	Author        string  `json:"author"`
	Description   string  `json:"description"`
	ImageUrl      string  `json:"image_url"`
	PublishedDate string  `json:"published_date"`
	Isbn          string  `json:"isbn"`
	Price         float32 `json:"price"`
}
