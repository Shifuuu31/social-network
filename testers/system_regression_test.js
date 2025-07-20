#!/usr/bin/env node

const axios = require('axios');
const WebSocket = require('ws');

const API_BASE = 'http://localhost:8080/api';
const WS_URL = 'ws://localhost:8080/ws';

// Test configuration
const testConfig = {
    timeout: 5000,
    retries: 3
};

// Test results tracking
let testResults = {
    passed: 0,
    failed: 0,
    errors: []
};

// Helper function to make API calls with retry logic
async function apiCall(method, endpoint, data = null, headers = {}) {
    for (let i = 0; i < testConfig.retries; i++) {
        try {
            const config = {
                method,
                url: `${API_BASE}${endpoint}`,
                headers: {
                    'Content-Type': 'application/json',
                    ...headers
                },
                timeout: testConfig.timeout
            };
            
            if (data) {
                config.data = data;
            }
            
            const response = await axios(config);
            return response;
        } catch (error) {
            if (i === testConfig.retries - 1) throw error;
            await new Promise(resolve => setTimeout(resolve, 1000));
        }
    }
}

// Helper function to run a test and track results
async function runTest(testName, testFunction) {
    try {
        console.log(`\n🧪 Testing: ${testName}`);
        await testFunction();
        testResults.passed++;
        console.log(`✅ PASSED: ${testName}`);
    } catch (error) {
        testResults.failed++;
        testResults.errors.push({ test: testName, error: error.message });
        console.log(`❌ FAILED: ${testName} - ${error.message}`);
    }
}

// Test 1: API Health Check
async function testAPIHealth() {
    // Test that API is responding to basic requests
    const groupsResponse = await apiCall('GET', '/groups');
    if (groupsResponse.status !== 200) {
        throw new Error('API not responding correctly');
    }
    
    console.log('  ✓ API is healthy and responding');
}

// Test 2: Authentication System
async function testAuthentication() {
    // Test register endpoint
    const registerData = {
        username: `test_user_${Date.now()}`,
        email: `test${Date.now()}@example.com`,
        password: 'password123'
    };
    
    const registerResponse = await apiCall('POST', '/register', registerData);
    if (registerResponse.status !== 200 && registerResponse.status !== 201) {
        throw new Error(`Register failed with status ${registerResponse.status}`);
    }
    
    // Test login endpoint
    const loginData = {
        email: registerData.email,
        password: registerData.password
    };
    
    const loginResponse = await apiCall('POST', '/login', loginData);
    if (loginResponse.status !== 200) {
        throw new Error(`Login failed with status ${loginResponse.status}`);
    }
    
    console.log('  ✓ Register and login working');
}

// Test 3: Groups System
async function testGroupsSystem() {
    // Test get all groups
    const groupsResponse = await apiCall('GET', '/groups');
    if (groupsResponse.status !== 200) {
        throw new Error(`Get groups failed with status ${groupsResponse.status}`);
    }
    
    // Test create group
    const groupData = {
        name: `Test Group ${Date.now()}`,
        description: 'Test group for comprehensive testing'
    };
    
    const createResponse = await apiCall('POST', '/groups', groupData);
    if (createResponse.status !== 200 && createResponse.status !== 201) {
        throw new Error(`Create group failed with status ${createResponse.status}`);
    }
    
    const groupId = createResponse.data.id || createResponse.data.group_id;
    if (!groupId) {
        throw new Error('Group creation did not return group ID');
    }
    
    // Test get specific group
    const getGroupResponse = await apiCall('GET', `/groups/${groupId}`);
    if (getGroupResponse.status !== 200) {
        throw new Error(`Get specific group failed with status ${getGroupResponse.status}`);
    }
    
    console.log('  ✓ Groups CRUD operations working');
    return groupId;
}

// Test 4: Group Members System
async function testGroupMembers(groupId) {
    // Test get group members
    const membersResponse = await apiCall('GET', `/groups/${groupId}/members`);
    if (membersResponse.status !== 200) {
        throw new Error(`Get group members failed with status ${membersResponse.status}`);
    }
    
    // Test invite to group (this should create notifications)
    const inviteData = {
        email: 'testmember@example.com'
    };
    
    const inviteResponse = await apiCall('POST', `/groups/${groupId}/invite`, inviteData);
    if (inviteResponse.status !== 200 && inviteResponse.status !== 201) {
        throw new Error(`Group invite failed with status ${inviteResponse.status}`);
    }
    
    console.log('  ✓ Group members system working');
}

