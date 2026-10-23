import { Component, OnInit, output, ChangeDetectionStrategy } from '@angular/core';
import { IonContent, IonHeader, IonTitle, IonToolbar, IonButtons, IonButton, IonCard, IonGrid, IonRow, IonCol, IonText, IonRippleEffect, IonInput, IonSearchbar } from '@ionic/angular/standalone';

import { FormsModule } from '@angular/forms';
import { Assets } from 'src/constants/Assets';
// import { OrderType } from 'src/enums/OrderType';
import { OrderType } from 'src/app/shared/services/client/shared/api/order_type_pb'
import { OrderTypeCards } from 'src/models/OrderTypeCards';

@Component({
  selector: 'app-start',
  templateUrl: './start.component.html',
  styleUrls: ['./start.component.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonContent, IonHeader, IonTitle, IonToolbar, FormsModule, IonButtons, IonButton, IonCard, IonGrid, IonRow, IonCol, IonText, IonRippleEffect, IonInput, IonSearchbar]
})
export class StartComponent  implements OnInit {
  orderType = OrderType

  output = output<{
    orderType: OrderType,
    image: string
    label: string
  }>()

  orderTypeCards: Array<{
    orderType: OrderType
    label: string
    image: string
  }> = Object.values(OrderTypeCards).filter(x => !(<Partial<Record<OrderType, true>>>{[OrderType.CESSPOOL]: true, [OrderType.CONSTRUCTION]: true})[x.orderType])
  
  // [
  //   {
  //     orderType: OrderType.VEHICLE_RENTAL,
  //     image: Assets.Images.rentalVehicle,
  //     label: "Vehicle Rental"
  //   },
  //   {
  //     orderType: OrderType.FREIGHT,
  //     image: Assets.Images.truck,
  //     label: "Moving"
  //   },
  //   {
  //     orderType: OrderType.WATER_PURIFIED,
  //     image: Assets.Images.water,
  //     label: "Drinking Water"
  //   },
  //   {
  //     orderType: OrderType.WATER_UNPURIFIED,
  //     image: Assets.Images.waterIrrigation,
  //     label: "Irrigation Water"
  //   },
  //   {
  //     orderType: OrderType.ICE,
  //     image: Assets.Images.ice,
  //     label: "Ice"
  //   },
  //   {
  //     orderType: OrderType.EGG,
  //     image: Assets.Images.eggs,
  //     label: "Eggs"
  //   },
  //   {
  //     orderType: OrderType.COCONUT,
  //     image: Assets.Images.coconuts,
  //     label: "Coconuts"
  //   },
  //   {
  //     orderType: OrderType.COOKING_GAS,
  //     image: Assets.Images.gasCylinder,
  //     label: "Cooking Gas"
  //   },
  //   {
  //     orderType: OrderType.HONEY,
  //     image: Assets.Images.honeyJar,
  //     label: "Honey"
  //   },
  //   {
  //     orderType: OrderType.BLOCK,
  //     image: Assets.Images.cinderBlocks,
  //     label: "Blocks"
  //   },
  //   {
  //     orderType: OrderType.CHARCOAL,
  //     image: Assets.Images.charcoalPile,
  //     label: "Charcoal"
  //   },
  //   // {
  //   //   orderType: OrderType.FREIGHT,
  //   //   image: Assets.Images.dumpTruck,
  //   //   label: "Dumping"
  //   // },
  //   // {
  //   //   orderType: OrderType.CESSPOOL,
  //   //   image: Assets.Images.cespoolTruck,
  //   //   label: "Cesspool"
  //   // },
  // ]

  constructor() { }

  ngOnInit() {
  }

}
