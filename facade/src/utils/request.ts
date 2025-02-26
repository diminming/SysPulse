import axios from 'axios'
import { notification, Modal } from 'ant-design-vue';
import { AlertOutlined, ExclamationCircleOutlined, StopOutlined } from "@ant-design/icons-vue";
import { h } from 'vue';
import { lock403Store } from "@/stores/lock403"

const BASE_URL = import.meta.env.VITE_BASE_URL

// create an axios instance
const service = axios.create({
  baseURL: BASE_URL, // url = base url + request url
  // withCredentials: true, // send cookies when cross-domain requests
  timeout: import.meta.env.VITE_REQUEST_TIMEOUT * 1000 // request timeout, unit is seconds.
}),
  lock403 = lock403Store()

// request interceptor
service.interceptors.request.use(
  config => {
    let token = localStorage.getItem("token")
    if(token) {
      config.headers["token"] = token
    }
    return config
  },
  error => {
    // do something with request error
    console.log(error) // for debug
    return Promise.reject(error)
  }
)

// response interceptor
service.interceptors.response.use(
  /**
   * If you want to get http information such as headers or status
   * Please return  response => response
  */

  /**
   * Determine the request status by custom code
   * Here is just an example
   * You can also judge the status by HTTP Status Code
   */
  response => {
    // const res = response.data
    // return res
    return Promise.resolve(response['data'])
  },
  error => {
    const status = error.response.status
    const data = error.response.data
    if (status === 403) {
      
      if(lock403.isLocked()) return

      lock403.lock()
      Modal.warning({
        title: '登录信息错误',
        icon: h(ExclamationCircleOutlined, {}),
        content: "您尚未登录或登陆信息已过期，请您重新登陆。",
        onOk() {
          setTimeout(()=>{
            location.href = import.meta.env.VITE_LOGIN_PAGE
          }, 500)
        },
        class: 'test',
      });
      console.error("request without login info.")
      // location.href = LOGIN_LOCATION
    } else if(status === 401) {
      Modal.warning({
        title: '权限不足',
        icon: h(StopOutlined, {}),
        content: `您没有权限访问此数据或功能，如有疑问请联系管理员。`,
        class: 'test',
      });
      console.error(error)
    }else {
      notification.error({
        message: `服务端错误: ${data["status"]}`,
        description: data["msg"],
        icon: () => h(AlertOutlined, { style: "color: #ff4d4f" }),
        onClick: () => {
          notification.destroy()
        }
      })
    }
    return Promise.reject(error)
  }
)

export default service