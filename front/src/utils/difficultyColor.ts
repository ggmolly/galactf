export function difficultyColor(difficulty: number) {
    switch (difficulty) {
        case 0:
            return 'bg-emerald-200 text-black'
        case 1:
            return 'bg-cyan-200 text-black'
        case 2:
            return 'bg-indigo-200 text-black'
        case 3:
            return 'bg-purple-200 text-black'
        case 4:
            return 'bg-pink-200 text-black'
        case 5:
            return 'bg-red-200 text-black'
        default:
            return 'bg-stone-200'
    }
}