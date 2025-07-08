import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// Font Awesome setup
import { library } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { faHome, faBell, faEnvelope, faUser, faUsers, faUserFriends } from '@fortawesome/free-solid-svg-icons'

library.add(faHome, faBell, faEnvelope, faUser, faUsers, faUserFriends)

const app = createApp(App)
app.component('font-awesome-icon', FontAwesomeIcon)
app.use(router)
app.mount('#app')