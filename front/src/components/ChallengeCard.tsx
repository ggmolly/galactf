import { ChallengeWithSolveRate } from "@/interfaces/challenge.interface";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { Skeleton } from "./ui/skeleton";
import { DifficultyBadge } from "./DifficultyBadge";
import { CategoryBadge } from "./CategoryBadge";
import { computePoints } from "@/utils/computePoints";
import { cn } from "@/lib/utils";
import { useEffect, useState } from "react";
import { durationFormat } from "@/utils/durationFormat";
import { LockIcon } from "lucide-react";

interface ChallengeCardProps {
  challenge: ChallengeWithSolveRate;
  selectChallenge: (challenge: ChallengeWithSolveRate) => void;
}

export function ChallengeCard({ challenge, selectChallenge }: ChallengeCardProps) {
  if (!!challenge.reveal_in) {
    return <ChallengeCardLocked revealIn={challenge.reveal_in} />;
  }

  return (
    <Card
      className={cn(
        "w-72 hover:scale-105 transition-all duration-150 cursor-pointer",
        challenge.solved ? "border-secondary" : null
      )}
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
          <span className="text-lg font-bold">{computePoints(challenge.solvers ?? 0)} points</span>
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
  );
}

interface ChallengeCardLockedProps {
  revealIn: number;
}

export function ChallengeCardLocked({ revealIn }: ChallengeCardLockedProps) {
  const [secondsRemaining, setSecondsRemaining] = useState<number>(revealIn);

  useEffect(() => {
    const timer = setInterval(() => {
      setSecondsRemaining((prevSeconds) => {
        if (prevSeconds <= 1) {
          clearInterval(timer);
          return 0;
        }
        return prevSeconds - 1;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, []);

  return (
    <Card className="w-72 hover:scale-105 transition-all duration-150 hover:shadow-xl cursor-not-allowed border-primary">
      <CardHeader>
        <CardTitle>
          <div className="flex gap-x-2 justify-between">
            <div className="flex gap-x-2">
              <LockIcon className="w-4 h-4 text-primary" />
              <span className="text-primary">Locked</span>
            </div>
            <span className="text-xs text-secondary">{durationFormat(secondsRemaining)} left</span>
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex flex-col gap-2">
          <span className="text-lg font-bold text-muted-foreground">{computePoints(0)} points</span>
          <span className="text-xs text-muted-foreground">Not available yet</span>
        </div>

        <hr className="my-2" />

        <div className="flex flex-wrap gap-2 mt-2">
          <span className="text-xs text-muted-foreground">
            This challenge is locked until{" "}
            {new Date(Date.now() + secondsRemaining * 1000).toLocaleString()}
          </span>
        </div>
      </CardContent>
    </Card>
  );
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
          {Array((index % 3) + 1)
            .fill(null)
            .map((_, i) => (
              <Skeleton key={i} className="w-16 h-4" />
            ))}
        </div>
      </CardContent>
    </Card>
  );
}

export default ChallengeCard;
