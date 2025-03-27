const possibleColors = [
  "bg-primary-foreground",
  "bg-secondary-foreground",
  "bg-accent-foreground",
  "text-primary",
  "text-secondary",
  "text-accent",
];

export function categoryColor(category: string) {
  return possibleColors[((category.charCodeAt(0) * 1.25) | 0) % (possibleColors.length / 2)];
}

export function categoryForegroundColor(category: string) {
  return possibleColors[
    (((category.charCodeAt(0) * 1.25) | 0) % (possibleColors.length / 2)) +
      possibleColors.length / 2
  ];
}
