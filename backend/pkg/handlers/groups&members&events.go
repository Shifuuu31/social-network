package handlers

import (
	"net/http"

	"social-network/pkg/models"
	"social-network/pkg/tools"

	"github.com/google/uuid"
)

func (rt *Root) NewGroupsHandler() (groupsMux *http.ServeMux) {
	groupsMux = http.NewServeMux()

	groupsMux.HandleFunc("POST /group/new", rt.NewGroup)
	groupsMux.HandleFunc("POST /group/invite", rt.InviteToJoinGroup)
	groupsMux.HandleFunc("POST /group/request", rt.RequestToJoinGroup)
	groupsMux.HandleFunc("POST /group/accept-decline-request", rt.AcceptDeclineRequestToJoin)
	groupsMux.HandleFunc("POST /group/accept-decline-invitation", rt.AcceptDeclineInvitation)
	groupsMux.HandleFunc("POST /browse", rt.BrowseGroups)
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

	group.CreatorID = rt.DL.GetRequesterID(w, r)
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

	// insert user into db
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
	requesterID := rt.DL.GetRequesterID(w, r)

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

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Member invited successfully",
		Metadata: map[string]any{
			"ip":   r.RemoteAddr,
			"path": r.URL.Path,
		},
	})
	// TODO Add notification

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

	// Insert into group_members as "pending_request"
	// member.Status = "requested"
	member.UserID = rt.DL.GetRequesterID(w, r)
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

	// TODO Notify group creator
}

// Notify group creator
func (rt *Root) AcceptDeclineRequestToJoin(w http.ResponseWriter, r *http.Request) {
	var member *models.GroupMember
	if err := tools.DecodeJSON(r, &member); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:    "ERROR",
			Message:  "Failed to decode group accept/decline JSON (creator)",
			Metadata: map[string]any{"ip": r.RemoteAddr, "path": r.URL.Path, "error": err.Error()},
		})
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Decoded join request for group"})

	if member.PrevStatus != "requested" {
		tools.RespondError(w, "Only group creator can act on 'requested' status", http.StatusBadRequest)
		return
	}

	requesterID := rt.DL.GetRequesterID(w, r)
	if err := rt.DL.Groups.IsUserCreator(member.GroupID, requesterID); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:    "ERROR",
			Message:  "Forbidden: requester isn't the group creator",
			Metadata: map[string]any{"ip": r.RemoteAddr, "path": r.URL.Path, "error": err.Error()},
		})
		tools.RespondError(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Accept = change to "member", Decline = "declined" or delete
	if err := rt.DL.Members.Upsert(member); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:    "ERROR",
			Message:  "Failed to update group membership",
			Metadata: map[string]any{"ip": r.RemoteAddr, "path": r.URL.Path, "err": err.Error()},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Group creator processed request"})
}

func (rt *Root) AcceptDeclineInvitation(w http.ResponseWriter, r *http.Request) {
	var member *models.GroupMember
	if err := tools.DecodeJSON(r, &member); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:    "ERROR",
			Message:  "Failed to decode group invitation accept/decline JSON (user)",
			Metadata: map[string]any{"ip": r.RemoteAddr, "path": r.URL.Path, "error": err.Error()},
		})
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Decoded invitation for group"})

	if member.PrevStatus != "invited" {
		tools.RespondError(w, "Only invited users can accept/decline", http.StatusBadRequest)
		return
	}

	requesterID := rt.DL.GetRequesterID(w, r)
	if requesterID != member.UserID {
		rt.DL.Logger.Log(models.LogEntry{
			Level:    "ERROR",
			Message:  "Forbidden: requester is not the invited user",
			Metadata: map[string]any{"ip": r.RemoteAddr, "path": r.URL.Path},
		})
		tools.RespondError(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := rt.DL.Members.Upsert(member); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:    "ERROR",
			Message:  "Failed to update group membership (invited user)",
			Metadata: map[string]any{"ip": r.RemoteAddr, "path": r.URL.Path, "err": err.Error()},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Invited user responded to invitation"})
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

	posts, err := rt.DL.Groups.GetGroups(rt.DL.GetRequesterID(w, r), payload)
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

	// TODO send notification to group members

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
