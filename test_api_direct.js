async function testPostFeedAPI() {
  console.log('🧪 Testing Post Feed API directly...\n');

  try {
    // Step 1: Sign in to get a session
    console.log('1. Signing in...');
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

    // Step 2: Get post feed and log the raw response
    console.log('\n2. Fetching post feed...');
    
    const feedResponse = await fetch('http://localhost:8080/post/feed', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': cookies
      },
      body: JSON.stringify({
        id: 0,
        type: 'public',
        start: 0,
        n_post: 3
      })
    });

    if (feedResponse.ok) {
      const posts = await feedResponse.json();
      console.log('✅ Post feed fetched successfully!');
      console.log('📊 Raw response:', JSON.stringify(posts, null, 2));
      
      if (posts.length > 0) {
        console.log('\n📋 First post structure:');
        console.log('- All keys:', Object.keys(posts[0]));
        console.log('- Has avatar_url:', 'avatar_url' in posts[0]);
        console.log('- avatar_url value:', posts[0].avatar_url);
      }
    } else {
      const error = await feedResponse.text();
      console.log('❌ Failed to fetch post feed!');
      console.log('📊 Error:', error);
      console.log('📊 Status:', feedResponse.status);
    }

  } catch (error) {
    console.log('❌ Test failed with error:', error.message);
  }
}

// Run the test
testPostFeedAPI(); 