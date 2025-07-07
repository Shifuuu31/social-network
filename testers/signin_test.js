const baseUrl = "http://localhost:8080/auth/signin";

// const email = `signin1750728753762@test.com`;
// const password = "GoodPassword123";

const email = `signout1750730580607@test.com`;
const password = "Logout123";
(async () => {
  console.log("\n--- SIGNIN TESTS ---");

    // ❌ Missing fields
  const missingFields = await fetch(baseUrl, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({})
  });
  console.log("❌ Signin missing fields:", missingFields.status >= 400);

  // ❌ Wrong password
  const badPassword = await fetch(baseUrl, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password: "WrongPass" })
  });
  console.log("❌ Signin wrong password:", badPassword.status === 401, badPassword.status);


  //  ✅ Correct credentials
  const goodSignin = await fetch(baseUrl, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password })
  });
  const cookie = goodSignin.headers.get("set-cookie");
  console.log("✅ Signin valid:", goodSignin.status === 200);
  console.log("🍪 Cookie present:", cookie?.includes("session_token"));
})();
