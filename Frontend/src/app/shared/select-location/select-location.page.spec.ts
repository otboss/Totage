import { ComponentFixture, TestBed } from '@angular/core/testing';
import { SelectLocationPage } from './select-location.page';

describe('SelectLocationPage', () => {
  let component: SelectLocationPage;
  let fixture: ComponentFixture<SelectLocationPage>;

  beforeEach(() => {
    fixture = TestBed.createComponent(SelectLocationPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
