<template>

    <div class="post-wrapper">
        <div class="post">
            <div class="user-picture medium"
                 :style="{ backgroundImage: `url(http://localhost:8081/${postData.author.avatar})` }"></div>
            <div class="post-content">
                <router-link :to="{name: 'Profile', params: {id: postData.author.id}}" class="post-author">{{ postData.author.nickname }}</router-link>

                <p class="post-body">{{ postData.content }}</p>
                <img v-if="postData.image" class="post-image" :src="'http://localhost:8081/' + postData.image" alt="">
                <button v-if="!isCommentsOpen" @click="toggleComments" class="btn ">Comments</button>

            </div>

        </div>

        <div v-if="isCommentsOpen">

            <div class="create-comment">
                <textarea v-model="this.comment.body" name="" id="" cols="30" rows="4"
                          placeholder="Add your comment here"></textarea>

                <div class="btns-wrapper">
                    <button class="btn outline hideCommentBtn" @click="toggleComments">Hide
                        comments</button>

                    <div class="add-image">
                        <div class="selected-image" v-if="fileAdded">
                            <p class="additional-info">{{ comment.image.name }}</p>
                            <i class="uil uil-times close" @click="removeImage"></i>
                        </div>

                        <p class="additional-info" v-else>No file chosen
                        </p>

                        <label :for="'upload-img-'+postData.id" >
                            <input type="file" accept="image/png, image/gif, image/jpeg" style=""
                                   @change="checkPicture" ref="fileUpload" :id="'upload-img-'+postData.id" />

                            <div></div>
                        </label>

                    </div>

                </div>

                <button class="btn submitCommentBtn" @click="submitComment(postData.id)">Comment</button>

            </div>

            <div class="comments" v-if="postData.comments">
                <div class="comment" lang="en" v-for="comment in postData.comments">
                    <div class="user-picture medium"
                         :style="{ backgroundImage: `url(http://localhost:8081/${comment.author.avatar})` }"></div>
                    <div class="comment-content">
                        <router-link :to="{name: 'Profile', params: {id: comment.author.id}}" class="comment-author">{{ comment.author.nickname }}</router-link>
                        <p class="comment-body">{{ comment.content }}</p>
                        <img class="comment-image" v-if="comment.image" :src="'http://localhost:8081/' + comment.image"
                             alt="">
                    </div>
                </div>
            </div>
        </div>
    </div>

</template>


<script>
export default {
    name: 'Post',
    data() {
        return {
            isCommentsOpen: false,
            comment: {
                body: "",
                image: {}
            },

        }
    },
    props: {
        postData: {
            type: Object,
            default() {
                return {}
            }
        }
    },

    computed: {
        fileAdded() {
            return this.comment.image.name !== undefined
        },
    },

    methods: {
        toggleComments() {
            this.isCommentsOpen = !this.isCommentsOpen
        },
        async submitComment(post_id) {

            let commentData = new FormData();
            commentData.set('postid', post_id);
            commentData.set('body', this.comment.body);
            commentData.set('image', this.comment.image);

            await fetch('http://localhost:8081/newComment', {
                method: 'POST',
                credentials: 'include',
                body: commentData
            })
            this.$store.dispatch('fetchPosts')
            this.$store.dispatch('fetchMyPosts')
            this.$store.dispatch('getGroupPosts')
            if (this.$route.path != "/main"){
                this.$parent.$parent.getPosts()
            }
            this.comment.body = "";
            this.removeImage();
            console.log('Comment submitted.');
        },
        showPostId(postId) {
            console.log('post id: ', postId);
        },
        checkPicture(e) {
            let files = e.target.files
            if (!files.length) {
                return;
            }
            const file = files[0]

            // console.log("File", file)
            const [extension] = file.type.split("/")
            if ((!(extension == "image"))) {
                console.log('File is not an image.');
                this.$toast.open({
                    message: 'File is not an image.',
                    type: 'error', //One of success, info, warning, error, default
                })
                return
            }
            if (file.size > 2048000) {
                console.log('File size is more than 2 MB.');
                this.$toast.open({
                    message: 'File size is more than 2 MB.',
                    type: 'error', //One of success, info, warning, error, default
                })
                return
            }
            this.comment.image = file;


        },

        removeImage() {
            this.comment.image = {};
            this.$refs.fileUpload.value = "";
        }
    }
}
</script>


<style scoped>
/* POST & COMMENT */
.post-wrapper {
    display: inline-block;
    box-shadow: var(--container-shadow);
    padding: 30px;
    background-color: var(--color-white);
    /* width: 550px; */
    width: 100%;
    border-radius: 10px;
}


.post {
    display: flex;
    gap: 10px;
}


.post-author,
.comment-author {
    font-weight: 500;
}

.post-content,
.comment-content {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 5px;
    flex-grow: 1;
}

.post-image,
.comment-image {
    width: 100%;
    margin: 10px 0 10px 0;
    border-radius: 5px;
}

.post-body,
.comment-body {
    overflow-wrap: anywhere;
}

.post-content button {
    align-self: flex-end;
    margin-top: 10px;
}


.comments>* {
    display: flex;
    gap: 10px;
    border-top: 1px solid #DDDDDD;
    padding-top: 30px;
}

.comments {
    display: flex;
    flex-direction: column;
    gap: 30px;
    margin-top: 30px;
}


.create-comment {
    padding-left: 58px;
}

.create-comment textarea {
    margin: 10px 0;
}


.btns-wrapper {
    display: flex;
    gap: 15px;
    align-items: center;
    justify-content: space-between;
    width: 100%;

}


.hideCommentBtn {
    flex-shrink: 0;
}

.submitCommentBtn {
    margin-left: auto;
    margin-top: 10px;
}

.additional-info {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;

    font-size: 14px;
}
</style>