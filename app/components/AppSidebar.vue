<template>
  <aside class="w-64 bg-[#0f172a] text-slate-300 flex flex-col h-screen shrink-0 border-r border-slate-800">
    <!-- Brand Title -->
    <div class="h-16 flex items-center px-6 border-b border-slate-800">
      <div class="flex items-center space-x-2.5">
        <!-- Glowing indigo bar -->
        <span class="w-[3px] h-4 bg-indigo-500 rounded-full shadow-[0_0_8px_#6366f1]"></span>
        <h1 class="text-sm font-extrabold text-white tracking-wide font-sans select-none">MCOP Investment Monitor</h1>
      </div>
    </div>
    
    <!-- Navigation Menus -->
    <nav class="flex-1 overflow-y-auto py-5 px-3 custom-scrollbar">
      <ul class="space-y-1.5">
        <!-- Root Menu List -->
        <li v-for="menu in sidebarMenus" :key="menu.name" class="space-y-1">
          
          <!-- Category Header (If has children, acts as Accordion Trigger) -->
          <div v-if="menu.children">
            <button
              @click="isCategoryEnabled(menu) ? toggleCategory(menu.name) : null"
              class="w-full flex items-center justify-between px-4 py-2.5 rounded-2xl text-xs font-bold transition-all duration-200 select-none text-slate-300"
              :class="isCategoryEnabled(menu)
                ? 'hover:text-white hover:bg-slate-800/60'
                : 'opacity-30 cursor-not-allowed pointer-events-none'"
            >
              <div class="flex items-center">
                <span class="mr-2.5 text-xs text-slate-500 font-bold select-none">{{ menu.icon }}</span>
                <span>{{ menu.name }}</span>
              </div>
              
              <!-- Accordion Arrow Icon -->
              <span 
                class="text-[9px] text-slate-500 transition-transform duration-200"
                :class="isExpanded(menu.name) ? 'rotate-90 text-slate-300' : ''"
                v-if="isCategoryEnabled(menu)"
              >
                ▶
              </span>
            </button>

            <!-- Collapsible Children Container -->
            <transition name="expand">
              <ul v-if="isExpanded(menu.name)" class="mt-1 ml-4 pl-3 border-l border-slate-800 space-y-1 overflow-hidden">
                <li v-for="child in menu.children" :key="child.name">
                  <NuxtLink
                    :to="child.to"
                    class="flex items-center px-3 py-2.5 rounded-xl text-[11px] font-semibold transition-all duration-150 select-none"
                    :class="[
                      !isMenuEnabled(child.to)
                        ? 'opacity-30 cursor-not-allowed pointer-events-none'
                        : '',
                      isActiveRoute(child.to)
                        ? 'bg-indigo-600 text-white font-bold shadow-md shadow-indigo-900/30'
                        : 'text-slate-400 hover:text-white hover:bg-slate-800/40'
                    ]"
                  >
                    <!-- Chevron › indicator -->
                    <span 
                      class="mr-2 text-xs font-bold transition-colors select-none"
                      :class="isActiveRoute(child.to) ? 'text-white' : 'text-slate-600'"
                    >›</span>
                    <span class="truncate">{{ child.name }}</span>
                  </NuxtLink>
                </li>
              </ul>
            </transition>
          </div>

          <!-- Direct Link (No children) -->
          <NuxtLink
            v-else
            :to="menu.to"
            class="flex items-center px-4 py-2.5 rounded-2xl text-xs font-bold transition-all duration-200 select-none"
            :class="[
              !isMenuEnabled(menu.to)
                ? 'opacity-30 cursor-not-allowed pointer-events-none'
                : '',
              isActiveRoute(menu.to)
                ? 'bg-indigo-600 text-white font-bold shadow-md shadow-indigo-900/30'
                : 'text-slate-300 hover:text-white hover:bg-slate-800/60'
            ]"
          >
            <span class="mr-2.5 text-xs text-slate-500 font-bold select-none" :class="isActiveRoute(menu.to) ? 'text-white' : ''">{{ menu.icon }}</span>
            {{ menu.name }}
          </NuxtLink>

        </li>
      </ul>
    </nav>
  </aside>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRoute } from '#app'

const route = useRoute()

