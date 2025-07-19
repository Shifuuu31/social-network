// api.js
import { useAuth } from '../composables/useAuth.js';

const API_BASE_URL = '/post';
const HEADERS = {
  JSON: { 'Content-Type': 'application/json' }
};

/**
 * Creates a new post by sending either JSON or form data to the backend.
 * @param {Object|FormData} postData - Post data including ownerId, content, imageFile, etc.
 * @returns {Promise<Object>} - Resolves with the server response (e.g., { message, post_id })
 */
export async function createPost(postData) {
  console.log("Creating post with data:", postData);

  // Fetch current user before creating post
  const { fetchCurrentUser, user } = useAuth();
  const userFetched = await fetchCurrentUser();
  
  if (!userFetched) {
    throw new Error('Failed to fetch current user. Please ensure you are logged in.');
  }

  const url = `${API_BASE_URL}/new`;
  let response;

  try {
    if (postData instanceof FormData) {
      // Handle FormData (for image uploads)
      // The owner_id should already be set in the FormData from the component
      
      response = await fetch(url, {
        method: 'POST',
        body: postData,
        // Don't set Content-Type header - let browser set it with boundary
      });
    } else {
      // Handle regular JSON payload
      const bodyPayload = {
        ownerId: user.value.id, // Use current user's ID
        groupId: postData.groupId || null,
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

// --- Close Friends API ---
const API_BASE = '/api';
const CLOSE_FRIENDS_BASE_URL = `${API_BASE}/users/close-friends`;

/**
 * Add a user to close friends
 * @param {number} friendId
 * @returns {Promise<Object>}
 */
export async function addToCloseFriends(friendId) {
  const response = await fetch(`${CLOSE_FRIENDS_BASE_URL}/add`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ friend_id: friendId })
  });
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || 'Failed to add close friend');
  }
  return response.json();
}

/**
 * Remove a user from close friends
 * @param {number} friendId
 * @returns {Promise<Object>}
 */
export async function removeFromCloseFriends(friendId) {
  const response = await fetch(`${CLOSE_FRIENDS_BASE_URL}/remove`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ friend_id: friendId })
  });
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || 'Failed to remove close friend');
  }
  return response.json();
}

/**
 * Fetch the current user's close friends
 * @returns {Promise<Array>}
 */
export async function fetchCloseFriends() {
  const response = await fetch(`${CLOSE_FRIENDS_BASE_URL}/list`, {
    method: 'GET',
    credentials: 'include',
    headers: { 'Accept': 'application/json' }
  });
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || 'Failed to fetch close friends');
  }
  return response.json();
}