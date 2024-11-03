<template>
  <a-layout-content :style="{
    background: '#fff',
    padding: '24px',
    margin: 0,
    minHeight: '800px',
  }">
    <div class="opBar">
      <a-button type="primary" class="op-item" @click="pathTo1">新增</a-button>
      <a-input-search class="op-item keyword" v-model:value="keyword" placeholder="请输入检索内容" enter-button @search="onSearch"/>
    </div>

    <a-table :data-source="tabData" :columns="columns" size="small" 
             :pagination="pgSetting"
             @change="onChange">
      <template #bodyCell="{ text, column, record }">
        <template v-if="column.key === 'hostname'">
          <a @click="gotoLinuxDetail(record)">
            {{ record.hostname }}
          </a>
        </template>
        <template v-else-if="column.key === 'bizName'">
          {{ record.biz.bizName }}
        </template>
        <template v-else-if="column.dataIndex === 'operation'">
          <span>
            <a-button type="link" @click="onEdite(record)" size="small">编辑</a-button>
            <a-divider type="vertical"/>
            <a-popconfirm title="是否确认删除该记录?" ok-text="确认" cancel-text="取消" @confirm="onDelete(record.id)">
              <a-button danger type="link" size="small">删除</a-button>
            </a-popconfirm>
          </span>
        </template>
      </template>
    </a-table>
  </a-layout-content>
</template>


<script setup lang="ts">
import {computed, onMounted, reactive, ref} from "vue";
import {useRouter} from "vue-router";
import request from "@/utils/request";
import type {TableColumnType, TableProps} from "ant-design-vue";
import type {Linux} from "@/views/linux/api"

const showEditeDialog = ref<boolean>(false),
    linuxObj = reactive({
      id: 0,
      hostname: "",
      linux_id: "",
      biz: {
        id: "",
        bizName: "",
        bizId: "",
        bizDesc: "",
      },
      agent_conn: "",
      createTimestamp: "",
      updateTimestamp: ""
    }),
    pagination = reactive({
      page: 0,
      pageSize: 10,
      total: 0,
    }),
    tabData = ref([]),
    state = reactive({
      searchText: "",
      searchedColumn: "",
    }),
    keyword = ref("");


const columns: TableColumnType<Linux>[] = [
  {
    title: "主机名",
    dataIndex: "hostname",
    key: "hostname",
    width: 200,
  },
  {
    title: "标识符",
    dataIndex: "linux_id",
    key: "linux_id",
    width: 200,
  },
  {
    title: "业务系统名称",
    dataIndex: "biz/bizName",
    key: "bizName",
    width: 300,
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation",
    width: 100,
  },
];

const rowSelection: TableProps["rowSelection"] = {};

const router = useRouter();
const pathTo1 = () => {
  router.push("/main/linux/add");
};
const gotoLinuxDetail = (data: Linux) => {
  // router.push({ name: "linuxSystem", query: { ...data.biz, ...data } });
  // router.push({path: "/main/linux/detail", query: {linuxId: data.id}});
  router.push({
    path: `/main/linux/${data.id}/detail`,
  });
};

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

const onSearch = () => {
  pagination.page = 0
  getPage()
};

const getPage = () => {
  request({
    url: "/linux/page",
    method: "GET",
    params: {
      page: pagination.page,
      pageSize: pagination.pageSize,
      kw: keyword.value,
    },
  }).then((resp) => {
    tabData.value = resp["data"]["lst"];
    pagination.total = resp["data"]["total"];
  });
};

onMounted(() => {
  getPage();
});

const onEdite = (linux: Linux) => {
  router.push({path: "/main/linux/edit", query: {linuxId: linux.id}});
};

const onDelete = (linux_id: number) => {
  return new Promise((resolve) => {
    request({
      url: "/linux",
      method: "delete",
      params: {
        linux_id: linux_id,
      },
    }).then(() => {
      resolve(true);
      getPage();
    });
  });
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
    showEditeDialog.value = false;
    getPage();
  });
};

const updateLinuxInfo = () => {
  request({
    url: "/linux",
    method: "PUT",
    data: {
      id: linuxObj.id,
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
    showEditeDialog.value = false;
    getPage();
  });
};

const saveLinuxInfo = () => {
  if (linuxObj.id === 0) {
    createLinuxInfo();
  } else {
    updateLinuxInfo();
  }
};
</script>


<style scoped>
.highlight {
  background-color: rgb(255, 192, 105);
  padding: 0px;
}

.opBar {
  display: flex;
  margin-bottom: 1rem;
  width: 100%;
  /* justify-content: space-between; */
}

.op-item {
  margin-right: 1rem;
}

.keyword {
  width: 30%;
}
</style>