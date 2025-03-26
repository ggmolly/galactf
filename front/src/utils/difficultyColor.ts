export function difficultyColor(difficulty: number) {
  switch (difficulty) {
    case 0:
      return "bg-muted-foreground text-muted";
    case 1:
      return "bg-accent-foreground text-accent";
    case 2:
      return "bg-secondary-foreground text-secondary";
    case 3:
      return "bg-primary-foreground text-primary";
    case 4:
      return "bg-indigo-950 text-indigo-300";
    case 5:
      return "bg-rose-950 text-rose-300";
    default:
      return "bg-stone-200";
  }
}
