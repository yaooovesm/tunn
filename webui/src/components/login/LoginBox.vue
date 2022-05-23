<template>
  <div v-loading="loading"
       class="box-outer"
       style="width: 420px;height:405px;padding-top:30px;padding-bottom: 8px;transition-duration: 0.5s">
    <div class="title">
      <div class="title-text">登录</div>
    </div>
    <el-row>
      <el-col :span="20" :offset="2">
        <el-form
            label-position="top"
            label-width="100px"
            :model="user"
            style="max-width: 460px"
        >
          <!--          <el-form-item label="服务器">-->
          <!--            <template #label>-->
          <!--              服务器-->
          <!--              <el-link :underline="false" style="color: #007BBB" @click="setting = true">-->
          <!--                {{ server.address === "" ? "设置" : "切换" }}-->
          <!--              </el-link>-->
          <!--            </template>-->
          <!--            <el-input v-model="serverShotCut" readonly style="cursor: pointer"/>-->
          <!--          </el-form-item>-->
          <el-form-item label="用户">
            <el-input v-model="user.account"/>
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="user.password" show-password/>
          </el-form-item>
        </el-form>
        <el-button type="primary" style="margin-top: 40px;width: 100%" @click="login">Login</el-button>
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
      Connect to Hub {{ '<' }} <span style="color: #007BBB">{{ serverShotCut }}</span> {{ '>' }}
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
              <el-form-item label="服务器地址">
                <el-input v-model="server.address"/>
              </el-form-item>
              <el-form-item label="服务器端口">
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
</template>

<script>
import axios from "axios";

export default {
  name: "LoginBox",
  data() {
    return {
      setting: false,
      loading: false,
      version: "",
      develop: false,
      serverShotCut: "未指定服务器",
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
    this.loadConfig()
  },
  methods: {
    login: function () {
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
            port: this.server.port
          },
          security: {
            cert: this.server.cert
          }
        }
      }).then(() => {
        this.loading = false
      }).catch((err) => {
        this.loading = false
        this.$utils.HandleError(err)
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
            port: this.server.port
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
  box-shadow: 1px 1px 5px rgba(50, 50, 50, 0.5);
  border-radius: 10px;
}

.footer {
  font-size: 13px;
  color: #808080;
  text-align: right;
  padding-right: 10px;
  opacity: 0.8;
}
</style>