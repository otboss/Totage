import { Injectable } from '@angular/core';
import { createClient } from '@connectrpc/connect';
import { GetBulkGoodTrucksNearbyQuery, GetFreightTrucksNearbyQuery, GetRentalVehiclesNearbyQuery, LogisticService } from 'src/app/shared/services/client/modules/logistic/api/logistic_pb';
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root',
})
export class ResultsService {
    async getFreightTrucksNearby(requestBody: GetFreightTrucksNearbyQuery){
        const client = createClient(LogisticService, environment.connectTransport)
        return await client.getFreightTrucksNearbyHandler(requestBody)
    }

    async getRentalVehiclesNearby(requestBody: GetRentalVehiclesNearbyQuery){
        const client = createClient(LogisticService, environment.connectTransport)
        return await client.getRentalVehiclesNearbyHandler(requestBody)
    }

    async getBulkGoodTrucksNearby(requestBody: GetBulkGoodTrucksNearbyQuery){
        const client = createClient(LogisticService, environment.connectTransport)
        return await client.getBulkGoodTrucksNearbyHandler(requestBody)
    }
}
