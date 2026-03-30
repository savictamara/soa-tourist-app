export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  role: string;
  firstName?: string;
  lastName?: string;
  profileImage?: string;
  biography?: string;
  motto?: string;
}
