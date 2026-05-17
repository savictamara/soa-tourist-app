import { AfterViewInit, Component, OnDestroy } from '@angular/core';
import * as L from 'leaflet';
import { User } from '../models/user.model';
import { TouristPosition } from '../models/tourist-position.model';
import { AuthStateService } from '../services/auth-state.service';
import { PositionSimulatorService } from '../services/position-simulator.service';

@Component({
  selector: 'app-position-simulator',
  templateUrl: './position-simulator.component.html',
  styleUrls: ['./position-simulator.component.css']
})
export class PositionSimulatorComponent implements AfterViewInit, OnDestroy {
  currentUser: User | null;
  currentPosition: TouristPosition | null = null;
  statusMessage = '';
  errorMessage = '';
  isSaving = false;

  private mapInstance: L.Map | null = null;
  private mapMarker: L.Marker | null = null;

  private readonly defaultMarkerIcon = L.icon({
    iconRetinaUrl: 'assets/leaflet/marker-icon-2x.png',
    iconUrl: 'assets/leaflet/marker-icon.png',
    shadowUrl: 'assets/leaflet/marker-shadow.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    shadowSize: [41, 41]
  });

  constructor(
    private readonly authStateService: AuthStateService,
    private readonly positionSimulatorService: PositionSimulatorService
  ) {
    this.currentUser = this.authStateService.currentUser;
  }

  ngAfterViewInit(): void {
    this.initializeMap();
    this.loadStoredPosition();
  }

  ngOnDestroy(): void {
    this.mapInstance?.remove();
    this.mapInstance = null;
  }

  clearPosition(): void {
    if (!this.currentUser || this.isSaving) {
      return;
    }

    this.isSaving = true;
    this.errorMessage = '';
    this.positionSimulatorService.clearPosition().subscribe({
      next: () => {
        this.currentPosition = null;
        this.statusMessage = 'Your saved position has been cleared.';
        if (this.mapMarker && this.mapInstance) {
          this.mapInstance.removeLayer(this.mapMarker);
          this.mapMarker = null;
        }
        this.isSaving = false;
      },
      error: (error) => {
        this.errorMessage = error?.error?.message ?? 'Could not clear the tourist position.';
        this.isSaving = false;
      }
    });
  }

  private initializeMap(): void {
    if (this.mapInstance) {
      this.mapInstance.remove();
    }

    const initialCoordinates: L.LatLngTuple = this.currentPosition
      ? [this.currentPosition.latitude, this.currentPosition.longitude]
      : [45.2671, 19.8335];
    const initialZoom = this.currentPosition ? 13 : 11;

    this.mapInstance = L.map('tourist-position-map').setView(initialCoordinates, initialZoom);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(this.mapInstance);

    if (this.currentPosition) {
      this.renderMarker(initialCoordinates);
    }

    this.mapInstance.on('click', (event: L.LeafletMouseEvent) => {
      if (!this.currentUser || this.isSaving) {
        return;
      }

      const latitude = Number(event.latlng.lat.toFixed(6));
      const longitude = Number(event.latlng.lng.toFixed(6));
      this.isSaving = true;
      this.errorMessage = '';
      this.positionSimulatorService.setPosition(latitude, longitude).subscribe({
        next: (position) => {
          this.currentPosition = position;
          this.statusMessage = 'Your current position has been updated.';
          this.renderMarker([position.latitude, position.longitude]);
          this.isSaving = false;
        },
        error: (error) => {
          this.errorMessage = error?.error?.message ?? 'Could not save the tourist position.';
          this.isSaving = false;
        }
      });
    });

    setTimeout(() => this.mapInstance?.invalidateSize(), 0);
  }

  private renderMarker(coordinates: L.LatLngTuple): void {
    if (!this.mapInstance) {
      return;
    }

    if (this.mapMarker) {
      this.mapMarker.setLatLng(coordinates);
    } else {
      this.mapMarker = L.marker(coordinates, { icon: this.defaultMarkerIcon }).addTo(this.mapInstance);
    }

    this.mapMarker.bindPopup('Current position').openPopup();
    this.mapInstance.panTo(coordinates);
  }

  private loadStoredPosition(): void {
    if (!this.currentUser) {
      return;
    }

    this.positionSimulatorService.getPosition().subscribe({
      next: (position) => {
        this.currentPosition = position;
        if (position) {
          this.statusMessage = 'Your last saved position is shown on the map.';
          this.renderMarker([position.latitude, position.longitude]);
        } else {
          this.statusMessage = 'Click on the map to choose your position.';
        }
      },
      error: (error) => {
        this.errorMessage = error?.error?.message ?? 'Could not load the tourist position.';
      }
    });
  }
}
