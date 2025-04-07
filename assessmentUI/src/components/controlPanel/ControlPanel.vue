<script lang="ts" setup>
import {
  computed,
  inject,
  onMounted,
  ref,
  unref,
  toRaw,
  watch,
  type Reactive,
  type Ref,
} from 'vue';
import UserPreferences from './UserPreferences.vue';
import type { Device } from '../../types/types';
import { addIcons } from 'oh-vue-icons';
import { MdGpsfixed, RiLoader5Line } from 'oh-vue-icons/icons';
import Address from './Address.vue';
import osgLogo from '../../assets/osg-logo.svg';
import { MAP_SYMBOL, type MapStore } from '../../compositions/useMap';
import {
  PREFERENCES_SYMBOL,
  type PreferencesStore,
} from '../../compositions/useUpdatePreference';
import {
  sortAlphaAscending,
  sortAlphaDescending,
  sortLeastActiveFirst,
  sortMostActiveFirst,
} from '../../lib/sortDevices';
import { getDriveStatus, getDriveStatusColor } from '../../lib/driveStatus';
import { Geocode } from '../../classes/Geocode';
import {
  FOLLOWING_SYMBOL,
  type FollowingDeviceStore,
} from '../../compositions/useFollowingDevice';

const devices = inject<Ref<Device[]>>('devices');
const mapStore = inject<MapStore>(MAP_SYMBOL);
const preferencesStore = inject<PreferencesStore>(PREFERENCES_SYMBOL);
const followingDeviceStore = inject<FollowingDeviceStore>(FOLLOWING_SYMBOL);
const geocode = inject<Reactive<Geocode>>('geocode');

if (!mapStore) {
  throw new Error('mapStore not provided');
} else if (!preferencesStore) {
  throw new Error('Preferences store not provided');
} else if (!followingDeviceStore) {
  throw new Error('Following store not provided');
}

const map = computed(() => mapStore?.map.value);
const preferences = computed(() => preferencesStore?.preferences.value);
const followingDeviceId = computed(
  () => followingDeviceStore.followingDeviceId.value,
);
const intervalId = ref<number | undefined>();

addIcons(MdGpsfixed, RiLoader5Line);

defineProps({
  loading: {
    type: Boolean,
    required: true,
    default: true,
  },
});

watch(
  () => followingDeviceId.value,
  (_, __, cleanUp) => {
    clearInterval(intervalId.value);

    if (followingDeviceId.value) {
      followSelectedDevice(followingDeviceId.value);

      const device = toRaw(devices)?.value.find(
        (d) => d.deviceId === followingDeviceId.value,
      );

      if (device) {
        // initial fetch
        geocode?.getLocation(device);

        // periodic fetch
        intervalId.value = setInterval(
          () => {
            const currentDevice = toRaw(devices)?.value.find(
              (d) => d.deviceId === followingDeviceId.value,
            );
            geocode?.getLocation(currentDevice || device);
          },
          (preferences.value.pollingFrequency ?? 5) * 1000,
        );
      }
    }

    cleanUp(() => clearInterval(intervalId.value));
  },
);

watch(
  () => devices?.value,
  () => {
    if (followingDeviceId.value) followSelectedDevice(followingDeviceId.value);
    sortDevices();
  },
);

watch(() => preferences.value, sortDevices);

onMounted(sortDevices);

function setFollowingDeviceId(deviceId: string): void {
  followingDeviceStore?.setFollowingDeviceId(deviceId);
}

function followSelectedDevice(deviceId: string): void {
  devices?.value.forEach((device) => {
    if (deviceId === device.deviceId) {
      const { lat, lng } = device.latestDevicePoint;
      centerSelectedDevice(lat, lng);
    }
  });
}

function centerSelectedDevice(lat: number, lng: number): void {
  try {
    const latlng = new google.maps.LatLng(lat, lng);
    map?.value?.setZoom(16);
    map?.value?.setCenter(latlng);
    map?.value?.fitBounds(map?.value.getBounds() as google.maps.LatLngBounds, {
      left: 150,
    });
  } catch (error) {
    console.error(error);
  }
}

function sortDevices() {
  let compareFn;
  switch (preferences.value?.sortOrder) {
    case 'name_alphabetical_ascending':
      compareFn = sortAlphaAscending;
      break;

    case 'name_alphabetical_descending':
      compareFn = sortAlphaDescending;
      break;

    case 'status_most_active_first':
      compareFn = sortMostActiveFirst;
      break;

    case 'status_least_active_first':
      compareFn = sortLeastActiveFirst;
  }
  devices?.value.sort(compareFn);
}

