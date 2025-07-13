async function testBasicSignup() {
  console.log('🧪 Testing Basic Signup (JSON)...\n');

  try {
    console.log('1. Sending basic signup request...');
    
    const response = await fetch('http://localhost:8080/auth/signup', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: 'testbasic@example.com',
        password: 'TestPass123!',
        first_name: 'Test',
        last_name: 'Basic',
        date_of_birth: '1990-01-01T00:00:00Z',
        gender: 'male',
        nickname: 'testbasic',
        about_me: 'Test user without avatar',
        is_public: true
      })
    });

    if (response.ok) {
      console.log('✅ Basic signup successful!');
    } else {
      const error = await response.text();
      console.log('❌ Basic signup failed!');
      console.log('📊 Error:', error);
      console.log('📊 Status:', response.status);
    }

  } catch (error) {
    console.log('❌ Test failed with error:', error.message);
  }
}

// Run the test
testBasicSignup(); 