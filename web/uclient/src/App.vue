<template>
  <el-progress
    v-if="globalProgress > 0 && globalProgress < 100"
    :percentage="globalProgress"
    :stroke-width="2"
    :show-text="false"
    :color="customColors"
    class="global-progress-bar"
  />
  <div id="app">
    <header class="grid-content header-color">
      <div class="header-content">
        <div class="brand">
          <el-dropdown trigger="click">
            <a href="#">{{ title }}</a>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleShowCheckVersionDialog"
                  >版本检测
                </el-dropdown-item>
                <el-dropdown-item @click="manusForm.show = true"
                  >手动升级
                </el-dropdown-item>
                <el-dropdown-item @click="handleClearData"
                  >清空数据
                </el-dropdown-item>
                <el-dropdown-item @click="handleReboot"
                  >重启应用
                </el-dropdown-item>
                <el-dropdown-item @click="showVersion"
                  >查看版本
                </el-dropdown-item>
                <el-dropdown-item @click="handleResetClients"
                  >重置列表
                </el-dropdown-item>
                <!--                <el-dropdown-item @click="handleAResetNetwork"-->
                <!--                  >重置网络-->
                <!--                </el-dropdown-item>-->
                <el-dropdown-item @click="handleShowStaticIpListDialog"
                  >静态列表
                </el-dropdown-item>
                <el-dropdown-item @click="handleWebhookSetting"
                  >信息配置
                </el-dropdown-item>
                <el-dropdown-item @click="handleGithub"
                  >Github
                </el-dropdown-item>
                <el-dropdown-item @click="showtmp">temp</el-dropdown-item>
                <el-dropdown-item @click="handleTest" v-if="false"
                  >test
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        <div class="dark-switch">
          <el-switch
            v-model="darkmodeSwitch"
            inline-prompt
            active-text="Dark"
            inactive-text="Light"
            @change="toggleDark"
            style="
              --el-switch-on-color: #444452;
              --el-switch-off-color: #589ef8;
            "
          />
        </div>
      </div>
    </header>
    <section>
      <el-main>
        <el-table
          :data="paginatedTableData"
          style="width: 100%"
          :border="true"
          :highlight-current-row="false"
          :preserve-expanded-content="true"
        >
          <el-table-column type="expand">
            <template #default="props">
              <ViewExpand :row="props.row" />
            </template>
          </el-table-column>
          <el-table-column
            prop="hostname"
            label="名称"
            class-name="no-wrap-column"
            :header-cell-class-name="() => 'no-wrap-header'"
            :cell-class-name="() => 'no-wrap-cell'"
            show-overflow-tooltip
            sortable
          >
            <template #default="props">
              <el-text
                :tag="
                  props.row.online
                    ? props.row.nick
                      ? props.row.nick.workType
                        ? props.row.nick.workType.webhookUrl !== ''
                          ? 'ins'
                          : 'p'
                        : 'p'
                      : 'p'
                    : 'p'
                "
                :type="props.row.online ? 'success' : 'danger'"
                >{{ getClientName(props.row) }}
              </el-text>
            </template>
          </el-table-column>
          <el-table-column
            prop="ip"
            label="IP"
            class-name="no-wrap-column"
            :header-cell-class-name="() => 'no-wrap-header'"
            :cell-class-name="() => 'no-wrap-cell'"
            show-overflow-tooltip
            sortable
          />
          <!--          <el-table-column prop="vendor" label="品牌" sortable>-->
          <!--            <template #default="props">-->
          <!--              {{ props.row.vendor }}-->
          <!--            </template>-->
          <!--          </el-table-column>-->
          <!--          <el-table-column prop="staType" label="类型" sortable />-->
          <el-table-column prop="signal" label="信号强度" sortable />
          <el-table-column prop="upRate" label="上传" sortable />
          <el-table-column prop="downRate" label="下载" sortable />
          <!--          <el-table-column-->
          <!--            prop="ssid"-->
          <!--            label="wifi名称"-->
          <!--            sortable-->
          <!--            v-if="!isMobile()"-->
          <!--          />-->
          <el-table-column
            prop="mac"
            label="Mac地址"
            show-overflow-tooltip
            sortable
            v-if="!isMobile()"
          />
          <el-table-column
            prop="starTime"
            label="连接时间"
            class-name="no-wrap-column"
            :header-cell-class-name="() => 'no-wrap-header'"
            :cell-class-name="() => 'no-wrap-cell'"
            show-overflow-tooltip
            sortable
            v-if="!isMobile()"
          >
            <template #default="props">
              <div class="no-wrap">
                {{ formatTimeStamp(props.row.starTime) }}
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="online" label="状态" sortable align="center">
            <template #default="scope">
              <el-tag v-if="scope.row.online" type="success">在线</el-tag>
              <el-tag v-else type="danger">离线</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" max="80" fixed="right" align="center">
            <template #default="{ row }">
              <el-dropdown trigger="click">
                <el-button type="text">菜单</el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <!--                    <el-dropdown-item @click="handleShowStaitcIpDialog(row)"-->
                    <!--                      >静态IP-->
                    <!--                    </el-dropdown-item>-->
                    <el-dropdown-item @click="handleShowDeviceSetting(row)"
                      >设备设置
                    </el-dropdown-item>
                    <el-dropdown-item @click="handleGoToTimeLineDialog(row)"
                      >时间表
                    </el-dropdown-item>
                    <el-dropdown-item @click="handleOfflineDevice(row)"
                      >强制下线
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <el-pagination
          style="margin-top: 20px"
          background
          layout="prev, pager, next"
          :total="filteredTableData.length"
          :page-size="pageSize"
          :current-page="currentPage"
          :pager-count="mobileLayout ? 3 : 7"
          @current-change="handlePageChange"
        />
      </el-main>
    </section>
    <footer></footer>
  </div>

  <!--  客户端程序升级-->
  <el-dialog v-model="manusForm.show" align-center width="500">
    <template #header><span>程序升级</span></template>
    <el-input
      v-model="manusForm.binUrl"
      autocomplete="off"
      placeholder="请输入程序Url地址～"
    />

    <template #footer>
      <div class="dialog-footer">
        <el-upload
          class="upload-demo"
          :http-request="handleUploadUpgradeBin"
          :limit="1"
        >
          <template #trigger>
            <el-button type="primary" :disabled="manusForm.binUrl.length > 0"
              >上传文件升级
            </el-button>
          </template>
          <!-- 添加额外按钮 -->
          <el-button
            style="margin-left: 10px"
            type="danger"
            @click="handleUpdate"
          >
            文件url升级
          </el-button>
        </el-upload>
      </div>
    </template>
  </el-dialog>

  <!-- 弹窗显示版本 -->
  <el-dialog v-model="versionDialogVisible" width="30%">
    <template #header><span>版本信息</span></template>
    <el-descriptions :column="1" :size="size" border>
      <el-descriptions-item width="100">
        <template #label>
          <div class="cell-item">软件名称</div>
        </template>
        {{ version?.appName }}
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">软件版本</div>
        </template>
        {{ version?.appVersion }}
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">编译时间</div>
        </template>
        {{ version?.buildTime }}
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">git版本</div>
        </template>
        {{ version?.gitRevision }}
      </el-descriptions-item>
      <el-descriptions-item>
        <template #label>
          <div class="cell-item">go编译版本</div>
        </template>
        {{ version?.goVersion }}
      </el-descriptions-item>
    </el-descriptions>
  </el-dialog>

  <PushSettingDialog ref="pushDialogRef" />
  <StaticIpListDialog ref="staticIpListDialogRef" />
  <!--  <ClientStaticIpSettingDialog ref="clientStaticIpDialogRef" />-->
  <ClientSettingDialog ref="deviceSettingDialogRef" />
  <UpgradeDialog ref="upgradeRef" />
  <ClientTimeLineDialog ref="clientTimeLineDialogRef" />
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useDark, useToggle } from '@vueuse/core'
import { Client, Version } from './utils/type.ts'
import ClientTimeLineDialog from './components/ClientTimeLineDialog.vue'
import {
  isMobile,
  showErrorTips,
  showLoading,
  showSucessTips,
  showTips,
  showWarmDialog,
  showWarmTips,
  formatTimeStamp,
  xhrPromise,
  formatToUTC8,
} from './utils/utils.ts'
import { EventAwareSSEClient } from './utils/sseclient.ts'
import ViewExpand from './components/expand/ViewExpand.vue'
import UpgradeDialog from './components/expand/UpgradeDialog.vue'
import StaticIpListDialog from './components/StaticIpListDialog.vue'
import ClientSettingDialog from './components/ClientSettingDialog.vue'
import { ComponentSize, ElNotification } from 'element-plus'
import PushSettingDialog from './components/PushSettingDialog.vue'
// import { testTableData } from './utils/data.ts'

