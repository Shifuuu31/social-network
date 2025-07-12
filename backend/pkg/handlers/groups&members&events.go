package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"social-network/pkg/models"
	"social-network/pkg/tools"

	"github.com/google/uuid"
)

func (rt *Root) NewGroupsHandler() (groupsMux *http.ServeMux) {
	groupsMux = http.NewServeMux()

	groupsMux.HandleFunc("POST /group/new", rt.NewGroup)
	groupsMux.HandleFunc("GET /group/{id}", rt.GetGroup)
	groupsMux.HandleFunc("POST /group/events", rt.GetGroupEvents)
	groupsMux.HandleFunc("POST /group/invite", rt.InviteToJoinGroup)
	groupsMux.HandleFunc("POST /group/request", rt.RequestToJoinGroup)
	groupsMux.HandleFunc("POST /group/accept-decline", rt.AcceptDeclineGroup) // TODO: ws
	groupsMux.HandleFunc("POST /group/browse", rt.BrowseGroups)
	groupsMux.HandleFunc("POST /group/event", rt.NewEvent)
	groupsMux.HandleFunc("POST /group/event/vote", rt.EventVote)

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Group routes registered",
		Metadata: map[string]any{
			"path": "/group/new, /group/invite, /group/request, /group/accept-decline, /group/browse, /group/event, /group/event/vote",
		},
	})

	return groupsMux
}

// GetGroup retrieves a group by its ID
func (rt *Root) GetGroup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetGroup called")
	groupID, err := strconv.Atoi(r.PathValue("id"))
	if groupID <= 0 || err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Group ID not provided in request",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
			},
		})
		tools.RespondError(w, "Group ID is required", http.StatusBadRequest)
		return
	}
	group := &models.Group{
		ID: groupID,
	}

	err = rt.DL.Groups.GetGroupByID(group)
	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to retrieve group from DB",
			Metadata: map[string]any{
				"group_id": groupID,
				"ip":       r.RemoteAddr,
				"path":     r.URL.Path,
				"err":      err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the current user is a member of this group
	// requesterID := rt.DL.GetRequesterID(w, r) // TODO enable this when we have auth
	// if requesterID <= 0 {
	// 	rt.DL.Logger.Log(models.LogEntry{
	// 		Level:   "ERROR",
	// 		Message: "Unauthorized: requester ID not found",
	// 		Metadata: map[string]any{
	// 			"ip":   r.RemoteAddr,
	// 			"path": r.URL.Path,
	// 		},
	// 	})
	// 	tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
	// 	return

	// }

	member := &models.GroupMember{
		GroupID: groupID,
		UserID:  1, // TODO enable this when we have auth
	}
	err = rt.DL.Members.GetMember(member)
	if err != nil {
		// User is not a member - this is not an error
		group.IsMember = ""
	} else {
		group.IsMember = member.Status
	}

	// if group == nil {    // this will never be nil
	// 	rt.DL.Logger.Log(models.LogEntry{
	// 		Level:   "INFO",
	// 		Message: "Group not found",
	// 		Metadata: map[string]any{
	// 			"group_id": groupID,
	// 			"ip":       r.RemoteAddr,
	// 			"path":     r.URL.Path,
	// 		},
	// 	})
	// 	http.NotFound(w, r)
	// 	return
	// }

	if err := tools.EncodeJSON(w, http.StatusOK, group); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send group response",
			Metadata: map[string]any{
				"group_id": groupID,
				"ip":       r.RemoteAddr,
				"path":     r.URL.Path,
				"err":      err.Error(),
			},
		})
	}
}

