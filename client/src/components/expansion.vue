<script lang="ts">
export default {
  name: 'expansion-componetn'
}
</script>
<template>
  <el-dialog v-model="props.expansionVisible" title="Expansion Conf" width="80%">
    <div>
        <div
                  class="resu"
                  style="
                    background-color: black;
                    height: 500px;
                    padding: 5px 10px;
                    overflow: scroll;
                  "
                >
                  <div style="color: aliceblue">
                    SimpLogServer :: {{ props.serverName }} :: created By leeks
                  </div>
                  <div
                    style="color: aliceblue; margin: 2px;white-space: pre;"
                    v-html="state.logger"
                  >
                  </div>
                </div></div>
    <template #footer>
      <div style="display: flex; align-items: center; justify-content: center">
        <el-button type="primary" @click="emits('closeDialog')">Close</el-button>
        <el-button type="success">Release</el-button>
        <el-button type="danger">Upload</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { getProxyList } from '@/api/nginx'
import { onMounted, reactive } from 'vue'

// 扩容组件
const props = defineProps<{
  expansionVisible: boolean
  serverName: string
}>()
const state = reactive({
//   httpConf: {},
  logger:''
})
const emits = defineEmits(["closeDialog"])

async function init() {
  const data = await getProxyList()
  state.logger = data.Data.conf.replace(/\n/g, '<br>')
}
onMounted(() => {
    init()
})
</script>

<style scoped></style>
