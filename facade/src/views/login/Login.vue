<template>
  <div class="loginPage">

    <div class="loginTitle">
     <img class="loginImg" :src="src" />
      <p>Insight</p>
    </div>

    <div class="loginBox">
      <!--      <a-tabs v-model:activeKey="activeKey">-->
      <!--        <a-tab-pane key="1" tab="账户密码登录">-->
      <!--      <p>账户密码登录</p>-->
      <a-form :model="formState" name="normal_login" class="login-form" @finish="onFinish"
              @finishFailed="onFinishFailed">

        <a-form-item label="账号" name="username" :rules="[{ required: true, message: '请输入用户名!' }]">
          <a-input v-model:value="formState.username" placeholder="请输入用户名">
            <template #prefix>
              <UserOutlined class="site-form-item-icon"/>
            </template>
          </a-input>
        </a-form-item>

        <a-form-item label="密码" name="password" :rules="[{ required: true, message: '请输入密码!' }]">
          <a-input-password v-model:value="formState.password" placeholder="请输入密码">
            <template #prefix>
              <LockOutlined class="site-form-item-icon"/>
            </template>
          </a-input-password>
        </a-form-item>
        <div class="lbTips">
          <a-checkbox v-model:checked="formState.remember">记住密码</a-checkbox>
          <a class="login-form-forgot" href="">忘记密码</a>
        </div>

        <a-form-item class="lbBtn">
          <a-button block="true" type="primary" html-type="submit" class="login-form-button">
            登录
          </a-button>
        </a-form-item>
      </a-form>
    </div>

  </div>
</template>

<script setup lang="ts">
import type {CSSProperties} from 'vue';
import {reactive, computed, h} from 'vue';
import {UserOutlined, LockOutlined, ExclamationCircleOutlined} from '@ant-design/icons-vue';
import {ref} from 'vue';

import request from "@/utils/request"
import {Modal} from 'ant-design-vue';
import {useRouter} from "vue-router";

const router = useRouter();
const src = import.meta.env.VITE_ICON_PATH

interface FormState {
  username: string;
  password: string;
  remember: boolean;
}

const activeKey = ref('1');
const formState = reactive<FormState>({
  username: '',
  password: '',
  remember: true,
});
const onFinish = (values: any) => {
  let params = {
    "username": values.username,
    "passwd": values.password
  }
  request({
    url: "/login",
    method: "POST",
    params
  }).then((resp) => {
    if (401 === resp['status']) {
      Modal.confirm({
        title: '登录信息错误',
        icon: h(ExclamationCircleOutlined, {}),
        content: "用户名或密码错误，请假查后重新登录。",
      });
    } else if (200 === resp['status']) {
      let data = resp['data'], token = data['token'], user = data['user']
      localStorage.setItem("token", token)
      localStorage.setItem("curr_user", JSON.stringify(user))

      router.push("/")

    }
  })
};

const onFinishFailed = (errorInfo: any) => {
  console.log('Failed:', errorInfo);
};
const disabled = computed(() => {
  return !(formState.username && formState.password);
});
const cardStyle: CSSProperties = {
  width: '620px',
};
const imgStyle: CSSProperties = {
  display: 'block',
  width: '273px',
};
</script>

<style scoped>
/*#components-form-demo-normal-login .login-form {*/
/*  max-width: 200px;*/
/*}*/
#components-form-demo-normal-login .login-form-forgot {
  float: right;
}

#components-form-demo-normal-login .login-form-button {
  width: 100%;
}

.loginPage {
  height: 100%;
  width: 100%;
  /*background-color:grey;*/
  background-image: url("@/assets/images/background.svg");
  overflow: hidden;
}

.loginTitle {
  margin-top: 200px;
  text-align: center;
  font-size: 40px;
  letter-spacing: 2px;
  font-weight: bold;
  display: flex;
  justify-content: center;
  align-items: center
}

.loginImg {
  width: 100px;
  margin-right: 20px;
  object-fit: cover;
}

.loginTitle p{
  margin-bottom: 0
}

.loginBox {
  width: 350px;
  border: none;
  margin: 50px auto 0;
}

.lbTips {
  margin: 20px 0 30px;
  display: flex;
  justify-content: space-between;
}
</style>
<style>
.loginPage .ant-input {
  height: 36px;
}

.loginPage .ant-btn {
  height: 40px;
}
</style>
