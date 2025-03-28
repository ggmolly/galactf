import { useLeaderboard } from "@/providers/leaderboard.provider";
import { useAuth } from "@/providers/auth.provider";
import { cn } from "@/lib/utils";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "./ui/table";

export const Leaderboard = () => {
  const leaderboard = useLeaderboard();
  const { user } = useAuth();

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-[100px]">Rank</TableHead>
          <TableHead>Name</TableHead>
          <TableHead>Score</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {leaderboard.users.map((u, index) => (
          <TableRow key={index}>
            <TableCell
              className={cn("font-medium", {
                "text-primary": index === 0,
                "text-secondary": index === 1,
                "text-accent": index === 2,
              })}
            >
              {index + 1}
            </TableCell>
            <TableCell className={cn({ "font-bold": user.id === u.id })}>
              {u.name}
            </TableCell>
            <TableCell className="">{u.score}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
};
