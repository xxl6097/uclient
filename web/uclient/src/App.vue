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
                  >webhook设置
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
            sortable
            min-width="130"
          >
            <template #default="props">
              <el-text
                :type="
                  props.row.online
                    ? props.row.nick
                      ? props.row.nick.workType
                        ? props.row.nick.workType.webhookUrl !== ''
                          ? 'warning'
                          : 'success'
                        : 'success'
                      : 'success'
                    : 'none'
                "
                >{{ getClientName(props.row) }}
              </el-text>
            </template>
          </el-table-column>
          <el-table-column prop="ip" label="IP" sortable min-width="135" />
          <el-table-column
            prop="mac"
            label="Mac地址"
            sortable
            v-if="!isMobile()"
          />
          <el-table-column
            prop="starTime"
            label="连接时间"
            sortable
            v-if="!isMobile()"
          >
            <template #default="props">
              {{ formatTimeStamp(props.row.starTime) }}
            </template>
          </el-table-column>
          <el-table-column
            prop="online"
            label="状态"
            sortable
            min-width="80"
            align="center"
          >
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

  <StaticIpListDialog ref="staticIpListDialogRef" />
  <!--  <ClientStaticIpSettingDialog ref="clientStaticIpDialogRef" />-->
  <ClientSettingDialog ref="deviceSettingDialogRef" />
  <UpgradeDialog ref="upgradeRef" />
  <ClientTimeLineDialog ref="clientTimeLineDialogRef" />
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useDark, useToggle } from '@vueuse/core'
import { Client } from './utils/type.ts'
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
  Prompt,
} from './utils/utils.ts'
import { EventAwareSSEClient } from './utils/sseclient.ts'
import ViewExpand from './components/expand/ViewExpand.vue'
import UpgradeDialog from './components/expand/UpgradeDialog.vue'
import StaticIpListDialog from './components/StaticIpListDialog.vue'
import ClientSettingDialog from './components/ClientSettingDialog.vue'
import { ElNotification } from 'element-plus'

const title = ref<string>('客户端列表')
const clientTimeLineDialogRef = ref<InstanceType<
  typeof ClientTimeLineDialog
> | null>(null)
// const clientStaticIpDialogRef = ref<InstanceType<
//   typeof ClientStaticIpSettingDialog
// > | null>(null)
const deviceSettingDialogRef = ref<InstanceType<
  typeof ClientSettingDialog
> | null>(null)
const staticIpListDialogRef = ref<InstanceType<
  typeof StaticIpListDialog
> | null>(null)

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

function renderTable(data: any) {
  tableData.value = data as Client[]
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
          title.value = `客户端列表 ${json.data.appVersion}`
        }
      }
    })
    .catch(() => {
      showErrorTips('失败')
    })
}

const handleWebhookSetting = () => {
  Prompt('请输入WebHook地址', 'webhook设置', '').then((result) => {
    console.log('handleWebhookSetting', result.value)
    const body = {
      webhookUrl: result.value,
    }
    fetch('../api/webhook/set', {
      credentials: 'include',
      method: 'POST',
      body: JSON.stringify(body),
    })
      .then((res) => {
        return res.json()
      })
      .then((json) => {
        if (json) {
          showTips(json.code, json.msg)
        }
      })
      .catch((error) => {
        showErrorTips(`失败:${JSON.stringify(error)}`)
      })
  })
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
        console.log(json)
        renderTable(json.data)
      }
    })
    .catch((error) => {
      console.error(error)
      showErrorTips(`${JSON.stringify(error)}`)
      // renderTable(testTableData)
    })
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
          staticIpListDialogRef.value.showDialogForm(json.data)
        }
      }
      showTips(json.code, json.msg)
    })
    .catch((error) => {
      console.log('error', error)
      showErrorTips(`获取失败${JSON.stringify(error)}`)
    })
    .finally(() => {
      loadings.close()
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
}

const connectSSE = () => {
  try {
    const sseUrl = `../api/client/sse`
    console.log('connectSSE', sseUrl)
    source.value = new EventAwareSSEClient(sseUrl)
    source.value.addEventListener('update', (data) => {
      console.log('update', data)
      renderTable(data)
    })
    source.value.addEventListener('update-one', (data) => {
      console.log('update-one', data)
      showNotifyMessage(data)
    })
    source.value.connect()
  } catch (e) {
    console.error('connectSSE err', e)
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
</style>