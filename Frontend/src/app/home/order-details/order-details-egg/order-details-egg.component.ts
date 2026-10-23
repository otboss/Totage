import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';

@Component({
  selector: 'app-order-details-egg',
  templateUrl: './order-details-egg.component.html',
  styleUrls: ['./order-details-egg.component.scss'],
  changeDetection: ChangeDetectionStrategy.Eager,
  standalone: true
})
export class OrderDetailsEggComponent  implements OnInit {

  constructor() { }

  ngOnInit() {}

}
