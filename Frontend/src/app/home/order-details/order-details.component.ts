import { Component, effect, input, Input, OnInit, output, signal, ChangeDetectionStrategy } from '@angular/core';
import { IonGrid, IonRow, IonCol, IonCard, IonRippleEffect, IonText, IonButton, IonIcon, IonContent } from "@ionic/angular/standalone";
import { OrderDetailsFreightComponent } from "./order-details-freight/order-details-freight.component";
import { OrderDetailsEggComponent } from "./order-details-egg/order-details-egg.component";
import { OrderDetailsBlockComponent } from "./order-details-block/order-details-block.component";
import { OrderDetailsCharcoalComponent } from "./order-details-charcoal/order-details-charcoal.component";
import { OrderDetailsVehicleRentalComponent } from "./order-details-vehicle-rental/order-details-vehicle-rental.component";
import { OrderDetailsWaterPurifiedComponent } from "./order-details-water-purified/order-details-water-purified.component";
import { OrderDetailsWaterUnpurifiedComponent } from "./order-details-water-unpurified/order-details-water-unpurified.component";
import { OrderDetailsCoconutComponent } from "./order-details-coconut/order-details-coconut.component";
import { addIcons } from 'ionicons';
import { chevronBackCircleOutline, chevronForwardCircleOutline } from 'ionicons/icons';
import { GetFreightTrucksNearbyQuery, GetRentalVehiclesNearbyQuery } from 'src/app/shared/services/client/modules/logistic/api/logistic_pb';
import { OrderDetailsFreightValidator } from './order-details-freight/order-details-freight.validator';
import { CompilationError, RuntimeError, ValidationError } from '@bufbuild/protovalidate';
import { OrderType } from 'src/app/shared/services/client/shared/api/order_type_pb';
import { OrderDetailsIceComponent } from "./order-details-ice/order-details-ice.component";

@Component({
  selector: 'app-order-details',
  templateUrl: './order-details.component.html',
  styleUrls: ['./order-details.component.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonGrid, IonRow, IonCol, IonCard, IonRippleEffect, IonText, OrderDetailsFreightComponent, OrderDetailsEggComponent, OrderDetailsBlockComponent, OrderDetailsCharcoalComponent, OrderDetailsVehicleRentalComponent, OrderDetailsWaterPurifiedComponent, OrderDetailsWaterUnpurifiedComponent, OrderDetailsCoconutComponent, IonButton, IonIcon, IonContent, OrderDetailsIceComponent],
})
export class OrderDetailsComponent  implements OnInit {

  goPreviousSlide = output<void>()
  goNextSlide = output<void>()
  isSubmitted: boolean = false
  formErrorsOutput: Partial<Record<any, string> & (ValidationError | CompilationError | RuntimeError)> = {}

  orderType = input<OrderType>(OrderType.FREIGHT)
  label = input<string>("")
  image = input<string>("")
  OrderType = OrderType

  output = output<
    GetFreightTrucksNearbyQuery |
    GetRentalVehiclesNearbyQuery
  >()

  // orderDetailsFreight = signal<Omit<GetFreightTrucksNearbyQuery, "$typeName" | "$unknown">>(null)
  // orderDetailsFreightOutput = output<Omit<GetFreightTrucksNearbyQuery, "$typeName" | "$unknown">>()

  vehicheRentalOrderDetails = output<{
    // TODO: 
  }>()



  constructor(private freightOrderValidator: OrderDetailsFreightValidator) {
    addIcons({ chevronForwardCircleOutline, chevronBackCircleOutline })
    // this.goPreviousSlide.subscribe(x => {
    //   this.isSubmitted = false
    // })
    // this.goNextSlide.subscribe(x => {
    //   this.isSubmitted = false
    // })
  }

  ngOnInit() {}

  formErrorsPresent(): boolean{
    return Object.keys(this.formErrorsOutput).length > 0
  }

  getIsSubmitted(){
    // GetFreightTruckClassesNearbyHandler
    return this.isSubmitted
  }

  setIsSubmitted(){
    this.isSubmitted = false
    setTimeout(() => {
      this.isSubmitted = true
    }, 0)
  }

  async onSubmit(){
    const coordsRegexp = /^-?\d+\.\d+,-?\d+\.\d+$/

    // TOOD: Check if addresses are coordinates
  }
}
