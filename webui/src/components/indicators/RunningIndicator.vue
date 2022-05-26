<template>
  <div style="display: inline-block;height: 25px;line-height: 25px">
    <div :class="flag?'circle online-circle':'circle offline-circle'"></div>
    <div :class="flag?'online-text':'offline-text'">
      {{ flag ? runningText : stoppedText }}
    </div>
  </div>
</template>

<script>
export default {
  name: "RunningIndicator",
  props: {
    flag: {
      type: Boolean,
      default: false
    },
    runningText: {
      type: String,
      default: "运行中"
    },
    stoppedText: {
      type: String,
      default: "未运行"
    }
  }
}
</script>

<style scoped>
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
</style>