import { Component, OnInit, ViewChild, ChangeDetectionStrategy } from '@angular/core';

import { FormsModule } from '@angular/forms';
import { IonContent, IonHeader, IonTitle, IonToolbar, IonTabBar, IonTabButton, IonIcon, IonTab, IonTabs, IonButton, IonButtons, IonBackButton, IonCard, IonText, IonRippleEffect, IonFab, IonFabButton, IonSegment, IonSegmentButton, IonLabel, IonSearchbar, IonRefresher, IonRefresherContent, IonSpinner, IonInfiniteScroll, IonInfiniteScrollContent } from '@ionic/angular/standalone';
import { addIcons } from 'ionicons';
import { add, albumsOutline, basketOutline, cardOutline, carOutline, eggOutline, libraryOutline, logoBitcoin, logoPaypal, peopleCircle, peopleOutline, personOutline, searchOutline, settingsOutline, starOutline, star, starHalf, waterOutline, refreshOutline, cubeOutline } from 'ionicons/icons';
import { AppRoutes } from 'src/constants/AppRoutes';
import { Assets } from 'src/constants/Assets';
import { ActivatedRoute, Router } from '@angular/router';
import { Account, GetAccountsQuerySchema, GetAccountsResponse, GetOrganizationQuerySchema, GetOrganizationResponse, GetOrganizationsResponse_Organization, GetOrganizationsResponse_OrganizationSchema, GetOrganizationsResponseSchema } from 'src/app/shared/services/client/modules/organization/api/organization_pb';
import { Location } from '@angular/common';
import { OrganizationService } from '../organization.service';
import { create } from '@bufbuild/protobuf';
import { InfiniteScrollCustomEvent, RefresherCustomEvent } from '@ionic/core';
import { GetVehiclesQuerySchema, GetVehiclesResponse, Vehicle } from 'src/app/shared/services/client/modules/logistic/api/logistic_pb';
import { GetOrganizationOrdersQuerySchema, GetOrganizationOrdersResponse } from 'src/app/shared/services/client/modules/order/api/order_pb';
import { OrderType } from 'src/app/shared/services/client/shared/api/order_type_pb';
// import { GetOrdersQuerySchema, GetOrdersResponse } from 'src/app/shared/services/client/modules/order/api/order_pb';

