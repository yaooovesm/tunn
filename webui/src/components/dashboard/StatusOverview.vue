<template>
  <div>
    <el-card shadow="always" body-style="padding:0">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px;">
        <div class="title-text">概况
        </div>
      </div>
      <div style="padding: 20px">
        <el-row :gutter="30">
          <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
            <el-descriptions
                direction="vertical"
                :column="3"
                size="small"
                border
            >
              <template #title>
                <div style="color: #606060;font-size: 14px">客户端</div>
              </template>
              <el-descriptions-item label-class-name="description-label" label="网卡状态" width="33.3%">
                <running-indicator style="line-height: 34px" :flag="properties.initialized" running-text="已初始化"
                                   stopped-text="未初始化"/>
              </el-descriptions-item>
              <el-descriptions-item label-class-name="description-label" label="在线状态" width="33.3%">
                <running-indicator style="line-height: 34px" :flag="properties.online" running-text="在线"
                                   stopped-text="离线"/>
              </el-descriptions-item>
              <el-descriptions-item label-class-name="description-label" label="运行状态" width="33.3%">
                <running-indicator style="line-height: 34px" :flag="properties.running" running-text="运行中"
                                   stopped-text="已停止"/>
              </el-descriptions-item>
              <el-descriptions-item label-class-name="description-label" label="版本" width="33.3%">
                {{ properties.config.runtime.app === "" ? "unknown" : properties.config.runtime.app }}
              </el-descriptions-item>
              <el-descriptions-item label-class-name="description-label" label="平台" width="66.6%">
                <template #label>
                  <el-tag size="small" effect="dark" :type="properties.config.runtime.os === ''?'info':''">
                    {{ properties.config.runtime.os === "" ? "未知平台" : properties.config.runtime.os }}
                  </el-tag>
                  <el-tag size="small" effect="dark" :type="properties.config.runtime.arch === ''?'info':''"
                          style="margin-left: 5px">
                    {{ properties.config.runtime.arch === "" ? "未知架构" : properties.config.runtime.arch }}
                  </el-tag>
                </template>
                {{ properties.config.runtime.platform === "" ? "unknown" : properties.config.runtime.platform }}
                {{ properties.config.runtime.version }}
              </el-descriptions-item>
            </el-descriptions>
          </el-col>
          <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
            <el-descriptions
                direction="vertical"
                :column="3"
                size="small"
                v-if="properties.online"
                border
            >
              <template #title>
                <div style="color: #606060;font-size: 14px">连接属性</div>
              </template>
              <el-descriptions-item label-class-name="description-label" label="内网地址" width="33.3%">
                {{ properties.config.device.cidr }}
              </el-descriptions-item>
              <el-descriptions-item label-class-name="description-label" label="传输协议" width="33.3%">
                {{ properties.config.global.protocol }}
              </el-descriptions-item>
              <el-descriptions-item label-class-name="description-label" label="数据处理" width="33.3%">
                <el-tag size="small" effect="dark" :type="properties.config.data_process.encrypt === ''?'info':''"
                        style="margin-left: 5px">
                  {{ properties.config.data_process.encrypt === "" ? "无" : properties.config.data_process.encrypt }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label-class-name="description-label" label="并发连接数" width="33.3%">
                {{ properties.config.global.multi_connection }}
              </el-descriptions-item>
              <el-descriptions-item label-class-name="description-label" label="MTU" width="33.3%">
                {{ properties.config.global.mtu }}
              </el-descriptions-item>
            </el-descriptions>
            <el-descriptions
                direction="vertical"
                :column="3"
                size="small"
                v-else
                border
            >
              <template #title>
                <div style="color: #606060;font-size: 14px">连接属性</div>
              </template>
              <el-descriptions-item label-class-name="description-label" label="" width="100%">
                <template #label>
                  属性
                </template>
                <span style="color: #909399;display: block;text-align: center">暂无数据</span>
              </el-descriptions-item>
            </el-descriptions>
          </el-col>
        </el-row>
      </div>
      <div style="font-size: 12px;color: #808080;text-align: right;padding: 5px 10px">
        更新于
        {{ $utils.FormatDate("YYYY/mm/dd HH:MM:SS", updateTime) }}&nbsp;
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
import RunningIndicator from "@/components/indicators/RunningIndicator";

export default {
  name: "StatusOverview",
  components: {RunningIndicator},
  data() {
    return {
      properties: {
        config: {
          global: {
            address: "",
            port: 0,
            protocol: "",
            mtu: 0,
            pprof: 0,
            multi_connection: 0
          },
          user: {
            Account: "",
            Password: ""
          },
          route: [],
          device: {
            cidr: "",
            dns: ""
          },
          auth: {
            Address: "",
            Port: 0
          },
          data_process: {
            encrypt: "",
            key: ""
          },
          security: {
            cert: ""
          },
          runtime: {
            os: "",
            version: "",
            arch: "",
            platform: "",
            app: ""
          },
          admin: {
            address: "",
            port: 0,
            user: "",
            password: ""
          }
        },
        error: "",
        initialized: false,
        online: false,
        running: false
      },
      loading: false,
      updateTime: new Date()
    }
  },
  mounted() {
    this.update()
  },
  methods: {
    update: function () {
      this.loading = true
      axios({
        method: "get",
        url: "/api/config/all",
        data: {}
      }).then((res) => {
        let response = res.data
        this.properties = response.data
        this.loading = false
      }).catch((err) => {
        this.loading = false
        this.$utils.HandleError(err)
      })
    }
  }
}
</script>

<style>
.description-label {
  font-size: 13px !important;
  font-weight: 500 !important;
}
</style>