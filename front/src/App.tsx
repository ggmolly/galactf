import { useEffect, useState } from "react";
import { apiClient } from "./lib/axios";
import { ChallengeWithSolveRate } from "./interfaces/challenge.interface";
import { toast } from "sonner";
import ChallengeCard, {
  ChallengeCardSkeleton,
} from "./components/ChallengeCard";
import ChallengeModal from "./components/ChallengeModal";
import { axiosErrorFactory } from "./utils/errorFactory";
import { useAuth } from "./providers/auth.provider";

function App() {
  const [challenges, setChallenges] = useState<
    ChallengeWithSolveRate[] | undefined
  >(undefined);
  const [selectedChallenge, setSelectedChallenge] = useState<
    ChallengeWithSolveRate | undefined
  >(undefined);
  const [modalVisible, setModalVisible] = useState(false);
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
    if (!user) {
      return;
    }
    loadChallenges();
  }, [user]);

  return (
    <>
      <div className="flex flex-wrap justify-around gap-4 p-4">
        {challenges === undefined
          ? Array(22)
              .fill(null)
              .map((_, i) => <ChallengeCardSkeleton key={i} index={i} />)
          : challenges.map((challenge) => (
              <ChallengeCard
                key={challenge.id}
                challenge={challenge}
                selectChallenge={(challenge) => {
                  setSelectedChallenge(challenge);
                  setModalVisible(true);
                }}
              />
            ))}
      </div>
      {selectedChallenge !== undefined && (
        <ChallengeModal
          challenge={selectedChallenge}
          onClose={() => {
            setModalVisible(false);
          }}
          open={modalVisible}
        />
      )}
    </>
  );
}

export default App;
