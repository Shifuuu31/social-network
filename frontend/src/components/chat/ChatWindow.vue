<template>
  <div class="chat-window">
    <div v-if="!selectedUserId" class="no-selection">
      <div class="placeholder">
        <div class="icon">💬</div>
        <h3>Select a conversation</h3>
        <p>Choose a conversation from the list to start chatting</p>
      </div>
    </div>

    <div v-else class="chat-container">
      <!-- Chat Header -->
      <div class="chat-header">
        <div class="user-info">
          <div class="avatar">
            <img 
              :src="getAvatarUrl(selectedUser?.avatar_url)" 
              :alt="selectedUser?.nickname || selectedUser?.first_name"
            />
          </div>
          <div class="user-details">
            <h4>{{ selectedUser?.nickname || selectedUser?.first_name }}</h4>
            <span class="status" :class="{ online: isUserOnline }">
              {{ isUserOnline ? 'Online' : 'Offline' }}
            </span>
          </div>
        </div>
      </div>

      <!-- Messages Area -->
      <div class="messages-container" ref="messagesContainer">
        <div v-if="loading" class="loading-messages">
          <div class="loading-spinner"></div>
          <p>Loading messages...</p>
        </div>

        <div v-else-if="error" class="error-messages">
          {{ error }}
        </div>

        <div v-else-if="messages.length === 0" class="no-messages">
          <p>No messages yet</p>
          <p class="hint">Start the conversation!</p>
        </div>

        <div v-else class="messages">
          <div
            v-for="message in messages"
            :key="message.id"
            class="message"
            :class="{ 
              'message-sent': message.sender_id === currentUserId,
              'message-received': message.sender_id !== currentUserId
            }"
          >
            <div class="message-content">
              <div class="message-text">{{ message.content }}</div>
              <div class="message-time">
                {{ formatMessageTime(message.created_at) }}
              </div>
            </div>
            <div v-if="message.sender_id === currentUserId && !isTemporaryMessage(message.id)" class="message-actions">
              <button 
                @click="deleteMessage(message.id)"
                class="delete-btn"
                title="Delete message"
              >
                ×
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Message Input -->
      <div class="message-input-container">
        <div class="input-wrapper">
          <textarea
            v-model="newMessage"
            @keydown.enter.prevent="sendMessage"
            @keydown.enter.ctrl="sendMessage"
            placeholder="Type a message..."
            class="message-input"
            rows="1"
            ref="messageInput"
          ></textarea>
          <button 
            type="button"
            class="emoji-btn"
            @click="toggleEmojiPicker"
            title="Add emoji"
          >
            😊
          </button>
          <button 
            @click="sendMessage"
            :disabled="!newMessage.trim() || sending"
            class="send-btn"
          >
            <span v-if="sending">...</span>
            <span v-else>Send</span>
          </button>
          <div v-if="showEmojiPicker" class="emoji-picker-dropdown" ref="emojiPicker">
            <div class="emoji-grid">
              <span v-for="emoji in emojis" :key="emoji" class="emoji-item" @click="insertEmoji(emoji)">{{ emoji }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useAuth } from '../../composables/useAuth.js'
import chatService from '../../services/chatService.js'

