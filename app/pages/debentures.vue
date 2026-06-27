<template>
  <div class="flex h-screen bg-transparent font-sans text-slate-800">
    <AppSidebar />

    <main class="flex-1 flex flex-col overflow-hidden">
      <div class="flex-1 overflow-y-auto p-8 space-y-8 custom-scrollbar">
        <div
          class="flex items-center justify-between pb-4 border-b border-slate-200"
        >
          <div>
            <span
              class="text-[10px] font-extrabold text-slate-400 uppercase tracking-widest block mb-1"
              >ตราสารหนี้</span
            >
            <h2
              class="text-2xl font-black text-slate-800 tracking-tight leading-none"
            >
              หุ้นกู้ (Debentures)
            </h2>
          </div>
          <button
            @click="openAddModal"
            class="bg-gradient-to-r from-indigo-600 to-indigo-500 hover:opacity-95 active:scale-95 px-5 py-2.5 rounded-full text-white text-xs font-extrabold shadow-[0_4px_14px_rgba(99,102,241,0.2)] transition-all flex items-center space-x-1.5 select-none"
          >
            <span>➕ แจ้งบันทึกบั๊กใหม่</span>
          </button>
        </div>

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
          <section class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
            <div
              v-for="(stat, index) in bugSummaryStats"
              :key="index"
              class="relative bg-white p-6 rounded-[22px] border border-slate-200/80 shadow-[0_4px_18px_rgba(0,0,0,0.015)] hover:shadow-[0_8px_24px_rgba(99,102,241,0.05)] hover:-translate-y-0.5 transition-all duration-300 flex flex-col justify-between overflow-hidden"
            >
              <span
                class="absolute top-0 left-0 right-0 h-[3px] bg-gradient-to-r"
                :class="stat.gradient"
              ></span>

              <div class="flex items-center justify-between mb-3">
                <p
                  class="text-xs font-bold text-slate-500 uppercase tracking-wider"
                >
                  {{ stat.label }}
                </p>
                <span class="text-base select-none">{{ stat.icon }}</span>
              </div>
              <div>
                <h3
                  class="text-2xl font-black text-slate-900 tracking-tight leading-none"
                >
                  {{ stat.value }}
                </h3>
              </div>
            </div>
          </section>

          <section
            class="bg-white rounded-[22px] border border-slate-200/80 shadow-[0_4px_18_rgba(0,0,0,0.015)] p-6 space-y-6"
          >
            <div
              class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 border-b border-slate-200 pb-5"
            >
              <h3 class="text-sm font-bold text-slate-800">
                รายการบั๊กทั้งหมด (Defect Tracking List)
              </h3>

              <div class="flex items-center gap-3 w-full sm:w-auto">
                <div
                  class="flex items-center bg-slate-50 border border-slate-200 rounded-full px-3.5 py-1.5 focus-within:bg-white focus-within:border-indigo-300 focus-within:shadow-[0_2px_10px_rgba(99,102,241,0.04)] transition-all duration-200 w-full sm:w-64"
                >
                  <span class="text-slate-500 text-xs mr-2">🔍</span>
                  <input
                    type="text"
                    v-model="searchQuery"
                    placeholder="ค้นหารหัสบั๊ก, ความคิดเห็น..."
                    class="bg-transparent border-none outline-none w-full text-xs text-slate-700 placeholder-slate-400 font-bold"
                  />
                </div>

                <div
                  class="relative min-w-[170px] shrink-0 custom-select-container"
                >
                  <button
                    type="button"
                    @click="filterDropdownOpen = !filterDropdownOpen"
                    class="w-full bg-slate-50 border border-slate-200 hover:bg-slate-100 hover:border-slate-300 rounded-full pl-4 pr-8 py-2 text-xs text-slate-700 outline-none font-bold cursor-pointer transition-all flex items-center justify-between select-none shadow-[0_1px_4px_rgba(0,0,0,0.015)]"
                  >
                    <span class="flex items-center">
                      <span
                        v-if="statusFilter"
                        class="px-2.5 py-0.5 rounded-full text-[9px] font-black mr-1.5 uppercase tracking-wider"
                        :class="statusClass(statusFilter)"
                      >
                        {{ statusFilter }}
                      </span>
                      <span v-else>สถานะทั้งหมด</span>
                    </span>
                    <span
                      class="text-[8px] text-slate-500 transition-transform duration-200"
                      :class="filterDropdownOpen ? 'rotate-180' : ''"
                      >▼</span
                    >
                  </button>

                  <div
                    v-if="filterDropdownOpen"
                    class="absolute right-0 top-full mt-2 w-full min-w-[190px] bg-white border border-slate-200 shadow-[0_10px_30px_rgba(0,0,0,0.08)] rounded-2xl py-2 z-50 animate-fade-in"
                  >
                    <button
                      @click="
                        statusFilter = '';
                        filterDropdownOpen = false;
                      "
                      class="w-full text-left px-4 py-2 text-xs font-bold text-slate-500 hover:bg-indigo-50/50 hover:text-indigo-655 transition-all"
                    >
                      สถานะทั้งหมด
                    </button>
                    <div
                      class="max-h-[220px] overflow-y-auto custom-scrollbar space-y-1 px-2 mt-1"
                    >
                      <button
                        v-for="st in [
                          'Defect',
                          'Demo',
                          'Demo Fail',
                          'Pass',
                          'Ready For Demo',
                          'Doing',
                        ]"
                        :key="st"
                        @click="
                          statusFilter = st;
                          filterDropdownOpen = false;
                        "
                        class="w-full text-left px-2 py-1.5 text-xs font-bold text-slate-700 hover:bg-slate-50 rounded-xl transition-all flex items-center"
                      >
                        <span
                          class="px-2.5 py-0.5 rounded-full text-[9px] font-black uppercase tracking-wider"
                          :class="statusClass(st)"
                        >
                          {{ st }}
                        </span>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="overflow-x-auto border border-slate-200/80 rounded-2xl">
              <table class="w-full text-left border-collapse min-w-[1000px]">
                <thead>
                  <tr
                    class="bg-slate-100/90 border-b border-slate-200 text-[10px] text-slate-700 font-extrabold uppercase tracking-wider"
                  >
                    <th class="py-3.5 pl-4 w-24">ID</th>
                    <th class="py-3.5 w-40">รายละเอียด</th>
                    <th class="py-3.5 w-64">ความคิดเห็น (Checklist)</th>
                    <th class="py-3.5 w-36">ความคาดหวัง</th>
                    <th class="py-3.5 w-28 text-center">ภาพประกอบ</th>
                    <th class="py-3.5 w-24 text-center">Figma</th>
                    <th class="py-3.5 w-28 text-center">สถานะ</th>
                    <th class="py-3.5 w-32">วันที่สร้าง</th>
                    <th class="py-3.5 w-32">วันที่แก้ไข</th>
                    <th class="py-3.5 w-32">วันที่ผ่าน</th>
                    <th class="py-3.5 text-center pr-4 w-40">การจัดการ</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-200">
                  <tr
                    v-for="item in paginatedBugs"
                    :key="item.id"
                    class="text-xs transition-all hover:bg-slate-50/70 border-b border-slate-100 group"
                  >
                    <td class="py-4 pl-4 font-black text-slate-900 align-top">
                      {{ item.bugId }}
                    </td>

                    <td class="py-4 font-bold text-slate-800 align-top">
                      {{ item.title }}
                    </td>

                    <td class="py-4 text-slate-800 align-top pr-4">
                      <ul class="space-y-1">
                        <li
                          v-for="(comment, cIdx) in item.comments"
                          :key="cIdx"
                          class="flex items-start space-x-1.5"
                        >
                          <span
                            class="text-[9px] mt-0.5 select-none text-indigo-500 font-bold"
                            >—</span
                          >
                          <span class="leading-relaxed font-semibold">{{
                            comment
                          }}</span>
                        </li>
                      </ul>
                    </td>

                    <td class="py-4 text-slate-700 font-bold align-top">
                      {{ item.expectation || "" }}
                    </td>

                    <td class="py-4 align-top text-center">
                      <a
                        v-if="item.driveLink"
                        :href="item.driveLink"
                        target="_blank"
                        class="text-indigo-700 hover:text-indigo-600 font-extrabold hover:underline inline-flex items-center space-x-0.5"
                      >
                        <span>📎 LINK</span>
                      </a>
                      <span v-else></span>
                    </td>

                    <td class="py-4 align-top text-center">
                      <a
                        v-if="item.figmaLink"
                        :href="item.figmaLink"
                        target="_blank"
                        class="text-indigo-700 hover:text-indigo-600 font-extrabold hover:underline"
                      >
                        LINK
                      </a>
                      <span v-else></span>
                    </td>

                    <td class="py-4 align-top text-center">
                      <span
                        class="px-2.5 py-0.5 rounded-full text-[9px] font-black inline-block uppercase tracking-wider"
                        :class="statusClass(item.status)"
                      >
                        {{ item.status }}
                      </span>
                    </td>

                    <td class="py-4 align-top text-slate-700 font-bold">
                      {{ item.createdDate || "" }}
                    </td>

                    <td class="py-4 align-top text-slate-700 font-bold">
                      {{ item.fixedDate || "" }}
                    </td>

                    <td class="py-4 align-top text-slate-700 font-bold">
                      {{ item.passedDate || "" }}
                    </td>

                    <td
                      class="py-4 align-top text-center pr-4 space-x-1.5 whitespace-nowrap"
                    >
                      <button
                        @click="openViewModal(item)"
                        class="bg-indigo-50 hover:bg-indigo-100 text-indigo-600 px-2.5 py-1.5 rounded-full text-[10px] font-extrabold transition-all duration-150 active:scale-95 shadow-[0_1px_4px_rgba(99,102,241,0.04)]"
                      >
                        ดูรายละเอียด
                      </button>
                      <button
                        @click="openEditModal(item)"
                        class="bg-amber-50 hover:bg-amber-100 text-amber-600 px-2.5 py-1.5 rounded-full text-[10px] font-extrabold transition-all duration-150 active:scale-95 shadow-[0_1px_4px_rgba(217,119,6,0.04)]"
                      >
                        แก้ไข
                      </button>
                      <button
                        @click="confirmDelete(item)"
                        class="bg-rose-50 hover:bg-rose-100 text-rose-600 px-2.5 py-1.5 rounded-full text-[10px] font-extrabold transition-all duration-150 active:scale-95 shadow-[0_1px_4px_rgba(225,29,72,0.04)]"
                      >
                        ลบ
                      </button>
                    </td>
                  </tr>
                  <tr v-if="paginatedBugs.length === 0">
                    <td
                      colspan="11"
                      class="py-8 text-center text-slate-400 font-medium"
                    >
                      ไม่พบรายการบั๊กที่ค้นหา
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Pagination Bar -->
            <div
              class="flex items-center justify-between mt-4 px-2"
              v-if="filteredBugs.length > 0"
            >
              <div
                class="text-[10px] font-extrabold text-slate-500 uppercase tracking-wide"
              >
                แสดง {{ (currentPage - 1) * itemsPerPage + 1 }} ถึง
                {{
                  Math.min(currentPage * itemsPerPage, filteredBugs.length)
                }}
                จากทั้งหมด {{ filteredBugs.length }} รายการ
              </div>
              <div class="flex items-center space-x-1">
                <button
                  @click="prevPage"
                  :disabled="currentPage === 1"
                  class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all"
                  :class="
                    currentPage === 1
                      ? 'text-slate-300 cursor-not-allowed'
                      : 'text-slate-600 hover:bg-slate-100 active:scale-95'
                  "
                >
                  ก่อนหน้า
                </button>

                <div class="flex space-x-1">
                  <button
                    v-for="page in visiblePages"
                    :key="page"
                    @click="page !== '...' ? goToPage(page) : null"
                    :disabled="page === '...'"
                    class="min-w-[32px] h-8 rounded-xl text-xs font-extrabold transition-all flex items-center justify-center px-2"
                    :class="[
                      currentPage === page
                        ? 'bg-indigo-600 text-white shadow-[0_4px_10px_rgba(99,102,241,0.3)]'
                        : 'text-slate-600 hover:bg-slate-100 active:scale-95',
                      page === '...'
                        ? 'cursor-default hover:bg-transparent text-slate-400'
                        : '',
                    ]"
                  >
                    {{ page }}
                  </button>
                </div>

                <button
                  @click="nextPage"
                  :disabled="currentPage === totalPages"
                  class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all"
                  :class="
                    currentPage === totalPages
                      ? 'text-slate-300 cursor-not-allowed'
                      : 'text-slate-600 hover:bg-slate-100 active:scale-95'
                  "
                >
                  ถัดไป
                </button>
              </div>
            </div>
          </section>
        </template>
      </div>
    </main>

    <!-- Loading Popup Overlay ขณะยิง API -->
    <div
      v-if="isProcessing"
      class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/40 backdrop-blur-md p-4 animate-fade-in"
    >
      <div
        class="bg-white/95 w-full max-w-sm rounded-[32px] border border-white/80 shadow-[0_24px_50px_rgba(99,102,241,0.06)] p-8 flex flex-col items-center justify-center space-y-6 text-center"
      >
        <div class="relative flex items-center justify-center">
          <div
            class="w-16 h-16 rounded-full border-2 border-indigo-500/10 border-t-indigo-500 animate-spin shadow-[0_0_15px_rgba(99,102,241,0.15)]"
          ></div>
          <div
            class="absolute w-8 h-8 rounded-full bg-gradient-to-br from-indigo-500 to-sky-400 animate-pulse shadow-[0_0_20px_rgba(99,102,241,0.5)] opacity-80"
          ></div>
        </div>
        <div class="space-y-1.5">
          <h3 class="text-sm font-extrabold text-slate-800 tracking-wide">
            {{ processingMessage }}
          </h3>
          <p class="text-[10px] font-bold text-slate-400 tracking-wider">
            กรุณารอสักครู่ ระบบกำลังเชื่อมต่อกับเซิร์ฟเวอร์...
          </p>
        </div>
      </div>
    </div>

    <div
      v-if="showFormModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/40 backdrop-blur-md p-4 animate-fade-in"
    >
      <div
        class="bg-white/95 w-full max-w-3xl rounded-[32px] border border-white/80 shadow-[0_24px_50px_rgba(99,102,241,0.06)] flex flex-col max-h-[90vh] overflow-hidden"
      >
        <div
          class="flex items-center justify-between px-6 py-5 border-b border-slate-100 bg-gradient-to-r from-indigo-50/40 via-white to-transparent shrink-0"
        >
          <div class="flex items-center space-x-2">
            <span class="w-[3px] h-4 bg-indigo-500 rounded-full"></span>
            <h3
              class="text-sm font-extrabold text-slate-800 tracking-wide font-sans"
            >
              {{
                formType === "add"
                  ? "แจ้งบันทึกบั๊กใหม่"
                  : `แก้ไขรายการบั๊ก ${formState.bugId}`
              }}
            </h3>
          </div>
          <button
            @click="closeFormModal"
            class="w-7 h-7 rounded-full bg-slate-50 border border-slate-100 hover:bg-rose-50 hover:border-rose-100 hover:text-rose-500 text-slate-400 font-bold transition-all duration-200 flex items-center justify-center text-[10px]"
          >
            ✕
          </button>
        </div>

        <div
          class="flex-1 overflow-y-auto px-6 py-5 space-y-4.5 custom-scrollbar text-xs"
        >
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label class="font-bold text-slate-500 tracking-wider"
                >รหัสบั๊ก (BUG ID) *</label
              >
              <input
                type="text"
                v-model="formState.bugId"
                readonly
                class="w-full bg-slate-100/70 border border-slate-200 rounded-2xl px-4 py-3 outline-none transition-all font-bold text-slate-500 cursor-not-allowed select-none"
              />
            </div>
            <div class="space-y-1.5">
              <label class="font-bold text-slate-500 tracking-wider"
                >รายละเอียดหัวข้อ(เส้นทาง) *</label
              >
              <input
                type="text"
                v-model="formState.title"
                placeholder="ตัวอย่าง ปรับให้ตรงตาม Figma การคำนวณอัตราดอกเบี้ยสุทธิ"
                class="w-full bg-slate-50/50 border border-slate-200 rounded-2xl px-4 py-3 outline-none focus:border-indigo-400 focus:bg-white focus:shadow-[0_4px_16_rgba(99,102,241,0.03)] transition-all font-semibold text-slate-700"
              />
            </div>
          </div>

          <div class="space-y-1.5">
            <label
              class="font-bold text-slate-500 tracking-wider flex items-center justify-between"
            >
              <span>รายการแก้ไข / ความคิดเห็น (CHECKLIST) *</span>
              <span class="text-[9px] font-medium text-slate-400"
                >แยกบรรทัดเพื่อเพิ่มข้อใหม่</span
              >
            </label>
            <textarea
              v-model="formCommentsRaw"
              placeholder="ตัวอย่าง ปรับให้ตรงตาม Figma:&#10;- สูตรคำนวณหักภาษี ณ ที่จ่ายคิดผิดไป 0.01 บาท&#10;- ตรวจสอบทศนิยมให้อยู่ในรูปแบบ 2 ตำแหน่ง"
              rows="4"
              class="w-full bg-slate-50/50 border border-slate-200 rounded-2xl px-4 py-3 outline-none focus:border-indigo-400 focus:bg-white focus:shadow-[0_4px_16_rgba(99,102,241,0.03)] transition-all font-semibold text-slate-700 resize-none leading-relaxed"
            ></textarea>
          </div>

          <div class="space-y-1.5">
            <label class="font-bold text-slate-500 tracking-wider"
              >ความคาดหวัง (EXPECTATION)</label
            >
            <input
              type="text"
              v-model="formState.expectation"
              placeholder="ตัวอย่าง ปรับให้ตรงตาม Figma "
              class="w-full bg-slate-50/50 border border-slate-200 rounded-2xl px-4 py-3 outline-none focus:border-indigo-400 focus:bg-white focus:shadow-[0_4px_16_rgba(99,102,241,0.03)] transition-all font-semibold text-slate-700"
            />
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label class="font-bold text-slate-500 tracking-wider"
                >ภาพประกอบ (DRIVE LINK)</label
              >
              <input
                type="text"
                v-model="formState.driveLink"
                placeholder="https://drive.google.com/..."
                class="w-full bg-slate-50/50 border border-slate-200 rounded-2xl px-4 py-3 outline-none focus:border-indigo-400 focus:bg-white focus:shadow-[0_4px_16_rgba(99,102,241,0.03)] transition-all font-semibold text-slate-700"
              />
            </div>
            <div class="space-y-1.5">
              <label class="font-bold text-slate-500 tracking-wider"
                >เอกสารเพิ่มเติม (FIGMA LINK)</label
              >
              <input
                type="text"
                v-model="formState.figmaLink"
                placeholder="https://figma.com/..."
                class="w-full bg-slate-50/50 border border-slate-200 rounded-2xl px-4 py-3 outline-none focus:border-indigo-400 focus:bg-white focus:shadow-[0_4px_16_rgba(99,102,241,0.03)] transition-all font-semibold text-slate-700"
              />
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div class="space-y-1.5">
              <label class="font-bold text-slate-500 tracking-wider"
                >วันที่สร้าง</label
              >
              <input
                type="text"
                v-model="formState.createdDate"
                placeholder="28/05/2026 11:45"
                class="w-full bg-slate-50/50 border border-slate-200 rounded-2xl px-4 py-3 outline-none focus:border-indigo-400 focus:bg-white focus:shadow-[0_4px_16_rgba(99,102,241,0.03)] transition-all font-semibold text-slate-700"
              />
            </div>
            <div class="space-y-1.5">
              <label class="font-bold text-slate-500 tracking-wider"
                >วันที่แก้ไข</label
              >
              <input
                type="text"
                v-model="formState.fixedDate"
                placeholder="29/05/2026 15:00"
                class="w-full bg-slate-50/50 border border-slate-200 rounded-2xl px-4 py-3 outline-none focus:border-indigo-400 focus:bg-white focus:shadow-[0_4px_16_rgba(99,102,241,0.03)] transition-all font-semibold text-slate-700"
              />
            </div>
            <div class="space-y-1.5">
              <label class="font-bold text-slate-500 tracking-wider"
                >วันที่ผ่าน</label
              >
              <input
                type="text"
                v-model="formState.passedDate"
                placeholder="30/05/2026 09:00"
                class="w-full bg-slate-50/50 border border-slate-200 rounded-2xl px-4 py-3 outline-none focus:border-indigo-400 focus:bg-white focus:shadow-[0_4px_16_rgba(99,102,241,0.03)] transition-all font-semibold text-slate-700"
              />
            </div>
          </div>

          <div class="space-y-1.5 custom-select-container">
            <label class="font-bold text-slate-500 tracking-wider"
              >สถานะ *</label
            >
            <div class="relative">
              <button
                type="button"
                @click="formStatusDropdownOpen = !formStatusDropdownOpen"
                class="w-full bg-slate-50/50 border border-slate-200 rounded-2xl pl-4 pr-10 py-3 outline-none text-left font-bold text-slate-700 focus:border-indigo-400 focus:bg-white focus:shadow-[0_4px_16_rgba(99,102,241,0.03)] transition-all flex items-center justify-between select-none"
              >
                <span
                  class="px-2.5 py-0.5 rounded-full text-[9px] font-black uppercase tracking-wider"
                  :class="statusClass(formState.status)"
                >
                  {{ formState.status }}
                </span>
                <span
                  class="text-[8px] text-slate-400 pointer-events-none transition-transform duration-200"
                  :class="formStatusDropdownOpen ? 'rotate-180' : ''"
                  >▼</span
                >
              </button>

              <div
                v-if="formStatusDropdownOpen"
                class="absolute left-0 right-0 bottom-full mb-1.5 bg-white border border-slate-200 shadow-[0_-12px_36px_rgba(0,0,0,0.08)] rounded-2xl py-2.5 z-55 animate-fade-in"
              >
                <div
                  class="max-h-[220px] overflow-y-auto custom-scrollbar space-y-1.5 px-2"
                >
                  <button
                    v-for="st in [
                      'Defect',
                      'Demo',
                      'Demo Fail',
                      'Pass',
                      'Ready For Demo',
                      'Doing',
                    ]"
                    :key="st"
                    type="button"
                    @click="
                      formState.status = st;
                      formStatusDropdownOpen = false;
                    "
                    class="w-full text-left px-3 py-2 text-xs font-bold text-slate-700 hover:bg-slate-50 rounded-xl transition-all flex items-center"
                  >
                    <span
                      class="px-2.5 py-0.5 rounded-full text-[9px] font-black uppercase tracking-wider"
                      :class="statusClass(st)"
                    >
                      {{ st }}
                    </span>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div
          class="flex justify-end space-x-3 px-6 py-4.5 border-t border-slate-100 bg-slate-50/50 shrink-0"
        >
          <button
            @click="closeFormModal"
            class="bg-white border border-slate-200 hover:bg-slate-100 hover:border-slate-300 active:scale-95 px-6 py-2.5 rounded-full text-xs font-bold text-slate-600 transition-all select-none shadow-[0_2px_6px_rgba(0,0,0,0.02)]"
          >
            ยกเลิก
          </button>
          <button
            @click="saveBug"
            class="bg-gradient-to-r from-indigo-600 to-indigo-500 hover:opacity-95 active:scale-95 px-7 py-2.5 rounded-full text-white text-xs font-extrabold shadow-[0_4px_14_rgba(99,102,241,0.2)] transition-all select-none"
          >
            บันทึกข้อมูล
          </button>
        </div>
      </div>
    </div>

    <div
      v-if="showViewModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/40 backdrop-blur-md p-4 animate-fade-in"
    >
      <div
        class="bg-white w-full max-w-2xl rounded-[32px] border border-white/80 shadow-[0_24px_50px_rgba(99,102,241,0.06)] flex flex-col max-h-[85vh] overflow-hidden"
      >
        <div
          class="px-6 py-5 border-b shrink-0 flex items-center justify-between"
          :class="[
            selectedBug.status === 'Pass'
              ? 'bg-emerald-50/60 border-emerald-100/50'
              : selectedBug.status === 'Demo Fail'
                ? 'bg-amber-50/60 border-amber-100/50'
                : 'bg-rose-50/60 border-rose-100/50',
          ]"
        >
          <div class="flex items-center space-x-3">
            <span
              class="w-2.5 h-2.5 rounded-full animate-pulse"
              :class="[
                selectedBug.status === 'Pass'
                  ? 'bg-emerald-500'
                  : selectedBug.status === 'Demo Fail'
                    ? 'bg-amber-500'
                    : 'bg-rose-500',
              ]"
            ></span>
            <div>
              <span
                class="text-[9px] uppercase tracking-wider font-extrabold text-slate-400 block mb-0.5"
                >รายละเอียดข้อบกพร่อง</span
              >
              <h3
                class="text-sm font-extrabold text-slate-800 font-sans tracking-tight leading-none"
              >
                บั๊กรหัส {{ selectedBug.bugId }}
              </h3>
            </div>
          </div>
          <button
            @click="closeViewModal"
            class="w-7 h-7 rounded-full bg-white border border-slate-100 hover:bg-slate-200/50 hover:text-slate-700 text-slate-400 font-bold transition-all duration-200 flex items-center justify-center text-[10px] shadow-[0_2px_6px_rgba(0,0,0,0.015)]"
          >
            ✕
          </button>
        </div>

        <div
          class="flex-1 overflow-y-auto px-6 py-6 space-y-5 custom-scrollbar text-xs"
        >
          <div
            class="grid grid-cols-1 sm:grid-cols-2 gap-4 bg-slate-50/50 border border-slate-100/50 rounded-2xl p-4"
          >
            <div class="space-y-1">
              <span
                class="font-bold text-slate-400 text-[10px] uppercase tracking-wider"
                >รายละเอียดหัวข้อ:</span
              >
              <p class="font-bold text-slate-800 text-xs leading-relaxed">
                {{ selectedBug.title }}
              </p>
            </div>
            <div class="space-y-1">
              <span
                class="font-bold text-slate-400 text-[10px] uppercase tracking-wider"
                >ความคาดหวัง (Expectation):</span
              >
              <p class="font-semibold text-slate-700 text-xs leading-relaxed">
                {{ selectedBug.expectation || "-" }}
              </p>
            </div>
          </div>

          <div class="space-y-2">
            <span
              class="font-bold text-slate-500 text-[10px] uppercase tracking-wider block"
              >รายการข้อคิดเห็น / ความต้องการแก้ (Checklist)</span
            >
            <div
              class="bg-white border border-slate-200 rounded-2xl p-4 shadow-[0_4px_16_rgba(0,0,0,0.005)]"
            >
              <ul class="space-y-3">
                <li
                  v-for="(c, cIdx) in selectedBug.comments"
                  :key="cIdx"
                  class="flex items-start space-x-3"
                >
                  <span
                    class="w-4.5 h-4.5 rounded-full bg-indigo-50 border border-indigo-150 flex items-center justify-center shrink-0 text-indigo-500 select-none text-[9px] font-bold"
                    >✓</span
                  >
                  <span
                    class="leading-relaxed font-semibold text-slate-700 text-xs"
                    >{{ c }}</span
                  >
                </li>
              </ul>
            </div>
          </div>

          <div class="space-y-2">
            <span
              class="font-bold text-slate-500 text-[10px] uppercase tracking-wider block"
              >เอกสารแนบและแบบฟิกมา (Attachments)</span
            >
            <div class="grid grid-cols-2 gap-3.5">
              <a
                v-if="selectedBug.driveLink"
                :href="selectedBug.driveLink"
                target="_blank"
                class="flex items-center justify-center p-3 rounded-2xl bg-indigo-50/70 border border-indigo-100 text-indigo-600 font-extrabold hover:bg-indigo-100 hover:text-indigo-700 transition-all text-center space-x-1.5 shadow-[0_2px_8px_rgba(99,102,241,0.03)]"
              >
                <span>📎 ดูรูปภาพประกอบ</span>
              </a>
              <div
                v-else
                class="flex items-center justify-center p-3 rounded-2xl bg-slate-50 border border-slate-100 text-slate-400 font-bold text-center"
              >
                ไม่มีรูปภาพประกอบ
              </div>

              <a
                v-if="selectedBug.figmaLink"
                :href="selectedBug.figmaLink"
                target="_blank"
                class="flex items-center justify-center p-3 rounded-2xl bg-slate-50 border border-slate-200 text-slate-600 font-extrabold hover:bg-slate-100 hover:text-slate-800 transition-all text-center space-x-1.5 shadow-[0_2px_8px_rgba(0,0,0,0.015)]"
              >
                <span>🔗 เปิดดีไซน์ Figma</span>
              </a>
              <div
                v-else
                class="flex items-center justify-center p-3 rounded-2xl bg-slate-50 border border-slate-100 text-slate-400 font-bold text-center"
              >
                ไม่มีเอกสารอ้างอิง Figma
              </div>
            </div>
          </div>

          <div
            class="border-t border-slate-100 pt-4 grid grid-cols-3 gap-3 text-center"
          >
            <div
              class="bg-slate-50/50 rounded-xl p-2.5 border border-slate-100/50"
            >
              <span
                class="text-[9px] uppercase tracking-wider font-extrabold text-slate-400 block mb-0.5"
                >วันที่ตรวจพบ</span
              >
              <span class="text-slate-700 font-bold text-[10px]">{{
                selectedBug.createdDate || "-"
              }}</span>
            </div>
            <div
              class="bg-slate-50/50 rounded-xl p-2.5 border border-slate-100/50"
            >
              <span
                class="text-[9px] uppercase tracking-wider font-extrabold text-slate-400 block mb-0.5"
                >วันที่แก้ไข</span
              >
              <span class="text-slate-700 font-bold text-[10px]">{{
                selectedBug.fixedDate || "-"
              }}</span>
            </div>
            <div
              class="bg-slate-50/50 rounded-xl p-2.5 border border-slate-100/50"
            >
              <span
                class="text-[9px] uppercase tracking-wider font-extrabold text-slate-400 block mb-0.5"
                >วันที่ผ่านตรวจ</span
              >
              <span class="text-slate-700 font-bold text-[10px]">{{
                selectedBug.passedDate || "-"
              }}</span>
            </div>
          </div>
        </div>

        <div
          class="flex justify-end px-6 py-4 border-t border-slate-100 bg-slate-50/50 shrink-0"
        >
          <button
            @click="closeViewModal"
            class="bg-slate-900 hover:bg-slate-800 text-white px-6 py-2.5 rounded-full text-xs font-extrabold shadow-sm active:scale-95 transition-all select-none"
          >
            ปิด
          </button>
        </div>
      </div>
    </div>

    <div
      v-if="showDeleteModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/40 backdrop-blur-md p-4 animate-fade-in"
    >
      <div
        class="bg-white w-full max-w-md rounded-[32px] border border-white/80 shadow-[0_24px_50px_rgba(99,102,241,0.06)] p-6 text-center"
      >
        <div
          class="mx-auto w-16 h-16 rounded-full bg-rose-50 border border-rose-100 text-rose-500 flex items-center justify-center mb-4 text-2xl select-none animate-bounce"
        >
          ⚠️
        </div>

        <h3 class="text-sm font-extrabold text-slate-800 mb-2">
          ยืนยันการลบรายการบั๊ก?
        </h3>
        <p class="text-xs text-slate-500 font-bold leading-relaxed mb-6">
          คุณกำลังจะลบรายการบั๊ก
          <span class="text-rose-600 font-black">{{ bugToDelete?.bugId }}</span>
          ({{ bugToDelete?.title }}) <br />
          การดำเนินการนี้จะไม่สามารถกู้คืนข้อมูลได้
        </p>

        <div class="flex items-center justify-center space-x-3">
          <button
            @click="closeDeleteModal"
            class="bg-white border border-slate-200 hover:bg-slate-100 hover:border-slate-350 active:scale-95 px-6 py-2.5 rounded-full text-xs font-bold text-slate-600 transition-all select-none shadow-[0_2px_6px_rgba(0,0,0,0.02)]"
          >
            ยกเลิก
          </button>
          <button
            @click="executeDelete"
            class="bg-rose-600 hover:bg-rose-700 active:scale-95 px-7 py-2.5 rounded-full text-white text-xs font-extrabold shadow-[0_4px_14px_rgba(244,63,94,0.2)] transition-all select-none"
          >
            ยืนยันลบรายการ
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from "vue";
import Swal from "sweetalert2";
import { useBugs } from "~/composables/useBugs";

