<template>
    <i class="uil uil-smile" @click="toggleShowEmojis"></i>
    <div class="emojis" v-show="showEmojis">
        <p class="emoji" v-for="emoji in getEmojis"
           @click.stop="addEmoji(emoji)">{{ emoji }}</p>
    </div>

</template>

<script>
export default {
    props: ["input", "messagebox"],
    data() {
        return {
            emojiCodes: [128512, 128514, 128519, 128520],
            showEmojis: false
        }
    },

    computed: {
        getEmojis() {
            return this.emojiCodes.map((code) => String.fromCodePoint(code))
        }
    },

    methods: {
        addEmoji(emoji) {
            this.input.value += emoji;
        },

        toggleShowEmojis() {
            this.showEmojis = !this.showEmojis
        },



    }
}
</script>

<style scoped>
i {
    font-size: 1.35em;
}

i:hover {
    color: var(--hover-color);
}



.emojis {
    display: flex;
    padding-top: 5px;
    padding-left: 7px;
    width: 100%;
    gap: 5px;
}

.emoji {
    font-size: 1.25em;
    cursor: pointer;
}
</style>