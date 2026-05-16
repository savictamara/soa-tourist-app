export interface FollowerUser {
  userId: string;
  username: string;
  role: string;
}

export interface FollowRequest {
  followerId: string;
  followingId: string;
}
