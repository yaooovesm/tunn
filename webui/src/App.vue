<template>
  <div id="app">
    <router-view/>
  </div>
</template>

<script>

export default {
  name: 'App',
  components: {},
  mounted() {
    window.onload = function () {
      if (!window.sessionStorage["login"]) {
        localStorage.removeItem("tunn")
        location.reload()
      } else {
        window.sessionStorage.removeItem("login");
      }
    };
    window.onunload = function () {
      window.sessionStorage["login"] = true;
    };
    window.onbeforeunload = function () {
      window.sessionStorage["login"] = true;
    };
  }
}
</script>

<style>
.el-form-item__label {
  height: 22px !important;
}

.icon {
  width: 1em;
  height: 1em;
  vertical-align: -0.15em;
  fill: currentColor;
  overflow: hidden;
}

.default-dialog {
  border-radius: 4px !important;
}

.default-dialog .el-dialog__header {
  padding-top: 20px;
  padding-left: 0;
  height: 50px;

}

.title {
  text-align: left;
  border-left: solid 8px rgba(0, 123, 187, 0.8);
  padding-left: 20px;
  height: 23px;
  line-height: 25px;
  margin-bottom: 50px;
}

.title-text {
  font-size: 18px;
}

.circle {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  display: inline-block;
  margin-bottom: 1px;
}

.online-text {
  color: #67C23A;
  height: 25px;
  margin-left: 3px;
  display: inline-block;
  font-size: 12px;
  transform: translateY(-1px);
}

.offline-text {
  color: #909399;
  height: 25px;
  margin-left: 3px;
  display: inline-block;
  font-size: 12px;
  transform: translateY(-1px);
}

.online-circle {
  background-color: #67C23A;
  box-shadow: 0 0 3px #67C23A;
}

.offline-circle {
  background-color: #909399;
  box-shadow: 0 0 3px #909399;
}

.default-drawer .el-drawer__header {
  padding-left: 0;
  color: #2c3e50;
}

body {
  -webkit-tap-highlight-color: transparent;
  -webkit-touch-callout: none;
  -webkit-user-select: none;
  padding: 0;
  margin: 0;
  height: 100%;
  width: 100%;
  position: absolute;
  user-select: none;
}

#app {
  height: 100%;
  width: 100%;
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
