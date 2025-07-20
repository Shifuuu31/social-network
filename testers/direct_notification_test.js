// Direct notification system test - bypasses user creation issues
const sqlite3 = require('sqlite3').verbose();

console.log('🧪 Direct Notification System Test');
console.log('===================================');

// Open database
const db = new sqlite3.Database('/home/mbakhcha/mokZwina/backend/pkg/db/data.db', (err) => {
  if (err) {
    console.error('❌ Database connection failed:', err.message);
    process.exit(1);
  }
  console.log('✅ Connected to database');
});

async function testNotificationSystem() {
  try {
    // Insert test users directly into database
    console.log('\n📝 Step 1: Creating test users in database...');
    
    await new Promise((resolve, reject) => {
      db.run(`INSERT OR IGNORE INTO users (id, email, password_hash, first_name, last_name, date_of_birth, image_uuid, nickname, about_me) 
              VALUES (1, 'user1@test.com', 'hashed_password1', 'Test', 'User1', '1990-01-01', 'uuid1', 'testuser1', 'Test user 1')`,
        (err) => {
          if (err) reject(err);
          else resolve();
        });
    });
    
    await new Promise((resolve, reject) => {
      db.run(`INSERT OR IGNORE INTO users (id, email, password_hash, first_name, last_name, date_of_birth, image_uuid, nickname, about_me) 
              VALUES (2, 'user2@test.com', 'hashed_password2', 'Test', 'User2', '1990-01-01', 'uuid2', 'testuser2', 'Test user 2')`,
        (err) => {
          if (err) reject(err);
          else resolve();
        });
    });
    
    console.log('✅ Test users created');

    // Insert test group
    console.log('\n👥 Step 2: Creating test group in database...');
    
    await new Promise((resolve, reject) => {
      db.run(`INSERT OR IGNORE INTO groups (id, title, description, creator_id, image_uuid) 
              VALUES (1, 'Test Notification Group', 'A group for testing notifications', 1, 'group_uuid1')`,
        (err) => {
          if (err) reject(err);
          else resolve();
        });
    });
    
    console.log('✅ Test group created');

    // Insert test notification
    console.log('\n📬 Step 3: Creating test notification in database...');
    
    await new Promise((resolve, reject) => {
      db.run(`INSERT INTO notifications (user_id, type, message, seen) 
              VALUES (2, 'group_invite', 'You have been invited to join Test Notification Group', 0)`,
        (err) => {
          if (err) reject(err);
          else resolve();
        });
    });
    
    console.log('✅ Test notification created');

    // Test API endpoints
    console.log('\n🔍 Step 4: Testing notification API endpoints...');
    
    const fetch = require('node-fetch');
    const baseUrl = 'http://localhost:8080';

    // Test get notifications
    console.log('Testing GET /api/notifications...');
    const notificationsRes = await fetch(`${baseUrl}/api/notifications`);
    const notificationsData = await notificationsRes.json();
    console.log('📋 Notifications response:', JSON.stringify(notificationsData, null, 2));

    // Test unread count
    console.log('Testing GET /api/notifications/unread-count...');
    const unreadRes = await fetch(`${baseUrl}/api/notifications/unread-count`);
    const unreadData = await unreadRes.json();
    console.log('🔢 Unread count response:', JSON.stringify(unreadData, null, 2));

    // Test mark as read
    console.log('Testing POST /api/notifications/mark-read...');
    const markReadRes = await fetch(`${baseUrl}/api/notifications/mark-read`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ notification_ids: [1] })
    });
    const markReadData = await markReadRes.json();
    console.log('✅ Mark as read response:', JSON.stringify(markReadData, null, 2));

    // Test notifications after marking as read
    console.log('Testing notifications after marking as read...');
    const finalNotificationsRes = await fetch(`${baseUrl}/api/notifications`);
    const finalNotificationsData = await finalNotificationsRes.json();
    console.log('📋 Final notifications response:', JSON.stringify(finalNotificationsData, null, 2));

    console.log('\n🎉 Notification system test completed successfully!');
    
  } catch (error) {
    console.error('❌ Test failed:', error.message);
  } finally {
    db.close((err) => {
      if (err) {
        console.error('❌ Database close error:', err.message);
      } else {
        console.log('✅ Database connection closed');
      }
    });
  }
}

testNotificationSystem();
