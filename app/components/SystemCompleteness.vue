<template>
  <div class="bg-white p-6 rounded-[22px] border border-slate-200/80 shadow-[0_4px_18px_rgba(0,0,0,0.015)] flex flex-col h-full">
    <h3 class="text-xs font-extrabold text-slate-800 uppercase tracking-wider mb-5">ความสมบูรณ์ของระบบ (Pass vs Issues)</h3>
    <div class="flex-1 flex flex-col justify-center gap-5">
      <div v-for="item in topModules" :key="item.name" class="w-full">
        <!-- Labels -->
        <div class="flex justify-between text-xs mb-2">
          <span class="font-extrabold text-slate-855">{{ item.name }}</span>
          <span class="text-[11px] font-bold text-slate-500">
            <span class="text-emerald-600 font-extrabold">{{ item.pass }} Pass</span>
            <span class="mx-1.5 text-slate-400">/</span>
            <span class="text-rose-600 font-extrabold">{{ item.defect + item.demoFail }} Issues</span>
          </span>
        </div>
        
        <!-- Progress Bar with rounded sections & smooth layout -->
        <div class="w-full h-2 bg-slate-50 rounded-full flex overflow-hidden shadow-inner">
          <div
            class="bg-gradient-to-r from-emerald-400 to-teal-500 h-full transition-all duration-500"
            :style="{ width: getPercentage(item.pass, item) + '%' }"
            title="Pass"
          ></div>
          <div
            class="bg-gradient-to-r from-amber-400 to-orange-500 h-full transition-all duration-500"
            :style="{ width: getPercentage(item.demoFail, item) + '%' }"
            title="Demo Fail"
          ></div>
          <div
            class="bg-gradient-to-r from-rose-400 to-red-500 h-full transition-all duration-500"
            :style="{ width: getPercentage(item.defect, item) + '%' }"
            title="Defect"
          ></div>
        </div>

        <!-- Custom percentage badges underneath -->
        <div class="flex justify-end space-x-2.5 text-[9px] font-bold mt-1.5 text-slate-400 tracking-wider uppercase">
          <span v-if="item.pass > 0" class="text-emerald-500">{{ Math.round(getPercentage(item.pass, item)) }}% Pass</span>
          <span v-if="item.demoFail > 0" class="text-amber-500">{{ Math.round(getPercentage(item.demoFail, item)) }}% Demo</span>
          <span v-if="item.defect > 0" class="text-rose-500">{{ Math.round(getPercentage(item.defect, item)) }}% Defect</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modules: {
    type: Array,
    required: true
  },
  limit: {
    type: Number,
    default: 3
  }
})

const topModules = computed(() => {
  return props.modules.slice(0, props.limit)
})

const getPercentage = (value, item) => {
  const total = item.defect + item.demoFail + item.pass
  if (total === 0) return 0
  return (value / total) * 100
}
</script>
