<template>
  <section class="space-y-4">
    <h2 class="text-xs font-extrabold text-slate-800 uppercase tracking-wider">
      สถานะรายโมดูล (Module Status)
    </h2>
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-5">
      <div
        v-for="module in modules"
        :key="module.name"
        @click="$emit('select-module', module)"
        class="bg-white rounded-[22px] border border-slate-200/80 shadow-[0_4px_18px_rgba(0,0,0,0.015)] hover:shadow-[0_12px_24px_rgba(99,102,241,0.06)] hover:-translate-y-1 hover:border-indigo-300 transition-all duration-300 cursor-pointer flex flex-col justify-between h-40 p-5 group"
      >
        <!-- Header: Module Name & Icon -->
        <div class="border-b border-slate-200/80 pb-3 mb-3 flex items-center justify-between">
          <span class="text-xs font-extrabold text-slate-800 group-hover:text-blue-600 transition-colors duration-250 truncate">{{ module.name }}</span>
          <span class="text-base select-none shrink-0">{{ module.icon }}</span>
        </div>
        
        <!-- Metrics List with soft hover backgrounds -->
        <div class="space-y-1 flex-1 flex flex-col justify-center">
          <div class="flex items-center justify-between text-[11px] font-medium py-0.5 px-1.5 rounded transition-colors hover:bg-rose-50/50">
            <span class="text-slate-600 font-bold flex items-center">
              <span class="w-1.5 h-1.5 rounded-full bg-rose-400 mr-2 shadow-[0_0_6px_#f43f5e]"></span>
              Defect
            </span>
            <span class="text-rose-500 font-bold">{{ module.defect }}</span>
          </div>
          
          <div class="flex items-center justify-between text-[11px] font-medium py-0.5 px-1.5 rounded transition-colors hover:bg-amber-50/50">
            <span class="text-slate-600 font-bold flex items-center">
              <span class="w-1.5 h-1.5 rounded-full bg-amber-400 mr-2 shadow-[0_0_6px_#f59e0b]"></span>
              Demo Fail
            </span>
            <span class="text-amber-500 font-bold">{{ module.demoFail }}</span>
          </div>
          
          <div class="flex items-center justify-between text-[11px] font-medium py-0.5 px-1.5 rounded transition-colors hover:bg-emerald-50/50">
            <span class="text-slate-600 font-bold flex items-center">
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 mr-2 shadow-[0_0_6px_#10b981]"></span>
              Pass
            </span>
            <span class="text-emerald-500 font-bold">{{ module.pass }}</span>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
defineProps({
  modules: {
    type: Array,
    required: true
  }
})

defineEmits(['select-module'])
</script>
