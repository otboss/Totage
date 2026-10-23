import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';

import { FormsModule } from '@angular/forms';
import { IonContent, IonHeader, IonTitle, IonToolbar, IonButtons, IonBackButton } from '@ionic/angular/standalone';
import { OrderType } from 'src/app/shared/services/client/shared/api/order_type_pb'
import { ActivatedRoute, Router } from '@angular/router';
import { AppRoutes } from 'src/constants/AppRoutes';
import { AddEditEggComponent } from "./add-edit-egg/add-edit-egg.component";
import { AddEditWaterPurifiedComponent } from "./add-edit-water-purified/add-edit-water-purified.component";
import { AddEditWaterUnpurifiedComponent } from "./add-edit-water-unpurified/add-edit-water-unpurified.component";
import { AddEditHoneyComponent } from "./add-edit-honey/add-edit-honey.component";
import { AddEditBlockComponent } from "./add-edit-block/add-edit-block.component";
import { AddEditCharcoalComponent } from "./add-edit-charcoal/add-edit-charcoal.component";
import { AddEditCoconutComponent } from "./add-edit-coconut/add-edit-coconut.component";
import { AddEditIceComponent } from "./add-edit-ice/add-edit-ice.component";
import { AddEditConstructionComponent } from "./add-edit-construction/add-edit-construction.component";
import { AddEditFreightComponent } from "./add-edit-freight/add-edit-freight.component";

@Component({
  selector: 'app-add-edit-order-type',
  templateUrl: './add-edit-order-type.page.html',
  styleUrls: ['./add-edit-order-type.page.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonContent, IonHeader, IonTitle, IonToolbar, FormsModule, AddEditEggComponent, AddEditWaterPurifiedComponent, AddEditWaterUnpurifiedComponent, AddEditHoneyComponent, AddEditBlockComponent, AddEditCharcoalComponent, AddEditCoconutComponent, AddEditIceComponent, AddEditConstructionComponent, IonButtons, IonBackButton, AddEditFreightComponent]
})
export class AddEditOrderTypePage implements OnInit {

  selectedOrderType: OrderType | null = null
  orderType = OrderType
  readonly dashboardPage: string = AppRoutes.Dashboard()

  constructor(
    private router: Router,
    private route: ActivatedRoute,
  ) { }

  ngOnInit() {
    const navigation = this.router.currentNavigation();
    const organizationId = this.route.snapshot.paramMap.get('organizationId')
    if(organizationId === null){
      this.router.navigateByUrl(AppRoutes.Dashboard())
      return
    }
    this.selectedOrderType = <OrderType>navigation?.extras?.state?.orderType ?? null
    if(this.selectedOrderType === null || OrderType[this.selectedOrderType] === undefined){
      this.router.navigateByUrl(AppRoutes.Dashboard())
      return
    }
  }



}
