import { AfterViewInit, Component, OnDestroy, OnInit } from '@angular/core';
import { KeyPoint } from '../models/key-point.model';
import { CreateKeyPointRequest, CreateTourRequest, Tour, UpdateKeyPointRequest } from '../models/tour.model';
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

  editKeyPointForm = {
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
  isUpdatingKeyPoint = false;
  isDeletingKeyPointId: string | null = null;
  isLoadingTours = false;
  isLoadingKeyPoints = false;
  editingKeyPoint: KeyPoint | null = null;
  successMessage = '';
  errorMessage = '';

  keyPointImageDragActive = false;
  keyPointImageFileName = '';
  editKeyPointImageFileName = '';
  editKeyPointImageDragActive = false;

  private mapInstance: any;
  private mapMarker: L.Marker | null = null;
  private editMarker: L.Marker | null = null;
  private tourMarkers: L.Marker[] = [];
  private tourPolyline: L.Polyline | null = null;
  mapMode: 'add' | 'edit' = 'add';

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

    this.tourApiService.createTour(payload).pipe(
      finalize(() => { this.isCreatingTour = false; })
    ).subscribe({
      next: (tour) => {
        const normalizedTour = this.normalizeTour(tour);
        this.tours = [normalizedTour, ...this.tours.filter(existing => existing.id !== normalizedTour.id)];
        this.selectTour(normalizedTour.id);
        this.createTourForm = { ...this.createTourForm, name: '', description: '', difficulty: 'easy', tagInput: '' };
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
      return;
    }
    this.isLoadingTours = true;
    this.errorMessage = '';

    this.tourApiService.getToursByAuthor(authorId).subscribe({
      next: (tours) => {
        this.tours = (tours ?? []).map(tour => this.normalizeTour(tour));
        if (this.tours.length > 0) {
          const tourToSelect = preferredTourId && this.tours.some(t => t.id === preferredTourId)
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

  loadTours(preferredTourId?: string): void {
    this.isLoadingTours = true;
    this.errorMessage = '';
    this.tourApiService.getTours().subscribe({
      next: (tours) => {
        this.tours = (tours ?? []).map(tour => this.normalizeTour(tour));
        if (this.tours.length > 0) {
          const tourToSelect = preferredTourId && this.tours.some(t => t.id === preferredTourId)
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
    this.cancelEdit();
    this.resetSelectedMapPoint();
    this.loadKeyPoints();
    this.ensureMapReady();
  }

  addKeyPoint(): void {
    if (!this.canCreateTours) {
      this.errorMessage = 'Only Guide or Author can add key points.';
      return;
    }
    if (!this.selectedTourId) {
      this.errorMessage = 'Select a tour first.';
      return;
    }
    if (!this.keyPointForm.name.trim() || !this.keyPointForm.description.trim() ||
      this.keyPointForm.latitude === null || this.keyPointForm.longitude === null ||
      !this.keyPointForm.imageUrl.trim()) {
      this.errorMessage = 'Name, description, image, latitude and longitude are required.';
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
      next: () => {
        this.keyPointForm = { name: '', description: '', latitude: null, longitude: null, imageUrl: '' };
        this.keyPointImageFileName = '';
        this.resetSelectedMapPoint();
        this.successMessage = 'Key point added successfully.';
        this.isAddingKeyPoint = false;
        this.loadToursByAuthor(this.getAuthorId(), this.selectedTourId);
      },
      error: (error) => {
        this.errorMessage = error.error?.error ?? 'Could not add key point.';
        this.isAddingKeyPoint = false;
      }
    });
  }

  startEditKeyPoint(kp: KeyPoint): void {
    this.editingKeyPoint = kp;
    this.mapMode = 'edit';
    this.editKeyPointForm = {
      name: kp.name,
      description: kp.description,
      latitude: kp.latitude,
      longitude: kp.longitude,
      imageUrl: kp.imageUrl ?? ''
    };
    this.editKeyPointImageFileName = '';
    this.errorMessage = '';
    this.successMessage = '';

    if (this.mapInstance) {
      if (this.editMarker) {
        this.mapInstance.removeLayer(this.editMarker);
      }
      this.editMarker = L.marker([kp.latitude, kp.longitude], { icon: this.defaultMarkerIcon }).addTo(this.mapInstance);
      this.mapInstance.panTo([kp.latitude, kp.longitude]);
    }
  }

  cancelEdit(): void {
    this.editingKeyPoint = null;
    this.mapMode = 'add';
    if (this.editMarker && this.mapInstance) {
      this.mapInstance.removeLayer(this.editMarker);
      this.editMarker = null;
    }
  }

  saveEditKeyPoint(): void {
    if (!this.editingKeyPoint || !this.selectedTourId) return;
    if (!this.editKeyPointForm.name.trim() || !this.editKeyPointForm.description.trim() ||
      this.editKeyPointForm.latitude === null || this.editKeyPointForm.longitude === null ||
      !this.editKeyPointForm.imageUrl.trim()) {
      this.errorMessage = 'All fields are required.';
      return;
    }

    this.isUpdatingKeyPoint = true;
    this.errorMessage = '';
    this.successMessage = '';

    const request: UpdateKeyPointRequest = {
      name: this.editKeyPointForm.name.trim(),
      description: this.editKeyPointForm.description.trim(),
      latitude: this.editKeyPointForm.latitude,
      longitude: this.editKeyPointForm.longitude,
      imageUrl: this.editKeyPointForm.imageUrl.trim()
    };

    this.tourApiService.updateKeyPoint(this.selectedTourId, this.editingKeyPoint.id, request).subscribe({
      next: () => {
        this.successMessage = 'Key point updated.';
        this.cancelEdit();
        this.loadKeyPoints();
        this.isUpdatingKeyPoint = false;
      },
      error: (error) => {
        this.errorMessage = error.error?.error ?? 'Could not update key point.';
        this.isUpdatingKeyPoint = false;
      }
    });
  }

  deleteKeyPoint(keyPointId: string): void {
    if (!this.selectedTourId) return;
    if (!confirm('Are you sure you want to delete this key point?')) return;
    this.isDeletingKeyPointId = keyPointId;
    this.errorMessage = '';
    this.successMessage = '';

    this.tourApiService.deleteKeyPoint(this.selectedTourId, keyPointId).subscribe({
      next: () => {
        this.successMessage = 'Key point deleted.';
        if (this.editingKeyPoint?.id === keyPointId) {
          this.cancelEdit();
        }
        this.isDeletingKeyPointId = null;
        this.loadKeyPoints();
      },
      error: (error) => {
        this.errorMessage = error.error?.error ?? 'Could not delete key point.';
        this.isDeletingKeyPointId = null;
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
        this.selectedTourKeyPoints = (keyPoints ?? []).map(kp => this.normalizeKeyPoint(kp));
        this.isLoadingKeyPoints = false;
        this.renderTourOnMap();
      },
      error: (error) => {
        this.errorMessage = error.error?.error ?? 'Could not load key points.';
        this.isLoadingKeyPoints = false;
      }
    });
  }

  private renderTourOnMap(): void {
    if (!this.mapInstance) return;
    this.clearTourLayers();

    const sorted = [...this.selectedTourKeyPoints].sort((a, b) => a.order - b.order);
    if (sorted.length === 0) return;

    const latlngs: L.LatLngTuple[] = sorted.map(kp => [kp.latitude, kp.longitude]);
    this.tourPolyline = L.polyline(latlngs, { color: '#ef7fa8', weight: 4, opacity: 0.8 }).addTo(this.mapInstance);

    sorted.forEach(kp => {
      const icon = L.divIcon({
        className: '',
        html: `<div class="kp-marker">${kp.order}</div>`,
        iconSize: [28, 28],
        iconAnchor: [14, 14],
        popupAnchor: [0, -16]
      });
      const marker = L.marker([kp.latitude, kp.longitude], { icon })
        .addTo(this.mapInstance)
        .bindPopup(`<strong>${kp.name}</strong><br>${kp.description}`);
      this.tourMarkers.push(marker);
    });

    if (latlngs.length === 1) {
      this.mapInstance.setView(latlngs[0], 13);
    } else {
      this.mapInstance.fitBounds(L.latLngBounds(latlngs), { padding: [40, 40] });
    }
  }

  private clearTourLayers(): void {
    this.tourMarkers.forEach(m => this.mapInstance?.removeLayer(m));
    this.tourMarkers = [];
    if (this.tourPolyline) {
      this.mapInstance?.removeLayer(this.tourPolyline);
      this.tourPolyline = null;
    }
  }

  private async initializeMap(): Promise<void> {
    const mapHost = document.getElementById('tour-key-point-map');
    if (!mapHost) return;

    if (this.mapInstance) {
      this.mapInstance.invalidateSize();
      return;
    }

    this.mapInstance = L.map('tour-key-point-map').setView([45.2671, 19.8335], 11);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(this.mapInstance);

    this.mapInstance.on('click', (e: any) => {
      const lat = Number(e.latlng.lat.toFixed(6));
      const lng = Number(e.latlng.lng.toFixed(6));

      if (this.mapMode === 'edit') {
        this.editKeyPointForm.latitude = lat;
        this.editKeyPointForm.longitude = lng;
        if (this.editMarker) {
          this.editMarker.setLatLng([lat, lng]);
        } else {
          this.editMarker = L.marker([lat, lng], { icon: this.defaultMarkerIcon }).addTo(this.mapInstance);
        }
      } else {
        this.keyPointForm.latitude = lat;
        this.keyPointForm.longitude = lng;
        if (this.mapMarker) {
          this.mapMarker.setLatLng([lat, lng]);
        } else {
          this.mapMarker = L.marker([lat, lng], { icon: this.defaultMarkerIcon }).addTo(this.mapInstance);
        }
      }
    });

    this.renderTourOnMap();
  }

  private ensureMapReady(): void {
    if (!this.selectedTour) return;
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

  addTagFromInput(event: Event): void {
    event.preventDefault();
    this.pushPendingTag();
  }

  removeTag(tag: string): void {
    this.tags = this.tags.filter(t => t !== tag);
  }

  private pushPendingTag(): void {
    const value = this.createTourForm.tagInput.trim();
    if (!value) return;
    if (!this.tags.includes(value)) {
      this.tags = [...this.tags, value];
    }
    this.createTourForm.tagInput = '';
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
    return { ...keyPoint };
  }

  // ── key point image (add form) ──────────────────────────

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
    if (file) this.readImageFile(file, 'add');
  }

  onKeyPointImageSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    const file = input.files?.item(0);
    if (!file) return;
    this.readImageFile(file, 'add');
    input.value = '';
  }

  removeKeyPointImage(): void {
    this.keyPointForm.imageUrl = '';
    this.keyPointImageFileName = '';
  }

  // ── key point image (edit form) ─────────────────────────

  onEditKeyPointImageDragOver(event: DragEvent): void {
    event.preventDefault();
    this.editKeyPointImageDragActive = true;
  }

  onEditKeyPointImageDragLeave(event: DragEvent): void {
    event.preventDefault();
    this.editKeyPointImageDragActive = false;
  }

  onEditKeyPointImageDrop(event: DragEvent): void {
    event.preventDefault();
    this.editKeyPointImageDragActive = false;
    const file = event.dataTransfer?.files?.item(0);
    if (file) this.readImageFile(file, 'edit');
  }

  onEditKeyPointImageSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    const file = input.files?.item(0);
    if (!file) return;
    this.readImageFile(file, 'edit');
    input.value = '';
  }

  removeEditKeyPointImage(): void {
    this.editKeyPointForm.imageUrl = '';
    this.editKeyPointImageFileName = '';
  }

  private readImageFile(file: File, target: 'add' | 'edit'): void {
    if (!file.type.startsWith('image/')) {
      this.errorMessage = 'Only image files are allowed.';
      return;
    }
    const reader = new FileReader();
    reader.onload = () => {
      const src = typeof reader.result === 'string' ? reader.result : '';
      if (!src) { this.errorMessage = 'Could not read the selected image.'; return; }
      if (target === 'edit') {
        this.editKeyPointForm.imageUrl = src;
        this.editKeyPointImageFileName = file.name;
      } else {
        this.keyPointForm.imageUrl = src;
        this.keyPointImageFileName = file.name;
      }
      this.errorMessage = '';
    };
    reader.onerror = () => { this.errorMessage = 'Could not read the selected image.'; };
    reader.readAsDataURL(file);
  }
}
