<template>
  <a-layout-content :style="style">
    <div class="opBar">
      <a-button type="primary" class="op-item" @click="toAddUser" size="small"
        v-if="props.stage === 'edition'">添加</a-button>
      <a-input-search v-model:value="searchKeyword" class="op-item keyword" size="small" enter-button
        placeholder="请输入检索内容" @search="search" />
    </div>
    <a-table size="small" :data-source="tabData" :columns="columns" :pagination="pagination">
      <template #bodyCell="{ text, column, record }">
        <template v-if="column.dataIndex === 'operation' && props.stage !== 'selecting'">
          <span>
            <a-popconfirm title="是否确认删除该记录?" ok-text="确认" cancel-text="取消" @confirm="onDelete(record.id)">
              <a-button danger type="link" size="small">删除</a-button>
            </a-popconfirm>

          </span>
        </template>
        <template v-else-if="column.key === 'parent'">
          {{ record.parentTitle }}
        </template>
        <template v-else-if="stage === 'selecting' && (column.key === 'identity' || column.key === 'title')">
            <a-button type="link" size="small" @click="onSelect(record)">{{ column.key === 'identity' ? record.identity : record.title }}</a-button>
        </template>
        <template v-else-if="column.key === 'createTimestamp'">
          {{ timestamp2DateString(record.createTimestamp) }}
        </template>
        <template v-else-if="column.key === 'updateTimestamp'">
          {{ timestamp2DateString(record.updateTimestamp) }}
        </template>
        <template v-else-if="column.key === 'permissionLst'">
          {{ showPermissionLst(record.permissionLst) }}
        </template>
      </template>
    </a-table>
  </a-layout-content>
</template>

<script lang="ts" setup>
import { computed, onMounted, reactive, ref, defineProps } from 'vue';
import type { TableColumnType } from "ant-design-vue";
import { useRoute, useRouter } from 'vue-router'

import { Permission, Role, showPermissionLst } from "./user"
import { timestamp2DateString } from '@/utils/common';

const route = useRoute(),
  router = useRouter(),
  props = defineProps({
    "stage": {
      type: String,
      default: "edition"
    }
  })

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
    title: "Identity",
    dataIndex: "identity",
    key: "identity",
    width: 200,
  },
  {
    title: "角色名称",
    dataIndex: "name",
    key: "name",
    width: 200,
  },
  {
    title: "权限列表",
    dataIndex: "permissionLst",
    key: "permissionLst",
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
];

if (props.stage === "edition") {
  columns.push({
    title: "操作",
    dataIndex: "operation",
    key: "operation",
    width: 260,
  })
}

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
  router.push("/main/role/add")
}

const search = () => {

}

const onDelete = (id: number) => {
  new Role({ id }).delete().then((resp) => {
    Role.loadPage(pgSetting).then((resp) => {
      const data = resp.data
      tabData.value = data
    })
  })
}

const emit = defineEmits(["select"])

const onSelect = (record: any) => {
  emit("select", new Role({id: record["id"], name: record["name"]}))
}

onMounted(() => {
  Role.loadPage(pgSetting).then((resp) => {
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