const { bugsList, isLoading, fetchBugs, addBug, updateBug, deleteBug } =
  useBugs("หุ้นกู้");

// SweetAlert Toast definition
let Toast = null;
onMounted(() => {
  fetchBugs();

  Toast = Swal.mixin({
    toast: true,
    position: "top-end",
    showConfirmButton: false,
    timer: 2500,
    timerProgressBar: true,
    customClass: {
      popup: "rounded-2xl border border-slate-100 shadow-lg text-xs",
    },
    didOpen: (toast) => {
      toast.addEventListener("mouseenter", Swal.stopTimer);
      toast.addEventListener("mouseleave", Swal.resumeTimer);
    },
  });
});

const triggerToast = (icon, title) => {
  if (Toast) {
    Toast.fire({ icon, title });
  }
};

// --- ฟังก์ชันสำหรับ Native Loading Popup ---
const isProcessing = ref(false);
const processingMessage = ref("กำลังประมวลผล...");

const showLoadingPopup = (titleText = "กำลังประมวลผล...") => {
  processingMessage.value = titleText;
  isProcessing.value = true;
};

const closeLoadingPopup = () => {
  isProcessing.value = false;
};

// --- จัดการ Dropdown ---
const filterDropdownOpen = ref(false);
const formStatusDropdownOpen = ref(false);

