<template>
  <div class="dialog-wrapper">
    <el-dialog
      :modal="true"
      :close-on-click-modal="true"
      :close-on-press-escape="true"
      :width="isMobile() ? '80%' : '30%'"
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
                    v-model="formData.first.name"
                    placeholder="请输入设备名称"
                  />
                </el-form-item>
                <el-form-item label="推送状态：">
                  <el-checkbox v-model="formData.first.isPush"></el-checkbox>
                </el-form-item>
                <el-form-item label="统计考勤：">
                  <!--                  <el-checkbox v-model="formData.first.isWork"></el-checkbox>-->
                  <div style="display: flex; align-items: center">
                    <el-time-select
                      v-model="formData.client.nick.workType.onWorkTime"
                      style="width: 140px"
                      :max-time="formData.client.nick.workType.offWorkTime"
                      class="mr-4"
                      placeholder="上班考勤"
                      start="06:00"
                      step="00:15"
                      end="20:00"
                    />
                    <el-time-select
                      v-model="formData.client.nick.workType.offWorkTime"
                      style="width: 140px"
                      :min-time="formData.client.nick.workType.onWorkTime"
                      placeholder="下班考勤"
                      start="06:00"
                      step="00:15"
                      end="20:00"
                    />
                  </div>
                </el-form-item>
              </el-form>
            </el-tab-pane>
            <el-tab-pane label="静态IP设置" name="second">
              <el-form label-width="90">
                <el-form-item label="设备名称：">
                  <el-input
                    v-model="formData.second.hostname"
                    placeholder="请输入设备名称"
                    :input-style="
                      formData.client.static ? { color: 'red' } : {}
                    "
                  />
                </el-form-item>
                <el-form-item label="设备Mac：">
                  <el-input
                    v-model="formData.second.mac"
                    placeholder="请输入设备Mac地址"
                  />
                </el-form-item>
                <el-form-item label="设备IP：">
                  <el-input
                    v-model="formData.second.ip"
                    placeholder="请输入设备IP"
                  />
                </el-form-item>
              </el-form>
            </el-tab-pane>
            <el-tab-pane label="统计" name="thrid">
              <div>
                <div style="margin-bottom: 10px">
                  <el-date-picker
                    v-model="value3"
                    type="datetime"
                    placeholder="Select date and time"
                    value-format="x"
                  />

                  <el-popconfirm title="确定补签吗?">
                    <template #reference>
                      <el-button
                        type="primary"
                        style="margin-left: 10px"
                        @click="handleChangeWorkTime"
                        >补签
                      </el-button>
                    </template>
                  </el-popconfirm>
                </div>
                <el-table :data="paginatedTableData" border>
                  <el-table-column type="expand">
                    <template #default="props">
                      <div m="4">
                        <el-table :data="props.row.workTime" border>
                          <el-table-column label="日期" prop="date" />
                          <el-table-column label="上班" prop="workTime1" />
                          <el-table-column label="下班" prop="workTime2" />
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
                              <el-tag v-if="scope.row.isWeekDay" type="danger"
                                >节假日
                              </el-tag>
                              <el-tag v-else type="success">工作日</el-tag>
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
          :loading="formData.loading"
          :loading-icon="Eleme"
          @click="handleDeleteStaticIp"
          >{{ formData.loading ? '删除中...' : '删除' }}
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
import { Client, WorkStatics } from '../utils/type.ts'
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
  loading: false,
  activeName: 'first',
  hideErrBtn: false,
  showFooter: true,
  title: '',
  first: {
    name: '',
    isPush: false,
  },
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
      workType: {
        onWorkTime: '',
        offWorkTime: '',
      },
    },
  } as Client,
})

const value3 = ref<number>(0)
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
  switch (tab.paneName) {
    case 'first':
      break
    case 'second':
      formData.value.hideErrBtn = true
      break
    case 'thrid':
      formData.value.showFooter = false
      break
  }
}

function handleChangeWorkTime() {
  console.log('handleChangeWorkTime', value3.value, formData.value.client.mac)
  const loadings = showLoading('补签申请中...')
  const row = {
    timestamp: value3.value,
    mac: formData.value.client.mac,
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
      console.log('handleChangeWorkTime', json)
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

const handleChangeNickName = () => {
  const row = {
    isPush: formData.value.first.isPush,
    name: formData.value.first.name,
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
  formData.value.loading = true
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
          formData.value.loading = false
          hideDialog()
        })
    })
    .catch(() => {})
}

function handleConfirm() {
  console.log('handleConfirm')
  switch (formData.value.activeName) {
    case 'first':
      handleChangeNickName()
      break
    case 'second':
      handleStaticSet()
      break
  }
}

