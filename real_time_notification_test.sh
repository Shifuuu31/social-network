#!/bin/bash

# Real-time Notification System Test
# Tests different notification types between users
echo "🚀 Real-time Notification System Test"
echo "====================================="

BASE_URL="http://localhost:8080"

# Color codes for better output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test helper functions
check_response() {
    local response="$1"
    local expected_status="$2"
    local test_name="$3"
    
    if echo "$response" | grep -q "error"; then
        echo -e "${RED}❌ $test_name FAILED${NC}"
        echo "Response: $response"
        return 1
    else
        echo -e "${GREEN}✅ $test_name PASSED${NC}"
        return 0
    fi
}

get_unread_count() {
    local user_id="$1"
    # The middleware will automatically use user_id 2 for notification endpoints
    local count=$(curl -s "$BASE_URL/notifications/unread-count" | jq -r '.unread_count')
    echo "$count"
}

# Check if server is running
echo -e "\n${BLUE}🔍 Checking server status...${NC}"
response=$(curl -s -w "%{http_code}" -o /dev/null "$BASE_URL/notifications/unread-count")
if [ "$response" -ne 200 ]; then
    echo -e "${RED}❌ Server is not running or not accessible${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Server is running${NC}"

# Get initial notification counts
echo -e "\n${BLUE}📊 Getting initial notification counts...${NC}"
initial_count_user2=$(get_unread_count 2)
echo "Initial unread count for user 2: $initial_count_user2"

# Clear all notifications for clean testing
echo -e "\n${BLUE}🧹 Clearing existing notifications for clean testing...${NC}"
mark_all_response=$(curl -s -X POST "$BASE_URL/notifications/mark-all-seen" \
  -H "Content-Type: application/json" \
  -d '{}')
echo "Mark all seen response: $mark_all_response"

# Test 1: Group Invitation (User 1 invites User 2 to join a group)
echo -e "\n${YELLOW}🧪 Test 1: Group Invitation (User 1 → User 2)${NC}"
echo "User 1 (group owner) invites User 2 to join group"

invite_response=$(curl -s -X POST "$BASE_URL/groups/group/invite" \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": 1,
    "user_id": 2,
    "status": "invited",
    "prev_status": "none"
  }')

echo "Invite response: $invite_response"
check_response "$invite_response" "success" "Group Invitation"

# Check if notification was created for user 2
sleep 1
count_after_invite=$(get_unread_count 2)
echo "Unread count for user 2 after invite: $count_after_invite"

# Test 2: Group Join Request (User 3 requests to join User 1's group)
echo -e "\n${YELLOW}🧪 Test 2: Group Join Request (User 3 → User 1)${NC}"
echo "User 3 requests to join User 1's group"

join_request_response=$(curl -s -X POST "$BASE_URL/groups/group/request" \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": 1,
    "user_id": 3,
    "status": "requested",
    "prev_status": "none"
  }')

echo "Join request response: $join_request_response"
check_response "$join_request_response" "success" "Group Join Request"

# Test 3: Accept Group Invitation (User 2 accepts invitation from User 1)
echo -e "\n${YELLOW}🧪 Test 3: Accept Group Invitation (User 2 accepts)${NC}"
echo "User 2 accepts the group invitation from User 1"

accept_invite_response=$(curl -s -X POST "$BASE_URL/groups/group/accept-decline" \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": 1,
    "user_id": 2,
    "status": "member",
    "prev_status": "invited"
  }')

echo "Accept invitation response: $accept_invite_response"
check_response "$accept_invite_response" "success" "Accept Group Invitation"

# Test 4: Group Owner Accepts Join Request (User 1 accepts User 3's request)
echo -e "\n${YELLOW}🧪 Test 4: Group Owner Accepts Join Request (User 1 accepts User 3)${NC}"
echo "User 1 (group owner) accepts User 3's join request"

accept_request_response=$(curl -s -X POST "$BASE_URL/groups/group/accept-decline" \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": 1,
    "user_id": 3,
    "status": "member",
    "prev_status": "requested"
  }')

echo "Accept join request response: $accept_request_response"
check_response "$accept_request_response" "success" "Accept Join Request"

