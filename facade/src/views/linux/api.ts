import echarts from '@/utils/echarts';
import { CustomChart, GraphChart } from 'echarts/charts';

import { JsonResponse } from "@/utils/common";
import request from "@/utils/request"
import dayjs from 'dayjs';
import 'dayjs/locale/zh-cn';
import { _cf } from 'ant-design-vue/es/_util/cssinjs/hooks/useStyleRegister';
import { v4 as uuidv4 } from 'uuid';


echarts.use([CustomChart, GraphChart]);

class Node {
  id: string
  name: string
  value: number
  children: Array<Node>
  constructor(id: string, name: string, value: number, children: Array<Node>) {
    this.id = id
    this.name = name
    this.value = value
    this.children = children
  }
}

export class Business {
  id: number;
  bizName?: string;
  bizId?: string;
  bizDesc?: string;
  createTimestamp?: number;
  updateTimestamp?: number;
  constructor(id: number, bizName?: string) {
    this.id = id;
    this.bizName = bizName;
    // this.bizId = bizId;
    // this.bizDesc = bizDesc;
    // this.createTimestamp = createTimestamp;
    // this.updateTimestamp = updateTimestamp;
  }
};

const JOB_STATUS_MAPPING = new Map<number, string>()
JOB_STATUS_MAPPING.set(1, "已创建")
JOB_STATUS_MAPPING.set(2, "运行中")
JOB_STATUS_MAPPING.set(3, "已完成")

const JOB_TYPE_MAPPING = new Map<string, string>()
JOB_TYPE_MAPPING.set("oncpu", "On CPU")
JOB_TYPE_MAPPING.set("offcpu", "Off CPU")

export class Job {

  id: number;
  category: string;
  job_name?: string;
  linux_id?: number;
  status?: number;
  create_timestamp?: number;

  constructor(id: number, category: string) {
    this.id = id;
    this.category = category
  }

  getJobStatusTxt() {
    if (this.status) {
      const txt = JOB_STATUS_MAPPING.get(this.status)
      if (txt) {
        return txt
      }
    }
    return "Invaild Value..."
  }

  createJob() {
    return request({
      url: "/job",
      method: "POST",
      data: this
    })
  }

  getResult() {
    return request({
      url: `/job/${this.id}/result`,
      method: "GET",
    })
  }
}

export class TrafficAnalyzationJob extends Job {
  identity?: string;
  ifName?: string;
  ipAddr?: string;
  direction?: [];
  count?: Number;
  port?: Number;

  constructor(id: number) {
    super(id, "traffic")
  }
}

export function GetTrafficJobLst(page: Number, pageSize: Number, linuxId: Number) {
  return request({
    url: "/job/traffic/page",
    method: "get",
    params: {
      "page": page,
      "pageSize": pageSize,
      "linuxId": linuxId
    }
  })
}

export class ProfilingJob extends Job {

  type?: string;
  startup_time?: number;
  pid?: number;
  duration?: string;
  immediately?: boolean;

  constructor(id: number) {
    super(id, "profiling")
  }

  getImmediatelyTxt() {
    if (this.immediately === false) {
      return "否"
    } else if (this.immediately === true) {
      return "是"
    }
    return "Invaild Value..."
  }

  getCreateTimestampTxt() {
    if (this.create_timestamp) {
      return dayjs.unix(this.create_timestamp).format("YYYY/MM/DD hh:mm:ss")
    }
  }

  getTypeTxt() {
    if (this.type) {
      if (JOB_TYPE_MAPPING.has(this.type)) {
        return JOB_TYPE_MAPPING.get(this.type)
      }

    }
    return "Invaild Value..."
  }

  getStartupTimeTxt() {
    if (this.startup_time) {
      return dayjs(this.startup_time).format("YYYY/M/D HH:mm")
    }
  }

  RenderStackTraceChart() {
    const dom = document.querySelector(".graph.flame") as HTMLElement;
    const chart = echarts.init(dom);
    chart.showLoading()
    request({
      url: `/job/${this.id}/result`,
      method: "GET",
    }).then(resp => {
      const data = resp.data, root: Node = new Node(uuidv4(), "root", 0, [])
      data.forEach((e1: any) => {
        e1["Profiling"]["OnCPU"]["stacks"].forEach((e2: any) => {
          const stackSymbols = e2['stackSymbols']
          stackSymbols.reverse()
          let currNode = root
          stackSymbols.forEach((symbol: string) => {
            let target = currNode.children.find((item: any) => {
              return item["name"] === symbol
            })
            if (!target) {
              target = new Node(uuidv4(), symbol, 0, [])
              currNode.children.push(target)
            }
            currNode.value += 1
            currNode = target
          })
        })
      });
      renderFlameChart(chart, root)
    })
  }
}

const ColorTypes = new Map<string, string>()
ColorTypes.set("root", '#8fd3e8')
ColorTypes.set("genunix", '#d95850')
ColorTypes.set("unix", '#eb8146')
ColorTypes.set("ufs", '#ffb248')
ColorTypes.set("FSS", '#f2d643')
ColorTypes.set("namefs", '#ebdba4')
ColorTypes.set("doorfs", '#fcce10')
ColorTypes.set("lofs", '#b5c334')
ColorTypes.set("zfs", '#1bca93')

