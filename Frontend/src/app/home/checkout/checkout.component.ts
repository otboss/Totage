import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { IonGrid, IonRow, IonCol } from "@ionic/angular/standalone";

@Component({
  selector: 'app-checkout',
  templateUrl: './checkout.component.html',
  styleUrls: ['./checkout.component.scss'],
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonGrid, IonRow, IonCol],
})
export class CheckoutComponent  implements OnInit {

  constructor() { }

  ngOnInit() {}

}
