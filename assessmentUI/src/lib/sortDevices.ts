import type { Device } from '../types/types';

/**
 * A comparator function that sorts a list of devices by their `displayName` property in ascending-alphanumeric order.
 * @param {Device} a - The first device object.
 * @param {Device} b - The second device object.
 * @returns {number} A number -- `1`, `0`, or `-1` -- representing where the first device should be in the list relative to the second device.
 */
export const sortAlphaAscending = (a: Device, b: Device) =>
  a.displayName.localeCompare(b.displayName);

/**
 * A comparator function that sorts a list of devices by their `displayName` property in descending-alphanumeric order.
 * @param {Device} a - The first device.
 * @param {Device} b - The second device.
 * @returns {number} A number -- `1`, `0`, or `-1` -- representing where the first device should be in the list relative to the second device.
 */
export const sortAlphaDescending = (a: Device, b: Device) =>
  -1 * a.displayName.localeCompare(b.displayName);

/**
 * A comparator function that sorts a list of devices by their `online` and `driveStatus` in most-active-first order.
 * @param {Device} a - The first device.
 * @param {Device} b - The second device.
 * @returns {number} A number -- `1`, `0`, or `-1` -- representing where the first device should be in the list relative to the second device.
 */
export const sortMostActiveFirst = (a: Device, b: Device) => {
  if (a.online === false) return 1;
  if (b.online === false) return -1;
  const aState = a.latestDevicePoint.deviceState.driveStatus;
  const bState = b.latestDevicePoint.deviceState.driveStatus;
  const aDuration = a.latestDevicePoint.deviceState.driveStatusDuration.value;
  const bDuration = b.latestDevicePoint.deviceState.driveStatusDuration.value;
  return aState.localeCompare(bState) || aDuration - bDuration;
};

/**
 * A comparator function that sorts a list of devices by their `online` and `driveStatus` in least-active-first order.
 * @param {Device} a - The first device.
 * @param {Device} b - The second device.
 * @returns {number} A number -- `1`, `0`, or `-1` -- representing where the first device should be in the list relative to the second device.
 */
export const sortLeastActiveFirst = (a: Device, b: Device) => {
  if (a.online === false) return -1;
  if (b.online === false) return 1;
  const aState = a.latestDevicePoint.deviceState.driveStatus;
  const bState = b.latestDevicePoint.deviceState.driveStatus;
  const aDuration = a.latestDevicePoint.deviceState.driveStatusDuration.value;
  const bDuration = b.latestDevicePoint.deviceState.driveStatusDuration.value;
  return -1 * aState.localeCompare(bState) || bDuration - aDuration;
};
