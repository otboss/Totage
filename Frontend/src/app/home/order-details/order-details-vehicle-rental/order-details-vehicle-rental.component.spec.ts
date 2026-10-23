import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { IonicModule } from '@ionic/angular';

import { OrderDetailsVehicleRentalComponent } from './order-details-vehicle-rental.component';

describe('OrderDetailsVehicleRentalComponent', () => {
  let component: OrderDetailsVehicleRentalComponent;
  let fixture: ComponentFixture<OrderDetailsVehicleRentalComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ OrderDetailsVehicleRentalComponent ],
      imports: [IonicModule.forRoot()]
    }).compileComponents();

    fixture = TestBed.createComponent(OrderDetailsVehicleRentalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  }));

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
