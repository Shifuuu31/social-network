// chatService.js
import { io } from 'socket.io-client';

const API_BASE = '/api';
const CHAT_BASE_URL = `${API_BASE}/chat`;

// Debug logging utility
const DEBUG = true;
const debugLog = (method, message, data = null) => {
  if (DEBUG) {
    console.log(`🔍 [ChatService.${method}] ${message}`, data || '');
  }
};

class ChatService {
  constructor() {
    this.socket = null;
    this.isConnected = false;
    this.messageHandlers = new Set();
    this.connectionHandlers = new Set();
    debugLog('constructor', 'ChatService initialized');
  }

  // Connect to WebSocket
  connect(userId) {
    debugLog('connect', `Attempting to connect with userId: ${userId}`);
    
    if (this.socket && this.isConnected) {
      debugLog('connect', 'Already connected, returning existing connection');
      return Promise.resolve();
    }

    return new Promise((resolve, reject) => {
      try {
        debugLog('connect', 'Creating WebSocket connection to ws://localhost:8080/chat/ws');
        // For now, we'll use a simple WebSocket connection
        // In production, you might want to use Socket.IO with proper authentication
        this.socket = new WebSocket(`ws://localhost:8080/chat/ws`);
        
        this.socket.onopen = () => {
          debugLog('connect', 'WebSocket connection opened successfully');
          console.log('WebSocket connected');
          this.isConnected = true;
          this.notifyConnectionHandlers(true);
          resolve();
        };

        this.socket.onmessage = (event) => {
          debugLog('onmessage', `Received WebSocket message: ${event.data}`);
          try {
            const message = JSON.parse(event.data);
            debugLog('onmessage', 'Parsed message:', message);
            this.notifyMessageHandlers(message);
          } catch (error) {
            console.error('Error parsing WebSocket message:', error);
            debugLog('onmessage', 'Failed to parse message:', error);
          }
        };

        this.socket.onclose = () => {
          debugLog('connect', 'WebSocket connection closed');
          console.log('WebSocket disconnected');
          this.isConnected = false;
          this.notifyConnectionHandlers(false);
        };

        this.socket.onerror = (error) => {
          debugLog('connect', 'WebSocket connection error:', error);
          console.error('WebSocket error:', error);
          this.isConnected = false;
          this.notifyConnectionHandlers(false);
          reject(error);
        };

      } catch (error) {
        debugLog('connect', 'Error creating WebSocket:', error);
        console.error('Error connecting to WebSocket:', error);
        reject(error);
      }
    });
  }

  // Disconnect from WebSocket
  disconnect() {
    debugLog('disconnect', 'Disconnecting WebSocket');
    if (this.socket) {
      this.socket.close();
      this.socket = null;
      this.isConnected = false;
    }
  }

  // Send a message via WebSocket
  sendMessage(receiverId, content) {
    debugLog('sendMessage', `Sending WebSocket message to ${receiverId}:`, { content });
    
    if (!this.socket || !this.isConnected) {
      const error = 'WebSocket not connected';
      debugLog('sendMessage', error);
      throw new Error(error);
    }

    const message = {
      type: 'message',
      receiver_id: receiverId,
      content: content
    };

    debugLog('sendMessage', 'Sending WebSocket message:', message);
    this.socket.send(JSON.stringify(message));
  }

