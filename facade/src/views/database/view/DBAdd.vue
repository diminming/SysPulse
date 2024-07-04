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
          :label-col="labelCol"
          :wrapper-col="wrapperCol"
      >
        <a-form-item><span class="ant-form-text">新增Linux资源</span></a-form-item>
        <a-form-item ref="name" label="数据库名称" name="name">
          <a-input v-model:value="dbObj.name"/>
        </a-form-item>
        <a-form-item ref="db_id" label="数据库标识" name="db_id">
          <a-input v-model:value="dbObj.db_id"/>
        </a-form-item>
        <a-form-item label="数据库类型" name="region">
          <a-select v-model:value="dbObj.type" placeholder="请选择数据库类型">
            <a-select-option value="shanghai">Zone one</a-select-option>
            <a-select-option value="beijing">Zone two</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="业务系统">
          <a-input placeholder="请选择业务系统"/>
          <a-button type="primary" @click="showModal1">选择</a-button>
        </a-form-item>
        <a-form-item label="Linux">
          <a-input placeholder="请选择Linux主机"/>
          <a-button type="primary" @click="showModal2">选择</a-button>
        </a-form-item>
        <a-form-item :wrapper-col="{ span: 6, offset: 2 }">
          <a-button type="primary" @click="createDBInfo">确定</a-button>
          <a-button style="margin-left: 10px" @click="resetForm">返回</a-button>
        </a-form-item>
      </a-form>

      <a-modal v-model:open="open1" width="1200px" @ok="handleOk1" style="top: 20%">
        <a-input-search
            class="op-item1"
            v-model:value="keyword"
            placeholder="请输入检索内容"
            enter-button
            @search=""
        />
        <a-table
            :data-source="tabData1"
            :pagination="pgSetting1"
            :scroll="{ x: 1000, y: 400 }"
        >
        </a-table>
      </a-modal>
      <a-modal v-model:open="open2" width="1200px" @ok="handleOk2" style="top: 20%">
        <a-input-search
            class="op-item2"
            v-model:value="keyword"
            placeholder="请输入检索内容"
            enter-button
            @search=""
        />

        <a-table
            :data-source="tabData2"
            size="middle"
            :pagination="pgSetting2"
            :scroll="{ x: 1000, y: 400 }"
        >
        </a-table>
      </a-modal>
    </div>
  </a-layout-content>
</template>

<script setup lang="ts">
import {computed, onMounted, reactive, ref, toRaw} from 'vue';
import {useRouter} from "vue-router";
// import type {Rule, TableColumnType, TableProps} from 'ant-design-vue/es/form';
import request from "@/utils/request";

interface Business {
  id: number;
  bizName: string;
  bizId: string;
  bizDesc: string;
  createTimestamp: number;
  updateTimestamp: number;
}
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
interface Database {
  id: number;
  name: string;
  dbId: string;
  type: string;
  linux: {
    id: number,
    hostname: string,
    linux_id: string,
  },
  biz: {
    id: number,
    bizName: string,
    bizId: string,
    bizDesc: string,
  },
  createTimestamp: number;
  updateTimestamp: number;
}

const bizObj = reactive({
      id: 0,
      bizName: "",
      bizId: "",
      bizDesc: "",
    }),
    pagination= reactive({
      page: 0,
      pageSize: 20,
      total: 0,
    }),
    tabData1 = ref([]),
    tabData2 = ref([]),
    state = reactive({
      searchText: "",
      searchedColumn: "",
    }),
    keyword = ref();
const dbObj = reactive({
  id: 0,
  name: "",
  db_id: "",
  type: "",
  linux: {
    id: "",
    hostname: "",
    linux_id: "",
  },
  biz: {
    id: "",
    bizName: "",
    bizId: "",
    bizDesc: "",
  },
  create_timestamp: "",
  update_timestamp: ""
})


// const columns1: TableColumnType<Business>[] = [
//   {
//     title: "业务名称",
//     dataIndex: "bizName",
//     key: "bizName",
//     width: 200,
//   },
//   {
//     title: "业务标识",
//     dataIndex: "bizId",
//     key: "bizId",
//     width: 200,
//   },
//   {
//     title: "业务描述",
//     dataIndex: "bizDesc",
//     key: "bizDesc",
//     width: 300,
//   },
//   {
//     title: "操作",
//     dataIndex: "operation",
//     key: "operation",
//     width: 100,
//   },
// ];
// const columns2: TableColumnType<Linux>[] = [
//   {
//     title: "主机名",
//     dataIndex: "hostname",
//     key: "hostname",
//     width: 200,
//   },
//   {
//     title: "标识符",
//     dataIndex: "linux_id",
//     key: "linux_id",
//     width: 200,
//   },
//   {
//     title: "业务系统名称",
//     dataIndex: "biz/bizName",
//     key: "bizName",
//     width: 300,
//   },
//   {
//     title: "操作",
//     dataIndex: "operation",
//     key: "operation",
//     width: 100,
//   },
// ];

