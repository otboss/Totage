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
  private tag: string | null = null

  getSelectedLocation(){
    return this.selectedLocation()
  }

  setSelectedLocation(selectedLocation: coords){
    this.selectedLocation.set(selectedLocation)
  }

  getTag(): string {
    return this.tag
  }

  setTag(tag: string) {
    this.tag = tag
  }
}
