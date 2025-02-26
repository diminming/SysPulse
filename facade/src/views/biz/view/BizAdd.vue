<template>
  <a-layout-content :style="{
    background: '#fff',
    padding: '24px',
    margin: 0,
    minHeight: '800px',
  }">
    <div class="systemPage">
      <a-form ref="formRef" :model="bizObj" :rules="rules" :label-col="labelCol" :wrapper-col="wrapperCol">
        <a-form-item label="业务名称">
          <a-input v-model:value="bizObj.bizName" placeholder="请输入业务名称" />
        </a-form-item>
        <a-form-item label="业务标识">
          <a-input v-model:value="bizObj.bizId" placeholder="请输入业务标识" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model:value="bizObj.bizDesc" placeholder="请输入备注" />
        </a-form-item>
        <a-form-item :wrapper-col="{ span: 6, offset: 2 }">
          <a-button type="primary" @click="saveBizInfo">确定</a-button>
          <a-button style="margin-left: 10px" @click="resetForm">返回</a-button>
        </a-form-item>
      </a-form>
    </div>
  </a-layout-content>

</template>

<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue';
import { notification } from 'ant-design-vue';
import type { Rule } from 'ant-design-vue/es/form';
import { useRoute, useRouter } from "vue-router";
import request from "@/utils/request";

const bizObj = reactive({
  id: 0,
  bizName: "",
  bizId: "",
  bizDesc: "",
});
const route = useRoute()
const router = useRouter();
const resetForm = () => {
  router.push("/main/biz/");
};

const formRef = ref();
const labelCol = { span: 2 };
const wrapperCol = { span: 6 };
const rules: Record<string, Rule[]> = {
  name: [
    { required: true, message: 'Please input Activity name', trigger: 'change' },
    { min: 3, max: 5, message: 'Length should be 3 to 5', trigger: 'blur' },
  ],
  region: [{ required: true, message: 'Please select Activity zone', trigger: 'change' }],
  date1: [{ required: true, message: 'Please pick a date', trigger: 'change', type: 'object' }],
  type: [
    {
      type: 'array',
      required: true,
      message: 'Please select at least one activity type',
      trigger: 'change',
    },
  ],
  resource: [{ required: true, message: 'Please select activity resource', trigger: 'change' }],
  desc: [{ required: true, message: 'Please input activity form', trigger: 'blur' }],
};
const saveBizInfo = () => {
  request({
    url: "/biz",
    method: bizObj.id === 0 ? "POST" : "PUT",
    data: bizObj,
  }).then((resp) => {

    notification.success({
      message: '保存成功',
      description:
        '记录已保存。',
      duration: 2,
      onClick: () => {
        // console.log('Notification Clicked!');
      },
    });
    // 接口返回数据之后的操作
    let data = resp.data;
    bizObj.id = 0;
    bizObj.bizName = "";
    bizObj.bizId = "";
    bizObj.bizDesc = "";
    router.push("/main/biz");
  });
};

onMounted(() => {
  if (route.query.bizId) {
    let id = parseInt(route.query.bizId as string)
    if (id > 0) {
      request({
        url: `/biz/${id}`,
        method: "get",
      }).then(resp => {
        const data = resp["data"]
        bizObj.id = data['id']
        bizObj.bizId = data['bizId']
        bizObj.bizName = data['bizName']
        bizObj.bizDesc = data['bizDesc']
      })
    }
  }

})

</script>

<style>
.systemPage {
  width: 100%;
}

.ant-form-text {
  font-weight: bold;
}
</style>