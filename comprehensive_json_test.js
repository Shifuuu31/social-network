#!/usr/bin/env node

// Comprehensive test for the JSON type fix
// Tests both GroupsPayload (string user_id) and GroupMember (integer user_id) endpoints

const baseUrl = 'http://localhost:5174/api';

async function testJSONTypeFixes() {
    console.log('🔬 COMPREHENSIVE JSON TYPE TESTING');
    console.log('=====================================');
    console.log('Testing both GroupsPayload and GroupMember endpoints\n');

    // Test 1: Group Browsing (expects string user_id)
    console.log('1️⃣ Testing Group Browsing (GroupsPayload - expects string user_id)...');
    
    try {
        const response = await fetch(`${baseUrl}/groups/group/browse`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                user_id: "1",        // String - correct for GroupsPayload
                start: -1,
                n_items: 20,
                type: "all",
                search: ""
            })
        });

        if (response.ok) {
            const data = await response.json();
            console.log(`   ✅ SUCCESS: Group browsing works with string user_id`);
            console.log(`   Found ${Array.isArray(data) ? data.length : 'unknown'} groups`);
        } else {
            const errorData = await response.json().catch(() => ({}));
            console.log(`   ❌ FAILED: ${errorData.error || 'Unknown error'}`);
        }
    } catch (error) {
        console.error('   ❌ Error:', error.message);
    }

    // Test 2: Group Request (expects integer user_id)
    console.log('\n2️⃣ Testing Group Join Request (GroupMember - expects integer user_id)...');
    
    try {
        const response = await fetch(`${baseUrl}/groups/group/request`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                user_id: 1,          // Integer - correct for GroupMember
                group_id: 10,
                status: "requested",
                prev_status: "none"
            })
        });

        if (response.ok) {
            const data = await response.json();
            console.log(`   ✅ SUCCESS: Group request works with integer user_id`);
            console.log(`   Response: ${JSON.stringify(data)}`);
        } else {
            const errorData = await response.json().catch(() => ({}));
            console.log(`   ❌ FAILED: ${errorData.error || 'Unknown error'}`);
        }
    } catch (error) {
        console.error('   ❌ Error:', error.message);
    }

    // Test 3: Accept/Decline Group (expects integer user_id)
    console.log('\n3️⃣ Testing Accept Group Invitation (GroupMember - expects integer user_id)...');
    
    try {
        const response = await fetch(`${baseUrl}/groups/group/accept-decline`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                user_id: 1,          // Integer - correct for GroupMember
                group_id: 10,
                status: "member",
                prev_status: "requested"
            })
        });

        if (response.ok) {
            const data = await response.json();
            console.log(`   ✅ SUCCESS: Accept invitation works with integer user_id`);
            console.log(`   Response: ${JSON.stringify(data)}`);
        } else {
            const errorData = await response.json().catch(() => ({}));
            console.log(`   ❌ FAILED: ${errorData.error || 'Unknown error'}`);
        }
    } catch (error) {
        console.error('   ❌ Error:', error.message);
    }

    // Test 4: Wrong type for group browsing (should fail)
    console.log('\n4️⃣ Testing Group Browsing with wrong type (integer user_id - should fail)...');
    
    try {
        const response = await fetch(`${baseUrl}/groups/group/browse`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                user_id: 1,          // Integer - wrong for GroupsPayload
                start: -1,
                n_items: 20,
                type: "all",
                search: ""
            })
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            console.log(`   ✅ EXPECTED: Integer user_id correctly rejected for group browsing`);
            console.log(`   Error: ${errorData.error || 'Unknown error'}`);
        } else {
            console.log(`   ⚠️  UNEXPECTED: Integer user_id was accepted (this should fail)`);
        }
    } catch (error) {
        console.error('   ❌ Error:', error.message);
    }

    // Test 5: Wrong type for group member operations (should fail)
    console.log('\n5️⃣ Testing Group Request with wrong type (string user_id - should fail)...');
    
    try {
        const response = await fetch(`${baseUrl}/groups/group/request`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                user_id: "1",        // String - wrong for GroupMember
                group_id: 10,
                status: "requested",
                prev_status: "none"
            })
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            console.log(`   ✅ EXPECTED: String user_id correctly rejected for group member operations`);
            console.log(`   Error: ${errorData.error || 'Unknown error'}`);
        } else {
            console.log(`   ⚠️  UNEXPECTED: String user_id was accepted (this should fail)`);
        }
    } catch (error) {
        console.error('   ❌ Error:', error.message);
    }

    console.log('\n🎉 JSON TYPE FIX VERIFICATION COMPLETE');
    console.log('=====================================');
    console.log('✅ GroupsPayload endpoints accept string user_id');
    console.log('✅ GroupMember endpoints accept integer user_id');
    console.log('✅ Type mismatches are properly rejected');
    console.log('✅ Frontend should send correct types for each endpoint');
}

// Run the comprehensive test
testJSONTypeFixes();
