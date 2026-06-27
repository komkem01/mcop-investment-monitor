<template>
  <Transition name="fade">
    <div 
      v-if="show" 
      class="fixed inset-0 z-50 flex items-center justify-center bg-slate-50/75 backdrop-blur-sm"
    >
      <div class="flex items-center space-x-3 bg-white px-5 py-3 rounded-full border border-slate-100 shadow-sm">
        <!-- Minimal Spinning Ring -->
        <svg class="animate-spin h-4 w-4 text-slate-800" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-10" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-80" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <span class="text-xs font-medium text-slate-700 tracking-wide">
          {{ displayMessage }}
        </span>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  message: {
    type: String,
    default: ''
  }
})

const defaultMessages = [
  'กำลังเตรียมข้อมูล...',
  'เชื่อมต่อฐานข้อมูล MCOP...',
  'จัดเตรียมหน้าแดชบอร์ด...'
]

const currentMessageIndex = ref(0)
let timer = null

onMounted(() => {
  timer = setInterval(() => {
    currentMessageIndex.value = (currentMessageIndex.value + 1) % defaultMessages.length
  }, 2500)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

const displayMessage = computed(() => {
  return props.message || defaultMessages[currentMessageIndex.value]
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
