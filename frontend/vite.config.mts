import path from 'path';
import fs from 'fs';
import dns from 'dns';
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import tailwindcss from '@tailwindcss/vite';
import svgr from 'vite-plugin-svgr';

dns.setDefaultResultOrder('verbatim');

const allHtmlEntries = fs
  .readdirSync('.')
  .filter((file) => path.extname(file) === '.html')
  .reduce<Record<string, string>>((acc, file) => {
    acc[path.basename(file, '.html')] = path.resolve(__dirname, file);
    return acc;
  }, {});

export default defineConfig({
  build: {
    rollupOptions: {
      input: allHtmlEntries,
    },
  },
  plugins: [
    react(),
    tailwindcss(),
    svgr(),
  ],
  resolve: {
    alias: {
      '@app': path.resolve(__dirname, './src/app'),
      '@api': path.resolve(__dirname, './src/api'),
      '@components': path.resolve(__dirname, './src/components'),
      '@stores': path.resolve(__dirname, './src/stores'),
      '@features': path.resolve(__dirname, './src/features'),
      '@lib': path.resolve(__dirname, './src/lib'),
      '@utils': path.resolve(__dirname, './src/utils'),
      '@i18n': path.resolve(__dirname, './src/i18n'),
    },
  },
  server: {
    port: 3000,
    allowedHosts: ['nsnz.ngrok.io'],
  },
});
