<template>
  <a-layout-content :style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '800px' }">
    <div class="opBar">
      <a-button type="primary" class="op-item" @click="gotoAddBiz" v-if="!isSelectingStage()">新增</a-button>
      <a-input-search class="op-item keyword" v-model:value="keyword" placeholder="请输入检索内容" enter-button
        @search="onSearch" />
    </div>

    <a-table :data-source="tabData" :columns="columns" size="small" :pagination="pgSetting" @change="onChange">
      <template #bodyCell="{ text, column, record }">
        <template v-if="column.key === 'bizName'">
          <a-button type="link" size="small" @click="isSelectingStage() ? selectBiz(record) : gotoBizDetail(record)">
            {{ record.bizName }}
          </a-button>
        </template>
        <template v-else-if="column.dataIndex === 'operation' && !isSelectingStage()">
          <span>
            <a-button type="link" @click="onEdite(record)" size="small">编辑</a-button>
            <a-divider type="vertical" />
            <a-popconfirm title="是否确认删除该记录?" ok-text="确认" cancel-text="取消" @confirm="onDelete(record.id)">
              <a-button danger type="link" size="small">删除</a-button>
            </a-popconfirm>
          </span>
        </template>
      </template>
    </a-table>
  </a-layout-content>
  <a-modal v-model:open="showEditeDialog" title="编辑" centered ok-text="确认" cancel-text="取消" @ok="saveBizInfo">
    <a-form layout="horizontal" :model="bizObj" :label-col="{ style: { width: '4rem' } }">
      <a-form-item label="业务名称">
        <a-input v-model:value="bizObj.bizName" placeholder="请输入业务名称" />
      </a-form-item>
      <a-form-item label="业务标识">
        <a-input v-model:value="bizObj.bizId" placeholder="请输入业务标识" />
      </a-form-item>
      <a-form-item label="备注">
        <a-textarea v-model:value="bizObj.bizDesc" placeholder="请输入备注" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { SearchOutlined } from "@ant-design/icons-vue";
import { computed, onMounted, reactive, ref, defineProps, defineEmits } from "vue";
import { useRouter } from "vue-router";
import request from "@/utils/request";
import type { TableColumnType, TableProps } from "ant-design-vue";
import { Business } from "@/views/linux/api";

const props = defineProps({
  stage: String
}), showEditeDialog = ref<boolean>(false),
  bizObj = reactive({
    id: 0,
    bizName: "",
    bizId: "",
    bizDesc: "",
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
const gotoAddBiz = () => {
  router.push("/main/biz/bizAdd");
};

const gotoBizDetail = (data: any) => {
  // router.push({ name: "bizSystem", query: { ...data } });
  // router.push({
  //   path: `/main/biz/${data.id}/detail`,
  // });
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
  },
  {
    // title: "操作",
    dataIndex: "operation",
    key: "operation",
    width: 180,
  },
];

const isSelectingStage = () => {
  return props.stage === "select"
}

const emit = defineEmits(["selectBiz"])

const selectBiz = (record: any) => {
  emit("selectBiz", new Business(record.id, record.bizName))
}

const saveBizInfo = () => {
  if (bizObj.id === 0) {
    createBizInfo();
  } else {
    updateBizInfo();
  }
};

const onEdite = (biz: Business) => {
  router.push({ path: "/main/biz/edit", query: { bizId: biz.id } });
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
  pagination.page = 0
  getPage()
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

<style scoped lang="css">
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