const POINTS_BASE = 500;
const LAMBDA = 0.035;

export function computePoints(solvers: number): number {
  return Math.min(Math.round(POINTS_BASE * Math.exp(-LAMBDA * (solvers-1))), 500);
}
