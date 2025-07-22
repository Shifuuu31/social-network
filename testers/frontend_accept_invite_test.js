// Frontend Accept Group Invite Test - Multiple Users
// This test simulates frontend behavior for accepting group invitations with different users

const baseUrl = 'http://localhost:8080/groups/group'

console.log('🎭 Testing frontend acceptGroupInvite functionality with different users...\n')

// Simulate the frontend store's acceptGroupInvite function for different users
async function simulateAcceptGroupInvite(groupId, userId) {
  console.log(`🔄 User ${userId} attempting to accept invitation to group ${groupId}...`)
  
  try {
    const response = await fetch(`${baseUrl}/accept-decline`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        group_id: groupId,
        user_id: userId,
        status: 'member',
        prev_status: 'invited'
      })
    })

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}))
      throw new Error(errorData.message || 'Failed to accept invitation')
    }

    const data = await response.json()
    console.log(`✅ User ${userId} successfully accepted invitation to group ${groupId}`)
    console.log(`   Response:`, data)
    
    return data
  } catch (err) {
    console.error(`❌ User ${userId} failed to accept invitation:`, err.message)
    throw err
  }
}

// Simulate the frontend store's fetchGroups function for a specific user
async function simulateFetchUserGroups(userId) {
  console.log(`📋 Fetching groups for user ${userId}...`)
  
  try {
    const response = await fetch(`${baseUrl}/browse`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        user_id: userId.toString(),
        start: -1,
        n_items: 20,
        type: "user"
      })
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const groups = await response.json()
    console.log(`   User ${userId} has ${groups?.length || 0} groups:`)
    groups?.forEach(g => {
      const status = g.is_member || (g.creator_id === userId ? 'creator' : 'none')
      console.log(`   - ${g.title} (${status})`)
    })
    
    return groups
  } catch (err) {
    console.error(`❌ Failed to fetch groups for user ${userId}:`, err.message)
    return []
  }
}

// Create test invitations for multiple users
async function createTestInvitations() {
  console.log('🏗️ Setting up test invitations...\n')
  
  const testCases = [
    { groupId: 4, userId: 3, groupName: "wesh a weldi group" },
    { groupId: 5, userId: 2, groupName: "Gamer's United" },
    { groupId: 1, userId: 4, groupName: "JavaScript Wizards" }
  ]
  
  for (const testCase of testCases) {
    try {
      console.log(`📧 Creating invitation: User ${testCase.userId} → Group ${testCase.groupId} (${testCase.groupName})`)
      
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
        console.log(`   ✅ Invitation created successfully`)
      } else {
        const error = await response.text()
        console.log(`   ⚠️ Invitation result: ${response.status} - ${error}`)
      }
    } catch (error) {
      console.error(`   ❌ Failed to create invitation:`, error.message)
    }
  }
  
  console.log('\n⏱️ Waiting 2 seconds for invitations to process...\n')
  await new Promise(resolve => setTimeout(resolve, 2000))
}

// Test accepting invitations with different users
async function testMultipleUsersAcceptInvites() {
  console.log('🎯 Testing accept invite functionality with multiple users...\n')
  
  const users = [2, 3, 4] // Different user IDs to test with
  
  for (const userId of users) {
    console.log(`\n👤 === TESTING USER ${userId} ===`)
    
    // 1. Check user's current groups and find invitations
    const userGroups = await simulateFetchUserGroups(userId)
    const invitedGroups = userGroups?.filter(g => g.is_member === 'invited') || []
    
    if (invitedGroups.length === 0) {
      console.log(`   ℹ️ User ${userId} has no pending invitations`)
      continue
    }
    
    console.log(`   📨 User ${userId} has ${invitedGroups.length} pending invitation(s)`)
    
    // 2. Accept the first invitation
    const groupToAccept = invitedGroups[0]
    console.log(`   🎯 Attempting to accept invitation to "${groupToAccept.title}"`)
    
    try {
      await simulateAcceptGroupInvite(groupToAccept.id, userId)
      
      // 3. Verify the change
      console.log(`   🔍 Verifying acceptance...`)
      await new Promise(resolve => setTimeout(resolve, 1000)) // Wait for DB update
      
      const updatedGroups = await simulateFetchUserGroups(userId)
      const acceptedGroup = updatedGroups?.find(g => g.id === groupToAccept.id)
      
      if (acceptedGroup && acceptedGroup.is_member === 'member') {
        console.log(`   ✅ SUCCESS: User ${userId} is now a member of "${acceptedGroup.title}"`)
        console.log(`   📊 Group member count: ${acceptedGroup.member_count}`)
      } else {
        console.log(`   ❌ VERIFICATION FAILED: Status not updated properly`)
      }
      
    } catch (error) {
      console.error(`   ❌ FAILED: User ${userId} could not accept invitation:`, error.message)
    }
  }
}

