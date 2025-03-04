import { ChallengeWithSolveRate } from "@/interfaces/challenge.interface";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { Skeleton } from "./ui/skeleton";
import { DifficultyBadge } from "./DifficultyBadge";
import { CategoryBadge } from "./CategoryBadge";

interface ChallengeCardProps {
  challenge: ChallengeWithSolveRate;
  selectChallenge: (challenge: ChallengeWithSolveRate) => void;
}

export function ChallengeCard({ challenge, selectChallenge }: ChallengeCardProps) {
    return (
        <Card
            className="w-72 hover:scale-105 transition-all duration-150 hover:shadow-xl cursor-pointer"
            onClick={() => {
                selectChallenge(challenge);
            }}
        >
            <CardHeader>
                <CardTitle>
                    <div className="flex gap-x-2 justify-between">
                        <span className="my-auto">{challenge.name}</span>
                        <DifficultyBadge difficulty={challenge.difficulty} />
                    </div>
                </CardTitle>
            </CardHeader>
            <CardContent>
                <div className="flex flex-col gap-2">
                    <span className="text-lg font-bold">
                        500 points
                    </span>
                    <span className="text-xs text-muted-foreground">
                        {(challenge.solve_rate * 100) | 0}% solved
                    </span>
                </div>

                <hr className="my-2" />

                <div className="flex flex-wrap gap-2 mt-2">
                    {challenge.categories.map((category, i) => (
                        <CategoryBadge key={i} category={category} />
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