<template>
  <a-layout-content :style="{
    background: '#fff',
    padding: '24px',
    margin: 0,
    minHeight: '800px',
  }">
    <a-form
        ref="formRef"

        :label-col="labelCol"
        :wrapper-col="wrapperCol"
    >
      <a-form-item label="业务实例" name="region">
        <a-select placeholder="please select your zone">
          <a-select-option value="shanghai">Zone one</a-select-option>
          <a-select-option value="beijing">Zone two</a-select-option>
        </a-select>
      </a-form-item>
    </a-form>
    <a-tabs v-model:activeKey="activeKey">
      <a-tab-pane key="1" tab="Tab 1">Content of Tab Pane 1</a-tab-pane>
      <a-tab-pane key="2" tab="Tab 2" force-render>
        Content of Tab Pane 2
        <div class="charts">
<!--          <flame-graph></flame-graph>-->
          <alert-stat></alert-stat>
        </div>
      </a-tab-pane>
      <a-tab-pane key="3" tab="Tab 3">
        <a-descriptions title="系统详情" layout="vertical" bordered>
          <a-descriptions-item label="系统名称">{{  }}</a-descriptions-item>
          <a-descriptions-item label="ID">{{  }}</a-descriptions-item>
          <a-descriptions-item label="业务标识">{{  }}</a-descriptions-item>
          <a-descriptions-item label="运行状态" :span="3">
            <a-badge status="processing" text="Running"/>
          </a-descriptions-item>
          <a-descriptions-item label="创建时间">{{  }}</a-descriptions-item>
          <a-descriptions-item label="更新时间" :span="2">{{  }}</a-descriptions-item>
          <a-descriptions-item label="描述信息">
            Data disk type: MongoDB
            <br/>
            Database version: 3.4
            <br/>
            Package: dds.mongo.mid
            <br/>
            Storage space: 10 GB
            <br/>
            Replication factor: 3
            <br/>
            Region: East China 1
            <br/>
          </a-descriptions-item>
        </a-descriptions>
      </a-tab-pane>
    </a-tabs>

  </a-layout-content>
</template>

<script setup lang="ts">
import request from "@/utils/request";
import {getCurrentInstance, reactive, onMounted, ref} from 'vue'
// import FlameGraph from "@/views/biz/widget/FlameGraph.vue";
import AlertStat from "@/views/dashboard/widget/AlertStat.vue";

const activeKey = ref('1');
const labelCol = {span: 1};
const wrapperCol = {span: 4};
// const {proxy} = getCurrentInstance();
const getSystemInfo = () => {
  request({
    url: "/biz/page", // 请求地址后缀
    method: "GET",  // 请求方法 常用post和get 注意大写
    params: {     // get方法参数直接跟在地址后用params post则用data
      page: 0,
      pageSize: 20,
    },
  }).then((resp) => {
    console.log('数据', resp.data)
  });
};
getSystemInfo()

const bizInfo = reactive({obj: {}})
onMounted(() => {
  // bizInfo.obj = proxy.$route.query
  console.log('点击数据的', bizInfo.obj)

});
</script>

<style scoped>
ant-descriptions .ant-descriptions-title {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  flex: auto;
  color: rgba(0, 0, 0, 0.88);
  font-weight: 600;
  font-size: 14px;
  line-height: 1.5;
}
</style>