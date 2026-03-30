import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { LoginRequest } from '../models/login-request.model';
import { AuthStateService } from '../services/auth-state.service';
import { StakeholdersApiService } from '../services/stakeholders-api.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
  form: LoginRequest = {
    username: '',
    password: ''
  };

  isSubmitting = false;
  errorMessage = '';

  constructor(
    private readonly apiService: StakeholdersApiService,
    private readonly authStateService: AuthStateService,
    private readonly router: Router
  ) {}

  submit(): void {
    this.isSubmitting = true;
    this.errorMessage = '';

    this.apiService.login(this.form).subscribe({
      next: (response) => {
        this.authStateService.setSession(response);
        this.isSubmitting = false;
        void this.router.navigateByUrl(response.user.role === 'Administrator' ? '/users' : '/profile');
      },
      error: (error) => {
        this.errorMessage = error.error?.message ?? 'Login failed.';
        this.isSubmitting = false;
      }
    });
  }
}
