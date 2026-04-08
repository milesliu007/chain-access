<template>
  <div v-if="result" class="result" :class="result.type">
    {{ result.message }}
  </div>
</template>

<script setup>
defineProps({
  result: {
    type: Object,
    default: null,
  },
})
</script>
