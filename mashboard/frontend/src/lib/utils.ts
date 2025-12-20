export const getTimeAgo = (timestamp: string = "") => {
  let date = new Date(timestamp.replace(" ", "T") + "Z");

  const now = new Date();
  const diff = Math.floor((now.getTime() - date.getTime()) / 1000); // Difference in seconds
  const days = Math.floor(diff / (60 * 60 * 24));
  const hours = Math.floor(diff / (60 * 60));
  const minutes = Math.floor(diff / 60);

  if (days > 0) return `${days}d ago`;
  if (hours > 0) return `${hours}h ago`;
  return `${minutes}m ago`;
};
