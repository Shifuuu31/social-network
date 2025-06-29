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




export async function getPosts(filter = {}) {
  // filter example structure:
  // {
  //   userId: number,        // required - the requesting user's ID
  //   groupId: number | null, // optional - filter by group
  //   limit: number,         // optional - number of posts to return
  //   offset: number,        // optional - for pagination
  // }

  try {
    const response = await fetch(`${API_BASE_URL}/feed`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(filter),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      const errorMsg = errorData.error || 'Failed to fetch posts';
      throw new Error(errorMsg);
    }

    return await response.json();
  } catch (err) {
    console.error('getPosts error:', err);
    throw err;
  }
}