package models

// defines  User data type
type User struct {
	ID          string `json:"id"`
	Email       string `json:"login"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Password    string `json:"password,omitempty"`
	Nickname    string `json:"nickname"`
	About       string `json:"about"`
	DateOfBirth string `json:"dateOfBirth"`
	ImagePath   string `json:"avatar"`
	Status      string `json:"status"`      // private / public
	CurrentUser bool   `json:"currentUser"` //returns true for current, false otherwise

	Follower  bool `json:"follower"`  //if this user is following another user
	Following bool `json:"following"` //if curr user is following this one
	FollowRequestPending bool `json:"requestPending"` // true if requested to follow
}

// Repository represent all possible actions availible to deal with User
// all db packages(in case of different db) should implement those function
type UserRepository interface {
	Add(User) error                           //save new user in db
	EmailNotTaken(email string) (bool, error) //returns true if not taken
	FindUserByEmail(email string) (User, error)

	GetAllAndFollowing(userID string) ([]User, error) //all users and follow info
	GetFollowers(userId string) ([]User, error)       //get client followers
	GetFollowing(userId string) ([]User, error)       //get who is following client
	SaveFollower(userId, followerId string) error      //save new follower
	DeleteFollower(userId, followerId string) error

	IsFollowing(userID, currentUserID string) (bool, error) //returns true if current is following
	GetDataMin(userID string) (User, error)                 // returns id, nickname and image (for comment or post author)
	ProfileStatus(userID string) (string, error)            //evaluates if profile public

	GetProfileMax(userID string) (User, error) //returns all data about user
	GetProfileMin(userID string) (User, error) //returns some data about user

	GetStatus(userID string) (string, error) //get current status
	SetStatus(User) error                    // change status (needs id and new status)
}
