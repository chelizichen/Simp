import { fileURLToPath, URL } from 'node:url'

import { defineConfig, splitVendorChunkPlugin,loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import viteCompression from 'vite-plugin-compression'
import { chunkSplitPlugin } from 'vite-plugin-chunk-split';

// https://vitejs.dev/config/
export default defineConfig(({ command, mode })=>{
  const env = loadEnv(mode, process.cwd(), '')
  return {
    plugins: [
      vue(),
      splitVendorChunkPlugin(),
      chunkSplitPlugin(),
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    },
    server: {
      port: 3000,
      proxy: {
        '/simpserver': {
          target: 'http://localhost:8511/',
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/simpserver/, '/simpserver/') // 不可以省略rewrite
        },
        "/simpexpansionserver":{
          target: 'http://localhost:8518/',
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/simpexpansionserver/, '/simpexpansionserver/') // 不可以省略rewrite
        },
        "/simpshellserver":{
          target: 'http://localhost:8515/',
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/simpshellserver/, '/simpshellserver/') // 不可以省略rewrite
        }
      }
    },
    base:env.BASE,
    build:{
      minify:true,
      terserOptions:{
        compress: {
          drop_console: true,
          drop_debugger: true
        }
      },
      rollupOptions:[
        viteCompression({
          verbose: true, // 是否在控制台中输出压缩结果
          disable: false,
          threshold: 10240, // 如果体积大于阈值，将被压缩，单位为b，体积过小时请不要压缩，以免适得其反
          algorithm: 'gzip', // 压缩算法，可选['gzip'，' brotliccompress '，'deflate '，'deflateRaw']
          ext: '.gz',
          deleteOriginFile: true // 源文件压缩后是否删除(我为了看压缩后的效果，先选择了true)
        }),
      ]
    }
  }

})
