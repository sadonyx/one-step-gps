<script lang="ts" setup>
import { type Ref } from 'vue';
import type { Device } from '../../../types/types';
import { getDriveStatus, getDriveStatusColor } from '../../../lib/driveStatus';

type Props = { device: Ref<Device> };

const props: Props = defineProps<Props>();
</script>

<template>
  <div class="info-window-header">
    <span class="info-window-header-column name-header">{{
      props.device.value.displayName
    }}</span>
    <span class="info-window-header-column">
      <span
        class="drive-status-header"
        :style="{ backgroundColor: getDriveStatusColor(props.device.value) }"
        >{{ getDriveStatus(props.device.value) }}</span
      >
    </span>
  </div>
</template>

<style>
.info-window-header {
  display: flex;
  flex-direction: row;
  align-items: center;

  padding-bottom: 15px;
  border-bottom: 1px solid whitesmoke;
}

.info-window-header-column {
  flex: 1 1 0;
}

.name-header {
  font-weight: bolder;
  font-size: large;
}

.drive-status-header {
  width: fit-content;
  border-radius: 5px;
  padding: 5px 10px;
  color: white;
  font-weight: bold;
  margin-left: -10px;
}

.gm-style-iw button {
  display: none !important;
}
</style>