const filterJson = (json: any, id: any) => {
  if (id == null) {
    return json;
  }
  const recur = (item: any, id: any) => {
    if (item.id === id) {
      return item;
    }
    for (const child of item.children || []) {
      const temp = recur(child, id);
      if (temp) {
        item.children = [temp];
        item.value = temp.value; // change the parents' values
        return item;
      }
    }
  };
  return recur(json, id) || json;
};
const recursionJson = (jsonObj: any, id: any) => {
  const data: any = [];
  const filteredJson = filterJson(structuredClone(jsonObj), id);
  const rootVal = filteredJson.value;
  const recur = (item: any, start = 0, level = 0) => {
    const temp = {
      name: item.id,
      // [level, start_val, end_val, name, percentage]
      value: [
        level,
        start,
        start + item.value,
        item.name,
        (item.value / rootVal) * 100
      ],
      itemStyle: {
        color: ColorTypes.get(item.name)
      }
    };
    data.push(temp);
    let prevStart = start;
    for (const child of item.children || []) {
      recur(child, prevStart, level + 1);
      prevStart = prevStart + child.value;
    }
  };
  recur(filteredJson);
  return data;
};
const heightOfJson = (json: any) => {
  const recur = (item: any, level = 0) => {
    if ((item.children || []).length === 0) {
      return level;
    }
    let maxLevel = level;
    for (const child of item.children) {
      const tempLevel = recur(child, level + 1);
      maxLevel = Math.max(maxLevel, tempLevel);
    }
    return maxLevel;
  };
  return recur(json);
};
const renderItem = (params: any, api: any) => {
  const level = api.value(0);
  const start = api.coord([api.value(1), level]);
  const end = api.coord([api.value(2), level]);
  const height = ((api.size && api.size([0, 1])) || [0, 20])[1];
  const width = end[0] - start[0];
  return {
    type: 'rect',
    transition: ['shape'],
    shape: {
      x: start[0],
      y: start[1] - height / 2,
      width,
      height: height - 2 /* itemGap */,
      r: 2
    },
    style: {
      fill: api.visual('color')
    },
    emphasis: {
      style: {
        stroke: '#000'
      }
    },
    textConfig: {
      position: 'insideLeft'
    },
    textContent: {
      style: {
        text: api.value(3),
        fontFamily: 'Verdana',
        fill: '#000',
        width: width - 4,
        overflow: 'truncate',
        ellipsis: '..',
        truncateMinChar: 1
      },
      emphasis: {
        style: {
          stroke: '#000',
          lineWidth: 0.5
        }
      }
    }
  };
};

const renderFlameChart = (chart: any, root: Node) => {
  chart.hideLoading();
  const levelOfOriginalJson = heightOfJson(root);
  const option = {
    backgroundColor: {
      type: 'linear',
      x: 0,
      y: 0,
      x2: 0,
      y2: 1,
      colorStops: [
        {
          offset: 0.05,
          color: '#eee'
        },
        {
          offset: 0.95,
          color: '#eeeeb0'
        }
      ]
    },
    tooltip: {
      formatter: (params: any) => {
        const samples = params.value[2] - params.value[1];
        return `${params.marker} ${params.value[3]
          }: (${echarts.format.addCommas(
            samples
          )} samples, ${+params.value[4].toFixed(2)}%)`;
      }
    },
    title: [
      {
        text: 'Flame Graph',
        left: 'center',
        top: 10,
        textStyle: {
          fontFamily: 'Verdana',
          fontWeight: 'normal',
          fontSize: 20
        }
      }
    ],
    toolbox: {
      feature: {
        restore: {}
      },
      right: 20,
      top: 10
    },
    xAxis: {
      show: false
    },
    yAxis: {
      show: false,
      max: levelOfOriginalJson
    },
    series: [
      {
        type: 'custom',
        renderItem,
        encode: {
          x: [0, 1, 2],
          y: 0
        },
        data: recursionJson(root, undefined)
      }
    ]
  };
  chart.setOption(option);
  chart.on('click', (params: any) => {
    const data = recursionJson(root, params.data.name);
    const rootValue = data[0].value[2];
    chart.setOption({
      xAxis: { max: rootValue },
      series: [{ data }]
    });
  });
}

export const CaclUsage = (target: any, value: any) => {
  return (target / (value.user + value.system + value.idle + value.nice + value.iowait + value.irq + value.softirq + value.steal + value.guest + value.guestnice) * 100).toFixed(2)
}


const renderChartInDashboard = (selector: string, series: Array<any>) => {
  let dom = document.querySelector(selector),
    myChart = echarts.getInstanceByDom(dom as HTMLElement) || echarts.init(dom as HTMLElement),
    options = {
      legend: {
        type: 'scroll',
        data: series.map(item => item.name)
      },
      tooltip: {
        trigger: 'axis'
      },
      xAxis: {
        type: 'time',
        show: true,
      },
      yAxis: {
        type: 'value',
        show: true,
      },
      icon: "circle",
      grid: {
        left: '10%',   // 左边距
        right: 0,  // 右边距
        top: '15%',    // 上边距
        bottom: '10%'  // 下边距
      },
      animation: true,
      series
    }
  myChart.setOption(options)
}

export class Linux {

  id: number;
  hostname?: string;
  linux_id?: string;
  biz?: Business | null;
  agent_conn?: string;
  createTimestamp?: number;
  updateTimestamp?: number;

  constructor(id: number, hostname?: string, linux_id?: string, biz?: Business | null, agentConn?: string, createTimestamp?: number, updateTimestamp?: number) {
    this.id = id;
    this.hostname = hostname;
    this.linux_id = linux_id;
    this.biz = biz;
    this.agent_conn = agentConn;
    this.createTimestamp = createTimestamp;
    this.updateTimestamp = updateTimestamp;
  }

