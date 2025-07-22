// Test the new group member endpoints
const baseUrl = 'http://localhost:8080/groups/group'

// Test getting available users for group ID 1
fetch(`${baseUrl}/1/available-users`, {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json'
  }
})
.then(async res => {
  const text = await res.text();
  console.log('Available users response:', res.status, text);
  try {
    const json = JSON.parse(text);
    console.log('Available users:', json);
  } catch (e) {
    console.error('Invalid JSON response:', e.message);
  }
})
.catch(err => console.error('Available users request failed:', err));

// Test searching users not in group ID 1
fetch(`${baseUrl}/1/search-users?q=john`, {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json'
  }
})
.then(async res => {
  const text = await res.text();
  console.log('Search users response:', res.status, text);
  try {
    const json = JSON.parse(text);
    console.log('Search results:', json);
  } catch (e) {
    console.error('Invalid JSON response:', e.message);
  }
})
.catch(err => console.error('Search users request failed:', err));

// Test browsing groups with type "not_joined" 
fetch(`${baseUrl}/browse`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    user_id: "1",
    start: -1,
    n_items: 10,
    type: "not_joined"
  })
})
.then(async res => {
  const text = await res.text();
  console.log('Browse not joined groups response:', res.status, text);
  try {
    const json = JSON.parse(text);
    console.log('Not joined groups:', json);
  } catch (e) {
    console.error('Invalid JSON response:', e.message);
  }
})
.catch(err => console.error('Browse not joined groups request failed:', err));
