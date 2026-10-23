import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';

@Component({
  selector: 'app-order-details-block',
  templateUrl: './order-details-block.component.html',
  styleUrls: ['./order-details-block.component.scss'],
  changeDetection: ChangeDetectionStrategy.Eager,
  standalone: true
})
export class OrderDetailsBlockComponent  implements OnInit {

  constructor() { }

  ngOnInit() {}

}
