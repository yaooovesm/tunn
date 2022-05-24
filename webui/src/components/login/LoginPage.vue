<template>
  <div>
    <el-container style="height: 100vh;background-color: #f2f2f2">
      <el-main>
        <el-row>
          <el-col :span="24">
            <!--            <login-box style="margin-top: 21vh;margin-left: 50%;transform: translateX(-50%);background-color: #ffffff"-->
            <!--                       :auto="true"/>-->
            <admin-login-box
                style="margin-top: 21vh;margin-left: 50%;transform: translateX(-50%);background-color: #ffffff"/>
          </el-col>
        </el-row>
      </el-main>
    </el-container>
  </div>
</template>

<script>
import AdminLoginBox from "@/components/login/AdminLoginBox";

export default {
  name: "LoginPage",
  components: {
    AdminLoginBox,
  },
  mounted() {
  },
  data() {
    return {}
  },
  methods: {}
}
</script>

<style scoped>

</style>