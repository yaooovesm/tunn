<template>
  <div>
    <el-card shadow="always" body-style="padding:0">
      <div class="title" style="margin-top: 20px;margin-bottom: 20px;">
        <div class="title-text">连接属性
        </div>
      </div>
      <div style="padding: 20px">
        <el-row :gutter="30">
        </el-row>
      </div>
      <div style="font-size: 12px;color: #808080;text-align: right;padding: 5px 10px">
        更新于
        {{ $utils.FormatDate("YYYY/mm/dd HH:MM:SS", updateTime) }}&nbsp;
        <el-button type="text" style="font-size: 12px;height: 12px;line-height: 13px" @click="update">刷新
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
        console.log(response)
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