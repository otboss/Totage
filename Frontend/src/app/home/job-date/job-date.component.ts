import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { IonDatetime, IonGrid, IonCol, IonRow, IonButton, IonChip, IonIcon } from "@ionic/angular/standalone";
import { addIcons } from 'ionicons';
import { chevronBack, chevronBackOutline, chevronForward, chevronForwardOutline } from 'ionicons/icons';

@Component({
  selector: 'app-job-date',
  templateUrl: './job-date.component.html',
  styleUrls: ['./job-date.component.scss'],
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonDatetime, IonGrid, IonCol, IonRow, IonButton, IonChip, IonIcon],
})
export class JobDateComponent  implements OnInit {

  constructor() { 
    addIcons({ chevronBackOutline, chevronForwardOutline });
  }

  ngOnInit() {}

}
