<template>
  <div class="dialog-wrapper">
    <el-dialog
      :modal="true"
      :close-on-click-modal="true"
      :close-on-press-escape="true"
      :width="isMobile() ? '80%' : 'fit-content'"
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
                    v-model="settings.webHookData.address"
                    placeholder="请输入设备名称"
                  />
                </el-form-item>
              </el-form>
            </el-tab-pane>
            <el-tab-pane label="ntfy设置" name="2">
              <el-form
                label-position="left"
                label-width="auto"
                style="max-width: 600px"
              >
                <el-form-item label="ntfy地址：">
                  <el-input
                    v-model="settings.pushMsgData.address"
                    placeholder="请输入ntfy地址"
                  />
                </el-form-item>
                <el-form-item label="req主题：">
                  <el-input
                    v-model="settings.pushMsgData.reqtopic"
                    placeholder="请输入ntfy订阅req主题"
                  />
                </el-form-item>
                <el-form-item label="res主题：">
                  <el-input
                    v-model="settings.pushMsgData.restopic"
                    placeholder="请输入ntfy订阅res主题"
                  />
                </el-form-item>
                <el-form-item label="用户名：">
                  <el-input
                    v-model="settings.pushMsgData.username"
                    placeholder="请输入ntfy用户名"
                  />
                </el-form-item>
                <el-form-item label="用户密码：">
                  <el-input
                    v-model="settings.pushMsgData.password"
                    placeholder="请输入ntfy用户密码"
                  />
                </el-form-item>
              </el-form>
            </el-tab-pane>

            <el-tab-pane label="微信设置" name="3">
              <el-form
                label-position="left"
                label-width="auto"
                style="max-width: 600px"
              >
                <el-form-item label="OpenID：">
                  <el-input
                    v-model="settings.wechatData.openid"
                    placeholder="请输入微信的openid"
                  />
                </el-form-item>
                <el-form-item label="UserId：">
                  <el-input
                    v-model="settings.wechatData.userid"
                    placeholder="请输入公众号userid"
                  />
                </el-form-item>
                <el-form-item label="template_id：">
                  <el-input
                    v-model="settings.wechatData.template_id"
                    placeholder="请输入公众号template_id"
                  />
                </el-form-item>
              </el-form>
            </el-tab-pane>

            <el-tab-pane label="系统设置" name="4">
              <el-form
                label-position="left"
                label-width="auto"
                style="max-width: 600px"
              >
                <el-form-item label="SysLog监听：">
                  <el-checkbox
                    v-model="settings.listenData.isSysLogListen"
                  ></el-checkbox>
                </el-form-item>
                <el-form-item label="ARP监听：">
                  <el-checkbox
                    v-model="settings.listenData.isArpListen"
                  ></el-checkbox>
                </el-form-item>
                <el-form-item label="Hostapd监听：">
                  <el-checkbox
                    v-model="settings.listenData.isHostApdListen"
                  ></el-checkbox>
                </el-form-item>
                <el-form-item label="Dnsmasq监听：">
                  <el-checkbox
                    v-model="settings.listenData.isDnsmasqListen"
                  ></el-checkbox>
                </el-form-item>
              </el-form>
            </el-tab-pane>
            <el-button type="primary" @click="handleSettingSave"
              >提交
            </el-button>
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
import { SettingsData } from '../utils/type.ts'

