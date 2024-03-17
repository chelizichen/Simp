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
