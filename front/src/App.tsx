import { useState } from "react";
import { ChallengeWithSolveRate } from "./interfaces/challenge.interface";
import ChallengeCard, {
  ChallengeCardSkeleton,
} from "./components/ChallengeCard";
import ChallengeModal from "./components/ChallengeModal";
import { useChallenges } from "./providers/challenges.provider";

function App() {
  const [selectedChallenge, setSelectedChallenge] = useState<
    ChallengeWithSolveRate | undefined
  >(undefined);
  const [modalVisible, setModalVisible] = useState(false);
  const { challenges } = useChallenges();

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
