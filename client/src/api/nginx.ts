import type { BasicResp } from '@/dto/dto'
import HttpReq from '@/utils/exp_req'

export function getProxyList() {
  return HttpReq({
    url: '/getProxyList',
    method: 'get'
  }) as unknown as BasicResp<any>
}

export function nginxExpansion(data) {
  return HttpReq({
    url: '/nginxExpansion',
    method: 'post',
    data
  }) as unknown as BasicResp<any>
}

export function nginxExpansionPreview(data) {
  return HttpReq({
    url: '/nginxExpansionPreview',
    method: 'post',
    data
  }) as unknown as BasicResp<any>
}

export function nginxReload() {
  return HttpReq({
    url: '/nginxReload',
    method: 'post'
  }) as unknown as BasicResp<any>
}

export function getBackupList() {
  return HttpReq({
    url: '/getBackupList',
    method: 'get'
  }) as unknown as BasicResp<any>
}

export function getBackupFile(params: any) {
  return HttpReq({
    url: '/getBackupFile',
    method: 'get',
    params
  }) as unknown as BasicResp<any>
}

export function backupNginx(params: any) {
  return HttpReq({
    url: '/backup',
    method: 'get',
    params
  }) as unknown as BasicResp<any>
}
