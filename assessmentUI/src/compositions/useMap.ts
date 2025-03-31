import { shallowRef, type ShallowRef } from 'vue';

export const MAP_SYMBOL = Symbol('map');

export type MapStore = {
  map: ShallowRef<google.maps.Map | undefined>;
  setMap: (newMap: google.maps.Map) => void;
};

export function createMapStore(): MapStore {
  // `shallowRef` MUST be used for google Map objects, as `ref` wraps the map in a `Proxy`
  const map = shallowRef<google.maps.Map | undefined>();

  return {
    map,
    setMap(newMap: google.maps.Map) {
      map.value = newMap;
    },
  };
}
