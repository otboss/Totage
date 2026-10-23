import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';

import { FormsModule } from '@angular/forms';
import { IonContent, IonHeader, IonTitle, IonToolbar,IonButtons, IonBackButton, IonButton, IonInput, IonGrid, IonRow, IonToggle, ToastController, NavController, IonRouterOutlet } from '@ionic/angular/standalone';
import { AppRoutes } from 'src/constants/AppRoutes';
import { AddEditAccountValidator } from './add-edit-account.validator';
import { Account, CreateAccountCommandSchema, GetAccountQuerySchema } from 'src/app/shared/services/client/modules/organization/api/organization_pb';
import { ActivatedRoute, Router } from '@angular/router';
import { AddEditAccountService } from './add-edit-account.service';
import { create } from '@bufbuild/protobuf';

@Component({
  selector: 'app-add-edit-account',
  templateUrl: './add-edit-account.page.html',
  styleUrls: ['./add-edit-account.page.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonContent, IonHeader, IonTitle, IonToolbar, FormsModule, IonButtons, IonBackButton, IonButton, IonInput, IonGrid, IonRow, IonToggle]
})
export class AddEditAccountPage implements OnInit {
  readonly dashboardPage: string = AppRoutes.Dashboard()

  updateForm: typeof this.validator.updateForm = (x, y) => this.validator.updateForm(x, y) 
  formErrors = this.validator.formErrors
  showFieldErrors = this.validator.showFieldErrors
  isSubmitAttempted: boolean = false
  isSubmitting: boolean = false

  account: Account | null = null

  constructor(
    private addEditAccountService: AddEditAccountService,
    public validator: AddEditAccountValidator,
    private router: Router,
    private route: ActivatedRoute,
    private navCtrl: NavController,
    private routerOutlet: IonRouterOutlet,
    private toastController: ToastController,
  ) { }

  async ngOnInit() {
    this.validator.reset()
    const navigation = this.router.currentNavigation();
    this.account = <Account>navigation?.extras?.state?.account ?? null

    const organizationId = this.route.snapshot.paramMap.get('organizationId')
    const accountId = this.route.snapshot.paramMap.get('accountId')

    if(accountId !== null && this.account === null){
      const accountResponse = await this.addEditAccountService.getAccount(
        create(GetAccountQuerySchema, {
          accountId: accountId,
          organizationId: BigInt(organizationId)
        })
      )
      this.account = accountResponse.account
    }

    if(this.account === null){
      return
    }

    this.validator.updateForm('email', this.account.email)
    this.validator.updateForm('isAdmin', this.account.isAdmin)
    this.validator.updateForm('isDriver', this.account.isDriver)
    this.validator.updateForm('isCustomerSupport', this.account.isCustomerSupport)
  }

  async onSubmit(){
    this.isSubmitting = true
    this.isSubmitAttempted = true
    try{
      if(Object.keys(this.formErrors()).length > 0){
        return
      }

      const requestBody = create(CreateAccountCommandSchema, {
        ...this.validator.formValues(),
      })

      await this.addEditAccountService.createAccount(requestBody)

      if(this.routerOutlet && this.routerOutlet.canGoBack()){
        this.navCtrl.pop()
      } else {
        this.navCtrl.navigateBack(AppRoutes.Dashboard())
      }

      const toast = await this.toastController.create({
        message: `Account ${this.account ? 'updated' : 'created'} successfully`,
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
}
