import { ComponentFixture, TestBed } from '@angular/core/testing';
import { OrganizationCreatePage } from './organization-create.page';

describe('OrganizationCreatePage', () => {
  let component: OrganizationCreatePage;
  let fixture: ComponentFixture<OrganizationCreatePage>;

  beforeEach(() => {
    fixture = TestBed.createComponent(OrganizationCreatePage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
