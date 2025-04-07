<script lang="ts" setup>
import { Loader } from '@googlemaps/js-api-loader';
import {
  computed,
  createApp,
  h,
  inject,
  onMounted,
  ref,
  shallowRef,
  toRef,
  watch,
  type Reactive,
  type Ref,
} from 'vue';
import type { Device, Markers } from '../../types/types';
import { MAP_SYMBOL, type MapStore } from '../../compositions/useMap';
import {
  PREFERENCES_SYMBOL,
  type PreferencesStore,
} from '../../compositions/useUpdatePreference';
import InfoWindowBody from './infoWindow/InfoWindowBody.vue';
import InfoWindowHeader from './infoWindow/InfoWindowHeader.vue';
import {
  FOLLOWING_SYMBOL,
  type FollowingDeviceStore,
} from '../../compositions/useFollowingDevice';
import MapUserIcon from './MapUserIcon.vue';
import type { Geocode } from '../../classes/Geocode';

const MAP_ZOOM = 12; // dynamically calculate based on average of all coords
const MAP_CENTER = {
  lat: 34.0549,
  lng: -118.2426,
};

const geocode = inject<Reactive<Geocode>>('geocode');
const mapStore = inject<MapStore>(MAP_SYMBOL);
const preferencesStore = inject<PreferencesStore>(PREFERENCES_SYMBOL);
const followingDeviceStore = inject<FollowingDeviceStore>(FOLLOWING_SYMBOL);

if (!mapStore) {
  throw new Error('mapStore not provided');
} else if (!preferencesStore) {
  throw new Error('Preferences store not provided');
} else if (!followingDeviceStore) {
  throw new Error('Following store not provided');
}

const preferences = computed(() => preferencesStore?.preferences.value);
const followingDeviceId = computed(
  () => followingDeviceStore.followingDeviceId.value,
);
const devices = inject<Ref<Device[]>>('devices');
const markers = ref<Markers>({});
const infoWindow = shallowRef<google.maps.InfoWindow>();
const infoWindowDeviceId = ref<string | undefined>();
const emit = defineEmits(['mounted']);

watch(
  () => [devices?.value, preferences.value.hiddenDevices],
  () => {
    if (Object.keys(markers.value).length) {
      updateMarkers();
    } else if (devices?.value.length) {
      createMarkers();
      updateMarkers();
    }
  },
  {
    deep: true,
  },
);

watch(
  () => preferences.value.hiddenDevices,
  () => {
    if (
      preferences.value.hiddenDevices &&
      preferences.value.hiddenDevices.includes(infoWindowDeviceId.value ?? '')
    ) {
      infoWindow.value?.close();
    }
  },
);

watch(
  () => followingDeviceId.value,
  () => {
    if (followingDeviceId?.value) {
      infoWindowDeviceId.value = followingDeviceId.value;
    }
  },
);

watch(
  () => [devices?.value, infoWindowDeviceId.value],
  () => {
    if (infoWindowDeviceId.value) {
      const device = devices?.value.find(
        (d) => d.deviceId === infoWindowDeviceId.value,
      );
      infoWindow?.value?.open({
        anchor: markers.value[infoWindowDeviceId.value],
        map: mapStore?.map.value,
      });
      setInfoWindow(device);
    }
  },
  {
    deep: true,
  },
);

onMounted(async () => {
  await initializeMap();
  emit('mounted');
});

