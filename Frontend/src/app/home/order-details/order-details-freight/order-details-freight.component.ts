import { Component, computed, CUSTOM_ELEMENTS_SCHEMA, effect, ElementRef, input, OnInit, output, signal, ViewChild, ChangeDetectionStrategy } from '@angular/core';
import { IonGrid, IonRow, IonCard, IonInput, IonItem, IonButton, IonIcon, IonSegment, IonSegmentButton, IonLabel, RadioGroupCustomEvent, ToastController } from "@ionic/angular/standalone";
import { addIcons } from 'ionicons';
import { locationOutline } from 'ionicons/icons';
import { LocationSelectionModalComponent } from 'src/app/shared/components/location-selection-modal/location-selection-modal.component';
import { FreightTruckClass, FreightTruckProperty } from 'src/app/shared/services/client/shared/api/freight_truck_class_pb';
import { Assets } from 'src/constants/Assets';
import { OrderDetailsFreightService } from './order-details-freight.service';
import { GeocodingQuerySchema, GetFreightTruckClassesNearbyQuerySchema, GetFreightTruckClassesNearbyResponse, GetFreightTrucksNearbyQuery, GetFreightTrucksNearbyQuerySchema } from 'src/app/shared/services/client/modules/logistic/api/logistic_pb';
import { OrderDetailsFreightValidator } from './order-details-freight.validator';
import { ɵInternalFormsSharedModule } from '@angular/forms';
import { create } from '@bufbuild/protobuf';
import { Capacitor } from '@capacitor/core'
import { Geolocation as CapacitorGeolocation } from '@capacitor/geolocation'
import { CompilationError, createValidator, RuntimeError, ValidationError } from "@bufbuild/protovalidate";
import { debouncedSignal } from 'src/utils/debouncedSignal';
import { Router } from '@angular/router';
import { AppRoutes } from 'src/constants/AppRoutes';
import { LocationSelectorService } from 'src/app/shared/components/location-selector/location-selector.service';

