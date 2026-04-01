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
    role: 'Tourist'
  };

  isSubmitting = false;
  successMessage = '';
  errorMessage = '';

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
          role: 'Tourist'
        };
        this.isSubmitting = false;
        void this.router.navigateByUrl('/profile');
      },
      error: (error) => {
        this.errorMessage = error.error?.message ?? 'Registration failed.';
        this.isSubmitting = false;
      }
    });
  }
}
