import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { IonicModule } from '@ionic/angular';

import { AddEditWaterUnpurifiedComponent } from './add-edit-water-unpurified.component';

describe('AddEditWaterUnpurifiedComponent', () => {
  let component: AddEditWaterUnpurifiedComponent;
  let fixture: ComponentFixture<AddEditWaterUnpurifiedComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ AddEditWaterUnpurifiedComponent ],
      imports: [IonicModule.forRoot()]
    }).compileComponents();

    fixture = TestBed.createComponent(AddEditWaterUnpurifiedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  }));

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
