import { ComponentFixture, TestBed } from '@angular/core/testing';
import { AddEditAccountPage } from './add-edit-account.page';

describe('AddEditAccountPage', () => {
  let component: AddEditAccountPage;
  let fixture: ComponentFixture<AddEditAccountPage>;

  beforeEach(() => {
    fixture = TestBed.createComponent(AddEditAccountPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
