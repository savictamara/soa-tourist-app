import { Component, OnInit } from '@angular/core';
import { StakeholdersApiService } from '../services/stakeholders-api.service';
import { User } from '../models/user.model';

@Component({
  selector: 'app-users',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.css']
})
export class UsersComponent implements OnInit {
  users: User[] = [];
  isLoading = false;
  errorMessage = '';
  adminUsername = '';

  constructor(private readonly apiService: StakeholdersApiService) {}

  ngOnInit(): void {}

  loadUsers(): void {
    if (!this.adminUsername.trim()) {
      this.errorMessage = 'Admin username is required.';
      this.users = [];
      return;
    }

    this.isLoading = true;
    this.errorMessage = '';

    this.apiService.getUsers(this.adminUsername.trim()).subscribe({
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
}
