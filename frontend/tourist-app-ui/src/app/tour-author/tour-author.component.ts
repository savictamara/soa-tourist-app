import { AfterViewInit, Component, OnDestroy, OnInit } from '@angular/core';
import { KeyPoint } from '../models/key-point.model';
import { CreateKeyPointRequest, CreateReviewRequest, CreateTourRequest, Review, Tour } from '../models/tour.model';
import { AuthStateService } from '../services/auth-state.service';
import { TourApiService } from '../services/tour-api.service';
import { finalize } from 'rxjs';
import * as L from 'leaflet';

@Component({
  selector: 'app-tour-author',
  templateUrl: './tour-author.component.html',
  styleUrls: ['./tour-author.component.css']
})
export class TourAuthorComponent implements OnInit, AfterViewInit, OnDestroy {
  createTourForm = {
    name: '',
    description: '',
    difficulty: 'easy',
    tagInput: ''
  };

  keyPointForm = {
    name: '',
    description: '',
    latitude: null as number | null,
    longitude: null as number | null,
    imageUrl: ''
  };

  tours: Tour[] = [];
  selectedTourId = '';
  selectedTour: Tour | null = null;
  selectedTourKeyPoints: KeyPoint[] = [];
  tags: string[] = [];
  isCreatingTour = false;
  isAddingKeyPoint = false;
  isLoadingTours = false;
  isLoadingKeyPoints = false;
  successMessage = '';
  errorMessage = '';
  private mapInstance: any;
  private mapMarker: L.Marker | null = null;
  keyPointImageDragActive = false;
  keyPointImageFileName = '';
  reviewForm = {
    rating: 5,
    comment: '',
    visitedDate: '',
    images: [] as string[]
  };
  reviews: Review[] = [];
  isLoadingReviews = false;
  isAddingReview = false;
  reviewImageDragActive = false;
  private readonly defaultMarkerIcon = L.icon({
    iconRetinaUrl: 'assets/leaflet/marker-icon-2x.png',
    iconUrl: 'assets/leaflet/marker-icon.png',
    shadowUrl: 'assets/leaflet/marker-shadow.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    tooltipAnchor: [16, -28],
    shadowSize: [41, 41]
  });

  constructor(
    private readonly tourApiService: TourApiService,
    private readonly authStateService: AuthStateService
  ) {}

  ngOnInit(): void {
    if (this.authStateService.currentUser?.role === 'Tourist') {
      this.loadTours();
      return;
    }
    this.loadToursByAuthor(this.getAuthorId());
  }

  async ngAfterViewInit(): Promise<void> {
    await this.initializeMap();
  }

  ngOnDestroy(): void {
    if (this.mapInstance) {
      this.mapInstance.remove();
      this.mapInstance = null;
      this.mapMarker = null;
    }
  }

  get currentAuthorId(): string {
    return this.getAuthorId();
  }

  private getAuthorId(): string {
    const currentUser = this.authStateService.currentUser;
    if (currentUser?.username?.trim()) {
      return currentUser.username.trim();
    }
    return 'guide-1';
  }

  get canCreateTours(): boolean {
    const role = this.authStateService.currentUser?.role;
    return role === 'Guide' || role === 'Author';
  }