export default {
  name: 'ChatWindow',
  props: {
    selectedUserId: {
      type: Number,
      default: null
    }
  },
  setup(props) {
    const { user } = useAuth()
    const messages = ref([])
    const selectedUser = ref(null)
    const loading = ref(false)
    const error = ref(null)
    const newMessage = ref('')
    const sending = ref(false)
    const isUserOnline = ref(false)
    const messagesContainer = ref(null)
    const messageInput = ref(null)
    const showEmojiPicker = ref(false)
    const emojis = [
      '😀','😁','😂','🤣','😊','😍','😘','😜','🤔','😎','😢','😭','😡','👍','🙏','👏','🙌','💪','🔥','🎉','🥳','❤️','💔','😇','😱','😴','🤗','😅','😉','😏','😬','😋','😆','😃','😄','😚','😙','😗','😐','😑','😶','🙄','😯','😲','😳','🥺','😤','😠','😩','😖','😞','😟','😔','😕','🙃','🤑','🤠','😷','🤒','🤕','🤢','🤮','🥵','🥶','😵','🤯','🤓','🧐','😺','😸','😹','😻','😼','😽','🙀','😿','😾'
    ]
    const emojiPicker = ref(null)

    const currentUserId = user.value?.id

    const loadMessages = async () => {
      if (!props.selectedUserId) {
        return;
      }

      loading.value = true
      error.value = null
      
      try {
        const response = await chatService.getConversation(props.selectedUserId, 50, 0)
        
        messages.value = (response.messages || []).reverse() // Show oldest first
        
        await nextTick()
        scrollToBottom()
      } catch (err) {
        console.error('Error loading messages:', err)
        error.value = err.message
      } finally {
        loading.value = false
      }
    }

    const loadSelectedUser = async () => {
      if (!props.selectedUserId) {
        selectedUser.value = null
        return
      }

      try {
        // You might want to create an API endpoint to get user by ID
        // For now, we'll use the existing user data from messages
        const response = await chatService.getConversation(props.selectedUserId, 1, 0)
        
        if (response.messages && response.messages.length > 0) {
          const message = response.messages[0]
          selectedUser.value = {
            id: props.selectedUserId,
            nickname: message.sender_id === props.selectedUserId ? message.sender_name : message.receiver_name,
            first_name: message.sender_id === props.selectedUserId ? message.sender_name : message.receiver_name,
            avatar_url: message.sender_id === props.selectedUserId ? message.sender_avatar : null
          }
        }
      } catch (err) {
        console.error('Error loading selected user:', err)
      }
    }

    const sendMessage = async () => {
      if (!newMessage.value.trim() || sending.value) {
        return;
      }

      const content = newMessage.value.trim()
      
      newMessage.value = ''
      sending.value = true

      try {
        // Add message to local state immediately for better UX
        const tempMessage = {
          id: Date.now(), // Temporary ID
          sender_id: currentUserId,
          receiver_id: props.selectedUserId,
          content: content,
          created_at: new Date().toISOString(),
          sender_name: user.value?.nickname || user.value?.first_name,
          sender_avatar: user.value?.avatar_url
        }
        messages.value.push(tempMessage)
        await nextTick()
        scrollToBottom()

        // Use HTTP API directly since it's working
        const response = await chatService.sendMessageAPI(props.selectedUserId, content)
        
        // Update the temporary message with the real ID from the database
        if (response.message_id) {
          const messageIndex = messages.value.findIndex(msg => msg.id === tempMessage.id)
          if (messageIndex !== -1) {
            messages.value[messageIndex].id = response.message_id
          }
        }
      } catch (err) {
        console.error('Error sending message:', err)
        // Remove the temporary message if sending failed
        messages.value = messages.value.filter(msg => msg.id !== tempMessage.id)
        // Restore message if sending failed
        newMessage.value = content
      } finally {
        sending.value = false
      }
    }

    const deleteMessage = async (messageId) => {
      
      // Don't allow deletion of temporary messages (messages with timestamp IDs)
      if (isTemporaryMessage(messageId)) {
        console.warn('Cannot delete temporary message')
        return
      }
      
      try {
        await chatService.deleteMessage(messageId)
        messages.value = messages.value.filter(msg => msg.id !== messageId)
      } catch (err) {
        console.error('Error deleting message:', err)
      }
    }

    const scrollToBottom = () => {
      if (messagesContainer.value) {
        messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
      }
    }

    const getAvatarUrl = (avatarUrl) => {
      if (!avatarUrl) return '/default-avatar.png'
      if (avatarUrl.startsWith('http')) return avatarUrl
      return `/api/images/${avatarUrl}`
    }

    const formatMessageTime = (timestamp) => {
      const date = new Date(timestamp)
      return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    }

    const isTemporaryMessage = (messageId) => {
      return typeof messageId === 'number' && messageId > 1000000000000
    }

    // WebSocket message handler
    const handleMessage = (message) => {
      
      if (message.sender_id === props.selectedUserId || message.receiver_id === props.selectedUserId) {
        
        // Check if this is a message from the current user (sent via WebSocket)
        if (message.sender_id === currentUserId) {
          
          // Find the temporary message and update it with the real ID
          const tempMessageIndex = messages.value.findIndex(msg => 
            msg.content === message.content && 
            msg.sender_id === message.sender_id &&
            msg.receiver_id === message.receiver_id &&
            isTemporaryMessage(msg.id)
          )
          
          if (tempMessageIndex !== -1) {
            // Update the temporary message with the real ID
            messages.value[tempMessageIndex].id = message.id || Date.now()
            messages.value[tempMessageIndex].created_at = message.timestamp
          }
        } else {
          
          // This is a message from the other user, add it to the conversation
          const newMsg = {
            id: message.id || Date.now(),
            sender_id: message.sender_id,
            receiver_id: message.receiver_id,
            content: message.content,
            created_at: message.timestamp,
            sender_name: message.sender_id === props.selectedUserId ? selectedUser.value?.nickname : user.value?.nickname,
            sender_avatar: message.sender_id === props.selectedUserId ? selectedUser.value?.avatar_url : user.value?.avatar_url
          }
          messages.value.push(newMsg)
          nextTick(() => scrollToBottom())
        }
      }
    }

    // Watch for selected user changes
    watch(() => props.selectedUserId, async (newUserId) => {
      
      if (newUserId) {
        await loadSelectedUser()
        await loadMessages()
      } else {
        messages.value = []
        selectedUser.value = null
      }
    })

    onMounted(async () => {
      
      if (props.selectedUserId) {
        await loadSelectedUser()
        await loadMessages()
      }
      
      // Set up WebSocket message handler
      chatService.onMessage(handleMessage)
    })

    onUnmounted(() => {
      // Cleanup is handled by the chat service
    })

    function toggleEmojiPicker() {
      showEmojiPicker.value = !showEmojiPicker.value
    }

    function insertEmoji(emoji) {
      // Insert emoji at cursor position in textarea
      const textarea = messageInput.value
      if (textarea) {
        const start = textarea.selectionStart
        const end = textarea.selectionEnd
        const text = newMessage.value
        newMessage.value = text.slice(0, start) + emoji + text.slice(end)
        nextTick(() => {
          textarea.focus()
          textarea.selectionStart = textarea.selectionEnd = start + emoji.length
        })
      } else {
        newMessage.value += emoji
      }
      showEmojiPicker.value = false
    }

    // Close emoji picker when clicking outside
    function handleClickOutside(event) {
      if (showEmojiPicker.value && emojiPicker.value && !emojiPicker.value.contains(event.target) && !event.target.classList.contains('emoji-btn')) {
        showEmojiPicker.value = false
      }
    }
    onMounted(() => {
      document.addEventListener('mousedown', handleClickOutside)
    })
    onUnmounted(() => {
      document.removeEventListener('mousedown', handleClickOutside)
    })

    return {
      messages,
      selectedUser,
      loading,
      error,
      newMessage,
      sending,
      isUserOnline,
      currentUserId,
      messagesContainer,
      messageInput,
      sendMessage,
      deleteMessage,
      getAvatarUrl,
      formatMessageTime,
      isTemporaryMessage,
      showEmojiPicker,
      emojis,
      toggleEmojiPicker,
      insertEmoji,
      emojiPicker
    }
  }
}
</script>