const closeAllDropdowns = (e) => {
  if (!e.target.closest(".custom-select-container")) {
    filterDropdownOpen.value = false;
    formStatusDropdownOpen.value = false;
  }
};

onMounted(() => {
  document.addEventListener("click", closeAllDropdowns);
});

onUnmounted(() => {
  document.removeEventListener("click", closeAllDropdowns);
});

// Summary statistics for Debentures module bugs
const bugSummaryStats = computed(() => {
  const total = bugsList.value.length;
  const pass = bugsList.value.filter((b) => b.status === "Pass").length;
  const defect = bugsList.value.filter((b) => b.status === "Defect").length;
  const demoFail = bugsList.value.filter(
    (b) => b.status === "Demo Fail",
  ).length;

  return [
    {
      label: "บั๊กทั้งหมด (Total)",
      value: `${total} รายการ`,
      icon: "📊",
      gradient: "from-sky-400 to-indigo-500",
    },
    {
      label: "ผ่านการทดสอบ (Pass)",
      value: `${pass} รายการ`,
      icon: "🟢",
      gradient: "from-emerald-400 to-teal-500",
    },
    {
      label: "รอแก้ไข (Defect)",
      value: `${defect} รายการ`,
      icon: "🔴",
      gradient: "from-rose-400 to-red-500",
    },
    {
      label: "ไม่ผ่าน (Demo Fail)",
      value: `${demoFail} รายการ`,
      icon: "🟠",
      gradient: "from-amber-400 to-orange-500",
    },
  ];
});

