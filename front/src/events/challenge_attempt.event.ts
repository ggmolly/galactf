import { ChallengeWithSolveRate } from "@/interfaces/challenge.interface";
import { User } from "@/interfaces/user.interface";
import { ChallengeAttempt } from "@/proto/challenge_attempt";
import { toast } from "sonner";

export function handleChalAttempt(
  event: ChallengeAttempt,
  parameters: {
    challenges: React.RefObject<ChallengeWithSolveRate[]>;
    setChallenges: React.Dispatch<React.SetStateAction<ChallengeWithSolveRate[]>>;
    user: User;
  }
) {
  const { challenges, setChallenges, user } = parameters;

  setChallenges((prevChallenges) =>
    prevChallenges.map((challenge) => {
      if (challenge.id !== event.challengeId) {
        return challenge;
      }

      const newTotalAttempts = challenge.total_attempts + 1;
      const newSolvers = event.success ? challenge.solvers + 1 : challenge.solvers;
      const newSolveRate = newSolvers / newTotalAttempts;

      return {
        ...challenge,
        total_attempts: newTotalAttempts,
        solvers: newSolvers,
        solve_rate: newSolveRate,
        solved: user.id === event.user!.id && event.success ? true : challenge.solved,
      };
    })
  );

  if (event.firstBlood && event.user!.id !== user.id) {
    const chal = challenges.current.find((c) => c.id === event.challengeId);
    toast.info(`${event.user!.name} got the first blood on ${chal?.name}!`, {
      duration: 60000,
      dismissible: true,
    });
  }
}
