<template>
  <a-layout-content :style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '800px' }">
    <a-form layout="horizontal" :model="linux" :label-col="{ style: { width: '6rem' } }" class="edit-form">
      <a-form-item label="主机名称">
        <a-input v-model:value="linux.hostname" placeholder="请输入主机名称" />
      </a-form-item>
      <a-form-item label="标识符">
        <a-input v-model:value="linux.linux_id" placeholder="请输入标识符" />
      </a-form-item>
      <a-form-item label="业务系统">
        <a-tag v-if="biz.id > 0" :closable="true" @close="Object.assign(biz, new BizObj(-1))">{{ biz.bizName }}</a-tag>
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
  <a-modal v-model:open="showSelectingBizDialog" title="选择业务系统" width="70%">
    <Business stage="select" @select-biz="onSelectBiz"></Business>
  </a-modal>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router'

// import { message } from 'ant-design-vue';
import { notification } from 'ant-design-vue';

import { Business as BizObj, Linux } from "@/views/linux/api"
import Business from '../biz/Business.vue';


const linux = reactive(new Linux(-1))
const route = useRoute(), router = useRouter(), showSelectingBizDialog = ref(false)

const biz = reactive(new BizObj(-1))

function onCancel() {
  router.push("/main/linux")
}

const onSelectBiz = (business: BizObj) => {
  Object.assign(biz, business)
  linux.biz = business
  showSelectingBizDialog.value = false
}

function onSelect() {
  showSelectingBizDialog.value = true
}

function onSave() {
  linux.save().then(() => {

    notification.success({
      message: '保存成功',
      description:
        '记录已保存。',
      duration: 2,
      onClick: () => {
        // console.log('Notification Clicked!');
      },
    });
    router.push("/main/linux")
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

      const bizData = data['biz']
      const biz0 = new BizObj(bizData['id'], bizData['bizName'])
      linux0.biz = biz0
      Object.assign(biz, biz0)

      Object.assign(linux, linux0)
    })
  }

})

</script>
<style scoped>
.edit-form {
  width: 30rem;
}
</style>