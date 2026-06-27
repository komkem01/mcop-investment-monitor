<template>
  <div class="bg-white p-6 rounded-[22px] border border-slate-200/80 shadow-[0_4px_18px_rgba(0,0,0,0.015)] flex flex-col h-full">
    <h3 class="text-xs font-extrabold text-slate-800 uppercase tracking-wider mb-5">สัดส่วน Defect แยกตามโมดูล</h3>
    <div class="flex-1 flex flex-col sm:flex-row items-center justify-center gap-8">
      
      <!-- SVG Donut Chart with Linear Gradients -->
      <div v-if="chartData.length > 0" class="relative w-36 h-36 shrink-0">
        <svg viewBox="0 0 42 42" class="w-full h-full transform -rotate-90">
          <defs>
            <!-- Define linear gradients for each color slice -->
            <linearGradient 
              v-for="(gradient, idx) in gradientDefinitions" 
              :key="idx" 
              :id="'donut-grad-' + idx"
              x1="0%" y1="0%" x2="100%" y2="100%"
            >
              <stop offset="0%" :stop-color="gradient.start" />
              <stop offset="100%" :stop-color="gradient.end" />
            </linearGradient>
          </defs>
          
          <!-- Background track -->
          <circle
            cx="21"
            cy="21"
            r="15.91549430918954"
            fill="transparent"
            stroke="#f8fafc"
            stroke-width="2.5"
          />
          
          <!-- Colored segments -->
          <circle
            v-for="(segment, idx) in chartData"
            :key="idx"
            cx="21"
            cy="21"
            r="15.91549430918954"
            fill="transparent"
            :stroke="'url(#donut-grad-' + (idx % gradientDefinitions.length) + ')'"
            stroke-width="2.5"
            :stroke-dasharray="segment.strokeDasharray"
            :stroke-dashoffset="segment.segmentOffset"
            stroke-linecap="round"
            class="transition-all duration-500 ease-out"
          />
        </svg>
        <div class="absolute inset-0 flex flex-col items-center justify-center">
          <span class="text-2xl font-black text-slate-800 tracking-tight leading-none">{{ totalDefectCount }}</span>
          <span class="text-[9px] text-slate-500 font-bold uppercase tracking-widest mt-1">Issues</span>
        </div>
      </div>
      
      <!-- Fallback empty state -->
      <div v-else class="h-36 flex items-center justify-center bg-slate-50/50 border border-slate-100 rounded text-slate-400 text-xs text-center w-full">
        ไม่มีข้อมูล Defect
      </div>

      <!-- Legend -->
      <div v-if="chartData.length > 0" class="flex-1 space-y-2.5 w-full max-h-36 overflow-y-auto pr-1 custom-scrollbar">
        <div 
          v-for="(item, idx) in chartData.slice(0, 5)" 
          :key="idx" 
          class="flex items-center justify-between text-xs transition-colors hover:bg-slate-50/50 p-1 rounded"
        >
          <div class="flex items-center space-x-2 truncate">
            <!-- Sleek color indicator with matching gradient background -->
            <span 
              class="w-2 h-2 rounded-full shrink-0" 
              :style="{ 
                background: `linear-gradient(135deg, ${gradientDefinitions[idx % gradientDefinitions.length].start}, ${gradientDefinitions[idx % gradientDefinitions.length].end})` 
              }"
            ></span>
            <span class="text-slate-700 font-bold truncate">{{ item.name }}</span>
          </div>
          <span class="text-slate-800 font-extrabold shrink-0 ml-2">
            {{ item.defect }} <span class="text-[10px] text-slate-500 font-bold">({{ item.percentage }}%)</span>
          </span>
        </div>
        <div v-if="chartData.length > 5" class="text-[10px] text-slate-500 text-right font-bold italic pr-1">
          และอื่นๆ อีก {{ chartData.length - 5 }} โมดูล
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
  }
})

// Premium theme gradient colors
const gradientDefinitions = [
  { start: '#fb7185', end: '#f43f5e' }, // Rose-400 to Rose-500
  { start: '#f87171', end: '#ef4444' }, // Red-400 to Red-500
  { start: '#fbbf24', end: '#f59e0b' }, // Amber-400 to Amber-500
  { start: '#60a5fa', end: '#3b82f6' }, // Blue-400 to Blue-500
  { start: '#34d399', end: '#10b981' }, // Emerald-400 to Emerald-500
  { start: '#c084fc', end: '#a855f7' }, // Purple-400 to Purple-500
  { start: '#818cf8', end: '#6366f1' }, // Indigo-400 to Indigo-500
]

const totalDefectCount = computed(() => {
  return props.modules.reduce((sum, m) => sum + m.defect, 0)
})

const chartData = computed(() => {
  const validModules = props.modules
    .filter(m => m.defect > 0)
    .map(m => ({
      name: m.name,
      defect: m.defect
    }))
    .sort((a, b) => b.defect - a.defect)
  
  if (totalDefectCount.value === 0) return []
  
  let currentOffset = 100
  
  return validModules.map((m, index) => {
    const percentage = (m.defect / totalDefectCount.value) * 100
    const offset = currentOffset
    currentOffset -= percentage
    
    return {
      name: m.name,
      defect: m.defect,
      percentage: Math.round(percentage),
      strokeDasharray: `${percentage} ${100 - percentage}`,
      segmentOffset: offset
    }
  })
})
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 3px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #e2e8f0;
  border-radius: 2px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #cbd5e1;
}
</style>
