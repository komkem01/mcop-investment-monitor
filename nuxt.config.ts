import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  // 1. นำ Vite Plugin ของ Tailwind มาใส่ในส่วนของ vite
  vite: {
    plugins: [
      tailwindcss(),
    ],
  },
  app: {
    baseURL: process.env.NUXT_PUBLIC_BASE_URL,
    head: {
      title: 'MCOP INVESTMENT MONITOR'
    }
  },

  // 2. ระบุไฟล์ CSS หลักที่จะให้ Nuxt โหลดใช้งาน
  css: ['~/assets/main.css'],
})