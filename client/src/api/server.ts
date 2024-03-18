import type { BasicResp } from '@/dto/dto'
import HttpReq from '@/utils/request'

const API = {
  UploadServer: function (data, url = '/uploadServer') {
    return HttpReq({
      url,
      data,
      timeout: 10 * 60 * 1000
    }) as unknown as BasicResp<any>
  },
  RestartServer: function (data, url = '/restartServer') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  GetServerList: function (url = '/getServers') {
    return HttpReq({
      url
    }) as unknown as BasicResp<any>
  },
  GetServerPackageList: function (data, url = '/getServerPackageList') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  CheckServer: function (data, url = '/checkServer') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  CreateServer: function (data, url = '/createServer') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  CheckConfig: function (data, url = '/checkConfig') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  DeleteAllPackages: function (data, url = '/deleteAllPackage') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  DeletePackage: function (data, url = '/deletePackage') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  ShutDownServer: function (data, url = '/shutdownServer') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  GetLogger: function (data, url = '/getServerLog') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  GetApiJson: function (data, url = '/getApiJson') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  GetDoc: function (data, url = '/getDoc') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  GetLogList: function (data, url = '/getLogList') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  CoverConfig: function (data, url = '/coverConfig') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  DeleteServer: function (data, url = '/deleteServer') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  GetMainLogList: function (data, url = '/main/getLogList') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  GetMainLogger: function (data, url = '/main/getServerLog') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  Login: function (data, url = '/login') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  getChildStats: function (data, url = '/getChildStats') {
    return HttpReq({
      url,
      data
    })
  }
}

export default API