const title = ref<string>('客户端列表')
const clientTimeLineDialogRef = ref<InstanceType<
  typeof ClientTimeLineDialog
> | null>(null)
const deviceSettingDialogRef = ref<InstanceType<
  typeof ClientSettingDialog
> | null>(null)

const staticIpListDialogRef = ref<InstanceType<
  typeof StaticIpListDialog
> | null>(null)

const pushDialogRef = ref<InstanceType<typeof PushSettingDialog> | null>(null)

const manusForm = ref({
  show: false,
  binUrl: '',
})
const customColors = [
  { color: '#f56c6c', percentage: 20 },
  { color: '#e6a23c', percentage: 40 },
  { color: '#5cb87a', percentage: 60 },
  { color: '#1989fa', percentage: 80 },
  { color: '#6f7ad3', percentage: 100 },
]

const size = ref<ComponentSize>('default')
const versionDialogVisible = ref(false)
const appinfo = ref<any>()
const globalProgress = ref(0)
const isDark = useDark()
const darkmodeSwitch = ref(isDark)
const toggleDark = useToggle(isDark)
const source = ref<EventAwareSSEClient | null>()
// 搜索关键字
const searchKeyword = ref<string>('')
const pageSize = ref<number>(50)
const currentPage = ref<number>(1)
const tableData = ref<Client[]>([])
const version = ref<Version>()
// 分页后的表格数
const paginatedTableData = computed<Client[]>(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredTableData.value.slice(start, end)
})
// 过滤后的表格数据（根据搜索关键字）
const filteredTableData = computed<Client[]>(() => {
  return tableData.value.filter(() => !searchKeyword.value)
})

