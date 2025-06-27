package handlers

import (
	"net/http"

	"social-network/pkg/models"
	"social-network/pkg/tools"
)

func (rt *Root) NewGroupsHandler() (groupsMux *http.ServeMux) {
	groupsMux = http.NewServeMux()

	groupsMux.HandleFunc("POST /group/new", rt.NewGroup)
	groupsMux.HandleFunc("POST /group/invite", rt.InviteToJoinGroup)
	groupsMux.HandleFunc("POST /group/request", rt.RequestToJoinGroup)
	groupsMux.HandleFunc("POST /group/accept-decline", rt.InviteToJoinGroup)
	groupsMux.HandleFunc("POST	/group/browse", rt.BrowseGroups)
	groupsMux.HandleFunc("POST	/group/event", rt.NewEvent)
	groupsMux.HandleFunc("POST	/group/event/vote", rt.EventVote)

	return groupsMux
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

	// insert user into db
	if err := rt.DL.Groups.InsertGroup(group); err != nil {
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

	// TODO validate payload

	// Check if requester is creator
	requesterID := rt.DL.GetRequesterID(w, r)

	err1 := rt.DL.Members.IsUserInGroup(member.GroupID, requesterID)
	err2 := rt.DL.Groups.IsUserCreator(member.GroupID, requesterID)

	if err1 != nil && err2 != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "forbidden requester isn't the grup creator",
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

	// Insert into group_members as "pending_invite"
	// member.Status = "invited"
	if err := rt.DL.Members.UpsertMember(member); err != nil {
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

	// TODO Add notification to user

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "member invited successfully",
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

	// TODO validate payload

	// Insert into group_members as "pending_request"
	// member.Status = "requested"
	if err := rt.DL.Members.UpsertMember(member); err != nil {
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
	// Notify group creator
}

func (rt *Root) AcceptDeclineGroup(w http.ResponseWriter, r *http.Request) {
	var member *models.GroupMember
	if err := tools.DecodeJSON(r, &member); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode group accept info JSON",
			Metadata: map[string]any{
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"error": err.Error(),
			},
		})
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// TODO validate payload

	// Check if requester is creator
	requesterID := rt.DL.GetRequesterID(w, r)

	switch member.Status {
	case "requested":
		if err := rt.DL.Groups.IsUserCreator(member.GroupID, requesterID); err != nil {
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "ERROR",
				Message: "forbidden requester isn't the grup creator",
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
				Message: "forbidden requester isn't invited to the group",
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
			Message: "invalid",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
			},
		})
		tools.RespondError(w, "invalid payload", http.StatusBadRequest)
		return

	}

	// Update status to "joined" or delete row
	if err := rt.DL.Members.UpsertMember(member); err != nil {
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
	// Notify group creator
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
	posts, err := rt.DL.Groups.GetGroups(payload)
	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to insert group into DB",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tools.EncodeJSON(w, http.StatusOK, posts); err != nil {
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
		Message: "groups  sent successfully",
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

	// verify group creation input
	if err := event.Validate(); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Invalid group input",
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

	// insert user into db
	if err := rt.DL.Events.InsertEvent(event); err != nil {
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

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "New group created successfully",
		Metadata: map[string]any{
			"event": event,
			"ip":    r.RemoteAddr,
			"path":  r.URL.Path,
		},
	})

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

	// verify group creation input
	if err := vote.Validate(); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Invalid group input",
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

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "New vote created successfully",
		Metadata: map[string]any{
			"vote": vote,
			"ip":   r.RemoteAddr,
			"path": r.URL.Path,
		},
	})

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