// Hierarchical Sidebar Menu Data using › chevron symbol to match UI mockup
// (Declared first at the top of script to avoid TDZ / Before Initialization error on immediate watch)
const sidebarMenus = [
  {
    name: 'Dashboard',
    icon: '›',
    children: [
      { name: 'ภาพรวมระบบ', to: '/' },
      { name: 'บั๊กหน้าแดชบอร์ด', to: '/dashboard-bugs' }
    ]
  },
  {
    name: 'เงินลงทุน',
    icon: '›',
    children: [
      { name: 'ตราสารหนี้: หุ้นกู้', to: '/debentures' },
      { name: 'ตราสารหนี้: พันธบัตร', to: '/bonds' },
      { name: 'ตราสารทุน: หุ้นสามัญ ที่อยู่ในความต้องการของตลาด', to: '/equity-in-demand' },
      { name: 'ตราสารทุน: หุ้นสามัญ ที่ไม่อยู่ในความต้องการของตลาด', to: '/equity-not-in-demand' },
      { name: 'เงินฝาก', to: '/deposits' },
      { name: 'กองทุน (Private fund)', to: '/private-fund' },
      { name: 'รายงาน: ตราสารหนี้', to: '/reports/debt' },
      { name: 'รายงาน: ตราสารทุน', to: '/reports/equity' },
      { name: 'รายงาน: เงินฝาก', to: '/reports/deposits' },
      { name: 'รายงาน: 🚫ตั๋วสัญญาใช้เงิน', to: '/reports/promissory-notes' },
      { name: 'รายงาน: กองทุน', to: '/reports/funds' }
    ]
  },
  {
    name: 'จัดหาเงิน',
    icon: '›',
    children: [
      { name: 'กู้ยืมเงิน: ธนาคาร', to: '/procurement/loans/banks' },
      { name: 'กู้ยืมเงิน: สหกรณ์', to: '/procurement/loans/cooperatives' },
      { name: 'กู้ยืมเงิน: บริษัทหลักทรัพย์', to: '/procurement/loans/securities' },
      { name: 'กู้ยืมเงิน: ตั๋วสัญญาใช้เงิน (ขาย)', to: '/procurement/loans/promissory-notes' },
      { name: 'กู้ยืมเงิน: รายงาน', to: '/procurement/loans/reports' },
      { name: 'ตั๋วสัญญาใช้เงิน', to: '/procurement/promissory-notes' }
    ]
  },
  {
    name: 'การเงิน',
    icon: '›',
    children: [
      { name: 'การเงินหลัก', to: '/finance' },
      { name: 'ทุนเรือนหุ้น', to: '/finance/share-capital' },
      { name: 'ปรับเปลี่ยนหุ้นรายเดือน', to: '/finance/monthly-shares' }
    ]
  },
  {
    name: 'ตั้งค่า',
    icon: '›',
    children: [
      { name: 'ตั้งค่าระบบ', to: '/settings/system' },
      { name: 'ตั้งค่าผู้ใช้งาน', to: '/settings/users' },
      { name: 'ประวัติการทำงาน', to: '/settings/logs' }
    ]
  }
]

// Categories expansion states (Only one category can be open at a time)
const expandedCategories = ref({
  'Dashboard': false,
  'เงินลงทุน': false,
  'จัดหาเงิน': false,
  'การเงิน': false,
  'ตั้งค่า': false
})

// Function to auto-expand category containing active route and collapse others
const autoExpandActiveCategory = () => {
  for (const menu of sidebarMenus) {
    if (menu.children) {
      const hasActiveChild = menu.children.some(c => c.to === route.path)
      if (hasActiveChild) {
        // Collapse all others
        for (const key in expandedCategories.value) {
          expandedCategories.value[key] = false
        }
        // Expand the active category
        expandedCategories.value[menu.name] = true
        break
      }
    }
  }
}

onMounted(() => {
  autoExpandActiveCategory()
})

watch(() => route.path, () => {
  autoExpandActiveCategory()
}, { immediate: true })

const toggleCategory = (name) => {
  const currentVal = expandedCategories.value[name]
  // Collapse all categories
  for (const key in expandedCategories.value) {
    expandedCategories.value[key] = false
  }
  // Toggle the clicked one
  expandedCategories.value[name] = !currentVal
}

const isExpanded = (name) => {
  return expandedCategories.value[name]
}

// Active Route checking logic
const isActiveRoute = (path) => {
  return route.path === path
}

// Function to check if a specific menu item is enabled
const isMenuEnabled = (to) => {
  return true
}

// Function to check if a category header is enabled (has at least one enabled child or is enabled itself)
const isCategoryEnabled = (menu) => {
  return true
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
  background: #1e293b;
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #334155;
}

/* Slide/Expand Animation */
.expand-enter-active,
.expand-leave-active {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  max-height: 500px;
}

.expand-enter-from,
.expand-leave-to {
  max-height: 0;
  opacity: 0;
}
</style>