async function updateDeviceVisibility(deviceId: string): Promise<void> {
  const currentlyHidden = unref(preferences.value.hiddenDevices);
  const index = currentlyHidden.indexOf(deviceId);
  if (index === -1) {
    currentlyHidden.push(deviceId);
    await preferencesStore?.updatePreferences({
      hiddenDevices: currentlyHidden,
    });
  } else {
    currentlyHidden.splice(index, 1);
    await preferencesStore?.updatePreferences({
      hiddenDevices: currentlyHidden,
    });
  }
}
</script>

<template>
  <div class="control-panel">
    <div class="osg-branding">
      <img :src="osgLogo" />
      <p>User</p>
    </div>
    <div class="devices-header">
      <v-icon class="gps-icon" name="md-gpsfixed" />
      <h2>Devices</h2>
    </div>
    <UserPreferences />
    <ul class="device-list" v-if="!loading">
      <li
        v-for="device in devices"
        @click="() => setFollowingDeviceId(device.deviceId)"
        :key="device.deviceId"
        :style="{
          'border-left-color': getDriveStatusColor(device),
        }"
      >
        <div class="device-list-item-top-container">
          <span>{{ device.displayName }}</span>
          <span class="device-drive-status">{{ getDriveStatus(device) }}</span>
          <span
            class="eye-con"
            v-if="preferences.showVisibilityControls"
            @click="
              (e) => {
                e.stopPropagation();
                updateDeviceVisibility(device.deviceId);
              }
            "
          >
            <v-icon
              :name="
                preferences.hiddenDevices.includes(device.deviceId)
                  ? 'md-visibilityoff-twotone'
                  : 'md-visibility-twotone'
              "
            />
          </span>
        </div>
        <br />
        <div class="device-map-info">
          <Address
            :devices="devices"
            :place="geocode?.places[device.deviceId]"
          />
          <p
            class="following-device"
            v-if="unref(followingDeviceId) === device.deviceId"
          >
            Following
          </p>
        </div>
      </li>
    </ul>
    <div class="loading-devices" v-else>
      <v-icon class="spinner" name="ri-loader-5-line" />
      <p>Fetching devices...</p>
    </div>
  </div>
</template>

<style lang="css">
.spinner {
  animation-name: spin;
  animation-duration: 750ms;
  animation-iteration-count: infinite;
  animation-timing-function: linear;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.loading-devices {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: fit-content;
  height: fit-content;
  margin: 200px auto;
}

.control-panel {
  position: relative;
  color: black;
  width: 300px;
  height: 92.5%;
  background-color: white;
  margin-top: auto;
  margin-bottom: auto;
  margin-left: 50px;
  border-radius: 5px;
  border: 1.5px solid rgb(38, 46, 66);
  display: flex;
  flex-direction: column;
}

.devices-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  font-size: 1.5em;
  padding: 20px 0;
}

.gps-icon {
  margin: 0px 10px;
  width: 40px;
  height: 40px;
}

.device-list {
  list-style-type: none;
  padding: 0;
  line-height: normal;
  flex-grow: 1;
  overflow-y: auto;
}

.device-list li {
  display: flex;
  flex-direction: column;
  padding: 5px 10px;

  border-left: 5px solid;
  border-top: 1px solid gray;
  height: fit-content;
}

.device-list li:last-child {
  border-bottom: 1px solid gray;
}

.device-list li:hover {
  background-color: whitesmoke;
  cursor: pointer;
}

.device-list-item-top-container {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  height: 25px;
}

.device-drive-status {
  font-size: x-small;
}

.device-map-info {
  display: flex;
  flex-direction: row;
  align-items: end;
  justify-content: space-between;
  font-size: small;
  width: 100%;
}

.following-device {
  color: rgb(46, 139, 231);
  opacity: 1;
  animation: fadeinout 5s linear infinite;
  font-size: x-small;
}

@-webkit-keyframes fadeinout {
  0%,
  100% {
    opacity: 0;
  }
  65% {
    opacity: 1;
  }
  35% {
    opacity: 1;
  }
}

@keyframes fadeinout {
  0%,
  100% {
    opacity: 0;
  }
  65% {
    opacity: 1;
  }
  35% {
    opacity: 1;
  }
}

.osg-branding {
  height: 40px;
  width: 100%;
  background-color: rgb(38, 46, 66);
  border-top-left-radius: 3px;
  border-top-right-radius: 3px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
}

.osg-branding img {
  width: 50%;
  margin-left: 15px;
}

.osg-branding p {
  color: white;
  margin-right: 15px;
}
</style>
