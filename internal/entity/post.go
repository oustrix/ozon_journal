package entity

type Post struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	PublishedAt int       `json:"published_at"`
	AuthorID    int       `json:"author_id"`
	Commentable bool      `json:"commentable"`
	Comments    []Comment `json:"comments"`
}
