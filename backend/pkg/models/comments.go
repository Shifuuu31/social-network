package models

type Comment struct {
	ID string `json:"id"`

	PostID    string `json:"postId"`
	Content   string `json:"content"`
	ImagePath string `json:"image"`
	AuthorID  string `json:"authorId"`
	// for sending back with author
	Author User `json:"author"`
}

type CommentRepository interface {
	// get comment based on postID
	Get(postID string) ([]Comment, error)
	New(Comment) error
}
