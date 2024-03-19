<script lang="ts">
export default {
  name: 'server-component'
}
</script>
<script setup lang="ts">
import { ElLoading, ElMessage, ElMessageBox, ElPopconfirm, type UploadUserFile } from 'element-plus'
import { onMounted, reactive, ref, watch } from 'vue'
import API from '../api/server'
import asideComponent from '@/components/aside.vue'
import mainLogger from '@/components/mainlogger.vue'
import { InfoFilled, Remove } from '@element-plus/icons-vue'
import { cloneDeep, reverse } from 'lodash'
import expansionComponent from '@/components/expansion.vue'
import { getProxyList } from '@/api/nginx'
const state = reactive({
  serverList: [],
  packageList: [],
  packagePortList: [],
  serverName: '',
  uploadVisible: false,
  status: {
    pid: 0,
    status: ''
  },
  configVisible: false,
  createServerVisible: false,
  createServerName: '',
  releaseVisible: false,
  shutdownVisible: false,
  selectRelease: '',
  logger: '',
  loggerList: [],
  loggerFile: '',
  pattern: '',
  activeName: 'logger',
  apis: [], // API接口 由gin生成
  doc: '',
  rows: '',
  mainLogList: [],
  expansionVisible: false
})
const uploadForm = ref({
  serverName: '',
  file: null,
  doc: ''
})
const config = ref({
  Name: '',
  Port: 0,
  Type: '',
  StaticPath: '',
  Storage: '',
  Proxy: null,
  Host: ''
})

const rowList = [50, 100, 500, 1000]

async function handleOpen(serverName: string) {
  state.serverName = serverName
  await getServerPackageList(serverName)
  uploadForm.value.serverName = serverName
}

async function getServerPackageList(serverName: string) {
  const formData = new FormData()
  formData.append('serverName', serverName)
  const resp = await API.GetServerPackageList(formData)
  state.packageList = resp.Data || []
  if (state.packageList.length) {
    const resp = await API.CheckServer(formData)
    const ret = resp.Data
    console.log('ret', ret)
    state.status = {
      status: ret.status
        ? '<div style="color:#55bd55">online</div>'
        : '<div style="color: rgb(207, 15, 124)">offline</div>',
      pid: ret.pid
    }
  } else {
    state.status = {
      status: '<div style="color: rgb(207, 15, 124)">offline</div>',
      pid: 0
    }
  }
}

async function restartServer() {
  if (!state.selectRelease || !state.serverName) {
    ElMessage({
      type: 'info',
      message: '发布失败！请选择指定的服务和发布包'
    })
    return
  }
  const loading = ElLoading.service({
    lock: true,
    text: 'Loading',
    spinner: 'el-icon-loading',
    background: 'rgba(0, 0, 0, 0.7)'
  })

  if (multpieNodesState.selectUpstream == SingleNode) {
    const formData = new FormData()
    formData.append('fileName', state.selectRelease)
    formData.append('serverName', state.serverName)
    const resp = await API.RestartServer(formData)
    state.status = {
      status: resp.Data.status
        ? '<div style="color:#55bd55">online</div>'
        : '<div style="color: rgb(207, 15, 124)">offline</div>',
      pid: resp.Data.pid
    }
    if (resp.Code) {
      ElMessage({
        type: 'error',
        message: '发布失败' + resp.Message
      })
    } else {
      ElMessage({
        type: 'success',
        message: '发布成功'
      })
      state.releaseVisible = false
    }
    loading.close()
  } else {
    console.log('multipie', multpieNodesState.selectHosts)
    if (!multpieNodesState.selectHosts.length) {
      ElMessage.error('请选择至少一个节点进行发布')
      return loading.close()
    }
    const portRegex = /:(\d+)/
    let ports = multpieNodesState.selectHosts.map((v: string) => {
      const port = Number(v.match(portRegex)[1])
      return port
    })
    const formData = new FormData()
    formData.append('serverName', state.serverName)
    const resp = await API.CheckConfig(formData)
    const mainPort = resp.Data.Server.Port
    console.log('mainPort', mainPort)
    console.log('ports', ports)
    const hasMainPort = ports.indexOf(mainPort)
    if (hasMainPort == -1) {
      ElMessage.error('必须包含主控节点')
      return loading.close()
    }
    ports.splice(hasMainPort, 1)
    {
      const formData = new FormData()
      formData.append('serverName', state.serverName)
      formData.append('fileName', state.selectRelease)
      const resp = await API.RestartServer(formData)
      state.status = {
        status: resp.Data.status
          ? '<div style="color:#55bd55">online</div>'
          : '<div style="color: rgb(207, 15, 124)">offline</div>',
        pid: resp.Data.pid
      }
    }
    if (!ports.length) {
      return loading.close()
    }
    {
      await Promise.all(
        ports.map(async (targetPort) => {
          const formData = new FormData()
          formData.append('serverName', state.serverName)
          formData.append('fileName', state.selectRelease)
          formData.append('targetPort', String(targetPort))
          const resp = await API.RestartServer(formData)
          return resp
        })
      )
    }
    state.releaseVisible = false
    loading.close()
    // formData.append('targetPort',)
  }
}

