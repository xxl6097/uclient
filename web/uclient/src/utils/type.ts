export interface Status {
  timestamp: number
  connected: boolean
}

export interface DHCPHost {
  index: string
  hostname: string
  mac: string
  ip: string
}

export interface NickEntry {
  Name: string
  name: string
  isPush: boolean
  mac: string
  ip: string
  starTime: string
  hostname: string
}

export interface Client {
  ip: string
  mac: string
  phy: string
  hostname: string
  nick: NickEntry
  static: DHCPHost
  starTime: number
  online: boolean
  statusList: Status[]
}
