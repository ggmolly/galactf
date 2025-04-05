import { cn } from "@/lib/utils";
import { difficultyColor } from "@/utils/difficultyColor";
import { Badge } from "./ui/badge";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "./ui/tooltip";

export interface DifficultyBadgeProps {
  difficulty: number;
}

const diffMap: Record<number, string> = {
  0: "Intro",
  1: "Very easy",
  2: "Easy",
  3: "Medium",
  4: "Hard",
  5: "Very hard",
  6: "Impossible",
};

const BONUS_DIFFICULTY = 255;

export function DifficultyBadge({ difficulty }: DifficultyBadgeProps) {
    console.log({ difficulty })
    if (difficulty === BONUS_DIFFICULTY) {
        return <Badge className={cn("text-xs font-bold my-auto", difficultyColor(difficulty))}>Bonus</Badge>
    }
  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger tabIndex={-1}>
          <Badge className={cn("text-xs font-bold my-auto", difficultyColor(difficulty))}>
            {difficulty}/5
          </Badge>
        </TooltipTrigger>
        <TooltipContent>
          <p>{diffMap[difficulty] ?? "Unknown"}</p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
}
