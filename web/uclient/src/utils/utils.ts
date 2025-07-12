import { ElLoading, ElMessage, ElMessageBox } from 'element-plus'

export function showWarmDialog(title: string, ok: any, cancel: any) {
  ElMessageBox.confirm(title, 'Warning', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      ok()
    })
    .catch(() => {
      cancel()
    })
}

export function showErrorTips(message: string) {
  ElMessage({
    showClose: true,
    message: message,
    type: 'error',
  })
}

export function showTips(code: any, message: string) {
  if (code === 0) {
    showSucessTips(message)
  } else {
    showWarmTips(message)
  }
}

export function showSucessTips(message: string) {
  ElMessage({
    showClose: true,
    message: message,
    type: 'success',
  })
}

export function showWarmTips(message: string) {
  ElMessage({
    showClose: true,
    message: message,
    type: 'warning',
  })
}

/**
 * 基于 Promise 封装的 XMLHttpRequest 请求
 * @param {Object} config - 请求配置
 * @param {string} config.url - 请求地址
 * @param {string} [config.method='GET'] - 请求方法
 * @param {Object} [config.headers] - 请求头
 * @param {any} [config.data] - 请求数据
 * @param {number} [config.timeout=0] - 超时时间（毫秒）
 * @param {string} [config.responseType] - 响应类型
 * @param {Function} [config.onUploadProgress] - 上传进度回调
 * @param {Function} [config.onDownloadProgress] - 下载进度回调
 * @returns {Promise} 返回 Promise 对象
 */
export function xhrPromise(config: any) {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()
    // 初始化请求
    xhr.open(config.method || 'GET', config.url)
    // 设置请求头
    if (config.headers) {
      Object.entries(config.headers).forEach(([key, value]) => {
        xhr.setRequestHeader(key, value as string)
      })
    }

    // 设置响应类型
    if (config.responseType) {
      xhr.responseType = config.responseType
    }

    // 设置超时
    if (config.timeout) {
      xhr.timeout = config.timeout
    }

    // 上传进度处理
    if (config.onUploadProgress) {
      xhr.upload.onprogress = (event) => {
        if (event.lengthComputable) {
          const percentComplete = (event.loaded / event.total) * 100
          console.log('--->', percentComplete + '%')
          config.onUploadProgress(percentComplete.toFixed(2))
        }
      }
    }

    // 下载进度处理
    if (config.onDownloadProgress) {
      xhr.onprogress = (e) => {
        config.onDownloadProgress({
          loaded: e.loaded,
          total: e.total,
          progress: e.loaded / e.total,
        })
      }
    }

    // 请求成功处理
    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        resolve({
          data: xhr.response,
          status: xhr.status,
          statusText: xhr.statusText,
          headers: xhr.getAllResponseHeaders(),
        })
      } else {
        reject(new Error(`请求失败：${xhr.status} ${xhr.statusText}`))
      }
    }

    // 错误处理
    xhr.onerror = () => reject(new Error('网络错误'))
    xhr.ontimeout = () => reject(new Error(`请求超时（${config.timeout}ms）`))
    xhr.onabort = () => reject(new Error('请求被中止'))

    // 发送请求
    try {
      xhr.send(config.data)
    } catch (err) {
      reject(err)
    }
  })
}

