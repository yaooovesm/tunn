<template>
  <div>
    <div style="box-shadow: 1px 1px 4px rgba(50, 50, 50, 0.2);border-radius: 4px;">
      <el-alert style="margin-bottom: 10px" v-if="error==='logout'||error === 'terminated'"
                :title="error==='logout'?'已断开连接':'已中止运行'" type="info"
                show-icon/>
      <el-alert style="margin-bottom: 10px" v-else-if="error==='client stopped'"
                title="服务器中断连接" type="info"
                show-icon/>
      <el-alert style="margin-bottom: 10px" v-else-if="error==='client restart'"
                title="连接重置中" type="info"
                show-icon/>
      <el-alert style="margin-bottom: 10px" v-else-if="error==='authentication failed'"
                title="验证失败" type="warning"
                show-icon/>
      <el-alert style="margin-bottom: 10px;text-align: left" v-else-if="error!==''" :title="error" type="error"
                show-icon></el-alert>
    </div>
    <div v-loading="loading"
         element-loading-text="处理中，请稍后..."
         class="box-outer"
         style="width: 100%;height:405px;padding-top:30px;padding-bottom: 8px;transition-duration: 0.5s">
      <div class="title" :style="online?'border-left: solid 8px #67C23A;':'border-left: solid 8px #909090'">
        <div class="title-text">{{ running ? online ? "已连接" : "连接中" : online ? "客户端异常" : "创建连接" }}</div>
      </div>
      <el-row>
        <el-col :span="20" :offset="2">
          <el-form
              label-position="top"
              label-width="100px"
              :model="user"
              :disabled="running"
              style="max-width: 460px"
          >
            <el-form-item label="用户">
              <el-input v-model="user.account"/>
            </el-form-item>
            <el-form-item label="密码">
              <el-input v-if="!running" v-model="user.password" show-password/>
              <el-input v-else v-model="displayPwd" show-password/>
            </el-form-item>
          </el-form>
          <el-button type="primary" v-if="!running" style="margin-top: 40px;width: 100%" @click="start">连接</el-button>
          <el-button type="danger" v-else style="margin-top: 40px;width: 100%" @click="stop">中止</el-button>
        </el-col>
      </el-row>
      <el-divider style="margin-top: 50px;margin-bottom: 10px"/>
      <div class="footer" style="text-align: center">
        <div style="margin-bottom: 5px">
          <el-link :underline="false" style="color: #007BBB;font-size: 13px"
                   @click="setting = true">
            连接设置
          </el-link>
        </div>
        Connect to Hub {{ '<' }}
        <span style="color: #007BBB" v-if="serverShotCut!==''">{{ serverShotCut }}</span>
        <span style="color: #E6A23C" v-else>未指定服务器</span>
        {{ '>' }}
      </div>
      <el-drawer v-model="setting" title="客户端设置" :direction="'rtl'" :append-to-body="true" :show-close="false"
                 custom-class="default-drawer">
        <template #header>
          <div class="title">
            <div class="title-text">客户端设置</div>
          </div>
        </template>
        <template #default>
          <el-row :gutter="10">
            <el-col :span="22" :offset="1">
              <el-form
                  label-position="top"
                  label-width="100px"
                  :model="server"
              >
                <el-form-item label="认证服务器地址">
                  <el-input v-model="server.address"/>
                </el-form-item>
                <el-form-item label="认证服务器端口">
                  <el-input v-model="server.port"/>
                </el-form-item>
                <el-form-item label="证书地址">
                  <el-input v-model="server.cert"/>
                </el-form-item>
              </el-form>
            </el-col>
          </el-row>
        </template>
        <template #footer>
          <div style="flex: auto">
            <el-button @click="setting=false">取消</el-button>
            <el-button type="primary" @click="saveConfig">确认</el-button>
          </div>
        </template>
      </el-drawer>
    </div>
    <!--    <div style="margin-top: 30px;display: none">-->
    <!--    </div>-->
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "LoginBox",
  components: {},
  data() {
    return {
      timer: undefined,
      displayPwd: "******",
      running: false,
      online: false,
      error: "",
      setting: false,
      loading: false,
      version: "",
      develop: false,
      serverShotCut: "",
      server: {
        address: "",
        port: "",
        cert: ""
      },
      user: {
        account: "",
        password: "",
      }
    }
  },
  mounted() {
    this.check()
    this.timer = setInterval(() => {
      this.check(true)
    }, 5000)
    this.loadConfig()
  },
  unmounted() {
    clearInterval(this.timer)
  },
  methods: {
    check: function (silence) {
      if (!silence) {
        this.loading = true
      }
      axios({
        method: "get",
        url: "/api/application/",
        data: {}
      }).then((res) => {
        let response = res.data
        if (this.error !== response.data.error ||
            this.running !== response.data.running ||
            this.online !== response.data.online) {
          this.$emit("updated")
        }
        this.error = response.data.error
        if (this.error !== "" && this.loading) {
          this.loading = false
        }
        this.running = response.data.running
        this.online = response.data.online
      }).catch((err) => {
        this.error = err.response.data.error
        if (this.error !== "" && this.loading) {
          this.loading = false
        }
      }).finally(() => {
        if (!silence) {
          this.loading = false
        }
      })
    },
    stop: function () {
      this.loading = true
      axios({
        method: "get",
        url: "/api/application/stop",
        data: {}
      }).then(() => {
      }).catch((err) => {
        this.$utils.HandleError(err)
      }).finally(() => {
        setTimeout(() => {
          this.check(false)
          this.$emit("updateCounter")
        }, 1500)
      })
    },
    start: function () {
      if (this.user.account === "") {
        this.$utils.Warning("提示", "用户名不能为空")
        return
      }
      if (this.user.password === "") {
        this.$utils.Warning("提示", "密码不能为空")
        return
      }
      this.loading = true
      axios({
        method: "post",
        url: "/api/application/start",
        data: {
          user: {
            account: this.user.account,
            password: this.user.password
          },
          auth: {
            address: this.server.address,
            port: Number(this.server.port)
          },
          security: {
            cert: this.server.cert
          }
        }
      }).then(() => {
      }).catch((err) => {
        this.loading = false
        this.$utils.HandleError(err)
      }).finally(() => {
        setTimeout(() => {
          this.check(false)
          this.$emit("updateCounter")
        }, 1500)
      })
    },
    loadConfig: function () {
      this.loading = true
      axios({
        method: "get",
        url: "/api/config",
        data: {}
      }).then(res => {
        let response = res.data
        this.server.address = response.data.auth.Address
        this.server.port = response.data.auth.Port
        this.server.cert = response.data.security.cert
        this.serverShotCut = this.server.address
        this.user.account = response.data.user.Account
        this.loading = false
      }).catch(() => {
        this.loading = false
        this.$utils.Error("获取版本失败", "请检查网络连接")
      })
    },
    saveConfig: function () {
      this.loading = true
      axios({
        method: "post",
        url: "/api/config/save",
        data: {
          user: {
            account: this.user.account
          },
          auth: {
            address: this.server.address,
            port: Number(this.server.port)
          },
          security: {
            cert: this.server.cert
          }
        }
      }).then(res => {
        let response = res.data
        this.$utils.Success("提示", response.msg)
        this.setting = false
        this.loading = false
        this.loadConfig()
      }).catch((err) => {
        this.loading = false
        this.$utils.HandleError(err)
      })
    },
  }
}
</script>

<style scoped>
.box-outer {
  box-shadow: 1px 1px 8px rgba(50, 50, 50, 0.2);
  border-radius: 4px;
}

.footer {
  font-size: 13px;
  color: #808080;
  text-align: right;
  padding-right: 10px;
  opacity: 0.8;
}
</style>