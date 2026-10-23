import { Component, effect, input, OnInit, output, OutputRefSubscription, ChangeDetectionStrategy } from '@angular/core';
import { GetBulkGoodTrucksNearbyQuery, GetFreightTrucksNearbyQuery, GetRentalVehiclesNearbyQuery, GetRentalVehiclesNearbyResponse, GetTrucksNearbyResponse } from 'src/app/shared/services/client/modules/logistic/api/logistic_pb';
import { create } from "@bufbuild/protobuf";
import { Truck } from 'src/app/shared/services/client/shared/api/truck_pb';
import { IonButton, IonIcon, IonGrid, IonRow, IonCol, IonSpinner, IonCard } from "@ionic/angular/standalone";
import { ResultsService } from './results.service';

@Component({
  selector: 'app-results',
  templateUrl: './results.component.html',
  styleUrls: ['./results.component.scss'],
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonButton, IonIcon, IonGrid, IonRow, IonCol, IonSpinner, IonCard],
})
export class ResultsComponent  implements OnInit {

  output = output<Truck>()
  input = input<
    GetFreightTrucksNearbyQuery |
    GetRentalVehiclesNearbyQuery |
    GetBulkGoodTrucksNearbyQuery
  >()
  selectedTruckSubscription: OutputRefSubscription
  goPreviousSlide = output<void>()
  goNextSlide = output<void>()
  
  isLoading: boolean = false
  isLoadingErrorEncountered: boolean = false
  resultsTruck: GetTrucksNearbyResponse
  resultsRentalVehicle: GetRentalVehiclesNearbyResponse

  isSubmitted: boolean = false

  constructor(private resultsService: ResultsService) {
    this.selectedTruckSubscription = this.output.subscribe(x => {
      this.goNextSlide.emit()
    })

    effect(async () => {
      if(this.input()?.$typeName !== undefined){
        switch(this.input().$typeName){
          case "v1.GetFreightTrucksNearbyQuery":
            this.resultsTruck = await this.resultsService.getFreightTrucksNearby(this.input() as GetFreightTrucksNearbyQuery)
            break
          case "v1.GetRentalVehiclesNearbyQuery":
            this.resultsRentalVehicle = await this.resultsService.getRentalVehiclesNearby(this.input() as GetRentalVehiclesNearbyQuery)
            break
          case "v1.GetBulkGoodTrucksNearbyQuery":
            this.resultsTruck = await this.resultsService.getBulkGoodTrucksNearby(this.input() as GetBulkGoodTrucksNearbyQuery)
            break
          default:
            const _exhaustiveCheck: any = this.input().$typeName
            throw new Error(`Unhandled case: ${_exhaustiveCheck}`)
        }
      }
    })
  }

  ngOnInit() {}

  ngOnDestroy() {
    this.selectedTruckSubscription.unsubscribe()
  }

  reload(){

  }

}
