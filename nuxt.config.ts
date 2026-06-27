import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  // 1. นำ Vite Plugin ของ Tailwind มาใส่ในส่วนของ vite
  vite: {
    plugins: [
      tailwindcss(),
    ],
  },
  app: {
    head: {
      title: 'MCOP INVESTMENT MONITOR'
    }
  },
  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || process.env.API_BASE || 'http://localhost:8080/api/data'
    }
  },

  // 2. ระบุไฟล์ CSS หลักที่จะให้ Nuxt โหลดใช้งาน
  css: ['~/assets/main.css'],
})