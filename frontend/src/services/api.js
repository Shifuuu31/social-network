// api.js
const API_BASE_URL = '/post';
const HEADERS = {
  JSON: { 'Content-Type': 'application/json' }
};

/**
 * Creates a new post by sending either JSON or form data to the backend.
 * @param {Object} postData - Post data including ownerId, content, imageFile, etc.
 * @returns {Promise<Object>} - Resolves with the server response (e.g., { message, post_id })
 */
export async function createPost(postData) {
  console.log("Creating post with data:", postData);

  const url = `${API_BASE_URL}/new`;
  let response;

  try {
    if (postData.imageFile) {
      // Handle image upload with FormData
      const formData = buildFormData(postData);
      
      response = await fetch(url, {
        method: 'POST',
        body: formData,
      });
    } else {
      // Handle regular JSON payload
      const bodyPayload = {
        // ownerId: postData.ownerId,
      owner_id:1,
        group_id: postData.groupId || undefined,
        content: postData.content,
        privacy: postData.privacy,
        // chosenUsersIds: postData.privacy === 'private' ? postData.chosenUsersIds : [],
      };

      response = await fetch(url, {
        method: 'POST',
        headers: HEADERS.JSON,
        body: JSON.stringify(bodyPayload),
      });
    }

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      const errorMsg = errorData.error || `HTTP error! Status: ${response.status}`;
      throw new Error(errorMsg);
    }

    return await response.json(); // e.g., { message: "Post created", post_id: 123 }

  } catch (err) {
    console.error('Error creating post:', err.message);
    throw err;
  }
}

/**
 * Builds FormData object from postData for file uploads
 * @param {Object} postData
 * @returns {FormData}
 */
function buildFormData(postData) {
  const formData = new FormData();
  formData.append('ownerId', postData.ownerId);
  
  if (postData.groupId !== null && postData.groupId !== undefined) {
    formData.append('groupId', postData.groupId);
  }

  formData.append('content', postData.content);
  formData.append('privacy', postData.privacy);

  // if (postData.privacy === 'private' && Array.isArray(postData.chosenUsersIds)) {
  //   postData.chosenUsersIds.forEach(id => formData.append('chosenUsersIds', id));
  // }

  if (postData.imageFile) {
    formData.append('image', postData.imageFile);
  }

  return formData;
}







///////////////


 /**
 * Fetches posts from the feed based on filter criteria.
 * @param {Object} [filter] - Optional filter object (userId is required)
 * @returns {Promise<Object>} Resolves with list of posts
 */
 


 export async function getPosts(filter = {}) {
  const url = `${API_BASE_URL}/feed`;

  // Normalize input to match Go backend expectations
  const requestFilter = {
    id: filter.id || 0,
    type: filter.type || "public",
    start: filter.start || 0,
    n_post: filter.nPost || 20
  };

  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        // Add any auth headers if needed
        // 'Authorization': `Bearer ${token}`
      },
      mode: 'cors', // Explicitly set CORS mode
      credentials: 'include', // If you need cookies/auth
      body: JSON.stringify(requestFilter)
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      const errorMsg = errorData.error || `HTTP error! Status: ${response.status}`;
      throw new Error(errorMsg);
    }

    const data = await response.json();
    
    // Safely clone the data to avoid XrayWrapper issues
    return JSON.parse(JSON.stringify(data));
    
  } catch (err) {
    console.error('Error fetching posts:', err.message);
    
    // Check if it's a CORS error
    if (err.message.includes('CORS') || err.message.includes('cross-origin')) {
      throw new Error('Network error: Please check if the server is running and CORS is configured');
    }
    
    throw err;
  }
}