<template>
  <a-layout-content :style="{
    background: '#fff',
    padding: '24px',
    margin: 0,
    minHeight: '800px',
  }">
    <div class="opBar">
      <a-button type="primary" class="op-item" @click="pathTo1">新增</a-button>
      <a-input-search class="op-item" v-model:value="keyword" placeholder="请输入检索内容" enter-button
                      @search="onSearch"/>
    </div>
    <a-table :data-source="tabData" :columns="columns" size="middle" :row-selection="rowSelection"
             :pagination="pgSetting" @change="onChange">
      <template #bodyCell="{ text, column, record }">
        <template v-if="column.key === 'name'">
          <a @click="pathTo2">
            {{ record.name }}
          </a>
        </template>
        <template v-else-if="column.key === 'db_id'">
          {{ record.db_id }}
        </template>
        <template v-else-if="column.key === 'type'">
          {{ record.type }}
        </template>
        <template v-else-if="column.key === 'hostname'">
          {{ record.linux.hostname }}
        </template>
        <template v-else-if="column.key === 'bizName'">
          {{ record.biz.bizName }}
        </template>
        <template v-else-if="column.dataIndex === 'operation'">
          <span>
            <a-button type="link" @click="onEdite(record)">编辑</a-button>
            <a-divider type="vertical"/>
            <a-popconfirm title="是否确认删除该记录?" ok-text="确认" cancel-text="取消" @confirm="onDelete(record.id)">
              <a-button danger type="link">删除</a-button>
            </a-popconfirm>
          </span>
        </template>
      </template>
    </a-table>
  </a-layout-content>

  <a-modal v-model:open="showEditeDialog" title="编辑" centered @ok="saveDBInfo">
    <a-form layout="horizontal" :model="dbObj" :label-col="{ style: { width: '5rem' } }">
      <a-form-item label="数据库名称">
        <a-input v-model:value="dbObj.name" placeholder="请输入数据库名称"/>
      </a-form-item>
      <a-form-item label="数据库标识">
        <a-input v-model:value="dbObj.db_id" placeholder="请输入数据库标识"/>
      </a-form-item>
      <a-form-item label="数据库类型">
        <a-textarea v-model:value="dbObj.type" placeholder="请输入数据库类型"/>
      </a-form-item>
    </a-form>
  </a-modal>
</template>


<script setup lang="ts">
import {SearchOutlined} from '@ant-design/icons-vue';
import {computed, onMounted, reactive, ref} from 'vue';
import request from "@/utils/request"
import type {TableColumnType, TableProps} from 'ant-design-vue';
import {useRouter} from "vue-router";

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

const showEditeDialog = ref<boolean>(false),
    dbObj = reactive({
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
    }),
    pagination = reactive({
      page: 0,
      pageSize: 20,
      total: 0,
    }),
    tabData = ref([]),
    state = reactive({
      searchText: '',
      searchedColumn: '',
    }),
    keyword = ref()

const onSearch = () => {
}

const columns: TableColumnType<Database>[] = [
  {
    title: "数据库名",
    dataIndex: "name",
    key: "name",
    width: 200,
  },
  {
    title: "数据库标识",
    dataIndex: "db_id",
    key: "db_id",
    width: 100,
  },
  {
    title: "数据库类型",
    dataIndex: "type",
    key: "type",
    width: 250,
  },
  {
    title: "业务系统名称",
    dataIndex: "biz.bizName",
    key: "bizName",
    width: 250,
  },
  {
    title: "数据库类型",
    dataIndex: "linux.hostname",
    key: "hostname",
    width: 200,
  },
  {
    title: '操作',
    dataIndex: 'operation',
    key: 'operation',
    width: 200,
  },
];

const rowSelection: TableProps['rowSelection'] = {};

const router = useRouter();
const pathTo1 = () => {
  router.push("/main/database/dbAdd");
};
const pathTo2 = (data: any) => {
  // router.push({ name: "dbSystem", query: { ...data.linux,...data.biz, ...data } });
  router.push({ path: "/main/database/dbSystem", query: { ...data.linux,...data.biz, ...data } });
};

const pgSetting = computed(() => ({
  total: pagination.total,
  current: pagination.page + 1,
  pageSize: pagination.pageSize,
  showQuickJumper: true,
  showSizeChanger: true,
  showTotal: ((total: number, range: number) => {
    return range + ", 共" + total
  }),
}))

const onChange = ((pg: { pageSize: number; current: number }) => {
  pagination.page = pg.current - 1
  pagination.pageSize = pg.pageSize
  getPage()
})



const onEdite = (db: Database) => {
  debugger
  dbObj.id = db['id']
  dbObj.name = db["name"]
  // dbObj.db_id = db["db_id"]
  // dbObj.biz.bizName = db["biz.biz_name"]
  // dbObj.linux.hostname = db["linux.hostname"]
  showEditeDialog.value = true
}

const getPage = () => {
  request({
    url: "/db/page",
    method: "GET",
    params: {
      page: pagination.page,
      pageSize: pagination.pageSize,
    }
  }).then(resp => {
    tabData.value = resp['data']['lst']
    pagination.total = resp['data']['total']
  })
}

onMounted(() => {
  getPage()
})

const onDelete = (db_id: number) => {
  return new Promise(resolve => {
    request({
      url: "/db/page",
      method: "delete",
      params: {
        "db_id": db_id
      }
    }).then(() => {
      resolve(true)
      getPage()
    })
  });
}

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
    dbObj.type=""
    dbObj.biz.bizName = ""
    dbObj.linux.hostname = ""
    showEditeDialog.value = false
    getPage()
  })
}

const updateDBInfo = () => {
  request({
    url: "/db/page",
    method: "PUT",
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
    dbObj.type=""
    dbObj.biz.bizName = ""
    dbObj.linux.hostname = ""
    showEditeDialog.value = false
    getPage()
  })
}

const saveDBInfo = () => {
  if (dbObj.id === 0) {
    createDBInfo()
  } else {
    updateDBInfo()
  }
}
</script>


<style scoped>
.highlight {
  background-color: rgb(255, 192, 105);
  padding: 0px;
}

.opBar {
  display: flex;
  margin-bottom: 1rem;
  width: 25%;
  justify-content: space-between;
}

.op-item {
  margin-right: 1rem;
}
</style>