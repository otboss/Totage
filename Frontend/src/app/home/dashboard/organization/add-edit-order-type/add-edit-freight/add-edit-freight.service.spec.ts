import { TestBed } from '@angular/core/testing';

import { AddEditFreightService } from './add-edit-freight.service';

describe('AddEditFreightService', () => {
  let service: AddEditFreightService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(AddEditFreightService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