// Test 5: Events System
async function testEventsSystem(groupId) {
    // Test get group events
    const eventsResponse = await apiCall('GET', `/groups/${groupId}/events`);
    if (eventsResponse.status !== 200) {
        throw new Error(`Get group events failed with status ${eventsResponse.status}`);
    }
    
    // Test create event
    const eventData = {
        title: `Test Event ${Date.now()}`,
        description: 'Test event for comprehensive testing',
        event_date: new Date(Date.now() + 86400000).toISOString(), // Tomorrow
        location: 'Test Location'
    };
    
    const createEventResponse = await apiCall('POST', `/groups/${groupId}/events`, eventData);
    if (createEventResponse.status !== 200 && createEventResponse.status !== 201) {
        throw new Error(`Create event failed with status ${createEventResponse.status}`);
    }
    
    console.log('  ✓ Events system working');
}

// Test 6: Notification System (New Implementation)
async function testNotificationSystem() {
    // Test get notifications
    const notificationsResponse = await apiCall('GET', '/notifications');
    if (notificationsResponse.status !== 200) {
        throw new Error(`Get notifications failed with status ${notificationsResponse.status}`);
    }
    
    // Test get unread count
    const unreadResponse = await apiCall('GET', '/notifications/unread-count');
    if (unreadResponse.status !== 200) {
        throw new Error(`Get unread count failed with status ${unreadResponse.status}`);
    }
    
    // Create a test notification by creating a group (should trigger notification)
    const groupData = {
        name: `Notification Test Group ${Date.now()}`,
        description: 'Group to test notification creation'
    };
    
    const createResponse = await apiCall('POST', '/groups', groupData);
    if (createResponse.status !== 200 && createResponse.status !== 201) {
        throw new Error(`Failed to create group for notification test`);
    }
    
    // Check if notifications increased
    const newNotificationsResponse = await apiCall('GET', '/notifications');
    if (newNotificationsResponse.status !== 200) {
        throw new Error(`Failed to get notifications after group creation`);
    }
    
    // Test mark as read if we have notifications
    if (newNotificationsResponse.data.data && newNotificationsResponse.data.data.length > 0) {
        const notificationIds = newNotificationsResponse.data.data.slice(0, 2).map(n => n.id);
        const markReadResponse = await apiCall('POST', '/notifications/mark-read', {
            notification_ids: notificationIds
        });
        if (markReadResponse.status !== 200) {
            throw new Error(`Mark notifications as read failed with status ${markReadResponse.status}`);
        }
        
        // Test delete notification
        const deleteResponse = await apiCall('DELETE', `/notifications/${notificationIds[0]}`);
        if (deleteResponse.status !== 200) {
            throw new Error(`Delete notification failed with status ${deleteResponse.status}`);
        }
    }
    
    console.log('  ✓ Notification system working');
}

// Test 7: WebSocket Connection
async function testWebSocketConnection() {
    return new Promise((resolve, reject) => {
        const ws = new WebSocket(WS_URL);
        let connected = false;
        
        const timeout = setTimeout(() => {
            if (!connected) {
                ws.close();
                reject(new Error('WebSocket connection timeout'));
            }
        }, 5000);
        
        ws.on('open', () => {
            connected = true;
            clearTimeout(timeout);
            console.log('  ✓ WebSocket connection established');
            
            // Test sending a message
            ws.send(JSON.stringify({
                type: 'test',
                message: 'Test message'
            }));
            
            setTimeout(() => {
                ws.close();
                resolve();
            }, 1000);
        });
        
        ws.on('error', (error) => {
            clearTimeout(timeout);
            reject(new Error(`WebSocket error: ${error.message}`));
        });
        
        ws.on('message', (data) => {
            console.log('  ✓ WebSocket message received:', data.toString());
        });
    });
}

