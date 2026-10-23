import { ComponentFixture, TestBed } from '@angular/core/testing';
import { OrganizationSettingsPage } from './organization-settings.page';

describe('OrganizationSettingsPage', () => {
  let component: OrganizationSettingsPage;
  let fixture: ComponentFixture<OrganizationSettingsPage>;

  beforeEach(() => {
    fixture = TestBed.createComponent(OrganizationSettingsPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
