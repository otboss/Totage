import { ComponentFixture, TestBed } from '@angular/core/testing';
import { AddEditOrderTypePage } from './add-edit-order-type.page';

describe('AddEditOrderTypePage', () => {
  let component: AddEditOrderTypePage;
  let fixture: ComponentFixture<AddEditOrderTypePage>;

  beforeEach(() => {
    fixture = TestBed.createComponent(AddEditOrderTypePage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
