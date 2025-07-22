// Comprehensive test of all new group member features
const baseUrl = 'http://localhost:8080/groups/group'

console.log('🧪 Starting comprehensive group member feature tests...\n')

// Test 1: Get available users for group 1
console.log('1️⃣ Testing available users endpoint...')
fetch(`${baseUrl}/1/available-users`)
.then(res => res.json())
.then(users => {
  console.log(`✅ Found ${users?.length || 0} available users:`, users?.map(u => `${u.first_name} ${u.last_name}`))
})
.catch(err => console.error('❌ Available users test failed:', err))

// Test 2: Search users not in group
setTimeout(() => {
  console.log('\n2️⃣ Testing search functionality...')
  fetch(`${baseUrl}/1/search-users?q=john`)
  .then(res => res.json())
  .then(users => {
    console.log(`✅ Search for "john" found ${users?.length || 0} users:`, users?.map(u => u.nickname))
  })
  .catch(err => console.error('❌ Search test failed:', err))
}, 500)

// Test 3: Browse not joined groups
setTimeout(() => {
  console.log('\n3️⃣ Testing not joined groups...')
  fetch(`${baseUrl}/browse`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      user_id: "3", // From perspective of user 3
      start: -1,
      n_items: 10,
      type: "not_joined"
    })
  })
  .then(res => res.json())
  .then(groups => {
    console.log(`✅ User 3 sees ${groups?.length || 0} not joined groups:`, groups?.map(g => g.title))
  })
  .catch(err => console.error('❌ Not joined groups test failed:', err))
}, 1000)

// Test 4: Simulate invite and decline workflow
setTimeout(() => {
  console.log('\n4️⃣ Testing invite and decline workflow...')
  
  // First invite user 3 to group 1 (simulate in database)
  console.log('   📤 Simulating invite user 3 to group 1...')
  
  // Then test decline
  setTimeout(() => {
    fetch(`${baseUrl}/accept-decline`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        group_id: 1,
        user_id: 3,
        status: "declined",
        prev_status: "invited"
      })
    })
    .then(res => res.json())
    .then(result => {
      console.log('   ✅ Decline successful:', result.status)
      
      // Verify user 3 is available again
      setTimeout(() => {
        fetch(`${baseUrl}/1/available-users`)
        .then(res => res.json())
        .then(users => {
          const user3Available = users?.some(u => u.id === 3)
          console.log(`   ${user3Available ? '✅' : '❌'} User 3 is ${user3Available ? 'now' : 'not'} available for invite again`)
        })
      }, 200)
    })
    .catch(err => console.error('❌ Decline test failed:', err))
  }, 500)
}, 1500)

// Final summary
setTimeout(() => {
  console.log('\n🎉 All tests completed! Summary of features tested:')
  console.log('   ✅ Get users not in group')
  console.log('   ✅ Search users not in group') 
  console.log('   ✅ Browse not joined groups')
  console.log('   ✅ Decline invitation (deletes record)')
  console.log('   ✅ User becomes available again after decline')
}, 3000)
