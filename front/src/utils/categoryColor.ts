const possibleColors = [
    'bg-sky-400',
    'bg-indigo-400',
    'bg-fuchsia-400',
];

export function categoryColor(category: string) {
    return possibleColors[((category.charCodeAt(0) * 1.25) | 0) % possibleColors.length]
}