async function GetMainLogList() {
  state.serverName = ''
  const List = await API.GetMainLogList()
  const Data = List.Data
  state.mainLogList = Data
  reverse(state.mainLogList)
  state.loggerFile = Data.length ? Data[0] : ''
}

async function GetServerLogger() {
  if (state.serverName) {
    const formData = new FormData()
    formData.append('serverName', state.serverName)
    formData.append('fileName', state.loggerFile)
    formData.append('pattern', state.pattern || '')
    formData.append('rows', state.rows || '50')
    const rest = await API.GetLogger(formData)
    state.logger = rest.Data.split('\n')
  } else {
    const formData = new FormData()
    formData.append('logFile', state.loggerFile)
    formData.append('pattern', state.pattern || '')
    formData.append('rows', state.rows || '50')
    const rest = await API.GetMainLogger(formData)
    state.logger = rest.Data.split('\n')
  }
}

async function showConfig() {
  const formData = new FormData()
  formData.append('serverName', state.serverName)
  const resp = await API.CheckConfig(formData)
  config.value = resp.Data.Server
  state.configVisible = true
}

const childServiceList = ref([])
const choseServices = ref([])
const childServiceObj = ref({})
watch(
  () => state.shutdownVisible,
  async function (newVal) {
    if (!newVal) {
      return
    }
    const formData = new FormData()
    formData.append('serverName', state.serverName)
    const list = await API.getChildStats(formData)
    childServiceObj.value = list.Data
    const arr = []
    for (let k in list.Data) {
      arr.push(k)
    }
    childServiceList.value = arr
  }
)

async function shutdownServers() {
  const loading = ElLoading.service({
    lock: true,
    text: 'Loading',
    spinner: 'el-icon-loading',
    background: 'rgba(0, 0, 0, 0.7)'
  })
  try {
    const arrs = cloneDeep(choseServices.value).map((v) => {
      return v.split('|')[0].trim()
    })
    if (!arrs.length) {
      return ElMessage.error(`server is not online`)
    }
    await Promise.all(
      arrs.map(async (v) => {
        const formData = new FormData()
        formData.append('serverName', v)
        const data = await API.ShutDownServer(formData)
        if (data.Code) {
          ElMessage({
            type: 'error',
            message: '关闭失败!' + data.Message + '| ' + v
          })
          return
        }
        ElMessage({
          type: 'success',
          message: '成功关闭节点!' + '| ' + v
        })
        console.log(arrs)
        state.shutdownVisible = false
      })
    )
  } catch (e) {
    ElMessage.error(e)
  } finally {
    choseServices.value = []
    getServerPackageList(state.serverName)
    loading.close()
  }

  // const pid = state.status.pid
  // if (pid == 0) {
  //   return ElMessage.error('关闭失败！该服务并未启动')
  // }
}

