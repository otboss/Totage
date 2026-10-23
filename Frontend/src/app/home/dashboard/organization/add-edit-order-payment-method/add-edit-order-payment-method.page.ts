import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';

import { FormsModule } from '@angular/forms';
import { IonContent, IonHeader, IonTitle, IonToolbar } from '@ionic/angular/standalone';

@Component({
  selector: 'app-add-edit-order-payment-method',
  templateUrl: './add-edit-order-payment-method.page.html',
  styleUrls: ['./add-edit-order-payment-method.page.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonContent, IonHeader, IonTitle, IonToolbar, FormsModule]
})
export class AddEditOrderPaymentMethodPage implements OnInit {

  constructor() { }

  ngOnInit() {
  }

}
