export interface Challenge {
    id: number;
    name: string;
    description: string;
    difficulty: number;
    categories: string[];
}

export interface ChallengeWithSolveRate extends Challenge {
    solve_rate: number;
}