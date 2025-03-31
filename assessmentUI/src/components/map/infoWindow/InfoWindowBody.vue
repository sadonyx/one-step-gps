<script lang="ts" setup>
import { type Ref } from 'vue';
import { addIcons, OhVueIcon } from 'oh-vue-icons';
import { RiLoaderLine } from 'oh-vue-icons/icons';
import type { Device } from '../../../types/types';
import { getCardinalDirection } from '../../../lib/getCardinalDirection';
import { Geocode } from '../../../classes/Geocode';

addIcons(RiLoaderLine);

type Props = { device: Ref<Device>; geocode: Geocode };

const props: Props = defineProps<Props>();

function showLoading(): boolean {
  return (
    !props.geocode.places[props.device.value.deviceId] ||
    props.device.value.latestDevicePoint.deviceState.driveStatus === 'driving'
  );
}

function text() {
  if (!props.geocode?.places[props.device.value.deviceId]) {
    return Geocode.splitAtFirstComma('Loading,');
  } else {
    return Geocode.splitAtFirstComma(
      props.geocode.places[props.device.value.deviceId].formattedAddress,
    );
  }
}
</script>

<template>
  <div class="info-window">
    <div class="info-window-column">
      <div class="info-window-field">
        <div class="location-head-container">
          <h6>Location</h6>
          <OhVueIcon
            v-if="showLoading()"
            class="spinner location-loading"
            name="ri-loader-line"
          />
        </div>
        <p v-for="info in text()">{{ info }}</p>
      </div>
      <div class="info-window-field">
        <h6>Coordinates</h6>
        <p>
          {{ `Lat: ${props.device.value.latestDevicePoint.lat.toFixed(4)}` }}
        </p>
        <p>
          {{ `Lng: ${props.device.value.latestDevicePoint.lng.toFixed(4)}` }}
        </p>
      </div>
    </div>
    <div class="info-window-column">
      <div class="info-window-field">
        <h6>Speed</h6>
        <p>
          {{
            props.device.value.latestDevicePoint.devicePointDetail.speed.display
          }}
        </p>
      </div>
      <div class="info-window-field">
        <h6>Direction</h6>
        <p>
          {{ getCardinalDirection(props.device.value.latestDevicePoint.angle) }}
        </p>
      </div>
      <div class="info-window-field">
        <h6>Odometer</h6>
        <p>
          {{
            props.device.value.latestDevicePoint.deviceState.odometer.display
          }}
        </p>
      </div>
    </div>
  </div>
</template>

<style>
.info-window {
  display: flex;
  flex-direction: row;
  width: 300px;
  height: 160px;
  padding-top: 10px;
}

.info-window-column {
  display: flex;
  flex-direction: column;
  flex: 1 1 0;
}

.info-window-field {
  display: flex;
  flex-direction: column;
  padding-bottom: 15px;
}

.info-window-field h6 {
  font-weight: normal;
  padding-bottom: 5px;
}

.info-window-field p {
  font-weight: bold;
}

.location-head-container {
  display: flex;
  flex-direction: row;
  align-content: center;
}

.location-loading {
  width: 14px;
  height: 14px;
  margin-left: 5px;
  animation-duration: 1250ms !important;
}
</style>
