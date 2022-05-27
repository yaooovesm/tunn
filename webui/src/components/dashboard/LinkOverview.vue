<template>
  <div>
    <el-row :gutter="30">
      <el-col :span="24" style="font-size: 12px;color: #909399">
        <div style="text-align: center">
          TunnClient
          <span v-if="version===''" style="color: #E6A23C">未知版本</span>
          <span v-else>V{{ version }}</span>
          <span v-if="develop">[开发版本]</span>
        </div>
        <div style="text-align: center;margin-top: 8px">
          <el-link style="font-size: 12px;color: #007bbb;transform: translateY(-1px);">
            文档
          </el-link>
          |
          <el-link style="font-size: 12px;color: #007bbb;transform: translateY(-1px);">
            Gitee
          </el-link>
          |
          <el-link style="font-size: 12px;color: #007bbb;transform: translateY(-1px);">
            Github
          </el-link>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "LinkOverview",
  data() {
    return {
      version: "",
      develop: false
    }
  },
  mounted() {
    this.info()
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
      }).catch((err) => {
        this.$utils.HandleError(err)
      })
    },
  }
}
</script>

<style scoped>

</style>