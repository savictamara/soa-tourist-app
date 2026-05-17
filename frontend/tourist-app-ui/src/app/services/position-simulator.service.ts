import { Injectable } from '@angular/core';
import { Observable, map } from 'rxjs';
import { TouristPosition } from '../models/tourist-position.model';
import { StakeholdersApiService } from './stakeholders-api.service';

@Injectable({
  providedIn: 'root'
})
export class PositionSimulatorService {
  constructor(private readonly stakeholdersApiService: StakeholdersApiService) {}

  getPosition(): Observable<TouristPosition | null> {
    return this.stakeholdersApiService.getMyPosition().pipe(
      map(position => {
        if (position.latitude == null || position.longitude == null) {
          return null;
        }

        return {
          latitude: position.latitude,
          longitude: position.longitude
        };
      })
    );
  }

  setPosition(latitude: number, longitude: number): Observable<TouristPosition> {
    const request = {
      latitude,
      longitude
    };

    return this.stakeholdersApiService.updateMyPosition(request).pipe(
      map(position => ({
        latitude: position.latitude ?? latitude,
        longitude: position.longitude ?? longitude
      }))
    );
  }

  clearPosition(): Observable<void> {
    const request = {
      latitude: null,
      longitude: null
    };

    return this.stakeholdersApiService.updateMyPosition(request).pipe(
      map(() => void 0)
    );
  }
}
