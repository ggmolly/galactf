import { cn } from "@/lib/utils";
import { categoryColor } from "@/utils/categoryColor";
import { Badge } from "./ui/badge";

export interface CategoryBadgeProps {
    category: string;
}

export function CategoryBadge({ category }: CategoryBadgeProps) {
    return (
        <Badge key={category} className={cn("text-xs", categoryColor(category))}>
            {category}
        </Badge>
    )
}