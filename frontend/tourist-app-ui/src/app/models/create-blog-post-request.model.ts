export interface CreateBlogPostRequest {
  title: string;
  descriptionMarkdown: string;
  imageUrls?: string[];
}
