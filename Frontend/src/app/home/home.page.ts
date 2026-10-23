import { Component, CUSTOM_ELEMENTS_SCHEMA, ElementRef, OnInit, output, ViewChild, ChangeDetectionStrategy } from '@angular/core';
import { LocationComponent } from './location/location.component';
import { CheckoutComponent } from './checkout/checkout.component';
import { JobDateComponent } from './job-date/job-date.component';
import { addIcons } from 'ionicons';
import { personCircleOutline } from 'ionicons/icons';
import { StartComponent } from './start/start.component';
import { OrderDetailsComponent } from './order-details/order-details.component';
import { AppRoutes } from 'src/constants/AppRoutes';
import { Router } from '@angular/router';
import { SwiperContainer } from 'swiper/element';
import { Truck } from 'src/app/shared/services/client/shared/api/truck_pb';
import { ResultsComponent } from "./results/results.component";
import getUnicodeFlagIcon from 'country-flag-icons/unicode'
import { OrderDetailsFreightValidator } from './order-details/order-details-freight/order-details-freight.validator';
import { GetFreightTrucksNearbyQuery, GetRentalVehiclesNearbyQuery } from '../shared/services/client/modules/logistic/api/logistic_pb';
import { OrderType } from '../shared/services/client/shared/api/order_type_pb';
import { UserSessionService } from '../user-session.service';
import { LoginService } from '../login/login.service';
// import * as flags from 'country-flag-icons/string/3x2'

@Component({
  selector: 'app-home',
  templateUrl: './home.page.html',
  styleUrls: ['./home.page.scss'],
  standalone: true,
  imports: [LocationComponent, JobDateComponent, CheckoutComponent, OrderDetailsComponent, StartComponent, ResultsComponent],
  changeDetection: ChangeDetectionStrategy.Eager,
  schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class HomePage implements OnInit {

  myValue: string | null = ""
  readonly dashboardPage: string = AppRoutes.Dashboard()
  isUserLoggedIn: boolean = false
  selectedOrderType: OrderType = OrderType.FREIGHT
  selectedOrderTypeImage: string = ''
  selectedOrderTypeLabel: string = ''
  regionFlag: string = `assets/flags/JM.svg`
  @ViewChild('swiper') swiperRef!: ElementRef<SwiperContainer>;

  
  orderDetails: GetFreightTrucksNearbyQuery | GetRentalVehiclesNearbyQuery
  selectedTruck: Truck

  constructor(private router: Router, private loginService: LoginService, private freightOrderValidator: OrderDetailsFreightValidator, private userSessionService: UserSessionService) { 
    addIcons({ personCircleOutline });
  }

  supportedOrderTypes: Record<string, boolean> = {

  }

  ngOnInit() {
    // if(this.userSession.getSession() != null){
    //   this.isUserLoggedIn = true
    // }
  }

 onChange(event: any) {
    switch(event.detail.value as "dashboard" | "signOut" | "signIn"){
      case "dashboard":
        this.router.navigateByUrl(AppRoutes.Dashboard(), {replaceUrl: true});
        break
      case "signOut":
        this.signOut()
        break
      case "signIn":
        this.signIn()
        break
    }
    // Logic to prevent action or reset value
    this.myValue = null; 
    event.target.value = null; // Reverts UI state [1, 2]
  }

  public isLoggedIn = this.userSessionService.isLoggedIn

  private async signOut(){
    const authToken = this.userSessionService.getAuthToken()
    switch(authToken.flow){
      case "GOOGLE":
        await this.loginService.logoutFromGoogle()
        break
      case "PASSWORD":
        // TODO: 
        break
      default:
        throw new Error("unhandled case: "+ authToken.flow) 
    }
    this.userSessionService.deleteAuthToken()
    window.location.reload()
  }

  private signIn(){
    this.router.navigateByUrl(AppRoutes.Login())
  }

  handleOrderTypeSelected(order: {
    orderType: OrderType;
    image: string;
    label: string
}){
    this.selectedOrderType = order.orderType
    this.selectedOrderTypeImage = order.image
    this.selectedOrderTypeLabel = order.label
    this.swiperRef.nativeElement.swiper.slideNext()
  }

  nextSlide(){
    this.swiperRef.nativeElement.swiper.slideNext()
  }

  previousSlide(){
    this.swiperRef.nativeElement.swiper.slidePrev()
  }

  getOrderDetails(){
    return this.orderDetails
  }
}
