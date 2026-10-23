import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import * as L from 'leaflet';
import { AppRoutes } from 'src/constants/AppRoutes';
import { IonHeader, IonToolbar, IonButtons, IonBackButton, IonTitle, IonContent } from "@ionic/angular/standalone";

@Component({
  selector: 'app-location-selector',
  templateUrl: './location-selector.component.html',
  styleUrls: ['./location-selector.component.scss'],
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonHeader, IonToolbar, IonButtons, IonBackButton, IonTitle, IonContent],
})
export class LocationSelectorComponent  implements OnInit {
  private map!: L.Map
  homePage: string = AppRoutes.Home()

  constructor() { }

  ngAfterViewInit(): void {
    
  }

  ngOnInit() {
    this.initMap();
  }

  private initMap(): void {
    this.map = L.map('map').setView([0, 0], 2);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 19,
      attribution: '© OpenStreetMap'
    }).addTo(this.map);
  }
}
