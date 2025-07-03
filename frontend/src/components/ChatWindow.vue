<template>
  <div class="chat-window">
    <!-- 🧠 Fixed user/group navbar -->
    <div class="chat-header">
      <template v-if="store.activeType === 'private'">
        <img :src="activeUser.avatarUrl || defaultAvatar" class="header-avatar" />
        <div class="header-info">
          <strong>{{ activeUser.name }}</strong>
          <span :class="['status-dot', activeUser.online ? 'online' : 'offline']"></span>
          <small>{{ activeUser.online ? 'Online' : 'Offline' }}</small>
        </div>
      </template>

      <template v-else-if="store.activeType === 'group'">
        <img :src="activeGroup.avatarUrl || defaultGroupAvatar" class="header-avatar" />
        <div class="header-info">
          <strong>{{ activeGroup.name }}</strong>
          <small>Group chat</small>
        </div>
      </template>
    </div>

    <!-- 💬 Scrollable message list -->
    <div class="chat-body">
      <div
        v-for="msg in store.activeMessages"
        :key="msg.id"
        :class="['chat-message', msg.sender_id === store.currentUser.id ? 'right' : 'left']"
      >
        <div class="message-content">{{ msg.content }}</div>
        <div class="message-meta">
          <small class="timestamp">{{ formatTimestamp(msg.created_at) }}</small>
        </div>
      </div>

      <div v-if="store.typingUsers && store.typingUsers.length" class="typing-indicator">
  {{ typingText }}
</div>

    </div>
  </div>
</template>


<script setup>
import { computed } from 'vue'
import { useChatStore } from '../stores/chatStore'

const store = useChatStore()

const defaultAvatar = '/src/assets/avatar.png'
const defaultGroupAvatar = '/src/assets/group_avatar.png'

const activeUser = computed(() => store.users.find(u => u.id === store.activeTargetId))
const activeGroup = computed(() => store.groups.find(g => g.id === store.activeTargetId))

function formatTimestamp(ts) {
  if (!ts) return ''
  const date = new Date(ts)
  const now = new Date()
  const diff = Math.floor((now - date) / 60000)

  if (diff < 1) return 'now'
  if (diff < 60) return `${diff}m ago`
  if (diff < 1440) return `${Math.floor(diff / 60)}h ago`
  return date.toLocaleDateString()
}

const typingText = computed(() => {
  if (!store.typingUsers.length) return ''
  return store.typingUsers.length === 1
    ? `${store.typingUsers[0].name} is typing...`
    : 'Several people are typing...'
})
</script>


<style scoped>
.chat-window {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.chat-header {
  display: flex;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #ddd;
  background: #fff;
  position: sticky;
  top: 0;
  z-index: 10;
}

.header-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 10px;
}

.header-info {
  display: flex;
  flex-direction: column;
}

.status-dot {
  display: inline-block;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  margin-right: 4px;
}

.status-dot.online {
  background: #28a745;
}

.status-dot.offline {
  background: #ccc;
}

.chat-body {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
}

.chat-message {
  max-width: 60%;
  margin: 0.5rem 0;
  padding: 0.6rem 1rem;
  display: inline-block;
  clear: both;
}

.chat-message.left {
  background: #e0e0e0;
  border-radius:  0 12px 12px 12px ;
  float: left;
}

.chat-message.right {
  background: #b3d4fc;
  border-radius: 12px 12px 0 12px ;
  float: right;
}

.message-meta {
  text-align: right;
  font-size: 0.75em;
  color: #888;
}

.typing-indicator {
  font-style: italic;
  color: #666;
  margin-top: 10px;
}
</style>

