<template>
  <div>
    <el-card shadow="always" body-style="padding:0">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px;">
        <div class="title-text">连接属性
        </div>
      </div>
      <div style="padding: 20px">
        <el-row :gutter="30">
          {{ properties }}
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

export default {
  name: "StatusOverview",
  data() {
    return {
      properties: {
        config: {
          global: {
            address: "",
            port: 0,
            protocol: "",
            mtu: 0,
            multi_connection: 0
          },
          user: {
            Account: "",
            Password: ""
          },
          route: null,
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
            key: null
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
        running: false
      },
      loading: false,
      updateTime: new Date()
    }
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

<style scoped>

</style>