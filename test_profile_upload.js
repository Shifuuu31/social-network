const fs = require('fs');

async function testProfileUpload() {
  console.log('🧪 Testing Profile Image Upload...\n');

  try {
    // Step 1: Create a new user
    console.log('1. Creating a new user...');
    const signupResponse = await fetch('http://localhost:8080/auth/signup', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: 'testprofile@example.com',
        password: 'TestPass123!',
        first_name: 'Test',
        last_name: 'User',
        date_of_birth: '1990-01-01T00:00:00Z',
        gender: 'male',
        nickname: 'testuser',
        about_me: 'Test user for profile upload',
        is_public: true
      })
    });

    if (!signupResponse.ok) {
      console.log('❌ Signup failed:', await signupResponse.text());
      return;
    }

    console.log('✅ User created successfully');

    // Step 2: Sign in to get a session
    console.log('\n2. Signing in...');
    const signinResponse = await fetch('http://localhost:8080/auth/signin', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: 'testprofile@example.com',
        password: 'TestPass123!'
      })
    });

    if (!signinResponse.ok) {
      console.log('❌ Signin failed:', await signinResponse.text());
      return;
    }

    const cookies = signinResponse.headers.get('set-cookie');
    console.log('✅ Signin successful');
    console.log('🍪 Cookies:', cookies);

    // Step 3: Upload profile image
    console.log('\n3. Uploading profile image...');
    
    const formData = new FormData();
    formData.append('image', fs.createReadStream('./backend/uploads/Erdos_head_budapest_fall_1992.jpg'));

    const uploadResponse = await fetch('http://localhost:8080/upload/profile', {
      method: 'POST',
      headers: {
        'Cookie': cookies
      },
      body: formData
    });

    if (uploadResponse.ok) {
      const result = await uploadResponse.json();
      console.log('✅ Profile image upload successful!');
      console.log('📊 Response:', result);
    } else {
      const error = await uploadResponse.text();
      console.log('❌ Profile image upload failed!');
      console.log('📊 Error:', error);
      console.log('📊 Status:', uploadResponse.status);
    }

  } catch (error) {
    console.log('❌ Test failed with error:', error.message);
  }
}

// Run the test
testProfileUpload(); 