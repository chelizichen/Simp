<script lang="ts">
export default {
  name: 'main-logger'
}
</script>
<script lang="ts" setup>
import { backupNginx, getBackupFile, getBackupList, nginxReload } from '@/api/nginx';
import { ElMessage, ElMessageBox } from 'element-plus';
import { reactive, watch } from 'vue';

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

const histroyState = reactive({
  historyListVisible:false,
  fileList:[],
  logger:"",
  fileName:"",
})
watch(()=>histroyState.historyListVisible,async function(newVal){
  if(!newVal){
    return
  }
  const data = await getBackupList()
  if(data.Code){
    return 
  }
  data.Data.push("origin")
  histroyState.fileList = data.Data.reverse()
  histroyState.fileList
})

async function reload() {
  ElMessageBox.prompt('Are you sure to reload', 'Confirm', {
    confirmButtonText: 'OK',
    cancelButtonText: 'Cancel',
    inputPlaceholder: 'input password'
  })
  .then(async ({ value }) => {
    if(value != "0504"){
      return ElMessage.error(`reload error:password`)
    }
    const stream = await nginxReload()
    if(stream.Code){
      return ElMessage.error(`reload error:${stream.Message}`)
    }
    ElMessage.success("reload success")
    histroyState.logger = stream.Data
  })
}

async function checkFile(){
  const data = await getBackupFile({fileName:histroyState.fileName})
  histroyState.logger = data.Data
}

async function backup() {
  const data = await backupNginx({fileName:histroyState.fileName})
  histroyState.logger = data.Data
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
        <div @click="histroyState.historyListVisible = true" style="color: rgb(207, 15, 124); cursor: pointer;font-weight: 700;">
          Histroy
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
  <el-dialog
    v-model="histroyState.historyListVisible"
    title="Expansion Conf"
    width="80%"
    @close="histroyState.historyListVisible = false"
    >
    <div style="display: flex">
      <div
        class="resu"
        style="background-color: black; height: 500px; padding: 5px 10px; overflow: scroll; flex: 6"
      >
        <div style="color: aliceblue">
          SimpMainLogServer :: created By leeks
        </div>
        <div style="color: aliceblue; margin: 2px; white-space: pre" v-html="histroyState.logger"></div>
      </div>
      <div style="flex: 3">
        <el-form :model="histroyState" label-width="auto" style="max-width: 600px">
          <el-form-item label="Upstream">
            <el-select
              v-model="histroyState.fileName"
              placeholder="choose file"
              filterable
              default-first-option
              :reserve-keyword="false"
              @change="checkFile"
            >
              <el-option
                v-for="item in histroyState.fileList"
                :label="item"
                :value="item"
                :key="item"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="Submit">
            <div style="display: flex; align-items: center; justify-content: center">
              <el-button type="primary" @click="histroyState.historyListVisible = false">Close</el-button>
              <el-button type="success" style="background-color: blue;" @click="reload">Reload</el-button>
              <el-button type="danger" @click="backup">Cover</el-button>
            </div>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </el-dialog>
</template>, nginxReloadimport { ElMessageBox, ElMessage } from 'element-plus';

