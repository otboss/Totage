import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { IonContent, IonGrid, IonRow, IonInput, IonToggle, IonButton, ToastController, NavController, IonRouterOutlet } from "@ionic/angular/standalone";
import { AddEditFreightValidator } from './add-edit-freight.validator';
import { ActivatedRoute } from '@angular/router';
import { AddEditFreightService } from './add-edit-freight.service';
import { GetFreightConfigQuerySchema, UpdateFreightConfigCommandSchema } from 'src/app/shared/services/client/modules/organization/api/organization_pb';
import { create } from '@bufbuild/protobuf';
import { Code, ConnectError } from '@connectrpc/connect';
import { AppRoutes } from 'src/constants/AppRoutes';
import { SelectLocationPage } from "src/app/shared/select-location/select-location.page";
import { LocationSelectionModalComponent } from "src/app/shared/components/location-selection-modal/location-selection-modal.component";

@Component({
  selector: 'app-add-edit-freight',
  templateUrl: './add-edit-freight.component.html',
  styleUrls: ['./add-edit-freight.component.scss'],
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [SelectLocationPage, LocationSelectionModalComponent],
})
export class AddEditFreightComponent  implements OnInit {
  updateForm: typeof this.validator.updateForm = (x, y) => this.validator.updateForm(x, y) 
  formErrors = this.validator.formErrors
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
  ) { }

  async ngOnInit() {
    this.validator.reset()
    const organizationId = this.route.snapshot.paramMap.get('organizationId');

    try{
      const getPotableWaterSalesConfigResponse  = await this.addEditFreightService.getFreightConfig(
        create(GetFreightConfigQuerySchema, {
          organizationId: BigInt(organizationId)
        })
      )

      this.validator.updateForm('solidWasteDumpSiteLatitude', getPotableWaterSalesConfigResponse.solidWasteDumpSiteLatitude)
      this.validator.updateForm('solidWasteDumpSiteLongitude', getPotableWaterSalesConfigResponse.solidWasteDumpSiteLongitude)
      this.validator.updateForm('liquidWasteDumpSiteLatitude', getPotableWaterSalesConfigResponse.liquidWasteDumpSiteLatitude)
      this.validator.updateForm('liquidWasteDumpSiteLongitude', getPotableWaterSalesConfigResponse.liquidWasteDumpSiteLongitude)
      this.validator.updateForm('isEnabled', getPotableWaterSalesConfigResponse.isEnabled)
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
      if(Object.keys(this.formErrors()).length > 0){
        return
      }
      this.validator.updateForm('organizationId', BigInt(this.route.snapshot.paramMap.get('organizationId')))
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

}
