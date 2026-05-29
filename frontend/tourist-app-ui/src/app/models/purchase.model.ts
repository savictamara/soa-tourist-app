export interface OrderItem {
  tourId: string;
  tourName: string;
  price: number;
}

export interface ShoppingCart {
  id: string;
  touristId: string;
  items: OrderItem[];
  totalPrice: number;
  createdAt: string;
  updatedAt: string;
}

export interface TourPurchaseToken {
  id: string;
  touristId: string;
  tourId: string;
  purchasedAt: string;
  token: string;
}

export interface PurchaseStatus {
  purchased: boolean;
}
