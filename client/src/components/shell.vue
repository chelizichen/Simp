<script lang="ts">
export default {
  name: "shell-componetn",
};
</script>
<template>
  <el-dialog
    v-model="props.shellVisible"
    title="SimpShell"
    width="80%"
    @close="emits('closeDialog')"
  >
    <iframe
      src="/simpshellserver/console"
      width="100%"
      style="min-height: 500px"
    ></iframe>
  </el-dialog>
</template>

<script setup lang="ts">
const props = defineProps<{
  shellVisible: boolean;
}>();
const emits = defineEmits(["closeDialog", "showReleaseDialog"]);
</script>

<style scoped>
.inp >>> .el-textarea__inner {
  background-color: black;
  border: none;
  box-shadow: none;
  color: white;
}
</style>
