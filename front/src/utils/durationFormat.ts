export function durationFormat(seconds: number) {
  if (seconds <= 0) {
    return "0s";
  }

  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secondsLeft = seconds % 60;

  let result = "";

  if (days > 0) {
    result += `${days}d `;
  }
  if (hours > 0 || days > 0) {
    result += `${hours}h `;
  }
  if (minutes > 0 || hours > 0 || days > 0) {
    result += `${minutes}m `;
  }

  result += `${secondsLeft}s`;

  return result.trim();
}
