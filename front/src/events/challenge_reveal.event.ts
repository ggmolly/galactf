import { ChallengeWithSolveRate } from "@/interfaces/challenge.interface";
import { ChallengeReveal } from "@/proto/challenge_reveal";

export function handleChalReveal(
  event: ChallengeReveal,
  parameters: {
    challenges: React.RefObject<ChallengeWithSolveRate[]>;
    setChallenges: React.Dispatch<React.SetStateAction<ChallengeWithSolveRate[]>>;
  }
) {
  const { setChallenges } = parameters;

  setChallenges((prevChallenges) =>
    prevChallenges.map((challenge) => {
      if (challenge.id !== event.id) {
        return challenge;
      }

      return {
        name: event.name,
        difficulty: event.difficulty,
        categories: event.categories,
        attachments: event.attachments
          .filter((attachment) => attachment.filename.trim() !== "")
          .map((attachment) => ({
            ...attachment,
            type: attachment.type as "url" | "file",
          })),
        solve_rate: 0,
        solved: false,
        solvers: 0,
        total_attempts: 0,
        id: challenge.id,
        description: challenge.description,
        reveal_in: undefined,
      };
    })
  );
}
