package models

type Group struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AdminID     string `json:"adminId"`

	Invitations []string `json:"invitations"`

	Member        bool `json:"member"` // true if current user is a member
	Administrator bool `json:"admin"`  // true if current user is admin
	RequestPending bool `json:"requestPending"` // true if request to join is pending
}

type GroupRepository interface {
	GetAllAndRelations(userId string) ([]Group, error)
	GetUserGroups(userId string) ([]Group, error)
	NewGroup(Group) error                               //create new group
	GetGroupData(groupId string) (Group, error)         //get info- name and desc
	GetGroupMembers(groupId string) ([]User, error)     // get all group members and admin
	GetGroupAdmin(groupId string) (string, error)       //get admin id
	IsGroupMember(groupId, userId string) (bool, error) //checks if user is a member
	IsGroupAdmin(groupId, userId string) (bool, error)  //checks if user is admin

	SaveGroupMember(userId, groupId string)error 
}
