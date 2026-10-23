import { ComponentFixture, TestBed } from '@angular/core/testing';
import { LocationSelectorPage } from './location-selector.page';

describe('LocationSelectorPage', () => {
  let component: LocationSelectorPage;
  let fixture: ComponentFixture<LocationSelectorPage>;

  beforeEach(() => {
    fixture = TestBed.createComponent(LocationSelectorPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
