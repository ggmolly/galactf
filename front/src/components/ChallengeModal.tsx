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
import { FlagIcon, SparkleIcon } from "lucide-react";
import { CategoryBadge } from "./CategoryBadge";
import { useEffect, useState } from "react";
import { Skeleton } from "./ui/skeleton";
import { toast } from "sonner";
import { axiosErrorFactory } from "@/utils/errorFactory";
import { apiClient } from "@/lib/axios";
import { Attempt } from "@/interfaces/attempt.interface";

  
export interface ChallengeModalProps {
    challenge: Challenge;
    onClose: () => void;
    open: boolean;
}

export default function ChallengeModal({ challenge, open, onClose }: ChallengeModalProps) {
    const [solvers, setSolvers] = useState<Attempt[] | undefined>(undefined);
    
    const fetchChallengeSolvers = (id: number) => {
        apiClient.get(`/challenge/${id}/solvers`)
            .then((res) => {
                setSolvers(res.data.data);
            })
            .catch((err) => {
                const errMessage = axiosErrorFactory(err);
                toast.error(errMessage, {
                    duration: 5000,
                    dismissible: true,
                    action: {
                        label: 'Retry',
                        onClick: () => fetchChallengeSolvers(id),
                    },
                });
            })
    }

    useEffect(() => {
        fetchChallengeSolvers(challenge.id);
    }, [challenge]);

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
                    <SparkleIcon className="mr-2 h-4 w-4" />
                    {solvers === undefined ? (
                        <SolversSkeleton />
                    ) : (
                        <>
                            <span className="text-sm">Solvers ({solvers?.length ?? 0}):</span>
                            {solvers?.length > 0 ? (
                                solvers.map((solver, i) => (
                                    <>
                                        {i === 0 ? (
                                            <span className="text-sm ml-2 text-destructive-foreground font-bold">
                                                {solver.user.name}
                                            </span>
                                        ) : (
                                            <span className="text-sm ml-2">
                                                {solver.user.name}
                                            </span>
                                        )}
                                        {i !== solvers.length - 1 ? "," : ""}
                                    </>
                                ))
                            ) : (
                                <span className="text-sm ml-2 text-muted-foreground">No solvers</span>
                            )}
                        </>
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

function SolversSkeleton() {
    return (
        <div className="flex flex-wrap gap-2 items-center mt-1">
            <span className="text-sm">Solvers</span>
            <Skeleton className="w-6 h-6" />
            <span className="text-sm">: </span>
            <Skeleton className="w-36 h-6" />
        </div>
    );
}