function renderTable(data: Client[]) {
  //tableData.value = data as Client[]
  //tableData.value = data
  tableData.value.length = 0
  tableData.value.push(...data)
  // console.log('tableData', tableData)
  // console.log('data', data)
}

function getClientName(row: Client): string {
  // row.nickName === ''
  //   ? row.hostname
  //   : row.hostname === '*'
  //     ? row.nickName
  //     : `${row.hostname}(${row.nickName})`
  if (row.nick) {
    if (row.nick?.name === '') {
      return row.hostname
    } else {
      return row.nick?.name
    }
  } else {
    return row.hostname
  }
}

const getVersion = () => {
  // versionDialogVisible.value = true
  fetch('../api/version', { credentials: 'include', method: 'GET' })
    .then((res) => {
      return res.json()
    })
    .then((json) => {
      if (json && json.code === 0 && json.data) {
        appinfo.value = json.data
        if (json.data && json.data.appVersion) {
          title.value = `uclient ${json.data.hostName} ${json.data.appVersion}`
          document.title = `uclient ${json.data.hostName}`
        }
      }
    })
    .catch(() => {
      showErrorTips('失败')
    })
}
const handleGithub = () => {
  window.open('https://github.com/xxl6097/uclient/releases')
}
const handleTest = () => {
  handleShowCheckVersionDialog()
}
const showtmp = () => {
  const host = window.origin
  window.open(`${host}/tmp/`)
}
const handleWebhookSetting = () => {
  // Prompt('请输入WebHook地址', 'webhook设置', '').then((result) => {
  //   console.log('handleWebhookSetting', result.value)
  //   const body = {
  //     webhookUrl: result.value,
  //   }
  //   fetch('../api/webhook/set', {
  //     credentials: 'include',
  //     method: 'POST',
  //     body: JSON.stringify(body),
  //   })
  //     .then((res) => {
  //       return res.json()
  //     })
  //     .then((json) => {
  //       if (json) {
  //         showTips(json.code, json.msg)
  //       }
  //     })
  //     .catch((error) => {
  //       showErrorTips(`失败:${JSON.stringify(error)}`)
  //     })
  // })
  if (pushDialogRef.value) {
    pushDialogRef.value.showDialogForm()
  }
}

