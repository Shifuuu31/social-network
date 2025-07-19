# Chat System Documentation

## Overview

This document describes the complete chat system implementation for the social network application, including both frontend (Vue.js) and backend (Go) components. The system supports real-time messaging between users who follow each other.

## System Architecture

### Backend (Go)

#### 1. Database Schema

**Messages Table:**
```sql
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER,
    group_id INTEGER,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(receiver_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(group_id) REFERENCES groups(id) ON DELETE CASCADE
);
```

**Key Features:**
- Auto-incrementing message IDs
- Foreign key constraints to users table
- Support for both direct messages and group messages
- Automatic timestamp creation

#### 2. Message Model (`backend/pkg/models/message.go`)

**Core Functions:**

- **`Validate()`**: Validates message data (sender/receiver IDs, content length, etc.)
- **`CanSendMessage()`**: Checks if users can send messages (must follow each other)
- **`CreateMessage()`**: Saves new message to database
- **`GetConversation()`**: Retrieves messages between two users
- **`GetRecentConversations()`**: Gets recent conversations for a user
- **`DeleteMessage()`**: Deletes a message (sender only)

**Debug Logging:**
- Comprehensive debug logging added to track message flow
- Logs validation, database operations, and error conditions

#### 3. Chat Handler (`backend/pkg/handlers/chat.go`)

**API Endpoints:**

- **`POST /api/chat/send`**: Send a new message
- **`POST /api/chat/conversation`**: Get conversation between users
- **`GET /api/chat/recent`**: Get recent conversations
- **`DELETE /api/chat/delete`**: Delete a message
- **`GET /api/chat/ws`**: WebSocket connection endpoint

**SendMessage Handler Flow:**
1. Validate HTTP method (POST only)
2. Extract requester ID from session
3. Decode request body (receiver_id, content)
4. Check if users can send messages (follow relationship)
5. Save message to database
6. Broadcast via WebSocket (if receiver online)
7. Return success response with message ID

#### 4. WebSocket Implementation (`backend/pkg/websocket/hub.go`)

**Components:**
- **Hub**: Manages all WebSocket connections
- **Client**: Represents individual WebSocket connections
- **Message**: WebSocket message structure

**Features:**
- Real-time message broadcasting
- Connection management (register/unregister)
- Message serialization/deserialization
- Error handling and recovery

#### 5. Middleware Integration

**Authentication:**
- Session-based authentication using cookies
- User ID extraction from session context
- Unauthorized request handling

**CORS:**
- Cross-origin resource sharing configuration
- Support for frontend development server
- Credential inclusion for authentication

### Frontend (Vue.js)

#### 1. Chat Service (`frontend/src/services/chatService.js`)

**Core Functions:**

- **`connect(userId)`**: Establishes WebSocket connection
- **`sendMessage(receiverId, content)`**: Sends message via WebSocket
- **`sendMessageAPI(receiverId, content)`**: Sends message via HTTP API
- **`getConversation(otherUserId, limit, offset)`**: Retrieves conversation
- **`getRecentConversations(limit)`**: Gets recent conversations
- **`deleteMessage(messageId)`**: Deletes a message

**Debug Features:**
- Comprehensive logging for all API calls
- Raw response text logging
- Error handling and fallback mechanisms

#### 2. Chat Components

**ChatList Component (`frontend/src/components/chat/ChatList.vue`):**
- Displays list of conversations
- Shows recent conversations and following users
- Handles conversation selection
- WebSocket connection status display

**ChatWindow Component (`frontend/src/components/chat/ChatWindow.vue`):**
- Message display and input
- Real-time message sending
- Message deletion functionality
- Auto-scroll to bottom
- Temporary message handling for better UX

**Chat View (`frontend/src/views/Chat.vue`):**
- Main chat interface
- Debug panel for troubleshooting
- Component coordination

#### 3. Message Flow

**Sending Messages:**
1. User types message and clicks send
2. Temporary message added to UI immediately
3. HTTP API call made to backend
4. Backend validates and saves message
5. Temporary message updated with real database ID
6. WebSocket broadcast to online users

**Receiving Messages:**
1. WebSocket receives new message
2. Message added to conversation if relevant
3. UI updates automatically
4. Auto-scroll to show new message

#### 4. Debug System

