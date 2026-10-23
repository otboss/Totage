import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { IonGrid, IonSelect, IonRow, IonSelectOption, IonButton, IonRadioGroup, IonRadio, IonCol } from "@ionic/angular/standalone";

@Component({
  selector: 'app-order-details-water-purified',
  templateUrl: './order-details-water-purified.component.html',
  styleUrls: ['./order-details-water-purified.component.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonGrid, IonSelect, IonRow, IonSelectOption, IonButton, IonRadioGroup, IonRadio, IonCol]
})
export class OrderDetailsWaterPurifiedComponent  implements OnInit {

  constructor() { }

  ngOnInit() {}

}
