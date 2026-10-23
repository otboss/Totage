import { ComponentFixture, TestBed } from '@angular/core/testing';
import { AddEditVehiclePage } from './add-edit-vehicle.page';

describe('AddEditVehiclePage', () => {
  let component: AddEditVehiclePage;
  let fixture: ComponentFixture<AddEditVehiclePage>;

  beforeEach(() => {
    fixture = TestBed.createComponent(AddEditVehiclePage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
