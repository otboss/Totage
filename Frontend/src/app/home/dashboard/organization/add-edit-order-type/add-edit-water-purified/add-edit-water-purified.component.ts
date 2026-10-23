import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { IonContent, IonGrid, IonRow, IonInput, IonToggle, IonButton, ToastController, NavController, IonRouterOutlet } from "@ionic/angular/standalone";
import { AddEditWaterPurifiedValidator } from './add-edit-water-purified.validator';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { AddEditWaterPurifiedService } from './add-edit-water-purified.service';
import { create } from '@bufbuild/protobuf';
import { GetPotableWaterSalesConfigQuerySchema, GetPotableWaterSalesConfigResponse, UpdatePotableWaterSalesConfigCommandSchema } from 'src/app/shared/services/client/modules/organization/api/organization_pb';
import { AppRoutes } from 'src/constants/AppRoutes';
import { GetPotableWaterOrderTotalQuerySchema } from 'src/app/shared/services/client/modules/order/api/order_pb';
import { ConnectError, Code } from '@connectrpc/connect';

@Component({
  selector: 'app-add-edit-water-purified',
  templateUrl: './add-edit-water-purified.component.html',
  styleUrls: ['./add-edit-water-purified.component.scss'],
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonContent, IonGrid, IonRow, IonInput, IonToggle, IonButton, FormsModule],
})
export class AddEditWaterPurifiedComponent  implements OnInit {

  updateForm: typeof this.validator.updateForm = (x, y) => this.validator.updateForm(x, y) 
  formErrors = this.validator.formErrors
  showFieldErrors = this.validator.showFieldErrors
  isSubmitAttempted: boolean = false
  isSubmitting: boolean = false

  constructor(
    private validator: AddEditWaterPurifiedValidator,
    private toastController: ToastController,
    private route: ActivatedRoute,
    private addEditWaterPurifiedService: AddEditWaterPurifiedService,
    private navCtrl: NavController,
    private routerOutlet: IonRouterOutlet,
  ) { }

  async ngOnInit() {
    this.validator.reset()
    const organizationId = this.route.snapshot.paramMap.get('organizationId');

    try{
      const getPotableWaterSalesConfigResponse  = await this.addEditWaterPurifiedService.getPotableWaterSalesConfig(
        create(GetPotableWaterSalesConfigQuerySchema, {
          organizationId: BigInt(organizationId)
        })
      )

      this.validator.updateForm('availableSupply', getPotableWaterSalesConfigResponse.availableSupply)
      this.validator.updateForm('bottleCost', getPotableWaterSalesConfigResponse.bottleCost)
      this.validator.updateForm('isEnabled', getPotableWaterSalesConfigResponse.isEnabled)
      this.validator.updateForm('minOrderQuantity', getPotableWaterSalesConfigResponse.minOrderQuantity)
      this.validator.updateForm('waterPrice', getPotableWaterSalesConfigResponse.waterPrice)
      this.validator.updateForm('waterPriceMin', getPotableWaterSalesConfigResponse.waterPriceMin)
      this.validator.updateForm('waterPricePerKm', getPotableWaterSalesConfigResponse.waterPricePerKm)
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
      await this.addEditWaterPurifiedService.updatePotableWaterSalesConfig(
        create(UpdatePotableWaterSalesConfigCommandSchema, this.validator.formValues())
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

  parseFloat(val: string){
    if(val === ""){
      return null
    }
    return parseFloat(val)
  }

  getIsEnabled(){
    return this.validator.formValues().isEnabled
  }
}
