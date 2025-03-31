import '@googlemaps/js-api-loader';

type ValueUnitDisplay = {
  value: number;
  display: string;
  unit: string;
};

type DevicePointDetail = {
  speed: ValueUnitDisplay;
};

type DeviceState = {
  driveStatus: 'driving' | 'idle' | 'off';
  driveStatusDuration: ValueUnitDisplay;
  odometer: ValueUnitDisplay;
};

export type LatestDevicePoint = {
  deviceDetail: any;
  angle: number;
  devicePointDetail: DevicePointDetail;
  deviceState: DeviceState;
  formattedAddress: string;
  lat: number;
  lng: number;
};

export type Device = {
  deviceId: string;
  displayName: string;
  latestDevicePoint: LatestDevicePoint;
  make: string;
  model: string;
  online: boolean;
};

export type Coordinates = {
  deviceId: string;
  lat: number;
  lng: number;
};

export type Markers = {
  [key: string]: google.maps.marker.AdvancedMarkerElement;
};

export type SortOrder =
  | 'name_alphabetical_ascending'
  | 'name_alphabetical_descending'
  | 'status_most_active_first'
  | 'status_least_active_first';

export type UserPreferences = {
  sortOrder: SortOrder;
  hiddenDevices: string[];
  visits: number;
  showVisibilityControls: boolean;
  pollingFrequency: number;
};

type SortOrderConfigOptions = {
  title: string;
  icon: string;
  value: SortOrder;
};

export type SortOrderConfig = {
  [key: string]: SortOrderConfigOptions;
};
