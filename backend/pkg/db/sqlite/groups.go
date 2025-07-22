package sqlite

import (
	"database/sql"

	"social-network/pkg/models"
)

type GroupRepository struct {
	DB *sql.DB
}

func (repo *GroupRepository) GetAllAndRelations(userID string) ([]models.Group, error) {
	var groups []models.Group
	rows, err := repo.DB.Query("SELECT group_id, name, (SELECT COUNT(*) FROM group_users WHERE group_users.group_id = groups.group_id AND group_users.user_id = '" + userID + "') as member, administrator = '" + userID + "' as admin FROM groups;")
	if err != nil {
		return groups, err
	}
	for rows.Next() {
		var group models.Group
		var member int
		var admin int
		rows.Scan(&group.ID, &group.Name, &member, &admin)
		if member != 0 {
			group.Member = true
		}
		if admin != 0 {
			group.Administrator = true
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (repo *GroupRepository) GetUserGroups(userID string) ([]models.Group, error) {
	var groups []models.Group
	rows, err := repo.DB.Query("SELECT group_id, name, administrator = '" + userID + "' as admin FROM groups WHERE (SELECT COUNT(*) FROM group_users WHERE group_users.group_id = groups.group_id AND group_users.user_id = '" + userID + "') = 1 OR administrator = '" + userID + "';")
	if err != nil {
		return groups, err
	}
	for rows.Next() {
		var group models.Group
		var admin int
		rows.Scan(&group.ID, &group.Name, &admin)
		if admin != 0 {
			group.Administrator = true
		} else {
			group.Member = true
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (repo *GroupRepository) NewGroup(group models.Group) error {
	stmt, err := repo.DB.Prepare("INSERT INTO groups (group_id, name,description,administrator) values (?,?,?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(group.ID, group.Name, group.Description, group.AdminID); err != nil {
		return err
	}
	return nil
}

func (repo *GroupRepository) GetGroupData(groupId string) (models.Group, error) {
	row := repo.DB.QueryRow("SELECT name, description,administrator FROM groups WHERE group_id = ? ", groupId)
	var group models.Group
	if err := row.Scan(&group.Name, &group.Description, &group.AdminID); err != nil {
		return group, err
	}
	group.ID = groupId
	return group, nil
}

func (repo *GroupRepository) GetGroupMembers(groupId string) ([]models.User, error) {
	var members []models.User
	rows, err := repo.DB.Query("SELECT user_id, IFNULL(nickname, first_name || ' ' || last_name), image FROM users WHERE (user_id = (SELECT administrator FROM groups WHERE group_id =?)) OR ((SELECT COUNT() FROM group_users WHERE group_id = ? AND  user_id = users.user_id )=1) ", groupId, groupId)
	if err != nil {
		return members, err
	}
	for rows.Next() {
		var member models.User
		rows.Scan(&member.ID, &member.Nickname, &member.ImagePath)
		members = append(members, member)
	}
	return members, nil
}

func (repo *GroupRepository) IsGroupMember(groupId, userId string) (bool, error) {
	row := repo.DB.QueryRow("SELECT COUNT() FROM group_users WHERE group_id = ? AND  user_id = ?", groupId, userId)
	var member int
	if err := row.Scan(&member); err != nil {
		return false, err
	}
	if member == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func (repo *GroupRepository) IsGroupAdmin(groupId, userId string) (bool, error) {
	row := repo.DB.QueryRow("SELECT COUNT() FROM groups WHERE group_id = ? AND  administrator = ?", groupId, userId)
	var admin int
	if err := row.Scan(&admin); err != nil {
		return false, err
	}
	if admin == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func (repo *GroupRepository) GetGroupAdmin(groupId string) (string, error) {
	row := repo.DB.QueryRow("SELECT administrator FROM groups WHERE group_id = ? ", groupId)
	var admin string
	if err := row.Scan(&admin); err != nil {
		return admin, err
	}
	return admin, nil
}

func (repo *GroupRepository) SaveGroupMember(userId, groupId string) error {
	stmt, err := repo.DB.Prepare("INSERT INTO group_users (group_id, user_id) values (?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(groupId, userId); err != nil {
		return err
	}
	return nil
}
