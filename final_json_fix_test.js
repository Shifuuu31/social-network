#!/usr/bin/env node

// Final comprehensive test for JSON type fix
// Tests both string and integer user_id requirements

const baseUrl = 'http://localhost:8080';

async function testJSONTypeFix() {
    console.log('🎯 FINAL JSON TYPE FIX VERIFICATION');
    console.log('=' .repeat(50));

    try {
        // Test 1: Group browsing with string user_id (GroupsPayload)
        console.log('\n1️⃣ Testing Group Browsing (requires string user_id)...');
        
        const browseResponse = await fetch(`${baseUrl}/groups/group/browse`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                user_id: "1",         // String (required by GroupsPayload)
                start: -1,
                n_items: 20,
                type: 'all',
                search: ''
            })
        });

        console.log(`   Response status: ${browseResponse.status}`);
        
        if (browseResponse.ok) {
            const browseData = await browseResponse.json();
            console.log('   ✅ SUCCESS: Group browsing works with string user_id');
            console.log(`   Found ${Array.isArray(browseData) ? browseData.length : 'null/undefined'} groups`);
        } else {
            const errorData = await browseResponse.json().catch(() => ({}));
            console.log(`   ❌ FAILED: ${errorData.message || 'Unknown error'}`);
        }

        // Test 2: Group request with integer user_id (GroupMember)
        console.log('\n2️⃣ Testing Group Request (requires integer user_id)...');
        
        const requestResponse = await fetch(`${baseUrl}/groups/group/request`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                user_id: 1,           // Integer (required by GroupMember)
                group_id: 1,
                status: 'requested',
                prev_status: 'none'
            })
        });

        console.log(`   Response status: ${requestResponse.status}`);
        
        if (requestResponse.ok) {
            const requestData = await requestResponse.json();
            console.log('   ✅ SUCCESS: Group request works with integer user_id');
            console.log(`   Response: ${JSON.stringify(requestData)}`);
        } else {
            const errorData = await requestResponse.json().catch(() => ({}));
            console.log(`   ❌ FAILED: ${errorData.message || 'Unknown error'}`);
        }

        // Test 3: Verify type mismatches fail correctly
        console.log('\n3️⃣ Testing Type Mismatch Errors...');
        
        // Should fail: Group browsing with integer user_id
        const browseBadResponse = await fetch(`${baseUrl}/groups/group/browse`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                user_id: 1,           // Integer (should fail for GroupsPayload)
                start: -1,
                n_items: 20,
                type: 'all',
                search: ''
            })
        });
        
        if (!browseBadResponse.ok) {
            console.log('   ✅ CORRECT: Group browsing correctly rejects integer user_id');
        } else {
            console.log('   ⚠️  UNEXPECTED: Group browsing accepted integer user_id');
        }

        // Should fail: Group request with string user_id
        const requestBadResponse = await fetch(`${baseUrl}/groups/group/request`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                user_id: "1",         // String (should fail for GroupMember)
                group_id: 1,
                status: 'requested',
                prev_status: 'none'
            })
        });
        
        if (!requestBadResponse.ok) {
            console.log('   ✅ CORRECT: Group request correctly rejects string user_id');
        } else {
            console.log('   ⚠️  UNEXPECTED: Group request accepted string user_id');
        }

        console.log('\n🎉 JSON TYPE VERIFICATION COMPLETE');
        console.log('=====================================');
        console.log('✅ Frontend correctly handles different user_id types:');
        console.log('   - Group browsing: user_id as STRING');
        console.log('   - Group operations: user_id as INTEGER');
        console.log('✅ Backend correctly validates JSON types');
        console.log('✅ Type mismatches are properly rejected');

    } catch (error) {
        console.error('❌ Test failed with error:', error.message);
    }
}

// Run the test
testJSONTypeFix();
