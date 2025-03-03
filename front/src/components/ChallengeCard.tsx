import { Challenge } from "@/interfaces/challenge.interface";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { Skeleton } from "./ui/skeleton";
import { Badge } from "./ui/badge";
import { categoryColor } from "@/utils/categoryColor";
import { cn } from "@/lib/utils";
import { difficultyColor } from "@/utils/difficultyColor";

interface ChallengeCardProps {
  challenge: Challenge;
}

export function ChallengeCard({ challenge }: ChallengeCardProps) {
    return (
        <Card className="w-72">
            <CardHeader>
                <CardTitle>
                    <div className="flex gap-x-2 justify-between">
                        <span className="my-auto">{challenge.name}</span>
                        <Badge className={cn("text-xs", "font-bold", "my-auto", difficultyColor(challenge.difficulty))}>
                            {challenge.difficulty}/5
                        </Badge>
                    </div>
                </CardTitle>
            </CardHeader>
            <CardContent>
                <div className="flex flex-col gap-2">
                    <span className="text-lg font-bold">
                        500 points
                    </span>
                    <span className="text-xs text-muted-foreground">
                        {(Math.random() * 100) | 0}% solved
                    </span>
                </div>

                <hr className="my-2" />

                <div className="flex flex-wrap gap-2 mt-2">
                    {challenge.categories.map((category) => (
                        <Badge key={category} className={cn("text-xs", categoryColor(category))}>
                            {category}
                        </Badge>
                    ))}
                </div>
            </CardContent>
        </Card>
    )
}

interface ChallengeCardSkeletonProps {
    index: number;
}

export function ChallengeCardSkeleton({ index }: ChallengeCardSkeletonProps) {
    return (
        <Card className="w-72">
            <CardHeader>
                <CardTitle>
                    <div className="flex gap-x-2 justify-between">
                        <Skeleton className="w-2/3 h-4 my-auto" />
                        <Skeleton className="w-12 h-6 my-auto" />
                    </div>
                </CardTitle>
            </CardHeader>
            <CardContent>
                <div className="flex flex-col gap-2">
                    <Skeleton className="w-3/5 h-8" />
                    <Skeleton className="w-1/3 h-4" />
                </div>

                <hr className="my-2" />

                <div className="flex flex-wrap gap-2 mt-2">
                    {Array((index % 3) + 1).fill(null).map((_, i) => (
                        <Skeleton key={i} className="w-16 h-4" />
                    ))}
                </div>
            </CardContent>
        </Card>
    )
}

export default ChallengeCard;