func (rt *Root) NewGroup(w http.ResponseWriter, r *http.Request) {
	var group *models.Group
	if err := tools.DecodeJSON(r, &group); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode new group JSON",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "New group JSON decoded"})

	// Set the creator ID from the authenticated user (for now using hardcoded 1)
	group.CreatorID = 1 // TODO: Uncomment when auth is ready: rt.DL.GetRequesterID(w, r)

	// verify group creation input
	if err := group.Validate(); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Invalid group input",
			Metadata: map[string]any{
				"group": group,
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"err":   err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Group input validated"})

	// generate a unique uuid
	group.ImgUUID = uuid.NewString()

	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "img uuid generated successfully"})

	// insert group into db
	if err := rt.DL.Groups.Insert(group); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to insert group into DB",
			Metadata: map[string]any{
				"group": group,
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"err":   err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Group inserted into DB"})

	// Add creator as a member of the group
	creator := &models.GroupMember{
		GroupID:    group.ID,
		UserID:     group.CreatorID,
		Status:     "member",
		PrevStatus: "none", // Creator starts with no previous status
	}
	if err := rt.DL.Members.Upsert(creator); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to add creator as group member",
			Metadata: map[string]any{
				"group_id": group.ID,
				"user_id":  group.CreatorID,
				"ip":       r.RemoteAddr,
				"path":     r.URL.Path,
				"err":      err.Error(),
			},
		})
		// Note: We don't return here as the group was successfully created
		// This is just a membership addition that failed
	} else {
		rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Creator added as group member"})
	}

	// Initialize group chat room in WebSocket hub
	rt.Hub.InitializeGroupChat(group.ID)
	rt.DL.Logger.Log(models.LogEntry{
		Level:   "DEBUG",
		Message: "Group chat initialized in WebSocket hub",
		Metadata: map[string]any{
			"group_id": group.ID,
		},
	})

	// If the creator has an active WebSocket connection, automatically join them to the group chat
	if conn, hasConnection := rt.Hub.Clients[group.CreatorID]; hasConnection {
		rt.Hub.JoinGroup(group.CreatorID, group.ID)
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "DEBUG",
			Message: "Creator automatically joined group chat",
			Metadata: map[string]any{
				"group_id": group.ID,
				"user_id":  group.CreatorID,
			},
		})

		// Send a welcome message to the creator
		if conn != nil {
			welcomeMsg := map[string]interface{}{
				"type":      "group_chat_created",
				"group_id":  group.ID,
				"message":   "Group chat has been created! You can now start chatting with group members.",
				"timestamp": group.CreatedAt,
			}
			if err := conn.WriteJSON(welcomeMsg); err != nil {
				rt.DL.Logger.Log(models.LogEntry{
					Level:   "WARN",
					Message: "Failed to send welcome message to group creator",
					Metadata: map[string]any{
						"group_id": group.ID,
						"user_id":  group.CreatorID,
						"error":    err.Error(),
					},
				})
			}
		}
	}

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "New group created successfully",
		Metadata: map[string]any{
			"group": group,
			"ip":    r.RemoteAddr,
			"path":  r.URL.Path,
		},
	})

	if err := tools.EncodeJSON(w, http.StatusCreated, group); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send group response",
			Metadata: map[string]any{
				"group": group,
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"err":   err.Error(),
			},
		})
	}
}

func (rt *Root) InviteToJoinGroup(w http.ResponseWriter, r *http.Request) {
	var member *models.GroupMember
	if err := tools.DecodeJSON(r, &member); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode group invite info JSON",
			Metadata: map[string]any{
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"error": err.Error(),
			},
		})
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Group invite JSON decoded"})

	// TODO validate payload

	// Check if requester is creator
	requesterID := 1
	// requesterID := rt.DL.GetRequesterID(w, r)

	err1 := rt.DL.Members.IsUserGroupMember(member.GroupID, requesterID)
	err2 := rt.DL.Groups.IsUserCreator(member.GroupID, requesterID)

	if err1 != nil && err2 != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Forbidden: requester isn't the group creator",
			Metadata: map[string]any{
				"ip":     r.RemoteAddr,
				"path":   r.URL.Path,
				"error1": err1.Error(),
				"error2": err2.Error(),
			},
		})
		tools.RespondError(w, "Forbidden", http.StatusForbidden)
		return

	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Requester authorized for group invite"})

	// Insert into group_members as "pending_invite"
	// member.Status = "invited"
	if err := rt.DL.Members.Upsert(member); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to insert member into DB",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Member inserted with pending invite status"})

	// Create notification for the invited user
	// Get group information for the notification message
	groupInfo := &models.Group{ID: member.GroupID}
	if err := rt.DL.Groups.GetGroupByID(groupInfo); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Failed to get group info for notification, using generic message",
			Metadata: map[string]any{
				"group_id": member.GroupID,
				"error":    err.Error(),
			},
		})
		groupInfo.Title = "a group" // fallback
	}

	notification := &models.Notification{
		UserID:  member.UserID,
		Type:    "group_invite",
		Message: fmt.Sprintf("You've been invited to join the group '%s'.", groupInfo.Title),
		Seen:    false,
	}

	if err := rt.DL.Notifications.Insert(notification); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to create notification for group invite",
			Metadata: map[string]any{
				"user_id":  member.UserID,
				"group_id": member.GroupID,
				"error":    err.Error(),
			},
		})
		// Don't fail the whole request if notification creation fails
	} else {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Notification created for group invite",
			Metadata: map[string]any{
				"user_id":         member.UserID,
				"group_id":        member.GroupID,
				"notification_id": notification.ID,
			},
		})
	}

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Member invited successfully",
		Metadata: map[string]any{
			"ip":   r.RemoteAddr,
			"path": r.URL.Path,
		},
	})

	if err := tools.EncodeJSON(w, http.StatusCreated, nil); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send invite response",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
	}
}

