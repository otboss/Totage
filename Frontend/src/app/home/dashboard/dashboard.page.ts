import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';

import { FormsModule } from '@angular/forms';
import { IonContent, IonHeader, IonTitle, IonToolbar, IonSpinner, RefresherCustomEvent, IonButton, IonCard, IonBadge, IonRippleEffect, IonText, IonGrid, IonRow, IonIcon, IonFab, IonFabButton, IonRefresher, IonRefresherContent, IonButtons, IonSearchbar, IonInput } from '@ionic/angular/standalone';
import { addIcons } from 'ionicons';
import { add, businessOutline, closeOutline, exitOutline, homeOutline, hourglassOutline, mailOutline, searchOutline, star, starHalf, starOutline } from 'ionicons/icons';
import { ActivatedRoute, Router } from '@angular/router';
import { AppRoutes } from 'src/constants/AppRoutes';
import { OrganizationService } from './organization.service';
import { GetOrganizationsQuerySchema, GetOrganizationsResponse, GetOrganizationsResponse_Organization } from 'src/app/shared/services/client/modules/organization/api/organization_pb';
import { create } from '@bufbuild/protobuf';
import { GetOrderQuerySchema } from 'src/app/shared/services/client/modules/order/api/order_pb';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.page.html',
  styleUrls: ['./dashboard.page.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonContent, IonHeader, IonTitle, IonToolbar, FormsModule, IonButton, IonCard, IonBadge, IonRippleEffect, IonText, IonGrid, IonRow, IonIcon, IonFab, IonFabButton, IonButtons, IonSearchbar, IonInput, IonRefresher, IonRefresherContent, IonSpinner]
})
export class DashboardPage implements OnInit {

  isSearchActive: boolean = false
  organizations: GetOrganizationsResponse | undefined
  isFetchingOrganizations: boolean = false

  constructor(private router: Router, private route: ActivatedRoute, private organizationService: OrganizationService) { 
    addIcons({ homeOutline, businessOutline, starOutline, star, starHalf, exitOutline, add, searchOutline, closeOutline, hourglassOutline, mailOutline });
  }

  ngOnInit() {
    this.refreshOrganizations()
  }

  async refreshOrganizations(){
    try{
      this.isFetchingOrganizations = true
      this.organizations = await this.organizationService.getOrganizations(
        create(GetOrganizationsQuerySchema, {
          pagination: {
            limit: BigInt(100)
          }
        })
      )
    }
    catch(err){
      // TODO: toast showing connection error
      console.error(err)
    }
    finally{
      this.isFetchingOrganizations = false
    }
  }

  goToOrganization(event: any, organization: GetOrganizationsResponse_Organization){
    this.router.navigateByUrl(
      AppRoutes.Organization().replace(":organizationId", organization.organizationId.toString()), 
      {
        state: organization
      }
    )
  }

  viewInvitations(){

  }

  async handleRefresherRefresh(event: RefresherCustomEvent){
    await this.refreshOrganizations()
    event.target.complete()
  }

  goHome(){
    this.router.navigateByUrl(AppRoutes.Home(), {replaceUrl: true})
  }

  createOrganization(){
    this.router.navigateByUrl(AppRoutes.OrganizationsCreate())
    // this.router.navigate([AppRoutes.OrganizationsCreate()], { relativeTo: this.route });

  }

  isHalfRating(rating: bigint, num: number): boolean{
    const ratingValue = Number(rating) - num
    return ratingValue > 0 && ratingValue < 1
  }

}
