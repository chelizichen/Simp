import type { BasicResp } from '@/dto/dto'
import HttpReq from '@/utils/exp_req'

export function getProxyList() {
  return HttpReq({
    url: '/getProxyList',
    method: 'get'
  }) as unknown as BasicResp<any>
}

export function nginxExpansion(){
  return HttpReq({
    url: '/nginxExpansion',
    method: 'post'
  }) as unknown as BasicResp<any>
}