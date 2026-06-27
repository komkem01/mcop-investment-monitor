<template>
  <div class="space-y-3 w-full">
    <div
      v-for="(item, index) in items"
      :key="index"
      class="bg-white border transition-all duration-300 overflow-hidden"
      :class="[
        isOpen(index) 
          ? 'border-indigo-300 shadow-[0_4px_16px_rgba(99,102,241,0.04)] rounded-[22px]' 
          : 'border-slate-200/80 hover:border-slate-300 rounded-[22px]'
      ]"
    >
      <!-- Title Toggle Button -->
      <button
        type="button"
        @click="toggle(index)"
        class="w-full flex items-center justify-between px-6 py-4.5 text-left outline-none font-semibold text-slate-800 text-xs sm:text-sm tracking-wide select-none"
      >
        <span>{{ item.title }}</span>
        
        <!-- Animated Chevron Indicator -->
        <span
          class="text-[10px] text-slate-400 transition-transform duration-300 shrink-0 ml-4"
          :class="isOpen(index) ? 'rotate-90 text-indigo-500' : 'text-slate-300'"
        >
          ▶
        </span>
      </button>

      <!-- Expandable Content Panel -->
      <div
        class="transition-all duration-300 ease-in-out overflow-hidden"
        :style="{
          maxHeight: isOpen(index) ? '300px' : '0px',
          opacity: isOpen(index) ? '1' : '0'
        }"
      >
        <div class="px-6 pb-5 text-xs sm:text-sm text-slate-500 font-medium leading-relaxed border-t border-slate-50 pt-3">
          {{ item.content }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  items: {
    type: Array,
    required: true
  },
  // If true, multiple panels can be open at the same time. If false, acts like accordion/radio
  multiple: {
    type: Boolean,
    default: false
  }
})

// Stores the indexes of currently expanded items
const activeIndexes = ref([0]) // default first item open

const isOpen = (index) => {
  return activeIndexes.value.includes(index)
}

const toggle = (index) => {
  if (props.multiple) {
    if (isOpen(index)) {
      activeIndexes.value = activeIndexes.value.filter(i => i !== index)
    } else {
      activeIndexes.value.push(index)
    }
  } else {
    if (isOpen(index)) {
      activeIndexes.value = []
    } else {
      activeIndexes.value = [index]
    }
  }
}
</script>
