import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';

@Component({
  selector: 'app-order-details-vehicle-rental',
  templateUrl: './order-details-vehicle-rental.component.html',
  styleUrls: ['./order-details-vehicle-rental.component.scss'],
  changeDetection: ChangeDetectionStrategy.Eager,
  standalone: true
})
export class OrderDetailsVehicleRentalComponent  implements OnInit {

  constructor() { }

  ngOnInit() {}

}
