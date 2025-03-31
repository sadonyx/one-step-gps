import type { Device } from '../types/types';

/**
 * The respective colors of all possible `driveStatus` values.
 */
export const COLORS = {
  driving: '#4CAF50',
  idle: '#FF9800',
  off: '#F44336',
  offline: '#000000',
};

/**
 * Returns a human-readable string describing the drive status of a specific device.
 * @param {Device} device - The device object
 * @returns {string} A formatted string like "Driving 3m 42s"
 */
export function getDriveStatus(device: Device) {
  const isOnline = device.online;
  const status = device.latestDevicePoint.deviceState.driveStatus;
  const duration =
    device.latestDevicePoint.deviceState.driveStatusDuration.display;
  if (isOnline) {
    switch (status) {
      case 'driving':
        return 'Driving ' + duration;
      case 'idle':
        return 'Engine idle ' + duration;
      case 'off':
        return 'Engine off ' + duration;
    }
  } else return 'Offline ' + duration;
}

/**
 * Returns a color visually describing drive status of a specific device.
 * @param {Device} device - The device object
 * @returns {string} A hex color code like "#000000"
 */
export function getDriveStatusColor(device: Device) {
  const status = device.latestDevicePoint.deviceState.driveStatus;
  return device.online ? COLORS[status] : COLORS.offline;
}
