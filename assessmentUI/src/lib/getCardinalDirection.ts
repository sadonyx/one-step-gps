/**
 * Converts an angle in degrees to a cardinal direction.
 * @param {number} angle - The angle in degrees (0-360), where 0 is north, 90 is east, 180 is south, 270 is west
 * @returns {string} A string representing the cardinal direction
 */
export function getCardinalDirection(angle: number): string {
  const normalizedAngle = ((angle % 360) + 360) % 360;

  if (normalizedAngle >= 337.5 || normalizedAngle < 22.5) {
    return 'North';
  } else if (normalizedAngle >= 22.5 && normalizedAngle < 67.5) {
    return 'Northeast';
  } else if (normalizedAngle >= 67.5 && normalizedAngle < 112.5) {
    return 'East';
  } else if (normalizedAngle >= 112.5 && normalizedAngle < 157.5) {
    return 'Southeast';
  } else if (normalizedAngle >= 157.5 && normalizedAngle < 202.5) {
    return 'South';
  } else if (normalizedAngle >= 202.5 && normalizedAngle < 247.5) {
    return 'Southwest';
  } else if (normalizedAngle >= 247.5 && normalizedAngle < 292.5) {
    return 'West';
  } else {
    return 'Northwest';
  }
}
