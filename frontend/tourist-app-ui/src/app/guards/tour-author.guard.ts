import { Injectable } from '@angular/core';
import { CanActivate, Router, UrlTree } from '@angular/router';
import { AuthStateService } from '../services/auth-state.service';

@Injectable({
  providedIn: 'root'
})
export class TourAuthorGuard implements CanActivate {
  constructor(
    private readonly authStateService: AuthStateService,
    private readonly router: Router
  ) {}

  canActivate(): boolean | UrlTree {
    if (!this.authStateService.isAuthenticated) {
      return this.router.parseUrl('/login');
    }

    const role = this.authStateService.currentUser?.role;
    return role === 'Guide' || role === 'Author'
      ? true
      : this.router.parseUrl('/profile');
  }
}
