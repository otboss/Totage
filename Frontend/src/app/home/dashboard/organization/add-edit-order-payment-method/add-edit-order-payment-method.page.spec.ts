import { ComponentFixture, TestBed } from '@angular/core/testing';
import { AddEditOrderPaymentMethodPage } from './add-edit-order-payment-method.page';

describe('AddEditOrderPaymentMethodPage', () => {
  let component: AddEditOrderPaymentMethodPage;
  let fixture: ComponentFixture<AddEditOrderPaymentMethodPage>;

  beforeEach(() => {
    fixture = TestBed.createComponent(AddEditOrderPaymentMethodPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
