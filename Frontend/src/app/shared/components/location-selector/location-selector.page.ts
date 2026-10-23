import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';

import { FormsModule } from '@angular/forms';
import { IonContent, IonHeader, IonTitle, IonToolbar, IonButtons, IonBackButton, IonIcon, NavController, IonRouterOutlet, IonButton } from '@ionic/angular/standalone';
import { AppRoutes } from 'src/constants/AppRoutes';
import * as L from 'leaflet';
import { addIcons } from 'ionicons';
import { location } from 'ionicons/icons';
import { LocationSelectorService } from './location-selector.service';

@Component({
  selector: 'app-location-selector',
  templateUrl: './location-selector.page.html',
  styleUrls: ['./location-selector.page.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonContent, IonHeader, IonTitle, IonToolbar, FormsModule, IonButtons, IonBackButton, IonIcon, IonButton]
})
export class LocationSelectorPage implements OnInit {
  private map!: L.Map
  readonly homePage: string = AppRoutes.Home()

  constructor(
    private locationSelectorService: LocationSelectorService,
    private navCtrl: NavController,
    private routerOutlet: IonRouterOutlet,
  ) {
    addIcons({location})
  }

  ngOnInit() {
    this.initMap();
  }

  ionViewDidEnter() {
    if (this.map) {
      setTimeout(() => {
        this.map.invalidateSize();
      }, 100)
    }
  }

  private initMap(): void {
    this.map = L.map('map').setView([this.locationSelectorService.getSelectedLocation()?.latitude ?? 0, this.locationSelectorService.getSelectedLocation()?.longitude ?? 0], this.locationSelectorService.getSelectedLocation() === null ? 2 : 18);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 19,
      attribution: '© OpenStreetMap'
    }).addTo(this.map)
  }

  selectLocation(){
    const center = this.map.getCenter()
    if(this.routerOutlet && this.routerOutlet.canGoBack()){
      this.navCtrl.pop()
    } else {
      this.navCtrl.navigateBack(AppRoutes.Dashboard())
    }
    setTimeout(() => {
      this.locationSelectorService.setSelectedLocation({
        latitude: center.lat,
        longitude: center.lng
      })
    }, 100)
  }
}
