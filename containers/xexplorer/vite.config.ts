import dns from 'dns';
import path from 'path';
import react from '@vitejs/plugin-react-swc';
import { defineConfig, loadEnv, splitVendorChunkPlugin } from 'vite';
import mkcert from 'vite-plugin-mkcert';
import { nodePolyfills } from 'vite-plugin-node-polyfills';
import svgr from 'vite-plugin-svgr';
import tsconfigPaths from 'vite-tsconfig-paths';

// https://vitejs.dev/config/
(dns as any).setDefaultResultOrder('verbatim');

export default ({ mode }) => {
  process.env = { ...process.env, ...loadEnv(mode, process.cwd()) };

  const shouldUseTSL =
    process.env.VITE_APP_USE_HTTPS?.toLowerCase() !== 'false';

  return defineConfig({
    plugins: [
      svgr(),
      react(),
      tsconfigPaths(),
      splitVendorChunkPlugin(),
      nodePolyfills({
        globals: { Buffer: true, global: true, process: true }
      }),
      ...(shouldUseTSL ? [mkcert()] : [])
    ],
    resolve: {
      alias: {
        '~bootstrap': path.resolve(__dirname, 'node_modules/bootstrap')
      }
    },
    define: {
      'process.env': process.env
    },
    build: {
      outDir: 'build',
      cssMinify: true,
      minify: true
    },
    server: {
      port: 3002,
      strictPort: true,
      https: shouldUseTSL,
      host: '0.0.0.0',
      hmr: {
        overlay: false
      },
      watch: {
        usePolling: false,
        useFsEvents: false
      }
    },
    preview: {
      port: 3002,
      strictPort: true,
      https: shouldUseTSL,
      host: '0.0.0.0'
    }
  });
};
