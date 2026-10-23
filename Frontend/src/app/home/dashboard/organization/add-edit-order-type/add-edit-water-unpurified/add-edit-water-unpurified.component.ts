import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { IonContent, IonGrid, IonRow, IonToggle, IonInput, IonButton } from "@ionic/angular/standalone";
import { AddEditWaterUnpurifiedValidator } from './add-edit-water-unpurified.validator';
import { FormsModule } from '@angular/forms';
import { AddEditWaterUnpurifiedService } from './add-edit-water-unpurified.service';

@Component({
  selector: 'app-add-edit-water-unpurified',
  templateUrl: './add-edit-water-unpurified.component.html',
  styleUrls: ['./add-edit-water-unpurified.component.scss'],
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonContent, IonGrid, IonRow, IonToggle, IonInput, IonButton, FormsModule],
})
export class AddEditWaterUnpurifiedComponent  implements OnInit {
  updateForm: typeof this.validator.updateForm = (x, y) => this.validator.updateForm(x, y) 
  formErrors = this.validator.formErrors
  showFieldErrors = this.validator.showFieldErrors
  isSubmitAttempted: boolean = false

  constructor(private validator: AddEditWaterUnpurifiedValidator, private addEditWaterUnpurifiedService: AddEditWaterUnpurifiedService) { }

  async ngOnInit() {
    // TODO: Get Water unpurified config
  }

  onSubmit(){
    // TODO: 
    // this.addEditWaterUnpurifiedService
  }

  getIsEnabled(){
    return false
  }

  parseFloat(val: string){
    if(val === ""){
      return null
    }
    return parseFloat(val)
  }
}
