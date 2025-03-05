import { User } from "./user.interface";

export interface Attempt {
    id:           number;
    user_id:      number;
    challenge_id: number;
    user:         User;
}
