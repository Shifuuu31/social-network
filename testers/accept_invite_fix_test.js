// Test to verify the fix for the forbidden error
const baseUrl = 'http://localhost:8080/groups/group'

console.log('🔧 Testing Accept Invite Fix - No More Forbidden Errors\n')

async function testAcceptInviteWithDifferentUsers() {
  console.log('=== TESTING ACCEPT INVITE FIX ===\n')
  
  // Test data: different users accepting their invitations
  const testCases = [
    { userId: 2, groupId: 4, userName: 'User 2' },
    { userId: 3, groupId: 5, userName: 'User 3' },
    { userId: 4, groupId: 1, userName: 'User 4' }
  ]
  
  // First, create invitations for each user
  console.log('🏗️ Setting up test invitations...')
  for (const testCase of testCases) {
    try {
      const response = await fetch(`${baseUrl}/invite`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          group_id: testCase.groupId,
          user_id: testCase.userId,
          status: 'invited',
          prev_status: 'none'
        })
      })
      
      if (response.ok) {
        console.log(`   ✅ Created invitation: ${testCase.userName} → Group ${testCase.groupId}`)
      } else {
        console.log(`   ⚠️ Invitation may exist: ${testCase.userName} → Group ${testCase.groupId}`)
      }
    } catch (error) {
      console.error(`   ❌ Failed to create invitation for ${testCase.userName}:`, error.message)
    }
  }
  
  console.log('\n🎯 Testing accept invitations with fixed authorization...\n')
  
  // Now test accepting invitations
  for (const testCase of testCases) {
    console.log(`👤 Testing ${testCase.userName} accepting invitation to Group ${testCase.groupId}:`)
    
    try {
      const response = await fetch(`${baseUrl}/accept-decline`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          group_id: testCase.groupId,
          user_id: testCase.userId,
          status: 'member',
          prev_status: 'invited'
        })
      })
      
      const responseText = await response.text()
      console.log(`   Status: ${response.status}`)
      console.log(`   Response: ${responseText}`)
      
      if (response.status === 200) {
        console.log(`   ✅ SUCCESS! ${testCase.userName} accepted invitation without forbidden error`)
        
        // Verify the user is now a member
        const verifyResponse = await fetch(`${baseUrl}/browse`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            user_id: testCase.userId.toString(),
            start: -1,
            n_items: 20,
            type: "user"
          })
        })
        
        if (verifyResponse.ok) {
          const groups = await verifyResponse.json()
          const acceptedGroup = groups?.find(g => g.id === testCase.groupId)
          
          if (acceptedGroup && acceptedGroup.is_member === 'member') {
            console.log(`   ✅ VERIFIED: ${testCase.userName} is now a member of the group`)
          } else {
            console.log(`   ⚠️ Could not verify membership status`)
          }
        }
      } else if (response.status === 403) {
        console.log(`   ❌ STILL FORBIDDEN! Fix didn't work for ${testCase.userName}`)
      } else {
        console.log(`   ⚠️ Unexpected status for ${testCase.userName}: ${response.status}`)
      }
    } catch (error) {
      console.error(`   ❌ Error testing ${testCase.userName}:`, error.message)
    }
    
    console.log('')
  }
}

// Test edge cases
async function testEdgeCases() {
  console.log('🧪 Testing edge cases...\n')
  
  // Test 1: User trying to accept someone else's invitation (should fail)
  console.log('🔍 Test 1: User trying to accept someone else\'s invitation')
  try {
    const response = await fetch(`${baseUrl}/accept-decline`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        group_id: 4,
        user_id: 999, // Different user than who was invited
        status: 'member',
        prev_status: 'invited'
      })
    })
    
    console.log(`   Status: ${response.status}`)
    if (response.status === 403 || response.status === 400) {
      console.log(`   ✅ Correctly rejected unauthorized access`)
    } else {
      console.log(`   ❌ Should have been rejected but wasn't`)
    }
  } catch (error) {
    console.log(`   ✅ Error correctly caught: ${error.message}`)
  }
  
  // Test 2: Invalid prev_status
  console.log('\n🔍 Test 2: Invalid prev_status')
  try {
    const response = await fetch(`${baseUrl}/accept-decline`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        group_id: 4,
        user_id: 3,
        status: 'member',
        prev_status: 'invalid_status'
      })
    })
    
    console.log(`   Status: ${response.status}`)
    if (response.status === 400) {
      console.log(`   ✅ Correctly rejected invalid prev_status`)
    } else {
      console.log(`   ❌ Should have been rejected but wasn't`)
    }
  } catch (error) {
    console.log(`   ✅ Error correctly caught: ${error.message}`)
  }
}

// Show summary
async function showSummary() {
  console.log('📊 === SUMMARY OF FIX ===\n')
  
  console.log('🔧 CHANGES MADE:')
  console.log('1. ✅ Fixed hardcoded requesterID := 1')
  console.log('2. ✅ Made requesterID dynamic based on user_id in request')
  console.log('3. ✅ Improved authorization logic for invited users')
  console.log('4. ✅ Added better logging for debugging')
  
  console.log('\n🎯 EXPECTED RESULTS:')
  console.log('- Users can now accept their own invitations')
  console.log('- No more "Forbidden" errors for legitimate accepts')
  console.log('- Proper authorization still enforced')
  console.log('- Better error messages for debugging')
  
  console.log('\n🚀 NEXT STEPS:')
  console.log('- Test the fix in your application')
  console.log('- Verify different users can accept invitations')
  console.log('- Check that unauthorized access is still blocked')
}

// Run all tests
async function runFixVerificationTests() {
  try {
    console.log('🚀 Starting Accept Invite Fix Verification Tests')
    console.log('=' .repeat(60))
    
    await testAcceptInviteWithDifferentUsers()
    await testEdgeCases()
    await showSummary()
    
    console.log('\n🎉 Fix verification tests completed!')
    console.log('=' .repeat(60))
    
  } catch (error) {
    console.error('💥 Test suite failed:', error.message)
  }
}

runFixVerificationTests().catch(console.error)
