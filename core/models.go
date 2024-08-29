package core

type User struct {
	Id             int64 // id
	Password       string
	Email          string
	ReadingHistory []string // book ids
}

type Book struct {
	Isbn               string  `db:"isbn"`
	Isbn13             string  `db:"isbn13"`
	Title              string  `db:"title"`
	TitleWithoutSeries string  `db:"title_without_series"`
	Edition            string  `db:"edition_information"`
	Description        string  `db:"description"`
	Publisher          string  `db:"publisher"`
	PublicationDay     int32   `db:"publication_day"`
	PublicationMonth   int32   `db:"publication_month"`
	PublicationYear    int32   `db:"publication_year"`
	Pages              int32   `db:"num_pages"`
	LanguageCode       string  `db:"language_code"`
	CountryCode        string  `db:"country_code"`
	Link               string  `db:"link"`
	Url                string  `db:"url"`
	ImageUrl           string  `db:"image_url"`
	Id                 int64   `db:"id"`
	WorkId             string  `db:"work_id"`
	Reviews            int32   `db:"text_reviews_count"`
	Ratings            int32   `db:"ratings_count"`
	Rating             float32 `db:"average_rating"`
}