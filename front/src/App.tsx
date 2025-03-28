import { useState } from "react";
import { ChallengeWithSolveRate } from "./interfaces/challenge.interface";
import ChallengeCard, {
  ChallengeCardSkeleton,
} from "./components/ChallengeCard";
import ChallengeModal from "./components/ChallengeModal";
import { Tabs, TabsList, TabsContent, TabsTrigger } from "./components/ui/tabs";
import { useChallenges } from "./providers/challenges.provider";
import { Leaderboard } from "./components/Leaderboard";

function App() {
  const [selectedChallenge, setSelectedChallenge] = useState<
    ChallengeWithSolveRate | undefined
  >(undefined);
  const [modalVisible, setModalVisible] = useState(false);
  const { challenges } = useChallenges();

  return (
    <Tabs
      defaultValue="challenges"
      className="size-full flex flex-col items-center"
    >
      <TabsList className="size-full bg-transparent p-4">
        <TabsTrigger value="challenges" className="text-xl">
          Challenges
        </TabsTrigger>
        <TabsTrigger value="leaderboard" className="text-xl">
          Leaderboard
        </TabsTrigger>
      </TabsList>
      <TabsContent value="challenges">
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
      </TabsContent>
      <TabsContent value="leaderboard" className="w-2/3 py-6">
        <Leaderboard />
      </TabsContent>
    </Tabs>
  );
}

export default App;
