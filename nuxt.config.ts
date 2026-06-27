import tailwindcss from '@tailwindcss/vite'
import fs from 'node:fs'
import path from 'node:path'

// --- VUE INLINE FIX START ---
try {
  const PAGES_DIR = path.resolve('./app/pages');

  function walk(dir) {
    const files = fs.readdirSync(dir);
    for (const file of files) {
      const fullPath = path.join(dir, file);
      const stat = fs.statSync(fullPath);
      if (stat.isDirectory()) {
        walk(fullPath);
      } else if (file.endsWith('.vue')) {
        let content = fs.readFileSync(fullPath, 'utf8');
        let modified = false;

        // Regex definitions with dotAll and multiline whitespace flags
        const fixedClickRegex = /@click="if\s*\(!formState\.fixedDate\)\s*formState\.fixedDate\s*=\s*new\s*Date\(\)\.toLocaleString\('th-TH'\)"/gs;
        const fixedFocusRegex = /@focus="if\s*\(!formState\.fixedDate\)\s*formState\.fixedDate\s*=\s*new\s*Date\(\)\.toLocaleString\('th-TH'\)"/gs;
        const passedClickRegex = /@click="if\s*\(!formState\.passedDate\)\s*formState\.passedDate\s*=\s*new\s*Date\(\)\.toLocaleString\('th-TH'\)"/gs;
        const passedFocusRegex = /@focus="if\s*\(!formState\.passedDate\)\s*formState\.passedDate\s*=\s*new\s*Date\(\)\.toLocaleString\('th-TH'\)"/gs;

        const replaced1 = content.replace(fixedClickRegex, `@click="!formState.fixedDate && (formState.fixedDate = new Date().toLocaleString('th-TH'))"`);
        const replaced2 = replaced1.replace(fixedFocusRegex, `@focus="!formState.fixedDate && (formState.fixedDate = new Date().toLocaleString('th-TH'))"`);
        const replaced3 = replaced2.replace(passedClickRegex, `@click="!formState.passedDate && (formState.passedDate = new Date().toLocaleString('th-TH'))"`);
        const replaced4 = replaced3.replace(passedFocusRegex, `@focus="!formState.passedDate && (formState.passedDate = new Date().toLocaleString('th-TH'))"`);

        if (replaced4 !== content) {
          content = replaced4;
          modified = true;
        }

        if (modified) {
          console.log(`[VUE INLINE FIX] Corrected syntax in: ${file}`);
          fs.writeFileSync(fullPath, content, 'utf8');
        }
      }
    }
  }

  walk(PAGES_DIR);
} catch (e) {
  console.error('[VUE INLINE FIX] Error:', e);
}
// --- VUE INLINE FIX END ---

export default defineNuxtConfig({
  // 1. นำ Vite Plugin ของ Tailwind มาใส่ในส่วนของ vite
  vite: {
    plugins: [
      tailwindcss(),
    ],
  },
  app: {
    head: {
      title: 'MCOP INVESTMENT MONITOR',
      link: [
        { rel: 'icon', type: 'image/png', href: '/favicon.png' }
      ]
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