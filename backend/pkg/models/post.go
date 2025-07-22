package models

type Post struct {
	ID         string `json:"id"`
	Content    string `json:"content"`
	ImagePath  string `json:"image"`
	AuthorID   string `json:"authorId"`
	Visibility string `json:"visibility"`
	GroupID    string `json:"groupId"`
	// for sending back with author
	Author   User      `json:"author"`
	Comments []Comment `json:"comments"`
}

type PostRepository interface {
	// Get all posts that user have access to
	GetAll(userID string) ([]Post, error)
	// get user posts that current user have access to
	GetUserPosts(userID, currentUserID string) ([]Post, error)
	// get group psts from specific group
	GetGroupPosts(groupId string)([]Post, error)
	
	New(Post) error

	SaveAccess(postId, userId string) error //save access for almost_private post
}
