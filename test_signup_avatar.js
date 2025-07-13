const fs = require('fs');

async function testSignupWithAvatar() {
  console.log('🧪 Testing Signup with Avatar Upload...\n');

  try {
    // Create FormData with user info and avatar
    const FormData = require('form-data');
    const form = new FormData();
    
    // Add user data
    form.append('email', 'testavatar@example.com');
    form.append('password', 'TestPass123!');
    form.append('first_name', 'Test');
    form.append('last_name', 'Avatar');
    form.append('date_of_birth', '1990-01-01T00:00:00Z');
    form.append('gender', 'male');
    form.append('nickname', 'testavatar');
    form.append('about_me', 'Test user with avatar');
    form.append('is_public', 'true');
    
    // Add avatar file
    form.append('avatar_file', fs.createReadStream('./backend/uploads/Erdos_head_budapest_fall_1992.jpg'));

    console.log('1. Sending signup request with avatar...');
    
    const response = await fetch('http://localhost:8080/auth/signup', {
      method: 'POST',
      body: form,
      // Don't set Content-Type - let browser set it with boundary
    });

    if (response.ok) {
      console.log('✅ Signup with avatar successful!');
      
      // Now test signin to verify the avatar was saved
      console.log('\n2. Testing signin to verify avatar...');
      
      const signinResponse = await fetch('http://localhost:8080/auth/signin', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: 'testavatar@example.com',
          password: 'TestPass123!'
        })
      });

      if (signinResponse.ok) {
        console.log('✅ Signin successful!');
        
        // Get user profile to check avatar_url
        const cookies = signinResponse.headers.get('set-cookie');
        
        const profileResponse = await fetch('http://localhost:8080/users/profile/me', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Cookie': cookies
          }
        });

        if (profileResponse.ok) {
          const user = await profileResponse.json();
          console.log('📋 User profile:');
          console.log('- Email:', user.email);
          console.log('- Avatar URL:', user.avatar_url || 'No avatar');
          
          if (user.avatar_url) {
            console.log('✅ Avatar was saved during signup!');
          } else {
            console.log('❌ Avatar was not saved during signup');
          }
        }
      }
    } else {
      const error = await response.text();
      console.log('❌ Signup failed!');
      console.log('📊 Error:', error);
      console.log('📊 Status:', response.status);
    }

  } catch (error) {
    console.log('❌ Test failed with error:', error.message);
  }
}

// Run the test
testSignupWithAvatar(); 