  // Send message via HTTP API (fallback)
  async sendMessageAPI(receiverId, content) {
    debugLog('sendMessageAPI', `Sending HTTP message to ${receiverId}:`, { content });
    
    const payload = {
      receiver_id: receiverId,
      content: content
    };
    
    debugLog('sendMessageAPI', `Making POST request to ${CHAT_BASE_URL}/send with payload:`, payload);
    debugLog('sendMessageAPI', `Request URL: ${CHAT_BASE_URL}/send`);
    debugLog('sendMessageAPI', `Request headers:`, { 'Content-Type': 'application/json' });
    debugLog('sendMessageAPI', `Request body:`, JSON.stringify(payload));
    
    try {
      const response = await fetch(`${CHAT_BASE_URL}/send`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(payload)
      });

      debugLog('sendMessageAPI', `Response status: ${response.status} ${response.statusText}`);
      debugLog('sendMessageAPI', `Response headers:`, Object.fromEntries(response.headers.entries()));
      
      // Get the raw response text first
      const responseText = await response.text();
      debugLog('sendMessageAPI', `Raw response text:`, responseText);
      
      if (!response.ok) {
        let errorData;
        try {
          errorData = JSON.parse(responseText);
        } catch (e) {
          errorData = { error: responseText };
        }
        debugLog('sendMessageAPI', 'Error response:', errorData);
        throw new Error(errorData.error || 'Failed to send message');
      }

      // Parse the JSON response
      let responseData;
      try {
        responseData = JSON.parse(responseText);
      } catch (e) {
        debugLog('sendMessageAPI', 'Failed to parse JSON response:', e);
        throw new Error('Invalid JSON response from server');
      }
      
      debugLog('sendMessageAPI', 'Success response:', responseData);
      return responseData;
    } catch (error) {
      debugLog('sendMessageAPI', 'Exception occurred:', error);
      throw error;
    }
  }

  // Get conversation between two users
  async getConversation(otherUserId, limit = 50, offset = 0) {
    debugLog('getConversation', `Fetching conversation with ${otherUserId}, limit: ${limit}, offset: ${offset}`);
    
    const payload = {
      other_user_id: otherUserId,
      limit: limit,
      offset: offset
    };
    
    debugLog('getConversation', `Making POST request to ${CHAT_BASE_URL}/conversation with payload:`, payload);
    
    try {
      const response = await fetch(`${CHAT_BASE_URL}/conversation`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(payload)
      });

      debugLog('getConversation', `Response status: ${response.status} ${response.statusText}`);
      debugLog('getConversation', `Response headers:`, Object.fromEntries(response.headers.entries()));
      
      // Get the raw response text first
      const responseText = await response.text();
      debugLog('getConversation', `Raw response text:`, responseText);
      
      if (!response.ok) {
        let errorData;
        try {
          errorData = JSON.parse(responseText);
        } catch (e) {
          errorData = { error: responseText };
        }
        debugLog('getConversation', 'Error response:', errorData);
        throw new Error(errorData.error || 'Failed to get conversation');
      }

      // Parse the JSON response
      let responseData;
      try {
        responseData = JSON.parse(responseText);
      } catch (e) {
        debugLog('getConversation', 'Failed to parse JSON response:', e);
        throw new Error('Invalid JSON response from server');
      }
      
      debugLog('getConversation', 'Success response:', responseData);
      return responseData;
    } catch (error) {
      debugLog('getConversation', 'Exception occurred:', error);
      throw error;
    }
  }

  // Get recent conversations
  async getRecentConversations(limit = 20) {
    debugLog('getRecentConversations', `Fetching recent conversations, limit: ${limit}`);
    
    const url = `${CHAT_BASE_URL}/recent?limit=${limit}`;
    debugLog('getRecentConversations', `Making GET request to ${url}`);
    
    try {
      const response = await fetch(url, {
        method: 'GET',
        credentials: 'include',
        headers: { 'Accept': 'application/json' }
      });

      debugLog('getRecentConversations', `Response status: ${response.status} ${response.statusText}`);
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        debugLog('getRecentConversations', 'Error response:', errorData);
        throw new Error(errorData.error || 'Failed to get recent conversations');
      }

      const responseData = await response.json();
      debugLog('getRecentConversations', 'Success response:', responseData);
      return responseData;
    } catch (error) {
      debugLog('getRecentConversations', 'Exception occurred:', error);
      throw error;
    }
  }

  // Delete a message
  async deleteMessage(messageId) {
    debugLog('deleteMessage', `Deleting message with ID: ${messageId}`);
    
    const url = `${CHAT_BASE_URL}/delete?message_id=${messageId}`;
    debugLog('deleteMessage', `Making DELETE request to ${url}`);
    
    try {
      const response = await fetch(url, {
        method: 'DELETE',
        credentials: 'include'
      });

      debugLog('deleteMessage', `Response status: ${response.status} ${response.statusText}`);
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        debugLog('deleteMessage', 'Error response:', errorData);
        throw new Error(errorData.error || 'Failed to delete message');
      }

      const responseData = await response.json();
      debugLog('deleteMessage', 'Success response:', responseData);
      return responseData;
    } catch (error) {
      debugLog('deleteMessage', 'Exception occurred:', error);
      throw error;
    }
  }

  // Add message handler
  onMessage(handler) {
    debugLog('onMessage', 'Adding message handler');
    this.messageHandlers.add(handler);
    return () => this.messageHandlers.delete(handler);
  }

  // Add connection status handler
  onConnectionChange(handler) {
    debugLog('onConnectionChange', 'Adding connection change handler');
    this.connectionHandlers.add(handler);
    return () => this.connectionHandlers.delete(handler);
  }

  // Notify message handlers
  notifyMessageHandlers(message) {
    debugLog('notifyMessageHandlers', `Notifying ${this.messageHandlers.size} handlers with message:`, message);
    this.messageHandlers.forEach(handler => {
      try {
        handler(message);
      } catch (error) {
        console.error('Error in message handler:', error);
        debugLog('notifyMessageHandlers', 'Handler error:', error);
      }
    });
  }

  // Notify connection handlers
  notifyConnectionHandlers(isConnected) {
    debugLog('notifyConnectionHandlers', `Notifying ${this.connectionHandlers.size} handlers, connected: ${isConnected}`);
    this.connectionHandlers.forEach(handler => {
      try {
        handler(isConnected);
      } catch (error) {
        console.error('Error in connection handler:', error);
        debugLog('notifyConnectionHandlers', 'Handler error:', error);
      }
    });
  }

  // Get connection status
  getConnectionStatus() {
    debugLog('getConnectionStatus', `Connection status: ${this.isConnected}`);
    return this.isConnected;
  }
}

// Create singleton instance
const chatService = new ChatService();
export default chatService; 