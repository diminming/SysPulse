<template>
  <a-layout-content :style="style">
    <div class="opBar">
      <a-button type="primary" class="op-item" @click="toAddUser" size="small">添加</a-button>
      <a-input-search v-model:value="searchKeyword" class="op-item keyword" size="small" enter-button
        placeholder="请输入检索内容" @search="search" />
    </div>
    <a-table size="small" :data-source="tabData" :columns="columns" :pagination="pagination">
      <template #bodyCell="{ text, column, record }">
        <template v-if="column.dataIndex === 'operation'">
          <span>
            
            <a-button type="dashed" @click="onEdite(record)" size="small">重置密码</a-button>
            
            <a-divider type="vertical" />
            
            <a-button type="text" @click="onEdite(record)" size="small">禁用</a-button>

            <a-divider type="vertical" />

            <a-popconfirm title="是否确认删除该记录?" ok-text="确认" cancel-text="取消" @confirm="onDelete(record.id)">
              <a-button danger type="link" size="small">删除</a-button>
            </a-popconfirm>
            
          </span>
        </template>
        <template v-else-if="column.key === 'roleLst'">
          {{ showRoleLst(record.roleLst) }}
        </template>
        <template v-else-if="column.key === 'createTimestamp'">
          {{ timestamp2DateString(record.createTimestamp) }}
        </template>
        <template v-else-if="column.key === 'updateTimestamp'">
          {{ timestamp2DateString(record.updateTimestamp) }}
        </template>
        <template v-else-if="column.key === 'isActive' && record.isActive === true">
          <a-tag color="success">启用</a-tag>
        </template>
        <template v-else-if="column.key === 'isActive' && record.isActive === false">
          <a-tag color="warning">禁用</a-tag>
        </template>
      </template>
    </a-table>
  </a-layout-content>
</template>

<script lang="ts" setup>
import { computed, onMounted, reactive, ref } from 'vue';
import type { TableColumnType } from "ant-design-vue";
import { useRoute, useRouter } from 'vue-router'

import { User, showRoleLst } from "./user"
import { timestamp2DateString } from '@/utils/common';

const route = useRoute(),
  router = useRouter()

const searchKeyword = ref("")
const pgSetting = reactive<{
  pageNum: number,
  pageSize: number,
  total: number
}>({
  pageNum: 0,
  pageSize: 15,
  total: 0
})

const pagination = computed(() => ({
  total: pgSetting.total,
  current: pgSetting.pageNum + 1,
  pageSize: pgSetting.pageSize,
  showQuickJumper: true,
  showSizeChanger: true,
  showTotal: (total: number, range: number) => {
    return range + ", 共" + total;
  },
}));

const tabData = ref([])
const columns: TableColumnType[] = [
  {
    title: "用户名",
    dataIndex: "username",
    key: "username",
    width: 200,
  },
  {
    title: "角色",
    dataIndex: "roleLst",
    key: "roleLst",
  },
  {
    title: "状态",
    dataIndex: "isActive",
    key: "isActive",
    width: 100,
  },
  {
    title: "最后修改时间",
    dataIndex: "updateTimestamp",
    key: "updateTimestamp",
    width: 200,
  },
  {
    title: "创建时间",
    dataIndex: "createTimestamp",
    key: "createTimestamp",
    width: 200,
  },
  {
    title: "操作",
    dataIndex: "operation",
    key: "operation",
    width: 260,
  }
];

const style = reactive<{
  background: String,
  padding: String,
  margin: number,
  minHeight: string | number | undefined
}>({
  background: '#fff',
  padding: '24px',
  margin: 0,
  minHeight: undefined
})

const toAddUser = () => {
  router.push("/main/user/add")
}

const search = () => {

}

const onDelete = (id: number) => {
  new User({id}).delete().then((resp) => {
    User.loadPage(pgSetting).then((resp) => {
    const data = resp.data
    tabData.value = data
  })
  })
}

onMounted(() => {
  User.loadPage(pgSetting).then((resp) => {
    const data = resp.data
    tabData.value = data
  })
})

</script>

<style scoped>
.main_content {
  background: '#fff';
  padding: '24px';
  margin: 0;
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