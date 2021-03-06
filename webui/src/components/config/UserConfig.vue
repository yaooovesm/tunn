<template>
  <div>
    <el-dialog
        v-model="dialogVisible"
        width="50%"
        top="18vh"
        :close-on-click-modal="false"
        custom-class="default-dialog"
        destroy-on-close
        :append-to-body="true"
        draggable
    >
      <template #header>
        <div class="title">
          <div class="title-text">用户配置
            <el-tag
                type=""
                effect="dark"
                style="transform: translateY(-2px);margin-left: 10px;height: 25px"
            >
              用户&nbsp;{{ account }}
            </el-tag>
          </div>
        </div>
      </template>
      <div v-loading="loading">
        <el-card shadow="always" body-style="padding:0">
          <div class="title" style="margin-top: 20px;margin-bottom: 20px">
            <div class="title-text">网络导入</div>
          </div>
          <div style="padding: 20px">
            <div v-if="importRoutes.length>0">
              <span v-for="route in importRoutes" :key="route">
                <el-popover
                    placement="top-start"
                    :width="150"
                    trigger="hover"
                >
                  <template #reference>
                    <el-tag
                        closable
                        effect="dark"
                        type="info"
                        :disable-transitions="false"
                        @close="handleDeleteImportRoute(route)"
                        style="margin-right: 10px;margin-bottom: 5px"
                    >
                      {{ route.network }}
                    </el-tag>
                  </template>
                  <template #default>
                    <div class="detail-unit">
                      <span>名称 </span> {{ route.name === '' ? '未命名' : route.name }}
                    </div>
                    <div class="detail-unit">
                      <span>网络 </span> {{ route.network }}
                    </div>
                  </template>
                </el-popover>
              </span>

            </div>
            <div v-else>
              <span style="font-size: 12px;color: #909399">还没有导入网络</span>
            </div>
            <route-import-selector
                style="margin-top: 8px"
                :imported="importRoutes"
                :exported="exportRoutes"
                @submit="handleAddImport"
            />
          </div>
        </el-card>
        <el-card shadow="always" body-style="padding:0" style="margin-top: 20px">
          <div class="title" style="margin-top: 20px;margin-bottom: 20px">
            <div class="title-text">网络暴露</div>
          </div>
          <div style="padding: 20px">
            <div v-if="exportRoutes.length>0">
              <span v-for="route in exportRoutes"
                    :key="route">
                <el-popover
                    placement="top-start"
                    :width="150"
                    trigger="hover"
                >
                  <template #reference>
                    <el-tag
                        closable
                        type="info"
                        effect="dark"
                        :disable-transitions="false"
                        @click="$refs.edit_export.show(route)"
                        @close="handleDeleteExportRoute(route)"
                        style="margin-right: 10px;margin-bottom: 5px"
                    >
                      {{ route.network }}
                    </el-tag>
                  </template>
                  <template #default>
                    <div class="detail-unit">
                      <span>名称 </span> {{ route.name === '' ? '未命名' : route.name }}
                    </div>
                    <div class="detail-unit">
                      <span>网络 </span> {{ route.network }}
                    </div>
                  </template>
                </el-popover>
              </span>

            </div>
            <div v-else>
              <span style="font-size: 12px;color: #909399">还没有暴露网络</span>
            </div>
            <div>
              <route-export-dialog ref="export_dialog" @submit="handleAddExport"/>
              <el-button @click="$refs.export_dialog.show()" size="small" style="margin-top: 8px">添加
              </el-button>
            </div>
          </div>
        </el-card>
      </div>
      <template #footer>
        <el-button :loading="loading" style="margin-left: 10px" @click="reset" type="warning">重置</el-button>
        <el-button :loading="loading" type="primary" @click="save">保存</el-button>
        <el-button :loading="loading" @click="close">关闭</el-button>
      </template>
      <route-export-edit-dialog ref="edit_export" @submit="handleExportModify"/>
    </el-dialog>
  </div>
</template>

