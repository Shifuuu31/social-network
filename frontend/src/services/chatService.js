class ChatService {
  constructor() {
    this.ws = null
    this.messageHandlers = []
    this.connectionHandlers = []
    this.connected = false
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
    this.baseUrl = 'ws://localhost:8080'
  }

  connect(userId) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      return
    }

    try {
      this.ws = new WebSocket(`${this.baseUrl}/ws/connect`)
      
      this.ws.onopen = () => {
        console.log('WebSocket connected')
        this.connected = true
        this.reconnectAttempts = 0
        this.notifyConnectionHandlers(true)
      }

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          this.notifyMessageHandlers(data)
        } catch (error) {
          console.error('Error parsing WebSocket message:', error)
        }
      }

      this.ws.onclose = () => {
        console.log('WebSocket disconnected')
        this.connected = false
        this.notifyConnectionHandlers(false)
        this.attemptReconnect()
      }

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error)
        this.connected = false
        this.notifyConnectionHandlers(false)
      }
    } catch (error) {
      console.error('Failed to connect to WebSocket:', error)
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
      this.connected = false
    }
  }

  attemptReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      setTimeout(() => {
        console.log(`Attempting to reconnect... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`)
        this.connect()
      }, 2000 * this.reconnectAttempts)
    }
  }

  sendMessage(message) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message))
    } else {
      console.error('WebSocket is not connected')
    }
  }

  onMessage(handler) {
    this.messageHandlers.push(handler)
  }

  onConnectionChange(handler) {
    this.connectionHandlers.push(handler)
  }

  notifyMessageHandlers(message) {
    this.messageHandlers.forEach(handler => {
      try {
        handler(message)
      } catch (error) {
        console.error('Error in message handler:', error)
      }
    })
  }

  notifyConnectionHandlers(connected) {
    this.connectionHandlers.forEach(handler => {
      try {
        handler(connected)
      } catch (error) {
        console.error('Error in connection handler:', error)
      }
    })
  }

  joinGroup(groupId) {
    this.sendMessage({
      type: 'join_group',
      group_id: groupId
    })
  }

  sendGroupMessage(groupId, content) {
    this.sendMessage({
      type: 'group',
      group_id: groupId,
      content: content
    })
  }

  sendPrivateMessage(receiverId, content) {
    this.sendMessage({
      type: 'private',
      receiver_id: receiverId,
      content: content
    })
  }
}

// Create a singleton instance
const chatService = new ChatService()
export default chatService
