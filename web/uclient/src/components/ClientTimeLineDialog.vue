<template>
  <el-dialog
    :modal="true"
    :close-on-click-modal="true"
    :close-on-press-escape="true"
    :width="isMobile() ? '80%' : '20%'"
    v-model="showClientDialog"
    :title="title"
  >
    <div class="upgrade-popup-content">
      <!--      <el-timeline style="max-width: 200px">-->
      <!--        <el-timeline-item-->
      <!--          v-for="(activity, index) in activities"-->
      <!--          :key="index"-->
      <!--          :color="activity.connected ? '#55f604' : 'red'"-->
      <!--          :hollow="false"-->
      <!--          :timestamp="formatToUTC8(activity.timestamp)"-->
      <!--        >-->
      <!--          <span :style="{ color: activity.connected ? '#55f604' : 'red' }">-->
      <!--            {{ activity.connected ? '在线' : '离线' }}-{{-->
      <!--              activities.length - index-->
      <!--            }}-->
      <!--          </span>-->
      <!--        </el-timeline-item>-->
      <!--      </el-timeline>-->

      <el-table :data="paginatedTableData" style="width: 90%" border>
        <el-table-column type="index" align="center" :index="indexMethod" />
        <el-table-column prop="timestamp" label="时间" width="200">
          <template #default="props">
            {{ formatTimeStamp(props.row.timestamp) }}
          </template>
        </el-table-column>
        <el-table-column prop="connected" label="状态" align="center">
          <template #default="scope">
            <el-tag v-if="scope.row.connected" type="success">在线</el-tag>
            <el-tag v-else type="danger">离线</el-tag>
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
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, defineExpose, computed } from 'vue'
import { Client, Status } from '../utils/type.ts'
import {
  formatTimeStamp,
  isMobile,
  showErrorTips,
  showTips,
} from '../utils/utils.ts'
import { ComponentSize } from 'element-plus'

const showClientDialog = ref(false)
const title = ref<string>()

const activities = ref<Status[]>([])
const currentPage = ref<number>(1)
const pageSize = ref<number>(10)
const size = ref<ComponentSize>('default')
const indexMethod = (index: number) => {
  return (
    activities.value.length - (index + (currentPage.value - 1) * pageSize.value)
  )
}
// 分页后的表格数
const paginatedTableData = computed<Status[]>(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return activities.value.slice(start, end)
})
// // 分页切换
const handlePageChange = (page: number) => {
  currentPage.value = page
}
const handleSizeChange = (val: number) => {
  console.log(`${val} items per page`)
  pageSize.value = val
}

function fetchData(mac: string) {
  fetch(`../api/get/status?mac=${mac}`, {
    credentials: 'include',
    method: 'GET',
  })
    .then((res) => res.json())
    .then((json) => {
      console.log('get/status', json)
      if (json && json.code === 0 && json.data) {
        console.log(json)
        showClientDialog.value = true
        activities.value = json.data
      } else {
        showTips(json.code, json.msg)
      }
    })
    .catch((error) => {
      console.log('获取失败', error)
      showErrorTips('获取失败')
      // showClientDialog.value = true
      // activities.value = testData
    })
}

function getTitle(row: Client): string {
  if (row.nick && row.nick.name != '') {
    return row.nick.name
  } else {
    return row.hostname
  }
}

const openClientDetailDialog = (row: Client) => {
  console.log('打开对话框，row:', row)
  title.value = `${getTitle(row)}状态时间表`
  // showClientDialog.value = true
  // activities.value = row.statusList
  fetchData(row.mac)
}
//
// const testData = [
//   {
//     timestamp: 1752496014739,
//     connected: true,
//   },
//   {
//     timestamp: 1752496013636,
//     connected: false,
//   },
//   {
//     timestamp: 1752481034289,
//     connected: true,
//   },
//   {
//     timestamp: 1752481032218,
//     connected: false,
//   },
//   {
//     timestamp: 1752481018359,
//     connected: true,
//   },
//   {
//     timestamp: 1752481018010,
//     connected: false,
//   },
//   {
//     timestamp: 1752480934839,
//     connected: true,
//   },
//   {
//     timestamp: 1752480923056,
//     connected: false,
//   },
//   {
//     timestamp: 1752480907559,
//     connected: true,
//   },
//   {
//     timestamp: 1752480905838,
//     connected: false,
//   },
//   {
//     timestamp: 1752479583858,
//     connected: true,
//   },
//   {
//     timestamp: 1752479583421,
//     connected: false,
//   },
//   {
//     timestamp: 1752471812870,
//     connected: true,
//   },
//   {
//     timestamp: 1752471803144,
//     connected: false,
//   },
//   {
//     timestamp: 1752460403098,
//     connected: true,
//   },
//   {
//     timestamp: 1752460329154,
//     connected: false,
//   },
//   {
//     timestamp: 1752329802,
//     connected: true,
//   },
//   {
//     timestamp: 1752303218,
//     connected: false,
//   },
//   {
//     timestamp: 1752270729,
//     connected: true,
//   },
//   {
//     timestamp: 1752265674,
//     connected: true,
//   },
// ]
// 暴露方法供父组件调用
defineExpose({
  openClientDialog: openClientDetailDialog,
})
</script>
<style scoped>
.upgrade-popup-header h3 {
  line-height: 2.5;
  margin: 0;
}

.upgrade-popup-content {
  padding-left: 20px;
  padding-right: 20px;
  padding-bottom: 20px;
  display: grid;
  place-items: center; /* 水平与垂直居中 */
}

.upgrade-popup-footer button {
  margin-left: 10px;
}

.log-container {
  height: auto;
  max-height: 500px;
  overflow-y: auto;
  margin-left: 20px;
}

.log-item {
  margin-bottom: 5px;
}

.autoWidth {
  width: auto;
  min-width: 250px; /* 初始最小宽度 */
  max-width: 400px; /* 初始最小宽度 */
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
</style>
