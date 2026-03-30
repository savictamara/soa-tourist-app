import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { StakeholdersApiService } from '../services/stakeholders-api.service';
import { RegisterRequest } from '../models/register-request.model';
import { AuthStateService } from '../services/auth-state.service';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent {
  form: RegisterRequest = {
    username: '',
    email: '',
    password: '',
    role: 'Tourist',
    firstName: '',
    lastName: '',
    profileImage: '',
    biography: '',
    motto: ''
  };

  isSubmitting = false;
  successMessage = '';
  errorMessage = '';
  imageDragActive = false;

  constructor(
    private readonly apiService: StakeholdersApiService,
    private readonly authStateService: AuthStateService,
    private readonly router: Router
  ) {}

  submit(): void {
    this.isSubmitting = true;
    this.successMessage = '';
    this.errorMessage = '';

    this.apiService.register(this.form).subscribe({
      next: (response) => {
        this.successMessage = `Registered ${response.user.username} as ${response.user.role}.`;
        this.authStateService.setSession(response);
        this.form = {
          username: '',
          email: '',
          password: '',
          role: 'Tourist',
          firstName: '',
          lastName: '',
          profileImage: '',
          biography: '',
          motto: ''
        };
        this.imageDragActive = false;
        this.isSubmitting = false;
        void this.router.navigateByUrl('/profile');
      },
      error: (error) => {
        this.errorMessage = error.error?.message ?? 'Registration failed.';
        this.isSubmitting = false;
      }
    });
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
    this.form.profileImage = '';
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
      this.form.profileImage = canvas.toDataURL('image/jpeg', 0.72);
    };

    image.onerror = () => {
      this.errorMessage = 'Could not process the selected image.';
    };

    image.src = imageSource;
  }
}