const upgradeRef = ref<InstanceType<typeof UpgradeDialog> | null>(null)

const handleShowCheckVersionDialog = () => {
  if (upgradeRef.value) {
    upgradeRef.value.openUpgradeDialog()
  }
}
// 自定义上传函数
const handleUploadUpgradeBin = (options: any) => {
  const { file } = options
  const formData = new FormData()
  formData.append('file', file)
  const loading = showLoading('程序更新中...')
  globalProgress.value = 0
  manusForm.value.show = false
  xhrPromise({
    url: '../api/upgrade',
    method: 'POST',
    data: formData,
    onUploadProgress: (progress: string) => {
      console.log(`上传进度：${progress}`)
      loading.setText(`程序更新中...${progress}%`)
      globalProgress.value = parseInt(progress)
    },
  })
    .then((data: any) => {
      console.log('请求成功', data)
      // 上传成功的回调
      const json = JSON.parse(data.data)
      if (json.code !== 0) {
        if (json.msg !== '') {
          showErrorTips(json.msg)
        }
      } else {
        if (json.msg !== '') {
          showSucessTips(json.msg)
        }
      }
    })
    .catch((error) => {
      console.error('请求失败', error)
      // 上传失败的回调
      //showErrorTips('上传失败的回调')
    })
    .finally(() => {
      setTimeout(function () {
        loading.close()
        globalProgress.value = 0
        manusForm.value.show = false
        window.location.reload()
      }, 4000)
    })
}

const handleUpdate = () => {
  if (manusForm.value.binUrl.length > 0) {
    const binUrl = manusForm.value.binUrl
    console.log('upgradeByUrl', binUrl)
    const loading = showLoading('程序升级中...')
    manusForm.value.show = false
    fetch('../api/upgrade', {
      credentials: 'include',
      method: 'PUT',
      body: binUrl,
    })
      .then((res) => {
        return res.json()
      })
      .then((json) => {
        showTips(json.code, json.msg)
      })
      .catch(() => {
        showWarmTips('更新失败')
      })
      .finally(() => {
        setTimeout(function () {
          loading.close()
          window.location.reload()
        }, 4000)
      })
  } else {
    showWarmTips('请正确输入url地址')
  }
}

const fetchData = () => {
  const timestamp = 1752266198

  console.log('fetchData', formatToUTC8(timestamp))
  fetch(`../api/clients/get`, {
    credentials: 'include',
    method: 'GET',
  })
    .then((res) => res.json())
    .then((json) => {
      console.log('fetchData', json)
      if (json && json.code === 0 && json.data) {
        console.log('api/clients/get', json)
        renderTable(json.data)
      }
    })
    .catch((error) => {
      console.error(error)
      showErrorTips(`${JSON.stringify(error)}`)
      // renderTable(testTableData)
    })
}
const showVersion = () => {
  // versionDialogVisible.value = true
  fetch('../api/version', { credentials: 'include', method: 'GET' })
    .then((res) => {
      return res.json()
    })
    .then((json) => {
      console.log('showVersion', json)
      if (json.code === 0 && json.data) {
        version.value = json.data
        versionDialogVisible.value = true
      }
      showTips(json.code, json.msg)
    })
    .catch(() => {
      showErrorTips('失败')
    })
}

const handleReboot = () => {
  showWarmDialog(
    `确定重启应用吗？`,
    () => {
      fetch('../api/reboot', { credentials: 'include', method: 'GET' })
        .then((res) => {
          return res.json()
        })
        .then((json) => {
          showTips(json.code, json.msg)
        })
        .catch(() => {
          showErrorTips('重启失败')
        })
    },
    () => {},
  )
}

const handleClearData = () => {
  showWarmDialog(
    `确定清空临时数据吗？`,
    () => {
      fetch('../api/clear', { credentials: 'include', method: 'DELETE' })
        .then((res) => {
          return res.json()
        })
        .then((json) => {
          showTips(json.code, json.msg)
        })
        .catch(() => {
          showErrorTips('清空失败')
        })
    },
    () => {},
  )
}

