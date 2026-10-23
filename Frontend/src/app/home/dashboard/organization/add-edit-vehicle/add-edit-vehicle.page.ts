import { Component, OnInit, ChangeDetectionStrategy, CUSTOM_ELEMENTS_SCHEMA, ViewChild, ElementRef } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { IonContent, IonHeader, IonTitle, IonToolbar, IonButtons, IonBackButton, IonGrid, IonRow, IonInput, IonButton, IonToggle, IonSelect, IonSelectOption, IonSegment, IonSegmentButton, IonCheckbox, IonModal, IonIcon, NavController, ToastController, IonRouterOutlet, IonCard, AlertController } from '@ionic/angular/standalone';
import { AppRoutes } from 'src/constants/AppRoutes';
import { CreateRentalVehicleCommandSchema, CreateTruckCommandSchema, CreateVehicleCommandSchema, CreateVehicleResponse, CreateVehicleResponseSchema, DeleteVehicleCommandSchema, FuelType, GetJobTimeslotsQuerySchema, GetRentalVehicleQuerySchema, GetTruckQuerySchema, GetVehicleQuerySchema, RentalVehicle, SetJobTimeslotCommandSchema, TogglePublicityRentalVehicleCommandSchema, TogglePublicityTruckTruckCommandSchema, UpdateVehicleDriverCommandSchema, Vehicle } from 'src/app/shared/services/client/modules/logistic/api/logistic_pb';
import { AddEditRentalVehicleValidator, AddEditTruckValidator, AddEditVehicleValidator, SetJobTimeslotValidator } from './add-edit-vehicle.validator';
import { create } from '@bufbuild/protobuf';
import { ActivatedRoute, Router } from '@angular/router';
import { AddEditVehicleService } from './add-edit-vehicle.service';
import { Truck, TruckCreate_BedSizeSchema } from 'src/app/shared/services/client/shared/api/truck_pb';
import { Account, GetDriversQuerySchema } from 'src/app/shared/services/client/modules/organization/api/organization_pb';
import { add, closeCircleOutline, trashBinOutline } from 'ionicons/icons';
import { addIcons } from 'ionicons';
import { FreightTruckProperty } from 'src/app/shared/services/client/shared/api/freight_truck_class_pb';
import { toTitleCase } from 'src/utils/toTitleCase';
import { GetEnrollmentQuerySchema, SetEnrollmentCommandSchema } from 'src/app/shared/services/client/modules/order/api/order_pb';
import { OrderType } from 'src/app/shared/services/client/shared/api/order_type_pb';
import { JobTimeslot, JobTimeslotSchema } from 'src/app/shared/services/client/shared/api/job_timeslot_pb';
import { OrderTypeCards } from 'src/models/OrderTypeCards';
import { Weekdays } from 'src/app/shared/services/client/shared/api/weekdays_pb';

