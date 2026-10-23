import { inject, Injectable } from "@angular/core";
import { createClient } from "@connectrpc/connect";
import { GetPotableWaterSalesConfigQuery, OrganizationService, UpdatePotableWaterSalesConfigCommand } from "src/app/shared/services/client/modules/organization/api/organization_pb";
import { UserSessionService } from "src/app/user-session.service";
import { environment } from "src/environments/environment";


@Injectable({
  providedIn: 'root',
})
export class AddEditWaterPurifiedService {
  private userSessionService: UserSessionService

  constructor(){
    this.userSessionService = inject(UserSessionService)
  }

  async updatePotableWaterSalesConfig(requestBody: UpdatePotableWaterSalesConfigCommand){
    const client = createClient(OrganizationService, environment.connectTransport)
    return await client.updatePotableWaterSalesConfigHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async getPotableWaterSalesConfig(requestBody: GetPotableWaterSalesConfigQuery){
    const client = createClient(OrganizationService, environment.connectTransport)
    return await client.getPotableWaterSalesConfigHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

}