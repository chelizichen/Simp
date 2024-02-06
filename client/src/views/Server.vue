<script setup lang="ts">
import { ElIcon, ElLoading, ElMessage, ElPopconfirm } from 'element-plus'
import { onMounted, reactive, ref, watch } from 'vue'
import API from '../api/server'

const state = reactive({
  serverList: [],
  packageList: [],
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
  selectRelease: '',
  logger: '',
  loggerList: [],
  loggerFile: '',
  pattern: '',
  activeName: 'logger',
  apis: [], // API接口 由gin生成
  doc: '',

  mainLogList: []
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
    state.status = {
      status: ret.status ? 'online' : 'offline',
      pid: ret.pid
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
  const formData = new FormData()
  formData.append('fileName', state.selectRelease)
  formData.append('serverName', state.serverName)

  const resp = await API.RestartServer(formData)
  state.status = {
    status: resp.Data.status ? 'online' : 'offline',
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
}

async function GetMainLogList() {
  state.serverName = ''
  const List = await API.GetMainLogList()
  const Data = List.Data
  state.mainLogList = Data
  state.loggerFile = Data.length ? Data[0] : ''
}

async function GetServerLogger() {
  if (state.serverName) {
    const formData = new FormData()
    formData.append('serverName', state.serverName)
    formData.append('fileName', state.loggerFile)
    formData.append('pattern', state.pattern)
    const rest = await API.GetLogger(formData)
    state.logger = rest.Data.split('\n')
  } else {
    const formData = new FormData()
    formData.append('logFile', state.loggerFile)
    formData.append('pattern', state.pattern)
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

async function ShutDownServer() {
  const formData = new FormData()
  formData.append('serverName', state.serverName)
  const data = await API.ShutDownServer(formData)
  if (data.Code) {
    ElMessage({
      type: 'error',
      message: '关闭失败!' + data.Message
    })
    return
  }
  ElMessage({
    type: 'success',
    message: '成功关闭节点!'
  })
  getServerPackageList(state.serverName)
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

function handleFileChange(file: any) {
  uploadForm.value.file = file.raw
}

async function DeletePackage(hash) {
  // ElPopconfirm('确认删除该发布包?', '提示', {
  //   confirmButtonText: '确定',
  //   cancelButtonText: '取消',
  //   type: 'warning'
  // }).then(async () => {
  const formData = new FormData()
  const fileName = `${state.serverName}_${hash}.tar.gz`
  formData.append('serverName', state.serverName)
  formData.append('fileName', fileName)
  const rest = await API.DeletePackage(formData)
  if (rest.Code) {
    ElMessage.error('删除失败' + rest.Message)
    return
  }
  ElMessage.success('删除成功')
  await getServerPackageList(state.serverName)
  state.selectRelease = ''
  // })
  // .catch(() => {
  //   ElMessage.info('已取消删除')
  // })
}

async function initLogger() {
  state.loggerList = []
  state.loggerFile = ''
  const formData = new FormData()
  formData.append('serverName', state.serverName)
  const data = await API.GetLogList(formData)
  state.loggerList = data.Data || []
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
      state.doc = data.Data
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
      state.doc = ''
      state.apis = []
      state.logger = ''
      state.packageList = []
      state.activeName = 'logger'
      state.pattern = ''
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
</script>

<template>
  <div>
    <el-container>
      <el-aside width="200px">
        <h1
          style="color: rgb(207, 15, 124); text-align: center; font-family: fantasy"
          class="app-bigger-size"
        >
        <el-icon style="color: rgb(207, 90, 124); font-size: 36px"><Help /></el-icon>
          Simp
        </h1>
        <el-menu
          class="el-menu-vertical-demo"
          active-text-color="rgb(207, 15, 124)"
          style="border: none"
        >
          <el-menu-item
            v-for="(item, index) in state.serverList"
            class="app-text-center"
            :index="item"
            :key="index"
            @click="handleOpen(item)"
          >
            <el-icon class="app-not-show"><Menu/></el-icon>
            <template #title>{{ item }}</template>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-main>
        <el-card shadow="hover">
          <div
            style="height: 70px; display: flex; align-items: center; justify-content: space-around"
          >
            <div class="flex-item">
              <div
                @click="state.createServerVisible = true"
                style="color: rgb(207, 15, 124); cursor: pointer"
              >
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
              <div style="color: rgb(207, 15, 124)">{{ state.serverList.length }}</div>
            </div>
          </div>
        </el-card>
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
              <div style="color: rgb(207, 15, 124)">{{ state.status.status || '--' }}</div>
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
                <div>
                  {{ state.doc }}
                </div>
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
                <i class="el-icon-remove" style="color: crimson; cursor: pointer"></i>
              </span>
              <span style="float: right; color: #8492a6; font-size: 13px">
                {{ item.Hash }} &nbsp;
              </span>
            </el-option>
          </el-select>
          <template #footer>
            <div style="display: flex; align-items: center; justify-content: center">
              <el-button type="primary" @click="state.releaseVisible = false">Close</el-button>
              <el-button type="success" @click="restartServer()">Release</el-button>
              <el-button type="danger" @click="state.uploadVisible = true">Upload</el-button>
            </div>
          </template>
        </el-dialog>
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
</style>
