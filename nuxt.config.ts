import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  // 1. นำ Vite Plugin ของ Tailwind มาใส่ในส่วนของ vite
  vite: {
    plugins: [
      tailwindcss(),
    ],
  },
  
  // 2. ระบุไฟล์ CSS หลักที่จะให้ Nuxt โหลดใช้งาน
  css: ['~/assets/main.css'],
})