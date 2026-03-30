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
  isLoading = false;
  errorMessage = '';

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

  ngOnInit(): void {
    if (this.isLoggedIn) {
      this.loadProfile();
    }
  }

  loadProfile(): void {
    if (!this.authStateService.isAuthenticated) {
      this.errorMessage = 'Authentication is required.';
      this.selectedProfile = null;
      return;
    }

    this.isLoading = true;
    this.errorMessage = '';
    this.selectedProfile = null;

    this.apiService.getMyProfile().subscribe({
      next: (profile) => {
        this.selectedProfile = profile;
        this.isLoading = false;
      },
      error: (error) => {
        this.errorMessage = error.error?.message ?? 'Could not load profile.';
        this.isLoading = false;
      }
    });
  }
}