  load() {
    return request({
      url: `/linux/${this.id}`,
      method: "GET",
    })
  }
  save() {
    return request({
      url: this.id > 0 ? `/linux/${this.id}` : '/linux',
      method: this.id > 0 ? "PUT" : "POST",
      data: this
    })
  }
  GetProcessLst() {
    return request({
      url: `/linux/${this.id}/procLst`,
      method: "get"
    })
  }
  RefreshProcessLst() {
    return request({
      url: `/linux/${this.id}/procLst?refresh=true`,
      method: "get"
    })
  }
  GetAnalyzationJobLst(pid: number) {
    return request({
      url: `/linux/${this.id}/proc/${pid}/analyze`,
      method: "get"
    })
  }
  GetInterfaceLst() {
    return request({
      url: `/linux/${this.id}/if/lst`,
      method: "get"
    })
  }
  RenderCpuPerfChart(start: number, end: number) {
    const linuxId = this.id,
      selector = ".gallery .graph.cpu",
      dom = document.querySelector(selector),
      chart = echarts.getInstanceByDom(dom as HTMLElement) || echarts.init(dom as HTMLElement)
    chart.showLoading()
    request({
      url: "/perf/cpu",
      method: "GET",
      params: {
        linuxId,
        start,
        end
      }
    }).then((resp: any) => {
      if (resp.data && resp.data != null) {
        let jsonResp = new JsonResponse(resp.data, resp.msg, resp.status);
        let userData: any = [], systemData: any = [], idleData: any = [], iowaitData: any = [], stealData: any = [], niceData: any = [], irqData: any = [], softirqData: any = [], guestData: any = [], guestniceData: any = []
        jsonResp.Data.forEach((item: any) => {
          let timestamp: Date = new Date(item.timestamp)
          userData.push([timestamp, CaclUsage(item.user, item)])
          systemData.push([timestamp, CaclUsage(item.system, item)])
          idleData.push([timestamp, CaclUsage(item.idle, item)])
          niceData.push([timestamp, CaclUsage(item.nice, item)])
          iowaitData.push([timestamp, CaclUsage(item.iowait, item)])
          irqData.push([timestamp, CaclUsage(item.irq, item)])
          softirqData.push([timestamp, CaclUsage(item.softirq, item)])
          stealData.push([timestamp, CaclUsage(item.steal, item)])
          guestData.push([timestamp, CaclUsage(item.guest, item)])
          guestniceData.push([timestamp, CaclUsage(item.guestnice, item)])
        })
        renderChartInDashboard(selector, [
          {
            name: "user",
            type: 'line',
            smooth: false,
            symbol: 'none',
            data: userData
          }, {
            name: "system",
            type: 'line',
            smooth: false,
            symbol: 'none',
            data: systemData
          }, {
            name: "idle",
            type: 'line',
            smooth: false,
            symbol: 'none',
            data: idleData
          }, {
            name: "nice",
            type: 'line',
            smooth: false,
            symbol: 'none',
            data: niceData
          }, {
            name: "iowait",
            type: 'line',
            smooth: false,
            symbol: 'none',
            data: iowaitData
          }, {
            name: "irq",
            type: 'line',
            smooth: false,
            symbol: 'none',
            data: irqData
          }, {
            name: "softirq",
            type: 'line',
            smooth: false,
            symbol: 'none',
            data: softirqData
          }, {
            name: "steal",
            type: 'line',
            smooth: false,
            symbol: 'none',
            data: stealData
          },

          {
            name: "guest",
            type: 'line',
            smooth: false,
            symbol: 'none',
            data: guestData
          },

          {
            name: "guestnice",
            type: 'line',
            smooth: false,
            symbol: 'none',
            data: guestniceData
          }
        ])
      }
      chart.hideLoading()
    })
  }
  RenderMemoryChart(start: number, end: number) {
    const linuxId = this.id,
      selector = ".gallery .graph.mem",
      dom = document.querySelector(selector),
      chart = echarts.getInstanceByDom(dom as HTMLElement) || echarts.init(dom as HTMLElement)
    chart.showLoading()
    request({
      url: "/perf/mem",
      method: "get",
      params: {
        linuxId,
        start,
        end
      }
    }).then((resp: any) => {
      if (resp.data && resp.data != null) {
        let jsonResp = new JsonResponse(resp.data, resp.msg, resp.status);
        let totalData: any = [], availableData: any = [], freeData: any = []
        jsonResp.Data.forEach((item: any) => {
          let timestamp: Date = new Date(item.timestamp)
          totalData.push([timestamp, item.total])
          availableData.push([timestamp, item.available])
          freeData.push([timestamp, item.free])
        })
        renderChartInDashboard(".gallery .graph.mem", [
          {
            name: "Total",
            type: 'line',
            smooth: false,
            symbol: 'none',
            data: totalData
          },
          {
            name: "Available",
            type: 'line',
            smooth: false,
            symbol: 'none',
            data: availableData
          },
          {
            name: "Free",
            type: 'line',
            smooth: false,
            symbol: 'none',
            data: freeData
          },
        ])
      }
      chart.hideLoading()
    })
  }
  RenderLoadPerfChart(start: number, end: number) {
    const linuxId = this.id,
      selector = ".gallery .graph.load",
      dom = document.querySelector(selector),
      chart = echarts.getInstanceByDom(dom as HTMLElement) || echarts.init(dom as HTMLElement)
    chart.showLoading()
    request({
      url: "/perf/load",
      method: "get",
      params: {
        linuxId,
        start,
        end
      }
    }).then((resp: any) => {
      if (resp.data && resp.data != null) {
        const jsonResp = new JsonResponse(resp.data, resp.msg, resp.status),
          load1Data: any = [],
          load5Data: any = [],
          load15Data: any = []
        jsonResp.Data.forEach((item: any) => {
          let timestamp: Date = new Date(item.timestamp)
          load1Data.push([timestamp, item.load1])
          load5Data.push([timestamp, item.load5])
          load15Data.push([timestamp, item.load15])
        })
        renderChartInDashboard(selector, [
          {
            name: "Load 1",
            type: 'line',
            smooth: false,
            symbol: 'none',
            data: load1Data
          },
          {
            name: "Load 5",
            type: 'line',
            smooth: false,
            symbol: 'none',
            data: load5Data
          },
          {
            name: "Load 15",
            type: 'line',
            smooth: false,
            symbol: 'none',
            data: load15Data
          },
        ])
      }
      chart.hideLoading()
    })
  }
  RenderExchangePerfChart(start: number, end: number) {
    const linuxId = this.id,
      selector1 = ".gallery .graph.swap",
      dom1 = document.querySelector(selector1),
      chart1 = echarts.getInstanceByDom(dom1 as HTMLElement) || echarts.init(dom1 as HTMLElement),
      select2 = ".gallery .graph.swap.exchange",
      dom2 = document.querySelector(select2),
      chart2 = echarts.getInstanceByDom(dom2 as HTMLElement) || echarts.init(dom2 as HTMLElement),
      select3 = ".gallery .graph.mem.exchange",
      dom3 = document.querySelector(select3),
      chart3 = echarts.getInstanceByDom(dom3 as HTMLElement) || echarts.init(dom3 as HTMLElement)

    chart1.showLoading()
    chart2.showLoading()
    chart3.showLoading()

    request({
      url: "/perf/swap",
      method: "get",
      params: {
        linuxId,
        start,
        end
      }
    }).then((resp: any) => {
      try {
        const jsonResp = new JsonResponse(resp.data, resp.msg, resp.status)
        RenderSwapPerf(jsonResp)
        RenderSwapExchangePerf(jsonResp)
        RenderMemoryExchangePerf(jsonResp)
      } catch (e) {
        console.error(e)
      } finally {
        chart1.hideLoading()
        chart2.hideLoading()
        chart3.hideLoading()
      }

    })
  }
  RenderFileSystemPerfChart(start: number, end: number) {
    const linuxId = this.id,
      selector1 = ".gallery .graph.fs",
      dom1 = document.querySelector(selector1),
      chart1 = echarts.getInstanceByDom(dom1 as HTMLElement) || echarts.init(dom1 as HTMLElement),
      selector2 = ".gallery .graph.inode",
      dom2 = document.querySelector(selector2),
      chart2 = echarts.getInstanceByDom(dom2 as HTMLElement) || echarts.init(dom2 as HTMLElement)
    chart1.showLoading()
    chart2.showLoading()
    request({
      url: "/perf/fs",
      method: "get",
      params: {
        linuxId,
        start,
        end
      }
    }).then((resp: any) => {
      try {
        const jsonResp = new JsonResponse(resp.data, resp.msg, resp.status)
        RenderFSPerf(jsonResp)
        RenderInodePerf(jsonResp)
      } catch (exception) {
        console.error(exception)
      } finally {
        chart1.hideLoading()
        chart2.hideLoading()
      }
    })
  }
  RenderDiskIOPerfChart(start: number, end: number) {
    const linuxId = this.id,
      selector1 = ".gallery .disk.throughput",
      dom1 = document.querySelector(selector1),
      chart1 = echarts.getInstanceByDom(dom1 as HTMLElement) || echarts.init(dom1 as HTMLElement),
      selector2 = ".gallery .disk.io",
      dom2 = document.querySelector(selector2),
      chart2 = echarts.getInstanceByDom(dom2 as HTMLElement) || echarts.init(dom2 as HTMLElement)
    chart1.showLoading()
    chart2.showLoading()
    request({
      url: "/perf/disk",
      method: "get",
      params: {
        linuxId,
        start,
        end
      }
    }).then((resp: any) => {
      try {
        if (resp.data && resp.data != null) {
          const jsonResp = new JsonResponse(resp.data, resp.msg, resp.status)
          RenderDiskThroughputPerf(jsonResp)
          RenderDiskIOPerf(jsonResp)
        }
      } catch (err) {
        console.error(err)
      } finally {
        chart1.hideLoading()
        chart2.hideLoading()
      }
    })
  }
  RenderNetIOPerfChart(start: number, end: number) {
    const linuxId = this.id,
      selector1 = ".gallery .net.throughput",
      dom1 = document.querySelector(selector1),
      chart1 = echarts.getInstanceByDom(dom1 as HTMLElement) || echarts.init(dom1 as HTMLElement),
      selector2 = ".gallery .net.io",
      dom2 = document.querySelector(selector2),
      chart2 = echarts.getInstanceByDom(dom2 as HTMLElement) || echarts.init(dom2 as HTMLElement)

    chart1.showLoading()
    chart2.showLoading()

    request({
      url: "/perf/net",
      method: "get",
      params: {
        linuxId,
        start,
        end
      }
    }).then((resp: any) => {
      try {
        if (resp.data && resp.data != null) {
          const jsonResp = new JsonResponse(resp.data, resp.msg, resp.status)
          RenderNetThroughputPerf(jsonResp)
          RenderNetIOPerf(jsonResp)
        }
      } catch (e) {
        console.error(e)
      } finally {
        chart1.hideLoading()
        chart2.hideLoading()
      }
    })
  }
  RenderTopoGraph(chart: any, nodes: any, links: any) {
    let option = {
      legend: [
        {
          data: ["Linux", "Process", "Business"]
        }
      ],
      series: [
        {
          type: 'graph',
          layout: 'force',
          animation: true,
          emphasis: {
            focus: 'adjacency',
            label: {
              position: 'right',
              show: true
            }
          },
          edgeSymbol: ['circle', 'arrow'],
          edgeSymbolSize: [4, 8],
          roam: true,
          label: {
            show: true,
            position: 'right',
            formatter: '{b}'
          },
          labelLayout: {
            hideOverlap: true
          },
          draggable: true,
          data: nodes,
          categories: [{
            "name": "Linux",
            "base": "Linux",
            "keyword": {}
          }, {
            "name": "Process",
            "base": "Process",
            "keyword": {}
          }, {
            "name": "Business",
            "base": "Business",
            "keyword": {}
          }],
          force: {
            edgeLength: 200,
            repulsion: 800,
            gravity: 0.1
          },
          lineStyle: {
            width: 0.8,
            curveness: 0.3,
            opacity: 0.7
          },
          edges: links
        }
      ]
    }
    chart.setOption(option)
  }
  RenderTopological(callback: any) {
    const dom = document.querySelector("div.topo"),
      chart = echarts.getInstanceByDom(dom as HTMLElement) || echarts.init(dom as HTMLElement)
    chart.showLoading();
    request({
      url: `/linux/${this.id}/topo`,
      method: "get",
    }).then((resp: any) => {

      chart.hideLoading()

      const verteies = new Map()
      const links: any[] = []
      const infoLst = resp["data"]

      infoLst.forEach((item: any) => {

        item['vertices'].forEach((v: any) => {
          const key = v['_id']

          if (!verteies.has(key)) {
            let category = -1, name = "", symbolSize = 5, detail = {}

            if (key.startsWith("process/")) {
              category = 1
              symbolSize = 20
              name = `${v["pid"]}:${v['info']['name']}`
              detail = {
                "_id": v["_id"],
                "pid": v["pid"],
                "name": v['info']['name'],
                "exec": v["info"]['exec'],
                "timestamp": v['info']['create_time']
              }
            } else if (key.startsWith("host/")) {
              category = 0
              symbolSize = 40
              name = v['info']['hostname']
              detail = {
                "_id": v["_id"],
                "pid": v["pid"],
                "name": v['info']['name'],
                "exec": v["info"]['exec'],
                "timestamp": v['info']['create_time']
              }
            } else if (key.startsWith("business/")) {
              category = 2
              symbolSize = 60
              name = v['bizName']
            }

            verteies.set(key, {
              "name": name,
              "symbolSize": symbolSize,
              "_detail": detail,
              "category": category
            })
          }
        })

        item['edges'].forEach((e: any) => {
          let color = "", type = ""
          if (e["_id"].startsWith("deployment/")) {
            color = "#1E90FF"
            type = "depl"
          } else if (e["_id"].startsWith("conn_tcp/")) {
            color = "#FF4500"
            type = "conn_tcp"
          }
          links.push({
            "source": verteies.get(e['_from'])["name"],
            "target": verteies.get(e['_to'])["name"],
            "lineStyle": {
              "color": color
            },
            "_detail": {
              "type": type
            }
          })
        })

      })

      if (callback)
        callback(chart, verteies, links)

    })
  }
}


