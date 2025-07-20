#!/bin/bash

# Comprehensive Notification System Test
echo "🔧 Starting Comprehensive Notification System Test"
echo "=================================================="

BASE_URL="http://localhost:8080"

# Test 1: Check if server is running
echo -e "\n1️⃣ Testing server connectivity..."
response=$(curl -s -w "%{http_code}" -o /dev/null "$BASE_URL/notifications/unread-count")
if [ "$response" -eq 200 ]; then
    echo "✅ Server is running and accessible"
else
    echo "❌ Server is not responding (HTTP $response)"
    exit 1
fi

# Test 2: Get unread count
echo -e "\n2️⃣ Testing unread count endpoint..."
unread_response=$(curl -s "$BASE_URL/notifications/unread-count")
echo "Response: $unread_response"

# Test 3: Test fetch notifications with different payloads
echo -e "\n3️⃣ Testing fetch notifications endpoint..."

echo "Testing with minimal payload..."
fetch_response1=$(curl -s -X POST "$BASE_URL/notifications/fetch" \
  -H "Content-Type: application/json" \
  -d '{}')
echo "Response: $fetch_response1"

echo "Testing with complete payload..."
fetch_response2=$(curl -s -X POST "$BASE_URL/notifications/fetch" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "2",
    "type": "all",
    "start": 0,
    "num_of_items": 10
  }')
echo "Response: $fetch_response2"

# Test 4: Test mark as seen
echo -e "\n4️⃣ Testing mark as seen endpoint..."
mark_seen_response=$(curl -s -X POST "$BASE_URL/notifications/mark-seen" \
  -H "Content-Type: application/json" \
  -d '{
    "notification_id": 1
  }')
echo "Response: $mark_seen_response"

# Test 5: Test mark all as seen
echo -e "\n5️⃣ Testing mark all as seen endpoint..."
mark_all_response=$(curl -s -X POST "$BASE_URL/notifications/mark-all-seen" \
  -H "Content-Type: application/json" \
  -d '{}')
echo "Response: $mark_all_response"

# Test 6: Check database directly
echo -e "\n6️⃣ Checking database contents..."
echo "Notifications for user_id 2:"
cd /home/mbakhcha/mokZwina/backend
sqlite3 pkg/db/data.db "SELECT id, user_id, type, message, seen FROM notifications WHERE user_id = 2 LIMIT 5;"

echo -e "\nNotifications for user_id 1:"
sqlite3 pkg/db/data.db "SELECT id, user_id, type, message, seen FROM notifications WHERE user_id = 1 LIMIT 5;"

echo -e "\nAll notifications:"
sqlite3 pkg/db/data.db "SELECT id, user_id, type, message, seen FROM notifications ORDER BY created_at DESC LIMIT 5;"

# Test 6b: Test direct database query simulation
echo -e "\n6️⃣b Testing SQL query simulation..."
echo "Query: SELECT id, user_id, type, message, seen, created_at FROM notifications WHERE user_id = 2 ORDER BY created_at DESC LIMIT 10 OFFSET 0"
sqlite3 pkg/db/data.db "SELECT id, user_id, type, message, seen, created_at FROM notifications WHERE user_id = 2 ORDER BY created_at DESC LIMIT 10 OFFSET 0;"

# Test 7: WebSocket connectivity test
echo -e "\n7️⃣ Testing WebSocket endpoint connectivity..."
ws_response=$(curl -s -w "%{http_code}" -o /dev/null "$BASE_URL/connect")
echo "WebSocket endpoint HTTP response: $ws_response"

echo -e "\n🏁 Test completed!"
echo "===================================="