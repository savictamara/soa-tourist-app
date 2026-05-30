import { Component, OnDestroy, OnInit } from '@angular/core';
import { forkJoin, of } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { KeyPoint } from '../models/key-point.model';
import { CreateReviewRequest, Review, Tour, TourExecution } from '../models/tour.model';
import { AuthStateService } from '../services/auth-state.service';
import { PositionSimulatorService } from '../services/position-simulator.service';
import { PurchaseApiService } from '../services/purchase-api.service';
import { TourApiService } from '../services/tour-api.service';

@Component({
  selector: 'app-tourist-tours',
  templateUrl: './tourist-tours.component.html',
  styleUrls: ['./tourist-tours.component.css']
})
export class TouristToursComponent implements OnInit, OnDestroy {
  tours: Tour[] = [];
  selectedTour: Tour | null = null;
  keyPoints: KeyPoint[] = [];
  reviews: Review[] = [];
  isLoadingTours = false;
  isLoadingDetails = false;
  isSubmittingReview = false;
  addingToCartTourId: string | null = null;
  purchasedTourIds = new Set<string>();
  activeExecution: TourExecution | null = null;
  latestExecution: TourExecution | null = null;
  executionStatusMessage = '';
  isStartingExecution = false;
  isCompletingExecution = false;
  isAbandoningExecution = false;
  successMessage = '';
  errorMessage = '';
  reviewImageDragActive = false;
  reviewForm = {
    rating: 5,
    comment: '',
    visitedDate: '',
    images: [] as string[]
  };

  private executionPollingId: number | null = null;

  constructor(
    private readonly tourApiService: TourApiService,
    private readonly purchaseApiService: PurchaseApiService,
    private readonly positionSimulatorService: PositionSimulatorService,
    private readonly authStateService: AuthStateService
  ) {}

  get todayString(): string {
    return new Date().toISOString().split('T')[0];
  }

  get touristId(): string {
    const user = this.authStateService.currentUser;
    return user?.id ? String(user.id) : user?.username?.trim() ?? '';
  }

  ngOnInit(): void {
    this.loadTours();
    this.loadActiveExecution();
  }

  ngOnDestroy(): void {
    this.stopExecutionPolling();
  }

  loadTours(): void {
    this.isLoadingTours = true;
    this.errorMessage = '';
    this.tourApiService.getAvailableTours().subscribe({
      next: (tours) => {
        this.tours = (tours ?? []).map(tour => ({
          ...tour,
          tags: Array.isArray(tour.tags) ? tour.tags : [],
          keyPoints: Array.isArray(tour.keyPoints) ? tour.keyPoints : [],
          reviews: Array.isArray(tour.reviews) ? tour.reviews : [],
          durations: Array.isArray(tour.durations) ? tour.durations : [],
          lengthKm: Number(tour.lengthKm ?? 0)
        }));
        this.loadPurchaseStatuses();
        this.isLoadingTours = false;
      },
      error: (error) => {
        this.errorMessage = error?.error?.error ?? 'Could not load tours.';
        this.isLoadingTours = false;
      }
    });
  }

  startTour(tour: Tour, event?: Event): void {
    event?.stopPropagation();
    if (!this.touristId) {
      this.errorMessage = 'Logged tourist is required.';
      return;
    }
    if (!this.isPurchased(tour.id)) {
      this.errorMessage = 'Tour must be purchased before it can be started.';
      return;
    }
    this.isStartingExecution = true;
    this.errorMessage = '';
    this.successMessage = '';
    this.positionSimulatorService.getPosition().subscribe({
      next: (position) => {
        if (!position) {
          this.errorMessage = 'Set your current position in Position simulator first.';
          this.isStartingExecution = false;
          return;
        }
        this.tourApiService.startTourExecution(tour.id, {
          touristId: this.touristId,
          latitude: position.latitude,
          longitude: position.longitude
        }).subscribe({
          next: (execution) => {
            this.activeExecution = this.normalizeExecution(execution);
            this.selectedTour = tour;
            this.loadAllKeyPointsAndReviews(tour.id);
            this.successMessage = 'Tour started.';
            this.executionStatusMessage = 'Active tour session is running.';
            this.isStartingExecution = false;
            this.startExecutionPolling();
            this.checkActiveExecutionLocation();
          },
          error: (error) => {
            this.errorMessage = error?.error?.error ?? 'Could not start tour.';
            this.isStartingExecution = false;
          }
        });
      },
      error: (error) => {
        this.errorMessage = error?.error?.message ?? 'Could not read current position.';
        this.isStartingExecution = false;
      }
    });
  }

