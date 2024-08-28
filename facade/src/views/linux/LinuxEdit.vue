<template>
  <a-layout-content :style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '800px' }">
    <a-form layout="horizontal" :model="linux" :label-col="{ style: { width: '6rem' } }" class="edit-form">
      <a-form-item label="主机名称">
        <a-input v-model:value="linux.hostname" placeholder="请输入主机名称" />
      </a-form-item>
      <a-form-item label="标识符">
        <a-input v-model:value="linux.linux_id" placeholder="请输入标识符" />
      </a-form-item>
      <a-form-item label="业务系统名称">
        <a-button type="primary" @click="onSelect">选择</a-button>
      </a-form-item>
      <a-form-item label="Agent地址">
        <a-input v-model:value="linux.agent_conn" placeholder="host_addr:port" />
      </a-form-item>
      <a-form-item>
        <div style="text-align: center;">
          <a-button type="primary" @click="onSave">保存</a-button>
          <a-button @click="onCancel" style="margin-left: 2rem;">取消</a-button>
        </div>
      </a-form-item>
    </a-form>
  </a-layout-content>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router'

import { message } from 'ant-design-vue';

import { Linux } from "@/views/linux/api"


const linux = ref<Linux>(new Linux(-1))
const route = useRoute(), router = useRouter()

function onCancel() {
  router.push("/main/linux")
}

function onSelect() {

}

function onSave() {
  linux.value.save().then(() => {
    message.success('保存完成，将返回列表页...', 3, () => {
      router.push("/main/linux")
    });
  })
}

onMounted(() => {
  if (route.query.linuxId) {
    let id = parseInt(route.query.linuxId as string)
    let linux0 = new Linux(id)
    linux0.load().then((resp) => {
      const data = resp.data
      linux0.agent_conn = data['agent_conn']
      linux0.hostname = data['hostname']
      linux0.linux_id = data['linux_id']
      linux.value = linux0
    })
  }

})

</script>
<style scoped>
.edit-form {
  width: 30rem;
}
</style>