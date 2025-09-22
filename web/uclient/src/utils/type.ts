export interface WorkType {
  onWorkTime: string
  offWorkTime: string
  webhookUrl: string
  isSaturdayWork: boolean
}

// export interface WorkTime {
//   date: string
//   workTime1: string
//   workTime2: string
//   isWeekDay: boolean
//   showSelect: boolean
//   weekday: number
//   dayType: number //0工作日，1休息日，2补班日
//   overWorkTimes: string
// }
//
// export interface WorkStatics1 {
//   month: string
//   overtime: string
//   workTime: WorkTime[]
// }
export interface DayData {
  date: string
  workTime1: string
  workTime2: string
  isWeekDay: boolean
  showSelect: boolean
  weekday: number
  dayType: number //0工作日，1休息日，2补班日
  // overHours: number
  soverHours: string
}

export interface MonthData {
  month: string
  weekCount: number
  dayCount: number
  // totalOverHours: number
  // workDayOverHours: number
  // workDayAveOverHours: number
  // saturdayOverHours: number
  // saturdayAveOverHours: number
  stotalOverHours: number
  sworkDayOverHours: number
  sworkDayAveOverHours: number
  ssaturdayOverHours: number
  ssaturdayAveOverHours: number
  saturdayCount: string[]
  dayDatas: DayData[]
}

export interface Status {
  timestamp: number
  connected: boolean
}

export interface TimeLine {
  timestamp: number
  dateTime: string
  ago: string
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
  isPush: boolean
  workType: WorkType
}

export interface Client {
  ip: string
  mac: string
  phy: string
  hostname: string
  staType: string
  ssid: string
  upRate: string
  downRate: string
  vendor: string
  signal: number
  freq: number
  nick: NickEntry
  static: DHCPHost
  starTime: number
  online: boolean
  statusList: Status[]
}

// 定义类型化的注入键
export interface Version {
  frpcVersion: string
  appName: string
  appVersion: string
  buildTime: string
  gitRevision: string
  gitBranch: string
  goVersion: string
  displayName: string
  description: string
  osType: string
  arch: string
  compiler: string
  gitTreeState: string
  gitCommit: string
  gitVersion: string
  gitReleaseCommit: string
  binName: string
  totalSize: string
  usedSize: string
  freeSize: string
  hostName: string
}
