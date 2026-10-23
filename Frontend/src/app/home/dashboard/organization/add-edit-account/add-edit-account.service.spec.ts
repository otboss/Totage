import { TestBed } from '@angular/core/testing';

import { AddEditAccountService } from './add-edit-account.service';

describe('AddEditAccountService', () => {
  let service: AddEditAccountService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(AddEditAccountService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