// const rowSelection1: TableProps["rowSelection"] = {};
// const rowSelection2: TableProps["rowSelection"] = {};

const router = useRouter();
const resetForm = () => {
  router.push("/main/database/");

};
const formRef = ref();
const labelCol = {span: 2};
const wrapperCol = {span: 6};
// const rules: Record<string, Rule[]> = {
//   name: [
//     {required: true, message: '请输入数据库名称', trigger: 'change'},
//     {min: 3, max: 5, message: 'Length should be 3 to 5', trigger: 'blur'},
//   ],
//   region: [{required: true, message: '请选择数据库类型', trigger: 'change'}],
//   date1: [{required: true, message: 'Please pick a date', trigger: 'change', type: 'object'}],
//   type: [
//     {
//       type: 'array',
//       required: true,
//       message: '请选择至少一个类型',
//       trigger: 'change',
//     },
//   ],
//   resource: [{required: true, message: 'Please select activity resource', trigger: 'change'}],
//   desc: [{required: true, message: 'Please input activity form', trigger: 'blur'}],
// };
const createDBInfo = () => {
  request({
    url: "/db/page",
    method: "POST",
    data: {
      "name": dbObj.name,
      "db_id": dbObj.db_id,
      "type": dbObj.type,
      "biz.bizName": dbObj.biz.bizName,
      "linux.hostname": dbObj.linux.hostname
    }
  }).then(resp => {
    let data = resp.data
    dbObj.id = 0
    dbObj.name = ""
    dbObj.db_id = ""
    dbObj.type = ""
    dbObj.biz.bizName = ""
    dbObj.linux.hostname = ""
    router.push("/main/db");
  });
};

const open1 = ref<boolean>(false);
const open2 = ref<boolean>(false);
const showModal1 = () => {
  open1.value = true;
};
const showModal2 = () => {
  open1.value = true;
};
const handleOk1 = (e: MouseEvent) => {
  console.log(e);
  open1.value = false;
};
const handleOk2 = (e: MouseEvent) => {
  console.log(e);
  open2.value = false;
};

// const rowSelection: TableProps["rowSelection"] = {};
const getPage1 = () => {
  request({
    url: "/biz/page",
    method: "GET",
    params: {
      page: pagination.page,
      pageSize: pagination.pageSize,
    },
  }).then((resp) => {
    tabData1.value = resp["data"]["lst"];
    pagination.total = resp["data"]["total"];
  });
};
const getPage2 = () => {
  request({
    url: "/linux/page",
    method: "GET",
    params: {
      page: pagination.page,
      pageSize: pagination.pageSize,
    },
  }).then((resp) => {
    tabData1.value = resp["data"]["lst"];
    pagination.total = resp["data"]["total"];
  });
};

onMounted(() => {
  getPage1();
  getPage2();
});
const pgSetting1 = computed(() => ({
  total: pagination.total,
  current: pagination.page + 1,
  pageSize: pagination.pageSize,
  showQuickJumper: true,
  showSizeChanger: true,
  showTotal: (total: number, range: number) => {
    return range + ", 共" + total;
  },
}));
const pgSetting2 = computed(() => ({
  total: pagination.total,
  current: pagination.page + 1,
  pageSize: pagination.pageSize,
  showQuickJumper: true,
  showSizeChanger: true,
  showTotal: (total: number, range: number) => {
    return range + ", 共" + total;
  },
}));
const onChange1 = (pg: { pageSize: number; current: number }) => {
  pagination.page = pg.current - 1;
  pagination.pageSize = pg.pageSize;
  getPage1();
};
const onChange2 = (pg: { pageSize: number; current: number }) => {
  pagination.page = pg.current - 1;
  pagination.pageSize = pg.pageSize;
  getPage2();
};

</script>

<style scoped>
.ant-form-text {
  font-weight: bold;
}

.op-item1, .op-item2 {
  margin-bottom: 15px;
  width: 325px;
}
</style>