import { inject, Injectable } from "@angular/core";
import { createClient } from "@connectrpc/connect";
import { GetIrrigationWaterSalesConfigQuery, OrganizationService, UpdateIrrigationWaterSalesConfigCommand } from "src/app/shared/services/client/modules/organization/api/organization_pb";
import { UserSessionService } from "src/app/user-session.service";
import { environment } from "src/environments/environment";


@Injectable({
  providedIn: 'root',
})
export class AddEditWaterUnpurifiedService {
  private userSessionService: UserSessionService

  constructor(){
    this.userSessionService = inject(UserSessionService)
  }

  async updateIrrigationWaterSalesConfig(requestBody: UpdateIrrigationWaterSalesConfigCommand){
    const client = createClient(OrganizationService, environment.connectTransport)
    return await client.updateIrrigationWaterSalesConfigHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async getIrrigationWaterSalesConfig(requestBody: GetIrrigationWaterSalesConfigQuery){
    const client = createClient(OrganizationService, environment.connectTransport)
    return await client.getIrrigationWaterSalesConfigHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

}