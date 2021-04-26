package posts

type Post struct {
	Id      int64  `json:"id,omitempty"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
