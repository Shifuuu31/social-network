import { useNotificationStore } from './stores/notificationStore'

const store = useNotificationStore()

// Create a follow request notification
store.createFollowRequest(
  { id: 1, username: 'john_doe' },
  'recipient123'
)

// Create a group invitation notification
store.createGroupInvitation(
  { id: 1, name: 'Study Group' },
  { id: 2, username: 'jane_doe' },
  'recipient123'
)

// Create a group join request notification
store.createGroupJoinRequest(
  { id: 1, name: 'Study Group' },
  { id: 3, username: 'bob_smith' },
  'recipient123'
)

// Create a group event notification
store.createGroupEvent(
  { id: 1, name: 'Study Group' },
  { id: 1, title: 'Study Session' },
  'recipient123'
)