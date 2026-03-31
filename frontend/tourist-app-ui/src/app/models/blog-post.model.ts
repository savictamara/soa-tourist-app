export interface BlogPost {
  id: number;
  title: string;
  descriptionMarkdown: string;
  createdAtUtc: string;
  authorUsername: string;
  imageUrls: string[];
}