const defaultSettings = {
  wechatData: {
    openid: '',
    userid: '',
    template_id: '',
  },
  webHookData: {
    address: '',
  },
  pushMsgData: {
    address: '',
    username: '',
    password: '',
    reqtopic: '',
    restopic: '',
  },
  listenData: {
    isSysLogListen: false,
    isArpListen: false,
    isHostApdListen: false,
    isDnsmasqListen: false,
  },
}
const settings = ref<SettingsData>({ ...defaultSettings })
const formData = ref({
  show: false,
  loading: false,
  title: '',
  activeName: '1',
  // webhookUrl: '',
  // openid: '',
  // ntfy: {
  //   address: '',
  //   topic: '',
  //   username: '',
  //   password: '',
  // },
  // settings: {
  //   isSysLogListen: true,
  //   isArpListen: true,
  //   isHostApdListen: true,
  //   isDnsmasqListen: true,
  // },
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

// const handleOpenIDSetting = () => {
//   if (settings.value.wechatData.openid === '') {
//     showErrorTips('请正确输入openid')
//     return
//   }
//   console.log('handleOpenIDSetting', settings.value.wechatData.openid)
//   const body = {
//     authcode: settings.value.wechatData.openid,
//   }
//   fetch('../api/auth/add', {
//     credentials: 'include',
//     method: 'POST',
//     body: JSON.stringify(body),
//   })
//     .then((res) => {
//       return res.json()
//     })
//     .then((json) => {
//       if (json) {
//         showTips(json.code, json.msg)
//       }
//     })
//     .catch((error) => {
//       showErrorTips(`失败:${JSON.stringify(error)}`)
//     })
// }

const handleSettingSave = () => {
  console.log('handleSetting', settings.value)
  fetch('../api/setting/set', {
    credentials: 'include',
    method: 'POST',
    body: JSON.stringify(settings.value),
  })
    .then((res) => {
      return res.json()
    })
    .then((json) => {
      if (json) {
        showTips(json.code, json.msg)
        formData.value.show = false
      }
    })
    .catch((error) => {
      showErrorTips(`失败:${JSON.stringify(error)}`)
    })
}

const getSettingsData = () => {
  fetch(`../api/setting/get`, {
    credentials: 'include',
    method: 'GET',
  })
    .then((res) => res.json())
    .then((json) => {
      console.log('fetchData', json)
      if (json && json.code === 0 && json.data) {
        console.log('api/setting/get', json)
        const data = json.data
        // ✅ 正确写法：合并对象，不直接覆盖
        settings.value = {
          ...defaultSettings, // 先保留默认结构
          ...data, // 再用后端数据覆盖有值的部分
          listenData: {
            // 嵌套对象也要合并！
            ...defaultSettings.listenData,
            ...data.listenData,
          },
          webHookData: {
            ...defaultSettings.webHookData,
            ...data.webHookData,
          },
          wechatData: {
            ...defaultSettings.wechatData,
            ...data.wechatData,
          },
          pushMsgData: {
            ...defaultSettings.pushMsgData,
            ...data.pushMsgData,
          },
        }
      }
    })
    .catch((error) => {
      console.error(error)
      showErrorTips(`${JSON.stringify(error)}`)
    })
}

// const handleWebhookSetting = () => {
//   if (formData.value.webhookUrl === '') {
//     showErrorTips('请正确输入webhook地址')
//     return
//   }
//   console.log('handleWebhookSetting', formData.value.webhookUrl)
//   const body = {
//     webhookUrl: formData.value.webhookUrl,
//   }
//   fetch('../api/webhook/set', {
//     credentials: 'include',
//     method: 'POST',
//     body: JSON.stringify(body),
//   })
//     .then((res) => {
//       return res.json()
//     })
//     .then((json) => {
//       if (json) {
//         showTips(json.code, json.msg)
//       }
//     })
//     .catch((error) => {
//       showErrorTips(`失败:${JSON.stringify(error)}`)
//     })
// }

// const handleNtfySetting = () => {
//   console.log('handleNtfySetting', formData.value.ntfy)
//   fetch('../api/ntfy/set', {
//     credentials: 'include',
//     method: 'POST',
//     body: JSON.stringify(formData.value.ntfy),
//   })
//     .then((res) => {
//       return res.json()
//     })
//     .then((json) => {
//       if (json) {
//         showTips(json.code, json.msg)
//       }
//     })
//     .catch((error) => {
//       showErrorTips(`失败:${JSON.stringify(error)}`)
//     })
// }

const showDialogForm = () => {
  console.log('打开对话框，row:')
  getSettingsData()
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
