<template>
  <div class="chat-view">
    <!-- Debug Panel -->
    <div v-if="showDebug" class="debug-panel">
      <div class="debug-header">
        <h4>🔍 Debug Panel</h4>
        <button @click="clearDebugLogs" class="clear-btn">Clear</button>
      </div>
      <div class="debug-content">
        <div v-for="(log, index) in debugLogs" :key="index" class="debug-log">
          <span class="debug-time">{{ log.time }}</span>
          <span class="debug-message">{{ log.message }}</span>
        </div>
      </div>
    </div>
    
    <div class="chat-container">
      <ChatList 
        :selectedUserId="selectedUserId"
        @select-conversation="selectConversation"
      />
      <ChatWindow 
        :selectedUserId="selectedUserId"
      />
    </div>
    
    <!-- Debug Toggle Button -->
    <button @click="toggleDebug" class="debug-toggle">
      {{ showDebug ? '🔍' : '🐛' }}
    </button>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import ChatList from '../components/chat/ChatList.vue'
import ChatWindow from '../components/chat/ChatWindow.vue'

export default {
  name: 'Chat',
  components: {
    ChatList,
    ChatWindow
  },
  setup() {
    const selectedUserId = ref(null)
    const showDebug = ref(false)
    const debugLogs = ref([])

    const selectConversation = (userId) => {
      selectedUserId.value = userId
      addDebugLog(`Selected conversation with user ${userId}`)
    }

    const toggleDebug = () => {
      showDebug.value = !showDebug.value
      addDebugLog(`Debug panel ${showDebug.value ? 'enabled' : 'disabled'}`)
    }

    const clearDebugLogs = () => {
      debugLogs.value = []
    }

    const addDebugLog = (message) => {
      const time = new Date().toLocaleTimeString()
      debugLogs.value.push({ time, message })
      // Keep only last 50 logs
      if (debugLogs.value.length > 50) {
        debugLogs.value = debugLogs.value.slice(-50)
      }
    }

    // Override console.log to capture debug messages
    onMounted(() => {
      const originalLog = console.log
      console.log = (...args) => {
        originalLog.apply(console, args)
        if (args[0] && typeof args[0] === 'string' && args[0].includes('🔍')) {
          addDebugLog(args.join(' '))
        }
      }
    })

    return {
      selectedUserId,
      showDebug,
      debugLogs,
      selectConversation,
      toggleDebug,
      clearDebugLogs
    }
  }
}
</script>

<style scoped>
.chat-view {
  height: 100vh;
  background: #f8f9fa;
  position: relative;
}

.chat-container {
  height: 100%;
  display: flex;
  background: white;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.debug-panel {
  position: fixed;
  top: 10px;
  right: 10px;
  width: 400px;
  max-height: 300px;
  background: #1a1a1a;
  color: #00ff00;
  border-radius: 8px;
  z-index: 1000;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.debug-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #2a2a2a;
  border-bottom: 1px solid #333;
}

.debug-header h4 {
  margin: 0;
  color: #00ff00;
  font-size: 14px;
}

.clear-btn {
  background: #ff4444;
  color: white;
  border: none;
  padding: 4px 8px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 10px;
}

.clear-btn:hover {
  background: #cc3333;
}

.debug-content {
  max-height: 250px;
  overflow-y: auto;
  padding: 8px;
}

.debug-log {
  margin-bottom: 4px;
  line-height: 1.3;
}

.debug-time {
  color: #888;
  margin-right: 8px;
}

.debug-message {
  color: #00ff00;
  word-break: break-all;
}

.debug-toggle {
  position: fixed;
  bottom: 20px;
  right: 20px;
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: #007bff;
  color: white;
  border: none;
  font-size: 20px;
  cursor: pointer;
  z-index: 1000;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  transition: all 0.3s ease;
}

.debug-toggle:hover {
  background: #0056b3;
  transform: scale(1.1);
}
</style> 