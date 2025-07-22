# Group Member Management Features

## Overview
This document describes the new features implemented for group member management in the social network application.

## Features Implemented

### 1. User Filtering for Invitations
- **Feature**: When fetching users for group invitations, only users who have never been invited or been members of a group are returned.
- **Endpoints**: 
  - `GET /group/{id}/available-users` - Get all users not in the specified group
  - `GET /group/{id}/search-users?q={query}` - Search users not in the specified group
- **Implementation**: Uses SQL `NOT IN` clause to exclude users who have any record in the `group_members` table for the specified group.

### 2. Decline Request Cleanup
- **Feature**: When a user declines a join request or invitation, their record is deleted from the `group_members` table instead of updating the status.
- **Endpoint**: `POST /group/accept-decline`
- **Implementation**: Modified `AcceptDeclineGroup` function to check if `member.Status == "declined"` and call `rt.DL.Members.Delete()` instead of `Upsert()`.

### 3. Search Functionality for User Invitations
- **Feature**: Added search capability when looking for users to invite to a group.
- **Endpoint**: `GET /group/{id}/search-users?q={query}`
- **Search criteria**: First name, last name, nickname, or full name (first + last)
- **Implementation**: Uses SQL `LIKE` clauses with pattern matching.

### 4. Not Joined Groups Filtering
- **Feature**: In the "not joined" section, only show groups where the user has never interacted (not present in `group_members` table).
- **Endpoint**: `POST /group/browse` with `type: "not_joined"`
- **Implementation**: Enhanced `GetGroups` function in the group model to filter out groups where the user has any record in `group_members`.

## Database Changes
No schema changes were required. The implementation uses existing tables:
- `users` table for user information
- `groups` table for group information  
- `group_members` table for membership tracking

## API Changes

### New Endpoints
1. `GET /group/{id}/available-users`
   - Returns users who can be invited to the group
   - Requires group membership or creator privileges

2. `GET /group/{id}/search-users?q={query}`
   - Searches available users by name/nickname
   - Requires group membership or creator privileges

### Modified Endpoints
1. `POST /group/accept-decline`
   - Now deletes records when status is "declined"
   - Uses `prev_status` to determine authorization

2. `POST /group/browse`
   - Enhanced with `type: "not_joined"` option
   - Filters groups where user has no interaction history

## Authorization
- Available users endpoints require the requester to be either:
  - A member of the group, OR
  - The creator of the group
- Accept/decline logic checks:
  - For "requested" → Must be group creator
  - For "invited" → Must be the invited user

## Testing
Test files created:
- `/testers/group_members_test.js` - Test new endpoints
- `/testers/decline_test.js` - Test decline functionality and cleanup

## Error Handling
- Proper HTTP status codes (400, 403, 500)
- Detailed logging for debugging
- Graceful fallbacks for notification failures

## Benefits
1. **Clean Data**: Declined users don't clutter the database
2. **Better UX**: Users can find available people easily with search
3. **Accurate Filtering**: "Not joined" only shows truly uninteracted groups
4. **Security**: Proper authorization checks for all operations
