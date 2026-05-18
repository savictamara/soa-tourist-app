import { Component, OnInit } from '@angular/core';
import { BlogApiService } from '../services/blog-api.service';
import { AuthStateService } from '../services/auth-state.service';
import { BlogLikeStatus } from '../models/blog-like-status.model';
import { BlogComment, BlogPost } from '../models/blog-post.model';
import { CreateBlogCommentRequest } from '../models/create-blog-comment-request.model';
import { CreateBlogPostRequest } from '../models/create-blog-post-request.model';
import { UpdateBlogCommentRequest } from '../models/update-blog-comment-request.model';

@Component({
  selector: 'app-blog-create',
  templateUrl: './blog-create.component.html',
  styleUrls: ['./blog-create.component.css']
})
export class BlogCreateComponent implements OnInit {
  form: CreateBlogPostRequest = {
    title: '',
    descriptionMarkdown: '',
    imageUrls: []
  };
  imageDragActive = false;
  isLoadingBlogs = false;
  isSubmitting = false;
  isCreateFormVisible = false;
  errorMessage = '';
  successMessage = '';
  blogs: BlogPost[] = [];
  currentImageIndexByBlogId: Record<number, number> = {};
  commentTextByBlogId: Record<number, string> = {};
  commentSubmittingByBlogId: Record<number, boolean> = {};
  editingCommentIdByBlogId: Record<number, number | null> = {};
  likeSubmittingByBlogId: Record<number, boolean> = {};

  constructor(
    private readonly blogApiService: BlogApiService,
    private readonly authStateService: AuthStateService
  ) {}

  ngOnInit(): void {
    this.loadBlogs();
  }

