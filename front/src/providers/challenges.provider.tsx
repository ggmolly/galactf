import { ChallengeWithSolveRate } from "@/interfaces/challenge.interface";
import {
  createContext,
  ReactNode,
  useContext,
  useEffect,
  useState,
} from "react";
import { useAuth } from "./auth.provider";
import { apiClient } from "@/lib/axios";
import { axiosErrorFactory } from "@/utils/errorFactory";
import { toast } from "sonner";

interface ChallengesProviderProps {
  children: ReactNode;
}

interface ChallengesContextType {
  challenges: ChallengeWithSolveRate[];
  setChallenges: (challenges: ChallengeWithSolveRate[]) => void;
}

const ChallengesContext = createContext<ChallengesContextType | undefined>(
  undefined
);

export const ChallengesProvider: React.FC<ChallengesProviderProps> = ({
  children,
}) => {
  const [challenges, setChallenges] = useState<ChallengeWithSolveRate[]>([]);

  const { user } = useAuth();

  const loadChallenges = () => {
    apiClient
      .get("/challenges")
      .then((res) => {
        setChallenges(res.data.data);
      })
      .catch((err) => {
        const errorMessage = axiosErrorFactory(err);
        toast.error(errorMessage, {
          duration: 5000,
          dismissible: true,
          action: {
            label: "Retry",
            onClick: () => loadChallenges(),
          },
        });
      });
  };

  useEffect(() => {
    if (!user.id) {
      return;
    }
    loadChallenges();
  }, [user]);

  return (
    <ChallengesContext.Provider value={{ challenges, setChallenges }}>
      {children}
    </ChallengesContext.Provider>
  );
};

export const useChallenges = (): ChallengesContextType => {
  const context = useContext(ChallengesContext);
  if (context === undefined) {
    throw new Error("useChallenges must be used within a ChallengesProvider");
  }
  return context;
};