@Component({
  selector: 'app-add-edit-vehicle',
  templateUrl: './add-edit-vehicle.page.html',
  styleUrls: ['./add-edit-vehicle.page.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonContent, IonHeader, IonTitle, IonToolbar, FormsModule, IonButtons, IonBackButton, IonGrid, IonRow, IonInput, IonButton, IonSelect, IonSelectOption, IonSegment, IonSegmentButton, IonCheckbox, IonToggle, IonModal, IonIcon, IonCard],
  schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class AddEditVehiclePage implements OnInit {
  readonly dashboardPage: string = AppRoutes.Dashboard()
  updateVehicleForm: typeof this.validatorVehicle.updateForm = (x, y) => this.validatorVehicle.updateForm(x, y)
  updateTruckForm: typeof this.validatorTruck.updateForm = (x, y) => this.validatorTruck.updateForm(x, y)
  updateRentalVehicleForm: typeof this.validatorRentalVehicle.updateForm = (x, y) => this.validatorRentalVehicle.updateForm(x, y)
  updateJobTimeslotForm<K extends keyof JobTimeslot & string>(
    index: number, 
    field: K, 
    event: Event,
    extraValue?: string
  ) {
    const input = event.target as HTMLInputElement
    const value = input.value

    const path = `jobTimeslots.${index}.${field}` as any

    if(field === 'weekdays'){
      const weekdays = new Set(this.validatorJobTimeslot.formValues().jobTimeslots[index].weekdays)
      weekdays.delete(Number(extraValue) as Weekdays)
      if((event as any).detail.checked){
        weekdays.add(Number(extraValue) as Weekdays)
      }
      this.validatorJobTimeslot.updateForm(path, Array.from(weekdays))
      return
    }

    this.validatorJobTimeslot.updateForm(path, !isNaN(Number(value)) ? Number(value) : value)
  }

  @ViewChild('swiper', { static: false }) swiperRef!: ElementRef

  vehicle: Vehicle | null = null
  organizationId: bigint | null = null

  rentalVehicle: RentalVehicle | null = null
  truck: Truck | null = null
  drivers: Array<Account> | null = null
  FuelType = FuelType
  Weekdays = Weekdays
  FreightTruckProperty = FreightTruckProperty

  formErrorsVehicle = this.validatorVehicle.formErrors
  formErrorsTruck = this.validatorTruck.formErrors
  formErrorsRentalVehicle = this.validatorRentalVehicle.formErrors
  formErrorsJobTimeslot = this.validatorJobTimeslot.formErrors

  showFieldErrorsVehicle = this.validatorVehicle.showFieldErrors
  showFieldErrorsTruck = this.validatorTruck.showFieldErrors
  showFieldErrorsRentalVehicle = this.validatorRentalVehicle.showFieldErrors
  showFieldErrorsJobTimeslot = this.validatorJobTimeslot.showFieldErrors

  toTitleCase = toTitleCase
  selectedEnrollments: Partial<Record<OrderType, true>> = {}

  isSubmitAttempted: boolean = false
  isSubmitting: boolean = false
  isSelectDriverModalOpen: boolean = false
  vehicleCategory: "commercial" | "private" = "commercial"

  compatibleEnrollment: Array<OrderType> | null = null

  OrderTypeCards = OrderTypeCards
  standardCommercialCards = Object.values(OrderTypeCards).filter(x => !(<Partial<Record<OrderType, true>>>{[OrderType.EGG]: true, [OrderType.VEHICLE_RENTAL]: true, [OrderType.CESSPOOL]: true, [OrderType.CONSTRUCTION]: true})[x.orderType])

  readonly fuelTypeOptions = Object.keys(FuelType).filter(x => isNaN(parseInt(x))).slice(1)
  readonly truckPropertiesOptions = () => Object.keys(FreightTruckProperty).filter(x => isNaN(parseInt(x))).slice(1).filter(x =>
    !!this.validatorTruck.formValues().truck.bedSize.heightMm === [FreightTruckProperty[FreightTruckProperty.FREEZER], FreightTruckProperty[FreightTruckProperty.REFRIGERATED], FreightTruckProperty[FreightTruckProperty.STANDARD_BOX]].includes(x)
  )

  constructor(
    private validatorVehicle: AddEditVehicleValidator,
    private validatorTruck: AddEditTruckValidator,
    private validatorRentalVehicle: AddEditRentalVehicleValidator,
    private validatorJobTimeslot: SetJobTimeslotValidator,
    private addEditVehicleService: AddEditVehicleService,
    private router: Router,
    private route: ActivatedRoute,
    private navCtrl: NavController,
    private toastController: ToastController,
    private routerOutlet: IonRouterOutlet,
    private alertController: AlertController
  ) {
    addIcons({ trashBinOutline, add, closeCircleOutline })
  }

  async ngOnInit() {
    try{
      const navigation = this.router.currentNavigation();
      this.vehicle = <Vehicle>navigation?.extras?.state?.vehicle ?? null
      this.organizationId = <bigint>navigation?.extras?.state?.organizationId ?? null
      const vehicleId = this.route.snapshot.paramMap.get('vehicleId')

      if(this.vehicle !== null){
        this.addEditVehicleService.getEnrollment(
          create(GetEnrollmentQuerySchema, {
            vehicleId: this.vehicle.vehicleId
          })
        ).then(resp => {
          resp.orderTypes.forEach(orderType => {
            this.selectedEnrollments[orderType] = true
          })
        })
        this.addEditVehicleService.getJobTimeslots(
          create(GetJobTimeslotsQuerySchema, {
            vehicleId: this.vehicle.vehicleId
          })
        ).then(resp => {
          this.validatorJobTimeslot.updateForm("jobTimeslots", resp.jobTimeslot)
        })
      } else {
        this.validatorJobTimeslot.updateForm("jobTimeslots", [create(JobTimeslotSchema)])
      }

      this.isSubmitAttempted = false
      this.validatorVehicle.reset()
      this.validatorRentalVehicle.reset()
      this.validatorTruck.reset()
      this.validatorJobTimeslot.reset()
      this.validatorVehicle.updateForm('organizationId', BigInt(this.route.snapshot.paramMap.get('organizationId')))
      

      this.addEditVehicleService.getDrivers(
        create(GetDriversQuerySchema, {
          organizationId: BigInt(this.validatorVehicle.formValues().organizationId),
        })
      )
      .then(getDriversResponse => this.drivers = getDriversResponse.drivers)
      .catch(err => { /** TODO: Handle error */})

      if(vehicleId !== null && this.vehicle === null){
        const vehicleResponse = await this.addEditVehicleService.getVehicle(
          create(GetVehicleQuerySchema, {
            vehicleId: BigInt(vehicleId)
          })
        )
        this.vehicle = vehicleResponse.vehicle
      }

      if(this.vehicle === null){
        return
      }
      this.updateVehicleForm('driverAccountId', this.vehicle.driverAccountId)
      this.updateVehicleForm('chassisNumber', this.vehicle.chassisNumber)
      this.updateVehicleForm('licensePlateNumber', this.vehicle.licensePlateNumber)
      this.updateVehicleForm('make', this.vehicle.make)
      this.updateVehicleForm('model', this.vehicle.model)
      this.updateVehicleForm('fuelType', this.vehicle.fuelType)
      this.updateVehicleForm('year', this.vehicle.year)
      this.updateVehicleForm('seatingCapacity', this.vehicle.seatingCapacity)
      this.updateVehicleForm('isManualTransmission', this.vehicle.isManualTransmission)
      // this.updateVehicleForm('organizationId', this.vehicle.vehicleId)

      if(this.vehicle?.rentalVehicleId){
        this.addEditVehicleService.getRentalVehicle(
          create(GetRentalVehicleQuerySchema, {
            vehicleId: BigInt(vehicleId)
          })
        )
        .then(rentalVehicleResponse => {
          this.rentalVehicle = rentalVehicleResponse.rentalVehicle ?? null

          this.validatorRentalVehicle.formValues()
          this.validatorRentalVehicle.updateForm('rentalVehicle', {
            ...this.validatorRentalVehicle.formValues().rentalVehicle,
            ...rentalVehicleResponse.rentalVehicle,
          })
        })
        .catch(err => { /** TODO: Handle error */  })
      }

      if(this.vehicle?.truckId){
        this.addEditVehicleService.getTruck(
          create(GetTruckQuerySchema, {
            vehicleId: BigInt(vehicleId)
          })
        )
        .then(truckResponse => {
          this.truck = truckResponse.truck ?? null

          this.validatorTruck.updateForm('truck', {
            ...this.validatorTruck.formValues().truck,
            cost: this.truck.cost,
            freightPricePerKm: this.truck.freightPricePerKm,
            property: this.truck.property,
            isPublic: this.truck.isPublic,
            bedSize: {
              ...create(TruckCreate_BedSizeSchema),
              lengthMm: this.truck.bedSize.lengthMm,
              widthMm: this.truck.bedSize.widthMm,
              heightMm: this.truck.bedSize.heightMm
            }
          })
        })
        .catch(err => { /** TODO: Handle error */  })
      }
    }
    catch(err){
      console.error(err)
    }
  }

  async onSubmit(){
    this.isSubmitting = true
    this.isSubmitAttempted = true
    try{
      if(!this.validatorVehicle.isValid()){
        return
      }
      if(this.validatorRentalVehicle.formValues().rentalVehicle.isPublic && !this.validatorRentalVehicle.isValid()){
        return
      }
      if(this.validatorTruck.formValues().truck.isPublic && !this.validatorTruck.isValid()){
        return
      }
      if(!this.validatorJobTimeslot.isValid() && this.validatorJobTimeslot.formErrors().vehicleId === null){
        return
      }
      for(const [index, timeslot] of this.validatorJobTimeslot.formValues().jobTimeslots.entries()){
        for(const fieldError of Object.values(this.validatorJobTimeslot.formErrors().jobTimeslots[index])){
          if(fieldError != null){
            return
          }
        }
        if(((timeslot.endHH * 60 + timeslot.endMM) > (timeslot.startHH * 60 + timeslot.startMM)) === false){
          const alert = await this.alertController.create({
            header: 'Error',
            message: 'Schedule end times must be strictly after their start times',
            buttons: ['OK'],
          })
          await alert.present()
          return
        }
      }

      const createVehicleResponse: CreateVehicleResponse = this.vehicle?.vehicleId != null ? 
        create(CreateVehicleResponseSchema, {
          vehicleId: this.vehicle!.vehicleId
        }) :
        await this.addEditVehicleService.createVehicle(
          create(CreateVehicleCommandSchema, {
            ...this.validatorVehicle.formValues(),
            fuelType: this.validatorVehicle.formValues().fuelType,
          })
        )

      this.validatorJobTimeslot.updateForm('vehicleId', createVehicleResponse.vehicleId)

      if(this.vehicle !== null && this.vehicle?.driverAccountId !== this.validatorVehicle.formValues().driverAccountId){
        this.addEditVehicleService.updateVehicleDriver(
          create(UpdateVehicleDriverCommandSchema, {
            vehicleId: this.vehicle.vehicleId
          })
        )
      }

      if(this.validatorTruck.formValues().truck.isPublic){
        await this.addEditVehicleService.createTruck(
          create(CreateTruckCommandSchema, {
            vehicleId: createVehicleResponse.vehicleId,
            truck: this.validatorTruck.formValues().truck,
          })
        )
        this.addEditVehicleService.setEnrollment(
          create(SetEnrollmentCommandSchema, {
            vehicleId: createVehicleResponse.vehicleId,
            orderTypes: Object.keys(this.selectedEnrollments).map(x => parseInt(x)) as Array<OrderType>
          })
        )
      }

      await this.addEditVehicleService.toggleTruckPublicity(
        create(TogglePublicityTruckTruckCommandSchema, {
          vehicleId: createVehicleResponse.vehicleId,
          IsPublic: this.validatorTruck.formValues().truck.isPublic
        })
      )

      if(this.validatorRentalVehicle.formValues().rentalVehicle.isPublic){
        await this.addEditVehicleService.createRentalVehicle(
          create(CreateRentalVehicleCommandSchema, {
            vehicleId: createVehicleResponse.vehicleId,
            rentalVehicle: this.validatorRentalVehicle.formValues().rentalVehicle,
          })
        )
      }

      await this.addEditVehicleService.toggleRentalVehiclePublicity(
        create(TogglePublicityRentalVehicleCommandSchema, {
          vehicleId: createVehicleResponse.vehicleId,
          IsPublic: this.validatorRentalVehicle.formValues().rentalVehicle.isPublic
        })
      )

      await this.addEditVehicleService.setJobTimeslots(
        create(SetJobTimeslotCommandSchema, this.validatorJobTimeslot.formValues())
      )

      if(this.routerOutlet && this.routerOutlet.canGoBack()){
        this.navCtrl.pop()
      } else {
        this.navCtrl.navigateBack(AppRoutes.Dashboard())
      }

      const toast = await this.toastController.create({
        message: `Vehicle ${this.vehicle ? 'updated' : 'created'} successfully`,
        duration: 2500,
        position: "bottom",
      })
      await toast.present()
    }
    catch(err){
      console.error(err)
      const toast = await this.toastController.create({
        message: 'Check internet conneection',
        duration: 2500,
        position: "bottom",
      })
      await toast.present()
    }
    finally{
      this.isSubmitting = false
    }
  }

  setVehicleCategory(category: any){
    this.vehicleCategory = category.toString()
  }

  BigInt = (val: string) => {
    if(val === ""){
      return null
    }
    return BigInt(val)
  }

  parseInt(val: string){
    if(val === ""){
      return null
    }
    return parseInt(val)
  }

  parseFloat(val: string){
    if(val === ""){
      return null
    }
    return parseFloat(val)
  }

  getIsPrivateEnabled(): boolean{
    return this.validatorRentalVehicle.formValues().rentalVehicle?.isPublic ?? false
  }

  getIsCommercialEnabled(): boolean{
    return this.validatorTruck.formValues().truck?.isPublic ?? false
  }

  getCommercialProperty(){
    return this.validatorTruck.formValues().truck.property
  }

  toggleSelectDriverModal(isOpen: boolean){
    this.isSelectDriverModalOpen = isOpen;
  }

  selectDriver(accountId: string){
    this.updateVehicleForm('driverAccountId', accountId)
    this.isSelectDriverModalOpen = false
  }

  getTruckHeight(){
    return this.validatorTruck.formValues().truck?.bedSize?.heightMm ?? 0 
  }

  getDriverAccountId(){
    return this.vehicle?.driverAccountId ?? this.validatorVehicle.formValues().driverAccountId ?? ''
  }

  getRentalVehicleDriverRequired(): boolean{
    return this.validatorRentalVehicle.formValues().rentalVehicle.isDriverMandatory
  }

  getFuelType(val: any){
    return FuelType[parseInt(val)] ?? null
  }

  getFuelTypeNumberFromFuelName(val: string): number{
    return FuelType[val as any] as any
  }

  getTruckProperty(val: any){
    return FreightTruckProperty[parseInt(val)] ?? null
  }

  getTruckPropertyFromPropertyName(val: string): number{
    return FreightTruckProperty[val as any] as any
  }

  goToNextSlide(isSectionValid: boolean){
    if(!isSectionValid){
      return
    }
    const swiperEl = this.swiperRef.nativeElement
    if (swiperEl && swiperEl.swiper) {
      swiperEl.swiper.slideNext()
      this.isSubmitAttempted = false
    }
  }

  goToPreviousSlide(){
    const swiperEl = this.swiperRef.nativeElement
    if (swiperEl && swiperEl.swiper) {
      swiperEl.swiper.slidePrev()
    }
  }

  isVehicleFormValid(){
    return this.validatorVehicle.isValid()
  }

  isCommercialFormValid(){
    return this.validatorTruck.isValid()
  }

  isPrivateVehicleFormValid(){
    return this.validatorRentalVehicle.isValid()
  }
  
  async deleteVehicle(){
    const alert = await this.alertController.create({
      header: 'Remove',
      message: 'Are you sure you want to remove this vehicle?',
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          cssClass: 'secondary',
          handler: () => {

          }
        },
        {
          text: 'Ok',
          handler: () => {
            this.addEditVehicleService.deleteVehicle(
              create(DeleteVehicleCommandSchema, {
                vehicleId: this.vehicle.vehicleId,
                organizationId: this.organizationId,
              })
            ).then(async res => {
              if(this.routerOutlet && this.routerOutlet.canGoBack()){
                this.navCtrl.pop()
              } else {
                this.navCtrl.navigateBack(AppRoutes.Dashboard())
              }
              const toast = await this.toastController.create({
                message: `Vehicle removed successfully`,
                duration: 2500,
                position: "bottom",
              })
              await toast.present()
            }).catch(async err => {
              const toast = await this.toastController.create({
                message: `Error while removing vehicle`,
                duration: 2500,
                position: "bottom",
              })
              await toast.present()
            })
          }
        }
      ]
    });

    await alert.present();
  }

  getJobTimeslots(){
    return this.validatorJobTimeslot.formValues().jobTimeslots
  }

  removeJobTimeslot(index: number){
    const jobTimeslot = this.validatorJobTimeslot.formValues().jobTimeslots
    jobTimeslot.splice(index, 1)
    this.validatorJobTimeslot.updateForm("jobTimeslots", jobTimeslot)
  }

  addJobTimeslot(){
    const jobTimeslot = this.validatorJobTimeslot.formValues().jobTimeslots
    jobTimeslot.push(create(JobTimeslotSchema))
    this.validatorJobTimeslot.updateForm("jobTimeslots", jobTimeslot)
  }

  toggleEnrollmentSelection(orderType: OrderType){
    if(this.selectedEnrollments[orderType]){
      delete this.selectedEnrollments[orderType]
      return
    }
    this.selectedEnrollments[orderType] = true
  }

  getSelectedEnrollmentsLength(): number{
    return Object.keys(this.selectedEnrollments).length
  }
}
