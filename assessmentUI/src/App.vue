<script setup lang="ts">
import GoogleMapLoader from './components/map/GoogleMap.vue';
import ControlPanel from './components/controlPanel/ControlPanel.vue';
import {
  onMounted,
  provide,
  reactive,
  ref,
  type Reactive,
  type Ref,
} from 'vue';
import type { Device } from './types/types';
import { createMapStore, MAP_SYMBOL } from './compositions/useMap';
import {
  createPreferencesStore,
  PREFERENCES_SYMBOL,
} from './compositions/useUpdatePreference';
import { Geocode } from './classes/Geocode';
import {
  createFollowingDeviceStore,
  FOLLOWING_SYMBOL,
} from './compositions/useFollowingDevice';

const devices = ref<Device[]>([]);
provide<Ref<Device[]>>('devices', devices);

const mapStore = createMapStore();
provide(MAP_SYMBOL, mapStore);

const preferencesStore = createPreferencesStore(devices);
provide(PREFERENCES_SYMBOL, preferencesStore);

const followingDeviceStore = createFollowingDeviceStore();
provide(FOLLOWING_SYMBOL, followingDeviceStore);

const geocode = reactive(new Geocode());
provide<Reactive<Geocode>>('geocode', geocode);

const showGoogleMap = ref(true);
const showControlPanel = ref(false);

function handleMapMounted() {
  showControlPanel.value = true;
}

onMounted(async () => {
  await preferencesStore.fetchPreferences();
  preferencesStore.establishEventSource(devices);
});
</script>

<template>
  <GoogleMapLoader v-if="showGoogleMap" @mounted="handleMapMounted" />
  <ControlPanel
    v-if="showControlPanel"
    :loading="preferencesStore.loading.value"
  />
</template>

<style>
@import '/stylesheets/whitespace-reset.css';
</style>
