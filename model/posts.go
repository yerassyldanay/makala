package model

type Post struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	Author    string  `json:"author"`
	Link      *string `json:"link"`
	Submakala string  `json:"submakala"`
	Content   *string `json:"content"`
	Score     *int64  `json:"score"`
	Promoted  bool    `json:"promoted"`
	Nsfw      bool    `json:"nsfw"`
}