  createTour(): void {
    console.log('Create tour clicked');

    if (!this.canCreateTours) {
      this.errorMessage = 'Only Guide or Author can create tours.';
      this.successMessage = '';
      this.isCreatingTour = false;
      return;
    }

    this.pushPendingTag();

    const authorId = this.getAuthorId();
    if (!authorId ||
      !this.createTourForm.name.trim() ||
      !this.createTourForm.description.trim() ||
      !this.createTourForm.difficulty.trim() ||
      this.tags.length === 0) {
      this.errorMessage = 'Name, description, difficulty and at least one tag are required.';
      this.successMessage = '';
      this.isCreatingTour = false;
      return;
    }

    this.isCreatingTour = true;
    this.errorMessage = '';
    this.successMessage = '';
    const payload: CreateTourRequest = {
      authorId,
      name: this.createTourForm.name.trim(),
      description: this.createTourForm.description.trim(),
      difficulty: this.createTourForm.difficulty.trim(),
      tags: [...this.tags]
    };
    console.log('Create tour payload', payload);

    this.tourApiService.createTour(payload).pipe(
      finalize(() => {
        this.isCreatingTour = false;
      })
    ).subscribe({
      next: (tour) => {
        const normalizedTour = this.normalizeTour(tour);
        this.tours = [normalizedTour, ...this.tours.filter(existing => existing.id !== normalizedTour.id)];
        this.selectTour(normalizedTour.id);
        this.createTourForm = {
          ...this.createTourForm,
          name: '',
          description: '',
          difficulty: 'easy',
          tagInput: ''
        };
        this.tags = [];
        this.successMessage = 'Tour created successfully.';
        this.loadToursByAuthor(authorId, normalizedTour.id);
      },
      error: (error) => {
        const statusSuffix = error?.status ? ` (HTTP ${error.status})` : '';
        this.errorMessage = `${error?.error?.error ?? 'Could not create tour.'}${statusSuffix}`;
      }
    });
  }

  loadToursByAuthor(authorId: string, preferredTourId?: string): void {
    if (!authorId) {
      this.errorMessage = 'Author ID is required to fetch tours.';
      this.successMessage = '';
      return;
    }

    this.isLoadingTours = true;
    this.errorMessage = '';

    this.tourApiService.getToursByAuthor(authorId).subscribe({
      next: (tours) => {
        this.tours = (tours ?? []).map(tour => this.normalizeTour(tour));
        if (this.tours.length > 0) {
          const tourToSelect = preferredTourId && this.tours.some(tour => tour.id === preferredTourId)
            ? preferredTourId
            : this.tours[0].id;
          this.selectTour(tourToSelect);
        } else {
          this.selectedTour = null;
          this.selectedTourId = '';
          this.selectedTourKeyPoints = [];
        }
        this.isLoadingTours = false;
      },
      error: (error) => {
        this.errorMessage = error.error?.error ?? 'Could not load tours.';
        this.isLoadingTours = false;
      }
    });
  }

  selectTour(tourId: string): void {
    this.selectedTourId = tourId;
    const foundTour = this.tours.find(tour => tour.id === tourId) ?? null;
    this.selectedTour = foundTour ? this.normalizeTour(foundTour) : null;
    this.resetSelectedMapPoint();
    this.loadKeyPoints();
    this.loadReviews();
    this.ensureMapReady();
  }

  loadTours(preferredTourId?: string): void {
    this.isLoadingTours = true;
    this.errorMessage = '';
    this.tourApiService.getTours().subscribe({
      next: (tours) => {
        this.tours = (tours ?? []).map(tour => this.normalizeTour(tour));
        if (this.tours.length > 0) {
          const tourToSelect = preferredTourId && this.tours.some(tour => tour.id === preferredTourId)
            ? preferredTourId
            : this.tours[0].id;
          this.selectTour(tourToSelect);
        } else {
          this.selectedTour = null;
          this.selectedTourId = '';
          this.selectedTourKeyPoints = [];
        }
        this.isLoadingTours = false;
      },
      error: (error) => {
        this.errorMessage = error.error?.error ?? 'Could not load tours.';
        this.isLoadingTours = false;
      }
    });
  }

