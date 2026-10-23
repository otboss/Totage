import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';

import { FormsModule } from '@angular/forms';
import { IonContent, IonHeader, IonTitle, IonToolbar, IonInput, IonGrid, IonRow, IonSelect, IonSelectOption, IonButtons, IonBackButton, IonButton } from '@ionic/angular/standalone';
import { OrganizationService } from '../organization.service';
import { Currency } from 'src/app/shared/services/client/shared/api/currency_pb';
import { CountryCodeAlpha2, CountryCodeAlpha3 } from 'src/app/shared/services/client/shared/api/country_pb';
import * as countryCodes from 'country-codes-list';
import { AppRoutes } from 'src/constants/AppRoutes';
import { OrganizationCreateValidator } from './organization-create.validator';
import { create } from '@bufbuild/protobuf';
import { CreateOrganizationCommandSchema } from 'src/app/shared/services/client/modules/organization/api/organization_pb';
import { countryToAlpha3 } from 'country-to-iso';
import { NavController } from '@ionic/angular';

@Component({
  selector: 'app-organization-create',
  templateUrl: './organization-create.page.html',
  styleUrls: ['./organization-create.page.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonContent, IonHeader, IonTitle, IonToolbar, FormsModule, IonInput, IonGrid, IonRow, IonSelect, IonSelectOption, IonButtons, IonBackButton, IonButton]
})
export class OrganizationCreatePage implements OnInit {

  readonly dashboardPage: string = AppRoutes.Dashboard()

  updateForm: typeof this.validator.updateForm = (x, y) => this.validator.updateForm(x, y) 
  formErrors = this.validator.formErrors
  showFieldErrors = this.validator.showFieldErrors
  isSubmitAttempted: boolean = false

  currencyOptions = Object.freeze(Object.keys(Currency).filter((key) => {
    // 1. Filter out the numeric reverse-mappings
    const isNumericKey = !isNaN(Number(key)); 
    
    // 2. Filter out the unspecified/default value
    const isUnspecified = key === 'CURRENCY_CODE_UNSPECIFIED'; 

    return !isNumericKey && !isUnspecified;
  }))

  countryCodeToNameMap = countryCodes.customList('countryCode', '{countryNameEn}');

  countryValues = Object.freeze(Object.values(CountryCodeAlpha2).filter((key) => {
    // 1. Filter out the numeric reverse-mapping keys
    const isNumericKey = !isNaN(Number(key));
    
    // 2. Filter out the proto3 default placeholder
    const isUnspecified = key === 'UNSPECIFIED';

    return !isNumericKey && !isUnspecified;
  }))

  constructor(private organizationService: OrganizationService, private validator: OrganizationCreateValidator, private navCtrl: NavController) { }

  ngOnInit() {
  }

  async onSubmit(){
    this.isSubmitAttempted = true
    try{
      if(Object.keys(this.validator.formErrors()).length > 0){
        return
      }
      await this.organizationService.createOrganization(
        create(CreateOrganizationCommandSchema, {
          ...this.validator.formValues(),
          country: CountryCodeAlpha3[countryToAlpha3(this.validator.formValues().country as any) as any] as CountryCodeAlpha3 | any,
          currency: Currency[this.validator.formValues().currency] as Currency | any
        })
      )
      this.navCtrl.navigateBack(AppRoutes.Dashboard())
    }
    catch(err){
      // TODO: Toast error message
      const error = err as Error
      if(error.message.indexOf("invalid address provided") !== -1){
        this.validator.customFieldError('address', "Please provide a valid full address")
      }
      console.error(err)
    }
  }

  Number(val: any){
    return Number(val)
  }

}