const searchQuery = ref("");
const statusFilter = ref("");

// --- ระบบแบ่งหน้า (Pagination) ---
const currentPage = ref(1);
const itemsPerPage = ref(10);

watch([searchQuery, statusFilter], () => {
  currentPage.value = 1;
});

const filteredBugs = computed(() => {
  return bugsList.value.filter((item) => {
    // Search filter
    const matchesSearch =
      item.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      item.bugId.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      item.comments.some((c) =>
        c.toLowerCase().includes(searchQuery.value.toLowerCase()),
      );

    // Status filter
    let matchesStatus = true;
    if (statusFilter.value) {
      matchesStatus = item.status === statusFilter.value;
    }

    return matchesSearch && matchesStatus;
  });
});

const totalPages = computed(() => {
  return Math.ceil(filteredBugs.value.length / itemsPerPage.value) || 1;
});

const paginatedBugs = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value;
  const end = start + itemsPerPage.value;
  return filteredBugs.value.slice(start, end);
});

const nextPage = () => {
  if (currentPage.value < totalPages.value) currentPage.value++;
};

const prevPage = () => {
  if (currentPage.value > 1) currentPage.value--;
};

const goToPage = (page) => {
  currentPage.value = page;
};

const visiblePages = computed(() => {
  const total = totalPages.value;
  const current = currentPage.value;
  const delta = 1;

  const range = [];
  const rangeWithDots = [];
  let l;

  for (let i = 1; i <= total; i++) {
    if (
      i === 1 ||
      i === total ||
      (i >= current - delta && i <= current + delta)
    ) {
      range.push(i);
    }
  }

  for (let i of range) {
    if (l) {
      if (i - l === 2) {
        rangeWithDots.push(l + 1);
      } else if (i - l !== 1) {
        rangeWithDots.push("...");
      }
    }
    rangeWithDots.push(i);
    l = i;
  }

  return rangeWithDots;
});

