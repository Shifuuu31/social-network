// Diagnostic test to find the exact cause of the forbidden error
const baseUrl = 'http://localhost:8080/groups/group'

console.log('🔍 DIAGNOSTIC TEST - Finding Forbidden Error Cause\n')

async function debugAcceptInvite() {
  console.log('=== DEBUGGING FORBIDDEN ERROR ===\n')
  
  const testUser = 3
  const testGroup = 5
  
  // Step 1: Check current user groups and status
  console.log(`📋 Step 1: Checking current status for User ${testUser}...`)
  try {
    const response = await fetch(`${baseUrl}/browse`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        user_id: testUser.toString(),
        start: -1,
        n_items: 20,
        type: "user"
      })
    })
    
    if (response.ok) {
      const groups = await response.json()
      console.log(`   User ${testUser} has ${groups?.length || 0} groups:`)
      groups?.forEach(g => {
        const status = g.is_member || (g.creator_id === testUser ? 'creator' : 'none')
        console.log(`   - Group ${g.id}: "${g.title}" (${status})`)
        if (g.id === testGroup) {
          console.log(`     🎯 TARGET GROUP FOUND! Status: ${status}`)
        }
      })
      
      const targetGroup = groups?.find(g => g.id === testGroup)
      if (!targetGroup) {
        console.log(`   ⚠️ User ${testUser} has no record for Group ${testGroup}`)
      } else if (targetGroup.is_member !== 'invited') {
        console.log(`   ⚠️ User ${testUser} is not invited to Group ${testGroup} (status: ${targetGroup.is_member})`)
      }
    } else {
      console.log(`   ❌ Failed to fetch user groups: ${response.status}`)
    }
  } catch (error) {
    console.error(`   ❌ Error fetching groups:`, error.message)
  }
  
  // Step 2: Create a fresh invitation
  console.log(`\n📧 Step 2: Creating fresh invitation for User ${testUser} to Group ${testGroup}...`)
  try {
    const response = await fetch(`${baseUrl}/invite`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        group_id: testGroup,
        user_id: testUser,
        status: 'invited',
        prev_status: 'none'
      })
    })
    
    const responseText = await response.text()
    console.log(`   Status: ${response.status}`)
    console.log(`   Response: ${responseText}`)
    
    if (response.status === 201) {
      console.log(`   ✅ Fresh invitation created successfully`)
    } else if (response.status === 403) {
      console.log(`   ❌ Forbidden to create invitation - check invite authorization`)
    }
  } catch (error) {
    console.error(`   ❌ Error creating invitation:`, error.message)
  }
  
  // Step 3: Try different payload variations
  console.log(`\n🧪 Step 3: Testing different payload variations...`)
  
  const payloadVariations = [
    {
      name: 'Standard payload',
      payload: {
        group_id: testGroup,
        user_id: testUser,
        status: 'member',
        prev_status: 'invited'
      }
    },
    {
      name: 'Lowercase status',
      payload: {
        group_id: testGroup,
        user_id: testUser,
        status: 'member',
        prev_status: 'invited'
      }
    },
    {
      name: 'Accept action',
      payload: {
        group_id: testGroup,
        user_id: testUser,
        status: 'accepted',
        prev_status: 'invited'
      }
    }
  ]
  
  for (const variation of payloadVariations) {
    console.log(`\n   🔍 Testing: ${variation.name}`)
    console.log(`   Payload:`, JSON.stringify(variation.payload, null, 2))
    
    try {
      const response = await fetch(`${baseUrl}/accept-decline`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(variation.payload)
      })
      
      const responseText = await response.text()
      console.log(`   Status: ${response.status}`)
      console.log(`   Response: ${responseText}`)
      
      if (response.status === 200) {
        console.log(`   ✅ SUCCESS with ${variation.name}!`)
        break
      } else if (response.status === 403) {
        console.log(`   ❌ Still forbidden with ${variation.name}`)
      }
    } catch (error) {
      console.error(`   ❌ Error with ${variation.name}:`, error.message)
    }
  }
  
  // Step 4: Check backend logs and analyze
  console.log(`\n📊 Step 4: Analysis and Next Steps`)
  console.log(`   If all variations failed with 403:`)
  console.log(`   1. Check if the authorization fix was properly applied`)
  console.log(`   2. Verify the user actually has 'invited' status in database`)
  console.log(`   3. Check if there are other authorization checks`)
  console.log(`   4. Restart the backend server to ensure changes are loaded`)
  
  console.log(`\n   Backend file to check: backend/pkg/handlers/groups&members&events.go`)
  console.log(`   Look for the AcceptDeclineGroup function around line 530-580`)
}

// Test simple user 1 case (should work with hardcoded logic)
async function testHardcodedCase() {
  console.log('\n🔧 Testing with User 1 (hardcoded case)...')
  
  // Create invitation for user 1
  try {
    const inviteResponse = await fetch(`${baseUrl}/invite`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        group_id: 6,
        user_id: 1,
        status: 'invited',
        prev_status: 'none'
      })
    })
    console.log(`   Invite status: ${inviteResponse.status}`)
    
    // Try to accept as user 1
    const acceptResponse = await fetch(`${baseUrl}/accept-decline`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        group_id: 6,
        user_id: 1,
        status: 'member',
        prev_status: 'invited'
      })
    })
    
    const responseText = await acceptResponse.text()
    console.log(`   Accept status: ${acceptResponse.status}`)
    console.log(`   Response: ${responseText}`)
    
    if (acceptResponse.status === 200) {
      console.log(`   ✅ User 1 CAN accept - fix is working for invited case!`)
      console.log(`   The issue might be that other users don't have proper invitations`)
    } else {
      console.log(`   ❌ Even User 1 can't accept - there's a deeper issue`)
    }
  } catch (error) {
    console.error(`   ❌ Error testing User 1:`, error.message)
  }
}

// Main diagnostic function
async function runDiagnostics() {
  try {
    console.log('🚀 Starting Forbidden Error Diagnostics')
    console.log('=' .repeat(60))
    
    await debugAcceptInvite()
    await testHardcodedCase()
    
    console.log('\n💡 RECOMMENDATIONS:')
    console.log('1. Check if backend server was restarted after the fix')
    console.log('2. Verify users actually have "invited" status in database')
    console.log('3. Check if there are other authorization layers')
    console.log('4. Look at backend logs for detailed error messages')
    
    console.log('\n🎉 Diagnostic tests completed!')
    console.log('=' .repeat(60))
    
  } catch (error) {
    console.error('💥 Diagnostic failed:', error.message)
  }
}

runDiagnostics().catch(console.error)
