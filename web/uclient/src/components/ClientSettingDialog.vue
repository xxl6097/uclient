<template>
  <div class="dialog-wrapper">
    <el-dialog
      :modal="true"
      :close-on-click-modal="true"
      :close-on-press-escape="true"
      :width="isMobile() ? '80%' : width"
      v-model="formData.show"
      :title="formData.title"
    >
      <template #default>
        <div class="upgrade-popup-content">
          <el-tabs v-model="formData.activeName" @tab-click="handleClick">
            <el-tab-pane label="设置设备名称" name="first">
              <el-form>
                <el-form-item label="设备名称：">
                  <el-input
                    v-model="formData.client.nick.name"
                    placeholder="请输入设备名称"
                  />
                </el-form-item>
                <el-form-item label="状态通知：">
                  <el-checkbox
                    v-model="formData.client.nick.isPush"
                  ></el-checkbox>
                </el-form-item>
                <el-form-item label="周六加班：">
                  <el-checkbox
                    v-model="formData.client.nick.workType.isSaturdayWork"
                  ></el-checkbox>
                </el-form-item>
                <el-form-item label="推送地址：">
                  <el-input
                    v-model="formData.client.nick.workType.webhookUrl"
                  ></el-input>
                </el-form-item>
                <el-form-item label="统计考勤：">
                  <div style="display: flex; align-items: center">
                    <el-time-select
                      v-model="formData.client.nick.workType.onWorkTime"
                      style="width: 140px"
                      :max-time="formData.client.nick.workType.offWorkTime"
                      class="mr-4"
                      format="HH:mm:ss"
                      placeholder="上班考勤"
                      start="07:00:00"
                      step="00:30:00"
                      end="20:00:00"
                    />
                    <el-time-select
                      v-model="formData.client.nick.workType.offWorkTime"
                      style="width: 140px"
                      :min-time="formData.client.nick.workType.onWorkTime"
                      placeholder="下班考勤"
                      format="HH:mm:ss"
                      start="07:00:00"
                      step="00:30:00"
                      end="20:00:00"
                    />
                  </div>
                </el-form-item>
              </el-form>
            </el-tab-pane>
            <!--            <el-tab-pane label="静态IP设置" name="second">-->
            <!--              <el-form label-width="90">-->
            <!--                <el-form-item label="设备名称：">-->
            <!--                  <el-input-->
            <!--                    v-model="formData.second.hostname"-->
            <!--                    placeholder="请输入设备名称"-->
            <!--                    :input-style="-->
            <!--                      formData.client.static ? { color: 'red' } : {}-->
            <!--                    "-->
            <!--                  />-->
            <!--                </el-form-item>-->
            <!--                <el-form-item label="设备Mac：">-->
            <!--                  <el-input-->
            <!--                    v-model="formData.second.mac"-->
            <!--                    placeholder="请输入设备Mac地址"-->
            <!--                  />-->
            <!--                </el-form-item>-->
            <!--                <el-form-item label="设备IP：">-->
            <!--                  <el-input-->
            <!--                    v-model="formData.second.ip"-->
            <!--                    placeholder="请输入设备IP"-->
            <!--                  />-->
            <!--                </el-form-item>-->
            <!--              </el-form>-->
            <!--            </el-tab-pane>-->
            <el-tab-pane label="统计" name="thrid" v-if="isThridShow()">
              <div>
                <div style="margin-bottom: 10px">
                  <el-date-picker
                    v-model="value3"
                    type="datetime"
                    placeholder="Select date and time"
                    value-format="x"
                  />

                  <el-popconfirm
                    title="确定补签上班吗?"
                    @confirm="handleAddWorkTime(true)"
                  >
                    <template #reference>
                      <el-button type="primary" style="margin-left: 10px" plain
                        >补签上班
                      </el-button>
                    </template>
                  </el-popconfirm>

                  <el-popconfirm
                    title="确定补签下班吗?"
                    @confirm="handleAddWorkTime(false)"
                  >
                    <template #reference>
                      <el-button type="primary" style="margin-left: 10px" plain
                        >补签下班
                      </el-button>
                    </template>
                  </el-popconfirm>

                  <el-button
                    type="warning"
                    style="margin-left: 10px"
                    plain
                    @click="fetchWorkData"
                    >刷新
                  </el-button>

                  <el-button
                    type="warning"
                    style="margin-left: 10px"
                    plain
                    @click="fetchWorkEvent"
                    >触发统计
                  </el-button>
                </div>
                <el-table
                  :data="paginatedTableData"
                  border
                  :preserve-expanded-content="false"
                >
                  <el-table-column type="expand">
                    <template #default="props">
                      <div m="4">
                        <el-table
                          :data="props.row.workTime"
                          border
                          :preserve-expanded-content="false"
                        >
                          <el-table-column label="日期" prop="date" sortable />
                          <el-table-column label="上班" prop="workTime1">
                            <template #default="scope">
                              <el-time-picker
                                v-model="scope.row.workTime1"
                                style="width: 100px"
                                arrow-control
                                value-format="HH:mm:ss"
                              />
                            </template>
                          </el-table-column>
                          <el-table-column label="下班" prop="workTime2">
                            <template #default="scope">
                              <el-time-picker
                                v-model="scope.row.workTime2"
                                style="width: 100px"
                                arrow-control
                                value-format="HH:mm:ss"
                              />
                            </template>
                          </el-table-column>
                          <el-table-column
                            label="加班时长"
                            prop="overWorkTimes"
                          />
                          <el-table-column label="星期" prop="weekday">
                            <template #default="scope">
                              {{ getWeekDay(scope.row.weekday) }}
                            </template>
                          </el-table-column>
                          <el-table-column label="类型" prop="isWeekDay">
                            <template #default="scope">
                              <el-tag
                                v-if="!scope.row.showSelect"
                                :type="getTagType(scope.row.dayType)"
                                @dblclick="handleShowSelect(scope.row)"
                                >{{ getTagName(scope.row.dayType) }}
                              </el-tag>
                              <el-select
                                v-else
                                v-model="scope.row.dayType"
                                @change="handleSelectChange(scope.row)"
                              >
                                <el-option
                                  v-for="item in options"
                                  :label="item.label"
                                  :key="item.label"
                                  :value="item.value"
                                />
                              </el-select>
                            </template>
                          </el-table-column>
                          <el-table-column
                            label="操作"
                            max="80"
                            fixed="right"
                            align="center"
                          >
                            <template #default="{ row }">
                              <el-dropdown
                                size="small"
                                split-button
                                type="primary"
                              >
                                <el-popconfirm
                                  title="确定修改吗"
                                  @confirm="handleChangeWorkTime(row)"
                                >
                                  <template #reference> 修改</template>
                                </el-popconfirm>
                                <template #dropdown>
                                  <el-dropdown-menu>
                                    <el-dropdown-item
                                      @click="handleDeleteWorkTime(row)"
                                      >删除
                                    </el-dropdown-item>
                                  </el-dropdown-menu>
                                </template>
                              </el-dropdown>
                            </template>
                          </el-table-column>
                        </el-table>
                      </div>
                    </template>
                  </el-table-column>
                  <el-table-column prop="month" label="月份"></el-table-column>
                  <el-table-column
                    prop="overtime"
                    label="累计时长"
                    align="center"
                  >
                    <template #default="scope">
                      <el-tag type="danger" size="large"
                        >{{ scope.row.overtime }}
                      </el-tag>
                    </template>
                  </el-table-column>
                </el-table>
                <!-- 分页 -->
                <el-pagination
                  style="margin-top: 20px"
                  v-model:page-size="pageSize"
                  v-model:current-page="currentPage"
                  :page-sizes="[10, 20, 50, 100, 1000]"
                  layout="sizes,prev, pager, next"
                  background
                  :size="size"
                  :total="activities.length"
                  @size-change="handleSizeChange"
                  @current-change="handlePageChange"
                />
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
      </template>
      <template #footer v-if="formData.showFooter">
        <el-button
          type="danger"
          v-if="formData.hideErrBtn"
          :loading="formData.deleteloading"
          :loading-icon="Eleme"
          @click="handleDeleteStaticIp"
          >{{ formData.deleteloading ? '删除中...' : '删除' }}
        </el-button>
        <el-button
          type="primary"
          :loading="formData.loading"
          :loading-icon="Eleme"
          @click="handleConfirm"
          >{{ formData.loading ? '设置中...' : '确定' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { Eleme } from '@element-plus/icons-vue'
import { ref, defineExpose, computed } from 'vue'
import {
  Client,
  NickEntry,
  WorkStatics,
  WorkTime,
  WorkType,
} from '../utils/type.ts'
import {
  getWeekDay,
  isMobile,
  showErrorTips,
  showLoading,
  showSucessTips,
  showTips,
} from '../utils/utils.ts'
import { ComponentSize, ElMessageBox, TabsPaneContext } from 'element-plus'

const formData = ref({
  show: false,
  deleteloading: false,
  loading: false,
  activeName: 'first',
  hideErrBtn: false,
  showFooter: true,
  title: '',
  second: {
    hostname: '',
    mac: '',
    ip: '',
  },
  client: {
    hostname: '',
    ip: '',
    mac: '',
    nick: {
      isPush: false,
      workType: {
        onWorkTime: '',
        offWorkTime: '',
        webhookUrl: '',
        isSaturdayWork: false,
      } as WorkType,
    } as NickEntry,
  } as Client,
})
const handleShowSelect = (row: WorkTime) => {
  row.showSelect = !row.showSelect
}
const handleSelectChange = (row: WorkTime) => {
  row.showSelect = false
  // showWarmDialog(`${JSON.stringify(row)}`, {}, {})
}

//activities.length > 0
function isThridShow(): boolean {
  if (
    formData.value &&
    formData.value.client &&
    formData.value.client.nick &&
    formData.value.client.nick.workType &&
    formData.value.client.nick.workType.onWorkTime != '' &&
    formData.value.client.nick.workType.offWorkTime != ''
  ) {
    return true
  }
  return false
}

const getTagType = (value: number) => {
  switch (value) {
    case 0:
      return 'success'
    case 1:
      return 'danger'
    case 2:
      return 'warning'
    case 3:
      return 'warning'
    default:
      return 'primary'
  }
}

const getTagName = (value: number) => {
  switch (value) {
    case 0:
      return '工作日'
    case 1:
      return '休息日'
    case 2:
      return '补班日'
    case 3:
      return '加班日'
    default:
      return '未知'
  }
}
const options = [
  {
    value: 0,
    label: '工作日',
  },
  {
    value: 1,
    label: '节假日',
  },
  {
    value: 2,
    label: '补班日',
  },
  {
    value: 3,
    label: '加班日',
  },
]
const value3 = ref<number>(0)
const width = ref<string>('30%')
const activities = ref<WorkStatics[]>([])
const currentPage = ref<number>(1)
const pageSize = ref<number>(10)
const size = ref<ComponentSize>('default')
// 分页后的表格数
const paginatedTableData = computed<WorkStatics[]>(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return activities.value.slice(start, end)
})
//  分页切换
const handlePageChange = (page: number) => {
  currentPage.value = page
}
const handleSizeChange = (val: number) => {
  console.log(`${val} items per page`)
  pageSize.value = val
}

const handleClick = (tab: TabsPaneContext) => {
  formData.value.hideErrBtn = false
  formData.value.showFooter = true
  console.log('handleClick', tab.paneName)
  width.value = '30%'
  switch (tab.paneName) {
    case 'first':
      break
    case 'second':
      formData.value.hideErrBtn = true
      break
    case 'thrid':
      formData.value.showFooter = false
      width.value = '50%'
      break
  }
}

const handleChangeWorkTime = (row: WorkTime) => {
  const loadings = showLoading('修改中...')
  const body = {
    mac: formData.value.client.mac,
    day: row.date,
    data: row,
  }
  console.log('handleChangeWorkTime', body)
  fetch('../api/work/update', {
    credentials: 'include',
    method: 'POST',
    body: JSON.stringify(body),
  })
    .then((res) => {
      return res.json()
    })
    .then((json) => {
      console.log('handleChangeWorkTime', json)
      showTips(json.code, json.msg)
    })
    .catch((error) => {
      console.log('error', error)
      showErrorTips('修改失败')
    })
    .finally(() => {
      loadings.close()
    })
}

const handleDeleteWorkTime = (row: WorkTime) => {
  const loadings = showLoading('修改中...')
  const body = {
    mac: formData.value.client.mac,
    day: row.date,
  }
  console.log('handleChangeWorkTime', body)
  fetch('../api/work/del', {
    credentials: 'include',
    method: 'POST',
    body: JSON.stringify(body),
  })
    .then((res) => {
      return res.json()
    })
    .then((json) => {
      console.log('handleChangeWorkTime', json)
      showTips(json.code, json.msg)
    })
    .catch((error) => {
      console.log('error', error)
      showErrorTips('修改失败')
    })
    .finally(() => {
      loadings.close()
    })
}
// 默认展开的行
// const defaultExpandedKeys = ref(['1001'])
//
// // 处理展开/折叠事件
// const handleExpandChange = (row: any, expanded: any) => {
//   console.log('展开状态变化:', row, '是否展开:', expanded)
//   // 可以在这里处理展开/折叠时的额外逻辑
//   if (expanded) {
//     // 展开时的操作，如加载子表格数据
//     // loadChildrenData(row.id)
//   }
// }

function initOnWorkTime() {
  formData.value.client = {} as Client
  formData.value.client.nick = {} as NickEntry
  formData.value.client.nick.workType = {
    onWorkTime: '',
    offWorkTime: '',
    webhookUrl: '',
    isSaturdayWork: false,
  } as WorkType
}

function checkOnWorkTime() {
  if (!formData.value.client) {
    formData.value.client = {} as Client
  }
  if (!formData.value.client.nick) {
    formData.value.client.nick = {} as NickEntry
  }
  if (!formData.value.client.nick.workType) {
    formData.value.client.nick.workType = {
      onWorkTime: '',
      offWorkTime: '',
      webhookUrl: '',
      isSaturdayWork: false,
    } as WorkType
  }
}

function handleAddWorkTime(isOnWork: boolean) {
  console.log('handleAddWorkTime', value3.value, formData.value.client.mac)
  const loadings = showLoading('补签申请中...')
  const row = {
    timestamp: value3.value,
    mac: formData.value.client.mac,
    isOnWork: isOnWork,
  }
  fetch('../api/work/add', {
    credentials: 'include',
    method: 'POST',
    body: JSON.stringify(row),
  })
    .then((res) => {
      return res.json()
    })
    .then((json) => {
      console.log('handleAddWorkTime', json)
      showTips(json.code, json.msg)
    })
    .catch((error) => {
      console.log('error', error)
      showErrorTips('补签失败')
    })
    .finally(() => {
      loadings.close()
    })
}

function fetchWorkEvent() {
  const row = {
    mac: formData.value.client.mac,
  }
  fetch('../api/work/tigger', {
    credentials: 'include',
    method: 'POST',
    body: JSON.stringify(row),
  })
    .then((res) => {
      return res.json()
    })
    .then((json) => {
      console.log('fetchWorkEvent', json)
      showTips(json.code, json.msg)
    })
    .catch((error) => {
      console.log('error', error)
    })
    .finally(() => {})
}

function fetchWorkData() {
  const row = {
    mac: formData.value.client.mac,
  }
  fetch('../api/work/get', {
    credentials: 'include',
    method: 'POST',
    body: JSON.stringify(row),
  })
    .then((res) => {
      return res.json()
    })
    .then((json) => {
      console.log('fetchWorkData', json)
      if (json.code === 0 && json.data) {
        activities.value = json.data
      } else {
        //showTips(json.code, json.msg)
      }
    })
    .catch((error) => {
      console.log('error', error)
    })
    .finally(() => {})
}

const handleNickSetting = () => {
  const row = {
    name: formData.value.client.nick.name,
    isPush: formData.value.client.nick.isPush,
    starTime: formData.value.client.starTime,
    mac: formData.value.client.mac,
    ip: formData.value.client.ip,
    hostname: formData.value.client.hostname,
    workType: formData.value.client.nick.workType,
  }
  console.log('handleChangeNickName', row)
  formData.value.loading = true
  fetch('../api/nick/set', {
    credentials: 'include',
    method: 'POST',
    body: JSON.stringify(row),
  })
    .then((res) => {
      return res.json()
    })
    .then((json) => {
      console.log('handleChangeNickName', json)
      showTips(json.code, json.msg)
    })
    .catch((error) => {
      console.log('error', error)
      showErrorTips('修改昵称失败')
    })
    .finally(() => {
      formData.value.loading = false
      setTimeout(function () {
        hideDialog()
      }, 300)
    })
}

const handleDeleteStaticIp = () => {
  console.log('handleDeleteStaticIp', formData.value.client.mac)
  const mac = formData.value.client.mac
  const name = formData.value.client.hostname
  formData.value.deleteloading = true
  ElMessageBox.confirm(`确定删除【${name}】静态IP吗?`, 'Warning', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      const loader = showLoading('删除中...')
      fetch(`../api/staticip/delete?mac=${mac}`, {
        credentials: 'include',
        method: 'DELETE',
      })
        .then((res) => {
          return res.json()
        })
        .then((json) => {
          console.log('handleDeleteStaticIp', json)
          showTips(json.code, json.msg)
        })
        .catch((error) => {
          console.log('error', error)
          showErrorTips('删除失败')
        })
        .finally(() => {
          loader.close()
          formData.value.deleteloading = false
          hideDialog()
        })
    })
    .catch(() => {})
}

