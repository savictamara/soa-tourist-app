export interface BlogComment {
  id: number;
  authorUsername: string;
  authorRole: string;
  text: string;
  createdAtUtc: string;
  lastModifiedAtUtc: string;
}

export interface BlogPost {
  id: number;
  title: string;
  descriptionMarkdown: string;
  createdAtUtc: string;
  authorUsername: string;
  imageUrls: string[];
  likesCount: number;
  isLikedByCurrentUser: boolean;
  comments: BlogComment[];
}