async function ShutDownServer() {
  state.shutdownVisible = true
  return
}

async function fetchServerList() {
  const resp = await API.GetServerList()
  state.serverList = resp.Data || []
}

async function uploadFile() {
  const loading = ElLoading.service({
    lock: true,
    text: 'Loading',
    spinner: 'el-icon-loading',
    background: 'rgba(0, 0, 0, 0.7)'
  })
  const formData = new FormData()
  formData.append('serverName', uploadForm.value.serverName)
  formData.append('file', uploadForm.value.file)
  formData.append('doc', uploadForm.value.doc)
  const data = await API.UploadServer(formData)
  uploadForm.value.file = null
  if (data.Code) {
    ElMessage({
      type: 'error',
      message: '上传失败' + data.Message
    })
  } else {
    ElMessage({
      type: 'success',
      message: '上传成功'
    })
    state.uploadVisible = false
  }
  getServerPackageList(state.serverName)
  loading.close()
}

async function createServer() {
  const formData = new FormData()
  formData.append('serverName', state.createServerName)
  const data = await API.CreateServer(formData)
  ElMessage({
    type: 'success',
    message: '添加成功!'
  })
  state.createServerVisible = false
  fetchServerList()
}

async function previewSubServer() {
  const conf = config.value
  const target = `http://${conf.Host}:${conf.Port}/${conf.StaticPath}/`
  // http://localhost:8511/web/
  // http://localhost:8511/web/
  window.open(target)
}

function MapKeysToUpper(obj: any) {
  var newObj = {} as any
  for (var key in obj) {
    if (typeof key === 'string') {
      newObj[key.charAt(0).toUpperCase() + key.slice(1)] = obj[key]
    } else {
      newObj[key] = obj[key]
    }
  }
  return newObj
}
async function uploadConfig() {
  const body = MapKeysToUpper({
    ServerName: state.serverName,
    Conf: {
      Server: config.value
    }
  })
  const ret = await API.CoverConfig(body)
  if (ret.Data) {
    ElMessage({
      type: 'error',
      message: '覆盖失败!' + ret.Message
    })
    return
  }
  ElMessage({
    type: 'success',
    message: '覆盖成功!'
  })
  state.configVisible = false
}

const fileList = ref<UploadUserFile[]>([])

function handleFileChange(file: any) {
  if (!file.name.includes(uploadForm.value.serverName)) {
    ElMessage.error(`请上传正确的服务包 [ ${uploadForm.value.serverName} ] `)
    uploadForm.value.file = null
    fileList.value = []
  } else {
    uploadForm.value.file = file.raw
    fileList.value = [file]
  }
}

async function DeletePackage(hash: string) {
  ElMessageBox.prompt('Are you sure to delete this package', 'Confirm', {
    confirmButtonText: 'OK',
    cancelButtonText: 'Cancel',
    inputPlaceholder: 'input password'
  })
    .then(async ({ value }) => {
      if (value != '0504') {
        return false
      }
      const formData = new FormData()
      const fileName = `${state.serverName}_${hash}.tar.gz`
      formData.append('serverName', state.serverName)
      formData.append('fileName', fileName)
      const rest = await API.DeletePackage(formData)
      if (rest.Code) {
        ElMessage({
          type: 'info',
          message: 'Delete canceled'
        })
        return
      }
      await getServerPackageList(state.serverName)
      state.selectRelease = ''
      ElMessage({
        type: 'success',
        message: `Delete Success`
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: 'Delete canceled'
      })
    })
}

async function DeleteServer() {
  ElMessageBox.prompt('Are you sure to delete this package', 'Confirm', {
    confirmButtonText: 'OK',
    cancelButtonText: 'Cancel',
    inputPlaceholder: 'input password'
  })
    .then(async ({ value }) => {
      if (value != '0504') {
        return false
      }
      const serverName = state.serverName
      const Data = new FormData()
      Data.append('serverName', serverName)
      const data = await API.DeleteServer(Data)
      if (data.Code) {
        return ElMessage.error('delete error |' + data.Message)
      }
      ElMessage.success('delete success')
      fetchServerList()
      GetMainLogList()
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: 'Delete canceled'
      })
    })
}