func (rt *Root) RequestToJoinGroup(w http.ResponseWriter, r *http.Request) {
	var member *models.GroupMember
	if err := tools.DecodeJSON(r, &member); err != nil {
		fmt.Println(err)
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode group join request JSON",
			Metadata: map[string]any{
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"error": err.Error(),
			},
		})
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Group join request JSON decoded"})
	// TODO validate payload
	fmt.Println("Request to join group:", member)
	// Insert into group_members as "pending_request"
	// member.Status = "requested"
	if err := rt.DL.Members.Upsert(member); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to insert member into DB",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Member inserted with pending join request status"})

	// Create notification for the group creator about the join request
	// Get group information for the notification message
	groupInfo := &models.Group{ID: member.GroupID}
	if err := rt.DL.Groups.GetGroupByID(groupInfo); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Failed to get group info for notification, using generic message",
			Metadata: map[string]any{
				"group_id": member.GroupID,
				"error":    err.Error(),
			},
		})
		groupInfo.Title = "your group" // fallback
		groupInfo.CreatorID = 1        // fallback to avoid breaking
	}

	// Get the user who requested to join for a more personalized message
	requestUser := &models.User{ID: member.UserID}
	userDisplayName := "Someone"
	if err := rt.DL.Users.GetUserByID(requestUser); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Failed to get user info for notification, using generic message",
			Metadata: map[string]any{
				"user_id": member.UserID,
				"error":   err.Error(),
			},
		})
	} else {
		userDisplayName = requestUser.FirstName + " " + requestUser.LastName
	}

	notification := &models.Notification{
		UserID:  groupInfo.CreatorID,
		Type:    "group_join_request",
		Message: fmt.Sprintf("%s requested to join your group '%s'.", userDisplayName, groupInfo.Title),
		Seen:    false,
	}

	if err := rt.DL.Notifications.Insert(notification); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to create notification for group join request",
			Metadata: map[string]any{
				"creator_id":   groupInfo.CreatorID,
				"group_id":     member.GroupID,
				"requester_id": member.UserID,
				"error":        err.Error(),
			},
		})
		// Don't fail the whole request if notification creation fails
	} else {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Notification created for group join request",
			Metadata: map[string]any{
				"creator_id":      groupInfo.CreatorID,
				"group_id":        member.GroupID,
				"requester_id":    member.UserID,
				"notification_id": notification.ID,
			},
		})
	}

	tools.EncodeJSON(w, http.StatusCreated, member.Status)
}

