<template>
  <a-layout-content
      :style="{
      background: '#fff',
      padding: '24px',
      margin: 0,
      minHeight: '800px',
    }"
  >
    <div class="opBar">
      <a-button type="primary" class="op-item" @click="pathTo1"
      >新增
      </a-button
      >
      <a-input-search
          class="op-item"
          v-model:value="keyword"
          placeholder="请输入检索内容"
          enter-button
          @search="onSearch"
      />
    </div>

    <a-table
        :data-source="tabData"
        :columns="columns"
        size="middle"
        :row-selection="rowSelection"
        :pagination="pgSetting"
        @change="onChange"
    >
      <template #bodyCell="{ text, column, record }">
        <template v-if="column.key === 'bizName'">
          <a @click="pathTo2">
            {{ record.bizName }}
          </a>
        </template>
        <template v-else-if="column.dataIndex === 'operation'">
          <span>
            <a-button type="link" @click="onEdite(record)">编辑</a-button>
            <a-divider type="vertical"/>
            <a-popconfirm
                title="是否确认删除该记录?"
                ok-text="确认"
                cancel-text="取消"
                @confirm="onDelete(record.id)"
            >
              <a-button danger type="link">删除</a-button>
            </a-popconfirm>
          </span>
        </template>
      </template>
    </a-table>
  </a-layout-content>
  <a-modal
      v-model:open="showEditeDialog"
      title="编辑"
      centered
      ok-text="确认"
      cancel-text="取消"
      @ok="saveBizInfo"
  >
    <a-form
        layout="horizontal"
        :model="bizObj"
        :label-col="{ style: { width: '4rem' } }"
    >
      <a-form-item label="业务名称">
        <a-input v-model:value="bizObj.bizName" placeholder="请输入业务名称"/>
      </a-form-item>
      <a-form-item label="业务标识">
        <a-input v-model:value="bizObj.bizId" placeholder="请输入业务标识"/>
      </a-form-item>
      <a-form-item label="备注">
        <a-textarea v-model:value="bizObj.bizDesc" placeholder="请输入备注"/>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import {SearchOutlined} from "@ant-design/icons-vue";
import {computed, onMounted, reactive, ref} from "vue";
import {useRouter} from "vue-router";
import request from "@/utils/request";
import type {TableColumnType, TableProps} from "ant-design-vue";

interface Business {
  id: number;
  bizName: string;
  bizId: string;
  bizDesc: string;
  createTimestamp: number;
  updateTimestamp: number;
}

const showEditeDialog = ref<boolean>(false),
    bizObj = reactive({
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

const router = useRouter();
const pathTo1 = () => {
  router.push("/main/biz/bizAdd");
};
const pathTo2 = (data: any) => {
  // router.push({ name: "bizSystem", query: { ...data } });
  router.push({ path: "/main/biz/bizSystem", query: { ...data } });
};



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

// const handleSearch = (selectedKeys, confirm, dataIndex) => {
//   confirm();
//   state.searchText = selectedKeys[0];
//   state.searchedColumn = dataIndex;
// };

// const handleReset = clearFilters => {
//   clearFilters({ confirm: true });
//   state.searchText = '';
// };

const saveBizInfo = () => {
  if (bizObj.id === 0) {
    createBizInfo();
  } else {
    updateBizInfo();
  }
};

const onEdite = (biz: Business) => {
  // debugger;
  bizObj.id = biz["id"];
  bizObj.bizName = biz["bizName"];
  bizObj.bizId = biz["bizId"];
  bizObj.bizDesc = biz["bizDesc"];
  showEditeDialog.value = true;
};

const createBizInfo = () => {
  request({
    url: "/biz",
    method: "POST",
    data: {
      bizName: bizObj.bizName,
      bizId: bizObj.bizId,
      bizDesc: bizObj.bizDesc,
    },
  }).then((resp) => {
    let data = resp.data;
    bizObj.id = 0;
    bizObj.bizName = "";
    bizObj.bizId = "";
    bizObj.bizDesc = "";
    showEditeDialog.value = false;
    getPage();
  });
};

const updateBizInfo = () => {
  request({
    url: "/biz",
    method: "PUT",
    data: {
      id: bizObj.id,
      bizName: bizObj.bizName,
      bizId: bizObj.bizId,
      bizDesc: bizObj.bizDesc,
    },
  }).then((resp) => {
    let data = resp.data;
    bizObj.id = 0;
    bizObj.bizName = "";
    bizObj.bizId = "";
    bizObj.bizDesc = "";
    showEditeDialog.value = false;
    getPage();
  });
};

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

const onDelete = (bizId: number) => {
  return new Promise((resolve) => {
    request({
      url: "/biz",
      method: "delete",
      params: {
        biz_id: bizId,
      },
    }).then(() => {
      resolve(true);
      getPage();
    });
  });
};

const onSearch = () => {
};


const rowSelection: TableProps["rowSelection"] = {
  // onChange: (selectedRowKeys: string[], selectedRows: Business[]) => {
  //   console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
  // },
  // getCheckboxProps: (record: Business) => ({
  //   disabled: record.name === 'Disabled User', // Column configuration not to be checked
  //   name: record.name,
  // }),
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
  width: 25%;
  justify-content: space-between;
}

.op-item {
  margin-right: 1rem;
}

.busPage {
  overflow: hidden;
  height: 800px;
}

.full-modal {

.ant-modal {
  max-width: 100%;
  top: 0;
  padding-bottom: 0;
  margin: 0;
}

.ant-modal-content {
  display: flex;
  flex-direction: column;
  height: calc(100vh);
}

.ant-modal-body {
  flex: 1;
}

}
</style>