async function initLogger() {
  state.loggerList = []
  state.loggerFile = ''
  const formData = new FormData()
  formData.append('serverName', state.serverName)
  const data = await API.GetLogList(formData)
  state.loggerList = data.Data || []
  reverse(state.loggerList)
  state.loggerFile = state.loggerList.length ? state.loggerList[0] : ''
}

onMounted(() => {
  fetchServerList()
  GetMainLogList()
})

watch(
  () => state.activeName,
  async (newVal) => {
    if (newVal == 'api') {
      const formData = new FormData()
      formData.append('serverName', state.serverName)
      const data = await API.GetApiJson(formData)
      state.apis = JSON.parse(data.Data)
    }
    if (newVal == 'doc') {
      const formData = new FormData()
      formData.append('serverName', state.serverName)
      const data = await API.GetDoc(formData)
      state.doc = data.Data.split('\n').reduce(
        (pre: Array<{ hash: string; doc: string }>, curr: string, index: number) => {
          if (index % 2 == 0) {
            pre.push({
              hash: curr,
              doc: ''
            })
          } else {
            pre[pre.length - 1].doc = curr
          }
          return pre
        },
        []
      )
    }
    if (newVal == 'logger') {
      await initLogger()
      // const
    }
  }
)
watch(
  () => state.serverName,
  async function (newVal, oldVal) {
    if (oldVal != newVal) {
      state.doc = []
      state.apis = []
      state.logger = ''
      state.packageList = []
      state.activeName = 'logger'
      state.pattern = ''
      fileList.value = []
      //
      await initLogger()
    }
  }
)

watch(
  () => state.releaseVisible,
  function (newVal) {
    if (newVal) {
      state.selectRelease = ''
      uploadForm.value = {
        serverName: state.serverName,
        file: null,
        doc: ''
      }
    }
  }
)

const multpieNodesState = reactive({
  upstreams: [] as any[],
  selectUpstream: '',
  hosts: [],
  selectHosts: []
})

const SingleNode = 'SingleNode'
watch(
  () => state.selectRelease,
  async function (newVal) {
    if (!newVal) {
      return
    }
    try{
      const data = await getProxyList()
      console.log('data', data.Data)
      multpieNodesState.upstreams = [{ key: SingleNode }].concat(data.Data.upstreams)
    }catch(e){
      ElMessage.error("Error! ExpansionServer is not active ",)
      multpieNodesState.upstreams = [{ key: SingleNode }]
    }
  }
)
watch(
  () => multpieNodesState.selectUpstream,
  async function (newVal) {
    if (!newVal || newVal == SingleNode) {
      multpieNodesState.hosts = []
      multpieNodesState.selectHosts = []
      multpieNodesState.upstreams = []
      return
    }
    multpieNodesState.hosts = []
    multpieNodesState.selectHosts = []
    const hosts = multpieNodesState.upstreams.find(
      (v) => v.key === multpieNodesState.selectUpstream
    ).value.server
    if (hosts instanceof Array) {
      multpieNodesState.hosts = hosts
      multpieNodesState.selectHosts = hosts
    } else {
      multpieNodesState.hosts = [hosts]
      multpieNodesState.selectHosts = [hosts]
    }
  }
)
</script>