  get canCreateBlog(): boolean {
    const role = this.authStateService.currentUser?.role;
    return role === 'Guide' || role === 'Tourist';
  }

  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];

    if (file) {
      this.loadImage(file);
    }
  }

  onDragOver(event: DragEvent): void {
    event.preventDefault();
    this.imageDragActive = true;
  }

  onDragLeave(event: DragEvent): void {
    event.preventDefault();
    this.imageDragActive = false;
  }

  onDrop(event: DragEvent): void {
    event.preventDefault();
    this.imageDragActive = false;

    const file = event.dataTransfer?.files?.[0];

    if (file) {
      this.loadImage(file);
    }
  }

  removeImageUrl(index: number): void {
    this.form.imageUrls = (this.form.imageUrls ?? []).filter((_, imageIndex) => imageIndex !== index);
  }

  submit(): void {
    if (!this.canCreateBlog) {
      this.errorMessage = 'Only Guide and Tourist users can create blogs.';
      this.successMessage = '';
      return;
    }

    this.isSubmitting = true;
    this.errorMessage = '';
    this.successMessage = '';

    const request: CreateBlogPostRequest = {
      title: this.form.title.trim(),
      descriptionMarkdown: this.form.descriptionMarkdown.trim(),
      imageUrls: (this.form.imageUrls ?? []).filter(url => !!url.trim())
    };

    this.blogApiService.createBlog(request).subscribe({
      next: (createdBlog) => {
        this.blogs = [{
          ...createdBlog,
          likesCount: createdBlog.likesCount ?? 0,
          isLikedByCurrentUser: createdBlog.isLikedByCurrentUser ?? false,
          comments: createdBlog.comments ?? []
        }, ...this.blogs];
        this.currentImageIndexByBlogId[createdBlog.id] = 0;
        this.successMessage = 'Blog created successfully.';
        this.form = {
          title: '',
          descriptionMarkdown: '',
          imageUrls: []
        };
        this.imageDragActive = false;
        this.isCreateFormVisible = false;
        this.isSubmitting = false;
      },
      error: (error) => {
        this.errorMessage = error.error?.message ?? 'Could not create blog.';
        this.isSubmitting = false;
      }
    });
  }

  markdownPreview(): string {
    return this.renderMarkdown(this.form.descriptionMarkdown);
  }

  toggleCreateForm(): void {
    if (!this.canCreateBlog) {
      return;
    }

    this.isCreateFormVisible = !this.isCreateFormVisible;
    this.errorMessage = '';
    this.successMessage = '';
  }

  cancelCreate(): void {
    this.isCreateFormVisible = false;
    this.imageDragActive = false;
    this.form = {
      title: '',
      descriptionMarkdown: '',
      imageUrls: []
    };
  }

  renderedBlogMarkdown(markdown: string): string {
    return this.renderMarkdown(markdown);
  }

  canManageComments(): boolean {
    return this.canCreateBlog;
  }

  toggleLike(blog: BlogPost): void {
    if (!this.canCreateBlog || this.likeSubmittingByBlogId[blog.id]) {
      return;
    }

    this.likeSubmittingByBlogId[blog.id] = true;
    this.errorMessage = '';
    this.successMessage = '';

    const request$ = blog.isLikedByCurrentUser
      ? this.blogApiService.unlikeBlog(blog.id)
      : this.blogApiService.likeBlog(blog.id);

    request$.subscribe({
      next: (status) => {
        this.applyLikeStatus(status);
        this.likeSubmittingByBlogId[blog.id] = false;
      },
      error: (error) => {
        this.likeSubmittingByBlogId[blog.id] = false;
        this.errorMessage = error.error?.message ?? 'Could not update like.';
      }
    });
  }

  commentDraft(blogId: number): string {
    return this.commentTextByBlogId[blogId] ?? '';
  }

  isEditingComment(blogId: number, commentId: number): boolean {
    return this.editingCommentIdByBlogId[blogId] === commentId;
  }

  canEditComment(comment: BlogComment): boolean {
    return comment.authorUsername === this.authStateService.currentUser?.username;
  }

  beginEditComment(blogId: number, comment: BlogComment): void {
    this.editingCommentIdByBlogId[blogId] = comment.id;
    this.commentTextByBlogId[blogId] = comment.text;
    this.errorMessage = '';
    this.successMessage = '';
  }

  cancelCommentEdit(blogId: number): void {
    this.editingCommentIdByBlogId[blogId] = null;
    this.commentTextByBlogId[blogId] = '';
  }

  submitComment(blog: BlogPost): void {
    if (!this.canManageComments()) {
      this.errorMessage = 'Only Guide and Tourist users can comment on blogs.';
      this.successMessage = '';
      return;
    }

    const text = this.commentDraft(blog.id).trim();

    if (!text) {
      this.errorMessage = 'Comment text is required.';
      this.successMessage = '';
      return;
    }

    this.commentSubmittingByBlogId[blog.id] = true;
    this.errorMessage = '';
    this.successMessage = '';

    const editingCommentId = this.editingCommentIdByBlogId[blog.id];

    if (editingCommentId) {
      const updateRequest: UpdateBlogCommentRequest = { text };
      this.blogApiService.updateComment(blog.id, editingCommentId, updateRequest).subscribe({
        next: (updatedComment) => {
          this.updateBlogComment(blog.id, updatedComment);
          this.commentSubmittingByBlogId[blog.id] = false;
          this.editingCommentIdByBlogId[blog.id] = null;
          this.commentTextByBlogId[blog.id] = '';
          this.successMessage = 'Comment updated successfully.';
        },
        error: (error) => {
          this.commentSubmittingByBlogId[blog.id] = false;
          this.errorMessage = error.error?.message ?? 'Could not update comment.';
        }
      });
      return;
    }

    const createRequest: CreateBlogCommentRequest = { text };
    this.blogApiService.createComment(blog.id, createRequest).subscribe({
      next: (createdComment) => {
        this.prependBlogComment(blog.id, createdComment);
        this.commentSubmittingByBlogId[blog.id] = false;
        this.commentTextByBlogId[blog.id] = '';
        this.successMessage = 'Comment created successfully.';
      },
      error: (error) => {
        this.commentSubmittingByBlogId[blog.id] = false;
        if (error.status === 403) {
          this.errorMessage = 'You must follow the author before commenting on this blog.';
        } else {
          this.errorMessage = error.error?.message ?? 'Could not create comment.';
        }
      }
    });
  }

  currentBlogImage(blog: BlogPost): string {
    const imageUrls = blog.imageUrls ?? [];

    if (imageUrls.length === 0) {
      return '';
    }

    const currentIndex = this.currentImageIndexByBlogId[blog.id] ?? 0;
    const safeIndex = ((currentIndex % imageUrls.length) + imageUrls.length) % imageUrls.length;
    this.currentImageIndexByBlogId[blog.id] = safeIndex;
    return imageUrls[safeIndex];
  }

  previousBlogImage(blog: BlogPost): void {
    if (!blog.imageUrls || blog.imageUrls.length <= 1) {
      return;
    }

    const currentIndex = this.currentImageIndexByBlogId[blog.id] ?? 0;
    this.currentImageIndexByBlogId[blog.id] = (currentIndex - 1 + blog.imageUrls.length) % blog.imageUrls.length;
  }

  nextBlogImage(blog: BlogPost): void {
    if (!blog.imageUrls || blog.imageUrls.length <= 1) {
      return;
    }

    const currentIndex = this.currentImageIndexByBlogId[blog.id] ?? 0;
    this.currentImageIndexByBlogId[blog.id] = (currentIndex + 1) % blog.imageUrls.length;
  }

  currentBlogImageIndex(blog: BlogPost): number {
    return this.currentImageIndexByBlogId[blog.id] ?? 0;
  }

  private renderMarkdown(markdown: string): string {
    const escaped = markdown
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;');

    const lines = escaped.split(/\r?\n/);
    const renderedLines = lines.map(line => {
      if (line.startsWith('### ')) {
        return `<h3>${this.renderInlineMarkdown(line.slice(4))}</h3>`;
      }

      if (line.startsWith('## ')) {
        return `<h2>${this.renderInlineMarkdown(line.slice(3))}</h2>`;
      }

      if (line.startsWith('# ')) {
        return `<h1>${this.renderInlineMarkdown(line.slice(2))}</h1>`;
      }

      if (line.trim().startsWith('- ')) {
        return `<li>${this.renderInlineMarkdown(line.trim().slice(2))}</li>`;
      }

      if (!line.trim()) {
        return '';
      }

      return `<p>${this.renderInlineMarkdown(line)}</p>`;
    });

    const withLists: string[] = [];
    let listBuffer: string[] = [];

    const flushList = () => {
      if (listBuffer.length > 0) {
        withLists.push(`<ul>${listBuffer.join('')}</ul>`);
        listBuffer = [];
      }
    };

    renderedLines.forEach(line => {
      if (line.startsWith('<li>')) {
        listBuffer.push(line);
      } else {
        flushList();
        if (line) {
          withLists.push(line);
        }
      }
    });

    flushList();
    return withLists.join('');
  }

  private renderInlineMarkdown(text: string): string {
    return text
      .replace(/`([^`]+)`/g, '<code>$1</code>')
      .replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
      .replace(/\*([^*]+)\*/g, '<em>$1</em>')
      .replace(/\[([^\]]+)\]\((https?:\/\/[^)\s]+)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer">$1</a>');
  }

  private loadBlogs(): void {
    this.isLoadingBlogs = true;
    this.errorMessage = '';

    const user = this.authStateService.currentUser;
    console.log('Blog current user', user);

    if (!user) {
      this.blogs = [];
      this.isLoadingBlogs = false;
      return;
    }

    const userId = user.id;
    console.log('Loading followed blogs for userId', userId);

    this.blogApiService.getFollowedBlogs(userId).subscribe({
      next: (blogs) => {
        this.blogs = blogs.map(blog => ({
          ...blog,
          likesCount: blog.likesCount ?? 0,
          isLikedByCurrentUser: blog.isLikedByCurrentUser ?? false,
          comments: blog.comments ?? []
        }));
        this.currentImageIndexByBlogId = {};
        blogs.forEach(blog => {
          this.currentImageIndexByBlogId[blog.id] = 0;
        });
        this.isLoadingBlogs = false;
      },
      error: (error) => {
        this.errorMessage = error.error?.message ?? 'Could not load blogs.';
        this.isLoadingBlogs = false;
      }
    });
  }

  private loadImage(file: File): void {
    if (!file.type.startsWith('image/')) {
      this.errorMessage = 'Only image files are allowed.';
      return;
    }

    const reader = new FileReader();
    reader.onload = () => {
      const imageSource = typeof reader.result === 'string' ? reader.result : '';

      if (!imageSource) {
        this.errorMessage = 'Could not read the selected image.';
        return;
      }

      this.resizeImage(imageSource);
    };
    reader.readAsDataURL(file);
  }

  private resizeImage(imageSource: string): void {
    const image = new Image();

    image.onload = () => {
      const maxSize = 640;
      const scale = Math.min(maxSize / image.width, maxSize / image.height, 1);
      const width = Math.max(1, Math.round(image.width * scale));
      const height = Math.max(1, Math.round(image.height * scale));
      const canvas = document.createElement('canvas');

      canvas.width = width;
      canvas.height = height;

      const context = canvas.getContext('2d');

      if (!context) {
        this.errorMessage = 'Could not process the selected image.';
        return;
      }

      context.drawImage(image, 0, 0, width, height);
      const imageData = canvas.toDataURL('image/jpeg', 0.74);
      this.form.imageUrls = [...(this.form.imageUrls ?? []), imageData];
    };

    image.onerror = () => {
      this.errorMessage = 'Could not process the selected image.';
    };

    image.src = imageSource;
  }

  private prependBlogComment(blogId: number, createdComment: BlogComment): void {
    this.blogs = this.blogs.map(blog =>
      blog.id === blogId
        ? { ...blog, comments: [createdComment, ...(blog.comments ?? [])] }
        : blog
    );
  }

  private updateBlogComment(blogId: number, updatedComment: BlogComment): void {
    this.blogs = this.blogs.map(blog =>
      blog.id === blogId
        ? {
            ...blog,
            comments: (blog.comments ?? []).map(comment =>
              comment.id === updatedComment.id ? updatedComment : comment
            )
          }
        : blog
    );
  }

  private applyLikeStatus(status: BlogLikeStatus): void {
    this.blogs = this.blogs.map(blog =>
      blog.id === status.blogPostId
        ? {
            ...blog,
            likesCount: status.likesCount,
            isLikedByCurrentUser: status.isLikedByCurrentUser
          }
        : blog
    );
  }
}
