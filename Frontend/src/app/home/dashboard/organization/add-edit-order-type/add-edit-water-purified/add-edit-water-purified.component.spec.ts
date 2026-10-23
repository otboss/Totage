import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { IonicModule } from '@ionic/angular';

import { AddEditWaterPurifiedComponent } from './add-edit-water-purified.component';

describe('AddEditWaterPurifiedComponent', () => {
  let component: AddEditWaterPurifiedComponent;
  let fixture: ComponentFixture<AddEditWaterPurifiedComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ AddEditWaterPurifiedComponent ],
      imports: [IonicModule.forRoot()]
    }).compileComponents();

    fixture = TestBed.createComponent(AddEditWaterPurifiedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  }));

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