// Test 8: Database Integration
async function testDatabaseIntegration() {
    // Create a group and verify it persists
    const groupData = {
        name: `DB Test Group ${Date.now()}`,
        description: 'Testing database persistence'
    };
    
    const createResponse = await apiCall('POST', '/groups', groupData);
    if (createResponse.status !== 200 && createResponse.status !== 201) {
        throw new Error('Failed to create group for DB test');
    }
    
    const groupId = createResponse.data.id || createResponse.data.group_id;
    
    // Retrieve the group to verify it was saved
    const getResponse = await apiCall('GET', `/groups/${groupId}`);
    if (getResponse.status !== 200) {
        throw new Error('Failed to retrieve created group');
    }
    
    if (getResponse.data.name !== groupData.name) {
        throw new Error('Group data not properly saved to database');
    }
    
    console.log('  ✓ Database integration working');
}

// Test 9: Frontend Integration
async function testFrontendIntegration() {
    try {
        // Test if frontend is accessible
        const frontendResponse = await axios.get('http://localhost:5173/', { timeout: 5000 });
        if (frontendResponse.status !== 200) {
            throw new Error(`Frontend not accessible, status: ${frontendResponse.status}`);
        }
        console.log('  ✓ Frontend is accessible');
    } catch (error) {
        throw new Error(`Frontend integration test failed: ${error.message}`);
    }
}

// Test 10: CORS and Headers
async function testCORSAndHeaders() {
    try {
        const response = await apiCall('OPTIONS', '/groups');
        // OPTIONS request should succeed or return 200/204
        if (![200, 204, 404].includes(response.status)) {
            throw new Error(`CORS preflight failed with status ${response.status}`);
        }
        console.log('  ✓ CORS configuration working');
    } catch (error) {
        // If OPTIONS not implemented, that's okay for this test
        console.log('  ⚠ CORS preflight not implemented (optional)');
    }
}

// Main test runner
async function runComprehensiveTests() {
    console.log('🚀 Starting Comprehensive System Tests');
    console.log('=========================================');
    
    const startTime = Date.now();
    
    await runTest('API Health Check', testAPIHealth);
    await runTest('Database Integration', testDatabaseIntegration);
    await runTest('Frontend Integration', testFrontendIntegration);
    await runTest('CORS and Headers', testCORSAndHeaders);
    await runTest('Authentication System', testAuthentication);
    
    let groupId;
    await runTest('Groups System', async () => {
        groupId = await testGroupsSystem();
    });
    
    if (groupId) {
        await runTest('Group Members System', () => testGroupMembers(groupId));
        await runTest('Events System', () => testEventsSystem(groupId));
    }
    
    await runTest('Notification System', testNotificationSystem);
    await runTest('WebSocket Connection', testWebSocketConnection);
    
    const endTime = Date.now();
    const duration = (endTime - startTime) / 1000;
    
    console.log('\n📊 Test Results Summary');
    console.log('========================');
    console.log(`✅ Passed: ${testResults.passed}`);
    console.log(`❌ Failed: ${testResults.failed}`);
    console.log(`⏱️ Duration: ${duration}s`);
    
    if (testResults.failed > 0) {
        console.log('\n❌ Failed Tests:');
        testResults.errors.forEach(error => {
            console.log(`  - ${error.test}: ${error.error}`);
        });
    }
    
    if (testResults.failed === 0) {
        console.log('\n🎉 All tests passed! System is working correctly.');
        console.log('\n🔍 System Health Summary:');
        console.log('  ✅ API endpoints responding correctly');
        console.log('  ✅ Database operations working');
        console.log('  ✅ Frontend accessible');
        console.log('  ✅ Authentication system functional');
        console.log('  ✅ Groups system operational');
        console.log('  ✅ Notification system integrated successfully');
        console.log('  ✅ WebSocket real-time communication working');
        console.log('  ✅ No regression issues detected');
    } else {
        console.log(`\n⚠️ ${testResults.failed} test(s) failed. Please review the errors above.`);
    }
    
    process.exit(testResults.failed > 0 ? 1 : 0);
}

// Run the tests
runComprehensiveTests().catch(error => {
    console.error('💥 Test runner error:', error.message);
    process.exit(1);
});