function handleConfirm() {
  console.log('handleConfirm')
  switch (formData.value.activeName) {
    case 'first':
      handleNickSetting()
      break
    case 'second':
      handleStaticSet()
      break
  }
}

function handleStaticSet() {
  formData.value.loading = true
  const body = {
    hostname: formData.value.second.hostname,
    ip: formData.value.second.ip,
    mac: formData.value.second.mac,
  }
  console.log('handleStaticSet', body)
  fetch(`../api/staticip/set`, {
    credentials: 'include',
    method: 'POST',
    body: JSON.stringify(body),
  })
    .then((res) => res.json())
    .then((json) => {
      console.log('static ip setting', json)
      if (json && json.code === 0) {
        showSucessTips(json.msg)
      } else {
        showTips(json.code, json.msg)
      }
    })
    .catch((error) => {
      console.log(error)
      showErrorTips(`获取失败${JSON.stringify(error)}`)
    })
    .finally(() => {
      formData.value.loading = false
      setTimeout(function () {
        hideDialog()
      }, 300)
    })
}

const showDialogForm = (row: Client) => {
  console.log('打开对话框，row:', row)
  initOnWorkTime()
  formData.value.title = `设备设置`
  formData.value.client = row
  formData.value.show = true

  formData.value.second.hostname = row.hostname
  formData.value.second.ip = row.ip
  formData.value.second.mac = row.mac
  if (row && row.static) {
    formData.value.second.hostname = row.static.hostname
    formData.value.second.ip = row.static.ip
    formData.value.second.mac = row.static.mac
  }
  // activities.value = testSettingData
  checkOnWorkTime()
  fetchWorkData()
}

function hideDialog() {
  formData.value.show = false
}

// 暴露方法供父组件调用
defineExpose({
  showDialogForm: showDialogForm,
})
</script>
<style scoped>
.upgrade-popup-header h3 {
  line-height: 2.5;
  margin: 0;
}

.upgrade-popup-content {
  height: auto;
  padding: 20px;
}

.upgrade-popup-footer button {
  margin-left: 10px;
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

/* 深度选择器 + 外层容器 */
.dialog-wrapper :deep(.el-dialog__header) {
  display: none;
}
</style>
