// Frontend invite functionality test
// This tests the new user filtering and search features

console.log('Testing frontend invite functionality...');

// Test 1: Available users endpoint
fetch('http://localhost:5174/api/groups/group/1/available-users')
  .then(response => response.json())
  .then(data => {
    console.log('✅ Available users for group 1:', data);
    console.log(`Found ${data.length} available users`);
    
    // Expected: Users 0 and 3 (not members of group 1)
    const expectedIds = [0, 3];
    const actualIds = data.map(user => user.id).sort();
    const isCorrect = JSON.stringify(expectedIds.sort()) === JSON.stringify(actualIds);
    
    console.log(isCorrect ? '✅ User filtering is working correctly!' : '❌ User filtering has issues');
    console.log('Expected IDs:', expectedIds, 'Actual IDs:', actualIds);
  })
  .catch(error => console.error('❌ Error testing available users:', error));

// Test 2: Search functionality
setTimeout(() => {
  fetch('http://localhost:5174/api/groups/group/1/search-users?q=test')
    .then(response => response.json())
    .then(data => {
      console.log('✅ Search results for "test":', data);
      console.log(`Found ${data.length} users matching "test"`);
      
      // Expected: Only user 0 (first_name: "test", last_name: "test")
      const hasTestUser = data.some(user => user.first_name === 'test' && user.last_name === 'test');
      
      console.log(hasTestUser ? '✅ Search functionality is working correctly!' : '❌ Search functionality has issues');
    })
    .catch(error => console.error('❌ Error testing search:', error));
}, 1000);

// Test 3: Empty search (should return empty array)
setTimeout(() => {
  fetch('http://localhost:5174/api/groups/group/1/search-users?q=')
    .then(response => response.json())
    .then(data => {
      console.log('✅ Empty search results:', data);
      const isEmpty = data.length === 0;
      console.log(isEmpty ? '✅ Empty search handling is correct!' : '❌ Empty search should return empty array');
    })
    .catch(error => console.error('❌ Error testing empty search:', error));
}, 2000);