  completeActiveExecution(): void {
    if (!this.activeExecution) {
      return;
    }
    this.isCompletingExecution = true;
    this.tourApiService.completeExecution(this.activeExecution.id).subscribe({
      next: (execution) => {
        this.activeExecution = this.normalizeExecution(execution);
        this.executionStatusMessage = 'Tour completed.';
        this.successMessage = 'Tour completed.';
        this.stopExecutionPolling();
        this.isCompletingExecution = false;
      },
      error: (error) => {
        this.errorMessage = error?.error?.error ?? 'Could not complete tour.';
        this.isCompletingExecution = false;
      }
    });
  }

  abandonActiveExecution(): void {
    if (!this.activeExecution) {
      return;
    }
    this.isAbandoningExecution = true;
    this.tourApiService.abandonExecution(this.activeExecution.id).subscribe({
      next: (execution) => {
        this.activeExecution = this.normalizeExecution(execution);
        this.executionStatusMessage = 'Tour abandoned.';
        this.successMessage = 'Tour abandoned.';
        this.stopExecutionPolling();
        this.isAbandoningExecution = false;
      },
      error: (error) => {
        this.errorMessage = error?.error?.error ?? 'Could not abandon tour.';
        this.isAbandoningExecution = false;
      }
    });
  }

  isExecutionForTour(tourId: string): boolean {
    return this.activeExecution?.tourId === tourId && this.activeExecution.status === 'active';
  }

  isKeyPointCompleted(keyPointId: string): boolean {
    return (this.activeExecution?.completedKeyPoints ?? []).some(item => item.keyPointId === keyPointId);
  }

  getStatusLabel(status: string): string {
    if (status === 'completed') {
      return 'Completed';
    }
    if (status === 'abandoned') {
      return 'Abandoned';
    }
    return 'Active';
  }

  selectTour(tourId: string): void {
    this.isLoadingDetails = true;
    this.errorMessage = '';
    this.latestExecution = null;
    const tour = this.tours.find(item => item.id === tourId) ?? null;
    if (!tour) {
      this.errorMessage = 'Could not load selected tour.';
      this.isLoadingDetails = false;
      return;
    }
    this.selectedTour = tour;
    this.keyPoints = (tour.keyPoints ?? []).slice(0, 1);
    this.loadLatestExecution(tourId);
    this.purchaseApiService.isPurchased(this.touristId, tourId).pipe(
      catchError(() => of({ purchased: false }))
    ).subscribe({
      next: (status) => {
        if (status.purchased) {
          this.purchasedTourIds.add(tourId);
          this.loadAllKeyPointsAndReviews(tourId);
        } else {
          this.loadReviews(tourId);
        }
      }
    });
  }

  addToCart(tour: Tour, event?: Event): void {
    event?.stopPropagation();
    if ((tour.price ?? 0) <= 0) {
      this.errorMessage = 'Tour price must be set by guide.';
      this.successMessage = '';
      return;
    }
    if (!this.touristId) {
      this.errorMessage = 'Logged tourist is required.';
      return;
    }
    this.addingToCartTourId = tour.id;
    this.errorMessage = '';
    this.successMessage = '';
    this.purchaseApiService.addItem(this.touristId, {
      tourId: tour.id,
      tourName: tour.name,
      price: tour.price
    }).subscribe({
      next: () => {
        this.successMessage = 'Tour added to shopping cart.';
        this.addingToCartTourId = null;
      },
      error: (error) => {
        this.errorMessage = error?.error?.error ?? 'Could not add tour to cart.';
        this.addingToCartTourId = null;
      }
    });
  }

