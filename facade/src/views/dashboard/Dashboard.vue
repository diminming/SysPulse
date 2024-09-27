<template>
  <div class="dashboard">
    <div class="tiles">
      <a-space size="large">
        <tile v-for="item in tileLst" :key="item.name" :color="item.color" :title="item.title"
          :icon="item.icon" :loadValue="item.loadValue" :clickHandler="item.clickHandler">
        </tile>
      </a-space>
    </div>
    <div class="charts">
      
      <a-flex justify="space-between" :align="'flex-end'">
        <alert-stat></alert-stat>
        <heat-map></heat-map>
      </a-flex>

    </div>
    <div>
      <a-card title="告警列表" size="small">
        <AlarmLst :stage="'dashboard'"></AlarmLst>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import AlarmLst from "@/views/notification/AlarmLst.vue"
import Tile from "@/views/dashboard/widget/Tile.vue"
import AlertStat from "@/views/dashboard/widget/AlertStat.vue"
import { GroupOutlined, DesktopOutlined, BarsOutlined, NotificationOutlined } from '@ant-design/icons-vue';
import HeatMap from "@/views/dashboard/widget/HeatMap.vue";

import request from "@/utils/request"
import router from "@/router";

const tileLst = [
  {
    "name": "grp",
    "title": "资源分组",
    "icon": GroupOutlined,
    "color": "#321FDB",
    "loadValue": () => {
      return new Promise((resolve, reject) => {
        request({
          url: "/biz/count",
          method: "GET"
        }).then((resp: any) => {
          resolve(resp["data"])
        })
      })
    },
    "clickHandler": () => {
      router.push({
        path: "/main/biz"
      })
    }
  },
  {
    "name": "linux",
    "title": "Linux",
    "icon": DesktopOutlined,
    "color": "#3399FF",
    "loadValue": () => {
      return new Promise((resolve, reject) => {
        request({
          url: "/linux/count",
          method: "GET"
        }).then((resp: any) => {
          resolve(resp["data"])
        })
      })
    },
    "clickHandler": () => {
      router.push({
        path: "/main/linux"
      })
    }
  },
  {
    "name": "job",
    "title": "分析任务",
    "icon": BarsOutlined,
    "color": "#F9B115",
    "loadValue": () => {
      return new Promise((resolve, reject) => {
        request({
          url: "/job/count",
          method: "GET"
        }).then((resp: any) => {
          resolve(resp["data"])
        })
      })
    },
    "clickHandler": () => {
      console.log(123)
    }
  },
  {
    "name": "alert",
    "title": "告警",
    "icon": NotificationOutlined,
    "color": "#E55353",
    "loadValue": () => {
      return new Promise((resolve, reject) => {
        request({
          url: "/alarm/count",
          method: "GET"
        }).then((resp: any) => {
          resolve(resp["data"])
        })
      })
    },
    "clickHandler": () => {
      router.push({
        path: "/main/notification"
      })
    }
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

.charts {
  display: inline-block;
}
</style>