**Debug Panel:**
- Real-time log display
- Toggle debug mode
- Clear logs functionality
- Network request monitoring

**Debug Logging:**
- Component-level logging
- API call tracking
- WebSocket event logging
- Error condition logging

## Key Features Implemented

### 1. Authentication & Authorization
- Session-based authentication
- User permission checking (follow relationships)
- Secure message sending (only between following users)

### 2. Real-time Messaging
- WebSocket connections for instant messaging
- HTTP API fallback for reliability
- Message broadcasting to online users

### 3. Message Management
- Message creation and storage
- Conversation retrieval with pagination
- Message deletion (sender only)
- Recent conversations list

### 4. User Experience
- Temporary message display for immediate feedback
- Auto-scroll to latest messages
- Loading states and error handling
- Responsive design

### 5. Debugging & Monitoring
- Comprehensive logging throughout the system
- Debug panel for frontend troubleshooting
- Network request monitoring
- Error tracking and reporting

## Technical Challenges & Solutions

### 1. Recursive Update Issue
**Problem:** Vue component causing infinite re-renders
**Solution:** Simplified `getOtherUser` function to avoid reactive dependencies

### 2. WebSocket vs HTTP API
**Problem:** WebSocket connection issues
**Solution:** Implemented HTTP API as primary method with WebSocket as enhancement

### 3. Message Synchronization
**Problem:** Temporary vs permanent message IDs
**Solution:** Use temporary IDs for immediate UI feedback, update with real IDs from backend

### 4. Authentication Flow
**Problem:** Session token management
**Solution:** Proper cookie handling and session validation

### 5. CORS Configuration
**Problem:** Cross-origin request issues
**Solution:** Configured CORS middleware for development environment

## Database Relationships

### Follow System
- Users must follow each other to send messages
- Checked via `follow_request` table with 'accepted' status
- Bidirectional relationship validation

### Message Storage
- Messages linked to sender and receiver via foreign keys
- Automatic timestamp creation
- Cascade deletion for user cleanup

## Security Considerations

### 1. Input Validation
- Message content length limits (1000 characters)
- User ID validation
- SQL injection prevention via parameterized queries

### 2. Authorization
- Session-based authentication
- User permission checking before message operations
- Sender-only message deletion

### 3. Data Integrity
- Foreign key constraints
- Transaction handling
- Error logging and monitoring

## Performance Optimizations

### 1. Database Queries
- Indexed foreign key relationships
- Pagination for conversation retrieval
- Efficient conversation grouping

### 2. Frontend Performance
- Temporary message display for immediate feedback
- Efficient component updates
- WebSocket connection pooling

### 3. Caching
- Session caching
- User data caching
- Conversation caching

## Testing & Debugging

### 1. Backend Testing
- Database schema validation
- API endpoint testing
- Authentication flow testing
- Message creation and retrieval testing

### 2. Frontend Testing
- Component rendering
- API integration testing
- WebSocket connection testing
- User interaction testing

### 3. Integration Testing
- End-to-end message flow
- Cross-browser compatibility
- Mobile responsiveness

## Future Enhancements

### 1. Real-time Features
- Typing indicators
- Message read receipts
- Online/offline status
- Push notifications

### 2. Advanced Messaging
- File attachments
- Image sharing
- Message reactions
- Message threading

### 3. Group Chat
- Multi-user conversations
- Group management
- Role-based permissions

### 4. Performance Improvements
- Message pagination
- Lazy loading
- WebSocket connection optimization
- Database query optimization

## Deployment Considerations

### 1. Environment Configuration
- Database connection settings
- WebSocket server configuration
- CORS settings for production
- Logging levels

### 2. Scaling
- Database connection pooling
- WebSocket server clustering
- Load balancing
- CDN integration

### 3. Monitoring
- Application performance monitoring
- Error tracking
- User analytics
- System health checks

## Conclusion

The chat system provides a robust, real-time messaging solution with comprehensive debugging capabilities. The implementation successfully handles authentication, message persistence, real-time communication, and user experience considerations. The modular architecture allows for easy maintenance and future enhancements.

Key achievements:
- ✅ Real-time messaging between authenticated users
- ✅ Robust error handling and debugging
- ✅ Secure message storage and retrieval
- ✅ Responsive and user-friendly interface
- ✅ Comprehensive logging and monitoring
- ✅ Scalable architecture for future growth 