// const handleChangeNickName = (row: Client) => {
//   console.log('handleChangeNickName', row)
//   ElMessageBox.prompt('请输入设备昵称', '修改昵称', {
//     confirmButtonText: '确定',
//     cancelButtonText: '取消',
//     inputValue: row.nickName,
//   }).then(({ value }) => {
//     row.nickName = value
//     fetch('../api/nick/set', {
//       credentials: 'include',
//       method: 'POST',
//       body: JSON.stringify(row),
//     })
//       .then((res) => {
//         return res.json()
//       })
//       .then((json) => {
//         console.log('handleChangeNickName', json)
//         showTips(json.code, json.msg)
//       })
//       .catch((error) => {
//         console.log('error', error)
//         showErrorTips('修改昵称失败')
//       })
//   })
// }

function handleResetClients() {
  console.log('handleResetClients')
  fetch(`../api/clients/reset`, {
    credentials: 'include',
    method: 'POST',
  })
    .then((res) => res.json())
    .then((json) => {
      console.log('重置列表', json)
      showTips(json.code, json.msg)
    })
    .catch((error) => {
      console.log('error', error)
      showErrorTips(`重置失败${JSON.stringify(error)}`)
    })
}

// function handleAResetNetwork() {
//   showWarmDialog(
//     `确定重置网络吗？`,
//     () => {
//       fetch('../api/network/reset', { credentials: 'include', method: 'POST' })
//         .then((res) => {
//           return res.json()
//         })
//         .then((json) => {
//           showTips(json.code, json.msg)
//         })
//         .catch((error) => {
//           console.log('error', error)
//           showErrorTips(`重置网络失败${JSON.stringify(error)}`)
//         })
//     },
//     () => {},
//   )
// }

function handleShowStaticIpListDialog() {
  console.log('查看静态IP列表')
  const loadings = showLoading('静态IP列表请求中...')
  fetch(`../api/staticip/list`, {
    credentials: 'include',
    method: 'GET',
  })
    .then((res) => res.json())
    .then((json) => {
      console.log('静态IP列表', json)
      if (json && json.code === 0) {
        console.log(json)
        if (staticIpListDialogRef.value && json.data) {
          staticIpListDialogRef.value.showDialogForm(json.data, tableData.value)
        }
      } else {
        if (staticIpListDialogRef.value) {
          staticIpListDialogRef.value.showDialogForm([], tableData.value)
        }
      }
      showTips(json.code, json.msg)
    })
    .catch((error) => {
      console.log('获取失败', error)
      showErrorTips(`获取失败${JSON.stringify(error)}`)
    })
    .finally(() => {
      loadings.close()
      if (staticIpListDialogRef.value) {
        // staticIpListDialogRef.value.showDialogForm(
        //   testStatics as DHCPHost[],
        //   tableData.value,
        // )
      }
    })
}

const handleShowDeviceSetting = (row: Client) => {
  console.log('handleShowDeviceSetting', row)
  if (deviceSettingDialogRef.value) {
    deviceSettingDialogRef.value.showDialogForm(row)
  }
}

// const handleShowStaitcIpDialog = (row: Client) => {
//   console.log('handleShowStaitcIpDialog', row)
//   if (clientStaticIpDialogRef.value) {
//     clientStaticIpDialogRef.value.showDialogForm(row)
//   }
// }

// 调整详情
const handleOfflineDevice = (row: Client) => {
  console.log('handleOfflineDevice', row)
  const body = {
    mac: row.mac,
  }
  fetch('../api/client/offline', {
    credentials: 'include',
    method: 'POST',
    body: JSON.stringify(body),
  })
    .then((res) => {
      return res.json()
    })
    .then((json) => {
      console.log('handleOfflineDevice', json)
      showTips(json.code, json.msg)
    })
    .catch((error) => {
      console.log('error', error)
    })
    .finally(() => {})
}

// 调整详情
const handleGoToTimeLineDialog = (row: Client) => {
  console.log('handleGoToTimeLineDialog', row)
  if (clientTimeLineDialogRef.value) {
    clientTimeLineDialogRef.value.openClientDialog(row)
  }
}
// 分页切换
const handlePageChange = (page: number) => {
  currentPage.value = page
}

// 响应式布局相关
const mobileLayout = ref(false)
const checkMobile = () => {
  mobileLayout.value = window.innerWidth < 768
}

