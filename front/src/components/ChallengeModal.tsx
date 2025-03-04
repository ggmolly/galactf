import { Challenge } from "@/interfaces/challenge.interface";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
  } from "@/components/ui/dialog"
import { DifficultyBadge } from "./DifficultyBadge";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import { FlagIcon, HeartPulseIcon, SparkleIcon } from "lucide-react";
import { CategoryBadge } from "./CategoryBadge";
import { formatDuration } from "@/utils/formatDuration";

  
export interface ChallengeModalProps {
    challenge: Challenge;
    onClose: () => void;
    open: boolean;
}

export default function ChallengeModal({ challenge, open, onClose }: ChallengeModalProps) {
    // TODO: Remove mocked data
    const firstBloodTime = 60 + 37;
    const solvers: string[] = ["TODO", "molly", "test", "john"];

    return (
        <Dialog
            open={open}
            onOpenChange={onClose}
        >
            <DialogContent>
                <DialogHeader>
                    <DialogTitle className="flex gap-x-2 items-center">
                        <span>{challenge?.name}</span>
                        <DifficultyBadge difficulty={challenge?.difficulty ?? 0} />
                    </DialogTitle>
                </DialogHeader>
                <DialogDescription>
                    {challenge?.description ?? "No description available"}
                </DialogDescription>

                <div className="flex flex-wrap gap-2 items-center">
                    <span className="text-sm">Categories ({challenge?.categories.length ?? 0}):</span>
                    {challenge?.categories.map((category, i) => (
                        <CategoryBadge key={i} category={category} />
                    ))}
                </div>

                <div className="flex flex-wrap items-center mt-1">
                    <HeartPulseIcon className="mr-2 h-4 w-4 text-destructive-foreground" />
                    <span className="text-sm text-destructive-foreground">First blood:</span>
                    <span className="text-sm font-bold ml-2">
                        {solvers.length > 0 ? (
                            solvers[0] + ` (${formatDuration(firstBloodTime)})`
                        ) : (
                            <span className="text-muted-foreground">No solvers</span>
                        )}
                    </span>
                </div>

                <div className="flex flex-wrap items-center mt-1">
                    <SparkleIcon className="mr-2 h-4 w-4" />
                    <span className="text-sm">Solvers ({solvers.length}):</span>
                    {solvers.length > 0 ? (
                        solvers.map((solver, i) => (
                            <span key={i} className="text-sm ml-2">
                                {solver}{i !== solvers.length - 1 ? "," : ""}
                            </span>
                        ))
                    ) : (
                        <span className="text-sm ml-2 text-muted-foreground">No solvers</span>
                    )}
                </div>

                <hr />
                <DialogFooter>
                    <Input
                        minLength={0}
                        maxLength={255}
                        placeholder="GALA{...}"
                        className="w-full"
                    />
                    <Button variant="outline">
                        <FlagIcon className="mr-2 h-4 w-4" />
                        Submit
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}