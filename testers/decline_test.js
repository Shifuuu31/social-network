// Test the decline functionality to ensure records are deleted
const baseUrl = 'http://localhost:8080/groups/group'

// Test declining a group invite/request
fetch(`${baseUrl}/accept-decline`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    group_id: 1,
    user_id: 3,
    status: "declined",
    prev_status: "invited" // or "requested"
  })
})
.then(async res => {
  const text = await res.text();
  console.log('Decline response:', res.status, text);
  try {
    const json = JSON.parse(text);
    console.log('Declined member:', json);
  } catch (e) {
    console.error('Invalid JSON response:', e.message);
  }
})
.catch(err => console.error('Decline request failed:', err));

// After declining, test that the user appears in available users again
setTimeout(() => {
  fetch(`${baseUrl}/1/available-users`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json'
    }
  })
  .then(async res => {
    const text = await res.text();
    console.log('Available users after decline:', res.status, text);
    try {
      const json = JSON.parse(text);
      console.log('Users available after decline:', json);
      // Check if user ID 3 is now in the list
      const user3 = json.find(user => user.id === 3);
      if (user3) {
        console.log('✅ User 3 is now available for invite again');
      } else {
        console.log('❌ User 3 is not in available users list');
      }
    } catch (e) {
      console.error('Invalid JSON response:', e.message);
    }
  })
  .catch(err => console.error('Available users check failed:', err));
}, 1000);
