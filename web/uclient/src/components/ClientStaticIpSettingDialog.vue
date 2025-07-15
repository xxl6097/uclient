<template>
  <el-dialog
    :modal="true"
    :close-on-click-modal="true"
    :close-on-press-escape="true"
    :width="isMobile() ? '80%' : '30%'"
    v-model="formData.show"
    :title="formData.title"
  >
    <div class="upgrade-popup-content">
      <el-form label-width="130px">
        <el-form-item label="设备名称：">
          <el-input
            v-model="formData.client.hostname"
            placeholder="请输入设备名称"
          />
        </el-form-item>
        <el-form-item label="设备Mac：">
          <el-input
            v-model="formData.client.mac"
            placeholder="请输入设备Mac地址"
          />
        </el-form-item>
        <el-form-item label="设备IP：">
          <el-input v-model="formData.client.ip" placeholder="请输入设备IP" />
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
      <el-button
        type="danger"
        :loading="formData.loading"
        :loading-icon="Eleme"
        @click="handleDeleteStaticIp"
        >{{ formData.loading ? '静态地址删除中...' : '删除' }}
      </el-button>
      <el-button
        type="primary"
        :loading="formData.loading"
        :loading-icon="Eleme"
        @click="handleConfirm"
        >{{ formData.loading ? '静态地址设置中...' : '确定' }}
      </el-button>
    </template>
  </el-dialog>
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
import { ElMessageBox } from 'element-plus'

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

const handleDeleteStaticIp = (row: Client) => {
  console.log('handleDeleteStaticIp', row)
  ElMessageBox.confirm(`确定删除【${row.hostname}】静态IP吗?`, 'Warning', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
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
        })
    })
    .catch(() => {})
}

function handleConfirm() {
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
        formData.value.show = false
      }, 500)
    })
}

const showDialogForm = (row: Client) => {
  console.log('打开对话框，row:', row)
  formData.value.title = `设置静态IP`
  formData.value.client = JSON.parse(JSON.stringify(row))
  formData.value.show = true
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
