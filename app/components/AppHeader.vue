<template>
  <header class="h-16 bg-white border-b border-slate-200 flex items-center justify-between px-8 shrink-0">
    <!-- Minimalist Search Box -->
    <div class="flex items-center bg-slate-50/70 border border-slate-200 rounded-lg px-3 py-1.5 focus-within:bg-white focus-within:border-slate-400 focus-within:shadow-[0_2px_8px_rgba(0,0,0,0.02)] transition-all duration-200 w-80">
      <span class="text-slate-400 text-xs mr-2 select-none">🔍</span>
      <input
        type="text"
        :value="searchQuery"
        @input="$emit('update:searchQuery', $event.target.value)"
        placeholder="ค้นหาชื่อโมดูล..."
        class="bg-transparent border-none outline-none w-full text-xs text-slate-700 placeholder-slate-400"
      >
    </div>
    
    <!-- Action Menu -->
    <div class="flex items-center space-x-5">
      <button
        @click="$emit('add-bug')"
        class="bg-slate-900 hover:bg-slate-800 active:scale-[0.98] text-white px-4 py-2 rounded-lg text-xs font-medium tracking-wide shadow-sm transition-all duration-150"
      >
        + แจ้งบั๊กใหม่
      </button>
      
      <button 
        @click="$emit('click-notifications')"
        class="w-8 h-8 rounded-lg flex items-center justify-center text-slate-400 hover:text-slate-700 hover:bg-slate-50 text-base transition-all duration-150"
      >
        🔔
      </button>
      
      <div class="w-px h-4 bg-slate-200"></div>
      
      <!-- Avatar profile with subtle gradient background -->
      <div
        @click="$emit('click-profile')"
        class="w-8 h-8 rounded-full bg-gradient-to-tr from-blue-500 to-indigo-600 hover:opacity-90 flex items-center justify-center text-white text-xs font-bold cursor-pointer shadow-sm transition-opacity"
      >
        A
      </div>
    </div>
  </header>
</template>

<script setup>
defineProps({
  searchQuery: {
    type: String,
    default: ''
  }
})

defineEmits(['update:searchQuery', 'add-bug', 'click-notifications', 'click-profile'])
</script>