  addKeyPoint(): void {
    if (!this.canCreateTours) {
      this.errorMessage = 'Only Guide or Author can add key points.';
      return;
    }
    if (!this.selectedTourId) {
      this.errorMessage = 'Select a tour first.';
      this.successMessage = '';
      return;
    }
    if (!this.keyPointForm.name.trim() ||
      !this.keyPointForm.description.trim() ||
      this.keyPointForm.latitude === null ||
      this.keyPointForm.longitude === null ||
      !this.keyPointForm.imageUrl.trim()) {
      this.errorMessage = 'Name, description, image, latitude and longitude are required.';
      this.successMessage = '';
      return;
    }
    if (this.keyPointForm.latitude < -90 || this.keyPointForm.latitude > 90) {
      this.errorMessage = 'Latitude must be between -90 and 90.';
      this.successMessage = '';
      return;
    }
    if (this.keyPointForm.longitude < -180 || this.keyPointForm.longitude > 180) {
      this.errorMessage = 'Longitude must be between -180 and 180.';
      this.successMessage = '';
      return;
    }

    this.isAddingKeyPoint = true;
    this.errorMessage = '';
    this.successMessage = '';

    const request: CreateKeyPointRequest = {
      name: this.keyPointForm.name.trim(),
      description: this.keyPointForm.description.trim(),
      latitude: this.keyPointForm.latitude,
      longitude: this.keyPointForm.longitude,
      imageUrl: this.keyPointForm.imageUrl.trim()
    };

    this.tourApiService.addKeyPoint(this.selectedTourId, request).subscribe({
      next: (keyPoint) => {
        this.selectedTourKeyPoints = [...this.selectedTourKeyPoints, this.normalizeKeyPoint(keyPoint)];
        this.keyPointForm = {
          name: '',
          description: '',
          latitude: null,
          longitude: null,
          imageUrl: ''
        };
        this.keyPointImageFileName = '';
        this.resetSelectedMapPoint();
        this.successMessage = 'Key point added successfully.';
        this.loadKeyPoints();
        this.isAddingKeyPoint = false;
      },
      error: (error) => {
        this.errorMessage = error.error?.error ?? 'Could not add key point.';
        this.isAddingKeyPoint = false;
      }
    });
  }

  private loadKeyPoints(): void {
    if (!this.selectedTourId) {
      this.selectedTourKeyPoints = [];
      return;
    }

    this.isLoadingKeyPoints = true;
    this.tourApiService.getKeyPoints(this.selectedTourId).subscribe({
      next: (keyPoints) => {
        this.selectedTourKeyPoints = (keyPoints ?? []).map(keyPoint => this.normalizeKeyPoint(keyPoint));
        this.isLoadingKeyPoints = false;
      },
      error: (error) => {
        this.errorMessage = error.error?.error ?? 'Could not load key points.';
        this.isLoadingKeyPoints = false;
      }
    });
  }

  addTagFromInput(event: Event): void {
    event.preventDefault();
    this.pushPendingTag();
  }

  removeTag(tag: string): void {
    this.tags = this.tags.filter(existingTag => existingTag !== tag);
  }

  private normalizeTour(tour: Tour): Tour {
    return {
      ...tour,
      tags: Array.isArray(tour.tags) ? tour.tags : [],
      keyPoints: Array.isArray(tour.keyPoints) ? tour.keyPoints : [],
      reviews: Array.isArray(tour.reviews) ? tour.reviews : []
    };
  }

  private normalizeKeyPoint(keyPoint: KeyPoint): KeyPoint {
    return {
      ...keyPoint
    };
  }

  private pushPendingTag(): void {
    const value = this.createTourForm.tagInput.trim();
    if (!value) {
      return;
    }
    if (!this.tags.includes(value)) {
      this.tags = [...this.tags, value];
    }
    this.createTourForm.tagInput = '';
  }

  private async initializeMap(): Promise<void> {
    const mapHost = document.getElementById('tour-key-point-map');
    if (!mapHost) {
      return;
    }

    if (this.mapInstance) {
      this.mapInstance.invalidateSize();
      return;
    }

    this.mapInstance = L.map('tour-key-point-map').setView([45.2671, 19.8335], 11);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(this.mapInstance);

    this.mapInstance.on('click', (e: any) => {
      const latitude = Number(e.latlng.lat.toFixed(6));
      const longitude = Number(e.latlng.lng.toFixed(6));

      this.keyPointForm.latitude = latitude;
      this.keyPointForm.longitude = longitude;

      if (this.mapMarker) {
        this.mapMarker.setLatLng([latitude, longitude]);
        return;
      }

      this.mapMarker = L.marker([latitude, longitude], { icon: this.defaultMarkerIcon }).addTo(this.mapInstance);
    });
  }

  private ensureMapReady(): void {
    if (!this.selectedTour) {
      return;
    }
    setTimeout(async () => {
      await this.initializeMap();
      if (this.mapInstance) {
        this.mapInstance.invalidateSize();
      }
    }, 0);
  }

