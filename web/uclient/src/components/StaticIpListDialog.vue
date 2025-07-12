<template>
  <el-dialog
    :modal="true"
    :close-on-click-modal="true"
    :close-on-press-escape="true"
    :width="isMobile() ? '80%' : '80%'"
    v-model="formData.show"
    :title="formData.title"
  >
    <div class="upgrade-popup-content">
      <section>
        <el-main>
          <el-table
            :data="paginatedTableData"
            style="width: 100%"
            :border="true"
            :preserve-expanded-content="true"
          >
            <el-table-column prop="index" label="索引" sortable />
            <el-table-column prop="hostname" label="名称" sortable />
            <el-table-column prop="ip" label="IP" sortable />
            <el-table-column
              prop="mac"
              label="Mac地址"
              sortable
              v-if="!isMobile()"
            />
            <el-table-column label="操作">
              <template #default="{ row }">
                <el-button size="small" type="text" @click="handleDelete(row)"
                  >删除
                </el-button>
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
            :pager-count="7"
            @current-change="handlePageChange"
          />
        </el-main>
      </section>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, defineExpose } from 'vue'
import { Client, DHCPHost } from '../utils/type.ts'
import {
  isMobile,
  showErrorTips,
  showLoading,
  showTips,
} from '../utils/utils.ts'

// 搜索关键字
const searchKeyword = ref<string>('')
const pageSize = ref<number>(50)
const currentPage = ref<number>(1)
const tableData = ref<DHCPHost[]>([])
// 分页后的表格数
const paginatedTableData = computed<DHCPHost[]>(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredTableData.value.slice(start, end)
})
// 过滤后的表格数据（根据搜索关键字）
const filteredTableData = computed<DHCPHost[]>(() => {
  return tableData.value.filter(() => !searchKeyword.value)
})
// 分页切换
const handlePageChange = (page: number) => {
  currentPage.value = page
}

function renderTable(data: any) {
  tableData.value = data as DHCPHost[]
}

const formData = ref({
  show: false,
  loading: false,
  title: '',
  client: {
    hostname: '',
    ip: '',
    mac: '',
  } as Client,
})

function refreshList() {
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
        if (json.data) {
          renderTable(json.data)
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

function handleDelete(row: any) {
  console.log(row)
  const loader = showLoading('删除中...')
  fetch(`../api/staticip/delete?mac=${row.mac}`, {
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
      refreshList()
    })
}

const showDialogForm = (list: DHCPHost[]) => {
  console.log('打开对话框，row:')
  formData.value.title = `静态IP列表`
  formData.value.show = true
  renderTable(list)
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
  padding-left: 20px;
  padding-right: 20px;
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
</style>
