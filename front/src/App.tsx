import { useEffect, useState } from 'react'
import { apiClient } from './lib/axios';
import { ChallengeWithSolveRate } from './interfaces/challenge.interface';
import { Toaster } from './components/ui/sonner';
import { toast } from 'sonner';
import axios from 'axios';
import ChallengeCard, { ChallengeCardSkeleton } from './components/ChallengeCard';
import ChallengeModal from './components/ChallengeModal';

function App() {
  const [challenges, setChallenges] = useState<ChallengeWithSolveRate[] | undefined>(undefined)
  const [selectedChallenge, setSelectedChallenge] = useState<ChallengeWithSolveRate | undefined>(undefined)
  const [modalVisible, setModalVisible] = useState(false);

  const loadChallenges = () => {
    apiClient.get("/api/v1/challenges").then((res) => {
      setChallenges(res.data.data);
    }).catch((err) => {
      let errorMessage = "Network error. Please check your internet connection.";
      if (!axios.isAxiosError(err)) {
        errorMessage = `HTTP ${err.response.status}, failed to load challenges`;
      }
      toast.error(errorMessage, {
        duration: 5000,
        dismissible: true,
        action: {
          label: 'Retry',
          onClick: () => loadChallenges(),
        },
      });
    })
  }

  useEffect(() => {
    loadChallenges()
  }, []);

  return (
    <>
      <div className='flex flex-wrap justify-around gap-4 p-4'>
        {challenges === undefined ? (
          Array(22).fill(null).map((_, i) => (
            <ChallengeCardSkeleton key={i} index={i} />
          ))
        ) : (
          challenges.map((challenge) => (
            <ChallengeCard
              key={challenge.id}
              challenge={challenge}
              selectChallenge={(challenge) => {
                setSelectedChallenge(challenge)
                setModalVisible(true)
              }}
            />
          ))
        )}
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
      <Toaster position='top-right' />
    </>
  )
}

export default App
