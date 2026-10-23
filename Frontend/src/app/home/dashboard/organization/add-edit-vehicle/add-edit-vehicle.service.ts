import { Injectable, inject } from "@angular/core"
import { createClient } from "@connectrpc/connect"
import { GetVehiclesQuery, LogisticService, CreateTruckCommand, CreateRentalVehicleCommand, CreateVehicleCommand, GetRentalVehicleQuery, GetTruckQuery, TogglePublicityRentalVehicleCommand, TogglePublicityTruckTruckCommand, GetVehicleQuery, UpdateVehicleDriverCommand, GetJobTimeslotsQuery, SetJobTimeslotCommand, DeleteVehicleCommand } from "src/app/shared/services/client/modules/logistic/api/logistic_pb"
import { GetCompatibleEnrollmentQuery, GetEnrollmentQuery, OrderService, SetEnrollmentCommand } from "src/app/shared/services/client/modules/order/api/order_pb"
import { GetDriversQuery, OrganizationService } from "src/app/shared/services/client/modules/organization/api/organization_pb"
import { UserSessionService } from "src/app/user-session.service"
import { environment } from "src/environments/environment"

@Injectable({
  providedIn: 'root',
})
export class AddEditVehicleService {

  private userSessionService: UserSessionService

  constructor(){
    this.userSessionService = inject(UserSessionService)
  }

  async createVehicle(requestBody: CreateVehicleCommand){
    const client = createClient(LogisticService, environment.connectTransport)
    return await client.createVehicleHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async getVehicle(requestBody: GetVehicleQuery){
    const client = createClient(LogisticService, environment.connectTransport)
    return await client.getVehicleHandler(requestBody,{
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

  async getRentalVehicle(requestBody: GetRentalVehicleQuery){
    const client = createClient(LogisticService, environment.connectTransport)
    return await client.getRentalVehicleHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async getTruck(requestBody: GetTruckQuery){
    const client = createClient(LogisticService, environment.connectTransport)
    return await client.getTruckHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async toggleTruckPublicity(requestBody: TogglePublicityTruckTruckCommand){
    const client = createClient(LogisticService, environment.connectTransport)
    return await client.togglePublicityTruckHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async toggleRentalVehiclePublicity(requestBody: TogglePublicityRentalVehicleCommand){
    const client = createClient(LogisticService, environment.connectTransport)
    return await client.togglePublicityRentalVehicleHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async getDrivers(requestBody: GetDriversQuery){
    const client = createClient(OrganizationService, environment.connectTransport)
    return await client.getDriversHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async updateVehicleDriver(requestBody: UpdateVehicleDriverCommand){
    const client = createClient(LogisticService, environment.connectTransport)
    return await client.updateVehicleDriverHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async getEnrollment(requestBody: GetEnrollmentQuery){
    const client = createClient(OrderService, environment.connectTransport)
    return await client.getEnrollmentHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async setEnrollment(requestBody: SetEnrollmentCommand){
    const client = createClient(OrderService, environment.connectTransport)
    return await client.setEnrollmentHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async getJobTimeslots(requestBody: GetJobTimeslotsQuery){
    const client = createClient(LogisticService, environment.connectTransport)
    return await client.getJobTimeslotsHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async setJobTimeslots(requestBody: SetJobTimeslotCommand){
    const client = createClient(LogisticService, environment.connectTransport)
    return await client.setJobTimeslotsHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }

  async deleteVehicle(requestBody: DeleteVehicleCommand){
    const client = createClient(LogisticService, environment.connectTransport)
    return await client.deleteVehicleHandler(requestBody,{
      headers: {
        "Authorization": `Bearer ${this.userSessionService.getAuthToken().jwtToken}`
      }
    })
  }
}