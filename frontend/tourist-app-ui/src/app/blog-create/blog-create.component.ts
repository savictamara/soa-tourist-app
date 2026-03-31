import { Component, OnInit } from '@angular/core';
import { BlogApiService } from '../services/blog-api.service';
import { AuthStateService } from '../services/auth-state.service';
import { BlogPost } from '../models/blog-post.model';
import { CreateBlogPostRequest } from '../models/create-blog-post-request.model';

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
        this.blogs = [createdBlog, ...this.blogs];
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

    this.blogApiService.getBlogs().subscribe({
      next: (blogs) => {
        this.blogs = blogs;
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
}
