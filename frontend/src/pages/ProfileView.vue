<template>
  <div class="profile-view bg-gray-50 min-h-screen">
    <!-- Banner + Avatar/Header -->
    <section class="relative">
      <img
        src="/banner.jpg"
        alt="Banner"
        class="w-full h-60 object-cover"
      />
      <div class="absolute -bottom-12 left-10 flex items-end space-x-4">
        <img
          :src="user.avatar_url || defaultAvatar"
          alt="Avatar"
          class="w-24 h-24 rounded-full border-4 border-white"
        />
        <div>
          <h1 class="text-2xl font-bold text-gray-900">
            {{ user.first_name }} {{ user.last_name }}
          </h1>
          <p class="text-purple-600">
            {{ user.nickname || 'No nickname' }}
          </p>
        </div>
      </div>
    </section>

    <!-- Actions: Follow / Schedule -->
    <section class="mt-16 px-10 flex justify-between items-center">
      <div></div>
      <div class="flex gap-4">
        <button
          @click="toggleFollow"
          class="bg-purple-600 text-white px-4 py-2 rounded-full hover:bg-purple-700"
        >
          {{ isFollowing ? 'Unfollow' : 'Follow' }}
        </button>
        <button
          class="border border-purple-600 text-purple-600 px-4 py-2 rounded-full hover:bg-purple-50"
        >
          Schedule a meeting
        </button>
      </div>
    </section>

    <!-- Main content -->
    <main class="mt-10 px-10 grid grid-cols-12 gap-6">
      <!-- Left Info Sidebar -->
      <aside class="col-span-3">
        <SidebarInfo :user="user" />
      </aside>

      <!-- Center Tabs and Content -->
      <section class="col-span-6">
        <ProfileTabs :userId="user.id" />
      </section>

      <!-- Right Suggestions/Active -->
      <aside class="col-span-3 space-y-6">
        <SuggestionsList />
        <ActiveUsers />
      </aside>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import SidebarInfo from '@/components/SidebarInfo.vue'
import ProfileTabs from '@/components/ProfileTabs.vue'
import SuggestionsList from '@/components/SuggestionsList.vue'
import ActiveUsers from '@/components/ActiveUsers.vue'

const route = useRoute()
const targetId = Number(route.params.id)

const defaultAvatar = '/default-avatar.png'
const user = ref({
  id: null,
  first_name: '',
  last_name: '',
  nickname: '',
  avatar_url: '',
  is_public: true,
  is_following: false
})
const isFollowing = ref(false)

const fetchProfile = async () => {
  try {
    const res = await fetch('http://localhost:8080/profile/info', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: targetId })
    })
    if (!res.ok) throw new Error(await res.text())
    const data = await res.json()
    user.value = data
    isFollowing.value = data.is_following
  } catch (err) {
    console.error('Profile load error:', err.message)
  }
}

const toggleFollow = async () => {
  const action = isFollowing.value ? 'unfollow' : 'follow'
  try {
    const res = await fetch('http://localhost:8080/follow/follow-unfollow', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ target_id: targetId, action })
    })
    if (!res.ok) throw new Error(await res.text())
    isFollowing.value = !isFollowing.value
  } catch (err) {
    console.error('Follow toggle failed:', err.message)
  }
}

onMounted(fetchProfile)
</script>

<style scoped>
.profile-view section.relative {
  margin-bottom: 4rem;
}
</style>
