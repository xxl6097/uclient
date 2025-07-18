<template>
  <el-dialog
    :modal="true"
    :close-on-click-modal="true"
    :close-on-press-escape="true"
    :width="dialogWidth"
    v-model="showClientDialog"
    :title="title"
  >
    <div class="upgrade-popup-content">
      <el-table :data="paginatedTableData" border>
        <el-table-column
          type="index"
          align="center"
          fixed="left"
          :index="indexMethod"
          min-width="50"
        />
        <el-table-column
          prop="dateTime"
          label="时间"
          min-width="170"
          align="left"
        />
        <el-table-column prop="ago" label="时长" align="left" min-width="90" />
        <el-table-column
          prop="connected"
          label="状态"
          fixed="right"
          align="center"
        >
          <template #default="scope">
            <el-tag v-if="scope.row.connected" type="success">在线</el-tag>
            <el-tag v-else type="danger">离线</el-tag>
          </template>
        </el-table-column>
      </el-table>
      <!-- 分页 -->
      <el-pagination
        style="padding: 20px"
        v-model:page-size="pageSize"
        v-model:current-page="currentPage"
        :pager-count="3"
        :page-sizes="[10, 20, 50, 100, 1000]"
        :layout="responsiveLayout"
        background
        :size="size"
        :total="activities.length"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        :hide-on-single-page="true"
      />
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, defineExpose, computed } from 'vue'
import { Client, TimeLine } from '../utils/type.ts'
import { showErrorTips, showTips } from '../utils/utils.ts'
import { ComponentSize } from 'element-plus'
//layout="sizes, prev, pager, next"
// :layout="responsiveLayout"
const responsiveLayout = ref('sizes, prev, pager, next') // 默认移动端布局
const showClientDialog = ref(false)
const title = ref<string>()

const activities = ref<TimeLine[]>([])
const currentPage = ref<number>(1)
const pageSize = ref<number>(10)
const size = ref<ComponentSize>('default')
const indexMethod = (index: number) => {
  return (
    activities.value.length - (index + (currentPage.value - 1) * pageSize.value)
  )
}
const dialogWidth = ref('80%')

// 分页后的表格数
const paginatedTableData = computed<TimeLine[]>(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return activities.value.slice(start, end)
})

const updateLayout = () => {
  const width = window.innerWidth
  console.log('====>updateLayout', width)
  if (width < 640) {
    // 小屏：简化布局
    responsiveLayout.value = 'prev, pager, next'
    dialogWidth.value = '80%'
  } else if (width < 1024) {
    // 中屏：增加跳转功能
    responsiveLayout.value = 'prev, pager, next'
    dialogWidth.value = '60%'
  } else {
    responsiveLayout.value = 'sizes, prev, pager, next, jumper'
    if (width < 1500) {
      dialogWidth.value = '40%'
      if (width < 1280) {
        dialogWidth.value = '45%'
      }
      if (width < 1140) {
        dialogWidth.value = '48%'
      }
      if (width < 1080) {
        dialogWidth.value = '50%'
      }
    } else {
      dialogWidth.value = '35%'
    }
  }
}
// 分页切换
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
      showClientDialog.value = true
      // activities.value = testTimeLine
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
  updateLayout()
  // showClientDialog.value = true
  // activities.value = row.statusList
  fetchData(row.mac)
}

const updateDialogWidth = () => {
  console.log('打开对话框，updateDialogWidth')
  updateLayout()
}

// 暴露方法供父组件调用
defineExpose({
  openClientDialog: openClientDetailDialog,
  updateDialogWidth: updateDialogWidth,
})
</script>
<style scoped>
.upgrade-popup-header h3 {
  line-height: 2.5;
  margin: 0;
}

.upgrade-popup-footer button {
  margin-left: 10px;
}

.upgrade-popup-content {
  padding-left: 20px;
  padding-right: 20px;
  padding-bottom: 20px;
  display: block;
  place-items: center; /* 水平与垂直居中 */
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
