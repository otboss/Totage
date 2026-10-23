import { Injectable, signal } from '@angular/core';


type coords = {
  latitude: number,
  longitude: number
}

@Injectable({
  providedIn: 'root',
})
export class LocationSelectorService {
  private selectedLocation = signal<coords>(null)

  getSelectedLocation(){
    return this.selectedLocation()
  }

  setSelectedLocation(selectedLocation: coords){
    this.selectedLocation.set(selectedLocation)
  }
}
