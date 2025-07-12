export interface Status {
  timestamp: number
  connected: boolean
  mac: string
}

export interface DHCPHost {
  index: string
  hostname: string
  mac: string
  ip: string
}

export interface Client {
  ip: string
  mac: string
  phy: string
  hostname: string
  nickName: string
  starTime: number
  online: boolean
  statusList: Status[]
}
