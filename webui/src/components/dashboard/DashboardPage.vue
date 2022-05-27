<template>
  <div style="height: 100vh">
    <el-container style="height: 100vh">
      <el-header style="padding: 0">
        <div class="header-box">
          <span>Tunn Client</span>
          <el-button link @click="logout" style="color:#007BBB;margin-top: 15px;float: right;margin-right: 10px">退出
          </el-button>
        </div>
      </el-header>
      <el-main style="padding: 50px 20vh;">
        <div>
          <el-row :gutter="30">
            <el-col :xs="24" :sm="24" :md="8" :lg="6" :xl="6">
              <el-row :gutter="10">
                <el-col :span="24">
                  <login-box ref="login_box" @updated="updateOverview"/>
                </el-col>
                <el-col :span="24" style="margin-top: 30px">
                  <!--                  <link-overview/>-->
                </el-col>
              </el-row>
            </el-col>
            <el-col :xs="24" :sm="24" :md="16" :lg="18" :xl="18">
              <el-row :gutter="10">
                <el-col :span="24">
                  <flow-overview :passive="false"/>
                </el-col>
                <el-col :span="24" style="margin-top: 30px">
                  <status-overview ref="overview"/>
                </el-col>
              </el-row>
            </el-col>
          </el-row>
        </div>
      </el-main>
      <el-footer>
        <link-overview/>
      </el-footer>
    </el-container>
  </div>
</template>

<script>
import StatusOverview from "@/components/dashboard/StatusOverview";
// import ClientControl from "@/components/dashboard/ClientControl";
// import ConfigOverview from "@/components/dashboard/ConfigOverview";
import LoginBox from "@/components/login/LoginBox";
import axios from "axios";
import FlowOverview from "@/components/dashboard/FlowOverview";
import LinkOverview from "@/components/dashboard/LinkOverview";
//import LinkOverview from "@/components/dashboard/LinkOverview";

export default {
  name: "DashboardPage",
  components: {LinkOverview, FlowOverview, LoginBox, StatusOverview},
  data() {
    return {
      timer: undefined,
      application: {}
    }
  },
  mounted() {
    // this.getAppInfo()
    // this.timer = setInterval(() => {
    //   this.getAppInfo()
    // }, 5000)
  },
  unmounted() {
    //clearInterval(this.timer)
  },
  methods: {
    updateOverview: function () {
      this.$refs.overview.update()
    },
    getAppInfo: function () {
      axios({
        method: "get",
        url: "/api/application/",
        data: {}
      }).then((res) => {
        let response = res.data
        this.application = response.data
      })
    },
    logout: function () {
      localStorage.removeItem("tunn")
      this.$router.push({path: "/login"})
    }
  }
}
</script>

<style scoped>
.header-box {
  border-bottom: solid 1px rgba(60, 60, 60, 0.1);
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  height: 60px
}

.header-box span {
  display: inline-block;
  float: left;
  margin-left: 30px;
  line-height: 60px;
  font-size: 20px;
}
</style>