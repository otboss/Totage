import { inject, Injectable } from '@angular/core';
import { createClient } from '@connectrpc/connect';
import { CreateAccountCommand, GetAccountQuery, OrganizationService } from 'src/app/shared/services/client/modules/organization/api/organization_pb';
import { UserSessionService } from 'src/app/user-session.service';
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root',
})
export class AddEditAccountService {

  private userSessionService: UserSessionService

  constructor(){
    this.userSessionService = inject(UserSessionService)
  }

  async createAccount(requestBody: CreateAccountCommand){
    const client = createClient(OrganizationService, environment.connectTransport)
    return await client.createAccountHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async getAccount(requestBody: GetAccountQuery){
    const client = createClient(OrganizationService, environment.connectTransport)
    return await client.getAccountHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }
}
