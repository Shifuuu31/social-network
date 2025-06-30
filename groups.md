Here's how to **break down the backend for the "Groups" feature** in your social network project. We'll follow a clear structure to help you implement everything step-by-step.

---

## 🔧 GROUPS FEATURE – BACKEND BREAKDOWN

### ✅ Core Functional Requirements

1. **Create Group**
2. **Invite User to Group**
3. **Accept/Reject Group Invite**
4. **Request to Join Group**
5. **Accept/Reject Join Request (by creator)**
6. **Browse Groups**
7. **Group Posts and Comments (just like user posts, but scoped to the group)**
8. **Create Event inside Group**
9. **Vote on Event (Going / Not Going)**
10. **Group Chat (WebSocket)**

---

## 1. 📦 Folder Structure (Suggested)

```
backend/
├── pkg/
│   ├── db/
│   │   └── models/
│   │       ├── group.go
│   │       ├── group_members.go
│   │       ├── group_events.go
│   │       ├── group_posts.go
│   │       └── ...
│   ├── handlers/
│   │   └── groups.go
│   └── ws/         # for WebSocket group chat
│       └── group_chat.go
```

---

## 2. 🗃️ Database Models

### `groups`

```sql
CREATE TABLE groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    creator_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (creator_id) REFERENCES users(id)
);
```

### `group_members`

```sql
CREATE TABLE group_members (
    user_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    role TEXT NOT NULL CHECK(role IN ('member', 'creator')),
    status TEXT NOT NULL CHECK(status IN ('pending_invite', 'pending_request', 'joined')),
    PRIMARY KEY (user_id, group_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);
```

> `status` helps distinguish invitations vs. join requests.

---

### `group_posts`

```sql
CREATE TABLE group_posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    image_uuid TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### `group_post_comments`

```sql
CREATE TABLE group_post_comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES group_posts(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

---

### `group_events`

```sql
CREATE TABLE group_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    creator_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    event_time DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (creator_id) REFERENCES users(id)
);
```

### `group_event_votes`

```sql
CREATE TABLE group_event_votes (
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    vote TEXT CHECK(vote IN ('going', 'not_going')),
    PRIMARY KEY (event_id, user_id),
    FOREIGN KEY (event_id) REFERENCES group_events(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

---

## 3. 🧠 Handlers Logic (Go)

### Group Creation

```go
func (rt *Root) CreateGroup(w http.ResponseWriter, r *http.Request) {
    // Parse title, description
    // Insert into groups table
    // Add creator to group_members as "creator", "joined"
}
```

### Invite Member

```go
func (rt *Root) InviteToGroup(w http.ResponseWriter, r *http.Request) {
    // Check if requester is creator or member
    // Insert into group_members as "pending_invite"
    // Add notification to user
}
```

### Accept/Reject Invite or Request

```go
func (rt *Root) HandleGroupInviteResponse(w http.ResponseWriter, r *http.Request) {
    // Update status to "joined" or delete row
}
```

### Request to Join

```go
func (rt *Root) RequestToJoinGroup(w http.ResponseWriter, r *http.Request) {
    // Insert into group_members as "pending_request"
    // Notify group creator
}
```

---

## 4. 🔔 Notifications

| Action                   | Notification Target |
| ------------------------ | ------------------- |
| Invited to group         | Invited User        |
| Someone requests to join | Group Creator       |
| Event Created            | All group members   |

---

## 5. 💬 WebSocket Group Chat

You can use a single handler per group like:

```go
wsHandler := func(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID") // or from query param
	// validate user is in group
	// add to group-specific broadcast hub
}
```

Use Gorilla WebSocket with a simple hub pattern where `hub[groupID]` keeps a list of clients per group.

---

## 6. 📦 Migrations

Example filenames:

```
000005_create_groups_table.up.sql
000006_create_group_members_table.up.sql
000007_create_group_posts_table.up.sql
...
```

---

## ✅ API Endpoints

| Method | Path                         | Description           |
| ------ | ---------------------------- | --------------------- |
| POST   | `/group/create`              | Create new group      |
| POST   | `/group/invite`              | Invite user           |
| POST   | `/group/request`             | Request to join       |
| POST   | `/group/accept`              | Accept invite/request |
| POST   | `/group/reject`              | Reject invite/request |
| GET    | `/group/browse`              | List groups           |
| POST   | `/group/post`                | Create group post     |
| POST   | `/group/comment`             | Comment on post       |
| POST   | `/group/event`               | Create group event    |
| POST   | `/group/event/vote`          | Vote on event         |
| GET    | `/group/chat/ws?groupID=123` | WebSocket group chat  |

---

## ✅ Next Steps

* [ ] Define JSON input/output structs
* [ ] Implement authentication guard (middleware)
* [ ] Add tests and seed data for groups
* [ ] Add notifications to group actions

---

Let me know if you want me to:

* generate SQL migration files,
* help with the Go structs & handlers,
* or implement frontend integration.

Would you like code templates for any of these parts?
