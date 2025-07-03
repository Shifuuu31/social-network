<template>
  <form class="chat-input" @submit.prevent="send">
    <button type="button" @click="showPicker = !showPicker">
            <svg fill="none"  height="24" width="24" xmlns="http://www.w3.org/2000/svg">
<circle stroke-linejoin="round" stroke-linecap="round" stroke-width="2.5" stroke="#707277" r="10" cy="12" cx="12"></circle>
<path stroke-linejoin="round" stroke-linecap="round" stroke-width="2.5" stroke="#707277" d="M8 15C8.91212 16.2144 10.3643 17 12 17C13.6357 17 15.0879 16.2144 16 15"></path>
<path stroke-linejoin="round" stroke-linecap="round" stroke-width="3" stroke="#707277" d="M8.00897 9L8 9M16 9L15.991 9"></path>
</svg>
          </button>
    
    <textarea
      v-model="text"
      placeholder="Type a message..."
      rows="1"
      @keydown.enter.prevent="handleEnter"
    ></textarea>
    <button type="submit">Send</button>

    <div v-if="showPicker" class="emoji-picker">
      <span v-for="emoji in store.emojis" :key="emoji" @click="addEmoji(emoji)">{{ emoji }}</span>
    </div>
  </form>
</template>

<script setup>
import { ref } from 'vue'
import { useChatStore } from '../stores/chatStore'

const text = ref('')
const store = useChatStore()
const showPicker = ref(false)


function send() {
  if (text.value.trim()) {
    store.sendMessage(text.value)
    text.value = ''
  }
  showPicker.value = false
}

function addEmoji(emoji) {
  text.value += emoji
  showPicker.value = false  // <--- hide picker here
}

// Handle enter key: send message if no Shift pressed, else insert newline
function handleEnter(e) {
  if (!e.shiftKey) {
    send()
  } else {
    // Shift+Enter => insert newline manually because default prevented
    text.value += '\n'
  }
}
</script>


<style scoped>

textarea {
    width: 100%;
    margin: 0 10px;
    border-radius: 24px;
    padding: 10px 15px;
    resize: none;
}
.emoji-picker {
    border: 1px solid #ccc;
    padding: 5px;
    max-width: 200px;
    max-height: 150px;  /* limit height */
    overflow-y: auto;   /* enable vertical scrolling */
    border-radius: 24px;
  background: white;
  position: absolute;
  z-index: 10;
  /* optional: add some padding or spacing */
  display: flex;
  flex-wrap: wrap;
  bottom: 8% ;
}

.emoji-picker span {
  cursor: pointer;
  font-size: 1.5rem;  /* slightly bigger for easier clicking */
  padding: 5px;
  user-select: none;
  /* optional: rounded background on hover */
  border-radius: 5px;
  transition: background-color 0.2s;
}
.emoji-picker span:hover {
  background-color: #eee;
}

</style>
