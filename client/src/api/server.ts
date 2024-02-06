import HttpReq from '@/utils/request'

const API = {
  UploadServer: function (data, url = '/uploadServer') {
    return HttpReq({
      url,
      data,
      timeout: 10 * 60 * 1000
    })
  },
  RestartServer: function (data, url = '/restartServer') {
    return HttpReq({
      url,
      data
    })
  },
  GetServerList: function (url = '/getServers') {
    return HttpReq({
      url
    })
  },
  GetServerPackageList: function (data, url = '/getServerPackageList') {
    return HttpReq({
      url,
      data
    })
  },
  CheckServer: function (data, url = '/checkServer') {
    return HttpReq({
      url,
      data
    })
  },
  CreateServer: function (data, url = '/createServer') {
    return HttpReq({
      url,
      data
    })
  },
  CheckConfig: function (data, url = '/checkConfig') {
    return HttpReq({
      url,
      data
    })
  },
  DeleteAllPackages: function (data, url = '/deleteAllPackage') {
    return HttpReq({
      url,
      data
    })
  },
  DeletePackage: function (data, url = '/deletePackage') {
    return HttpReq({
      url,
      data
    })
  },
  ShutDownServer: function (data, url = '/shutdownServer') {
    return HttpReq({
      url,
      data
    })
  },
  GetLogger: function (data, url = '/getServerLog') {
    return HttpReq({
      url,
      data
    })
  },
  GetApiJson: function (data, url = '/getApiJson') {
    return HttpReq({
      url,
      data
    })
  },
  GetDoc: function (data, url = '/getDoc') {
    return HttpReq({
      url,
      data
    })
  },
  GetLogList: function (data, url = '/getLogList') {
    return HttpReq({
      url,
      data
    })
  },
  CoverConfig: function (data, url = '/coverConfig') {
    return HttpReq({
      url,
      data
    })
  },
  GetMainLogList: function (data, url = '/main/getLogList') {
    return HttpReq({
      url,
      data
    })
  },
  GetMainLogger: function (data, url = '/main/getServerLog') {
    return HttpReq({
      url,
      data
    })
  }
}

export default API
