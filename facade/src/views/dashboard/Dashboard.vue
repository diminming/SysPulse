<template>
  <div class="dashboard">
    <div class="tiles">
      <a-space size="large">
        <tile v-for="item in tile_lst" :key="item.name" :color="item.color" :title="item.title" :value="item.value" :icon="item.icon"></tile>
      </a-space>
    </div>
    <div class="charts">
        <a-flex justify="space-between" align="flex-end">
          <heat-map></heat-map>
          <alert-stat></alert-stat>
        </a-flex>


    </div>
    <div>
      <a-table :dataSource="dataSource" :columns="columns">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'level'">
          <span>
            <a-tag :color="colorMap[record.level]">{{ record.level }}</a-tag>
          </span>
        </template>
      </template>
    </a-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import Tile from "@/views/dashboard/widget/Tile.vue"
import AlertStat from "@/views/dashboard/widget/AlertStat.vue"
import { GroupOutlined, DesktopOutlined, BarsOutlined, NotificationOutlined } from '@ant-design/icons-vue';
import HeatMap from "@/views/dashboard/widget/HeatMap.vue";

const tile_lst = [
  {
    "name": "grp",
    "title": "资源分组",
    "value": 12345,
    "icon": GroupOutlined,
    "color": "#321FDB"
  },
  {
    "name": "linux",
    "title": "Linux",
    "value": 23456,
    "icon": DesktopOutlined,
    "color": "#3399FF"
  },
  {
    "name": "item",
    "title": "指标项",
    "value": 34567,
    "icon": BarsOutlined,
    "color": "#F9B115"
  },
  {
    "name": "alert",
    "title": "告警",
    "value": 45678,
    "icon": NotificationOutlined,
    "color": "#E55353"
  }
]

const colorMap: any = {
  正常: "green",
  告警: "orange",
  严重告警: "red",
};

const dataSource = [
  {
    time: "1",
    age: "Node1",
    level: "正常",
    contant: "12111",
  },
  {
    time: "1",
    age: "Node2",
    level: "告警",
    contant: "12111",
  },
  {
    time: "1",
    age: "Node3",
    level: "严重告警",
    contant: "12111",
  },
  {
    time: "1",
    age: "Node4",
    level: "正常",
    contant: "12111",
  },
];

const columns = [
  {
    title: "时间",
    dataIndex: "time",
    key: "time",
  },
  {
    title: "对象",
    dataIndex: "age",
    key: "age",
  },
  {
    title: "告警级别",
    dataIndex: "level",
    key: "level",
  },
  {
    title: "告警内容",
    dataIndex: "contant",
    key: "contant",
  },
];

</script>

<style scoped>

.dashboard {
  background-color: #fff;
  padding: 1rem;
}
.cardList {
  width: 95%;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.cardItem {
  width: 23%;
  height: 140px;
  border: 1px solid #000;
  border-radius: 5px;
}

.tiles {
  width: 95%;
  margin-bottom: 20px;
}

.charts{
  display: inline-block;
}
</style>