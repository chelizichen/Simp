new Vue({
    el: '#app',
    data: {
        uploadForm: {
            serverName: '',
            file: null,
            doc: ''
        },
        serverList: [],
        packageList: [],
        serverName: '',
        uploadVisible: false,
        status: {
            pid: 0,
            status: false,
        },
        configVisible: false,
        config: {
            Name: "",
            Port: 0,
            Type: "",
            StaticPath: "",
            Storage: "",
            Proxy: null,
            Host:"",
        },
        createServerVisible: false,
        createServerName: '',
        releaseVisible: false,
        selectRelease: '',
        logger: '',
        loggerList: [],
        loggerFile: '',
        pattern: '',
        activeName: 'logger',
        apis: '', // API接口 由gin生成
        doc: '',

        mainLogList: [],
    },
    methods: {
        handleFileChange(file) {
            this.uploadForm.file = file.raw;
        },
        async uploadFile() {
            const loading = this.$loading({
                lock: true,
                text: 'Loading',
                spinner: 'el-icon-loading',
                background: 'rgba(0, 0, 0, 0.7)'
            });
            const formData = new FormData();
            formData.append('serverName', this.uploadForm.serverName);
            formData.append('file', this.uploadForm.file);
            formData.append('doc', this.uploadForm.doc);
            const data = await API.UploadServer(formData)
            this.uploadForm.file = null;
            if (data.Code) {
                this.$message({
                    type: 'error',
                    message: '上传失败' + resp.Message
                });
            } else {
                this.$message({
                    type: 'success',
                    message: '上传成功'
                });
                this.uploadVisible = false
            }
            this.getServerPackageList(this.serverName)
            loading.close()
        },
        async restartServer() {
            if (!this.selectRelease || !this.serverName) {
                this.$message({
                    type: 'info',
                    message: '发布失败！请选择指定的服务和发布包'
                });
                return
            }
            const loading = this.$loading({
                lock: true,
                text: 'Loading',
                spinner: 'el-icon-loading',
                background: 'rgba(0, 0, 0, 0.7)'
            });
            const formData = new FormData();
            formData.append('fileName', this.selectRelease);
            formData.append('serverName', this.serverName);

            const resp = await API.RestartServer(formData)
            this.status = {
                status: resp.Data.status ? 'online' : 'offline',
                pid: resp.Data.pid
            }
            if (resp.Code) {
                this.$message({
                    type: 'error',
                    message: '发布失败' + resp.Message
                });
            } else {
                this.$message({
                    type: 'success',
                    message: '发布成功'
                });
                this.releaseVisible = false
            }
            loading.close()
        },
        async fetchServerList() {
            const resp = await API.GetServerList();
            this.serverList = resp.Data || [];
        },
        async getServerPackageList(serverName) {
            const formData = new FormData();
            formData.append('serverName', serverName);
            const vm = this;
            const resp = await API.GetServerPackageList(formData)
            vm.packageList = resp.Data || [];
            if (vm.packageList.length) {
                const resp = await API.CheckServer(formData)
                const ret = resp.Data;
                vm.status = {
                    status: ret.status ? 'online' : 'offline',
                    pid: ret.pid
                }
            }
        },
        async handleOpen(serverName, keyPath) {
            this.serverName = serverName;
            await this.getServerPackageList(serverName)
            this.uploadForm.serverName = serverName;
        },
        async showConfig() {
            const formData = new FormData();
            formData.append('serverName', this.serverName);
            const resp = await API.CheckConfig(formData)
            this.config = resp.Data.Server
            this.configVisible = true
        },
        async DeleteAllPackages() {
            if (!this.serverName) {
                this.$message.error("请先选择服务")
                return
            }
            this.$confirm('此操作将永久删除所有目录, 是否继续?', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(async () => {
                const formData = new FormData();
                formData.append('serverName', this.serverName);
                const data = await API.DeleteAllPackages(formData)
                this.$message({
                    type: 'success',
                    message: '删除成功!'
                });
            }).catch(() => {
                this.$message({
                    type: 'info',
                    message: '已取消删除'
                });
            });
        },
        async createServer() {
            const formData = new FormData();
            formData.append('serverName', this.createServerName);
            const data = await API.CreateServer(formData)
            this.$message({
                type: 'success',
                message: '添加成功!'
            });
            this.createServerVisible = false
            this.createServerName = ''
            this.fetchServerList()
        },
        async DeletePackage(hash) {
            this.$confirm('确认删除该发布包?', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(async () => {
                const formData = new FormData();
                const fileName = `${this.serverName}_${hash}.tar.gz`
                formData.append('serverName', this.serverName);
                formData.append('fileName', fileName);
                const rest = await API.DeletePackage(formData)
                if (rest.Code) {
                    this.$message.error("删除失败" + rest.Message)
                    return
                }
                this.$message.success("删除成功")
                await this.getServerPackageList(this.serverName)
                this.selectRelease = ""
            }).catch(() => {
                this.$message.info("已取消删除")
            })
        },
        async GetServerLogger() {
            if (this.serverName) {
                const formData = new FormData();
                formData.append('serverName', this.serverName);
                formData.append('fileName', this.loggerFile);
                formData.append('pattern', this.pattern);
                const rest = await API.GetLogger(formData)
                this.logger = rest.Data.split("\n");
            } else {
                const formData = new FormData();
                formData.append('logFile', this.loggerFile);
                formData.append('pattern', this.pattern);
                const rest = await API.GetMainLogger(formData)
                this.logger = rest.Data.split("\n");
            }
        },
        async initLogger() {
            this.loggerList = [];
            this.loggerFile = '';
            const formData = new FormData();
            formData.append('serverName', this.serverName);
            const data = await API.GetLogList(formData);
            this.loggerList = data.Data || [];
            this.loggerFile = this.loggerList.length ? this.loggerList[0] : ''
        },
        previewSubServer() {
            const conf = this.config
            const target = `http://${conf.Host}:${conf.Port}/${conf.StaticPath}/`
            // http://localhost:8511/web/
            // http://localhost:8511/web/
            window.open(target)
        },
        async ShutDownServer() {
            const formData = new FormData();
            formData.append('serverName', this.serverName);
            const data = await API.ShutDownServer(formData);
            if (data.Code) {
                this.$message({
                    type: 'error',
                    message: '关闭失败!' + data.Message
                });
                return;
            }
            this.$message({
                type: 'success',
                message: '成功关闭节点!'
            });
            this.getServerPackageList(this.serverName)
        },
        async uploadConfig() {
            const body = this.MapKeysToUpper({
                ServerName: this.serverName,
                Conf: {
                    Server: this.config
                }
            })
            const ret = await API.CoverConfig(body)
            if (ret.Data) {
                this.$message({
                    type: 'error',
                    message: '覆盖失败!' + data.Message
                });
                return;
            }
            this.$message({
                type: 'success',
                message: '覆盖成功!'
            });
            this.configVisible = false
        },
        MapKeysToUpper(obj) {
            var newObj = {};
            for (var key in obj) {
                if (typeof key === 'string') {
                    newObj[key.charAt(0).toUpperCase() + key.slice(1)] = obj[key];
                } else {
                    newObj[key] = obj[key];
                }
            }
            return newObj;
        },
        async GetMainLogList() {
            this.serverName = ""
            const List = await API.GetMainLogList()
            const Data = List.Data;
            this.mainLogList = Data;
            this.loggerFile = Data.length ? Data[0] : ''
        }
    },
    mounted() {
        // Fetch initial server list
        this.fetchServerList();
        this.GetMainLogList()
    },
    watch: {
        activeName: async function (val) {
            if (val == 'api') {
                const formData = new FormData();
                formData.append('serverName', this.serverName);
                const data = await API.GetApiJson(formData)
                this.apis = JSON.parse(data.Data)
            }
            if (val == 'doc') {
                const formData = new FormData();
                formData.append('serverName', this.serverName);
                const data = await API.GetDoc(formData)
                this.doc = data.Data
            }
            if (val == "logger") {
                await this.initLogger()
                // const 
            }
        },
        serverName: async function () {
            this.doc = ''
            this.apis = [];
            this.logger = ''
            this.packageList = []
            this.activeName = 'logger'
            this.pattern = ''
            //
            await this.initLogger()
        },
        releaseVisible:function(newVal){
            if(newVal){
                this.selectRelease = ''
                this.uploadForm= {
                    serverName: this.serverName,
                    file: null,
                    doc: ''
                }
            }
        }
    },
});