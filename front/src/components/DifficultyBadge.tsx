import { cn } from "@/lib/utils";
import { difficultyColor } from "@/utils/difficultyColor";
import { Badge } from "./ui/badge";

export interface DifficultyBadgeProps {
  difficulty: number;
}

export function DifficultyBadge({ difficulty }: DifficultyBadgeProps) {
  return (
    <Badge className={cn("text-xs font-bold my-auto", difficultyColor(difficulty))}>
      {difficulty}/5
    </Badge>
  );
}
