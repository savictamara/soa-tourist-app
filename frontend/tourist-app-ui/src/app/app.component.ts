import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { User } from './models/user.model';
import { AuthStateService } from './services/auth-state.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'Tourist App';

  constructor(
    private readonly authStateService: AuthStateService,
    private readonly router: Router
  ) {}

  get currentUser(): User | null {
    return this.authStateService.currentUser;
  }

  clearCurrentUser(): void {
    this.authStateService.clear();
    void this.router.navigateByUrl('/login');
  }
}
