import { Challenge } from "@/interfaces/challenge.interface";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { DifficultyBadge } from "./DifficultyBadge";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import {
  ExternalLinkIcon,
  FileIcon,
  FlagIcon,
  LightbulbIcon,
  LinkIcon,
  SparklesIcon,
} from "lucide-react";
import { useEffect, useState } from "react";
import { Skeleton } from "./ui/skeleton";
import { toast } from "sonner";
import { axiosErrorFactory } from "@/utils/errorFactory";
import { apiClient } from "@/lib/axios";
import { Attempt } from "@/interfaces/attempt.interface";
import { formatBytes } from "@/utils/formatBytes";
import { cn } from "@/lib/utils";
import { useAuth } from "@/providers/auth.provider";

const flagRegex = /^GALA{[A-Za-z0-9_-]{24,48}}$/;
export interface ChallengeModalProps {
  challenge: Challenge;
  onClose: () => void;
  open: boolean;
}

export default function ChallengeModal({
  challenge,
  open,
  onClose,
}: ChallengeModalProps) {
  const [solvers, setSolvers] = useState<Attempt[] | undefined>(undefined);
  const [flag, setFlag] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState<boolean>(false);

  const { user } = useAuth();

  const fetchChallengeSolvers = (id: number) => {
    apiClient
      .get(`/challenge/${id}/solvers`)
      .then((res) => {
        setSolvers(res.data.data);
      })
      .catch((err) => {
        const errMessage = axiosErrorFactory(err);
        toast.error(errMessage, {
          duration: 5000,
          dismissible: true,
          action: {
            label: "Retry",
            onClick: () => fetchChallengeSolvers(id),
          },
        });
      });
  };

  const submitFlag = () => {
    if (flag === null) {
      return;
    }
    setSubmitting(true);
    apiClient
      .post(`/challenge/${challenge.id}/submit`, { flag })
      .then((response) => {
        switch (response.status) {
          case 200:
            toast.success(response.data.message);
            onClose();
            break;
          case 201:
            toast.error(response.data.message);
            break;
        }
      })
      .catch((err) => {
        const errMessage = axiosErrorFactory(err);
        toast.error(errMessage, {
          duration: 5000,
          dismissible: true,
        });
      })
      .finally(() => {
        setSubmitting(false);
      });
  };

  useEffect(() => {
    fetchChallengeSolvers(challenge.id);
  }, [challenge]);

  return (
    <Dialog open={open} onOpenChange={onClose}>
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

        <div className="flex flex-wrap items-center mt-1 text-xs gap-2">
          <SparklesIcon className="mr-2 h-4 w-4 text-destructive-foreground" />
          <p>First blood:</p>
          {solvers === undefined ? (
            <Skeleton className="w-16 h-6" />
          ) : solvers.length > 0 ? (
            <span className={cn("font-bold", user.id === solvers[0].user.id)}>
              {solvers[0].user.name}
            </span>
          ) : (
            <span className="font-bold">Nobody</span>
          )}
        </div>

        <div className="flex flex-wrap items-center mt-1 text-xs">
          <LightbulbIcon className="mr-2 h-4 w-4" />
          {solvers === undefined ? (
            <SolversSkeleton />
          ) : (
            <div className="flex flex-col gap-2">
              <span>Solvers ({solvers?.length ?? 0}):</span>
              <div className="flex flex-wrap gap-2 items-center">
                {solvers.map((solver, i) => (
                  <span
                    key={i}
                    className={cn("font-bold", user.id === solver.user.id)}
                  >
                    {solver.user.name + (i === solvers.length - 1 ? "" : ", ")}
                  </span>
                ))}

                {solvers.length === 0 && (
                  <span className="font-bold">Nobody</span>
                )}
              </div>
            </div>
          )}
        </div>
        {challenge?.attachments.length > 0 && (
          <>
            <hr />
            <span className="text-sm">
              Attachments ({challenge?.attachments.length}):
            </span>
            <div className="flex flex-col">
              {challenge?.attachments.map((attachment, i) => (
                <div key={i} className="flex gap-x-2 items-center">
                  <a
                    className="flex gap-x-2 text-sm items-center"
                    href={import.meta.env.VITE_API_URL + attachment.url}
                    target="_blank"
                  >
                    {attachment.type === "url" ? (
                      <LinkIcon className="h-4 w-4" />
                    ) : (
                      <FileIcon className="h-4 w-4" />
                    )}
                    {attachment.filename + " "}
                    {attachment.type === "file" &&
                      attachment.size > 0 &&
                      "(" + formatBytes(attachment.size) + ")"}
                    <ExternalLinkIcon className="h-4 w-4 mb-2" />
                  </a>
                </div>
              ))}
            </div>
          </>
        )}
        <hr />
        <DialogFooter>
          <Input
            minLength={0}
            maxLength={255}
            placeholder={
              challenge.solved ? "You've solved this challenge!" : "GALA{...}"
            }
            className="w-full"
            value={flag ?? ""}
            onChange={(e) => setFlag(e.target.value)}
            disabled={submitting || challenge.solved}
          />
          <Button
            variant="outline"
            onClick={submitFlag}
            disabled={
              (flag !== null && !flag.match(flagRegex)) ||
              submitting ||
              challenge.solved
            }
          >
            <FlagIcon className="mr-2 h-4 w-4" />
            {challenge.solved
              ? "Solved"
              : submitting
                ? "Submitting..."
                : "Submit"}
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