<script>
import axios from "axios";
import {ElMessageBox} from "element-plus";
import RouteImportSelector from "@/components/config/RouteImportSelector";
import RouteExportDialog from "@/components/config/RouteExportDialog";
import RouteExportEditDialog from "@/components/config/RouteExportEditDialog";

export default {
  name: "UserConfig",
  components: {RouteExportEditDialog, RouteExportDialog, RouteImportSelector},
  data() {
    return {
      loading: false,
      dialogVisible: false,
      addImportValue: "",
      addExportValue: "",
      configId: "",
      account: "",
      importRoutes: [],
      exportRoutes: [],
    }
  },
  methods: {
    handleExportModify: function (r) {
      for (let i = 0; i < this.exportRoutes.length; i++) {
        let origin = this.exportRoutes[i]
        if (origin.name === r.name) {
          this.exportRoutes[i] = r
          break
        }
      }
      this.$refs.edit_export.close()
    },
    reset: function () {
      this.loading = true
      ElMessageBox.confirm(
          '是否重置该配置文件',
          '警告',
          {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: 'warning',
          }
      ).then(() => {
        axios({
          method: "get",
          url: "/api/remote/route/reset/",
          data: {}
        }).then(res => {
          let response = res.data
          this.$utils.Success("操作成功", response.msg)
        }).catch((err) => {
          this.$utils.HandleError(err)
        }).finally(() => {
          this.loading = false
          this.load()
        })
      }).catch(() => {
      }).finally(() => {
        this.loading = false
      })
    },
    show: function (account) {
      if (account === "") {
        this.$utils.Error("配置丢失", "请联系管理员")
        return
      }
      this.account = account
      this.dialogVisible = true
      this.load()
    },
    close: function () {
      this.dialogVisible = false
      this.addExportValue = ""
      this.addImportValue = ""
    },
    save: function () {
      this.loading = true
      console.log([...this.importRoutes, ...this.exportRoutes])
      axios({
        method: "post",
        url: "/api/remote/route/save",
        data: [...this.importRoutes, ...this.exportRoutes]
      }).then(() => {
        this.$utils.Success("提示", "更新配置成功")
        this.load()
        this.loading = false
      }).catch((err) => {
        this.$utils.HandleError(err)
        this.loading = false
      })
    },
    load: function () {
      this.loading = true
      axios({
        method: "get",
        url: "/api/remote/config/",
        data: {}
      }).then(res => {
        let response = res.data
        let routes = response.data.routes
        let imports = []
        let exports = []
        for (let i in routes) {
          if (routes[i].option === 'import') {
            imports.push(routes[i])
          } else if (routes[i].option === 'export') {
            exports.push(routes[i])
          }
        }
        this.importRoutes = imports
        this.exportRoutes = exports
        this.loading = false
      }).catch(() => {
        ElMessageBox.alert('加载用户配置失败', '错误', {
          confirmButtonText: '确认',
          callback: () => {
            this.dialogVisible = false
          },
        })
        this.loading = false
      })
    },
    handleDeleteImportRoute: function (route) {
      for (let i in this.importRoutes) {
        if (this.importRoutes[i].network === route.network) {
          this.importRoutes.splice(i, 1)
        }
      }
    },
    handleAddImport: function (routes) {
      for (let i = 0; i < routes.length; i++) {
        this.importRoutes.push({
          network: routes[i].network,
          option: "import"
        })
      }
    },
    handleDeleteExportRoute: function (route) {
      for (let i in this.exportRoutes) {
        if (this.exportRoutes[i].network === route.network) {
          this.exportRoutes.splice(i, 1)
        }
      }
    },
    handleAddExport: function (route) {
      this.exportRoutes.push(route)
      this.$refs.export_dialog.close()
    }
  }
}
</script>

<style scoped>
.detail-unit {
  text-align: right;
  font-size: 12px;
  color: #007bbb;
}

.detail-unit span {
  color: #404040;
  float: left;
  display: inline-block;
}
</style>