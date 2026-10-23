import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { IonicModule } from '@ionic/angular';

import { OrderDetailsWaterPurifiedComponent } from './order-details-water-purified.component';

describe('OrderDetailsWaterPurifiedComponent', () => {
  let component: OrderDetailsWaterPurifiedComponent;
  let fixture: ComponentFixture<OrderDetailsWaterPurifiedComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ OrderDetailsWaterPurifiedComponent ],
      imports: [IonicModule.forRoot()]
    }).compileComponents();

    fixture = TestBed.createComponent(OrderDetailsWaterPurifiedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  }));

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
