import { ref, type Ref } from 'vue';
import { type Device, type UserPreferences } from '../types/types';
import { keysToCamelCase } from '../lib/toCamelCase';

export const PREFERENCES_SYMBOL = Symbol('preferences');

export type PreferencesStore = {
  loading: Ref<boolean>;
  preferences: Ref<UserPreferences>;
  setPreferences: (newPreferences: UserPreferences) => void;
  fetchPreferences: () => Promise<void>;
  updatePreferences: (updated: Partial<UserPreferences>) => Promise<void>;
  updatePollingAndRestart: (pollingValue: number) => Promise<void>;
  establishEventSource: (devices: Ref<Device[]>) => void;
};

export function createPreferencesStore(
  devices: Ref<Device[]>,
): PreferencesStore {
  const loading = ref(true);
  const shortLivedToken = ref<string>();
  const eventSource = ref<EventSource>();
  const preferences = ref<UserPreferences>({
    sortOrder: 'name_alphabetical_ascending',
    hiddenDevices: [],
    visits: 0,
    showVisibilityControls: false,
    pollingFrequency: 5, // seconds
  });

  function setPreferences(newPreferences: UserPreferences) {
    preferences.value = newPreferences;
  }

  async function fetchPreferences() {
    try {
      const response = await fetch('http://localhost:8080/user-preferences', {
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      shortLivedToken.value =
        response.headers.get('Authorization') ?? undefined;
      setPreferences(keysToCamelCase(await response.json()));
    } catch (error) {
      console.error('Failed to get user preferences:', error);
    }
  }

  async function updatePreferences(updated: Partial<UserPreferences>) {
    const oldPreferences = preferences.value;
    try {
      const response = await fetch('http://localhost:8080/user-preferences', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },

        body: JSON.stringify({ ...oldPreferences, ...updated }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      preferences.value = keysToCamelCase(data);
    } catch (error) {
      console.error('Failed to update preferences:', error);
    }
  }

  async function updatePollingAndRestart(pollingValue: number) {
    await updatePreferences({ pollingFrequency: pollingValue });
    eventSource.value?.close();
    if (eventSource.value?.CLOSED) {
      await fetchPreferences();
      establishEventSource(devices);
    }
  }

  function establishEventSource(devices: Ref<Device[]>) {
    eventSource.value = new EventSource(
      `http://127.0.0.1:8080/?slt=${shortLivedToken.value}`,
    );

    eventSource.value.onmessage = (event) => {
      try {
        const response = JSON.parse(event.data);
        if (response['result_list'] && response['result_list'].length > 0) {
          devices.value = response['result_list'].map((device: any) =>
            keysToCamelCase(device),
          );
        } else console.log('No results.');
        loading.value = false;
      } catch (error) {
        console.error('Error retreiving or parsing data:', error);
      }
    };
  }

  return {
    loading,
    preferences,
    setPreferences,
    fetchPreferences,
    updatePreferences,
    updatePollingAndRestart,
    establishEventSource,
  };
}