func (rt *Root) AcceptDeclineGroup(w http.ResponseWriter, r *http.Request) {
	var member *models.GroupMember
	if err := tools.DecodeJSON(r, &member); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode group accept/decline JSON",
			Metadata: map[string]any{
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"error": err.Error(),
			},
		})
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Group accept/decline JSON decoded"})

	// TODO validate payload

	// Check if requester is creator
	requesterID := rt.DL.GetRequesterID(w, r)

	switch member.Status {
	case "requested":
		if err := rt.DL.Groups.IsUserCreator(member.GroupID, requesterID); err != nil {
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "ERROR",
				Message: "Forbidden: requester isn't the group creator",
				Metadata: map[string]any{
					"ip":    r.RemoteAddr,
					"path":  r.URL.Path,
					"error": err.Error(),
				},
			})
			tools.RespondError(w, "Forbidden", http.StatusForbidden)
			return
		}
	case "invited":
		if requesterID != member.ID {
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "ERROR",
				Message: "Forbidden: requester isn't invited to the group",
				Metadata: map[string]any{
					"ip":   r.RemoteAddr,
					"path": r.URL.Path,
				},
			})
			tools.RespondError(w, "Forbidden", http.StatusForbidden)
			return

		}
	default:
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Invalid status in payload",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
			},
		})
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return

	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Requester authorized for accept/decline"})

	// Update status to "joined" or delete row
	if err := rt.DL.Members.Upsert(member); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to insert member into DB",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Member status updated in DB"})

	// If the member status is now "member", automatically join them to the group chat
	if member.Status == "member" {
		// Ensure group chat exists in hub
		rt.Hub.InitializeGroupChat(member.GroupID)

		// If the new member has an active WebSocket connection, join them to the group chat
		if conn, hasConnection := rt.Hub.Clients[member.UserID]; hasConnection {
			rt.Hub.JoinGroup(member.UserID, member.GroupID)
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "DEBUG",
				Message: "New member automatically joined group chat",
				Metadata: map[string]any{
					"group_id": member.GroupID,
					"user_id":  member.UserID,
				},
			})

			// Send a welcome message to the new member
			if conn != nil {
				welcomeMsg := map[string]interface{}{
					"type":      "group_joined",
					"group_id":  member.GroupID,
					"message":   "Welcome to the group! You can now participate in group chat.",
					"timestamp": "",
				}
				if err := conn.WriteJSON(welcomeMsg); err != nil {
					rt.DL.Logger.Log(models.LogEntry{
						Level:   "WARN",
						Message: "Failed to send welcome message to new group member",
						Metadata: map[string]any{
							"group_id": member.GroupID,
							"user_id":  member.UserID,
							"error":    err.Error(),
						},
					})
				}
			}
		}
	}

	// TODO Notify group creator
	if err := tools.EncodeJSON(w, http.StatusOK, member); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send accept/decline response",
			Metadata: map[string]any{
				"member": member,
				"ip":     r.RemoteAddr,
				"path":   r.URL.Path,
				"error":  err.Error(),
			},
		})
	}
}

func (rt *Root) BrowseGroups(w http.ResponseWriter, r *http.Request) {
	var payload *models.GroupsPayload
	if err := tools.DecodeJSON(r, &payload); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode get groups JSON",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})

		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Get groups JSON decoded"})

	groups, err := rt.DL.Groups.GetGroups(payload)
	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to get groups from DB",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Groups retrieved from DB"})

	if err := tools.EncodeJSON(w, http.StatusOK, groups); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send groups response",
			Metadata: map[string]any{
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"error": err.Error(),
			},
		})
		return
	}

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Groups sent successfully",
		Metadata: map[string]any{
			"ip":   r.RemoteAddr,
			"path": r.URL.Path,
		},
	})
}

