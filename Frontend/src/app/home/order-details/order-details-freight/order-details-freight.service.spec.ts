import { TestBed } from '@angular/core/testing';

import { OrderDetailsFreightService } from './order-details-freight.service';

describe('OrderDetailsFreightService', () => {
  let service: OrderDetailsFreightService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(OrderDetailsFreightService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
