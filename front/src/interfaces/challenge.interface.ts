import { Attachment } from "./attachment.interface";

export interface Challenge {
  id: number;
  name: string;
  description: string;
  difficulty: number;
  categories: string[];
  attachments: Attachment[];
  reveal_in?: number;
}

export interface ChallengeWithSolveRate extends Challenge {
  solve_rate: number;
  solved: boolean;
  solvers: number;
  total_attempts: number;
}
