<template>
  <el-dialog
    :modal="true"
    :close-on-click-modal="true"
    :close-on-press-escape="true"
    :width="isMobile() ? '80%' : '30%'"
    v-model="showClientDialog"
    :title="title"
  >
    <div class="upgrade-popup-content">
      <el-timeline style="max-width: 200px">
        <el-timeline-item
          v-for="(activity, index) in activities"
          :key="index"
          :color="activity.connected ? '#55f604' : 'red'"
          :hollow="false"
          :timestamp="formatToUTC8(activity.timestamp)"
        >
          <span :style="{ color: activity.connected ? '#55f604' : 'red' }">
            {{ activity.connected ? '在线' : '离线' }}-{{ activities.length - index }}
          </span>
        </el-timeline-item>
      </el-timeline>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, defineExpose } from 'vue'
import { Client, Status } from '../utils/type.ts'
import {
  formatToUTC8,
  isMobile,
  showErrorTips,
  showTips,
} from '../utils/utils.ts'

const showClientDialog = ref(false)
const title = ref<string>()

const activities = ref<Status[]>([])

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
    .catch(() => {
      showErrorTips('获取失败')
    })
}

function getTitle(row: Client): string {
  if (row.nickName != '') {
    return row.nickName
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
