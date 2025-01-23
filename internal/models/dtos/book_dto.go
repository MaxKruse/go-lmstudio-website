package dtos

// @Description		Book dto
type Book struct {
	Id            int32   `json:"id" example:"1"`
	Title         string  `json:"title" example:"The Great Gatsby"`
	Author        string  `json:"author" example:"F. Scott Fitzgerald"`
	Description   string  `json:"description" example:"A novel about the decadence of the Roaring Twenties."`
	ImageUrl      string  `json:"image_url" example:"https://example.com/image.jpg"`
	PublishedDate string  `json:"published_date" example:"1925-12-01"`
	Isbn          string  `json:"isbn" example:"978-1-84953-745-2"`
	Price         float32 `json:"price" example:"29.99"`
}