export function syntaxHighlight(json: string): string {
  // 转义特殊字符防止 XSS
  json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

  // 正则匹配 JSON 元素并分配类名
  return json.replace(
    /("(\\u[\dA-Fa-f]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+-]?\d+)?)/g,
    (match) => {
      let cls = 'number'
      if (/^"/.test(match)) {
        cls = match.endsWith(':') ? 'key' : 'string' // 键名与字符串区分
      } else if (/true|false/.test(match)) {
        cls = 'boolean'
      } else if (/null/.test(match)) {
        cls = 'null'
      }
      return `<span class="${cls}">${match}</span>` // 直接内联类名判断[1,6](@ref)
    },
  )
}

export function isMobile() {
  if (isMobileWidth() || isMobilePhone()) {
    return true
  }
  return false
}

export function isMobileWidth() {
  if (window.innerWidth >= 992) {
    return false
  }
  return true
}

export function isMobilePhone() {
  const ua = navigator.userAgent.toLowerCase()
  return !!ua.match(/iOS|iPhone|Android|windows Phone|BB\d+/i)
}

export function showLoading(title: string) {
  return ElLoading.service({
    lock: true,
    text: title,
    background: 'rgba(0, 0, 0, 0.7)',
  })
}

export function markdownToHtml(markdown: string): string {
  let lines: string[] = markdown.split('\n')
  let html: string = ''
  let inList: boolean = false
  let listItems: string[] = []
  let inCodeBlock: boolean = false
  let codeBlockContent: string = ''

  for (let i = 0; i < lines.length; i++) {
    let line: string = lines[i].trim()

    // 处理代码块开始
    if (line.startsWith('```')) {
      if (inCodeBlock) {
        html += `<pre><code>${codeBlockContent}</code></pre>`
        inCodeBlock = false
        codeBlockContent = ''
      } else {
        inCodeBlock = true
      }
      continue
    }

    if (inCodeBlock) {
      codeBlockContent += line + '\n'
      continue
    }

    // 处理标题
    if (/^(#+) (.*)$/.test(line)) {
      let [, hashes, content] = line.match(/^(#+) (.*)$/)!
      let level: number = hashes.length
      if (inList) {
        html += `<ul>${listItems.join('')}</ul>`
        inList = false
        listItems = []
      }
      html += `<h${level}>${content}</h${level}>`
    }
    // 处理无序列表
    else if (/^([*-]) (.*)$/.test(line)) {
      let [, , content] = line.match(/^([*-]) (.*)$/)!
      if (!inList) {
        inList = true
      }
      listItems.push(`<li>${content}</li>`)
    }
    // 处理段落
    else {
      if (inList) {
        html += `<ul>${listItems.join('')}</ul>`
        inList = false
        listItems = []
      }
      if (line) {
        html += `<p>${line}</p>`
      }
    }
  }

  // 如果最后处于列表状态，闭合列表
  if (inList) {
    html += `<ul>${listItems.join('')}</ul>`
  }

  // 如果最后处于代码块状态，闭合代码块
  if (inCodeBlock) {
    html += `<pre><code>${codeBlockContent}</code></pre>`
  }

  // 处理加粗
  html = html.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')

  // 处理斜体
  html = html.replace(/\*(.*?)\*/g, '<em>$1</em>')

  return html
}

// export function formatTimeStamp(timestamp: number): string {
//   const date = new Date(timestamp * 1000) // 秒转毫秒[3,5](@ref)
//   const year = date.getFullYear()
//   const month = String(date.getMonth() + 1).padStart(2, '0') // 月份从0开始[1,7](@ref)
//   const day = String(date.getDate()).padStart(2, '0')
//   const hours = String(date.getHours()).padStart(2, '0')
//   const minutes = String(date.getMinutes()).padStart(2, '0')
//   const seconds = String(date.getSeconds()).padStart(2, '0')
//   return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
// }

// /**
//  * 将Unix时间戳格式化为东八区时间字符串
//  * @param {number} timestamp - Unix时间戳（秒级）
//  * @returns {string} 格式为 "YYYY-MM-DD HH:mm:ss" 的东八区时间字符串
//  */
export function formatToUTC8001(timestamp: number): string {
  // 将秒级时间戳转换为毫秒
  const date = new Date(timestamp * 1000)

  // 调整到东八区时间
  date.setHours(date.getHours() -8)

  // 提取时间分量并补零
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')

  // 组合成目标格式
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

/**
 * 将Unix时间戳格式化为东八区时间字符串（修正版本）
 * @param {number} timestamp - Unix时间戳（秒级）
 * @returns {string} 格式为 "YYYY-MM-DD HH:mm:ss" 的东八区时间字符串
 */
export function formatToUTC80002(timestamp: number): string {
  // 创建UTC时间对象（时间戳本质是UTC时间）
  const utcDate = new Date(timestamp * 1000)

  // 直接转换为东八区时间（无需额外加8小时）
  const date = new Date(utcDate.getTime())

  // 使用ISO日期方法直接获取本地时间分量
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

function formatTime(timestamp: number): string {
  // 判断时间戳类型（秒级或毫秒级）
  const date = timestamp.toString().length > 10
      ? new Date(timestamp)
      : new Date(timestamp * 1000); // 秒级需转为毫秒

  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}

export function formatTimeStamp(timestamp: number): string {
  return formatTime(timestamp)
}

export function formatToUTC8(timestamp: number): string {
  return formatTime(timestamp)
}