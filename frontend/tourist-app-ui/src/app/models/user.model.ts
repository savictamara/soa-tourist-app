export interface User {
  id: number;
  username: string;
  email: string;
  role: string;
  isBlocked: boolean;
  firstName?: string | null;
  lastName?: string | null;
  profileImage?: string | null;
  biography?: string | null;
  motto?: string | null;
}