// Modal states
const showFormModal = ref(false);
const showViewModal = ref(false);
const showDeleteModal = ref(false);
const formType = ref("add"); // 'add' or 'edit'
const selectedBug = ref(null);
const bugToDelete = ref(null);

const formState = ref({
  id: null,
  sheetRowIndex: null,
  bugId: "",
  title: "",
  comments: [],
  expectation: "",
  driveLink: "",
  figmaLink: "",
  status: "Defect",
  createdDate: "",
  fixedDate: "",
  passedDate: "",
});

const formCommentsRaw = ref("");

// ฟังก์ชันสุ่ม/หารหัสบั๊กถัดไป
const generateNextBugId = () => {
  if (bugsList.value.length === 0) return "UBIV1";

  let maxNum = 0;

  bugsList.value.forEach((bug) => {
    const match = bug.bugId.match(/^UBIV(\d+)$/i);
    if (match) {
      const num = parseInt(match[1], 10);
      if (num > maxNum) {
        maxNum = num;
      }
    }
  });

  return `UBIV${maxNum + 1}`;
};

const openAddModal = () => {
  formType.value = "add";
  formCommentsRaw.value = "";
  formState.value = {
    id: null,
    sheetRowIndex: null,
    bugId: generateNextBugId(), // เรียกใช้ฟังก์ชันดึงเลขบั๊กอัตโนมัติที่นี่
    title: "",
    comments: [],
    expectation: "",
    driveLink: "",
    figmaLink: "",
    status: "Defect",
    createdDate: new Date().toLocaleString("th-TH"),
    fixedDate: "",
    passedDate: "",
  };
  showFormModal.value = true;
};

