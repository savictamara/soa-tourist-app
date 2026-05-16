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

export interface PublicUser {
  id: number;
  username: string;
  role: string;
  profileImage?: string | null;
  bio?: string | null;
}
