<script lang="ts">
export default {
  name: 'expansion-componetn'
}
</script>
<template>
  <el-dialog
    v-model="props.expansionVisible"
    title="Expansion Conf"
    width="80%"
    @close="emits('closeDialog')"
  >
    <div style="display: flex">
      <div
        class="resu"
        style="background-color: black; height: 500px; padding: 5px 10px; overflow: scroll; flex: 6"
      >
        <div style="color: aliceblue">
          SimpLogServer :: {{ props.serverName }} :: created By leeks
        </div>
        <div style="color: aliceblue; margin: 2px; white-space: pre" v-html="state.logger"></div>
      </div>
      <div style="flex: 3">
        <el-form :model="body" label-width="auto" style="max-width: 600px">
          <el-form-item label="ServerName">
            <el-input :disabled="true" v-model="body.serverName" />
          </el-form-item>
          <el-form-item label="Proxy">
            <el-select v-model="body.locationName" placeholder="Proxy">
              <el-option
                v-for="item in state.servers"
                :label="item.key"
                :value="item.key"
                :key="item.key"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="Upstream">
            <el-select
              v-model="body.upstreamName"
              placeholder="Proxy"
              @change="changeHosts"
              filterable
              allow-create
              default-first-option
              :reserve-keyword="false"
            >
              <el-option
                v-for="item in state.upstreams"
                :label="item.key"
                :value="item.key"
                :key="item.key"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="Hosts" v-if="body.upstreamName">
            <el-input v-model="host">
              <template #append>
                <el-button :icon="CirclePlus" @click="addHost()" />
              </template>
            </el-input>
            <el-input v-for="item in hosts" :disabled="true" :model-value="item" :key="item">
              <template #append>
                <el-button :icon="Delete" @click="deleteHost(item)" />
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="Submit">
            <div style="display: flex; align-items: center; justify-content: center">
              <el-button type="primary" @click="emits('closeDialog')">Close</el-button>
              <el-button type="success" style="background-color: blue;" @click="reload">Reload</el-button>
              <el-button type="success" @click="releaseExpandConf">Release</el-button>
              <el-button type="danger" @click="uploadExpandConf">Preview</el-button>
            </div>
          </el-form-item>
          <el-form-item label="Tips">
            <div style="color: red;">1.reload  重启 nginx服务</div>
            <div style="color: red;">2.release 覆盖 nginx 并且 test 检测语法</div>
            <div style="color: red;">3.preview 预览修改后的配置</div>
          </el-form-item>
        </el-form>
       
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { getProxyList, nginxExpansion, nginxExpansionPreview, nginxReload } from '@/api/nginx'
import { CirclePlus, Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus';
import { reactive, ref, watch } from 'vue'

// 扩容组件
const props = defineProps<{
  expansionVisible: boolean
  serverName: string
}>()
const hosts = ref<string[]>([])
const host = ref('')

function changeHosts(v: string) {
  const findItem = state.upstreams.find((e) => e.key == v)
  console.log('findItem',findItem);
  if (findItem.value.server instanceof Array) {
    hosts.value = findItem.value.server
    return 
  }
  hosts.value = [findItem.value.server]
}
function deleteHost(item) {
  hosts.value = hosts.value.filter((v) => v !== item)
}
function addHost() {
  if(host.value == ''){
    return
  }
  hosts.value.push(host.value)
  host.value = ''
}
async function uploadExpandConf() {
  body.server = hosts.value
  const conf = await nginxExpansionPreview(body)
  state.logger = conf.Data
  console.log('body', body)
}

async function releaseExpandConf() {
  body.server = hosts.value
  const conf = await nginxExpansion(body)
  if(conf.Code){
    return ElMessage.error(`error:${conf.Message}`)
  }
  ElMessage.success("release success")
  state.logger = conf.Data
}

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
    state.logger = stream.Data
    emits('showReleaseDialog')
  })
}

const state = reactive({
  //   httpConf: {},
  logger: '',
  servers: [],
  upstreams: []
})
const emits = defineEmits(['closeDialog','showReleaseDialog'])

const body = reactive({
  upstreamName: '',
  server: [''],
  locationName: '',
  serverName: ''
})

async function init() {
  const data = await getProxyList()
  state.logger = data.Data.conf.replace(/\n/g, '<br>')
  state.servers = data.Data.servers
  state.upstreams = data.Data.upstreams
}
watch(props, (newVal) => {
  if (!newVal.expansionVisible) {
    return
  }
  init()
  body.serverName = props.serverName
})
</script>

<style scoped></style>