const openEditModal = (item) => {
  formType.value = "edit";
  formCommentsRaw.value = item.comments.join("\n");
  formState.value = { ...item };
  showFormModal.value = true;
};

const openViewModal = (item) => {
  selectedBug.value = item;
  showViewModal.value = true;
};

const closeFormModal = () => {
  showFormModal.value = false;
};

const closeViewModal = () => {
  showViewModal.value = false;
};

const saveBug = async () => {
  if (!formState.value.bugId || !formCommentsRaw.value) {
    triggerToast("error", "กรุณากรอกฟิลด์ที่จำเป็น (*) ให้ครบถ้วน");
    return;
  }

  // Parse comments
  formState.value.comments = formCommentsRaw.value
    .split("\n")
    .map((c) => c.trim())
    .filter((c) => c.length > 0);

  // ปิด Modal form ก่อน แล้วค่อยแสดง loading
  closeFormModal();

  if (formType.value === "add") {
    showLoadingPopup("กำลังบันทึกข้อมูลบั๊กใหม่...");
    const success = await addBug({
      bugId: formState.value.bugId,
      title: formState.value.title,
      comments: formState.value.comments,
      expectation: formState.value.expectation,
      driveLink: formState.value.driveLink,
      figmaLink: formState.value.figmaLink,
      status: formState.value.status,
      createdDate: formState.value.createdDate,
      fixedDate: formState.value.fixedDate,
      passedDate: formState.value.passedDate,
    });
    closeLoadingPopup();
    if (success) {
      triggerToast("success", "เพิ่มบันทึกบั๊กสำเร็จแล้ว");
    } else {
      triggerToast("error", "ไม่สามารถเพิ่มข้อมูลได้");
    }
  } else {
    showLoadingPopup("กำลังอัปเดตข้อมูลบั๊ก...");
    const success = await updateBug(formState.value);
    closeLoadingPopup();
    if (success) {
      triggerToast("success", "แก้ไขข้อมูลบั๊กสำเร็จแล้ว");
    } else {
      triggerToast("error", "ไม่สามารถอัปเดตข้อมูลได้");
    }
  }
};

