import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { IonicModule } from '@ionic/angular';

import { AddEditBitcoinLightningComponent } from './add-edit-bitcoin-lightning.component';

describe('AddEditBitcoinLightningComponent', () => {
  let component: AddEditBitcoinLightningComponent;
  let fixture: ComponentFixture<AddEditBitcoinLightningComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ AddEditBitcoinLightningComponent ],
      imports: [IonicModule.forRoot()]
    }).compileComponents();

    fixture = TestBed.createComponent(AddEditBitcoinLightningComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  }));

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
