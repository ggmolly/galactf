import { Attachment } from "./attachment.interface";

export interface Challenge {
  id: number;
  name: string;
  description: string;
  difficulty: number;
  categories: string[];
  attachments: Attachment[];
  solved?: boolean;
  solvers?: number;
  reveal_in?: number;
}

export interface ChallengeWithSolveRate extends Challenge {
  solve_rate: number;
}