function handleStaticSet() {
  formData.value.loading = true
  fetch(`../api/staticip/set`, {
    credentials: 'include',
    method: 'POST',
    body: JSON.stringify(formData.value.client),
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
  formData.value.title = `设备设置`
  // formData.value.client = JSON.parse(JSON.stringify(row))
  formData.value.client = row
  formData.value.show = true

  formData.value.first.isPush = true
  formData.value.first.name = row.hostname
  if (row && row.nick) {
    formData.value.first.isPush = row.nick.isPush
    formData.value.first.name = row.nick.name
  }
  formData.value.second.hostname = row.hostname
  formData.value.second.ip = row.ip
  formData.value.second.mac = row.mac
  if (row && row.static) {
    formData.value.second.hostname = row.static.hostname
    formData.value.second.ip = row.static.ip
    formData.value.second.mac = row.static.mac
  }

  activities.value = testData
}

function hideDialog() {
  formData.value.show = false
}

const testData: WorkStatics[] = [
  {
    month: '2025-01',
    overtime: '47h6m26s',
    workTime: [
      {
        date: '2025-01-01',
        weekday: 3,
        workTime1: '08:57:27',
        workTime2: '19:37:04',
        overWorkTimes: '1h9m37s',
        isWeekDay: false,
      },
      {
        date: '2025-01-02',
        weekday: 4,
        workTime1: '08:14:44',
        workTime2: '19:04:15',
        overWorkTimes: '1h19m31s',
        isWeekDay: false,
      },
      {
        date: '2025-01-03',
        weekday: 5,
        workTime1: '08:00:56',
        workTime2: '19:34:36',
        overWorkTimes: '2h3m40s',
        isWeekDay: false,
      },
      {
        date: '2025-01-04',
        weekday: 6,
        workTime1: '08:09:30',
        workTime2: '19:07:48',
        overWorkTimes: '1h28m18s',
        isWeekDay: false,
      },
      {
        date: '2025-01-05',
        weekday: 0,
        workTime1: '08:19:24',
        workTime2: '19:10:21',
        overWorkTimes: '1h20m57s',
        isWeekDay: false,
      },
      {
        date: '2025-01-06',
        weekday: 1,
        workTime1: '08:14:41',
        workTime2: '19:46:22',
        overWorkTimes: '2h1m41s',
        isWeekDay: false,
      },
      {
        date: '2025-01-07',
        weekday: 2,
        workTime1: '08:00:38',
        workTime2: '19:21:45',
        overWorkTimes: '1h51m7s',
        isWeekDay: false,
      },
      {
        date: '2025-01-08',
        weekday: 3,
        workTime1: '08:52:29',
        workTime2: '19:17:09',
        overWorkTimes: '54m40s',
        isWeekDay: false,
      },
      {
        date: '2025-01-09',
        weekday: 4,
        workTime1: '08:13:38',
        workTime2: '19:05:55',
        overWorkTimes: '1h22m17s',
        isWeekDay: false,
      },
      {
        date: '2025-01-10',
        weekday: 5,
        workTime1: '08:18:04',
        workTime2: '19:03:40',
        overWorkTimes: '1h15m36s',
        isWeekDay: false,
      },
      {
        date: '2025-01-11',
        weekday: 6,
        workTime1: '08:01:39',
        workTime2: '19:57:36',
        overWorkTimes: '2h25m57s',
        isWeekDay: false,
      },
      {
        date: '2025-01-12',
        weekday: 0,
        workTime1: '08:38:15',
        workTime2: '19:23:07',
        overWorkTimes: '1h14m52s',
        isWeekDay: false,
      },
      {
        date: '2025-01-13',
        weekday: 1,
        workTime1: '08:02:16',
        workTime2: '19:38:57',
        overWorkTimes: '2h6m41s',
        isWeekDay: false,
      },
      {
        date: '2025-01-14',
        weekday: 2,
        workTime1: '08:43:39',
        workTime2: '19:08:04',
        overWorkTimes: '54m25s',
        isWeekDay: false,
      },
      {
        date: '2025-01-15',
        weekday: 3,
        workTime1: '08:45:49',
        workTime2: '19:14:08',
        overWorkTimes: '58m19s',
        isWeekDay: false,
      },
      {
        date: '2025-01-16',
        weekday: 4,
        workTime1: '08:33:07',
        workTime2: '19:03:34',
        overWorkTimes: '1h0m27s',
        isWeekDay: false,
      },
      {
        date: '2025-01-17',
        weekday: 5,
        workTime1: '08:02:53',
        workTime2: '19:04:51',
        overWorkTimes: '1h31m58s',
        isWeekDay: false,
      },
      {
        date: '2025-01-18',
        weekday: 6,
        workTime1: '08:35:48',
        workTime2: '19:17:25',
        overWorkTimes: '1h11m37s',
        isWeekDay: false,
      },
      {
        date: '2025-01-19',
        weekday: 0,
        workTime1: '08:39:49',
        workTime2: '19:59:46',
        overWorkTimes: '1h49m57s',
        isWeekDay: false,
      },
      {
        date: '2025-01-20',
        weekday: 1,
        workTime1: '08:16:16',
        workTime2: '19:54:53',
        overWorkTimes: '2h8m37s',
        isWeekDay: false,
      },
      {
        date: '2025-01-21',
        weekday: 2,
        workTime1: '08:26:48',
        workTime2: '19:47:15',
        overWorkTimes: '1h50m27s',
        isWeekDay: false,
      },
      {
        date: '2025-01-22',
        weekday: 3,
        workTime1: '08:30:05',
        workTime2: '19:28:10',
        overWorkTimes: '1h28m5s',
        isWeekDay: false,
      },
      {
        date: '2025-01-23',
        weekday: 4,
        workTime1: '08:12:44',
        workTime2: '19:53:41',
        overWorkTimes: '2h10m57s',
        isWeekDay: false,
      },
      {
        date: '2025-01-24',
        weekday: 5,
        workTime1: '08:41:36',
        workTime2: '19:25:31',
        overWorkTimes: '1h13m55s',
        isWeekDay: false,
      },
      {
        date: '2025-01-25',
        weekday: 6,
        workTime1: '08:10:36',
        workTime2: '19:10:57',
        overWorkTimes: '1h30m21s',
        isWeekDay: false,
      },
      {
        date: '2025-01-26',
        weekday: 0,
        workTime1: '08:11:17',
        workTime2: '19:20:51',
        overWorkTimes: '1h39m34s',
        isWeekDay: false,
      },
      {
        date: '2025-01-27',
        weekday: 1,
        workTime1: '08:00:09',
        workTime2: '19:07:24',
        overWorkTimes: '1h37m15s',
        isWeekDay: false,
      },
      {
        date: '2025-01-28',
        weekday: 2,
        workTime1: '08:06:09',
        workTime2: '19:02:31',
        overWorkTimes: '1h26m22s',
        isWeekDay: false,
      },
      {
        date: '2025-01-29',
        weekday: 3,
        workTime1: '08:22:22',
        workTime2: '19:14:37',
        overWorkTimes: '1h22m15s',
        isWeekDay: false,
      },
      {
        date: '2025-01-30',
        weekday: 4,
        workTime1: '08:33:41',
        workTime2: '19:07:42',
        overWorkTimes: '1h4m1s',
        isWeekDay: false,
      },
      {
        date: '2025-01-31',
        weekday: 5,
        workTime1: '08:19:54',
        workTime2: '19:22:54',
        overWorkTimes: '1h33m0s',
        isWeekDay: false,
      },
    ],
  },
  {
    month: '2025-02',
    overtime: '43h35m47s',
    workTime: [
      {
        date: '2025-02-01',
        weekday: 6,
        workTime1: '08:44:55',
        workTime2: '19:03:32',
        overWorkTimes: '48m37s',
        isWeekDay: false,
      },
      {
        date: '2025-02-02',
        weekday: 0,
        workTime1: '08:20:36',
        workTime2: '19:35:14',
        overWorkTimes: '1h44m38s',
        isWeekDay: false,
      },
      {
        date: '2025-02-03',
        weekday: 1,
        workTime1: '08:52:35',
        workTime2: '19:59:02',
        overWorkTimes: '1h36m27s',
        isWeekDay: false,
      },
      {
        date: '2025-02-04',
        weekday: 2,
        workTime1: '08:41:53',
        workTime2: '19:20:23',
        overWorkTimes: '1h8m30s',
        isWeekDay: false,
      },
      {
        date: '2025-02-05',
        weekday: 3,
        workTime1: '08:30:29',
        workTime2: '19:14:03',
        overWorkTimes: '1h13m34s',
        isWeekDay: false,
      },
      {
        date: '2025-02-06',
        weekday: 4,
        workTime1: '08:25:08',
        workTime2: '19:37:30',
        overWorkTimes: '1h42m22s',
        isWeekDay: false,
      },
      {
        date: '2025-02-07',
        weekday: 5,
        workTime1: '08:45:58',
        workTime2: '19:23:51',
        overWorkTimes: '1h7m53s',
        isWeekDay: false,
      },
      {
        date: '2025-02-08',
        weekday: 6,
        workTime1: '08:00:17',
        workTime2: '19:44:16',
        overWorkTimes: '2h13m59s',
        isWeekDay: false,
      },
      {
        date: '2025-02-09',
        weekday: 0,
        workTime1: '08:21:03',
        workTime2: '19:42:14',
        overWorkTimes: '1h51m11s',
        isWeekDay: false,
      },
      {
        date: '2025-02-10',
        weekday: 1,
        workTime1: '08:21:37',
        workTime2: '19:48:44',
        overWorkTimes: '1h57m7s',
        isWeekDay: false,
      },
      {
        date: '2025-02-11',
        weekday: 2,
        workTime1: '08:37:19',
        workTime2: '19:03:15',
        overWorkTimes: '55m56s',
        isWeekDay: false,
      },
      {
        date: '2025-02-12',
        weekday: 3,
        workTime1: '08:11:48',
        workTime2: '19:41:16',
        overWorkTimes: '1h59m28s',
        isWeekDay: false,
      },
      {
        date: '2025-02-13',
        weekday: 4,
        workTime1: '08:22:23',
        workTime2: '19:44:46',
        overWorkTimes: '1h52m23s',
        isWeekDay: false,
      },
      {
        date: '2025-02-14',
        weekday: 5,
        workTime1: '08:44:14',
        workTime2: '19:19:09',
        overWorkTimes: '1h4m55s',
        isWeekDay: false,
      },
      {
        date: '2025-02-15',
        weekday: 6,
        workTime1: '08:40:27',
        workTime2: '19:41:14',
        overWorkTimes: '1h30m47s',
        isWeekDay: false,
      },
      {
        date: '2025-02-16',
        weekday: 0,
        workTime1: '08:22:07',
        workTime2: '19:01:03',
        overWorkTimes: '1h8m56s',
        isWeekDay: false,
      },
      {
        date: '2025-02-17',
        weekday: 1,
        workTime1: '08:07:07',
        workTime2: '19:02:19',
        overWorkTimes: '1h25m12s',
        isWeekDay: false,
      },
      {
        date: '2025-02-18',
        weekday: 2,
        workTime1: '08:33:05',
        workTime2: '19:36:19',
        overWorkTimes: '1h33m14s',
        isWeekDay: false,
      },
      {
        date: '2025-02-19',
        weekday: 3,
        workTime1: '08:08:56',
        workTime2: '19:42:46',
        overWorkTimes: '2h3m50s',
        isWeekDay: false,
      },
      {
        date: '2025-02-20',
        weekday: 4,
        workTime1: '08:04:17',
        workTime2: '19:35:57',
        overWorkTimes: '2h1m40s',
        isWeekDay: false,
      },
      {
        date: '2025-02-21',
        weekday: 5,
        workTime1: '08:05:30',
        workTime2: '19:18:17',
        overWorkTimes: '1h42m47s',
        isWeekDay: false,
      },
      {
        date: '2025-02-22',
        weekday: 6,
        workTime1: '08:04:03',
        workTime2: '19:58:00',
        overWorkTimes: '2h23m57s',
        isWeekDay: false,
      },
      {
        date: '2025-02-23',
        weekday: 0,
        workTime1: '08:56:17',
        workTime2: '19:20:29',
        overWorkTimes: '54m12s',
        isWeekDay: false,
      },
      {
        date: '2025-02-24',
        weekday: 1,
        workTime1: '08:39:40',
        workTime2: '19:58:18',
        overWorkTimes: '1h48m38s',
        isWeekDay: false,
      },
      {
        date: '2025-02-25',
        weekday: 2,
        workTime1: '08:04:05',
        workTime2: '19:23:03',
        overWorkTimes: '1h48m58s',
        isWeekDay: false,
      },
      {
        date: '2025-02-26',
        weekday: 3,
        workTime1: '08:33:13',
        workTime2: '19:07:38',
        overWorkTimes: '1h4m25s',
        isWeekDay: false,
      },
      {
        date: '2025-02-27',
        weekday: 4,
        workTime1: '08:26:29',
        workTime2: '19:05:20',
        overWorkTimes: '1h8m51s',
        isWeekDay: false,
      },
      {
        date: '2025-02-28',
        weekday: 5,
        workTime1: '08:05:01',
        workTime2: '19:18:21',
        overWorkTimes: '1h43m20s',
        isWeekDay: false,
      },
    ],
  },
  {
    month: '2025-03',
    overtime: '41h56m33s',
    workTime: [
      {
        date: '2025-03-01',
        weekday: 6,
        workTime1: '08:56:09',
        workTime2: '19:42:34',
        overWorkTimes: '1h16m25s',
        isWeekDay: false,
      },
      {
        date: '2025-03-02',
        weekday: 0,
        workTime1: '08:34:59',
        workTime2: '19:13:03',
        overWorkTimes: '1h8m4s',
        isWeekDay: false,
      },
      {
        date: '2025-03-03',
        weekday: 1,
        workTime1: '08:58:39',
        workTime2: '19:25:57',
        overWorkTimes: '57m18s',
        isWeekDay: false,
      },
      {
        date: '2025-03-04',
        weekday: 2,
        workTime1: '08:27:38',
        workTime2: '19:15:47',
        overWorkTimes: '1h18m9s',
        isWeekDay: false,
      },
      {
        date: '2025-03-05',
        weekday: 3,
        workTime1: '08:49:53',
        workTime2: '19:41:18',
        overWorkTimes: '1h21m25s',
        isWeekDay: false,
      },
      {
        date: '2025-03-06',
        weekday: 4,
        workTime1: '08:28:05',
        workTime2: '19:12:48',
        overWorkTimes: '1h14m43s',
        isWeekDay: false,
      },
      {
        date: '2025-03-07',
        weekday: 5,
        workTime1: '08:40:56',
        workTime2: '19:21:41',
        overWorkTimes: '1h10m45s',
        isWeekDay: false,
      },
      {
        date: '2025-03-08',
        weekday: 6,
        workTime1: '08:56:21',
        workTime2: '19:10:14',
        overWorkTimes: '43m53s',
        isWeekDay: false,
      },
      {
        date: '2025-03-09',
        weekday: 0,
        workTime1: '08:12:14',
        workTime2: '19:21:09',
        overWorkTimes: '1h38m55s',
        isWeekDay: false,
      },
      {
        date: '2025-03-10',
        weekday: 1,
        workTime1: '08:05:31',
        workTime2: '19:22:38',
        overWorkTimes: '1h47m7s',
        isWeekDay: false,
      },
      {
        date: '2025-03-11',
        weekday: 2,
        workTime1: '08:38:33',
        workTime2: '19:21:14',
        overWorkTimes: '1h12m41s',
        isWeekDay: false,
      },
      {
        date: '2025-03-12',
        weekday: 3,
        workTime1: '08:01:17',
        workTime2: '19:04:02',
        overWorkTimes: '1h32m45s',
        isWeekDay: false,
      },
      {
        date: '2025-03-13',
        weekday: 4,
        workTime1: '08:35:45',
        workTime2: '19:51:54',
        overWorkTimes: '1h46m9s',
        isWeekDay: false,
      },
      {
        date: '2025-03-14',
        weekday: 5,
        workTime1: '08:02:05',
        workTime2: '19:20:40',
        overWorkTimes: '1h48m35s',
        isWeekDay: false,
      },
      {
        date: '2025-03-15',
        weekday: 6,
        workTime1: '08:55:23',
        workTime2: '19:34:32',
        overWorkTimes: '1h9m9s',
        isWeekDay: false,
      },
      {
        date: '2025-03-16',
        weekday: 0,
        workTime1: '08:42:57',
        workTime2: '19:49:45',
        overWorkTimes: '1h36m48s',
        isWeekDay: false,
      },
      {
        date: '2025-03-17',
        weekday: 1,
        workTime1: '08:17:11',
        workTime2: '19:00:15',
        overWorkTimes: '1h13m4s',
        isWeekDay: false,
      },
      {
        date: '2025-03-18',
        weekday: 2,
        workTime1: '08:22:48',
        workTime2: '19:21:53',
        overWorkTimes: '1h29m5s',
        isWeekDay: false,
      },
      {
        date: '2025-03-19',
        weekday: 3,
        workTime1: '08:28:50',
        workTime2: '19:29:59',
        overWorkTimes: '1h31m9s',
        isWeekDay: false,
      },
      {
        date: '2025-03-20',
        weekday: 4,
        workTime1: '08:55:51',
        workTime2: '19:04:30',
        overWorkTimes: '38m39s',
        isWeekDay: false,
      },
      {
        date: '2025-03-21',
        weekday: 5,
        workTime1: '08:53:57',
        workTime2: '19:24:00',
        overWorkTimes: '1h0m3s',
        isWeekDay: false,
      },
      {
        date: '2025-03-22',
        weekday: 6,
        workTime1: '08:35:56',
        workTime2: '19:51:52',
        overWorkTimes: '1h45m56s',
        isWeekDay: false,
      },
      {
        date: '2025-03-23',
        weekday: 0,
        workTime1: '08:16:22',
        workTime2: '19:04:02',
        overWorkTimes: '1h17m40s',
        isWeekDay: false,
      },
      {
        date: '2025-03-24',
        weekday: 1,
        workTime1: '08:07:50',
        workTime2: '19:24:22',
        overWorkTimes: '1h46m32s',
        isWeekDay: false,
      },
      {
        date: '2025-03-25',
        weekday: 2,
        workTime1: '08:39:22',
        workTime2: '19:03:48',
        overWorkTimes: '54m26s',
        isWeekDay: false,
      },
      {
        date: '2025-03-26',
        weekday: 3,
        workTime1: '08:53:32',
        workTime2: '19:26:18',
        overWorkTimes: '1h2m46s',
        isWeekDay: false,
      },
      {
        date: '2025-03-27',
        weekday: 4,
        workTime1: '08:00:45',
        workTime2: '19:01:42',
        overWorkTimes: '1h30m57s',
        isWeekDay: false,
      },
      {
        date: '2025-03-28',
        weekday: 5,
        workTime1: '08:38:01',
        workTime2: '19:20:25',
        overWorkTimes: '1h12m24s',
        isWeekDay: false,
      },
      {
        date: '2025-03-29',
        weekday: 6,
        workTime1: '08:19:12',
        workTime2: '19:31:03',
        overWorkTimes: '1h41m51s',
        isWeekDay: false,
      },
      {
        date: '2025-03-30',
        weekday: 0,
        workTime1: '08:08:06',
        workTime2: '19:20:34',
        overWorkTimes: '1h42m28s',
        isWeekDay: false,
      },
      {
        date: '2025-03-31',
        weekday: 1,
        workTime1: '08:28:39',
        workTime2: '19:25:21',
        overWorkTimes: '1h26m42s',
        isWeekDay: false,
      },
    ],
  },
  {
    month: '2025-04',
    overtime: '44h46m30s',
    workTime: [
      {
        date: '2025-04-01',
        weekday: 2,
        workTime1: '08:20:53',
        workTime2: '19:26:24',
        overWorkTimes: '1h35m31s',
        isWeekDay: false,
      },
      {
        date: '2025-04-02',
        weekday: 3,
        workTime1: '08:29:32',
        workTime2: '19:33:47',
        overWorkTimes: '1h34m15s',
        isWeekDay: false,
      },
      {
        date: '2025-04-03',
        weekday: 4,
        workTime1: '08:30:00',
        workTime2: '19:38:43',
        overWorkTimes: '1h38m43s',
        isWeekDay: false,
      },
      {
        date: '2025-04-04',
        weekday: 5,
        workTime1: '08:55:05',
        workTime2: '19:07:36',
        overWorkTimes: '42m31s',
        isWeekDay: false,
      },
      {
        date: '2025-04-05',
        weekday: 6,
        workTime1: '08:02:30',
        workTime2: '19:13:47',
        overWorkTimes: '1h41m17s',
        isWeekDay: false,
      },
      {
        date: '2025-04-06',
        weekday: 0,
        workTime1: '08:17:08',
        workTime2: '19:48:10',
        overWorkTimes: '2h1m2s',
        isWeekDay: false,
      },
      {
        date: '2025-04-07',
        weekday: 1,
        workTime1: '08:26:10',
        workTime2: '19:45:21',
        overWorkTimes: '1h49m11s',
        isWeekDay: false,
      },
      {
        date: '2025-04-08',
        weekday: 2,
        workTime1: '08:11:55',
        workTime2: '19:22:31',
        overWorkTimes: '1h40m36s',
        isWeekDay: false,
      },
      {
        date: '2025-04-09',
        weekday: 3,
        workTime1: '08:41:38',
        workTime2: '19:50:53',
        overWorkTimes: '1h39m15s',
        isWeekDay: false,
      },
      {
        date: '2025-04-10',
        weekday: 4,
        workTime1: '08:55:38',
        workTime2: '19:17:12',
        overWorkTimes: '51m34s',
        isWeekDay: false,
      },
      {
        date: '2025-04-11',
        weekday: 5,
        workTime1: '08:18:23',
        workTime2: '19:50:01',
        overWorkTimes: '2h1m38s',
        isWeekDay: false,
      },
      {
        date: '2025-04-12',
        weekday: 6,
        workTime1: '08:15:43',
        workTime2: '19:02:18',
        overWorkTimes: '1h16m35s',
        isWeekDay: false,
      },
      {
        date: '2025-04-13',
        weekday: 0,
        workTime1: '08:13:07',
        workTime2: '19:34:35',
        overWorkTimes: '1h51m28s',
        isWeekDay: false,
      },
      {
        date: '2025-04-14',
        weekday: 1,
        workTime1: '08:26:25',
        workTime2: '19:16:45',
        overWorkTimes: '1h20m20s',
        isWeekDay: false,
      },
      {
        date: '2025-04-15',
        weekday: 2,
        workTime1: '08:55:16',
        workTime2: '19:36:48',
        overWorkTimes: '1h11m32s',
        isWeekDay: false,
      },
      {
        date: '2025-04-16',
        weekday: 3,
        workTime1: '08:14:51',
        workTime2: '19:10:33',
        overWorkTimes: '1h25m42s',
        isWeekDay: false,
      },
      {
        date: '2025-04-17',
        weekday: 4,
        workTime1: '08:24:26',
        workTime2: '19:38:47',
        overWorkTimes: '1h44m21s',
        isWeekDay: false,
      },
      {
        date: '2025-04-18',
        weekday: 5,
        workTime1: '08:35:11',
        workTime2: '19:15:24',
        overWorkTimes: '1h10m13s',
        isWeekDay: false,
      },
      {
        date: '2025-04-19',
        weekday: 6,
        workTime1: '08:30:32',
        workTime2: '19:03:36',
        overWorkTimes: '1h3m4s',
        isWeekDay: false,
      },
      {
        date: '2025-04-20',
        weekday: 0,
        workTime1: '08:19:16',
        workTime2: '19:07:01',
        overWorkTimes: '1h17m45s',
        isWeekDay: false,
      },
      {
        date: '2025-04-21',
        weekday: 1,
        workTime1: '08:44:28',
        workTime2: '19:27:41',
        overWorkTimes: '1h13m13s',
        isWeekDay: false,
      },
      {
        date: '2025-04-22',
        weekday: 2,
        workTime1: '08:11:52',
        workTime2: '19:34:08',
        overWorkTimes: '1h52m16s',
        isWeekDay: false,
      },
      {
        date: '2025-04-23',
        weekday: 3,
        workTime1: '08:22:38',
        workTime2: '19:03:23',
        overWorkTimes: '1h10m45s',
        isWeekDay: false,
      },
      {
        date: '2025-04-24',
        weekday: 4,
        workTime1: '08:08:15',
        workTime2: '19:08:15',
        overWorkTimes: '1h30m0s',
        isWeekDay: false,
      },
      {
        date: '2025-04-25',
        weekday: 5,
        workTime1: '08:38:49',
        workTime2: '19:45:27',
        overWorkTimes: '1h36m38s',
        isWeekDay: false,
      },
      {
        date: '2025-04-26',
        weekday: 6,
        workTime1: '08:01:22',
        workTime2: '19:08:23',
        overWorkTimes: '1h37m1s',
        isWeekDay: false,
      },
      {
        date: '2025-04-27',
        weekday: 0,
        workTime1: '08:14:08',
        workTime2: '19:18:42',
        overWorkTimes: '1h34m34s',
        isWeekDay: false,
      },
      {
        date: '2025-04-28',
        weekday: 1,
        workTime1: '08:52:31',
        workTime2: '19:44:23',
        overWorkTimes: '1h21m52s',
        isWeekDay: false,
      },
      {
        date: '2025-04-29',
        weekday: 2,
        workTime1: '08:40:39',
        workTime2: '19:27:54',
        overWorkTimes: '1h17m15s',
        isWeekDay: false,
      },
      {
        date: '2025-04-30',
        weekday: 3,
        workTime1: '08:04:09',
        workTime2: '19:30:32',
        overWorkTimes: '1h56m23s',
        isWeekDay: false,
      },
    ],
  },
  {
    month: '2025-05',
    overtime: '47h3m5s',
    workTime: [
      {
        date: '2025-05-01',
        weekday: 4,
        workTime1: '08:31:15',
        workTime2: '19:59:32',
        overWorkTimes: '1h58m17s',
        isWeekDay: false,
      },
      {
        date: '2025-05-02',
        weekday: 5,
        workTime1: '08:23:18',
        workTime2: '19:51:23',
        overWorkTimes: '1h58m5s',
        isWeekDay: false,
      },
      {
        date: '2025-05-03',
        weekday: 6,
        workTime1: '08:59:43',
        workTime2: '19:06:30',
        overWorkTimes: '36m47s',
        isWeekDay: false,
      },
      {
        date: '2025-05-04',
        weekday: 0,
        workTime1: '08:49:06',
        workTime2: '19:47:38',
        overWorkTimes: '1h28m32s',
        isWeekDay: false,
      },
      {
        date: '2025-05-05',
        weekday: 1,
        workTime1: '08:14:05',
        workTime2: '19:07:13',
        overWorkTimes: '1h23m8s',
        isWeekDay: false,
      },
      {
        date: '2025-05-06',
        weekday: 2,
        workTime1: '08:12:47',
        workTime2: '19:02:59',
        overWorkTimes: '1h20m12s',
        isWeekDay: false,
      },
      {
        date: '2025-05-07',
        weekday: 3,
        workTime1: '08:31:46',
        workTime2: '19:04:43',
        overWorkTimes: '1h2m57s',
        isWeekDay: false,
      },
      {
        date: '2025-05-08',
        weekday: 4,
        workTime1: '08:16:10',
        workTime2: '19:34:37',
        overWorkTimes: '1h48m27s',
        isWeekDay: false,
      },
      {
        date: '2025-05-09',
        weekday: 5,
        workTime1: '08:29:34',
        workTime2: '19:33:41',
        overWorkTimes: '1h34m7s',
        isWeekDay: false,
      },
      {
        date: '2025-05-10',
        weekday: 6,
        workTime1: '08:20:00',
        workTime2: '19:30:24',
        overWorkTimes: '1h40m24s',
        isWeekDay: false,
      },
      {
        date: '2025-05-11',
        weekday: 0,
        workTime1: '08:45:57',
        workTime2: '19:30:43',
        overWorkTimes: '1h14m46s',
        isWeekDay: false,
      },
      {
        date: '2025-05-12',
        weekday: 1,
        workTime1: '08:38:54',
        workTime2: '19:06:53',
        overWorkTimes: '57m59s',
        isWeekDay: false,
      },
      {
        date: '2025-05-13',
        weekday: 2,
        workTime1: '08:20:21',
        workTime2: '19:38:07',
        overWorkTimes: '1h47m46s',
        isWeekDay: false,
      },
      {
        date: '2025-05-14',
        weekday: 3,
        workTime1: '08:39:43',
        workTime2: '19:15:29',
        overWorkTimes: '1h5m46s',
        isWeekDay: false,
      },
      {
        date: '2025-05-15',
        weekday: 4,
        workTime1: '08:28:12',
        workTime2: '19:09:13',
        overWorkTimes: '1h11m1s',
        isWeekDay: false,
      },
      {
        date: '2025-05-16',
        weekday: 5,
        workTime1: '08:34:27',
        workTime2: '19:59:42',
        overWorkTimes: '1h55m15s',
        isWeekDay: false,
      },
      {
        date: '2025-05-17',
        weekday: 6,
        workTime1: '08:14:41',
        workTime2: '19:46:36',
        overWorkTimes: '2h1m55s',
        isWeekDay: false,
      },
      {
        date: '2025-05-18',
        weekday: 0,
        workTime1: '08:09:19',
        workTime2: '19:59:00',
        overWorkTimes: '2h19m41s',
        isWeekDay: false,
      },
      {
        date: '2025-05-19',
        weekday: 1,
        workTime1: '08:48:30',
        workTime2: '19:23:35',
        overWorkTimes: '1h5m5s',
        isWeekDay: false,
      },
      {
        date: '2025-05-20',
        weekday: 2,
        workTime1: '08:37:31',
        workTime2: '19:52:36',
        overWorkTimes: '1h45m5s',
        isWeekDay: false,
      },
      {
        date: '2025-05-21',
        weekday: 3,
        workTime1: '08:13:09',
        workTime2: '19:43:24',
        overWorkTimes: '2h0m15s',
        isWeekDay: false,
      },
      {
        date: '2025-05-22',
        weekday: 4,
        workTime1: '08:39:37',
        workTime2: '19:34:43',
        overWorkTimes: '1h25m6s',
        isWeekDay: false,
      },
      {
        date: '2025-05-23',
        weekday: 5,
        workTime1: '08:25:40',
        workTime2: '19:55:07',
        overWorkTimes: '1h59m27s',
        isWeekDay: false,
      },
      {
        date: '2025-05-24',
        weekday: 6,
        workTime1: '08:54:46',
        workTime2: '19:47:15',
        overWorkTimes: '1h22m29s',
        isWeekDay: false,
      },
      {
        date: '2025-05-25',
        weekday: 0,
        workTime1: '08:06:00',
        workTime2: '19:48:48',
        overWorkTimes: '2h12m48s',
        isWeekDay: false,
      },
      {
        date: '2025-05-26',
        weekday: 1,
        workTime1: '08:58:21',
        workTime2: '19:51:18',
        overWorkTimes: '1h22m57s',
        isWeekDay: false,
      },
      {
        date: '2025-05-27',
        weekday: 2,
        workTime1: '08:48:18',
        workTime2: '19:25:12',
        overWorkTimes: '1h6m54s',
        isWeekDay: false,
      },
      {
        date: '2025-05-28',
        weekday: 3,
        workTime1: '08:37:19',
        workTime2: '19:24:10',
        overWorkTimes: '1h16m51s',
        isWeekDay: false,
      },
      {
        date: '2025-05-29',
        weekday: 4,
        workTime1: '08:04:18',
        workTime2: '19:40:12',
        overWorkTimes: '2h5m54s',
        isWeekDay: false,
      },
      {
        date: '2025-05-30',
        weekday: 5,
        workTime1: '08:55:59',
        workTime2: '19:10:00',
        overWorkTimes: '44m1s',
        isWeekDay: false,
      },
      {
        date: '2025-05-31',
        weekday: 6,
        workTime1: '08:57:54',
        workTime2: '19:39:02',
        overWorkTimes: '1h11m8s',
        isWeekDay: false,
      },
    ],
  },
  {
    month: '2025-06',
    overtime: '46h3m30s',
    workTime: [
      {
        date: '2025-06-01',
        weekday: 0,
        workTime1: '08:28:07',
        workTime2: '19:13:06',
        overWorkTimes: '1h14m59s',
        isWeekDay: false,
      },
      {
        date: '2025-06-02',
        weekday: 1,
        workTime1: '08:15:16',
        workTime2: '19:19:47',
        overWorkTimes: '1h34m31s',
        isWeekDay: false,
      },
      {
        date: '2025-06-03',
        weekday: 2,
        workTime1: '08:02:49',
        workTime2: '19:29:11',
        overWorkTimes: '1h56m22s',
        isWeekDay: false,
      },
      {
        date: '2025-06-04',
        weekday: 3,
        workTime1: '08:47:05',
        workTime2: '19:57:51',
        overWorkTimes: '1h40m46s',
        isWeekDay: false,
      },
      {
        date: '2025-06-05',
        weekday: 4,
        workTime1: '08:15:56',
        workTime2: '19:40:36',
        overWorkTimes: '1h54m40s',
        isWeekDay: false,
      },
      {
        date: '2025-06-06',
        weekday: 5,
        workTime1: '08:41:25',
        workTime2: '19:06:41',
        overWorkTimes: '55m16s',
        isWeekDay: false,
      },
      {
        date: '2025-06-07',
        weekday: 6,
        workTime1: '08:24:00',
        workTime2: '19:34:15',
        overWorkTimes: '1h40m15s',
        isWeekDay: false,
      },
      {
        date: '2025-06-08',
        weekday: 0,
        workTime1: '08:18:14',
        workTime2: '19:00:21',
        overWorkTimes: '1h12m7s',
        isWeekDay: false,
      },
      {
        date: '2025-06-09',
        weekday: 1,
        workTime1: '08:29:33',
        workTime2: '19:28:29',
        overWorkTimes: '1h28m56s',
        isWeekDay: false,
      },
      {
        date: '2025-06-10',
        weekday: 2,
        workTime1: '08:34:17',
        workTime2: '19:51:58',
        overWorkTimes: '1h47m41s',
        isWeekDay: false,
      },
      {
        date: '2025-06-11',
        weekday: 3,
        workTime1: '08:57:30',
        workTime2: '19:23:40',
        overWorkTimes: '56m10s',
        isWeekDay: false,
      },
      {
        date: '2025-06-12',
        weekday: 4,
        workTime1: '08:17:29',
        workTime2: '19:12:15',
        overWorkTimes: '1h24m46s',
        isWeekDay: false,
      },
      {
        date: '2025-06-13',
        weekday: 5,
        workTime1: '08:08:48',
        workTime2: '19:01:21',
        overWorkTimes: '1h22m33s',
        isWeekDay: false,
      },
      {
        date: '2025-06-14',
        weekday: 6,
        workTime1: '08:28:50',
        workTime2: '19:34:03',
        overWorkTimes: '1h35m13s',
        isWeekDay: false,
      },
      {
        date: '2025-06-15',
        weekday: 0,
        workTime1: '08:17:16',
        workTime2: '19:17:21',
        overWorkTimes: '1h30m5s',
        isWeekDay: false,
      },
      {
        date: '2025-06-16',
        weekday: 1,
        workTime1: '08:00:30',
        workTime2: '19:16:44',
        overWorkTimes: '1h46m14s',
        isWeekDay: false,
      },
      {
        date: '2025-06-17',
        weekday: 2,
        workTime1: '08:37:41',
        workTime2: '19:34:34',
        overWorkTimes: '1h26m53s',
        isWeekDay: false,
      },
      {
        date: '2025-06-18',
        weekday: 3,
        workTime1: '08:16:26',
        workTime2: '19:58:19',
        overWorkTimes: '2h11m53s',
        isWeekDay: false,
      },
      {
        date: '2025-06-19',
        weekday: 4,
        workTime1: '08:30:56',
        workTime2: '19:22:28',
        overWorkTimes: '1h21m32s',
        isWeekDay: false,
      },
      {
        date: '2025-06-20',
        weekday: 5,
        workTime1: '08:10:32',
        workTime2: '19:32:15',
        overWorkTimes: '1h51m43s',
        isWeekDay: false,
      },
      {
        date: '2025-06-21',
        weekday: 6,
        workTime1: '08:35:07',
        workTime2: '19:58:17',
        overWorkTimes: '1h53m10s',
        isWeekDay: false,
      },
      {
        date: '2025-06-22',
        weekday: 0,
        workTime1: '08:21:05',
        workTime2: '19:19:49',
        overWorkTimes: '1h28m44s',
        isWeekDay: false,
      },
      {
        date: '2025-06-23',
        weekday: 1,
        workTime1: '08:48:30',
        workTime2: '19:44:01',
        overWorkTimes: '1h25m31s',
        isWeekDay: false,
      },
      {
        date: '2025-06-24',
        weekday: 2,
        workTime1: '08:54:49',
        workTime2: '19:21:30',
        overWorkTimes: '56m41s',
        isWeekDay: false,
      },
      {
        date: '2025-06-25',
        weekday: 3,
        workTime1: '08:28:57',
        workTime2: '19:14:47',
        overWorkTimes: '1h15m50s',
        isWeekDay: false,
      },
      {
        date: '2025-06-26',
        weekday: 4,
        workTime1: '08:38:44',
        workTime2: '19:50:23',
        overWorkTimes: '1h41m39s',
        isWeekDay: false,
      },
      {
        date: '2025-06-27',
        weekday: 5,
        workTime1: '08:58:14',
        workTime2: '19:03:21',
        overWorkTimes: '35m7s',
        isWeekDay: false,
      },
      {
        date: '2025-06-28',
        weekday: 6,
        workTime1: '08:31:08',
        workTime2: '19:49:56',
        overWorkTimes: '1h48m48s',
        isWeekDay: false,
      },
      {
        date: '2025-06-29',
        weekday: 0,
        workTime1: '08:00:32',
        workTime2: '19:31:27',
        overWorkTimes: '2h0m55s',
        isWeekDay: false,
      },
      {
        date: '2025-06-30',
        weekday: 1,
        workTime1: '08:06:54',
        workTime2: '19:41:24',
        overWorkTimes: '2h4m30s',
        isWeekDay: false,
      },
    ],
  },
  {
    month: '2025-07',
    overtime: '46h35m31s',
    workTime: [
      {
        date: '2025-07-01',
        weekday: 2,
        workTime1: '08:29:09',
        workTime2: '19:43:49',
        overWorkTimes: '1h44m40s',
        isWeekDay: false,
      },
      {
        date: '2025-07-02',
        weekday: 3,
        workTime1: '08:21:16',
        workTime2: '19:55:33',
        overWorkTimes: '2h4m17s',
        isWeekDay: false,
      },
      {
        date: '2025-07-03',
        weekday: 4,
        workTime1: '08:27:08',
        workTime2: '19:47:56',
        overWorkTimes: '1h50m48s',
        isWeekDay: false,
      },
      {
        date: '2025-07-04',
        weekday: 5,
        workTime1: '08:23:34',
        workTime2: '19:24:29',
        overWorkTimes: '1h30m55s',
        isWeekDay: false,
      },
      {
        date: '2025-07-05',
        weekday: 6,
        workTime1: '08:50:04',
        workTime2: '19:04:51',
        overWorkTimes: '44m47s',
        isWeekDay: false,
      },
      {
        date: '2025-07-06',
        weekday: 0,
        workTime1: '08:12:33',
        workTime2: '19:01:29',
        overWorkTimes: '1h18m56s',
        isWeekDay: false,
      },
      {
        date: '2025-07-07',
        weekday: 1,
        workTime1: '08:22:55',
        workTime2: '19:23:51',
        overWorkTimes: '1h30m56s',
        isWeekDay: false,
      },
      {
        date: '2025-07-08',
        weekday: 2,
        workTime1: '08:43:36',
        workTime2: '19:43:24',
        overWorkTimes: '1h29m48s',
        isWeekDay: false,
      },
      {
        date: '2025-07-09',
        weekday: 3,
        workTime1: '08:30:59',
        workTime2: '19:20:44',
        overWorkTimes: '1h19m45s',
        isWeekDay: false,
      },
      {
        date: '2025-07-10',
        weekday: 4,
        workTime1: '08:43:05',
        workTime2: '19:12:51',
        overWorkTimes: '59m46s',
        isWeekDay: false,
      },
      {
        date: '2025-07-11',
        weekday: 5,
        workTime1: '08:25:48',
        workTime2: '19:06:11',
        overWorkTimes: '1h10m23s',
        isWeekDay: false,
      },
      {
        date: '2025-07-12',
        weekday: 6,
        workTime1: '08:03:04',
        workTime2: '19:18:37',
        overWorkTimes: '1h45m33s',
        isWeekDay: false,
      },
      {
        date: '2025-07-13',
        weekday: 0,
        workTime1: '08:32:26',
        workTime2: '19:33:02',
        overWorkTimes: '1h30m36s',
        isWeekDay: false,
      },
      {
        date: '2025-07-14',
        weekday: 1,
        workTime1: '08:26:45',
        workTime2: '19:26:58',
        overWorkTimes: '1h30m13s',
        isWeekDay: false,
      },
      {
        date: '2025-07-15',
        weekday: 2,
        workTime1: '08:45:07',
        workTime2: '19:32:00',
        overWorkTimes: '1h16m53s',
        isWeekDay: false,
      },
      {
        date: '2025-07-16',
        weekday: 3,
        workTime1: '08:24:21',
        workTime2: '19:03:28',
        overWorkTimes: '1h9m7s',
        isWeekDay: false,
      },
      {
        date: '2025-07-17',
        weekday: 4,
        workTime1: '08:03:44',
        workTime2: '19:50:57',
        overWorkTimes: '2h17m13s',
        isWeekDay: false,
      },
      {
        date: '2025-07-18',
        weekday: 5,
        workTime1: '08:49:08',
        workTime2: '19:32:10',
        overWorkTimes: '1h13m2s',
        isWeekDay: false,
      },
      {
        date: '2025-07-19',
        weekday: 6,
        workTime1: '08:04:03',
        workTime2: '19:03:22',
        overWorkTimes: '1h29m19s',
        isWeekDay: false,
      },
      {
        date: '2025-07-20',
        weekday: 0,
        workTime1: '08:23:19',
        workTime2: '19:51:30',
        overWorkTimes: '1h58m11s',
        isWeekDay: false,
      },
      {
        date: '2025-07-21',
        weekday: 1,
        workTime1: '08:08:21',
        workTime2: '19:37:37',
        overWorkTimes: '1h59m16s',
        isWeekDay: false,
      },
      {
        date: '2025-07-22',
        weekday: 2,
        workTime1: '08:02:39',
        workTime2: '19:46:34',
        overWorkTimes: '2h13m55s',
        isWeekDay: false,
      },
      {
        date: '2025-07-23',
        weekday: 3,
        workTime1: '08:56:04',
        workTime2: '19:03:22',
        overWorkTimes: '37m18s',
        isWeekDay: false,
      },
      {
        date: '2025-07-24',
        weekday: 4,
        workTime1: '08:10:18',
        workTime2: '19:02:58',
        overWorkTimes: '1h22m40s',
        isWeekDay: false,
      },
      {
        date: '2025-07-25',
        weekday: 5,
        workTime1: '08:20:23',
        workTime2: '19:06:29',
        overWorkTimes: '1h16m6s',
        isWeekDay: false,
      },
      {
        date: '2025-07-26',
        weekday: 6,
        workTime1: '08:46:21',
        workTime2: '19:49:16',
        overWorkTimes: '1h32m55s',
        isWeekDay: false,
      },
      {
        date: '2025-07-27',
        weekday: 0,
        workTime1: '08:28:35',
        workTime2: '19:04:02',
        overWorkTimes: '1h5m27s',
        isWeekDay: false,
      },
      {
        date: '2025-07-28',
        weekday: 1,
        workTime1: '08:19:49',
        workTime2: '19:33:24',
        overWorkTimes: '1h43m35s',
        isWeekDay: false,
      },
      {
        date: '2025-07-29',
        weekday: 2,
        workTime1: '08:04:03',
        workTime2: '19:30:00',
        overWorkTimes: '1h55m57s',
        isWeekDay: false,
      },
      {
        date: '2025-07-30',
        weekday: 3,
        workTime1: '08:55:33',
        workTime2: '19:27:39',
        overWorkTimes: '1h2m6s',
        isWeekDay: false,
      },
      {
        date: '2025-07-31',
        weekday: 4,
        workTime1: '08:34:19',
        workTime2: '19:55:27',
        overWorkTimes: '1h51m8s',
        isWeekDay: false,
      },
    ],
  },
]
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
