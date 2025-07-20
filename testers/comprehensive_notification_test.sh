#!/bin/bash

echo "🧪 Comprehensive Notification System Test"
echo "========================================="

BASE_URL="http://localhost:8080"

echo "📝 Step 1: Creating test users..."

# Create user 1
echo "Creating user 1..."
curl -s -X POST ${BASE_URL}/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user1@test.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User1",
    "date_of_birth": "1990-01-01",
    "avatar": "",
    "nickname": "testuser1",
    "about_me": "Test user 1"
  }' | jq '.'

# Create user 2
echo "Creating user 2..."
curl -s -X POST ${BASE_URL}/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user2@test.com",
    "password": "password123", 
    "first_name": "Test",
    "last_name": "User2",
    "date_of_birth": "1990-01-01",
    "avatar": "",
    "nickname": "testuser2",
    "about_me": "Test user 2"
  }' | jq '.'

echo ""
echo "👥 Step 2: Creating a test group..."

# Create a group
curl -s -X POST ${BASE_URL}/groups/group/new \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Notification Group",
    "description": "A group for testing notifications",
    "type": "public"
  }' | jq '.'

echo ""
echo "📬 Step 3: Testing notification endpoints..."

# Test get notifications (should be empty)
echo "Getting notifications for user 2..."
curl -s -X GET ${BASE_URL}/api/notifications \
  -H "Accept: application/json" | jq '.'

# Test unread count
echo "Getting unread count for user 2..."
curl -s -X GET ${BASE_URL}/api/notifications/unread-count \
  -H "Accept: application/json" | jq '.'

echo ""
echo "💌 Step 4: Creating notification via group invite..."

# Send group invite (this should create a notification)
curl -s -X POST ${BASE_URL}/groups/group/invite \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": 1,
    "user_id": 2,
    "status": "invited",
    "prev_status": "none"
  }' | jq '.'

echo ""
echo "🔍 Step 5: Checking notifications after invite..."

# Check notifications again (should have new notification)
echo "Getting notifications for user 2 after invite..."
curl -s -X GET ${BASE_URL}/api/notifications \
  -H "Accept: application/json" | jq '.'

# Check unread count
echo "Getting unread count for user 2 after invite..."
curl -s -X GET ${BASE_URL}/api/notifications/unread-count \
  -H "Accept: application/json" | jq '.'

echo ""
echo "✅ Step 6: Testing mark as read functionality..."

# Mark notification as read
curl -s -X POST ${BASE_URL}/api/notifications/mark-read \
  -H "Content-Type: application/json" \
  -d '{
    "notification_ids": [1]
  }' | jq '.'

echo ""
echo "🏁 Test completed!"
echo "If you see notifications created and status changes, the system is working!"
