import { Component, OnInit, ChangeDetectionStrategy, effect } from '@angular/core';
import { IonContent, IonGrid, IonRow, IonInput, IonToggle, IonButton, ToastController, NavController, IonRouterOutlet, IonIcon } from "@ionic/angular/standalone";
import { AddEditFreightValidator } from './add-edit-freight.validator';
import { ActivatedRoute, Router } from '@angular/router';
import { AddEditFreightService } from './add-edit-freight.service';
import { GetFreightConfigQuerySchema, UpdateFreightConfigCommandSchema } from 'src/app/shared/services/client/modules/organization/api/organization_pb';
import { create } from '@bufbuild/protobuf';
import { Code, ConnectError } from '@connectrpc/connect';
import { AppRoutes } from 'src/constants/AppRoutes';
import { SelectLocationPage } from "src/app/shared/select-location/select-location.page";
import { LocationSelectionModalComponent } from "src/app/shared/components/location-selection-modal/location-selection-modal.component";
import { trashBinOutline } from 'ionicons/icons';
import { addIcons } from 'ionicons';
import { LocationSelectorService } from 'src/app/shared/components/location-selector/location-selector.service';

@Component({
  selector: 'app-add-edit-freight',
  templateUrl: './add-edit-freight.component.html',
  styleUrls: ['./add-edit-freight.component.scss'],
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [SelectLocationPage, LocationSelectionModalComponent, IonInput, IonContent, IonGrid, IonRow, IonButton, IonIcon, IonToggle],
})
export class AddEditFreightComponent  implements OnInit {
  updateForm: typeof this.validator.updateForm = (x, y) => this.validator.updateForm(x, y) 
  formErrors = this.validator.formErrors
  formValues = this.validator.formValues
  showFieldErrors = this.validator.showFieldErrors
  isSubmitAttempted: boolean = false
  isSubmitting: boolean = false

  constructor(
    private validator: AddEditFreightValidator,
    private toastController: ToastController,
    private route: ActivatedRoute,
    private addEditFreightService: AddEditFreightService,
    private navCtrl: NavController,
    private routerOutlet: IonRouterOutlet,
    private router: Router,
    private locationSelectorService: LocationSelectorService
  ) {
    addIcons({ trashBinOutline })

    effect(() => {
      const selectedLocation = this.locationSelectorService.getSelectedLocation()
      if(selectedLocation != null){
        switch(this.locationSelectorService.getTag()){
          case "solid-waste":
            this.validator.updateForm('solidWasteDumpSiteLatitude', selectedLocation.latitude)
            this.validator.updateForm('solidWasteDumpSiteLongitude', selectedLocation.longitude)
            break
          case "liquid-waste":
            this.validator.updateForm('liquidWasteDumpSiteLatitude', selectedLocation.latitude)
            this.validator.updateForm('liquidWasteDumpSiteLongitude', selectedLocation.longitude)
            break
        }
        this.locationSelectorService.setTag("")
      }
    })

  }

  async ngOnInit() {
    this.validator.reset()
    const organizationId = this.route.snapshot.paramMap.get('organizationId');

    try{
      const getSalesConfig  = await this.addEditFreightService.getFreightConfig(
        create(GetFreightConfigQuerySchema, {
          organizationId: BigInt(organizationId)
        })
      )

      this.validator.updateForm('organizationId', BigInt(organizationId))
      this.validator.updateForm('solidWasteDumpSiteLatitude', getSalesConfig.solidWasteDumpSiteLatitude)
      this.validator.updateForm('solidWasteDumpSiteLongitude', getSalesConfig.solidWasteDumpSiteLongitude)
      this.validator.updateForm('liquidWasteDumpSiteLatitude', getSalesConfig.liquidWasteDumpSiteLatitude)
      this.validator.updateForm('liquidWasteDumpSiteLongitude', getSalesConfig.liquidWasteDumpSiteLongitude)
      this.validator.updateForm('isEnabled', getSalesConfig.isEnabled)
    }
    catch(err){
      if (err instanceof ConnectError) {
        if (err.code !== Code.NotFound) {
            // TODO: 
        }
      }
    }
  }

  async onSubmit(){
    this.isSubmitting = true
    this.isSubmitAttempted = true
    try{
      console.log(this.formErrors())
      if(Object.keys(this.formErrors()).length > 0){
        return
      }
      await this.addEditFreightService.updateFreightConfig(
        create(UpdateFreightConfigCommandSchema, this.validator.formValues())
      )
      if(this.routerOutlet && this.routerOutlet.canGoBack()){
        this.navCtrl.pop()
      } else {
        this.navCtrl.navigateBack(AppRoutes.Dashboard())
      }
      const toast = await this.toastController.create({
        message: 'Settings updated',
        duration: 2500,
        position: "bottom",
      })
      await toast.present()
    }
    finally{
      this.isSubmitting = false
    }
  }


  getIsEnabled(){
    return this.validator.formValues().isEnabled
  }


  setLocationTag(tag: "solid-waste" | "liquid-waste"){
    this.locationSelectorService.setTag(tag)
  }

  viewLocation(latitude: number, longitude: number){
    this.locationSelectorService.setSelectedLocation({
      latitude,
      longitude
    })
    // this.navigateToLocationSelector()
    this.router.navigateByUrl(AppRoutes.LocationSelector(), {
      skipLocationChange: true,
      state: {
        viewOnly: true
      }
    })
  }

  navigateToLocationSelector(){
    this.router.navigateByUrl(AppRoutes.LocationSelector(), {
      skipLocationChange: true
    })
  }

}
