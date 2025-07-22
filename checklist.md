### Back-end

#### pkg/db

- [ ] sqlite.go - Finalize InitDB() function (db name)
- [ ] sqlite.go - Add function for creating tables if not exist
- [ ] sqlite.go - delete temporary database function

- [ ] events.go - create GetAll query
- [ ] notifications.go - Save new notification fn
- [ ] notifications.go - finish GetGroupRequests fn

#### pkg/handlers

- [ ] posts.go - on saving new post also save 'ALMOST_PRIVATE' list
- [ ] users.go - /UserStatus - Deal with pending follow requests after status change to PUBLIC
- [ ] groups.go - NewGroup -> deal with member invitations when creating group
- [ ] groups.go - NewGroupRequest - inform websocket about new notification 
- [ ] groups.go - responseGroupRequest - route needs testing
- [ ] groups.go - responseGroupRequest - maybe send notifications if user request approved/denied 
- [ ] groups.go - newGroupPost - route needs testing
- [ ] groups.go - NewGroupInvite - route needs testing / notify websocket abot new notif

- [ ] groups.go - ResponseInviteRequest - route needs testing / notify websocket abot deleted notif
- [ ] users.go - ResponseFollowRequest - route needs testing / notify websocket abot deleted notif

- [ ] events.go - NewEvent - notify websocket abot newEvent created

#### Websockets

- [ ] Plan out notification table (db)
- [ ] Plan chat message handling