async function initializeMap() {
  try {
    const loader = new Loader({
      apiKey: import.meta.env.VITE_GOOGLE_MAPS_JS_API_KEY,
      version: 'weekly',
    });

    await loader.importLibrary('maps');
    await loader.importLibrary('marker');

    const mapDiv: HTMLElement | null = document.getElementById('map');
    if (!mapDiv) {
      console.error('Map div not found');
      return;
    }

    if (!google.maps.marker) {
      console.error('Google Maps Marker library not loaded');
      return;
    }

    mapStore?.setMap(
      new google.maps.Map(mapDiv, {
        center: MAP_CENTER,
        zoom: MAP_ZOOM,
        mapId: 'map',

        fullscreenControl: false,
        streetViewControl: false,
        mapTypeControlOptions: {
          style: google.maps.MapTypeControlStyle.DEFAULT,
          mapTypeIds: ['roadmap', 'terrain'],
          position: google.maps.ControlPosition.TOP_RIGHT,
        },
      }),
    );

    infoWindow.value = new google.maps.InfoWindow({});

    if (mapStore?.map.value) {
      google.maps.event.addListener(mapStore.map.value, 'click', () => {
        infoWindow.value?.close();
        infoWindowDeviceId.value = undefined;
      });

      // Drag map to stop following device
      google.maps.event.addListener(mapStore.map.value, 'drag', () =>
        followingDeviceStore?.setFollowingDeviceId(undefined),
      );
    }
  } catch (error) {
    console.error('Error initializing map:', error);
  }
}

function createMarkers() {
  if (!mapStore?.map.value) {
    return;
  }

  // clean-up old pointers
  for (const k in markers.value) {
    markers.value[k].map = null;
  }
  markers.value = {};

  // ensure devices exist and have a value
  if (!devices?.value || devices.value.length === 0) {
    return;
  }

  const newMarkers: Markers = {};
  devices?.value.forEach((device) => {
    if (
      preferences.value.hiddenDevices &&
      !preferences.value.hiddenDevices.includes(device.deviceId)
    )
      try {
        const marker = new google.maps.marker.AdvancedMarkerElement({
          map: mapStore?.map.value,
          position: {
            lat: device.latestDevicePoint.lat,
            lng: device.latestDevicePoint.lng,
          },
          title: device.displayName,
          gmpClickable: true,
          content: null,
        });

        // showw info window for clicked marker
        marker.addEventListener('gmp-click', (e) => {
          e.stopPropagation();
          infoWindowDeviceId.value = device.deviceId;
          followingDeviceStore?.setFollowingDeviceId(device.deviceId);
        });

        newMarkers[device.deviceId] = marker;
      } catch (error) {
        console.error(
          'Error creating marker for device:',
          device.deviceId,
          error,
        );
      }
  });
  markers.value = newMarkers;
}

function updateMarkers(): void {
  devices?.value.forEach((device) => {
    try {
      const { lat, lng } = device.latestDevicePoint;
      const marker = markers.value[device.deviceId];
      if (!preferences.value.hiddenDevices.includes(device.deviceId)) {
        marker.map = mapStore?.map.value;
        marker.position = {
          lat: lat,
          lng: lng,
        };

        const glyph = document.createElement('div');
        const iconInstance = createApp({
          render: () => h(MapUserIcon, { device: toRef(device) }),
        });
        iconInstance.mount(glyph);
        marker.content = glyph;
      } else {
        marker.map = null;
      }
    } catch (error) {
      // console.error(`Could not update marker for ${device.deviceId}:`, error);
    }
  });
}

function setInfoWindow(device: Device | undefined): void {
  const contentBody = document.createElement('div');
  const contentHeader = document.createElement('div');
  const iWBody = createApp(InfoWindowBody, { device: toRef(device), geocode });
  const iWHeader = createApp(InfoWindowHeader, { device: toRef(device) });
  iWBody.mount(contentBody);
  iWHeader.mount(contentHeader);

  infoWindow.value?.setHeaderContent(contentHeader);
  infoWindow.value?.setContent(contentBody);
}
</script>

<template>
  <div class="map-container">
    <div id="map" ref="map"></div>
  </div>
</template>

<style lang="css" scoped>
.marker-overlay {
  position: absolute;
  background-color: white;
  border: 1px solid #ccc;
  border-radius: 8px;
  padding: 10px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  min-width: 200px;
  z-index: 1000;
}

.map-container {
  position: absolute;
  width: 100vw;
  height: 100vh;
}

#map {
  width: 100%;
  height: 100%;
}
</style>