const RenderNetIOPerf = (resp: JsonResponse) => {
  let mapping = new Map<string, Array<any>>()
  resp.Data.forEach((item: any) => {
    let timestamp: Date = new Date(item.timestamp),
      name: string = item.name
    Object.keys(item).forEach((key: string) => {
      if (key === "name" || key === "timestamp" || key === "id" || key === "bytessent" || key === "bytesrecv")
        return
      let mapKey = name + " - " + key
      if (mapping.has(mapKey)) {
        mapping.get(mapKey)?.push([timestamp, item[key]])
      } else {
        mapping.set(mapKey, [])
      }
    })
  })

  const series: Array<any> = []
  mapping.forEach((dataLst, path) => {
    series.push({
      name: path,
      type: 'line',
      smooth: false,
      symbol: 'none',
      data: dataLst
    })
  })
  series.sort((a: any, b: any) => {
    return a['name'].toLowerCase().localeCompare(b['name'].toLowerCase());
  })
  renderChartInDashboard(".gallery .net.io", series)
}

const RenderNetThroughputPerf = (resp: JsonResponse) => {
  let mapping = new Map<string, Array<any>>()
  resp.Data.forEach((item: any) => {
    let timestamp: Date = new Date(item.timestamp),
      key1 = item.name + " - Sent", key2 = item.name + " - Recive"
    if (mapping.has(key1)) {
      mapping.get(key1)?.push([timestamp, item["bytessent"]])
    } else {
      mapping.set(key1, [])
    }
    if (mapping.has(key2)) {
      mapping.get(key2)?.push([timestamp, item["bytesrecv"]])
    } else {
      mapping.set(key2, [])
    }
  })

  const series: Array<any> = []
  mapping.forEach((dataLst, path) => {
    series.push({
      name: path,
      type: 'line',
      smooth: false,
      symbol: 'none',
      data: dataLst
    })
  })
  series.sort((a: any, b: any) => {
    return a['name'].toLowerCase().localeCompare(b['name'].toLowerCase());
  })
  renderChartInDashboard(".gallery .net.throughput", series)
}

