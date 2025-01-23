package entities

type Book struct {
	Id            int32   `db:"id"`
	Title         string  `db:"title"`
	Author        string  `db:"author"`
	Description   string  `db:"description"`
	ImageUrl      string  `db:"image_url"`
	PublishedDate string  `db:"published_date"`
	Isbn          string  `db:"isbn"`
	Price         float32 `db:"price"`
	CreatedAt     string  `db:"created_at"`
	UpdatedAt     string  `db:"updated_at"`
	DeletedAt     string  `db:"deleted_at"`
}
