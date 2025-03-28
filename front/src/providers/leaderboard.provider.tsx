import { ChallengeWithSolveRate } from "@/interfaces/challenge.interface";
import { useChallenges } from "./challenges.provider";
import {
  createContext,
  ReactNode,
  useContext,
  useEffect,
  useState,
  useMemo,
  Dispatch,
} from "react";
import { useAuth } from "./auth.provider";
import { apiClient } from "@/lib/axios";
import { axiosErrorFactory } from "@/utils/errorFactory";
import { toast } from "sonner";
import { computePoints } from "@/utils/computePoints";

interface LeaderboardProviderProps {
  children: ReactNode;
}

export interface Solver {
    user: {
        id: number;
        name: string;
    }
    challengeId: number;
}

interface LeaderboardContextType {
  users: Array<{ id: number; name: string; score: number }>;
  setSolvers: Dispatch<React.SetStateAction<Solver[]>>;
}

const LeaderboardContext = createContext<LeaderboardContextType | undefined>(
  undefined,
);

export const LeaderboardProvider: React.FC<LeaderboardProviderProps> = ({
  children,
}) => {
  const { challenges } = useChallenges();
  const [solvers, setSolvers] = useState<Solver[]>([]);
  const { user } = useAuth();

  const loadSolvers = () => {
    apiClient
      .get("/leaderboard")
      .then((res) => {
        setSolvers(res.data.data.map((d: any): Solver => ({
            challengeId: d.challenge_id,
            user: {
                id: d.user.id,
                name: d.user.name,
            }
        })));
      })
      .catch((err) => {
        const errorMessage = axiosErrorFactory(err);
        toast.error(errorMessage, {
          duration: 5000,
          dismissible: true,
          action: {
            label: "Retry",
            onClick: () => loadSolvers(),
          },
        });
      });
  };

  useEffect(() => {
    if (!user.id) {
      return;
    }
    loadSolvers();
  }, [user]);

  const givenPointsByChallengeId = useMemo(() => getGivenPointsByChallengeId(challenges), [challenges]);
  const users = useMemo(
    () =>
      [
        ...solvers
          .reduce((acc, s) => {
            const givenPoints =
              givenPointsByChallengeId.get(s.challengeId) ?? 0;
            const existingUser = acc.get(s.user.id);

            acc.set(s.user.id, {
              id: s.user.id,
              name: s.user.name,
              score:
                existingUser !== undefined
                  ? existingUser.score + givenPoints
                  : givenPoints,
            });

            return acc;
          }, new Map<number, { id: number; name: string; score: number }>())
          .values(),
      ].sort((a, b) => b.score - a.score),
    [solvers, givenPointsByChallengeId]
  );

  return (
    <LeaderboardContext.Provider value={{ users, setSolvers }}>
      {children}
    </LeaderboardContext.Provider>
  );
};

export const useLeaderboard = (): LeaderboardContextType => {
  const context = useContext(LeaderboardContext);
  if (context === undefined) {
    throw new Error("useLeaderboard must be used within a LeaderboardProvider");
  }
  return context;
};

const getGivenPointsByChallengeId = (
  challenges: ChallengeWithSolveRate[],
): Map<number, number> => {
  const challengesMap = new Map<number, number>();
  challenges.forEach((challenge) => {
    challengesMap.set(challenge.id, computePoints(challenge.solvers));
  });
  return challengesMap;
};
