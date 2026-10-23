import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { IonicModule } from '@ionic/angular';

import { OrderDetailsWaterUnpurifiedComponent } from './order-details-water-unpurified.component';

describe('OrderDetailsWaterUnpurifiedComponent', () => {
  let component: OrderDetailsWaterUnpurifiedComponent;
  let fixture: ComponentFixture<OrderDetailsWaterUnpurifiedComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ OrderDetailsWaterUnpurifiedComponent ],
      imports: [IonicModule.forRoot()]
    }).compileComponents();

    fixture = TestBed.createComponent(OrderDetailsWaterUnpurifiedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  }));

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
