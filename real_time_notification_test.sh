#!/bin/bash

echo "🔄 TESTING REAL-TIME NOTIFICATION DELIVERY"
echo "========================================="
echo ""

# First, let's create a new notification via the database to simulate a group invite
echo "📝 Creating test notification..."
sqlite3 /home/mbakhcha/mokZwina/backend/pkg/db/data.db <<EOF
INSERT INTO notifications (user_id, type, message, seen) 
VALUES (2, 'group_invite', 'REAL-TIME TEST: You have been invited to join Test Group $(date +%H:%M:%S)', 0);
EOF

NOTIFICATION_ID=$(sqlite3 /home/mbakhcha/mokZwina/backend/pkg/db/data.db "SELECT id FROM notifications ORDER BY id DESC LIMIT 1;")

echo "✅ Created notification with ID: $NOTIFICATION_ID"
echo ""

# Test the API to verify the notification was created
echo "🔍 Verifying via API..."
NEW_COUNT=$(curl -s http://localhost:8080/api/notifications/unread-count | grep -o '"unread_count":[0-9]*' | cut -d: -f2)
echo "✅ Current unread count: $NEW_COUNT"

# Check if we can fetch the new notification
echo ""
echo "📋 Latest notification from API:"
curl -s http://localhost:8080/api/notifications | head -10

echo ""
echo ""
echo "🌐 FRONTEND TESTING INSTRUCTIONS:"
echo "=================================="
echo "1. Open your browser to: http://localhost:5175/"
echo "2. Open browser DevTools (F12) -> Console tab"
echo "3. Look for notification updates in the console"
echo "4. Check if the notification appears in the UI"
echo ""
echo "🧪 Manual WebSocket Test:"
echo "1. Go to: http://localhost:5175/test-notifications.html"
echo "2. Click 'Fetch Notifications' to see the new notification"
echo "3. Click 'Get Unread Count' to verify count updated"
echo ""
echo "✅ Test notification created successfully!"
echo "   The frontend should automatically display this notification"
echo "   if WebSocket is working properly."
