import { Component, OnInit } from '@angular/core';
import { StakeholdersApiService } from '../services/stakeholders-api.service';
import { AuthStateService } from '../services/auth-state.service';
import { User } from '../models/user.model';

@Component({
  selector: 'app-users',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.css']
})
export class UsersComponent implements OnInit {
  users: User[] = [];
  isLoading = false;
  isBlocking = false;
  errorMessage = '';
  blockMessage = '';

  constructor(
    private readonly apiService: StakeholdersApiService,
    private readonly authStateService: AuthStateService
  ) {}

  get currentUser(): User | null {
    return this.authStateService.currentUser;
  }

  ngOnInit(): void {
    if (this.currentUser?.role === 'Administrator') {
      this.loadUsers();
    }
  }

  loadUsers(): void {
    if (this.currentUser?.role !== 'Administrator') {
      this.errorMessage = 'Administrator access is required.';
      this.users = [];
      return;
    }

    this.isLoading = true;
    this.errorMessage = '';
    this.blockMessage = '';

    this.apiService.getUsers().subscribe({
      next: (users) => {
        this.users = users;
        this.isLoading = false;
      },
      error: (error) => {
        this.errorMessage = error.error?.message ?? 'Could not load users.';
        this.isLoading = false;
      }
    });
  }

  blockUser(user: User): void {
    if (this.currentUser?.role !== 'Administrator') {
      this.errorMessage = 'Administrator access is required.';
      return;
    }

    this.isBlocking = true;
    this.errorMessage = '';
    this.blockMessage = '';

    this.apiService.blockUser(user.id).subscribe({
      next: (response) => {
        this.users = this.users.map(existingUser =>
          existingUser.id === user.id ? { ...existingUser, isBlocked: true } : existingUser
        );
        this.blockMessage = response.message;
        this.isBlocking = false;
      },
      error: (error) => {
        this.errorMessage = error.error?.message ?? 'Could not block user.';
        this.isBlocking = false;
      }
    });
  }
}
