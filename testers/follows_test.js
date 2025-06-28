const baseUrl = 'http://localhost:8080/user/follow'

fetch(`${baseUrl}/follow-unfollow`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({ target_id: 2, action: 'follow' })
})
.then(res => res.text().then(data => console.log(res.status, data)))
.catch(err => console.error('Follow request failed:', err));

// To test unfollow, change action to 'unfollow'

fetch(`${baseUrl}/accept-decline`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({ target_id: 2, action: 'accepted' }) // or 'declined'
})
.then(res => res.text().then(data => console.log(res.status, data)))
.catch(err => console.error('Accept/Decline failed:', err));


