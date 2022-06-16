<template>
  <div>
    <el-card shadow="always" body-style="padding:0" v-loading="loading">
      <div class="title" style="margin-top: 20px;margin-bottom: 7px;">
        <div class="title-text">流量表
        </div>
      </div>

      <el-popover
          placement="bottom"
          :width="300"
          :hide-after="0"
          :title="error?'':'流量统计详情'"
          trigger="hover"
      >
        <template #default>
          <div v-if="error" v-loading="loading">
            <div style="font-size: 12px;padding: 10px;color: #bbbbbb;text-align: center">无数据</div>
          </div>
          <div v-else v-loading="loading">
            <div class="detail-unit">
              <span>发送流量 </span>
              {{ $utils.FormatBytesSizeM(flow.tx) }}
            </div>
            <div class="detail-unit">
              <span>接收流量 </span>
              {{ $utils.FormatBytesSizeM(flow.rx) }}
            </div>
            <div class="detail-unit">
              <span>流量总计 </span>
              {{ $utils.FormatBytesSizeM(flow.rx + flow.tx) }}
            </div>
            <el-divider style="margin-top: 10px;margin-bottom: 10px"/>
            <div class="detail-unit">
              <span>流量限制 </span>
              {{ limit === 0 ? "无限制" : $utils.FormatBytesSizeM(limit) }}
            </div>
            <div class="detail-unit">
              <span>流量剩余 </span>
              {{ limit === 0 ? "无限制" : $utils.FormatBytesSizeM(limit - (flow.rx + flow.tx)) }}
            </div>
            <el-divider style="margin-top: 10px;margin-bottom: 10px"/>
            <div class="detail-unit">
              <span style="color: #909399">
                更新于
              {{ $utils.FormatDate("YYYY/mm/dd HH:MM:SS", updateTime) }}
              </span>
              <el-button text
                         style="font-size: 12px;height: 12px;line-height: 13px;padding: 8px 2px;transform: translateY(-1px);float: right"
                         @click="update">刷新
              </el-button>
            </div>
          </div>
        </template>
        <template #reference>
          <div style="padding: 20px 20px 30px;">
            <el-progress :color="customColors" :percentage="percentage" :show-text="false">
              <!--              <template #default="{ percentage }">-->
              <!--                <div style="font-size: 12px;color: #909399">-->
              <!--                  <div v-if="error">-->
              <!--                    <span>无数据</span>-->
              <!--                  </div>-->
              <!--                  <div v-else>-->
              <!--                    <span v-if="this.limit!==0">剩余 {{ percentage }}% 可用</span>-->
              <!--                    <span v-else>无限制</span>-->
              <!--                  </div>-->
              <!--                </div>-->
              <!--              </template>-->
            </el-progress>
            <div style="color: #909399;font-size: 12px;text-align: left">
              <div v-if="!error && limit !== 0" style="margin-top: 3px">
                可用流量 <span style="color: #007bbb">{{
                  $utils.FormatBytesSizeM(limit - (flow.rx + flow.tx)).replaceAll("M", "")
                }} </span>M ({{ percentage }}%)
              </div>
              <div v-else style="margin-top: 3px">
                <span v-if="error">
                  无数据
                </span>
                <span v-else>
                  流量无限制
                </span>
              </div>
            </div>
          </div>
        </template>
      </el-popover>
    </el-card>

  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "FlowCounter",
  data() {
    return {
      loading: false,
      timer: undefined,
      error: false,
      limit: 0,
      flow: {
        rx: 0,
        tx: 0,
      },
      percentage: 0,
      updateTime: new Date(),
      customColors: [
        {color: '#f56c6c', percentage: 20},
        {color: '#e6a23c', percentage: 40},
        {color: '#1989fa', percentage: 60},
        {color: '#5cb87a', percentage: 80},
        {color: '#5cb87a', percentage: 100},
      ]
    }
  },
  mounted() {
    this.update()
    this.timer = setInterval(() => {
      this.update()
    }, 60000)
  },
  unmounted() {
    clearInterval(this.timer)
  },
  methods: {
    update: function () {
      this.loading = true
      axios({
        method: "get",
        url: "/api/remote/flow",
        data: {}
      }).then(res => {
        let response = res.data
        this.flow = response.data
        if (this.limit === 0) {
          this.percentage = 100
        } else {
          let pct = ((this.limit - (this.flow.rx + this.flow.tx)) / this.limit) * 100
          pct = pct < 0 ? 0 : pct > 100 ? 100 : pct.toFixed(1)
          this.percentage = Number(pct)
        }
        this.error = false
        this.loading = false
      }).catch(() => {
        this.percentage = 0
        this.flow = {
          rx: 0,
          tx: 0
        }
        this.error = true
        this.loading = false
      })
    }
  }
}
</script>

<style scoped>
.detail-unit {
  text-align: right;
  font-size: 12px;
  color: #007bbb;
}

.detail-unit span {
  color: #404040;
  float: left;
  display: inline-block;
}
</style>