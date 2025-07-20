#!/bin/bash

echo "🏗️ GROUP CREATION AND MEMBERSHIP TEST"
echo "===================================="
echo "Date: $(date)"
echo ""

BASE_URL="http://localhost:8080"

echo "📝 1. Creating a new group..."
GROUP_RESPONSE=$(curl -s -X POST "$BASE_URL/groups/group/new" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Auto-Membership Group",
    "description": "Testing automatic creator membership"
  }')

echo "✅ Group creation response:"
echo "$GROUP_RESPONSE" | python3 -c "import sys, json; data=json.load(sys.stdin); print(f'Group ID: {data.get(\"id\", \"N/A\")}'); print(f'Title: {data.get(\"title\", \"N/A\")}'); print(f'Creator ID: {data.get(\"creator_id\", \"N/A\")}')" 2>/dev/null || echo "$GROUP_RESPONSE"

# Extract group ID for next test
GROUP_ID=$(echo "$GROUP_RESPONSE" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('id', ''))" 2>/dev/null)

if [ -n "$GROUP_ID" ] && [ "$GROUP_ID" != "" ]; then
    echo ""
    echo "📋 2. Checking group details and membership..."
    
    GROUP_DETAILS=$(curl -s "$BASE_URL/groups/group/$GROUP_ID")
    echo "✅ Group details response:"
    echo "$GROUP_DETAILS" | python3 -c "import sys, json; data=json.load(sys.stdin); print(f'Group ID: {data.get(\"id\", \"N/A\")}'); print(f'Title: {data.get(\"title\", \"N/A\")}'); print(f'Creator ID: {data.get(\"creator_id\", \"N/A\")}'); print(f'Is Member: {data.get(\"is_member\", \"N/A\")}')" 2>/dev/null || echo "$GROUP_DETAILS"
    
    echo ""
    echo "📊 3. Checking database for membership record..."
    sqlite3 /home/mbakhcha/mokZwina/backend/pkg/db/data.db <<EOF
SELECT 'Group Members for Group $GROUP_ID:' as info;
SELECT id, group_id, user_id, status, created_at 
FROM group_members 
WHERE group_id = $GROUP_ID;
EOF

else
    echo "❌ Failed to extract group ID from response"
fi

echo ""
echo "📊 4. Checking all group members in database..."
sqlite3 /home/mbakhcha/mokZwina/backend/pkg/db/data.db <<EOF
SELECT 'All Group Members:' as info;
SELECT gm.id, gm.group_id, gm.user_id, gm.status, g.title as group_title, gm.created_at
FROM group_members gm
LEFT JOIN groups g ON gm.group_id = g.id
ORDER BY gm.created_at DESC
LIMIT 10;
EOF

echo ""
echo "🔍 5. Testing middleware user ID consistency..."
echo "Current middleware returns user ID: 1 (from GetRequesterID fallback)"
echo "All group operations should use this consistent ID"

echo ""
echo "🎯 DIAGNOSIS:"
echo "- Group creation should automatically add creator as member"
echo "- GetGroup should show is_member status correctly"
echo "- All user IDs should be consistent across operations"
echo ""
echo "✅ Test complete!"
