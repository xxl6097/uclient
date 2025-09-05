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
            <el-tab-pane label="webhook设置" name="1">
              <el-form>
                <el-form-item label="webhook地址：">
                  <el-input
                    v-model="formData.webhookUrl"
                    placeholder="请输入设备名称"
                  />
                </el-form-item>
              </el-form>
              <el-button type="primary" @click="handleWebhookSetting"
                >提交
              </el-button>
            </el-tab-pane>
            <el-tab-pane label="ntfy设置" name="2">
              <el-form
                label-position="left"
                label-width="auto"
                style="max-width: 600px"
              >
                <el-form-item label="ntfy地址：">
                  <el-input
                    v-model="formData.ntfy.address"
                    placeholder="请输入ntfy地址"
                  />
                </el-form-item>
                <el-form-item label="ntfy订阅主题：">
                  <el-input
                    v-model="formData.ntfy.topic"
                    placeholder="请输入ntfy订阅主题"
                  />
                </el-form-item>
                <el-form-item label="ntfy用户名称：">
                  <el-input
                    v-model="formData.ntfy.username"
                    placeholder="请输入ntfy用户名"
                  />
                </el-form-item>
                <el-form-item label="ntfy用户密码：">
                  <el-input
                    v-model="formData.ntfy.password"
                    placeholder="请输入ntfy用户密码"
                  />
                </el-form-item>
              </el-form>
              <el-button type="danger" @click="handleNtfySetting"
                >提交</el-button
              >
            </el-tab-pane>

            <el-tab-pane label="微信设置" name="1">
              <el-form>
                <el-form-item label="OpenID：">
                  <el-input
                    v-model="formData.openid"
                    placeholder="请输入微信的openid"
                  />
                </el-form-item>
              </el-form>
              <el-button type="primary" @click="handleOpenIDSetting"
              >提交
              </el-button>
            </el-tab-pane>
          </el-tabs>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, defineExpose } from 'vue'
import { isMobile, showErrorTips, showTips } from '../utils/utils.ts'
import { TabsPaneContext } from 'element-plus'

const formData = ref({
  show: false,
  loading: false,
  activeName: '1',
  webhookUrl: '',
  openid: '',
  ntfy: {
    address: '',
    topic: '',
    username: '',
    password: '',
  },
  title: '',
})

const handleClick = (tab: TabsPaneContext) => {
  console.log('handleClick', tab.paneName)
  switch (tab.paneName) {
    case 'first':
      break
    case 'second':
      break
    case 'thrid':
      break
  }
}

const handleOpenIDSetting = () => {
  if (formData.value.openid === '') {
    showErrorTips('请正确输入openid')
    return
  }
  console.log('handleOpenIDSetting', formData.value.openid)
  const body = {
    authcode: formData.value.openid,
  }
  fetch('../api/auth/add', {
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
}

const handleWebhookSetting = () => {
  if (formData.value.webhookUrl === '') {
    showErrorTips('请正确输入webhook地址')
    return
  }
  console.log('handleWebhookSetting', formData.value.webhookUrl)
  const body = {
    webhookUrl: formData.value.webhookUrl,
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
}

const handleNtfySetting = () => {
  console.log('handleNtfySetting', formData.value.ntfy)
  fetch('../api/ntfy/set', {
    credentials: 'include',
    method: 'POST',
    body: JSON.stringify(formData.value.ntfy),
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
}

const showDialogForm = () => {
  console.log('打开对话框，row:')
  formData.value.title = `设备设置`
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
