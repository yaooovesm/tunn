<template>
  <div v-loading="loading"
       class="box-outer"
       style="width: 420px;height:400px;padding-top:30px;padding-bottom: 8px;transition-duration: 0.5s">
    <div class="title">
      <div class="title-text">登录</div>
    </div>
    <el-row>
      <el-col :span="20" :offset="2">
        <el-form
            label-position="top"
            label-width="100px"
            :model="loginData"
            style="max-width: 460px"
        >
          <el-form-item label="用户">
            <el-input v-model="loginData.user" @keyup.enter="login"/>
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="loginData.password" @keyup.enter="login" show-password/>
          </el-form-item>
        </el-form>
        <el-button type="primary" @click="login" style="margin-top: 40px;width: 100%">Login</el-button>
      </el-col>
    </el-row>
    <el-divider style="margin-top: 50px;margin-bottom: 30px"/>
    <div class="footer">
      Tunn v{{ version }}
      <span v-if="develop">[开发版本]
        <el-tooltip
            effect="dark"
            content="当前版本为开发版，可能存在缺陷，请勿使用"
            placement="bottom-end"
        >
        <i style="font-size: 8px;color: rgba(0,123,187,0.8)" class="iconfont icon-question-circle"></i>
        </el-tooltip>
      </span>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "AdminLoginBox",
  data() {
    return {
      loading: false,
      version: "",
      develop: false,
      loginData: {
        user: "",
        password: "",
      }
    }
  },
  async mounted() {
    await this.info()
    this.autoLogin()
  },
  methods: {
    info: function () {
      axios({
        method: "get",
        url: "/api/admin/",
        data: {}
      }).then((res) => {
        let response = res.data
        this.version = response.data.version
        this.develop = response.data.develop
        if (!response.data.require_login) {
          this.$router.push({path: "/dashboard"})
        } else {
          this.$router.push({path: "/login"})
        }
      }).catch((err) => {
        this.$utils.HandleError(err)
      })
    },
    login: function () {
      if (this.loginData.user === "") {
        this.$utils.Warning("提示", "用户名不能为空")
        return
      }
      if (this.loginData.password === "") {
        this.$utils.Warning("提示", "密码不能为空")
        return
      }
      this.loading = true
      axios({
        method: "post",
        url: "/api/admin/login",
        data: {
          user: this.loginData.user,
          password: this.loginData.password
        }
      }).then(() => {
        localStorage.setItem("tunn", new Date().getTime().toString())
        this.loading = false
        this.$router.push({path: "/dashboard"})
      }).catch((err) => {
        this.loading = false
        this.$utils.Error("登录失败", err.response.data.error)
      })
    },
    autoLogin: function () {
      let lo = localStorage.getItem("tunn")
      if (lo !== "" && lo !== undefined && lo !== null) {
        //检查是否过期
        let sto = Number(lo)
        let now = new Date().getTime()
        if ((now - sto) > 60 * 60 * 1000) {
          localStorage.removeItem("tunn")
        } else {
          this.$router.push({path: "/dashboard"})
        }
      }
    }
  }
}
</script>

<style scoped>
.box-outer {
  box-shadow: 1px 1px 5px rgba(50, 50, 50, 0.5);
  border-radius: 10px;
}

.footer {
  font-size: 11px;
  color: #808080;
  text-align: right;
  padding-right: 10px;
  opacity: 0.5;
}
</style>