import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';
import { createClient } from "@connectrpc/connect"
import { GeocodingQuery, GetFreightTruckClassesNearbyQuery, LogisticService } from 'src/app/shared/services/client/modules/logistic/api/logistic_pb';

@Injectable({
  providedIn: 'root',
})
export class OrderDetailsFreightService {

  async getFreightTruckClassesNearby(requestBody: GetFreightTruckClassesNearbyQuery) {
      const client = createClient(LogisticService, environment.connectTransport)
      return await client.getFreightTruckClassesNearbyHandler(requestBody)
  }

  async geocodingQuery(requestBody: GeocodingQuery){
      const client = createClient(LogisticService, environment.connectTransport)
      return await client.geocodingHandler(requestBody)
  }
}
