const fs = require('fs');

async function testPostFeedWithAvatar() {
  console.log('🧪 Testing Post Feed with Avatar URL...\n');

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

    // Step 2: Get post feed
    console.log('\n2. Fetching post feed...');
    
    const feedResponse = await fetch('http://localhost:8080/post/feed', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': cookies
      },
      body: JSON.stringify({
        id: 0, // Current user ID
        type: 'public',
        start: 0,
        n_post: 5
      })
    });

    if (feedResponse.ok) {
      const posts = await feedResponse.json();
      console.log('✅ Post feed fetched successfully!');
      console.log('📊 Number of posts:', posts.length);
      
      if (posts.length > 0) {
        console.log('\n📋 First post details:');
        console.log('- ID:', posts[0].id);
        console.log('- Owner:', posts[0].owner);
        console.log('- Owner ID:', posts[0].owner_id);
        console.log('- Avatar URL:', posts[0].avatar_url || 'No avatar URL');
        console.log('- Content:', posts[0].content?.substring(0, 50) + '...');
        
        if (posts[0].avatar_url) {
          console.log('✅ Avatar URL is included in post feed!');
        } else {
          console.log('❌ Avatar URL is missing from post feed');
        }
      } else {
        console.log('📝 No posts found in feed');
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
testPostFeedWithAvatar(); 