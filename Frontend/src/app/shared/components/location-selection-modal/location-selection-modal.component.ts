import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { Router } from '@angular/router';
import { IonHeader, IonButtons, IonToolbar, IonModal, IonButton, IonTitle, IonContent, IonItem, IonInput, ModalController, IonIcon, IonBackButton, NavController } from "@ionic/angular/standalone";
import { addIcons } from 'ionicons';
import { closeOutline } from 'ionicons/icons';
import { Geolocation as CapacitorGeolocation } from '@capacitor/geolocation'
import { Capacitor } from '@capacitor/core';

@Component({
  selector: 'app-location-selection-modal',
  templateUrl: './location-selection-modal.component.html',
  styleUrls: ['./location-selection-modal.component.scss'],
  imports: [IonHeader, IonButtons, IonToolbar, IonButton, IonTitle, IonContent, IonItem, IonIcon, IonBackButton, IonModal],
  changeDetection: ChangeDetectionStrategy.Eager,
  standalone: true
})
export class LocationSelectionModalComponent  implements OnInit {

  constructor(private modalCtrl: ModalController) {
    addIcons({ closeOutline })
  }

  ngOnInit() {}

  cancel(){
    this.modalCtrl.dismiss(undefined)
  }

  confirm() {
    const data = { name: 'John Doe', role: 'Admin' }
    this.modalCtrl.dismiss(data, 'confirm')
  }

  async setCurrentUserCoordinates() {
    if (Capacitor.isPluginAvailable('Geolocation')) {
        const currentLocation = await CapacitorGeolocation.getCurrentPosition()
        // this.outputValues.update(currentValue => ({
        //   ...currentValue,
        //   pickupAddres: {
        //     ...currentValue.pickupAddress,
        //     latitude: currentLocation.coords.latitude,
        //     longitude: currentLocation.coords.longitude,
        //   }
        // }))
        return
    }

    navigator.geolocation.getCurrentPosition((currentLocation) => {
      // this.outputValues.update(currentValue => ({
      //   ...currentValue,
      //   pickupAddres: {
      //     ...currentValue.pickupAddress,
      //     latitude: currentLocation.coords.latitude,
      //     longitude: currentLocation.coords.longitude,
      //   }
      // }))
    })
  }

}
