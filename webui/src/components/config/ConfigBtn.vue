<template>
  <div>
    <div v-if="account!==''">
      <el-button size="small" type="primary" @click="$refs.config.show(account)">网络配置</el-button>
      <user-config ref="config"/>
    </div>
    <div v-else>
      <el-button size="small" disabled>网络配置</el-button>
      <el-tooltip
          effect="dark"
          content='配置丢失请联系管理员'
          placement="bottom"
      >
        <i class="iconfont icon-exclamation-circle"
           style="color: #F56C6C;font-size: 12px;margin-left: 5px;line-height: 24px"></i>
      </el-tooltip>
    </div>
  </div>
</template>

<script>
import UserConfig from "@/components/config/UserConfig";

export default {
  name: "ConfigBtn",
  components: {UserConfig},
  props: {
    account: {
      type: String,
    }
  },
}
</script>

<style scoped>

</style>