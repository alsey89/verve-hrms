<template>
    <div v-html="renderedMarkdown" class="prose"></div>
</template>
  
<script setup>
import { watchEffect } from 'vue';
import MarkdownIt from 'markdown-it';


const props = defineProps({
    content: String
});

const md = new MarkdownIt({
    breaks: true // Convert '\n' in paragraphs into <br>
});
const renderedMarkdown = ref('');

watchEffect(() => {
    renderedMarkdown.value = md.render(props.content || '');
});
</script>