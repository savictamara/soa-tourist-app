import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { OrderItem, PurchaseStatus, ShoppingCart, TourPurchaseToken } from '../models/purchase.model';

@Injectable({
  providedIn: 'root'
})
export class PurchaseApiService {
  private readonly baseUrl = 'http://localhost:8080/api/purchases';

  constructor(private readonly http: HttpClient) {}

  getCart(touristId: string): Observable<ShoppingCart> {
    return this.http.get<ShoppingCart>(`${this.baseUrl}/cart/${touristId}`);
  }

  addItem(touristId: string, item: OrderItem): Observable<ShoppingCart> {
    return this.http.post<ShoppingCart>(`${this.baseUrl}/cart/${touristId}/items`, item);
  }

  removeItem(touristId: string, tourId: string): Observable<ShoppingCart> {
    return this.http.delete<ShoppingCart>(`${this.baseUrl}/cart/${touristId}/items/${tourId}`);
  }

  checkout(touristId: string): Observable<{ tokens: TourPurchaseToken[] }> {
    return this.http.post<{ tokens: TourPurchaseToken[] }>(`${this.baseUrl}/cart/${touristId}/checkout`, {});
  }

  getTokens(touristId: string): Observable<TourPurchaseToken[]> {
    return this.http.get<TourPurchaseToken[]>(`${this.baseUrl}/tokens/${touristId}`);
  }

  isPurchased(touristId: string, tourId: string): Observable<PurchaseStatus> {
    return this.http.get<PurchaseStatus>(`${this.baseUrl}/tokens/${touristId}/tour/${tourId}`);
  }
}
