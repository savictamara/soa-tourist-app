import { Component } from '@angular/core';
import { StakeholdersApiService } from '../services/stakeholders-api.service';
import { RegisterRequest } from '../models/register-request.model';

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

  constructor(private readonly apiService: StakeholdersApiService) {}

  submit(): void {
    this.isSubmitting = true;
    this.successMessage = '';
    this.errorMessage = '';

    this.apiService.register(this.form).subscribe({
      next: (user) => {
        this.successMessage = `Registered ${user.username} as ${user.role}.`;
        this.form = {
          username: '',
          email: '',
          password: '',
          role: 'Tourist'
        };
        this.isSubmitting = false;
      },
      error: (error) => {
        this.errorMessage = error.error?.message ?? 'Registration failed.';
        this.isSubmitting = false;
      }
    });
  }
}
