import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { BlogPost } from '../models/blog-post.model';
import { CreateBlogPostRequest } from '../models/create-blog-post-request.model';

@Injectable({
  providedIn: 'root'
})
export class BlogApiService {
  private readonly baseUrl = 'http://localhost:5002/api';

  constructor(private readonly http: HttpClient) {}

  getBlogs(): Observable<BlogPost[]> {
    return this.http.get<BlogPost[]>(`${this.baseUrl}/blogs`);
  }

  createBlog(request: CreateBlogPostRequest): Observable<BlogPost> {
    return this.http.post<BlogPost>(`${this.baseUrl}/blogs`, request);
  }
}