const RenderDiskIOPerf = (resp: JsonResponse) => {
  let mapping = new Map<string, Array<any>>()
  resp.Data.forEach((item: any) => {
    let timestamp: Date = new Date(item.timestamp),
      name: string = item.name
    Object.keys(item).forEach((key: string) => {
      if (key === "name" || key === "timestamp" || key === "id" || key === "readbytes" || key === "writebytes")
        return
      let mapKey = name + " - " + key
      if (mapping.has(mapKey)) {
        mapping.get(mapKey)?.push([timestamp, item[key]])
      } else {
        mapping.set(mapKey, [])
      }
    })
  })

  const series: Array<any> = []
  mapping.forEach((dataLst, path) => {
    series.push({
      name: path,
      type: 'line',
      smooth: false,
      symbol: 'none',
      data: dataLst
    })
  })
  series.sort((a: any, b: any) => {
    return a['name'].toLowerCase().localeCompare(b['name'].toLowerCase());
  })
  renderChartInDashboard(".gallery .disk.io", series)
}

const RenderDiskThroughputPerf = (resp: JsonResponse) => {
  let mapping = new Map<string, Array<any>>()
  resp.Data.forEach((item: any) => {
    let timestamp: Date = new Date(item.timestamp),
      key1 = item.name + " - Read", key2 = item.name + " - Write"
    if (mapping.has(key1)) {
      mapping.get(key1)?.push([timestamp, item["readbytes"]])
    } else {
      mapping.set(key1, [])
    }
    if (mapping.has(key2)) {
      mapping.get(key2)?.push([timestamp, item["writebytes"]])
    } else {
      mapping.set(key2, [])
    }
  })

  const series: Array<any> = []
  mapping.forEach((dataLst, path) => {
    series.push({
      name: path,
      type: 'line',
      smooth: false,
      symbol: 'none',
      data: dataLst
    })
  })
  series.sort((a: any, b: any) => {
    return a['name'].toLowerCase().localeCompare(b['name'].toLowerCase());
  })
  renderChartInDashboard(".gallery .disk.throughput", series)
}

