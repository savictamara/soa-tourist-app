import { Component, OnInit } from '@angular/core';
import { StakeholdersApiService } from '../services/stakeholders-api.service';
import { UserProfile } from '../models/user-profile.model';
import { AuthStateService } from '../services/auth-state.service';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {
  selectedProfile: UserProfile | null = null;
  editableProfile: UserProfile = {};
  isEditing = false;
  imageDragActive = false;
  isLoading = false;
  errorMessage = '';
  successMessage = '';

  constructor(
    private readonly apiService: StakeholdersApiService,
    private readonly authStateService: AuthStateService
  ) {}

  get currentUsername(): string {
    return this.authStateService.currentUser?.username ?? '';
  }

  get isLoggedIn(): boolean {
    return this.authStateService.currentUser !== null;
  }

  get isAdmin(): boolean {
    return this.authStateService.currentUser?.role === 'Administrator';
  }

  ngOnInit(): void {
    if (this.isLoggedIn) {
      this.loadProfile();
    }
  }

  loadProfile(): void {
    if (!this.authStateService.isAuthenticated) {
      this.errorMessage = 'Authentication is required.';
      this.successMessage = '';
      this.selectedProfile = null;
      return;
    }

    this.isLoading = true;
    this.errorMessage = '';
    this.successMessage = '';
    this.isEditing = false;
    this.selectedProfile = null;

    this.apiService.getMyProfile().subscribe({
      next: (profile) => {
        this.selectedProfile = profile;
        this.editableProfile = { ...profile };
        this.isLoading = false;
      },
      error: (error) => {
        this.errorMessage = error.error?.message ?? 'Could not load profile.';
        this.isLoading = false;
      }
    });
  }

  saveProfile(): void {
    if (this.isAdmin) {
      this.errorMessage = 'Administrator profile editing is not available.';
      this.successMessage = '';
      return;
    }

    if (!this.authStateService.isAuthenticated) {
      this.errorMessage = 'Authentication is required.';
      this.successMessage = '';
      return;
    }

    this.isLoading = true;
    this.errorMessage = '';
    this.successMessage = '';

    this.apiService.updateMyProfile(this.editableProfile).subscribe({
      next: (profile) => {
        this.selectedProfile = profile;
        this.editableProfile = { ...profile };
        this.isEditing = false;
        this.successMessage = 'Profile updated successfully.';
        this.isLoading = false;
      },
      error: (error) => {
        this.errorMessage = error.error?.message ?? 'Could not update profile.';
        this.isLoading = false;
      }
    });
  }

  startEdit(): void {
    if (!this.selectedProfile || this.isAdmin) {
      return;
    }

    this.editableProfile = { ...this.selectedProfile };
    this.isEditing = true;
    this.errorMessage = '';
    this.successMessage = '';
  }

  cancelEdit(): void {
    this.isEditing = false;
    this.imageDragActive = false;
    this.editableProfile = this.selectedProfile ? { ...this.selectedProfile } : {};
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

  removeImage(): void {
    this.editableProfile.profileImage = null;
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
      const maxSize = 320;
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
      this.editableProfile.profileImage = canvas.toDataURL('image/jpeg', 0.72);
    };

    image.onerror = () => {
      this.errorMessage = 'Could not process the selected image.';
    };

    image.src = imageSource;
  }
}
