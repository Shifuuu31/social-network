// chatService.js
// Functional, module-based WebSocket chat service (no class, no this)

const API_BASE = '/api';
const CHAT_BASE_URL = `${API_BASE}/chat`;

// Module-level state
let socket = null;
let isConnected = false;
const messageHandlers = new Set();
const connectionHandlers = new Set();

// Get the current host and port for WebSocket connection
const getWebSocketUrl = () => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  // If running on Vite dev server, use backend port
  const host = window.location.hostname + ':8080';

  // Get session token from cookies
  const getCookie = (name) => {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
  };
  
  const sessionToken = getCookie('session_token');
  const tokenParam = sessionToken ? `?token=${encodeURIComponent(sessionToken)}` : '';

  return `${protocol}//${host}/chat/ws${tokenParam}`;
};

// Send authentication message to identify the user
const sendAuthMessage = (userId) => {
  if (socket && isConnected) {
    const authMessage = {
      type: 'auth',
      user_id: userId
    };
    socket.send(JSON.stringify(authMessage));
  }
};

// Notify message handlers
const notifyMessageHandlers = (message) => {
  messageHandlers.forEach(handler => {
    try {
      handler(message);
    } catch (error) {
      console.error('Error in message handler:', error);
    }
  });
};

// Notify connection handlers
const notifyConnectionHandlers = (connected) => {
  connectionHandlers.forEach(handler => {
    try {
      handler(connected);
    } catch (error) {
      console.error('Error in connection handler:', error);
    }
  });
};

// Connect to WebSocket
const connect = (userId) => {
  if (socket && isConnected) {
    return Promise.resolve();
  }
  return new Promise((resolve, reject) => {
    try {
      const wsUrl = getWebSocketUrl();
      
      // Create WebSocket with credentials
      socket = new WebSocket(wsUrl);
      
      socket.onopen = () => {
        console.log("WebSocket connected successfully");
        isConnected = true;
        notifyConnectionHandlers(true);
        // No need to send auth message since backend handles auth via session
        resolve();
      };
      
      socket.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          notifyMessageHandlers(message);
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };
      
      socket.onclose = (event) => {
        console.log("WebSocket connection closed:", event.code, event.reason);
        isConnected = false;
        notifyConnectionHandlers(false);
      };
      
      socket.onerror = (error) => {
        console.error('WebSocket connection error:', error);
        isConnected = false;
        notifyConnectionHandlers(false);
        reject(error);
      };
    } catch (error) {
      console.error('Error creating WebSocket connection:', error);
      reject(error);
    }
  });
};

// Disconnect from WebSocket
const disconnect = () => {
  if (socket) {
    socket.close(1000, 'User initiated disconnect');
    socket = null;
    isConnected = false;
  }
};

// Send a message via WebSocket
const sendMessage = (receiverId, content) => {
  if (!socket || !isConnected) {
    throw new Error('WebSocket not connected');
  }
  const message = {
    type: 'message',
    receiver_id: receiverId,
    content: content
  };
  socket.send(JSON.stringify(message));
};

// Send message via HTTP API (fallback)
const sendMessageAPI = async (receiverId, content) => {
  const payload = {
    receiver_id: receiverId,
    content: content
  };
  try {
    const response = await fetch(`${CHAT_BASE_URL}/send`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(payload)
    });
    const responseText = await response.text();
    if (!response.ok) {
      let errorData;
      try {
        errorData = JSON.parse(responseText);
      } catch (e) {
        errorData = { error: responseText };
      }
      throw new Error(errorData.error || 'Failed to send message');
    }
    let responseData;
    try {
      responseData = JSON.parse(responseText);
    } catch (e) {
      throw new Error('Invalid JSON response from server');
    }
    return responseData;
  } catch (error) {
    throw error;
  }
};

// Get conversation between two users
const getConversation = async (otherUserId, limit = 50, offset = 0) => {
  const payload = {
    other_user_id: otherUserId,
    limit: limit,
    offset: offset
  };
  try {
    const response = await fetch(`${CHAT_BASE_URL}/conversation`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(payload)
    });
    const responseText = await response.text();
    if (!response.ok) {
      let errorData;
      try {
        errorData = JSON.parse(responseText);
      } catch (e) {
        errorData = { error: responseText };
      }
      throw new Error(errorData.error || 'Failed to get conversation');
    }
    let responseData;
    try {
      responseData = JSON.parse(responseText);
    } catch (e) {
      throw new Error('Invalid JSON response from server');
    }
    return responseData;
  } catch (error) {
    throw error;
  }
};

// Get recent conversations
const getRecentConversations = async (limit = 20) => {
  console.log('getRecentConversations', limit);
  const url = `${CHAT_BASE_URL}/recent?limit=${limit}`;
  try {
    const response = await fetch(url, {
      method: 'GET',
      credentials: 'include',
      headers: { 'Accept': 'application/json' }
    });
    if (response.status === 401) {
      throw new Error('Please sign in to see your recent chats.');
    }
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || 'Failed to get recent conversations');
    }
    const responseData = await response.json();
    return responseData;
  } catch (error) {
    throw error;
  }
};

// Delete a message
const deleteMessage = async (messageId) => {
  const url = `${CHAT_BASE_URL}/delete?message_id=${messageId}`;
  try {
    const response = await fetch(url, {
      method: 'DELETE',
      credentials: 'include'
    });
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || 'Failed to delete message');
    }
    const responseData = await response.json();
    return responseData;
  } catch (error) {
    throw error;
  }
};

// Event handling
const onMessage = (handler) => {
  messageHandlers.add(handler);
  return () => messageHandlers.delete(handler); // Return unsubscribe function
};

const onConnectionChange = (handler) => {
  connectionHandlers.add(handler);
  return () => connectionHandlers.delete(handler); // Return unsubscribe function
};

// Get connection status
const getConnectionStatus = () => {
  return {
    isConnected,
    socket: socket ? 'connected' : 'disconnected'
  };
};

// Export all as a functional module
export default {
  connect,
  disconnect,
  sendMessage,
  sendMessageAPI,
  getConversation,
  getRecentConversations,
  deleteMessage,
  onMessage,
  onConnectionChange,
  getConnectionStatus
}; 