const RenderInodePerf = (resp: JsonResponse) => {
  let mapping = new Map<string, Array<any>>()
  resp.Data.forEach((item: any) => {
    let timestamp: Date = new Date(item.timestamp),
      path = item.path
    if (mapping.has(path)) {
      mapping.get(path)?.push([timestamp, item["inodesused"]])
    } else {
      mapping.set(path, [])
    }
  })
  const series: Array<any> = []
  mapping.forEach((dataLst, path) => {
    if (path.startsWith("/sys") || path.startsWith("/snap") || path.length === 0) {
      return
    }
    series.push({
      name: path,
      type: 'line',
      smooth: false,
      symbol: 'none',
      data: dataLst
    })
  })
  series.sort((a: any, b: any) => {
    return a['name'].toLowerCase().localeCompare(b['name'].toLowerCase());
  })
  renderChartInDashboard(".gallery .graph.inode", series)
}

const RenderFSPerf = (resp: JsonResponse) => {
  let mapping = new Map<string, Array<any>>()
  resp.Data.forEach((item: any) => {
    let timestamp: Date = new Date(item.timestamp),
      path = item.path
    if (mapping.has(path)) {
      mapping.get(path)?.push([timestamp, item["used"]])
    } else {
      mapping.set(path, [])
    }
  })
  const series: Array<any> = []
  mapping.forEach((dataLst, path) => {
    if (path.startsWith("/sys") || path.startsWith("/snap") || path.length === 0) {
      return
    }
    series.push({
      name: path,
      type: 'line',
      smooth: false,
      symbol: 'none',
      data: dataLst
    })
  })
  series.sort((a: any, b: any) => {
    return a['name'].toLowerCase().localeCompare(b['name'].toLowerCase());
  })
  renderChartInDashboard(".gallery .graph.fs", series)
}

const RenderSwapPerf = (resp: JsonResponse) => {
  let totalData: any = [], usedData: any = [], freeData: any = []
  resp.Data.forEach((item: any) => {
    let timestamp: Date = new Date(item.timestamp)
    totalData.push([timestamp, item.total])
    usedData.push([timestamp, item.used])
    freeData.push([timestamp, item.free])
  })
  renderChartInDashboard(".gallery .graph.swap", [
    {
      name: "Total",
      type: 'line',
      smooth: false,
      symbol: 'none',
      data: totalData
    },
    {
      name: "Used",
      type: 'line',
      smooth: false,
      symbol: 'none',
      data: usedData
    },
    {
      name: "Free",
      type: 'line',
      smooth: false,
      symbol: 'none',
      data: freeData
    },
  ])
}
const RenderSwapExchangePerf = (resp: JsonResponse) => {
  let sinData: any = [], soutData: any = []
  resp.Data.forEach((item: any) => {
    let timestamp: Date = new Date(item.timestamp)
    sinData.push([timestamp, item.sin])
    soutData.push([timestamp, item.sout])
  })
  renderChartInDashboard(".gallery .graph.swap.exchange", [
    {
      name: "Swap In",
      type: 'line',
      smooth: false,
      symbol: 'none',
      data: sinData
    },
    {
      name: "Swap Out",
      type: 'line',
      smooth: false,
      symbol: 'none',
      data: soutData
    },
  ])

}
const RenderMemoryExchangePerf = (resp: JsonResponse) => {
  let pginData: any = [], pgoutData: any = [], pgfaultData: any = []
  resp.Data.forEach((item: any) => {
    // debugger
    let timestamp: Date = new Date(item.timestamp)
    pginData.push([timestamp, item.pgin])
    pgoutData.push([timestamp, item.pgout])
    pgfaultData.push([timestamp, item.pgfault])
  })
  renderChartInDashboard(".gallery .graph.mem.exchange", [
    {
      name: "Page In",
      type: 'line',
      smooth: false,
      symbol: 'none',
      data: pginData
    },
    {
      name: "Page Out",
      type: 'line',
      smooth: false,
      symbol: 'none',
      data: pgoutData
    }, {
      name: "Page Fault",
      type: 'line',
      smooth: false,
      symbol: 'none',
      data: pgfaultData
    },
  ])
}

