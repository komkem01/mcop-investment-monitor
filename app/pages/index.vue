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
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

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
const globalStats = ref([
  { label: 'สำเร็จ (Pass)', value: '127', icon: '🟢', color: 'text-green-600' },
  { label: 'รอแก้ไข (Defect)', value: '57', icon: '🔴', color: 'text-red-600' },
  { label: 'ไม่ผ่าน (Demo Fail)', value: '8', icon: '🟠', color: 'text-orange-600' },
  { label: 'กำลังดำเนินการ (Doing)', value: '0', icon: '🔵', color: 'text-blue-600' }
])

// Raw Module Statistics data
const moduleStats = ref([
  { name: 'หุ้นกู้', icon: '📈', defect: 21, demoFail: 3, pass: 59 },
  { name: 'พันธบัตร', icon: '📜', defect: 17, demoFail: 0, pass: 16 },
  { name: 'ตั้งค่า', icon: '⚙️', defect: 10, demoFail: 4, pass: 29 },
  { name: 'รายงาน', icon: '📑', defect: 10, demoFail: 0, pass: 0 },
  { name: 'อยู่ในความต้องการ', icon: '📊', defect: 3, demoFail: 0, pass: 12 },
  { name: 'ไม่อยู่ในความต้องการ', icon: '📉', defect: 1, demoFail: 0, pass: 6 },
  { name: 'เงินฝาก', icon: '💰', defect: 0, demoFail: 4, pass: 1 },
  { name: 'ตั๋วสัญญาเงิน', icon: '📝', defect: 2, demoFail: 1, pass: 4 },
  { name: 'ทุนเรือนหุ้น', icon: '🏢', defect: 1, demoFail: 0, pass: 0 },
  { name: 'ปรับเปลี่ยนหุ้นรายเดือน', icon: '🔄', defect: 1, demoFail: 0, pass: 1 },
  { name: 'จัดหาเงิน', icon: '🤝', defect: 3, demoFail: 0, pass: 0 },
  { name: 'การเงิน', icon: '💵', defect: 2, demoFail: 0, pass: 0 }
])

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
