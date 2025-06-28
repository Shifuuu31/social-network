const baseUrl = "http://localhost:8080/auth";

const email = `signout1750730580607@test.com`;
const password = "Logout123";

(async () => {
  console.log("\n--- SIGNOUT TESTS ---");


  const signin = await fetch(`${baseUrl}/signin`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password })
  });
  const cookie = signin.headers.get("set-cookie")?.split(";")[0];
console.log(cookie);

  // ✅ Valid signout
  const signout = await fetch(`${baseUrl}/signout`, {
    method: "DELETE",
    headers: { Cookie: cookie }
  });
  console.log("✅ Signout with session:", signout.status === 200, signout.status);

  // ❌ No session cookie
  const signoutNoCookie = await fetch(`${baseUrl}/signout`, {
    method: "DELETE"
  });
  console.log("❌ Signout without session:", signoutNoCookie.status >= 400, signoutNoCookie.status);
})();
