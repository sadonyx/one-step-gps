/**
 * Calculates how long ago a timestamp was relative to now.
 * @param timestamp - A timestamp in milliseconds
 * @returns An object containing the elapsed time in seconds, minutes, and hours
 */
function timeAgo(timestamp: number): {
  seconds: number;
  minutes: number;
  hours: number;
} {
  const now = Date.now();

  const diffMs = now - timestamp;

  if (diffMs < 0) {
    return { seconds: 0, minutes: 0, hours: 0 };
  }

  const seconds = Math.floor(diffMs / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);

  return {
    seconds,
    minutes,
    hours,
  };
}

/**
 * Returns a human-readable string describing how long ago a timestamp was.
 * @param timestamp - A timestamp in milliseconds (from Date.now())
 * @returns A formatted string like "2 hours ago" or "45 seconds ago"
 */
export function formatTimeAgo(timestamp: number): string {
  const { seconds, minutes, hours } = timeAgo(timestamp);

  const secondsDisplay = seconds % 60;
  const minutesDisplay = minutes % 60;

  let result = '';

  if (hours > 0) {
    result = `${hours} hour${hours !== 1 ? 's' : ''}`;
  } else if (minutesDisplay > 0) {
    if (result.length > 0) result += ', ';
    result = `${minutesDisplay} minute${minutesDisplay !== 1 ? 's' : ''}`;
  } else if (secondsDisplay > 0 || result.length === 0) {
    if (result.length > 0) result += ', ';
    result = `${secondsDisplay} second${secondsDisplay !== 1 ? 's' : ''}`;
  }

  return `${result} ago`;
}
