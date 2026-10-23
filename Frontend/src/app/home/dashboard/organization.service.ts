import { inject, Injectable } from '@angular/core';
import { createClient } from '@connectrpc/connect';
import { environment } from 'src/environments/environment';
import { CreateOrganizationCommand, GetAccountsQuery, GetOrganizationQuery, GetOrganizationsQuery, OrganizationService as OrganizationServ } from '../../shared/services/client/modules/organization/api/organization_pb';
import { CreateRentalVehicleCommand, CreateTruckCommand, GeocodingQuery, GetRentalVehiclesQuery, GetTrucksQuery, GetVehiclesQuery, LogisticService, ReverseGeocodingQuery } from 'src/app/shared/services/client/modules/logistic/api/logistic_pb';
import { UserSessionService } from 'src/app/user-session.service';
import { GetOrganizationOrdersQuery, OrderService } from 'src/app/shared/services/client/modules/order/api/order_pb';

@Injectable({
  providedIn: 'root',
})
export class OrganizationService {

  private userSessionService: UserSessionService

  constructor(){
    this.userSessionService = inject(UserSessionService)
  }

  async createOrganization(requestBody: CreateOrganizationCommand){
    const client = createClient(OrganizationServ, environment.connectTransport)
    return await client.createOrganizationHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async getOrganization(requestBody: GetOrganizationQuery) {
    const client = createClient(OrganizationServ, environment.connectTransport)
    return await client.getOrganizationHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async getOrganizations(requestBody: GetOrganizationsQuery) {
    // const userSessionService = inject(UserSessionService)
    const client = createClient(OrganizationServ, environment.connectTransport)
    return await client.getOrganizationsHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async getCountryfromAddress(requestBody: GeocodingQuery){
    const client = createClient(LogisticService, environment.connectTransport)
    return await client.geocodingHandler(requestBody)
  }

  async getAccounts(requestBody: GetAccountsQuery){
    const client = createClient(OrganizationServ, environment.connectTransport)
    return await client.getAccountsHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async getVehicles(requestBody: GetVehiclesQuery){
    const client = createClient(LogisticService, environment.connectTransport)
    return await client.getVehiclesHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async createTruck(requestBody: CreateTruckCommand){
    const client = createClient(LogisticService, environment.connectTransport)
    return await client.createTruckHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async createRentalVehicle(requestBody: CreateRentalVehicleCommand){
    const client = createClient(LogisticService, environment.connectTransport)
    return await client.createRentalVehicleHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async getOrders(requestBody: GetOrganizationOrdersQuery){
    const client = createClient(OrderService, environment.connectTransport)
    return await client.getOrganizationOrdersHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    }) 
  }
  
}
