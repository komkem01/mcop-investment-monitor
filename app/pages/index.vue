<template>
  <div class="flex h-screen bg-transparent font-sans text-slate-800">
    <!-- Sidebar Navigation -->
    <AppSidebar :menus="sidebarMenus" @select-menu="handleSelectMenu" />

    <!-- Main Content Area -->
    <main class="flex-1 flex flex-col overflow-hidden">
      <!-- Content Scroll Container -->
      <div class="flex-1 overflow-y-auto p-8 space-y-8 custom-scrollbar">
        
        <!-- Search & Action Row (Inline, No Header Bar) -->
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <!-- Search input box - Cute Pill Rounded -->
          <div class="flex items-center bg-white border border-slate-200/80 rounded-full px-4.5 py-2.5 focus-within:border-indigo-300 focus-within:shadow-[0_4px_16px_rgba(99,102,241,0.06)] transition-all duration-200 w-full sm:w-80">
            <span class="text-slate-400 text-xs mr-2 select-none">🔍</span>
            <input
              type="text"
              v-model="searchQuery"
              placeholder="ค้นหาชื่อโมดูล..."
              class="bg-transparent border-none outline-none w-full text-xs text-slate-600 placeholder-slate-400 font-medium"
            >
          </div>

          <!-- Add Bug Action Button - Cute Pill Rounded -->
          <button
            @click="handleAddBug"
            class="bg-indigo-600 hover:bg-indigo-500 active:scale-[0.96] text-white px-5 py-2.5 rounded-full text-xs font-semibold tracking-wide shadow-[0_4px_14px_rgba(99,102,241,0.25)] hover:shadow-[0_6px_20px_rgba(99,102,241,0.35)] transition-all duration-200 shrink-0 self-start sm:self-auto"
          >
            + แจ้งบั๊กใหม่
          </button>
        </div>

        <!-- Loading Spinner (Same as other pages) -->
        <div
          v-if="isLoading"
          class="flex flex-col items-center justify-center py-40 space-y-6"
        >
          <div class="relative flex items-center justify-center">
            <div
              class="w-16 h-16 rounded-full border-2 border-indigo-500/10 border-t-indigo-500 animate-spin shadow-[0_0_15px_rgba(99,102,241,0.15)]"
            ></div>
            <div
              class="absolute w-8 h-8 rounded-full bg-gradient-to-br from-indigo-500 to-sky-400 animate-pulse shadow-[0_0_20px_rgba(99,102,241,0.5)] opacity-80"
            ></div>
          </div>
          <div class="flex flex-col items-center space-y-1">
            <span
              class="text-xs font-black tracking-[0.25em] text-transparent bg-clip-text bg-gradient-to-r from-slate-600 via-indigo-500 to-sky-500 uppercase animate-pulse"
            >
              Loading Data
            </span>
            <span class="text-[9px] font-bold text-slate-400/80 tracking-wider"
              >กำลังดึงข้อมูลจากระบบหลัก...</span
            >
          </div>
        </div>

        <template v-else>
          <!-- Stats summary cards -->
          <StatsGrid :stats="globalStats" />

          <!-- Double columns section -->
          <section class="grid grid-cols-1 lg:grid-cols-2 gap-8">
            <ModuleChart :modules="moduleStats" />
            <SystemCompleteness :modules="moduleStats" />
          </section>

          <!-- Module Status grids -->
          <ModuleStatusGrid
            :modules="filteredModuleStats"
            @select-module="handleSelectModule"
          />
        </template>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const isLoading = ref(true)
const rawGlobalStats = ref({ pass: 0, defect: 0, demoFail: 0, doing: 0 })
const rawModulesList = ref([])

const fetchDashboardStats = async () => {
  try {
    isLoading.value = true
    const config = useRuntimeConfig()
    const API_URL = `${config.public.apiBase.replace(/\/data$/, '')}/dashboard/stats`
    
    const data = await $fetch(API_URL)
    if (data) {
      if (data.global) rawGlobalStats.value = data.global
      if (data.modules) rawModulesList.value = data.modules
    }
  } catch (e) {
    console.error("Failed to load dashboard stats:", e)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchDashboardStats()
})

// Sidebar Menu items
const sidebarMenus = ref([
  { name: 'หน้าแดชบอร์ด', icon: '📊', active: true },
  { name: 'หุ้นกู้', icon: '📈', active: false },
  { name: 'พันธบัตร', icon: '📜', active: false },
  { name: 'ไม่อยู่ในความต้องการ', icon: '📉', active: false },
  { name: 'อยู่ในความต้องการ', icon: '📊', active: false },
  { name: 'ทุนเรือนหุ้น', icon: '🏢', active: false },
  { name: 'ปรับเปลี่ยนหุ้นรายเดือน', icon: '🔄', active: false },
  { name: 'เงินฝาก', icon: '💰', active: false },
  { name: 'ตั๋วสัญญาเงิน', icon: '📝', active: false },
  { name: 'จัดหาเงิน', icon: '🤝', active: false },
  { name: 'การเงิน', icon: '💵', active: false },
  { name: 'รายงาน', icon: '📑', active: false },
  { name: 'ตั้งค่า', icon: '⚙️', active: false }
])