export default {
  GetLoad1Line: (linuxId: number, start: number, end: number) => {
    return request({
      url: "/perf/load/load1",
      method: "GET",
      params: {
        linuxId,
        start,
        end
      },
    })
  },
  RenderLoad1Line: (resp: JsonResponse) => {
    let data = resp.Data.map((value: any, index: number, array: any[]) => {
      return [new Date(value.timestamp), value.load1]
    })
    let dom = document.querySelector(".graph.load-mini"),
      myChart = echarts.getInstanceByDom(dom as HTMLElement) || echarts.init(dom as HTMLElement),
      options = {
        xAxis: {
          type: 'time',
          show: false,
        },
        yAxis: {
          type: 'value',
          show: false,
        },
        grid: {
          left: '0%',   // 左边距
          right: '0%',  // 右边距
          top: '0%',    // 上边距
          bottom: '0%'  // 下边距
        },
        animation: false,
        series: [
          {
            name: 'Load 1',
            type: 'line',
            smooth: true,
            symbol: 'none',
            data: data,
            silent: true
          }
        ]
      }
    myChart.setOption(options)
  },
  GetSwapUsedLine: (linuxId: number, start: number, end: number) => {
    return request({
      url: "/perf/swap/used",
      method: "GET",
      params: {
        linuxId,
        start,
        end
      }
    })
  },
  RenderSwapUsedLine: (resp: JsonResponse) => {
    let data = resp.Data.map((value: any, index: number, array: any[]) => {
      return [new Date(value.timestamp), value.used]
    })
    let dom = document.querySelector(".graph.swap-mini"),
      myChart = echarts.getInstanceByDom(dom as HTMLElement) || echarts.init(dom as HTMLElement),
      options = {
        xAxis: {
          type: 'time',
          show: false,
        },
        yAxis: {
          type: 'value',
          show: false,
        },
        grid: {
          left: '0%',   // 左边距
          right: '0%',  // 右边距
          top: '0%',    // 上边距
          bottom: '0%'  // 下边距
        },
        animation: false,
        series: [
          {
            name: 'Load 1',
            type: 'line',
            smooth: true,
            symbol: 'none',
            data: data,
            silent: true
          }
        ]
      }
    myChart.setOption(options)
  },
  GetAvailableMemoryLine: (linuxId: number, start: number, end: number) => {
    return request({
      url: "/perf/mem/available",
      method: "GET",
      params: {
        linuxId,
        start,
        end
      }
    })
  },
  RenderAvailableMemoryLine: (resp: JsonResponse) => {
    let data = resp.Data.map((value: any, index: number, array: any[]) => {
      return [new Date(value.timestamp), value.available]
    })
    let dom = document.querySelector(".graph.mem-available-mini"),
      myChart = echarts.getInstanceByDom(dom as HTMLElement) || echarts.init(dom as HTMLElement),
      options = {
        xAxis: {
          type: 'time',
          show: false,
        },
        yAxis: {
          type: 'value',
          show: false,
        },
        grid: {
          left: '0%',   // 左边距
          right: '0%',  // 右边距
          top: '0%',    // 上边距
          bottom: '0%'  // 下边距
        },
        animation: false,
        series: [
          {
            name: 'Load 1',
            type: 'line',
            smooth: true,
            symbol: 'none',
            data: data,
            silent: true
          }
        ]
      }
    myChart.setOption(options)
  },
  GetCpuUsageLine: (linuxId: number, start: number, end: number) => {
    return request({
      url: "/perf/cpu/usage",
      method: "GET",
      params: {
        linuxId,
        start,
        end
      }
    })
  },
  RenderCpuUsageLine: (resp: JsonResponse) => {
    let data = resp.Data.map((value: any, index: number, array: any[]) => {
      return [new Date(value.timestamp), CaclUsage(value["user"], value)]
    })
    let dom = document.querySelector(".graph.cpu-usage-mini"),
      myChart = echarts.getInstanceByDom(dom as HTMLElement) || echarts.init(dom as HTMLElement),
      options = {
        xAxis: {
          type: 'time',
          show: false,
        },
        yAxis: {
          type: 'value',
          show: false,
        },
        grid: {
          left: '0%',   // 左边距
          right: '0%',  // 右边距
          top: '0%',    // 上边距
          bottom: '0%'  // 下边距
        },
        animation: false,
        series: [
          {
            name: 'Load 1',
            type: 'line',
            smooth: true,
            symbol: 'none',
            data: data,
            silent: true
          }
        ]
      }
    myChart.setOption(options)
  },
  GetDiskIOLine: (linuxId: number, start: number, end: number) => {
    return request({
      url: "/perf/disk/iocount",
      method: "GET",
      params: {
        linuxId,
        start,
        end
      }
    })
  },
  RenderDiskWriteCounterLine: (resp: JsonResponse) => {
    let data = resp.Data.map((value: any, index: number, array: any[]) => {
      return [new Date(value.timestamp), value.writecount]
    })
    let dom = document.querySelector(".graph.disk-write-counter-mini"),
      myChart = echarts.getInstanceByDom(dom as HTMLElement) || echarts.init(dom as HTMLElement),
      options = {
        xAxis: {
          type: 'time',
          show: false,
        },
        yAxis: {
          type: 'value',
          show: false,
        },
        grid: {
          left: '0%',   // 左边距
          right: '0%',  // 右边距
          top: '0%',    // 上边距
          bottom: '0%'  // 下边距
        },
        animation: false,
        series: [
          {
            name: 'Load 1',
            type: 'line',
            smooth: true,
            symbol: 'none',
            data: data,
            silent: true
          }
        ]
      }
    myChart.setOption(options)
  },
  RenderDiskReadCounterLine: (resp: JsonResponse) => {
    let data = resp.Data.map((value: any, index: number, array: any[]) => {
      return [new Date(value.timestamp), value.readcount]
    })
    let dom = document.querySelector(".graph.disk-read-counter-mini"),
      myChart = echarts.getInstanceByDom(dom as HTMLElement) || echarts.init(dom as HTMLElement),
      options = {
        xAxis: {
          type: 'time',
          show: false,
        },
        yAxis: {
          type: 'value',
          show: false,
        },
        grid: {
          left: '0%',   // 左边距
          right: '0%',  // 右边距
          top: '0%',    // 上边距
          bottom: '0%'  // 下边距
        },
        animation: false,
        series: [
          {
            name: 'Load 1',
            type: 'line',
            smooth: true,
            symbol: 'none',
            data: data,
            silent: true
          }
        ]
      }
    myChart.setOption(options)
  },
  RenderDiskIOCounterLine: function (resp: JsonResponse) {
    this.RenderDiskReadCounterLine(resp)
    this.RenderDiskWriteCounterLine(resp)
  },
  GetIfIOLine: (linuxId: number, start: number, end: number) => {
    return request({
      url: "/perf/if/iocount",
      method: "GET",
      params: {
        linuxId,
        start,
        end
      }
    })
  },
  RenderIfRecvCounterLine: (resp: JsonResponse) => {
    let data = resp.Data.map((value: any, index: number, array: any[]) => {
      return [new Date(value.timestamp), value.bytesrecv]
    })
    let dom = document.querySelector(".graph.if-recv-mini"),
      myChart = echarts.getInstanceByDom(dom as HTMLElement) || echarts.init(dom as HTMLElement),
      options = {
        xAxis: {
          type: 'time',
          show: false,
        },
        yAxis: {
          type: 'value',
          show: false,
        },
        grid: {
          left: '0%',   // 左边距
          right: '0%',  // 右边距
          top: '0%',    // 上边距
          bottom: '0%'  // 下边距
        },
        animation: false,
        series: [
          {
            name: 'Load 1',
            type: 'line',
            smooth: true,
            symbol: 'none',
            data: data,
            silent: true
          }
        ]
      }
    myChart.setOption(options)
  },
  RenderIfSentCounterLine: (resp: JsonResponse) => {
    let data = resp.Data.map((value: any, index: number, array: any[]) => {
      return [new Date(value.timestamp), value.bytessent]
    })
    let dom = document.querySelector(".graph.if-sent-mini"),
      myChart = echarts.getInstanceByDom(dom as HTMLElement) || echarts.init(dom as HTMLElement),
      options = {
        xAxis: {
          type: 'time',
          show: false,
        },
        yAxis: {
          type: 'value',
          show: false,
        },
        grid: {
          left: '0%',   // 左边距
          right: '0%',  // 右边距
          top: '0%',    // 上边距
          bottom: '0%'  // 下边距
        },
        animation: false,
        series: [
          {
            name: 'Load 1',
            type: 'line',
            smooth: true,
            symbol: 'none',
            data: data,
            silent: true
          }
        ]
      }
    myChart.setOption(options)
  },
  RenderIfIOCounterLine: function (resp: JsonResponse) {
    this.RenderIfSentCounterLine(resp)
    this.RenderIfRecvCounterLine(resp)
  },

  GetFileSystemPerf: (linuxId: number, start: number, end: number) => {
    return request({
      url: "/perf/fs",
      method: "get",
      params: {
        linuxId,
        start,
        end
      }
    })
  },
  RenderFsUsagePerf: (resp: JsonResponse) => {
    let pginData: any = [], pgoutData: any = [], pgfaultData: any = []
    resp.Data.forEach((item: any) => {
      let timestamp: Date = new Date(item.timestamp)
      pginData.push([timestamp, item.pgin])
      pgoutData.push([timestamp, item.pgout])
      pgfaultData.push([timestamp, item.pgfault])
    })
    renderChartInDashboard(".gallery .graph.fs", [
      {
        name: "Page In",
        type: 'line',
        smooth: false,
        symbol: 'none',
        data: pginData
      },
      {
        name: "Page Out",
        type: 'line',
        smooth: false,
        symbol: 'none',
        data: pgoutData
      }, {
        name: "Page Fault",
        type: 'line',
        smooth: false,
        symbol: 'none',
        data: pgfaultData
      },
    ])
  },

}