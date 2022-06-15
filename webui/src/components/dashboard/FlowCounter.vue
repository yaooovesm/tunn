<template>
  <div>
    <el-card shadow="always" body-style="padding:0" v-loading="loading">
      <div class="title" style="margin-top: 20px;margin-bottom: 7px;">
        <div class="title-text">流量表
        </div>
      </div>

      <el-popover
          placement="bottom"
          :width="200"
          :hide-after="0"
          title="流量统计详情"
          trigger="hover"
      >
        <template #default>
          <div>
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
          </div>
        </template>
        <template #reference>
          <div style="padding: 20px 20px 20px;">
            <el-progress :color="customColors" :percentage="percentage">
              <template #default="{ percentage }">
                <div style="font-size: 12px;color: #909399">
                  <span v-if="this.limit!==0">剩余 {{ percentage }}% 可用</span>
                  <span v-else>无限制</span>
                </div>
              </template>
            </el-progress>
          </div>
        </template>
      </el-popover>
      <div style="font-size: 12px;color: #808080;text-align: right;padding: 5px 10px">
        <!--        更新于-->
        <!--        {{ $utils.FormatDate("YYYY/mm/dd HH:MM:SS", updateTime) }}&nbsp;-->
        <el-button text
                   style="font-size: 12px;height: 12px;line-height: 13px;padding: 8px 2px;transform: translateY(-1px)"
                   @click="update">刷新
        </el-button>
      </div>
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
        this.loading = false
      }).catch(() => {
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