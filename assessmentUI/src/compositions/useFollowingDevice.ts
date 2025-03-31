import { ref, type Ref } from 'vue';

export const FOLLOWING_SYMBOL = Symbol('following');

export type FollowingDeviceStore = {
  followingDeviceId: Ref<string | undefined>;
  setFollowingDeviceId: (newFollowingDeviceId: string | undefined) => void;
};

export function createFollowingDeviceStore(): FollowingDeviceStore {
  // `shallowRef` MUST be used for google Map objects, as `ref` wraps the map in a `Proxy`
  const followingDeviceId = ref<string | undefined>();

  return {
    followingDeviceId,
    setFollowingDeviceId(newFollowingDeviceId: string | undefined) {
      followingDeviceId.value = newFollowingDeviceId;
    },
  };
}
