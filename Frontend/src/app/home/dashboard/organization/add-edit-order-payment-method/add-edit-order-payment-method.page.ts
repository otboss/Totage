import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { IonContent, IonHeader, IonTitle, IonToolbar, IonButtons, IonBackButton } from '@ionic/angular/standalone';
import { PaymentMethod } from 'src/app/shared/services/client/modules/order/api/order_pb';
import { AppRoutes } from 'src/constants/AppRoutes';
import { AddEditBitcoinComponent } from './add-edit-bitcoin/add-edit-bitcoin.component';
import { AddEditBitcoinLightningComponent } from './add-edit-bitcoin-lightning/add-edit-bitcoin-lightning.component';
import { AddEditPaypalComponent } from './add-edit-paypal/add-edit-paypal.component';

@Component({
  selector: 'app-add-edit-order-payment-method',
  templateUrl: './add-edit-order-payment-method.page.html',
  styleUrls: ['./add-edit-order-payment-method.page.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonContent, IonHeader, IonTitle, IonToolbar, FormsModule, IonButtons, IonBackButton, AddEditBitcoinComponent, AddEditBitcoinLightningComponent, AddEditPaypalComponent]
})
export class AddEditOrderPaymentMethodPage implements OnInit {

  selectedPaymentMethod: PaymentMethod | null  = null
  paymentMethod = PaymentMethod
  readonly dashboardPage: string = AppRoutes.Dashboard()

  constructor() { }

  ngOnInit() {
  }

}
