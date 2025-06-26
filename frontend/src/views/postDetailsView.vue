<template>
    <div class="post-detail-view">
        <div class="back-header">
            <button @click="$router.go(-1)" class="back-btn">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M20 11H7.83l5.59-5.59L12 4l-8 8 8 8 1.41-1.41L7.83 13H20v-2z" />
                </svg>
                <span>Back</span>
            </button>
            <h2>Post</h2>
        </div>

        <div v-if="loading" class="loading-spinner">Loading...</div>

        <div v-if="error" class="error-message">{{ error }}</div>

        <div v-if="post && !loading" class="post-detail-container">
            <PostItem :post="post" :is-detail="true" @reply="focusCommentForm" />

            <CommentForm ref="commentFormRef" :post-id="post.id" @comment-added="handleCommentAdded" />

            <div class="comments-section">
                <h3 v-if="comments.length > 0">
                    {{ comments.length }} {{ comments.length === 1 ? 'Reply' : 'Replies' }}
                </h3>

                <div v-if="loadingComments" class="loading-comments">
                    Loading comments...
                </div>

                <div v-if="comments.length === 0 && !loadingComments" class="no-comments">
                    No replies yet. Be the first to reply!
                </div>

                <CommentItem v-for="comment in comments" :key="comment.id" :comment="comment"
                    @reply="handleReplyToComment" />
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import PostItem from '@/components/posts/PostItem.vue'
import CommentForm from '@/components/comments/CommentForm.vue'
import CommentItem from '@/components/comments/CommentItem.vue'

const route = useRoute()
const commentFormRef = ref(null)

const post = ref(null)
const comments = ref([])
const loading = ref(false)
const loadingComments = ref(false)
const error = ref(null)

const fetchPost = async () => {

    try {
        loading.value = true
        error.value = null
        const response = await fetch(`/api/post/feed`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ id: +route.params.id, type: 'single', start: 1, n_post: 1 })
        })
        if (!response.ok) throw new Error('Failed to fetch post')
        const data = await response.json()
        console.log('Fetched post data:', data);
        console.log('Post ID from route:', route.params.id);

        post.value = Array.isArray(data) ? data[0] : data

    } catch (err) {
        error.value = err.message
        console.error('Failed to fetch post:', err)
    } finally {
        loading.value = false
    }
}

const fetchComments = async () => {
    try {
        loadingComments.value = true
        const response = await fetch(`/api/comments/${route.params.id}/comment`)
        if (!response.ok) throw new Error('Failed to fetch comments')
        const data = await response.json()
        comments.value = Array.isArray(data) ? data : []
    } catch (err) {
        console.error('Failed to load comments:', err)
        comments.value = []
    } finally {
        loadingComments.value = false
    }
}

const handleCommentAdded = (newComment) => {
    const normalizedComment = {
        id: newComment.id,
        content: newComment.content || '',
        author: newComment.author || newComment.owner_name || 'Anonymous',
        createdAt: newComment.createdAt || newComment.created_at || new Date().toISOString(),
        likes: newComment.likes || 0,
        ...newComment
    }

    comments.value.unshift(normalizedComment)

    if (post.value) {
        post.value.replies = (post.value.replies || 0) + 1
    }
}

const handleReplyToComment = (comment) => {
    if (commentFormRef.value) {
        commentFormRef.value.replyTo(comment.author || 'Anonymous')
    }
}

const focusCommentForm = () => {
    if (commentFormRef.value) {
        commentFormRef.value.focus()
    }
}

onMounted(() => {
    fetchPost()
    fetchComments()
})

watch(
    () => route.params.id,
    async (newId, oldId) => {
        console.log(`Route param changed from ${oldId} to ${newId}`);

        if (newId !== oldId) {
            try {
                post.value = null
                comments.value = []
                error.value = null
                loading.value = true
                loadingComments.value = true

                await Promise.all([fetchPost(), fetchComments()])
            } catch (err) {
                error.value = 'Failed to load post data'
                console.error(err)
            } finally {
                loading.value = false
                loadingComments.value = false
            }
        }
    },
    { immediate: true }
)
</script>

<style scoped>
.post-detail-view {
    overflow: auto;
    width: 100%;
    margin: 0 auto;
    border-left: 1px solid var(--twitter-extra-light-gray);
    border-right: 1px solid var(--twitter-extra-light-gray);
    min-height: 100vh;
    background: white;
}

.back-header {
    position: sticky;
    top: 0;
    background: rgba(255, 255, 255, 0.8);
    backdrop-filter: blur(12px);
    border-bottom: 1px solid var(--twitter-extra-light-gray);
    padding: 1rem;
    display: flex;
    align-items: center;
    gap: 2rem;
    z-index: 10;
}

.back-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: none;
    border: none;
    color: var(--twitter-dark);
    font-size: 1rem;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 9999px;
    transition: background-color 0.2s;
}

.back-btn:hover {
    background-color: var(--twitter-extra-extra-light-gray);
}

.back-header h2 {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 700;
}

.loading-spinner {
    text-align: center;
    padding: 2rem;
    color: var(--twitter-gray);
}

.error-message {
    background: #fee;
    color: #c00;
    padding: 1rem;
    margin: 1rem;
    border-radius: 8px;
    border: 1px solid #fcc;
}

.post-detail-container {
    background: white;
}

.comments-section {
    border-top: 1px solid var(--twitter-extra-light-gray);
}

.comments-section h3 {
    padding: 1rem;
    margin: 0;
    font-size: 1.1rem;
    font-weight: 700;
    color: var(--twitter-dark);
    border-bottom: 1px solid var(--twitter-extra-light-gray);
}

.loading-comments {
    text-align: center;
    padding: 2rem;
    color: var(--twitter-gray);
}

.no-comments {
    text-align: center;
    padding: 2rem;
    color: var(--twitter-gray);
    font-style: italic;
}

:root {
    --twitter-blue: #1da1f2;
    --twitter-dark: #14171a;
    --twitter-gray: #657786;
    --twitter-light-gray: #aab8c2;
    --twitter-extra-light-gray: #e1e8ed;
    --twitter-extra-extra-light-gray: #f5f8fa;
}
</style>