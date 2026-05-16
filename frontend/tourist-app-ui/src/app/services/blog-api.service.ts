import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { BlogComment, BlogPost } from '../models/blog-post.model';
import { BlogLikeStatus } from '../models/blog-like-status.model';
import { CreateBlogPostRequest } from '../models/create-blog-post-request.model';
import { CreateBlogCommentRequest } from '../models/create-blog-comment-request.model';
import { UpdateBlogCommentRequest } from '../models/update-blog-comment-request.model';

@Injectable({
  providedIn: 'root'
})
export class BlogApiService {
  private readonly baseUrl = 'http://localhost:8080/api';

  constructor(private readonly http: HttpClient) {}

  getBlogs(): Observable<BlogPost[]> {
    return this.http.get<BlogPost[]>(`${this.baseUrl}/blogs`);
  }

  createBlog(request: CreateBlogPostRequest): Observable<BlogPost> {
    return this.http.post<BlogPost>(`${this.baseUrl}/blogs`, request);
  }

  createComment(blogId: number, request: CreateBlogCommentRequest): Observable<BlogComment> {
    return this.http.post<BlogComment>(`${this.baseUrl}/blogs/${blogId}/comments`, request);
  }

  updateComment(blogId: number, commentId: number, request: UpdateBlogCommentRequest): Observable<BlogComment> {
    return this.http.put<BlogComment>(`${this.baseUrl}/blogs/${blogId}/comments/${commentId}`, request);
  }

  likeBlog(blogId: number): Observable<BlogLikeStatus> {
    return this.http.post<BlogLikeStatus>(`${this.baseUrl}/blogs/${blogId}/likes`, {});
  }

  unlikeBlog(blogId: number): Observable<BlogLikeStatus> {
    return this.http.delete<BlogLikeStatus>(`${this.baseUrl}/blogs/${blogId}/likes`);
  }
}
