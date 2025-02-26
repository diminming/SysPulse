export const pathDecomposition = (pathLst: []) => {

    const verteies = new Map()
    const edges: any[] = []

    pathLst.forEach((item: any) => {

        item['vertices'].forEach((v: any) => {
            const key = v['_id']
            const info = v["info"]

            if (!verteies.has(key)) {
                let category = -1, name = "", symbolSize = 5, detail = {}

                if (key.startsWith("process/")) {
                    category = 1
                    symbolSize = 20
                    name = `${v["pid"]}:${v['info']['name']}@${v['host_identity']}`
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
                    name = info ? v['info']['hostname'] : `[${v['name']}]`
                    // detail = {
                    //   "_id": v["_id"],
                    //   "pid": v["pid"],
                    //   "name": v['info']['name'],
                    //   "exec": v["info"]['exec'],
                    //   "timestamp": v['info']['create_time']
                    // }
                } else if (key.startsWith("business/")) {
                    category = 2
                    symbolSize = 60
                    name = v['bizName']
                } else {
                  console.error(key, info)
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
            } else if(e["_id"].startsWith("res_consumption/")) {
                color = "#FFC000"
                type = "res_consumption"
            } else {
              console.error(e)
            }
            edges.push({
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

    return {
        verteies,
        edges
    }

}

export const renderTopoGraph = (chart: any, nodes: any, links: any) => {
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
          // draggable: true,
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
            gravity: 0.1,
            layoutAnimation: false,
          },
          lineStyle: {
            width: 0.8,
            // curveness: 0.3,
            opacity: 0.7
          },
          edges: links
        }
      ]
    }
    chart.setOption(option)
  }