const confirmDelete = (item) => {
  bugToDelete.value = item;
  showDeleteModal.value = true;
};

const closeDeleteModal = () => {
  showDeleteModal.value = false;
  bugToDelete.value = null;
};

const executeDelete = async () => {
  console.log("executeDelete triggered for:", bugToDelete.value);
  if (bugToDelete.value) {
    const targetRowIndex = bugToDelete.value.sheetRowIndex;
    console.log("sheetRowIndex is:", targetRowIndex);
    // ปิด Modal ยืนยันก่อน แล้วแสดง loading
    closeDeleteModal();
    showLoadingPopup("กำลังลบข้อมูลบั๊ก...");

    try {
      const success = await deleteBug(targetRowIndex);
      console.log("deleteBug response:", success);
      closeLoadingPopup();

      if (success) {
        triggerToast("success", "ลบข้อมูลบั๊กสำเร็จแล้ว");
        // ถ้าลบตัวสุดท้ายของหน้านี้ ให้ถอยกลับ 1 หน้า
        if (paginatedBugs.value.length === 0 && currentPage.value > 1) {
          currentPage.value--;
        }
      } else {
        triggerToast("error", "ไม่สามารถลบข้อมูลได้");
      }
    } catch (err) {
      console.error("Error in executeDelete:", err);
      closeLoadingPopup();
      triggerToast("error", "เกิดข้อผิดพลาดขณะลบข้อมูล");
    }
  } else {
    console.warn("bugToDelete.value is null or undefined");
  }
};

const statusClass = (status) => {
  if (status === "Defect")
    return "bg-[#ffe4e6] text-[#b91c1c] border border-red-200";
  if (status === "Demo")
    return "bg-[#fef3c7] text-[#b45309] border border-amber-200";
  if (status === "Demo Fail")
    return "bg-[#fee2e2] text-[#dc2626] border border-red-200";
  if (status === "Pass")
    return "bg-[#d1fae5] text-[#065f46] border border-emerald-200";
  if (status === "Ready For Demo")
    return "bg-[#dbeafe] text-[#1e40af] border border-blue-200";
  if (status === "Doing")
    return "bg-[#f3e8ff] text-[#6b21a8] border border-purple-200";
  return "bg-slate-100 text-slate-700 border border-slate-200";
};
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
  height: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #f1f5f9;
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #e2e8f0;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.98);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.animate-fade-in {
  animation: fadeIn 0.2s ease-out forwards;
}
</style>
