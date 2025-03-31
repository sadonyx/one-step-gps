<script lang="ts" setup>
import { addIcons, OhVueIcon } from 'oh-vue-icons';
import {
  MdAssistantnavigation,
  MdPausecirclefilledRound,
  IoCloudOfflineSharp,
  MdStopcircleRound,
} from 'oh-vue-icons/icons';
import { getDriveStatusColor } from '../../lib/driveStatus';
import type { Device } from '../../types/types';
import type { Ref } from 'vue';

addIcons(
  MdAssistantnavigation,
  MdPausecirclefilledRound,
  IoCloudOfflineSharp,
  MdStopcircleRound,
);

type Props = { device: Ref<Device | undefined> };

const props: Props = defineProps<Props>();

function rotateIcon() {
  if (
    props.device.value?.online &&
    props.device.value.latestDevicePoint.deviceState.driveStatus === 'driving'
  ) {
    return props.device.value.latestDevicePoint.angle;
  } else {
    return 0;
  }
}
</script>

<template>
  <div
    class="user-icon"
    :style="{
      color: props.device.value && getDriveStatusColor(props.device.value),
      rotate: `${rotateIcon()}deg`,
    }"
  >
    <template v-if="props.device.value?.online">
      <OhVueIcon
        v-if="
          props.device.value?.latestDevicePoint.deviceState.driveStatus ===
          'driving'
        "
        name="md-assistantnavigation"
      />
      <OhVueIcon
        v-else-if="
          props.device.value?.latestDevicePoint.deviceState.driveStatus ===
          'idle'
        "
        name="md-pausecirclefilled-round"
      />
      <OhVueIcon
        v-else-if="
          props.device.value?.latestDevicePoint.deviceState.driveStatus ===
          'off'
        "
        name="md-stopcircle-round"
      />
    </template>
    <OhVueIcon v-else name="io-cloud-offline-sharp" />
  </div>
</template>

<style>
.user-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: white;
  cursor: pointer;
}

.user-icon svg {
  width: 30px;
  height: 30px;
}
</style>
