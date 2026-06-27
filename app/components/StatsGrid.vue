<template>
  <section class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
    <div
      v-for="(stat, index) in stats"
      :key="index"
      class="relative bg-white p-6 rounded-[22px] border border-slate-200/80 shadow-[0_4px_18px_rgba(0,0,0,0.015)] hover:shadow-[0_8px_24px_rgba(99,102,241,0.05)] hover:-translate-y-1 transition-all duration-300 flex flex-col justify-between overflow-hidden"
    >
      <!-- Top colored gradient line -->
      <span class="absolute top-0 left-0 right-0 h-[3px] bg-gradient-to-r" :class="gradientClass(stat.color)"></span>

      <div class="flex items-center justify-between mb-3">
        <p class="text-xs font-bold text-slate-500 uppercase tracking-wider">{{ stat.label }}</p>
        <span class="text-lg opacity-85 select-none">{{ stat.icon }}</span>
      </div>
      <div>
        <h3 class="text-3xl font-extrabold text-slate-800 tracking-tight leading-none">
          {{ stat.value }}
        </h3>
      </div>
    </div>
  </section>
</template>

<script setup>
defineProps({
  stats: {
    type: Array,
    required: true
  }
})

const gradientClass = (textColorClass) => {
  if (textColorClass.includes('green')) return 'from-emerald-400 to-teal-500'
  if (textColorClass.includes('red')) return 'from-rose-400 to-red-500'
  if (textColorClass.includes('orange')) return 'from-amber-400 to-orange-500'
  if (textColorClass.includes('blue')) return 'from-sky-400 to-indigo-500'
  return 'from-slate-300 to-slate-400'
}
</script>
