<script lang="ts" setup>
import { computed, inject, ref, type Ref } from 'vue';
import { addIcons } from 'oh-vue-icons';
import {
  CoSortAlphaDown,
  CoSortAlphaUp,
  CoSortAscending,
  CoSortDescending,
  CoSettings,
  MdVisibilityTwotone,
  MdVisibilityoffTwotone,
} from 'oh-vue-icons/icons';
import { stringToCamelCase } from '../../lib/toCamelCase';
import type { Device, SortOrderConfig } from '../../types/types';
import {
  PREFERENCES_SYMBOL,
  type PreferencesStore,
} from '../../compositions/useUpdatePreference';

addIcons(
  CoSortAlphaDown,
  CoSortAlphaUp,
  CoSortAscending,
  CoSortDescending,
  CoSettings,
  MdVisibilityTwotone,
  MdVisibilityoffTwotone,
);

const pollingValues = [3, 5, 10];

const sortOrderConfig: SortOrderConfig = {
  nameAlphabeticalAscending: {
    title: 'Name: A - Z',
    icon: 'co-sort-alpha-down',
    value: 'name_alphabetical_ascending',
  },
  nameAlphabeticalDescending: {
    title: 'Name: Z - A',
    icon: 'co-sort-alpha-up',
    value: 'name_alphabetical_descending',
  },
  statusMostActiveFirst: {
    title: 'Status: Most active first',
    icon: 'co-sort-descending',
    value: 'status_most_active_first',
  },
  statusLeastActiveFirst: {
    title: 'Status: Least active first',
    icon: 'co-sort-ascending',
    value: 'status_least_active_first',
  },
};

const devices = inject<Ref<Device[]>>('devices');
const preferencesStore = inject<PreferencesStore>(PREFERENCES_SYMBOL);
if (!preferencesStore) {
  throw new Error('Preferences store not provided');
}
const preferences = computed(() => preferencesStore?.preferences.value);
const showSortMenu = ref<boolean>(false);
const showPreferencesMenu = ref<boolean>(false);

function setShowSortMenu() {
  showSortMenu.value = !showSortMenu.value;
}

function setShowPreferencesMenu() {
  showPreferencesMenu.value = !showPreferencesMenu.value;
}

function getSortConfig() {
  let key: string = 'nameAlphabeticalAscending';
  if (preferences.value.sortOrder) {
    key = stringToCamelCase(preferences.value?.sortOrder);
  }
  return sortOrderConfig[key];
}

function getHiddenDevicesCount() {
  if (preferences.value.hiddenDevices) {
    return `${preferences.value.hiddenDevices.length}/${devices?.value.length}`;
  }
  return 'n/a';
}
</script>

<template>
  <div class="preferences">
    <div class="preference-controls">
      <button
        id="preference-button"
        class="preference-button"
        @click="setShowSortMenu"
      >
        <v-icon
          class="preference-icon"
          v-if="preferences?.sortOrder"
          :name="getSortConfig().icon"
        />
        {{ getSortConfig().title }}
      </button>
      <button class="preference-button" @click="setShowPreferencesMenu">
        <v-icon class="preference-icon" name="co-settings" />Preferences
      </button>
      <div class="visibility-options-container">
        <span
          @click="
            async () => {
              await preferencesStore.updatePreferences({
                showVisibilityControls: !preferences?.showVisibilityControls,
              });
            }
          "
        >
          <v-icon
            :name="
              preferences.showVisibilityControls
                ? 'md-visibility-twotone'
                : 'md-visibilityoff-twotone'
            "
          />
        </span>
        <p>
          {{ getHiddenDevicesCount() }}
        </p>
      </div>
    </div>
    <ul id="sort-menu" class="menu-modal" v-if="showSortMenu">
      <li
        class="sort-menu"
        v-for="k in sortOrderConfig"
        @click="
          async () => {
            await preferencesStore.updatePreferences({
              sortOrder: k.value,
            });
            setShowSortMenu();
          }
        "
      >
        <v-icon class="sort-menu-icon" :name="k.icon" />
        <div>{{ k.title }}</div>
      </li>
    </ul>
    <div id="pref-memu" class="menu-modal" v-if="showPreferencesMenu">
      <fieldset class="pf-fieldset">
        <legend>Set update frequency (seconds):</legend>
        <div v-for="pv in pollingValues">
          <input
            type="radio"
            :id="`pf-${pv}`"
            name="polling-frequency"
            :value="pv"
            :checked="pv === Number(preferences.pollingFrequency)"
            @change="
              async () => await preferencesStore.updatePollingAndRestart(pv)
            "
          />
          <label :for="`pf-${pv}`">{{ pv }}</label>
        </div>
      </fieldset>
    </div>
  </div>
</template>

<style>
.preferences {
  width: 100%;
  height: fit-content;
  padding: 0 0 20px 0;
}

.preference-controls {
  display: flex;
  flex-direction: row;
  justify-content: space-evenly;
  align-items: center;
}

.preference-button {
  background-color: white;
  color: black;
  border-radius: 5px;
  border: 1px solid gray;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-evenly;
  width: 40%;
  padding: 0.6em 0;
  position: relative;
  font-size: 0.8em;
  font-weight: 500;
  font-family: inherit;
  cursor: pointer;
  height: 40px;
  transition: border-color 0.25s;
}

.preference-button:hover {
  background-color: whitesmoke;
}

.preference-icon {
  color: black;
}

.menu-modal {
  width: fit-content;
  position: absolute;
  background-color: white;
  list-style-type: none;
  padding: 5px;
  gap: 0;
  padding: 10px 10px;
  font-size: 0.8em;
  font-weight: 500;

  -webkit-box-shadow: 0px 0px 9px 3px rgba(41, 41, 41, 0.25);
  -moz-box-shadow: 0px 0px 9px 3px rgba(41, 41, 41, 0.25);
  box-shadow: 0px 0px 9px 3px rgba(41, 41, 41, 0.25);
}

.sort-menu {
  display: flex;
  flex-direction: row;
  align-items: center;
}

.sort-menu:hover {
  background-color: whitesmoke;
  cursor: pointer;
}

.sort-menu-icon {
  width: 2.25em;
  height: 2.25em;
  margin-right: 15px;
}

.visibility-options-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  font-size: 0.75em;
  width: fit-content;
  height: 25px;
}

.pf-fieldset {
  display: flex;
  flex-direction: row;
  justify-content: space-around;
}

.pf-fieldset legend {
  margin-bottom: 15px;
}

:root {
  --form-control-color: rgb(38, 46, 66);
}

*,
*:before,
*:after {
  box-sizing: border-box;
}

input[type='radio'] {
  -webkit-appearance: none;
  appearance: none;
  background-color: var(--form-background);
  cursor: pointer;
  margin: 0;

  font: inherit;
  color: currentColor;
  width: 1.15em;
  height: 1.15em;
  border: 0.15em solid currentColor;
  border-radius: 50%;
  transform: translateY(-0.075em);

  display: grid;
  place-content: center;
}

input[type='radio']::before {
  content: '';
  width: 0.65em;
  height: 0.65em;
  border-radius: 50%;
  transform: scale(0);
  transition: 120ms transform ease-in-out;
  box-shadow: inset 1em 1em var(--form-control-color);
  background-color: CanvasText;
}

input[type='radio']:checked::before {
  transform: scale(1);
}

input[type='radio']:focus {
  outline: max(2px, 0.15em) solid currentColor;
  outline-offset: max(2px, 0.15em);
}
</style>
