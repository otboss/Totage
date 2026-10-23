import { inject, Injectable } from '@angular/core';
import { createClient } from '@connectrpc/connect';
import { GetFreightConfigQuery, OrganizationService, UpdateFreightConfigCommand } from 'src/app/shared/services/client/modules/organization/api/organization_pb';
import { UserSessionService } from 'src/app/user-session.service';
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root',
})
export class AddEditFreightService {
  private userSessionService: UserSessionService

  constructor(){
    this.userSessionService = inject(UserSessionService)
  }

  async getFreightConfig(requestBody: GetFreightConfigQuery){
    const client = createClient(OrganizationService, environment.connectTransport)
    return await client.getFreightConfigHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async updateFreightConfig(requestBody: UpdateFreightConfigCommand){
    const client = createClient(OrganizationService, environment.connectTransport)
    return await client.updateFreightConfigHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

}
