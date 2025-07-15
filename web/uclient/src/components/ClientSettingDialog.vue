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
          </el-tabs>
        </div>
      </template>
      <template #footer>
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
import { ref, defineExpose } from 'vue'
import { Client } from '../utils/type.ts'
import {
  isMobile,
  showErrorTips,
  showLoading,
  showSucessTips,
  showTips,
} from '../utils/utils.ts'
import { ElMessageBox, TabsPaneContext } from 'element-plus'

const formData = ref({
  show: false,
  loading: false,
  activeName: 'first',
  hideErrBtn: false,
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
  } as Client,
})

const handleClick = (tab: TabsPaneContext) => {
  formData.value.hideErrBtn = false
  switch (tab.paneName) {
    case 'first':
      break
    case 'second':
      formData.value.hideErrBtn = true
      break
  }
}

const handleChangeNickName = () => {
  const row = {
    isPush: formData.value.first.isPush,
    name: formData.value.first.name,
    starTime: formData.value.client.starTime,
    mac: formData.value.client.mac,
    ip: formData.value.client.ip,
    hostname: formData.value.client.hostname,
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