// Test edge cases and error scenarios
async function testEdgeCases() {
  console.log('\n🧪 Testing edge cases and error scenarios...\n')
  
  // Test 1: Try to accept invitation that doesn't exist
  console.log('🔍 Test 1: Accept non-existent invitation')
  try {
    await simulateAcceptGroupInvite(9999, 3) // Non-existent group
  } catch (error) {
    console.log(`   ✅ Expected error caught: ${error.message}`)
  }
  
  // Test 2: Try to accept invitation for wrong user
  console.log('\n🔍 Test 2: Accept invitation for different user')
  try {
    await simulateAcceptGroupInvite(4, 999) // Non-existent user
  } catch (error) {
    console.log(`   ✅ Expected error caught: ${error.message}`)
  }
  
  // Test 3: Try to accept already accepted invitation
  console.log('\n🔍 Test 3: Accept already processed invitation')
  try {
    // First check if user 3 has any member groups
    const groups = await simulateFetchUserGroups(3)
    const memberGroup = groups?.find(g => g.is_member === 'member')
    
    if (memberGroup) {
      await simulateAcceptGroupInvite(memberGroup.id, 3)
    } else {
      console.log(`   ℹ️ No member groups found for user 3 to test with`)
    }
  } catch (error) {
    console.log(`   ✅ Expected error caught: ${error.message}`)
  }
}

// Summary function to show final state
async function showFinalSummary() {
  console.log('\n📊 === FINAL SUMMARY ===')
  
  const users = [1, 2, 3, 4]
  
  for (const userId of users) {
    console.log(`\n👤 User ${userId} final groups:`)
    const groups = await simulateFetchUserGroups(userId)
    
    const memberGroups = groups?.filter(g => g.is_member === 'member') || []
    const invitedGroups = groups?.filter(g => g.is_member === 'invited') || []
    const requestedGroups = groups?.filter(g => g.is_member === 'requested') || []
    const createdGroups = groups?.filter(g => g.creator_id === userId) || []
    
    console.log(`   📈 Member of: ${memberGroups.length} groups`)
    console.log(`   📧 Pending invites: ${invitedGroups.length}`)
    console.log(`   ⏳ Pending requests: ${requestedGroups.length}`)
    console.log(`   👑 Created: ${createdGroups.length} groups`)
  }
}

// Main test runner
async function runFrontendTests() {
  try {
    console.log('🚀 Starting Frontend Accept Group Invite Tests')
    console.log('=' .repeat(60))
    
    // Step 1: Set up test data
    await createTestInvitations()
    
    // Step 2: Test multiple users accepting invites
    await testMultipleUsersAcceptInvites()
    
    // Step 3: Test edge cases
    await testEdgeCases()
    
    // Step 4: Show final summary
    await showFinalSummary()
    
    console.log('\n🎉 Frontend Accept Group Invite tests completed!')
    console.log('=' .repeat(60))
    
  } catch (error) {
    console.error('💥 Test suite failed:', error.message)
  }
}

// Run the tests
runFrontendTests().catch(console.error)