  private resetSelectedMapPoint(): void {
    this.keyPointForm.latitude = null;
    this.keyPointForm.longitude = null;
    if (this.mapMarker && this.mapInstance) {
      this.mapInstance.removeLayer(this.mapMarker);
      this.mapMarker = null;
    }
  }

  onKeyPointImageDragOver(event: DragEvent): void {
    event.preventDefault();
    this.keyPointImageDragActive = true;
  }

  onKeyPointImageDragLeave(event: DragEvent): void {
    event.preventDefault();
    this.keyPointImageDragActive = false;
  }

  onKeyPointImageDrop(event: DragEvent): void {
    event.preventDefault();
    this.keyPointImageDragActive = false;
    const file = event.dataTransfer?.files?.item(0);
    if (file) {
      this.readKeyPointImageFile(file);
    }
  }

  onKeyPointImageSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    const file = input.files?.item(0);
    if (!file) {
      return;
    }
    this.readKeyPointImageFile(file);
    input.value = '';
  }

  removeKeyPointImage(): void {
    this.keyPointForm.imageUrl = '';
    this.keyPointImageFileName = '';
  }

  private readKeyPointImageFile(file: File): void {
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
      this.keyPointForm.imageUrl = imageSource;
      this.keyPointImageFileName = file.name;
      this.errorMessage = '';
    };
    reader.onerror = () => {
      this.errorMessage = 'Could not read the selected image.';
    };
    reader.readAsDataURL(file);
  }

  private loadReviews(): void {
    if (!this.selectedTourId) {
      this.reviews = [];
      return;
    }
    this.isLoadingReviews = true;
    this.tourApiService.getReviews(this.selectedTourId).subscribe({
      next: (reviews) => {
        this.reviews = reviews ?? [];
        this.isLoadingReviews = false;
      },
      error: (error) => {
        this.errorMessage = error.error?.error ?? 'Could not load reviews.';
        this.isLoadingReviews = false;
      }
    });
  }

  addReview(): void {
    if (!this.selectedTourId) {
      this.errorMessage = 'Select a tour first.';
      return;
    }
    const currentUser = this.authStateService.currentUser;
    const touristUsername = currentUser?.username?.trim() ?? '';
    const touristId = currentUser?.id ? String(currentUser.id) : touristUsername;
    if (!touristUsername || !touristId) {
      this.errorMessage = 'Logged user is required to add review.';
      return;
    }
    if (!this.reviewForm.comment.trim() || !this.reviewForm.visitedDate) {
      this.errorMessage = 'Comment and visited date are required.';
      return;
    }
    if (this.reviewForm.rating < 1 || this.reviewForm.rating > 5) {
      this.errorMessage = 'Rating must be between 1 and 5.';
      return;
    }

    this.isAddingReview = true;
    this.errorMessage = '';
    this.successMessage = '';
    const payload: CreateReviewRequest = {
      rating: this.reviewForm.rating,
      comment: this.reviewForm.comment.trim(),
      visitedDate: this.reviewForm.visitedDate,
      touristId,
      touristUsername,
      images: [...this.reviewForm.images]
    };

    this.tourApiService.addReview(this.selectedTourId, payload).subscribe({
      next: (review) => {
        this.reviews = [review, ...this.reviews];
        this.reviewForm = { rating: 5, comment: '', visitedDate: '', images: [] };
        this.successMessage = 'Review added successfully.';
        this.isAddingReview = false;
      },
      error: (error) => {
        this.errorMessage = error.error?.error ?? 'Could not add review.';
        this.isAddingReview = false;
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
        this.readReviewImageFile(file);
      }
    }
  }

  onReviewImagesSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    const files = input.files;
    if (!files) {
      return;
    }
    for (let i = 0; i < files.length; i += 1) {
      const file = files.item(i);
      if (file) {
        this.readReviewImageFile(file);
      }
    }
    input.value = '';
  }

  removeReviewImage(index: number): void {
    this.reviewForm.images = this.reviewForm.images.filter((_, imageIndex) => imageIndex !== index);
  }

  private readReviewImageFile(file: File): void {
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
      this.reviewForm.images = [...this.reviewForm.images, imageSource];
    };
    reader.onerror = () => {
      this.errorMessage = 'Could not read the selected image.';
    };
    reader.readAsDataURL(file);
  }
}
