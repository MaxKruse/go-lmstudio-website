package requestdtos

// @Description		Create book request
type CreateBookRequest struct {
	Title         string  `json:"title" db:"title"`
	Author        string  `json:"author" db:"author"`
	Description   string  `json:"description" db:"description"`
	ImageUrl      string  `json:"image_url" db:"image_url"`
	PublishedDate string  `json:"published_date" db:"published_date"`
	Isbn          string  `json:"isbn" db:"isbn"`
	Price         float32 `json:"price" db:"price"`
}