<template>
  <div>
    <el-container>
      <el-aside width="200px">
        <aside-component
          :server-list="state.serverList"
          @handle-open="handleOpen"
        ></aside-component>
      </el-aside>
      <el-main>
        <main-logger
          :server-list="state.serverList"
          :create-server-visible="state.configVisible"
          @get-main-log-list="GetMainLogList()"
          @change-visible="
            (bool: boolean) => {
              state.createServerVisible = bool
            }
          "
        ></main-logger>
        <el-card shadow="hover" v-if="!state.serverName">
          <div style="display: flex; height: 700px" v-if="state.activeName == 'logger'">
            <div style="width: 13%; margin-right: 2%">
              <el-select v-model="state.loggerFile">
                <el-option
                  v-for="item in state.mainLogList"
                  :key="item"
                  :value="item"
                  :label="item"
                ></el-option>
              </el-select>
              <br />
              <br />
              <el-select v-model="state.rows" placeholder="select row">
                <el-option
                  v-for="item in rowList"
                  :key="item"
                  :value="item"
                  :label="item"
                ></el-option>
              </el-select>
              <br />
              <br />
              <el-input v-model="state.pattern"></el-input>
              <br />
              <br />
              <el-button @click="GetServerLogger()" type="primary">Search</el-button>
            </div>
            <div
              class="resu"
              style="
                background-color: black;
                height: 700px;
                padding: 5px 10px;
                width: 85%;
                overflow: scroll;
              "
            >
              <div style="color: aliceblue">
                SimpMainControlLogServer :: {{ state.serverName }} :: created By leeks
              </div>
              <div style="color: aliceblue; margin: 2px" v-for="item in state.logger" :key="item">
                {{ item }}
              </div>
            </div>
          </div>
        </el-card>
        <el-card shadow="hover" v-if="state.serverName">
          <div
            style="height: 70px; display: flex; align-items: center; justify-content: space-around"
          >
            <div class="flex-item">
              <div style="font-weight: 700">PackageCounts</div>
              <div style="color: rgb(207, 15, 124)">{{ state.packageList.length }}</div>
            </div>
            <div class="flex-item">
              <div style="font-weight: 700">ServerName</div>
              <div style="color: rgb(207, 15, 124); cursor: pointer" @click="showConfig()">
                {{ state.serverName || '--' }}
              </div>
            </div>
            <div class="flex-item">
              <div style="font-weight: 700">Pid</div>
              <div style="color: rgb(207, 15, 124)">{{ state.status.pid || '--' }}</div>
            </div>
            <div class="flex-item">
              <div style="font-weight: 700">Status</div>
              <div v-html="state.status.status"></div>
            </div>
            <div class="flex-item">
              <div
                style="font-weight: 700; color: rgb(207, 15, 124); cursor: pointer"
                @click="state.expansionVisible = true"
              >
                Expansion
              </div>
            </div>
            <div class="flex-item">
              <div
                @click="state.uploadVisible = true"
                style="color: rgb(207, 15, 124); cursor: pointer"
              >
                Upload
              </div>
            </div>
            <div class="flex-item">
              <div @click="ShutDownServer()" style="color: rgb(207, 15, 124); cursor: pointer">
                Shutdown
              </div>
            </div>
            <div class="flex-item">
              <div
                @click="state.releaseVisible = true"
                style="color: rgb(207, 15, 124); cursor: pointer"
              >
                Release
              </div>
            </div>
            <div class="flex-item">
              <div style="cursor: pointer; color: red" @click="DeleteServer">Delete</div>
            </div>
          </div>
        </el-card>
        <el-card shadow="hover" v-if="state.serverName" style="padding: 0">
          <el-tabs v-model="state.activeName">
            <el-tab-pane label="logger" name="logger">
              <div style="display: flex; height: 700px" v-if="state.activeName == 'logger'">
                <div style="width: 13%; margin-right: 2%">
                  <el-select v-model="state.loggerFile">
                    <el-option
                      v-for="item in state.loggerList"
                      :key="item"
                      :value="item"
                      :label="item"
                    ></el-option>
                  </el-select>
                  <br />
                  <br />
                  <el-select v-model="state.rows" placeholder="select row">
                    <el-option
                      v-for="item in rowList"
                      :key="item"
                      :value="item"
                      :label="item"
                    ></el-option>
                  </el-select>
                  <br />
                  <br />
                  <el-input v-model="state.pattern"></el-input>
                  <br />
                  <br />
                  <el-button @click="GetServerLogger()" type="primary">Search</el-button>
                </div>
                <div
                  class="resu"
                  style="
                    background-color: black;
                    height: 700px;
                    padding: 5px 10px;
                    width: 85%;
                    overflow: scroll;
                  "
                >
                  <div style="color: aliceblue">
                    SimpLogServer :: {{ state.serverName }} :: created By leeks
                  </div>
                  <div
                    style="color: aliceblue; margin: 2px"
                    v-for="item in state.logger"
                    :key="item"
                  >
                    {{ item }}
                  </div>
                </div>
              </div>
            </el-tab-pane>
            <el-tab-pane label="api" name="api">
              <div v-if="state.activeName == 'api'">
                <div v-for="item in state.apis" :key="item">
                  <div
                    style="
                      display: flex;
                      align-items: center;
                      justify-content: space-between;
                      margin: 5px;
                    "
                  >
                    <el-button type="text">{{ item.method }} | {{ item.path }}</el-button>
                    <el-button type="primary">invoke</el-button>
                  </div>
                </div>
              </div>
            </el-tab-pane>
            <el-tab-pane label="doc" name="doc">
              <div v-if="state.activeName == 'doc'">
                <el-table :data="state.doc" stripe style="width: 100%" border>
                  <el-table-column prop="doc" label="doc" width="180" align="center" />
                  <el-table-column prop="hash" label="hash" align="center" />
                </el-table>
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-card>
        <!-- 文件上传表单 -->
        <el-dialog append-to-body v-model="state.uploadVisible" width="50%" title="Release">
          <el-form :model="uploadForm" label-width="150px">
            <el-form-item label="Server Name" required>
              <el-input v-model="uploadForm.serverName" disabled></el-input>
            </el-form-item>
            <el-form-item label="Document" required>
              <el-input v-model="uploadForm.doc" type="textarea" row="5"></el-input>
            </el-form-item>
            <el-form-item label="File" required>
              <el-upload
                :file-list="fileList"
                :show-file-list="true"
                :on-change="handleFileChange"
                :auto-upload="false"
                action="/upload"
              >
                <el-button slot="trigger" size="small">Choose File</el-button>
              </el-upload>
            </el-form-item>
          </el-form>
          <span slot="footer">
            <div style="display: flex; align-items: center; justify-content: center">
              <el-button type="primary" @click="state.uploadVisible = false">Close</el-button>
              <el-button type="success" @click="uploadFile">Upload</el-button>
            </div>
          </span>
        </el-dialog>
        <el-dialog
          append-to-body
          v-model="state.configVisible"
          title="Server Configuration"
          width="60%"
        >
          <el-form :model="config" label-width="100px">
            <el-form-item label="Name">
              <el-input v-model="config.Name" disabled></el-input>
            </el-form-item>
            <el-form-item label="Host">
              <el-input v-model="config.Host"></el-input>
            </el-form-item>
            <el-form-item label="Port">
              <el-input v-model="config.Port" disabled></el-input>
            </el-form-item>
            <el-form-item label="Type">
              <el-input v-model="config.Type" disabled></el-input>
            </el-form-item>
            <el-form-item label="Static Path">
              <el-input v-model="config.StaticPath"></el-input>
            </el-form-item>
            <el-form-item label="Storage">
              <el-input v-model="config.Storage"></el-input>
            </el-form-item>
          </el-form>
          <span slot="footer">
            <div style="display: flex; align-items: center; justify-content: center">
              <el-button type="primary" @click="state.configVisible = false">Close</el-button>
              <el-button type="success" @click="previewSubServer()">Preivew</el-button>
              <el-button type="danger" @click="uploadConfig()">Upload</el-button>
            </div>
          </span>
        </el-dialog>

        <el-dialog
          append-to-body
          v-model="state.createServerVisible"
          title="Create Server"
          width="60%"
        >
          <el-form :model="config" label-width="100px">
            <el-form-item label="Name">
              <el-input v-model="state.createServerName"></el-input>
            </el-form-item>
          </el-form>
          <span slot="footer">
            <div style="display: flex; align-items: center; justify-content: center">
              <el-button type="primary" @click="state.createServerVisible = false">Close</el-button>
              <el-button type="success" @click="createServer()">Create</el-button>
            </div>
          </span>
        </el-dialog>

        <el-dialog append-to-body v-model="state.releaseVisible" title="Release Server" width="60%">
          <el-select v-model="state.selectRelease" placeholder="请选择" style="width: 100%">
            <el-option
              v-for="item in state.packageList"
              :key="item.Hash"
              :label="state.serverName + '_' + item.Hash + '.tar.gz'"
              :value="state.serverName + '_' + item.Hash + '.tar.gz'"
            >
              <span style="float: left">{{ state.serverName }}</span>
              <span style="float: right" @click="DeletePackage(item.Hash)">
                <el-icon style="color: crimson; cursor: pointer"><Remove /></el-icon>
              </span>
              <span style="float: right; color: #8492a6; font-size: 13px">
                {{ item.Hash }} &nbsp;
              </span>
            </el-option>
          </el-select>
          <br />
          <br />
          <el-select
            v-show="state.selectRelease"
            v-model="multpieNodesState.selectUpstream"
            placeholder="请选择"
            style="width: 100%"
          >
            <el-option
              v-for="(item, index) in multpieNodesState.upstreams"
              :key="index"
              :value="item.key"
            >
              <div v-if="item.key == SingleNode" style="color: blue; font-weight: 700">
                Release
              </div>
              <div v-else>{{ item.key }}</div>
            </el-option>
          </el-select>
          <br />
          <br />
          <el-checkbox-group
            v-model="multpieNodesState.selectHosts"
            v-show="state.selectRelease && state.selectRelease != SingleNode"
          >
            <el-checkbox
              v-for="item in multpieNodesState.hosts"
              :label="item"
              :value="item"
              :key="item"
              style="display: block"
            />
          </el-checkbox-group>
          <template #footer>
            <div style="display: flex; align-items: center; justify-content: center">
              <el-button type="primary" @click="state.releaseVisible = false">Close</el-button>
              <el-button type="success" @click="restartServer()">Release</el-button>
              <el-button type="danger" @click="state.uploadVisible = true">Upload</el-button>
            </div>
          </template>
        </el-dialog>
        <expansionComponent
          :expansion-visible="state.expansionVisible"
          :server-name="state.serverName"
          @close-dialog="() => (state.expansionVisible = false)"
          @showReleaseDialog="() => (state.releaseVisible = true)"
        ></expansionComponent>
        <el-dialog v-model="state.shutdownVisible" title="Shutdown Services">
          <el-checkbox-group v-model="choseServices">
            <template v-for="(item, index) in childServiceList" :key="item">
              <el-checkbox
                v-if="childServiceObj[item].status"
                :value="item"
                :label="item + `  | ` + childServiceObj[item].pid + ` | online`"
                style="display: block"
              ></el-checkbox>
            </template>
          </el-checkbox-group>
          <el-button @click="shutdownServers" type="danger">Shutdown</el-button>
        </el-dialog>
        <el-footer>
          <el-divider content-position="center">
            <div style="color: rgb(207, 15, 124); font-size: 18px">Copyright © 2023-2024</div>
          </el-divider>
          <el-divider content-position="center">
            <div style="color: rgb(207, 15, 124); font-size: 18px">
              SimpServer Started on AliCloud Platform
            </div>
          </el-divider>
        </el-footer>
      </el-main>
    </el-container>
  </div>
</template>

<style>
.flex-item {
  text-align: center;
  width: 15%;
  padding: 10px;
}
.el-container {
  min-height: 100vh;
}
</style>
