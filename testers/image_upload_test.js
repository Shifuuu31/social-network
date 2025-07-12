// Image Upload Test
// This test verifies that image upload functionality works correctly

const fs = require('fs');
const path = require('path');

// Test data
const testImagePath = './test-image.jpg'; // You'll need to create this test image
const testPostData = {
  content: 'Test post with image upload',
  privacy: 'public',
  ownerId: 1
};

async function testImageUpload() {
  console.log('🧪 Testing Image Upload Functionality...\n');

  try {
    // Check if test image exists
    if (!fs.existsSync(testImagePath)) {
      console.log('⚠️  Test image not found. Creating a dummy test...');
      console.log('📝 To test with real image, create a test-image.jpg file');
      return;
    }

    // Create FormData
    const FormData = require('form-data');
    const form = new FormData();
    
    form.append('content', testPostData.content);
    form.append('privacy', testPostData.privacy);
    form.append('ownerId', testPostData.ownerId.toString());
    form.append('image', fs.createReadStream(testImagePath));

    // Make request to backend
    const response = await fetch('http://localhost:8080/post/new', {
      method: 'POST',
      body: form,
      // Don't set Content-Type - let it be set automatically with boundary
    });

    if (response.ok) {
      const result = await response.json();
      console.log('✅ Image upload successful!');
      console.log('📊 Response:', result);
      console.log('🖼️  Image path saved:', result.image_url || 'No image path returned');
    } else {
      const error = await response.text();
      console.log('❌ Image upload failed!');
      console.log('📊 Error:', error);
    }

  } catch (error) {
    console.log('❌ Test failed with error:', error.message);
  }
}

// Test image serving
async function testImageServing() {
  console.log('\n🧪 Testing Image Serving...\n');

  try {
    // Test serving an image (you'll need to know a valid image path)
    const testImagePath = 'uploads/test-image.jpg'; // Adjust this path
    
    const response = await fetch(`http://localhost:8080/images/serve/${testImagePath}`, {
      method: 'GET'
    });

    if (response.ok) {
      console.log('✅ Image serving successful!');
      console.log('📊 Content-Type:', response.headers.get('content-type'));
      console.log('📊 Content-Length:', response.headers.get('content-length'));
    } else {
      console.log('❌ Image serving failed!');
      console.log('📊 Status:', response.status);
    }

  } catch (error) {
    console.log('❌ Test failed with error:', error.message);
  }
}

// Run tests
async function runTests() {
  console.log('🚀 Starting Image Upload Tests...\n');
  
  await testImageUpload();
  await testImageServing();
  
  console.log('\n✨ Image upload tests completed!');
}

// Run if this file is executed directly
if (require.main === module) {
  runTests().catch(console.error);
}

module.exports = { testImageUpload, testImageServing, runTests }; 