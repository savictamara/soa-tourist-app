import { Injectable } from '@angular/core';
import { CanActivate, Router, UrlTree } from '@angular/router';
import { AuthStateService } from '../services/auth-state.service';

@Injectable({
  providedIn: 'root'
})
export class AdminGuard implements CanActivate {
  constructor(
    private readonly authStateService: AuthStateService,
    private readonly router: Router
  ) {}

  canActivate(): boolean | UrlTree {
    if (!this.authStateService.isAuthenticated) {
      return this.router.parseUrl('/login');
    }

    return this.authStateService.currentUser?.role === 'Administrator'
      ? true
      : this.router.parseUrl('/profile');
  }
}