func (rt *Root) NewEvent(w http.ResponseWriter, r *http.Request) {
	var event *models.Event
	if err := tools.DecodeJSON(r, &event); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode new event JSON",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "New event JSON decoded"})

	// verify group creation input
	if err := event.Validate(); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Invalid event input",
			Metadata: map[string]any{
				"event": event,
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"err":   err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Event input validated"})

	// insert user into db
	if err := rt.DL.Events.Insert(event); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to insert event into DB",
			Metadata: map[string]any{
				"event": event,
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"err":   err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Event inserted into DB"})

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "New event created successfully",
		Metadata: map[string]any{
			"event": event,
			"ip":    r.RemoteAddr,
			"path":  r.URL.Path,
		},
	})

	// Send notifications to all group members about the new event
	// Get group information for the notification message
	groupInfo := &models.Group{ID: event.GroupId}
	if err := rt.DL.Groups.GetGroupByID(groupInfo); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Failed to get group info for event notification, using generic message",
			Metadata: map[string]any{
				"group_id": event.GroupId,
				"error":    err.Error(),
			},
		})
		groupInfo.Title = "a group" // fallback
	}

	// Get all group members to notify them
	groupMembers, err := rt.DL.Members.GetGroupMembers(event.GroupId)
	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to get group members for event notification",
			Metadata: map[string]any{
				"group_id": event.GroupId,
				"event_id": event.ID,
				"error":    err.Error(),
			},
		})
	} else {
		// Create notifications for all group members (except the event creator)
		eventCreatorID := 1 // TODO: Get from auth when available
		notificationCount := 0

		for _, member := range groupMembers {
			// Skip the event creator and only notify actual members
			if member.UserID != eventCreatorID && member.Status == "member" {
				notification := &models.Notification{
					UserID:  member.UserID,
					Type:    "group_event_created",
					Message: fmt.Sprintf("A new event '%s' has been created in group '%s'.", event.Title, groupInfo.Title),
					Seen:    false,
				}

				if err := rt.DL.Notifications.Insert(notification); err != nil {
					rt.DL.Logger.Log(models.LogEntry{
						Level:   "ERROR",
						Message: "Failed to create event notification for group member",
						Metadata: map[string]any{
							"user_id":  member.UserID,
							"group_id": event.GroupId,
							"event_id": event.ID,
							"error":    err.Error(),
						},
					})
				} else {
					notificationCount++
				}
			}
		}

		rt.DL.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Event notifications created for group members",
			Metadata: map[string]any{
				"group_id":           event.GroupId,
				"event_id":           event.ID,
				"notifications_sent": notificationCount,
				"total_members":      len(groupMembers),
			},
		})
	}

	if err := tools.EncodeJSON(w, http.StatusCreated, event); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send event response",
			Metadata: map[string]any{
				"event": event,
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"err":   err.Error(),
			},
		})
	}
}

func (rt *Root) EventVote(w http.ResponseWriter, r *http.Request) {
	var vote *models.EventVote
	if err := tools.DecodeJSON(r, &vote); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode new vote JSON",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "New vote JSON decoded"})

	// verify group creation input
	if err := vote.Validate(); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Invalid vote input",
			Metadata: map[string]any{
				"vote": vote,
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Vote input validated"})

	// insert user into db
	if err := rt.DL.Votes.UpsertVote(vote); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to insert vote into DB",
			Metadata: map[string]any{
				"vote": vote,
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Vote inserted/updated in DB"})

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "New vote created successfully",
		Metadata: map[string]any{
			"vote": vote,
			"ip":   r.RemoteAddr,
			"path": r.URL.Path,
		},
	})

	// TODO send notification to group members

	if err := tools.EncodeJSON(w, http.StatusCreated, vote); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send vote response",
			Metadata: map[string]any{
				"vote": vote,
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
	}
}

// GetGroupEvents retrieves events for a specific group
func (rt *Root) GetGroupEvents(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetGroupEvents called")
	payload := &models.EventsPayload{}
	if err := tools.DecodeJSON(r, &payload); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode events payload JSON",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if payload.GroupID <= 0 {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Group ID not provided in request",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
			},
		})
		tools.RespondError(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	// Set user ID for vote status (TODO: Get from auth when available)
	payload.UserID = 1

	events, err := rt.DL.Events.GetEventsByGroup(payload)
	fmt.Println("Retrieved events:", events)
	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to retrieve events from DB",
			Metadata: map[string]any{
				"group_id": payload.GroupID,
				"ip":       r.RemoteAddr,
				"path":     r.URL.Path,
				"err":      err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Events retrieved successfully",
		Metadata: map[string]any{
			"group_id": payload.GroupID,
			"count":    len(events),
			"ip":       r.RemoteAddr,
			"path":     r.URL.Path,
		},
	})

	if err := tools.EncodeJSON(w, http.StatusOK, events); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send events response",
			Metadata: map[string]any{
				"group_id": payload.GroupID,
				"ip":       r.RemoteAddr,
				"path":     r.URL.Path,
				"err":      err.Error(),
			},
		})
	}
}