  isPurchased(tourId: string): boolean {
    return this.purchasedTourIds.has(tourId);
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

  private loadAllKeyPointsAndReviews(tourId: string): void {
    this.tourApiService.getKeyPoints(tourId).subscribe({
      next: (keyPoints) => {
        this.keyPoints = keyPoints ?? [];
        this.loadReviews(tourId);
      },
      error: (error) => {
        this.errorMessage = error?.error?.error ?? 'Could not load key points.';
        this.isLoadingDetails = false;
      }
    });
  }

  private refreshSelectedTour(tourId: string): void {
    this.tourApiService.getAvailableTours().subscribe({
      next: (tours) => {
        this.tours = (tours ?? []).map(tour => ({
          ...tour,
          tags: Array.isArray(tour.tags) ? tour.tags : [],
          keyPoints: Array.isArray(tour.keyPoints) ? tour.keyPoints : [],
          reviews: Array.isArray(tour.reviews) ? tour.reviews : [],
          durations: Array.isArray(tour.durations) ? tour.durations : [],
          lengthKm: Number(tour.lengthKm ?? 0)
        }));
        this.loadPurchaseStatuses();
        this.selectedTour = this.tours.find(tour => tour.id === tourId) ?? this.selectedTour;
        if (this.purchasedTourIds.has(tourId)) {
          this.loadAllKeyPointsAndReviews(tourId);
        } else {
          this.keyPoints = (this.selectedTour?.keyPoints ?? []).slice(0, 1);
          this.loadReviews(tourId);
        }
      },
      error: () => {}
    });
  }

  private loadPurchaseStatuses(): void {
    if (!this.touristId || this.tours.length === 0) {
      this.purchasedTourIds = new Set<string>();
      return;
    }
    const checks = this.tours.map(tour =>
      this.purchaseApiService.isPurchased(this.touristId, tour.id).pipe(
        catchError(() => of({ purchased: false }))
      )
    );
    forkJoin(checks).subscribe(statuses => {
      const purchased = new Set<string>();
      statuses.forEach((status, index) => {
        if (status.purchased) {
          purchased.add(this.tours[index].id);
        }
      });
      this.purchasedTourIds = purchased;
    });
  }

  private loadActiveExecution(): void {
    if (!this.touristId) {
      return;
    }
    this.tourApiService.getActiveExecution(this.touristId).pipe(
      catchError(() => of(null))
    ).subscribe(execution => {
      this.activeExecution = execution ? this.normalizeExecution(execution) : null;
      if (this.activeExecution?.status === 'active') {
        this.executionStatusMessage = 'Active tour session is running.';
        this.startExecutionPolling();
      }
    });
  }

  private loadLatestExecution(tourId: string): void {
    if (!this.touristId) {
      return;
    }
    this.tourApiService.getLatestExecution(tourId, this.touristId).pipe(
      catchError(() => of(null))
    ).subscribe(execution => {
      this.latestExecution = execution ? this.normalizeExecution(execution) : null;
    });
  }

  private checkActiveExecutionLocation(): void {
    if (!this.activeExecution || this.activeExecution.status !== 'active') {
      return;
    }
    this.positionSimulatorService.getPosition().subscribe({
      next: (position) => {
        if (!position || !this.activeExecution) {
          this.executionStatusMessage = 'Set your current position in Position simulator to continue tracking.';
          return;
        }
        this.tourApiService.checkExecutionLocation(this.activeExecution.id, position).subscribe({
          next: (response) => {
            this.activeExecution = this.normalizeExecution(response.execution);
            this.executionStatusMessage = response.keyPointReached && response.reachedKeyPoint
              ? `Reached key point: ${response.reachedKeyPoint.keyPointName}.`
              : `Last activity updated. Progress: ${response.completedCount}/${response.totalKeyPointCount}.`;
          },
          error: (error) => {
            this.errorMessage = error?.error?.error ?? 'Could not check active tour location.';
          }
        });
      },
      error: (error) => {
        this.errorMessage = error?.error?.message ?? 'Could not read current position.';
      }
    });
  }

  private startExecutionPolling(): void {
    this.stopExecutionPolling();
    this.executionPollingId = window.setInterval(() => this.checkActiveExecutionLocation(), 10000);
  }

  private stopExecutionPolling(): void {
    if (this.executionPollingId != null) {
      window.clearInterval(this.executionPollingId);
      this.executionPollingId = null;
    }
  }

  private normalizeExecution(execution: TourExecution): TourExecution {
    return {
      ...execution,
      completedKeyPoints: Array.isArray(execution.completedKeyPoints) ? execution.completedKeyPoints : []
    };
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
