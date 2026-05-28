import { Component, OnInit } from '@angular/core';
import { KeyPoint } from '../models/key-point.model';
import { CreateReviewRequest, Review, Tour } from '../models/tour.model';
import { AuthStateService } from '../services/auth-state.service';
import { TourApiService } from '../services/tour-api.service';

@Component({
  selector: 'app-tourist-tours',
  templateUrl: './tourist-tours.component.html',
  styleUrls: ['./tourist-tours.component.css']
})
export class TouristToursComponent implements OnInit {
  tours: Tour[] = [];
  selectedTour: Tour | null = null;
  keyPoints: KeyPoint[] = [];
  reviews: Review[] = [];
  isLoadingTours = false;
  isLoadingDetails = false;
  isSubmittingReview = false;
  successMessage = '';
  errorMessage = '';
  reviewImageDragActive = false;
  reviewForm = {
    rating: 5,
    comment: '',
    visitedDate: '',
    images: [] as string[]
  };

  constructor(
    private readonly tourApiService: TourApiService,
    private readonly authStateService: AuthStateService
  ) {}

  get todayString(): string {
    return new Date().toISOString().split('T')[0];
  }

  ngOnInit(): void {
    this.loadTours();
  }

  loadTours(): void {
    this.isLoadingTours = true;
    this.errorMessage = '';
    this.tourApiService.getPublishedTours().subscribe({
      next: (tours) => {
        this.tours = (tours ?? []).map(tour => ({
          ...tour,
          tags: Array.isArray(tour.tags) ? tour.tags : [],
          keyPoints: Array.isArray(tour.keyPoints) ? tour.keyPoints : [],
          reviews: Array.isArray(tour.reviews) ? tour.reviews : [],
          durations: Array.isArray(tour.durations) ? tour.durations : [],
          lengthKm: Number(tour.lengthKm ?? 0)
        }));
        this.isLoadingTours = false;
      },
      error: (error) => {
        this.errorMessage = error?.error?.error ?? 'Could not load tours.';
        this.isLoadingTours = false;
      }
    });
  }

  selectTour(tourId: string): void {
    this.isLoadingDetails = true;
    this.errorMessage = '';
    const tour = this.tours.find(item => item.id === tourId) ?? null;
    if (!tour) {
      this.errorMessage = 'Could not load selected tour.';
      this.isLoadingDetails = false;
      return;
    }
    this.selectedTour = tour;
    this.keyPoints = (tour.keyPoints ?? []).slice(0, 1);
    this.loadReviews(tourId);
  }

  submitReview(): void {
    if (!this.selectedTour?.id) {
      this.errorMessage = 'Select a tour first.';
      return;
    }
    const currentUser = this.authStateService.currentUser;
    const touristUsername = currentUser?.username?.trim() ?? '';
    const touristId = currentUser?.id ? String(currentUser.id) : touristUsername;
    if (!touristId || !touristUsername) {
      this.errorMessage = 'Logged user is required.';
      return;
    }
    if (!this.reviewForm.comment.trim() || !this.reviewForm.visitedDate) {
      this.errorMessage = 'Comment and visited date are required.';
      return;
    }
    if (this.reviewForm.visitedDate > this.todayString) {
      this.errorMessage = 'Visited date cannot be in the future.';
      return;
    }

    this.isSubmittingReview = true;
    this.errorMessage = '';
    this.successMessage = '';
    const payload: CreateReviewRequest = {
      rating: this.reviewForm.rating,
      comment: this.reviewForm.comment.trim(),
      touristId,
      touristUsername,
      visitedDate: this.reviewForm.visitedDate,
      images: [...this.reviewForm.images]
    };
    this.tourApiService.addReview(this.selectedTour.id, payload).subscribe({
      next: () => {
        this.reviewForm = { rating: 5, comment: '', visitedDate: '', images: [] };
        this.successMessage = 'Review submitted successfully.';
        this.refreshSelectedTour(this.selectedTour!.id);
        this.isSubmittingReview = false;
      },
      error: (error) => {
        this.errorMessage = error?.error?.error ?? 'Could not submit review.';
        this.isSubmittingReview = false;
      }
    });
  }

  onReviewImageDragOver(event: DragEvent): void {
    event.preventDefault();
    this.reviewImageDragActive = true;
  }

  onReviewImageDragLeave(event: DragEvent): void {
    event.preventDefault();
    this.reviewImageDragActive = false;
  }

  onReviewImageDrop(event: DragEvent): void {
    event.preventDefault();
    this.reviewImageDragActive = false;
    const files = event.dataTransfer?.files;
    if (!files) {
      return;
    }
    for (let i = 0; i < files.length; i += 1) {
      const file = files.item(i);
      if (file) {
        this.readReviewImage(file);
      }
    }
  }

  onReviewImageSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    const files = input.files;
    if (!files) {
      return;
    }
    for (let i = 0; i < files.length; i += 1) {
      const file = files.item(i);
      if (file) {
        this.readReviewImage(file);
      }
    }
    input.value = '';
  }

  removeReviewImage(index: number): void {
    this.reviewForm.images = this.reviewForm.images.filter((_, i) => i !== index);
  }

  private loadReviews(tourId: string): void {
    this.tourApiService.getReviews(tourId).subscribe({
      next: (reviews) => {
        this.reviews = (reviews ?? []).map(review => ({
          ...review,
          images: Array.isArray(review.images) ? review.images : []
        }));
        this.isLoadingDetails = false;
      },
      error: (error) => {
        this.errorMessage = error?.error?.error ?? 'Could not load reviews.';
        this.isLoadingDetails = false;
      }
    });
  }

  private refreshSelectedTour(tourId: string): void {
    this.tourApiService.getPublishedTours().subscribe({
      next: (tours) => {
        this.tours = (tours ?? []).map(tour => ({
          ...tour,
          tags: Array.isArray(tour.tags) ? tour.tags : [],
          keyPoints: Array.isArray(tour.keyPoints) ? tour.keyPoints : [],
          reviews: Array.isArray(tour.reviews) ? tour.reviews : [],
          durations: Array.isArray(tour.durations) ? tour.durations : [],
          lengthKm: Number(tour.lengthKm ?? 0)
        }));
        this.selectedTour = this.tours.find(tour => tour.id === tourId) ?? this.selectedTour;
        this.keyPoints = (this.selectedTour?.keyPoints ?? []).slice(0, 1);
        this.loadReviews(tourId);
      },
      error: () => {}
    });
  }

  private readReviewImage(file: File): void {
    if (!file.type.startsWith('image/')) {
      this.errorMessage = 'Only image files are allowed.';
      return;
    }
    const reader = new FileReader();
    reader.onload = () => {
      const data = typeof reader.result === 'string' ? reader.result : '';
      if (!data) {
        this.errorMessage = 'Could not read image.';
        return;
      }
      this.reviewForm.images = [...this.reviewForm.images, data];
    };
    reader.onerror = () => {
      this.errorMessage = 'Could not read image.';
    };
    reader.readAsDataURL(file);
  }
}
