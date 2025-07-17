export interface WorkType {
  onWorkTime: string
  offWorkTime: string
  webhookUrl: string
  isSaturdayWork: boolean
}

export interface WorkTime {
  date: string
  workTime1: string
  workTime2: string
  isWeekDay: boolean
  showSelect: boolean
  weekday: number
  dayType: number //0工作日，1休息日，2补班日
  overWorkTimes: string
}

export interface WorkStatics {
  month: string
  overtime: string
  workTime: WorkTime[]
}

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
  name: string
  mac: string
  ip: string
  starTime: string
  hostname: string
  workType: WorkType
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