// 弹窗宽度控制
const dialogWidth = ref('500px')
const updateDialogWidth = () => {
  checkMobile()
  dialogWidth.value = mobileLayout.value ? '90%' : '500px'
  if (clientTimeLineDialogRef.value) {
    clientTimeLineDialogRef.value.updateDialogWidth()
  }
  if (upgradeRef.value) {
    upgradeRef.value.updateDialogWidth()
  }
}

const connectSSE = () => {
  try {
    const sseUrl = `../api/client/sse`
    console.log('connectSSE', sseUrl)
    source.value = new EventAwareSSEClient(sseUrl)
    source.value.addEventListener('updateAll', (data) => {
      console.log('updateAll', data)
      renderTable(data)
    })
    source.value.addEventListener('showNotify', (data) => {
      console.log('showNotify', data)
      updateTableByOne(data)
      showNotifyMessage(data)
    })
    source.value.addEventListener('updateOne', (data) => {
      console.log('update-status', data)
      updateTableByOne(data)
    })
    source.value.connect()
  } catch (e) {
    console.error('connectSSE err', e)
  }
}

function updateTableByOne(cls: Client) {
  console.log('updateTableByOne', cls)
  if (tableData.value) {
    // const newTableData = tableData.value.map((item: Client) => {
    //   item.mac === cls.mac ? { ...item, online: cls.online } : item
    // })
    // console.log('updateTableByOne', newTableData)
    // renderTable(newTableData)
    const index = tableData.value.findIndex((item) => item.mac === cls.mac)
    if (index !== -1) {
      tableData.value.forEach((item: Client) => {
        if (item.mac === cls.mac) {
          item.starTime = cls.starTime
          item.online = cls.online
          item.freq = cls.freq
          item.signal = cls.signal
        }
      })
    } else {
      tableData.value.push(cls)
    }
  }
}

function showNotifyMessage(cls: Client) {
  if (cls) {
    let name = ''
    if (cls.nick) {
      if (cls.nick.name !== '') {
        name = cls.nick.name
      }
    }

    if (name === '') {
      if (cls.online) {
        ElNotification({
          title: `未知设备上线了`,
          message: `mac地址:${cls.mac}`,
          type: 'success',
        })
      } else {
        ElNotification({
          title: `未知设备离线了`,
          message: `mac地址:${cls.mac}`,
          type: 'warning',
        })
      }
    } else {
      if (cls.online) {
        ElNotification({
          title: `${name}上线了`,
          message: `mac地址:${cls.mac}`,
          type: 'success',
        })
      } else {
        ElNotification({
          title: `${name}离线了`,
          message: `mac地址:${cls.mac}`,
          type: 'warning',
        })
      }
    }
  }
}

// 初始化监听
onMounted(() => {
  window.addEventListener('resize', updateDialogWidth)
  updateDialogWidth()
})

onUnmounted(() => {
  window.removeEventListener('resize', updateDialogWidth)
})
getVersion()
connectSSE()
fetchData()
</script>

<style>
body {
  margin: 0px;
  font-family:
    -apple-system,
    BlinkMacSystemFont,
    Helvetica Neue,
    sans-serif;
}

header {
  width: 100%;
  height: 60px;
}

.header-color {
  background: #4d83ac;
}

html.dark .header-color {
  background: #395c74;
}

.header-content {
  display: flex;
  align-items: center;
}

#content {
  margin-top: 20px;
  padding-right: 40px;
}

.brand {
  display: flex;
  justify-content: flex-start;
}

.brand a {
  color: #fff;
  background-color: transparent;
  margin-left: 20px;
  line-height: 25px;
  font-size: 25px;
  padding: 15px 15px;
  height: 30px;
  text-decoration: none;
}

.dark-switch {
  display: flex;
  justify-content: flex-end;
  flex-grow: 1;
  padding-right: 40px;
}

.global-progress-bar {
  position: fixed;
  top: 0;
  left: 0;
  z-index: 9999;
  width: 100%;
}

@media screen and (max-width: 1180px) {
  .main-width {
    width: 30%;
  }
}

@media screen and (max-width: 968px) {
  .main-width {
    width: 80%;
  }
}

/* 确保内容不换行 */
.no-wrap-column .cell {
  white-space: nowrap;
}

.no-wrap {
  white-space: nowrap; /* 禁止换行 */
}

.no-wrap-header,
.no-wrap-cell {
  white-space: nowrap !important;
}
</style>
