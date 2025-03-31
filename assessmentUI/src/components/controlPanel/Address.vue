<script lang="tsx" setup>
import { watch } from 'vue';
import { Geocode, type Place } from '../../classes/Geocode';
import { formatTimeAgo } from '../../lib/timeAgo';
import type { Device } from '../../types/types';

type Props = {
  devices: Device[] | undefined;
  place: Place | undefined;
};

const props = defineProps<Props>();

function showAddress() {
  if (!props.place) {
    return Geocode.splitAtFirstComma('Click to follow, and view live address');
  } else {
    return Geocode.splitAtFirstComma(props.place.formattedAddress);
  }
}

watch(
  () => props.devices,
  () => {
    if (props.place) props.place.lastUpdated = props.place.lastUpdated;
  },
);
</script>

<template>
  <div class="address">
    <p v-for="info in showAddress()">
      {{ info }}
    </p>
    <!-- <p v-if="props.place.lastUpdated">
      {{ formatTimeAgo(props.place.lastUpdated) }}
    </p> -->
    <p v-if="props.place" class="last-updated">
      Last updated {{ formatTimeAgo(props.place.lastUpdated) }}
    </p>
  </div>
</template>

<style>
.last-updated {
  font-size: x-small;
  font-style: italic;
  color: gray;
}
</style>
