import { ChallengeWithSolveRate } from "@/interfaces/challenge.interface";
import { User } from "@/interfaces/user.interface";
import { ChallengeAttempt } from "@/proto/challenge_attempt";

export function handleChalAttempt(
  event: ChallengeAttempt,
  parameters: {
    setChallenges: React.Dispatch<React.SetStateAction<ChallengeWithSolveRate[]>>;
    user: User;
  }
) {
  const { setChallenges, user } = parameters;

  setChallenges((prevChallenges) =>
    prevChallenges.map((challenge) => {
      if (challenge.id !== event.challengeId) {
        return challenge;
      }

      return {
        ...challenge,
        total_attempts: challenge.total_attempts + 1,
        solvers: event.success ? challenge.solvers + 1 : challenge.solvers,
        solved: user.id === event.user!.id && event.success ? true : challenge.solved,
      };
    })
  );
}
