/**
 * Returns a formatted duration string from a given number of seconds, example:
 * 61 = 1 minute, 1 second
 * 37 = 37 seconds
 * 3600 = 1 hour
 * @param seconds 
 */
export function formatDuration(seconds: number) {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secondsLeft = seconds % 60;
    const parts = [];

    if (hours > 0) parts.push(`${hours}h`);
    if (minutes > 0 || (hours > 0 && secondsLeft > 0)) parts.push(`${minutes}m`);
    if (secondsLeft > 0 || (hours === 0 && minutes === 0)) parts.push(`${secondsLeft}s`);

    return parts.join(' ');
}