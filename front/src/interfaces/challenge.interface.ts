import { Attachment } from "./attachment.interface";

export interface Challenge {
  id: number;
  name: string;
  description: string;
  difficulty: number;
  categories: string[];
  attachments: Attachment[];
}

export interface ChallengeWithSolveRate extends Challenge {
  solve_rate: number;
}