@Component({
  selector: 'app-order-details-freight',
  templateUrl: './order-details-freight.component.html',
  styleUrls: ['./order-details-freight.component.scss'],
  standalone: true,
  imports: [IonGrid, IonRow, IonCard, IonInput, IonItem, IonButton, IonIcon, LocationSelectionModalComponent, IonSegment, IonSegmentButton, IonLabel, ɵInternalFormsSharedModule],
  changeDetection: ChangeDetectionStrategy.Eager,
  schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class OrderDetailsFreightComponent  implements OnInit {  
  @ViewChild('form') form: ElementRef
  
  isSubmitAttempted = input<boolean>(false)
  readonly freightTruckProperty = FreightTruckProperty
  
  subcategory: 'freight' | 'solid-waste' | 'liquid-waste' = 'freight'

  private outputValues = signal<GetFreightTrucksNearbyQuery>(create(GetFreightTrucksNearbyQuerySchema))
  
  output = output<GetFreightTrucksNearbyQuery>()
  errors = output<Partial<Record<keyof typeof this.validator.type, string> & (ValidationError | CompilationError | RuntimeError)>>({})

  updateForm: typeof this.validator.updateForm = (x, y) => this.validator.updateForm(x, y) 
  formErrors = this.validator.formErrors
  showFieldErrors = this.validator.showFieldErrors
  formValues = this.validator.formValues
  pickupAddressSelectedOnMap: boolean = false

  freightTruckClasses: GetFreightTruckClassesNearbyResponse["freightTruckClasses"] = undefined

  selectedVehicleClassIndex: number = 0
  vehicleClassSlides: () => ReadonlyArray<{
    image: string
    label: string
    type: "OPEN" | "CLOSED"
    class: FreightTruckClass
  }> = () => [
    ...(this.subcategory === 'freight' ? [{
      label: "Van",
      image: Assets.Images.van,
      type: "CLOSED",
      class: FreightTruckClass.VAN
    }] : [] as any),
    {
      label: this.subcategory === 'liquid-waste' ? "Small Waste Truck" : "Open Truck Small",
      image: this.subcategory === 'liquid-waste' ? Assets.Images.truckOpenLiquidWasteSmall : Assets.Images.truckOpenCargoSmall,
      type: "OPEN",
      class: FreightTruckClass.OPEN_TRUCK_SMALL
    },
    ...(this.subcategory === 'freight' ? [{
      label: "Box Truck Small",
      image: Assets.Images.truckCargoSmail,
      type: "CLOSED",
      class: FreightTruckClass.BOX_TRUCK_SMALL
    }] : [] as any),
    {
      label: this.subcategory === 'liquid-waste' ? "Medium Waste Truck" : "Open Truck Large",
      image: this.subcategory === 'liquid-waste' ? Assets.Images.truckOpenLiquidWasteLarge : Assets.Images.truckOpenCargoLarge,
      type: "OPEN",
      class: FreightTruckClass.OPEN_TRUCK_LARGE
    },
    ...(this.subcategory === 'freight' ? [{
      label: "Box Truck Large",
      image: Assets.Images.truckCargoLarge,
      type: "CLOSED",
      class: FreightTruckClass.BOX_TRUCK_LARGE
    }] : [] as any),
    ...(this.subcategory === 'liquid-waste' ? [{
      label: "Large Waste Truck",
      image: Assets.Images.cespoolTruck,
      type: "CLOSED",
      class: FreightTruckClass.OPEN_TRUCK_LARGE
    }] : [] as any),
  ]

  isFetchingFreightTruckClasses: boolean = false

  
  constructor(
    private orderDetailsFreightService: OrderDetailsFreightService,
    private validator: OrderDetailsFreightValidator, 
    private toastController: ToastController,
    private router: Router,
    private locationSelectorService: LocationSelectorService,
  ) { 

    const protoValidationResult = createValidator().validate(GetFreightTrucksNearbyQuerySchema, this.outputValues())

    addIcons({ locationOutline })

    effect(() => {
      if(this.isSubmitAttempted()){
        this.debouncedValidation()
      }
    })

    effect(() => {
      const selectedLocation = this.locationSelectorService.getSelectedLocation()
      if(selectedLocation != null){
        this.validator.updateForm('pickupAddress', JSON.stringify(`${selectedLocation.latitude},${selectedLocation.longitude}`))
      }
    })

    effect(() => {
      // TODO: explore using protovalion instead of fluentvalidation for form validation

      const protoValidationResult = createValidator().validate(GetFreightTrucksNearbyQuerySchema, this.outputValues())
      
      this.errors.emit({
          ...protoValidationResult.error,
          ...this.validator.formErrors()
      })

      if (Object.keys(this.validator.formErrors()).length === 0 && protoValidationResult.kind  === "valid"){
        this.output.emit(this.outputValues())
      }
    })
  }

  ngOnInit() {
  }


  setSubcategory(val: any){
    this.subcategory = val as any
  }

  private async debouncedValidation(){
      if (this.validator.formValues().pickupAddress === "") return
      try{
        await this.geocodingPickupAddress()
        await this.getFreightTruckClassesNearby()
      }
      catch(err){
        const toast = await this.toastController.create({
          message: 'Please check pickup address or internet connection',
          duration: 2500,
          position: "bottom",
        })
        await toast.present()
      }
  }

  setSelectedVehicleClass(event: any){
    this.selectedVehicleClassIndex = event.detail[0].activeIndex
    if(!this.freightTruckClasses){
      return
    }
    this.outputValues.update(currentValue => ({
      ...currentValue,
      freightTruckClass: this.freightTruckClasses.find(vehicleClass => {
        return vehicleClass.class === this.vehicleClassSlides()[this.selectedVehicleClassIndex].class
      })
    }))
  }

  private async geocodingPickupAddress(){
     const geocodingResult = await this.orderDetailsFreightService.geocodingQuery(create(GeocodingQuerySchema, {
        address: this.validator.formValues().pickupAddress
      }))
      this.outputValues.update(currentValue => ({
        ...currentValue,
        pickupAddress: {
          ...currentValue.pickupAddress,
          latitude: geocodingResult.coordinates.latitude,
          longitude: geocodingResult.coordinates.longitude,
        },
        freightTruckClass: {
          ...currentValue.freightTruckClass,
          property: null
        }
      }))
  }

  private async getFreightTruckClassesNearby(){
    try{
      this.isFetchingFreightTruckClasses = true

      const response = await this.orderDetailsFreightService.getFreightTruckClassesNearby(
        create(GetFreightTruckClassesNearbyQuerySchema, {
          pickupAddress: this.outputValues().pickupAddress,
        })
      )

      this.freightTruckClasses = response.freightTruckClasses

      if(this.freightTruckClasses.length > 0){
        this.outputValues.update(currentValue => ({
          ...currentValue,
          freightTruckClass: this.freightTruckClasses[0]
        }))
      }
    }
    finally {
      this.isFetchingFreightTruckClasses = false
    }
  }

  async setCurrentUserCoordinates() {
    if (Capacitor.isPluginAvailable('Geolocation')) {
        const currentLocation = await CapacitorGeolocation.getCurrentPosition()
        this.outputValues.update(currentValue => ({
          ...currentValue,
          pickupAddres: {
            ...currentValue.pickupAddress,
            latitude: currentLocation.coords.latitude,
            longitude: currentLocation.coords.longitude,
          }
        }))
        return
    }

    navigator.geolocation.getCurrentPosition((currentLocation) => {
      this.outputValues.update(currentValue => ({
        ...currentValue,
        pickupAddres: {
          ...currentValue.pickupAddress,
          latitude: currentLocation.coords.latitude,
          longitude: currentLocation.coords.longitude,
        }
      }))
    })
  }

  isGeolocationSupported(): boolean{
    return Capacitor.isPluginAvailable('Geolocation') || (navigator.geolocation !== undefined)
  }

  handleTruckPropertyChange(event: RadioGroupCustomEvent<FreightTruckProperty>){
    this.outputValues.update((currentValue) => ({
      ...currentValue,
      freightTruckClass: {
        ...currentValue.freightTruckClass,
        property: event.detail.value,
      }
    }))
  }

  getSupportedVehicleClasses() {
  return this.vehicleClassSlides()
    // switch(this.subcategory){
    //   case 'freight':
    //     return this.vehicleClassSlides
    //   case 'solid-waste':
    //     return this.vehicleClassSlides.filter(x => x.class === FreightTruckClass.OPEN_TRUCK_SMALL || x.class === FreightTruckClass.OPEN_TRUCK_LARGE)
    //   case 'liquid-waste':
    //     return this.vehicleClassSlides.filter(x => x.class === FreightTruckClass.OPEN_TRUCK_SMALL || x.class === FreightTruckClass.OPEN_TRUCK_LARGE)
    // }
    // return this.vehicleClassSlides.filter(x => {
    //   return this.freightTruckClasses === undefined ? true : this.freightTruckClasses.filter(y => y.class === x.class)
    // })
  }

  getSupportedVehicleProperties(freightTruckClass: FreightTruckClass): Array<FreightTruckProperty>{
    return this.freightTruckClasses === undefined 
      ? Object.values(FreightTruckProperty) as Array<FreightTruckProperty>
      : this.freightTruckClasses.filter(x => x.class === freightTruckClass).map(x => x.property)
  }

  navigateToLocationSelector(){
    this.router.navigateByUrl(AppRoutes.LocationSelector(), {
      skipLocationChange: true
    })
  }

}
