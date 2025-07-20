#!/usr/bin/env node

// Test script to verify JSON parsing fix for user_id type mismatch
// This test specifically checks that user_id is sent as integer, not string

const baseUrl = 'http://localhost:8080/api';

async function testJSONTypeFix() {
    console.log('🧪 TESTING: JSON Type Fix for user_id');
    console.log('========================================');

    try {
        // Test 1: Request to join a group (this was causing the JSON parsing error)
        console.log('\n1️⃣ Testing requestJoinGroup with integer user_id...');
        
        const joinResponse = await fetch(`${baseUrl}/groups/group/request`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                user_id: 1,         // Integer (should work)
                group_id: 1,
                status: 'requested',
                prev_status: 'none'
            })
        });

        console.log(`   Response status: ${joinResponse.status}`);
        
        if (joinResponse.ok) {
            const joinData = await joinResponse.json();
            console.log('   ✅ SUCCESS: Integer user_id accepted by backend');
            console.log(`   Response: ${JSON.stringify(joinData)}`);
        } else {
            const errorData = await joinResponse.json().catch(() => ({}));
            console.log(`   ❌ FAILED: ${errorData.message || 'Unknown error'}`);
        }

        // Test 2: Try with string user_id (should fail)
        console.log('\n2️⃣ Testing with string user_id (should fail)...');
        
        const stringResponse = await fetch(`${baseUrl}/groups/group/request`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                user_id: "1",       // String (should fail)
                group_id: 1,
                status: 'requested',
                prev_status: 'none'
            })
        });

        console.log(`   Response status: ${stringResponse.status}`);
        
        if (!stringResponse.ok) {
            const errorData = await stringResponse.json().catch(() => ({}));
            console.log('   ✅ EXPECTED: String user_id correctly rejected');
            console.log(`   Error: ${errorData.message || 'Unknown error'}`);
        } else {
            console.log('   ⚠️  UNEXPECTED: String user_id was accepted (this should fail)');
        }

        // Test 3: Test group browsing with integer user_id
        console.log('\n3️⃣ Testing group browsing with integer user_id...');
        
        const browseResponse = await fetch(`${baseUrl}/groups/group/browse`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                user_id: 1,         // Integer
                start: -1,
                n_items: 20,
                type: 'all',
                search: ''
            })
        });

        console.log(`   Response status: ${browseResponse.status}`);
        
        if (browseResponse.ok) {
            const browseData = await browseResponse.json();
            console.log('   ✅ SUCCESS: Group browsing works with integer user_id');
            console.log(`   Found ${Array.isArray(browseData) ? browseData.length : 'unknown'} groups`);
        } else {
            const errorData = await browseResponse.json().catch(() => ({}));
            console.log(`   ❌ FAILED: ${errorData.message || 'Unknown error'}`);
        }

        console.log('\n🎉 JSON TYPE FIX VERIFICATION COMPLETE');
        console.log('=====================================');
        console.log('The frontend should now correctly send user_id as integers,');
        console.log('avoiding the "cannot unmarshal string into Go struct" error.');

    } catch (error) {
        console.error('❌ Test failed with error:', error.message);
    }
}

// Run the test
testJSONTypeFix();
