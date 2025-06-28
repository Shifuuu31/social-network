
const baseUrl = "http://localhost:8080/auth/signup";

(async () => {
  const   nickname = `test${Date.now()}`
  // const email = `${nickname}@example.com`;

  // const password = "ValidPass123";
  const email = `signout1750730580607@test.com`;
const password = "Logout12@3";

  console.log("\n--- SIGNUP TESTS ---");

  // ✅ Success case
  const success = await fetch(baseUrl, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      email,
      password,
      first_name: "Test",
      last_name: "User",
      nickname,
      date_of_birth: new Date("2000-01-01").toISOString()

    })
  });
  console.log("✅ Signup valid:", success.status === 201, success.status);

  // ❌ Missing email
  const missingEmail = await fetch(baseUrl, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      password,
      first_name: "Test",
      last_name: "User",
      nickname: "tester",
      date_of_birth: new Date("2000-01-01").toISOString()

    })
  });
  console.log("❌ Signup missing email:", missingEmail.status === 400);

  // ❌ Weak password
  const weakPassword = await fetch(baseUrl, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      email: `weak${Date.now()}@test.com`,
      password: "123",
      first_name: "Weak",
      last_name: "Pass",
      nickname: "weakie",
      date_of_birth: new Date("2000-01-01").toISOString()

    })
  });
  console.log("❌ Signup weak password:", weakPassword.status === 400);
})();
