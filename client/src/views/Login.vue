<script lang="ts">
export default {
  name: 'login-vue'
}
</script>
<script lang="ts" setup>
import API from '@/api/server'
import { ElMessage } from 'element-plus'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { md5 } from 'js-md5'
const router = useRouter()
const token = ref('')
async function saveToken() {
  const data = new FormData()
  const tkn = md5(token.value)
  data.append('token', tkn)
  const ret = await API.Login(data)
  if (ret.Data) {
    ElMessage.error('Please enter a valid token.')
  } else {
    router.push('/server')
    localStorage.setItem('token', tkn)
  }
}
</script>
<template>
  <div class="body">
    <div class="container">
      <div class="hello">Hello Simp!</div>
      <el-form :inline="true">
        <el-form-item label="Token">
          <el-input v-model="token" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="saveToken()">Submit</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
  <el-footer>
    <el-divider content-position="center">
      <div style="color: rgb(207, 15, 124); font-size: 18px">Copyright©2023-2024</div>
    </el-divider>
    <el-divider content-position="center">
      <div style="color: rgb(207, 15, 124); font-size: 18px">
        SimpServer Started on AliCloud Platform
      </div>
    </el-divider>
  </el-footer>
</template>

<style lang="less">
.body {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 90vh;
  background-color: #ffffff;
  .container {
    text-align: center;
    width: 25vw;
    padding: 20px;
    border: 1px solid #ccc;
    border-radius: 5px;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
    .hello {
      margin-bottom: 10px;
      font-size: 26px;
      color: var(--el-button-hover-bg-color);
    }
  }
}
</style>