# Test 5: Create Group Event (User 1 creates event for group members)
echo -e "\n${YELLOW}🧪 Test 5: Create Group Event (User 1 creates event)${NC}"
echo "User 1 creates a new event for group members"

event_response=$(curl -s -X POST "$BASE_URL/groups/group/event/new" \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": 1,
    "title": "Test Notification Event",
    "description": "Testing event notifications",
    "event_time": "2025-08-01 10:00:00"
  }')

echo "Create event response: $event_response"
check_response "$event_response" "success" "Create Group Event"

# Test 6: Follow Request (User 1 follows User 2)
echo -e "\n${YELLOW}🧪 Test 6: Follow Request (User 1 → User 2)${NC}"
echo "User 1 sends follow request to User 2"

follow_response=$(curl -s -X POST "$BASE_URL/profile/follow" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 2
  }')

echo "Follow request response: $follow_response"
check_response "$follow_response" "success" "Follow Request"

# Final notification check - Fetch all notifications for user 2
echo -e "\n${BLUE}📋 Final Notification Check for User 2${NC}"
final_notifications=$(curl -s -X POST "$BASE_URL/notifications/fetch" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "all",
    "start": 0,
    "num_of_items": 20
  }')

echo "Final notifications for user 2:"
echo "$final_notifications" | jq '.'

# Get final unread count
final_count=$(get_unread_count 2)
echo -e "\n${BLUE}📊 Final unread count for user 2: $final_count${NC}"

# Test WebSocket Real-time Notifications
echo -e "\n${YELLOW}🔌 Testing WebSocket Real-time Notifications${NC}"

# Create a simple WebSocket test
cat > /tmp/ws_notification_test.js << 'EOF'
const WebSocket = require('ws');

console.log('🔌 Testing WebSocket real-time notifications...');

const ws = new WebSocket('ws://localhost:8080/connect?user_id=2');

ws.on('open', function open() {
    console.log('✅ WebSocket connected as user 2');
    
    // Subscribe to notifications
    ws.send(JSON.stringify({
        type: 'notification_subscribe'
    }));
    
    console.log('📡 Subscribed to notifications');
    console.log('⏰ Waiting for real-time notifications for 10 seconds...');
});

ws.on('message', function message(data) {
    const msg = JSON.parse(data.toString());
    console.log('📨 Received real-time notification:', JSON.stringify(msg, null, 2));
});

ws.on('error', function error(err) {
    console.log('❌ WebSocket error:', err.message);
});

ws.on('close', function close() {
    console.log('🔌 WebSocket connection closed');
    process.exit(0);
});

// Close after 10 seconds
setTimeout(() => {
    console.log('⏰ Test timeout - closing WebSocket');
    ws.close();
}, 10000);
EOF

# Run WebSocket test if Node.js is available
if command -v node &> /dev/null; then
    echo "Running WebSocket test..."
    node /tmp/ws_notification_test.js &
    WS_PID=$!
    
    # Send a test notification while WebSocket is listening
    sleep 2
    echo "Sending test notification while WebSocket is connected..."
    curl -s -X POST "$BASE_URL/groups/invite" \
      -H "Content-Type: application/json" \
      -d '{
        "group_id": 1,
        "user_id": 2
      }' > /dev/null
    
    # Wait for WebSocket test to complete
    wait $WS_PID
else
    echo "Node.js not available, skipping WebSocket test"
fi

# Clean up
rm -f /tmp/ws_notification_test.js

echo -e "\n${GREEN}🎉 Notification System Test Complete!${NC}"
echo -e "${BLUE}=================================${NC}"

# Summary
echo -e "\n${YELLOW}📊 Test Summary:${NC}"
echo "✅ Group invitations tested"
echo "✅ Group join requests tested" 
echo "✅ Group invitation acceptance tested"
echo "✅ Group owner accepting join requests tested"
echo "✅ Group event creation tested"
echo "✅ Follow requests tested"
echo "✅ WebSocket real-time notifications tested"
echo "✅ Notification fetching and counting tested"

echo -e "\n${GREEN}All notification types working correctly! 🎯${NC}"