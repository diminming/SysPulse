<template>
  <a-layout-content
      :style="{
      background: '#fff',
      padding: '24px',
      margin: 0,
      minHeight: '800px',
    }"
  >
    <div class="systemPage ">
      <a-form
          ref="formRef"
          :rules="rules"
          :label-col="labelCol"
          :wrapper-col="wrapperCol"

      >
        <a-form-item><span class="ant-form-text">新增Linux资源</span></a-form-item>
        <a-form-item ref="hostname" label="主机名称" name="hostname">
          <a-input v-model:value="linuxObj.hostname"/>
        </a-form-item>
        <a-form-item ref="linux_id" label="主机标识" name="linux_id">
          <a-input v-model:value="linuxObj.linux_id"/>
        </a-form-item>
        <a-form-item label="业务系统" :wrapper-col="{ span: 3, offset: 0 }">
          <a-input placeholder="请选择业务系统"/>
          <a-button type="primary" @click="showModal">选择</a-button>
        </a-form-item>
        <!-- <a-form-item label="描述" name="desc">
          <a-textarea v-model:value="linuxObj.desc"/>
        </a-form-item> -->
        <a-form-item :wrapper-col="{ span: 6, offset: 2 }">
          <a-button type="primary" @click="createLinuxInfo">确定</a-button>
          <a-button style="margin-left: 10px" @click="resetForm">返回</a-button>
        </a-form-item>
      </a-form>

      <!--      <a-modal v-model:open="open" width="1200px" wrap-class-name="h180" @ok="handleOk" style="top: 20%" :body-style="bodystyle">-->
      <a-modal v-model:open="open" width="1200px" @ok="handleOk" style="top: 20%">
        <a-input-search
            class="op-item"
            v-model:value="keyword"
            placeholder="请输入检索内容"
            enter-button
            @search=""
        />

        <a-table
            :data-source="tabData"
            :columns="columns"
            size="middle"
            :row-selection="rowSelection"
            :scroll="{ x: 1000, y: 400 }"
            :pagination="pgSetting"
        >
        </a-table>
      </a-modal>
    </div>
  </a-layout-content>
</template>

<script setup lang="ts">
import {Dayjs} from 'dayjs';
import {computed, onMounted, reactive, ref, toRaw} from 'vue';
import type {UnwrapRef} from 'vue';
import type {Rule} from 'ant-design-vue/es/form';
import request from "@/utils/request";
import {useRouter} from "vue-router";
import type {TableColumnType, TableProps} from "ant-design-vue";

interface Linux {
  id: number;
  hostname: string;
  linux_id: string;
  biz: {
    id: number;
    bizName: string;
    bizId: string;
    bizDesc: string;
    createTimestamp: number;
    updateTimestamp: number;
  };
  createTimestamp: number;
  updateTimestamp: number;
}

interface Business {
  id: number;
  bizName: string;
  bizId: string;
  bizDesc: string;
  createTimestamp: number;
  updateTimestamp: number;
}

const bizObj = reactive({
      id: 0,
      bizName: "",
      bizId: "",
      bizDesc: "",
    }),
    pagination = reactive({
      page: 0,
      pageSize: 20,
      total: 0,
    }),
    tabData = ref([]),
    state = reactive({
      searchText: "",
      searchedColumn: "",
    }),
    keyword = ref();
const columns: TableColumnType<Business>[] = [
  {
    title: "业务名称",
    dataIndex: "bizName",
    key: "bizName",
    width: 200,
  },
  {
    title: "业务标识",
    dataIndex: "bizId",
    key: "bizId",
    width: 200,
  },
  {
    title: "业务描述",
    dataIndex: "bizDesc",
    key: "bizDesc",
    width: 300,
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation",
    width: 100,
  },
];

const linuxObj = reactive({
  id: 0,
  hostname: "",
  linux_id: "",
  biz: {
    id: "",
    bizName: "",
    bizId: "",
    bizDesc: "",
  },
  createTimestamp: "",
  updateTimestamp: ""
});

const rowSelection: TableProps["rowSelection"] = {};
const getPage = () => {
  request({
    url: "/biz/page",
    method: "GET",
    params: {
      page: pagination.page,
      pageSize: pagination.pageSize,
    },
  }).then((resp) => {
    tabData.value = resp["data"]["lst"];
    pagination.total = resp["data"]["total"];
  });
};

onMounted(() => {
  getPage();
});
const pgSetting = computed(() => ({
  total: pagination.total,
  current: pagination.page + 1,
  pageSize: pagination.pageSize,
  showQuickJumper: true,
  showSizeChanger: true,
  showTotal: (total: number, range: number) => {
    return range + ", 共" + total;
  },
}));

const onChange = (pg: { pageSize: number; current: number }) => {
  pagination.page = pg.current - 1;
  pagination.pageSize = pg.pageSize;
  getPage();
};
const formRef = ref();
const labelCol = {span: 2};
const wrapperCol = {span: 6};
const rules: Record<string, Rule[]> = {
  name: [
    {required: true, message: 'Please input Activity name', trigger: 'change'},
    {min: 3, max: 5, message: 'Length should be 3 to 5', trigger: 'blur'},
  ],
  region: [{required: true, message: 'Please select Activity zone', trigger: 'change'}],
  date1: [{required: true, message: 'Please pick a date', trigger: 'change', type: 'object'}],
  type: [
    {
      type: 'array',
      required: true,
      message: 'Please select at least one activity type',
      trigger: 'change',
    },
  ],
  resource: [{required: true, message: 'Please select activity resource', trigger: 'change'}],
  desc: [{required: true, message: 'Please input activity form', trigger: 'blur'}],
};

const router = useRouter();
const resetForm = () => {
  router.push("/main/linux");
};

const createLinuxInfo = () => {
  request({
    url: "/linux/page",
    method: "POST",
    data: {
      hostname: linuxObj.hostname,
      linux_id: linuxObj.linux_id,
      bizName: linuxObj.biz.bizName,
    },
  }).then((resp) => {
    let data = resp.data;
    linuxObj.id = 0;
    linuxObj.hostname = "";
    linuxObj.linux_id = "";
    linuxObj.biz.bizName = "";
    router.push("/main/linux");
  });
};

const open = ref<boolean>(false);

const showModal = () => {
  open.value = true;
};

const handleOk = (e: MouseEvent) => {
  console.log(e);
  open.value = false;
};

// const bodystyle = {
//   height: '480px',
//   overflow: 'hidden',
//   overflowY: 'scroll',
// }
</script>

<style scoped>
.ant-form-text {
  font-weight: bold;
}

.systemPage {
  width: 100%;
}

.op-item {
  margin-bottom: 15px;
  width: 325px;
}


</style>