<style scoped>
.chat-window {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: white;
}

.no-selection {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8f9fa;
}

.placeholder {
  text-align: center;
  color: #6c757d;
}

.placeholder .icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.placeholder h3 {
  margin: 0 0 0.5rem 0;
  color: #1a1a1a;
}

.chat-container {
  display: flex;
  width: 100%;    
  flex-direction: column;
   height: 100%;
    overflow: auto;
  margin: 0 auto;
  background: white;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
}
.chat-header {
  padding: 1rem;
  border-bottom: 1px solid #e1e5e9;
  background: white;
}

.user-info {
  display: flex;
  align-items: center;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  margin-right: 1rem;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-details h4 {
  margin: 0 0 0.25rem 0;
  color: #1a1a1a;
}

.status {
  font-size: 0.8rem;
  color: #6c757d;
}

.status.online {
  color: #28a745;
}

.messages-container {
  max-height: 70vh;
  overflow-y: auto;
  padding: 1rem;
  background: #f8f9fa;
}

.loading-messages, .error-messages, .no-messages {
  text-align: center;
  color: #6c757d;
  padding: 2rem;
}

.loading-spinner {
  border: 2px solid #f3f3f3;
  border-top: 2px solid #007bff;
  border-radius: 50%;
  width: 20px;
  height: 20px;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.messages {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.message {
  display: flex;
  align-items: flex-end;
  gap: 0.5rem;
}

.message-sent {
  justify-content: flex-end;
}

.message-received {
  justify-content: flex-start;
}

.message-content {
  max-width: 70%;
  padding: 0.75rem 1rem;
  border-radius: 1rem;
  position: relative;
}

.message-sent .message-content {
  background: #007bff;
  color: white;
  border-bottom-right-radius: 0.25rem;
}

.message-received .message-content {
  background: white;
  color: #1a1a1a;
  border: 1px solid #e1e5e9;
  border-bottom-left-radius: 0.25rem;
}

.message-text {
  margin-bottom: 0.25rem;
  word-wrap: break-word;
}

.message-time {
  font-size: 0.75rem;
  opacity: 0.7;
}

.message-actions {
  opacity: 0;
  transition: opacity 0.2s;
}

.message:hover .message-actions {
  opacity: 1;
}

.delete-btn {
  background: none;
  border: none;
  color: #dc3545;
  cursor: pointer;
  font-size: 1.2rem;
  padding: 0.25rem;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.delete-btn:hover {
  background: #f8d7da;
}

.message-input-container {
  padding: 1rem;
  border-top: 1px solid #e1e5e9;
  background: white;
  position: relative; /* Added for emoji picker positioning */
}

.input-wrapper {
  display: flex;
  gap: 0.5rem;
  align-items: flex-end;
}

.message-input {
  flex: 1;
  border: 1px solid #e1e5e9;
  border-radius: 1rem;
  padding: 0.75rem 1rem;
  resize: none;
  font-family: inherit;
  font-size: 0.9rem;
  line-height: 1.4;
  max-height: 100px;
  min-height: 40px;
}

.message-input:focus {
  outline: none;
  border-color: #007bff;
}

.send-btn {
  background: #007bff;
  color: white;
  border: none;
  border-radius: 1rem;
  padding: 0.75rem 1.5rem;
  cursor: pointer;
  font-weight: 600;
  transition: background-color 0.2s;
  min-width: 80px;
}

.send-btn:hover:not(:disabled) {
  background: #0056b3;
}

.send-btn:disabled {
  background: #6c757d;
  cursor: not-allowed;
}

.hint {
  font-size: 0.9rem;
  color: #adb5bd;
  margin-top: 0.5rem;
}

.emoji-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  margin-right: 0.5rem;
  transition: background 0.2s;
  border-radius: 50%;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.emoji-btn:hover {
  background: #f0f0f0;
}
.emoji-picker-dropdown {
  position: absolute;
  bottom: 48px;
  left: 0;
  z-index: 10;
  background: white;
  border: 1px solid #e1e5e9;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  padding: 8px;
  width: 320px;
  max-width: 90vw;
}
.emoji-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  max-height: 180px;
  overflow-y: auto;
}
.emoji-item {
  font-size: 1.3rem;
  cursor: pointer;
  padding: 4px;
  border-radius: 6px;
  transition: background 0.15s;
}
.emoji-item:hover {
  background: #f0f0f0;
}
</style> 