@Component({
  selector: 'app-organization',
  templateUrl: './organization.page.html',
  styleUrls: ['./organization.page.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonContent, IonHeader, IonTitle, IonToolbar, FormsModule, IonTabBar, IonTabButton, IonIcon, IonTab, IonTabs, IonButton, IonButtons, IonBackButton, IonCard, IonText, IonRippleEffect, IonFab, IonFabButton, IonSegment, IonSegmentButton, IonLabel, IonSearchbar, IonRefresher, IonRefresherContent, IonSpinner, IonInfiniteScroll, IonInfiniteScrollContent]
})
export class OrganizationPage implements OnInit {

  @ViewChild('myTabs', { static: false }) tabs!: IonTabs;
  activeTabName?: "Accounts" | "Vehicles" | "Categories" | "Payment" | "Orders" = 'Accounts';
  readonly coconutImage: string = Assets.Images.coconut
  readonly cinderBlockImage: string = Assets.Images.cinderBlock
  readonly charcoalImage: string = Assets.Images.charcoal
  readonly honeyImage: string = Assets.Images.honey
  readonly cesspoolImage: string = Assets.Images.cesspool
  readonly dashboardPage: string = AppRoutes.Dashboard()
  readonly organizationSettingsPage: string = AppRoutes.OrganizationSettings()
  // readonly organizationCreateAccountPage: string = AppRoutes.OrganizationCreateAccount()
  isSearchOn: boolean = false
  

  limit = BigInt(100)
  offsetId = BigInt(0)
  filterText: string = ""
  // organization: GetOrganizationsResponse_Organization | null = null
  organization: GetOrganizationResponse | null = null

  accounts: GetAccountsResponse | undefined
  isFetchingAccounts: boolean = false
  vehicles: GetVehiclesResponse | undefined
  isFetchingVehicles: boolean = false
  orders: GetOrganizationOrdersResponse | undefined
  isFetchingOrders: boolean = false
  OrderType = OrderType

  constructor(private location: Location, private router: Router, private route: ActivatedRoute, private organizationService: OrganizationService) {
    addIcons({starOutline, star, starHalf, cubeOutline, settingsOutline, basketOutline, carOutline, peopleOutline, libraryOutline, personOutline, add, albumsOutline, waterOutline, eggOutline, cardOutline, logoPaypal, logoBitcoin, searchOutline, refreshOutline })
  }

  async ngOnInit() {
    const organizationId = this.route.snapshot.paramMap.get('organizationId')
    const requestBody = create(GetOrganizationQuerySchema, {
      organizationId: BigInt(organizationId),
    })
    this.organization = await this.organizationService.getOrganization(requestBody)

    // this.organization = this.location.getState() as GetOrganizationsResponse_Organization
    this.refreshAccounts()
    // if(this.organization === undefined){

    // }
  }

  async ngAfterViewInit(){
    const organizationId = this.route.snapshot.paramMap.get('id');
    // TODO: Fetch organization]
    // const response = await this.organizationService.getOrganization(create(GetOrganizationQuerySchema, {
    //   organizationId: BigInt(organizationId),
    // }))

    // response.responseBody
  }

  async refresh(){
    this.refreshAccounts()
    // TODO: 
  }

  async refreshAccounts(){
    try{
      if(this.accounts === undefined){
        this.isFetchingAccounts = true
      }
      this.accounts = await this.organizationService.getAccounts(
        create(
            GetAccountsQuerySchema, {
            organizationId: this.organization.organization.organizationId!,
            pagination: {
              limit: this.limit,
              filterText: this.isSearchOn ? this.filterText : "",
              offsetId: BigInt((this.accounts?.accounts[this.accounts?.accounts.length - 1]?.accountId) ?? 0)
            }
          }
        )
      )
    }
    catch(err){
      // TODO: Show toast about error
    }
    finally{
      this.isFetchingAccounts = false
    }
  }

  async refreshVehicles(){
    try{
      if(this.vehicles === undefined){
        this.isFetchingVehicles = true
      }
      this.vehicles = await this.organizationService.getVehicles(
        create(
            GetVehiclesQuerySchema, {
            organizationId: this.organization.organization.organizationId!,
            pagination: {
              limit: BigInt(100),
              filterText: this.isSearchOn ? this.filterText : "",
              offsetId: BigInt((this.vehicles?.vehicles[this.vehicles?.vehicles.length - 1]?.vehicleId) ?? 0)
            }
          }
        )
      )
    }
    catch(err){
      // TODO: Show toast about error
    }
    finally{
      this.isFetchingVehicles = false
    }
  }

  async refreshOrders(){
    try{
      if(this.orders === undefined){
        this.isFetchingOrders = true
      }
      this.orders = await this.organizationService.getOrders(
        create(
            GetOrganizationOrdersQuerySchema, {
            organizationId: this.organization.organization.organizationId!,
            pagination: {
              limit: this.limit,
              filterText: this.isSearchOn ? this.filterText : "",
              offsetId: BigInt((this.orders?.orders[this.orders?.orders.length - 1]?.orderId) ?? 0)
            },
          }
        )
      )
    }
    catch(err){
      // TODO: Show toast about error
    }
    finally{
      this.isFetchingOrders = false
    }
  }

  async handleRefresherRefresh(event: RefresherCustomEvent){
    await this.refreshAccounts() // TODO: REfresh active tab contents only
    event.target.complete()
  }

  setCurrentTab(event: any){
    this.activeTabName = event.target.value
  }

  goToSettings(){
    this.router.navigateByUrl(AppRoutes.OrganizationSettings())
  }

  goToAddVehiclePage(){
    this.router.navigateByUrl(AppRoutes.VehiclesCreate().replace(":organizationId", this.organization.organization.organizationId!.toString())),
    {
      state: {
        organization: this.organization
      }
    }
  }

  goToEditVehiclePage(vehicle: Vehicle){
    this.router.navigateByUrl(`${AppRoutes.VehiclesEdit(vehicle.vehicleId).replace(":organizationId", this.organization.organization.organizationId!.toString())}`, 
    {
      state: {
        organizationId: this.organization.organization.organizationId, // TODO: Set organization id correctlu
        // organization: this.organization,
        vehicle,
      }
    })
  }

  getOrganization(){
    return this.organization
  }

  goToOrderDetailsPage(orderId: string){
    // this.router.navigateByUrl(`${AppRoutes.Order(vehicleId)}`)
  }

  isHalfRating(rating: bigint, num: number): boolean{
    const ratingValue = Number(rating) - num
    return ratingValue > 0 && ratingValue < 1
  }

  createAccount(){
    this.router.navigateByUrl(
      AppRoutes.OrganizationCreateAccount().replace(":organizationId", this.organization.organization.organizationId.toString())
    )
  }

  editAccount(account: Account){
    this.router.navigateByUrl(
      AppRoutes.OrganizationEditAccount().replace(":organizationId", this.organization.organization.organizationId.toString()).replace(":accountId", account.accountId), {
        state: account
      }
    )
  }

  onIonInfinite(event: InfiniteScrollCustomEvent) {
    if(this.activeTabName === "Accounts"){
      this.refreshAccounts()
    } else
    if(this.activeTabName === "Orders"){
      this.refreshOrders()
    } else 
    if(this.activeTabName === "Vehicles"){
      this.refreshVehicles()
    }

    setTimeout(() => {
      event.target.complete();
    }, 500);
  }

  navigateToAddEdit(orderType: OrderType){
    this.router.navigateByUrl(AppRoutes.OrganizationAddEditOrderType().replace(":organizationId", this.organization.organization.organizationId!.toString()), {
      state: { orderType }
    })
  }
}
