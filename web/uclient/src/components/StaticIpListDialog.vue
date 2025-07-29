<template>
  <el-dialog
    :modal="true"
    :close-on-click-modal="true"
    :close-on-press-escape="true"
    :width="isMobile() ? '80%' : '60%'"
    v-model="formData.show"
    :title="formData.title"
  >
    <div class="upgrade-popup-content">
      <section>
        <el-main>
          <div style="display: flex; margin-bottom: 20px">
            <div style="display: flex">
              <span style="min-width: 70px; align-content: center"
                >主机名：</span
              >
              <el-input
                placeholder="请输入主机名"
                v-model="valueName"
              ></el-input>
            </div>

            <div style="display: flex; margin-left: 20px">
              <span style="min-width: 70px; align-content: center"
                >Mac地址：</span
              >
              <el-select
                v-model="valueMac"
                placeholder="请选择Mac地址"
                style="width: 240px"
              >
                <el-option
                  v-for="item in citiesMac"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
                <template #footer>
                  <el-button
                    v-if="!isAddingMac"
                    text
                    bg
                    size="small"
                    @click="onAddOptionMac"
                  >
                    Add an option
                  </el-button>
                  <template v-else>
                    <el-input
                      v-model="optionNameMac"
                      class="option-input"
                      placeholder="input option name"
                      size="small"
                    />
                    <el-button
                      type="primary"
                      size="small"
                      @click="onConfirmMac"
                    >
                      confirm
                    </el-button>
                    <el-button size="small" @click="clearMac">cancel</el-button>
                  </template>
                </template>
              </el-select>
            </div>

            <div style="display: flex; margin-left: 20px">
              <span style="min-width: 70px; align-content: center"
                >IP地址：</span
              >
              <el-select
                v-model="valueIp"
                placeholder="请选择IP地址"
                style="width: 240px"
              >
                <el-option
                  v-for="item in citiesIp"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                  :disabled="item.disabled"
                />
                <template #footer>
                  <el-button
                    v-if="!isAddingIp"
                    text
                    bg
                    size="small"
                    @click="onAddOptionIp"
                  >
                    Add an option
                  </el-button>
                  <template v-else>
                    <el-input
                      v-model="optionNameIp"
                      class="option-input"
                      placeholder="input option name"
                      size="small"
                    />
                    <el-button type="primary" size="small" @click="onConfirmIp">
                      confirm
                    </el-button>
                    <el-button size="small" @click="clearIp">cancel</el-button>
                  </template>
                </template>
              </el-select>
            </div>

            <el-button style="margin-left: 20px" @click="handleAdd"
              >新增
            </el-button>
          </div>
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
  showSucessTips,
  showTips,
} from '../utils/utils.ts'
import { CheckboxValueType } from 'element-plus'

interface Option {
  label: string
  value: string
  disabled: boolean
}

const valueName = ref<string>('')
const isAddingMac = ref(false)
const valueMac = ref<CheckboxValueType[]>([])
const optionNameMac = ref('')
const citiesMac = ref<Option[]>([])
const onAddOptionMac = () => {
  isAddingMac.value = true
}
const onConfirmMac = () => {
  if (optionNameMac.value) {
    citiesMac.value.push({
      label: optionNameMac.value,
      value: optionNameMac.value,
      disabled: false,
    })
    clearMac()
  }
}

const clearMac = () => {
  optionNameMac.value = ''
  isAddingMac.value = false
}

const isAddingIp = ref(false)
const valueIp = ref<CheckboxValueType[]>([])
const optionNameIp = ref('')
const citiesIp = ref<Option[]>([])

const onAddOptionIp = () => {
  isAddingIp.value = true
}

const onConfirmIp = () => {
  if (optionNameIp.value) {
    citiesIp.value.push({
      label: optionNameIp.value,
      value: optionNameIp.value,
      disabled: false,
    })
    clearIp()
  }
}

const clearIp = () => {
  optionNameIp.value = ''
  isAddingIp.value = false
}

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

function renderTable(data: DHCPHost[]) {
  tableData.value = data
}

function renderClients(clients: Client[], data: DHCPHost[]) {
  valueName.value = ''
  valueIp.value.splice(0, valueName.value.length)
  valueMac.value.splice(0, valueMac.value.length)
  citiesMac.value.splice(0, citiesMac.value.length)
  citiesIp.value.splice(0, citiesIp.value.length)
  isAddingMac.value = false
  isAddingIp.value = false
  optionNameIp.value = ''
  optionNameMac.value = ''
  if (data) {
    data.forEach((item) => {
      citiesIp.value.push({
        label: `${item.ip}(${item.hostname})`,
        value: item.ip,
        disabled: true,
      })
    })
  }

  if (clients) {
    clients.forEach((client) => {
      citiesMac.value.push({
        label: `${client.mac}(${client.hostname})`,
        value: client.mac,
        disabled: false,
      })

      let has = false
      citiesIp.value.forEach((ip) => {
        if (ip.value === client.ip) {
          has = true
        }
      })
      if (!has) {
        citiesIp.value.push({
          label: `${client.ip}(${client.hostname})`,
          value: client.ip,
          disabled: false,
        })
      }
      // citiesIp.value.push({
      //   label: `${client.ip}(${client.hostname})`,
      //   value: client.ip,
      //   disabled: false,
      // })
    })
  }
}

function handleAdd() {
  console.log('handleAdd', valueName.value, valueMac.value, valueIp.value)
  const body = {
    hostname: valueName.value,
    ip: valueIp.value,
    mac: valueMac.value,
  }
  const loading = showLoading('静态IP设置中...')
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
      loading.close()
      refreshList()
    })
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

const showDialogForm = (list: DHCPHost[], clients: Client[]) => {
  console.log('打开对话框，row:')
  formData.value.title = `静态IP列表`
  formData.value.show = true
  renderTable(list)
  renderClients(clients, list)
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
