<script lang="ts">
export default {
  name: 'main-logger'
}
</script>
<script lang="ts" setup>
const props = defineProps<{
  serverList: any[]
  createServerVisible: boolean
}>()
const emits = defineEmits(['getMainLogList', 'changeVisible'])
function GetMainLogList() {
  emits('getMainLogList')
}
function changeVisible(bool: boolean) {
  emits('changeVisible', bool)
}
</script>

<template>
  <el-card shadow="hover">
    <div style="height: 70px; display: flex; align-items: center; justify-content: space-around">
      <div class="flex-item">
        <div @click="changeVisible(true)" style="color: rgb(207, 15, 124); cursor: pointer">
          CreateServer
        </div>
      </div>
      <div class="flex-item">
        <div @click="GetMainLogList()" style="color: rgb(207, 15, 124); cursor: pointer">
          CheckLog
        </div>
      </div>
      <div class="flex-item">
        <div style="font-weight: 700">ServerCounts</div>
        <div style="color: rgb(207, 15, 124)">{{ props.serverList.length }}</div>
      </div>
    </div>
  </el-card>
</template>
