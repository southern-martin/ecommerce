import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 3000,
    allowedHosts: true,
    proxy: {
      // ── Auth & User services ──
      '/api/v1/auth': { target: 'http://localhost:28090', changeOrigin: true },
      '/api/v1/sellers': { target: 'http://localhost:28091', changeOrigin: true },
      '/api/v1/addresses': { target: 'http://localhost:28091', changeOrigin: true },
      '/api/v1/users': { target: 'http://localhost:28091', changeOrigin: true },
      '/api/v1/wishlist': { target: 'http://localhost:28091', changeOrigin: true },

      // ── Product service ──
      '/api/v1/categories': { target: 'http://localhost:28081', changeOrigin: true },
      '/api/v1/products': { target: 'http://localhost:28081', changeOrigin: true },

      // ── Core commerce ──
      '/api/v1/cart': { target: 'http://localhost:28082', changeOrigin: true },
      '/api/v1/orders': { target: 'http://localhost:28083', changeOrigin: true },
      '/api/v1/payments': { target: 'http://localhost:28084', changeOrigin: true },

      // ── Supporting services ──
      '/api/v1/search': { target: 'http://localhost:28085', changeOrigin: true },
      '/api/v1/reviews': { target: 'http://localhost:28086', changeOrigin: true },
      '/api/v1/notifications': { target: 'http://localhost:28087', changeOrigin: true },
      '/api/v1/chat': { target: 'http://localhost:28088', changeOrigin: true },
      '/api/v1/media': { target: 'http://localhost:28089', changeOrigin: true },
      '/api/v1/promotions': { target: 'http://localhost:28093', changeOrigin: true },
      '/api/v1/returns': { target: 'http://localhost:28094', changeOrigin: true },
      '/api/v1/shipments': { target: 'http://localhost:28095', changeOrigin: true },
      '/api/v1/shipping': { target: 'http://localhost:28095', changeOrigin: true },
      '/api/v1/loyalty': { target: 'http://localhost:28096', changeOrigin: true },
      '/api/v1/affiliate': { target: 'http://localhost:28097', changeOrigin: true },
      '/api/v1/tax': { target: 'http://localhost:28098', changeOrigin: true },
      '/api/v1/cms': { target: 'http://localhost:28099', changeOrigin: true },

      // ── Seller routes (specific service backends) ──
      '/api/v1/seller/products': { target: 'http://localhost:28081', changeOrigin: true },
      '/api/v1/seller/orders': { target: 'http://localhost:28083', changeOrigin: true },
      '/api/v1/seller/coupons': { target: 'http://localhost:28093', changeOrigin: true },
      '/api/v1/seller/shipments': { target: 'http://localhost:28095', changeOrigin: true },
      '/api/v1/seller/carriers': { target: 'http://localhost:28095', changeOrigin: true },
      '/api/v1/seller/returns': { target: 'http://localhost:28094', changeOrigin: true },

      // ── Admin routes (specific service backends, BEFORE catch-all) ──
      '/api/v1/admin/promotions': { target: 'http://localhost:28093', changeOrigin: true },
      '/api/v1/admin/categories': { target: 'http://localhost:28081', changeOrigin: true },
      '/api/v1/admin/attributes': { target: 'http://localhost:28081', changeOrigin: true },
      '/api/v1/admin/carriers': { target: 'http://localhost:28095', changeOrigin: true },
      '/api/v1/admin/banners': { target: 'http://localhost:28099', changeOrigin: true },
      '/api/v1/admin/pages': { target: 'http://localhost:28099', changeOrigin: true },
      '/api/v1/admin/content': { target: 'http://localhost:28099', changeOrigin: true },
      '/api/v1/admin/reviews': { target: 'http://localhost:28086', changeOrigin: true },
      '/api/v1/admin/disputes': { target: 'http://localhost:28094', changeOrigin: true },
      '/api/v1/admin/affiliates': { target: 'http://localhost:28097', changeOrigin: true },
      '/api/v1/admin/tax': { target: 'http://localhost:28098', changeOrigin: true },
      // Admin catch-all (user service — sellers approval, user management)
      '/api/v1/admin': { target: 'http://localhost:28091', changeOrigin: true },

      // ── WebSocket ──
      '/ws': { target: 'ws://localhost:28088', ws: true },
    },
  },
});