// Global Statistics summary numbers
const globalStats = computed(() => [
  { label: 'สำเร็จ (Pass)', value: String(rawGlobalStats.value.pass), icon: '🟢', color: 'text-green-600' },
  { label: 'รอแก้ไข (Defect)', value: String(rawGlobalStats.value.defect), icon: '🔴', color: 'text-red-600' },
  { label: 'ไม่ผ่าน (Demo Fail)', value: String(rawGlobalStats.value.demoFail), icon: '🟠', color: 'text-orange-600' },
  { label: 'กำลังดำเนินการ (Doing)', value: String(rawGlobalStats.value.doing), icon: '🔵', color: 'text-blue-600' }
])

// Raw Module Statistics data
const moduleStats = computed(() => {
  const nameMapping = {
    "คำร้องขอปรับเปลี่ยนหุ้นรายเดือน": "ปรับเปลี่ยนหุ้นรายเดือน",
    "ตั๋วสัญญาใช้เงิน": "ตั๋วสัญญาเงิน"
  }
  const iconsMapping = {
    "หน้าแดชบอร์ด": "📊",
    "หุ้นกู้": "📈",
    "พันธบัตร": "📜",
    "ไม่อยู่ในความต้องการ": "📉",
    "อยู่ในความต้องการ": "📊",
    "ทุนเรือนหุ้น": "🏢",
    "ปรับเปลี่ยนหุ้นรายเดือน": "🔄",
    "เงินฝาก": "💰",
    "ตั๋วสัญญาเงิน": "📝",
    "จัดหาเงิน": "🤝",
    "การเงิน": "💵",
    "รายงาน": "📑",
    "ตั้งค่า": "⚙️"
  }

  const sidebarOrder = [
    'หุ้นกู้', 'พันธบัตร', 'ตั้งค่า', 'รายงาน', 'อยู่ในความต้องการ', 
    'ไม่อยู่ในความต้องการ', 'เงินฝาก', 'ตั๋วสัญญาเงิน', 'ทุนเรือนหุ้น', 
    'ปรับเปลี่ยนหุ้นรายเดือน', 'จัดหาเงิน', 'การเงิน'
  ]

  const mapped = rawModulesList.value.map(m => {
    const displayName = nameMapping[m.name] || m.name
    return {
      name: displayName,
      icon: iconsMapping[displayName] || '›',
      defect: m.defect,
      demoFail: m.demoFail,
      pass: m.pass,
      doing: m.doing
    }
  }).filter(m => m.name !== 'หน้าแดชบอร์ด')

  mapped.sort((a, b) => {
    const idxA = sidebarOrder.indexOf(a.name)
    const idxB = sidebarOrder.indexOf(b.name)
    return (idxA !== -1 ? idxA : 999) - (idxB !== -1 ? idxB : 999)
  })

  return mapped
})

// Search state
const searchQuery = ref('')

// Filter module stats based on search query
const filteredModuleStats = computed(() => {
  if (!searchQuery.value) return moduleStats.value
  const query = searchQuery.value.toLowerCase().trim()
  return moduleStats.value.filter(m => m.name.toLowerCase().includes(query))
})

// Event Handlers
const handleSelectMenu = (index) => {
  sidebarMenus.value.forEach((menu, idx) => {
    menu.active = idx === index
  })
}

const handleAddBug = () => {
  const router = useRouter()
  router.push('/dashboard-bugs')
}

const handleNotifications = () => {
  // Silent/No-op
}

const handleProfile = () => {
  // Silent/No-op
}

const handleSelectModule = (module) => {
  const router = useRouter()
  if (module.name === 'หุ้นกู้') router.push('/debentures')
  else if (module.name === 'พันธบัตร') router.push('/bonds')
  else if (module.name === 'ตั้งค่า') router.push('/settings/system')
  else if (module.name === 'รายงาน') router.push('/reports/debt')
  else if (module.name === 'อยู่ในความต้องการ') router.push('/equity-in-demand')
  else if (module.name === 'ไม่อยู่ในความต้องการ') router.push('/equity-not-in-demand')
  else if (module.name === 'เงินฝาก') router.push('/deposits')
  else if (module.name === 'ตั๋วสัญญาเงิน') router.push('/procurement/promissory-notes')
  else if (module.name === 'ทุนเรือนหุ้น') router.push('/finance/share-capital')
  else if (module.name === 'ปรับเปลี่ยนหุ้นรายเดือน') router.push('/finance/monthly-shares')
  else if (module.name === 'จัดหาเงิน') router.push('/procurement/loans/banks')
  else if (module.name === 'การเงิน') router.push('/finance')
}
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #f1f5f9;
  border-radius: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #e2e8f0;
}
</style>
