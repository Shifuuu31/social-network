### Prerequisites
if want to edit u need to make sure that u are in main branch and theres no one working on this file
install todo tree ext link: https://marketplace.visualstudio.com/items?itemName=Gruntfuggly.todo-tree


### Manual

  * 
  >>> to-do`
  * `TODO <username> >>> on going`
  * `[x] >>> done`
---

## 🖥️ BACKEND TASKS (Go + SQLite + WebSocket)

### 📦 Setup

Initialize Go project structure
Create `main.go` or `server.go` entrypoint
Install required packages:

  * Gorilla WebSocket
  * `golang-migrate`
  * `sqlite3`
  * `bcrypt`
  * `uuid` (gofrs or google)

### 🧩 Folder Structure

```
backend/
├── server.go
├── pkg/
│   ├── db/
│   │   ├── sqlite.go
│   │   └── migrations/sqlite/
│   │       ├── 000001_create_users_table.up.sql
│   │       └── ...
│   ├── models/
│   ├── handlers/
│   └── middleware/
```

### 🗃️ Database & Migrations

Design ERD
Implement DB connection (sqlite.go)
Add migrations:

  * Users, Sessions, Followers, Posts, Comments
  Groups, Messages, Notifications, Events

### 🔐 Authentication

Register: hash password, validate input
Login: verify password, start session
Sessions & cookies (custom or middleware)
Logout endpoint
Middleware to protect routes

### 👤 Profiles & Follow System

Get user profile info
Update profile visibility (public/private)
Follow/unfollow requests
Accept/decline follow requests

### 📝 Posts & Comments

Create/edit/delete post
Add post image (JPEG, PNG, GIF)
Post privacy handling (public, followers-only, selected)
Create comments on posts
Serve feed with filtered access

### 👥 Groups

Create group
Invite or request to join
Accept/decline invites and join requests
Group posts and comments
Create/view group events
RSVP: Going / Not Going

### 💬 Chat & WebSocket

Set up WebSocket hub (Gorilla)
Private chat between users (followed or public)
Group chat system
Store messages in DB
Send emojis via WebSocket

### 🔔 Notifications

Notification table/model
Trigger notification on:
Follow request
Group invite or join request
New event in a group
Fetch unread notifications

### 🐳 Docker

Dockerfile for backend
Expose port 8080
Connect DB (volume or internal file)
Allow frontend to communicate (CORS)

---

## 🌐 FRONTEND TASKS (Next.js)

### 🚀 Setup

`npx create-next-app@latest frontend`
Install packages
Create layout (Header/Navbar/Pages)
Configure environment variables

### 📁 Pages & Routing

`/` - Home feed
`/auth/signup`
`/auth/login`
`/profile/[username]`
`/groups`
`/groups/[id]`
`/chat`
`/notifications`

### 🔄 API & Session

Manage auth via cookies or tokens
Global user state with React Context or hook
Handle redirects for unauthorized pages

### 👤 Profile UI

Registration & login forms
View/edit profile
Toggle public/private
Follow/unfollow users
Accept/decline requests

### 📝 Post & Feed UI

Create/edit/delete post (with image upload)
Comment on posts
Show feed with privacy filtering

### 👥 Groups UI

Create & browse groups
Join or invite flow
Group posts & comments
Create/join events

### 💬 Chat UI

WebSocket setup
Private chat
Group chat
Emoji picker

### 🔔 Notifications UI

Notification list on navbar
Mark as read
Real-time updates via polling or socket

### 💅 Styling & UX

Responsive layout
Button/loaders/feedback
Form validation
Toast messages for errors/success

### 🐳 Docker

Dockerfile for frontend
Expose frontend on port 3000
Set backend URL via ENV
Communicate with backend

---

## 🔗 BACKEND + FRONTEND INTEGRATION

### 🌐 Communication

Set backend API base URL in `.env.local`
Handle session cookies (send `credentials: include`)
CORS middleware in Go to allow frontend origin

### 🔁 Auth Flow

Register → Create session cookie
Login → Set session
Frontend reads login state and updates UI
Logout → Clear session cookie

### 📤 Forms & File Uploads

Use `FormData` for avatar/post image upload
Backend handles multipart/form-data
Return image URLs to display in UI

### 🧠 Real-time Chat & Notifications

Frontend connects to WebSocket
Auth handshake over WS if needed
Push messages to open chat
Push notifications for events

---

