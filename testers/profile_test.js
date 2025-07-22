const baseUrl = 'http://localhost:8080/user/profile'
const email = `signout1750730580607@test.com`;
const password = "Logout123";

const signin = await fetch(`${baseUrl}/signin`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password })
  });
  const cookie = signin.headers.get("set-cookie")?.split(";")[0];
// console.log(cookie);

fetch(`${baseUrl}/info`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    Cookie: cookie ,
  },
  body: JSON.stringify({ id: 1 })
})
.then(async res => {
  const text = await res.text();
  // console.log(res.status, text);
  try {
    const json = JSON.parse(text);
    // console.log('Parsed JSON:', json);
  } catch (e) {
    console.error('Invalid JSON response:', e.message);
  }
})
.catch(err => console.error('Request failed:', err));


// fetch(`${baseUrl}/activity`, {
//   method: 'POST',
//   headers: {
//     'Content-Type': 'application/json'
//   },
//   body: JSON.stringify({ id: 1 })
// })
// .then(res => res.text().then(data => console.log(res.status, data)))
// .catch(err => console.error('Request failed:', err));

// fetch(`${baseUrl}/connections`, {
//   method: 'POST',
//   headers: {
//     'Content-Type': 'application/json'
//   },
//   body: JSON.stringify({ id: 1 })
// })
// .then(res => res.json().then(data => console.log(res.status, data)))
// .catch(err => console.error('Request failed:', err));

// fetch(`${baseUrl}/visibility`, {
//   method: 'POST',
//   headers: {
//     // 'Authorization': 'Bearer <your-token>' // add if required
//   }
// })
// .then(res => res.json().then(data => console.log(res.status, data)))
// .catch(err => console.error('Request failed:', err));
