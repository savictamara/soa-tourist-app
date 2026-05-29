import { Component, OnInit } from '@angular/core';
import { ShoppingCart } from '../models/purchase.model';
import { AuthStateService } from '../services/auth-state.service';
import { PurchaseApiService } from '../services/purchase-api.service';

@Component({
  selector: 'app-shopping-cart',
  templateUrl: './shopping-cart.component.html',
  styleUrls: ['./shopping-cart.component.css']
})
export class ShoppingCartComponent implements OnInit {
  cart: ShoppingCart | null = null;
  isLoading = false;
  isCheckingOut = false;
  removingTourId: string | null = null;
  successMessage = '';
  errorMessage = '';

  constructor(
    private readonly purchaseApiService: PurchaseApiService,
    private readonly authStateService: AuthStateService
  ) {}

  ngOnInit(): void {
    this.loadCart();
  }

  get touristId(): string {
    const user = this.authStateService.currentUser;
    return user?.id ? String(user.id) : user?.username?.trim() ?? '';
  }

  loadCart(): void {
    if (!this.touristId) {
      this.errorMessage = 'Logged tourist is required.';
      return;
    }
    this.isLoading = true;
    this.errorMessage = '';
    this.purchaseApiService.getCart(this.touristId).subscribe({
      next: (cart) => {
        this.cart = this.normalizeCart(cart);
        this.isLoading = false;
      },
      error: (error) => {
        this.errorMessage = error?.error?.error ?? 'Could not load shopping cart.';
        this.isLoading = false;
      }
    });
  }

  removeItem(tourId: string): void {
    if (!this.touristId) return;
    this.removingTourId = tourId;
    this.errorMessage = '';
    this.successMessage = '';
    this.purchaseApiService.removeItem(this.touristId, tourId).subscribe({
      next: (cart) => {
        this.cart = this.normalizeCart(cart);
        this.successMessage = 'Item removed.';
        this.removingTourId = null;
      },
      error: (error) => {
        this.errorMessage = error?.error?.error ?? 'Could not remove item.';
        this.removingTourId = null;
      }
    });
  }

  checkout(): void {
    if (!this.touristId) return;
    this.isCheckingOut = true;
    this.errorMessage = '';
    this.successMessage = '';
    this.purchaseApiService.checkout(this.touristId).subscribe({
      next: (response) => {
        this.successMessage = `${response.tokens?.length ?? 0} tour purchase token(s) created.`;
        this.isCheckingOut = false;
        this.loadCart();
      },
      error: (error) => {
        this.errorMessage = error?.error?.error ?? 'Checkout failed.';
        this.isCheckingOut = false;
      }
    });
  }

  private normalizeCart(cart: ShoppingCart): ShoppingCart {
    return {
      ...cart,
      items: Array.isArray(cart.items) ? cart.items : [],
      totalPrice: Number(cart.totalPrice ?? 0)
    };
  }
}
