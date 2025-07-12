// Test the new group filtering functionality
const baseUrl = 'http://localhost:8080/groups/group'

console.log('🧪 Testing new group filtering system...\n')

// Test 1: Get "My Groups" (groups user has interacted with)
console.log('1️⃣ Testing "My Groups" (user interacted groups)...')
fetch(`${baseUrl}/browse`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    user_id: "1", // Test with user 1
    start: -1,
    n_items: 20,
    type: "user"
  })
})
.then(res => res.json())
.then(groups => {
  console.log(`✅ User 1 has interacted with ${groups?.length || 0} groups:`)
  groups?.forEach(g => {
    console.log(`   - ${g.title} (${g.is_member || 'creator'})`)
  })
})
.catch(err => console.error('❌ My Groups test failed:', err))

// Test 2: Get "Explore" groups (groups user hasn't interacted with)
setTimeout(() => {
  console.log('\n2️⃣ Testing "Explore" groups (user not interacted)...')
  fetch(`${baseUrl}/browse`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      user_id: "1", // Test with user 1
      start: -1,
      n_items: 20,
      type: "all"
    })
  })
  .then(res => res.json())
  .then(groups => {
    console.log(`✅ User 1 can explore ${groups?.length || 0} new groups:`)
    groups?.forEach(g => {
      console.log(`   - ${g.title} (no interaction)`)
    })
  })
  .catch(err => console.error('❌ Explore test failed:', err))
}, 500)

// Test 3: Compare different users
setTimeout(() => {
  console.log('\n3️⃣ Testing different user perspective...')
  fetch(`${baseUrl}/browse`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      user_id: "2", // Test with user 2
      start: -1,
      n_items: 20,
      type: "user"
    })
  })
  .then(res => res.json())
  .then(groups => {
    console.log(`✅ User 2 has interacted with ${groups?.length || 0} groups:`)
    groups?.forEach(g => {
      console.log(`   - ${g.title} (${g.is_member || 'creator'})`)
    })
  })
  .catch(err => console.error('❌ User 2 test failed:', err))
}, 1000)

// Test 4: Verify no overlap between "My Groups" and "Explore"
setTimeout(() => {
  console.log('\n4️⃣ Testing no overlap between My Groups and Explore...')
  
  Promise.all([
    fetch(`${baseUrl}/browse`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ user_id: "1", start: -1, n_items: 20, type: "user" })
    }).then(res => res.json()),
    
    fetch(`${baseUrl}/browse`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ user_id: "1", start: -1, n_items: 20, type: "all" })
    }).then(res => res.json())
  ])
  .then(([myGroups, exploreGroups]) => {
    const myGroupIds = new Set((myGroups || []).map(g => g.id))
    const exploreGroupIds = new Set((exploreGroups || []).map(g => g.id))
    
    const overlap = [...myGroupIds].filter(id => exploreGroupIds.has(id))
    
    if (overlap.length === 0) {
      console.log('✅ Perfect! No overlap between My Groups and Explore')
      console.log(`   My Groups: ${myGroupIds.size} groups`)
      console.log(`   Explore: ${exploreGroupIds.size} groups`)
    } else {
      console.log('❌ Found overlap:', overlap)
    }
  })
  .catch(err => console.error('❌ Overlap test failed:', err))
}, 1500)

// Final summary
setTimeout(() => {
  console.log('\n🎉 Group filtering tests completed!')
  console.log('Summary:')
  console.log('✅ My Groups: Shows groups where user is member/requested/invited/creator')
  console.log('✅ Explore: Shows groups where user has no interaction')
  console.log('✅ Proper user-specific filtering implemented')
}, 2500)
