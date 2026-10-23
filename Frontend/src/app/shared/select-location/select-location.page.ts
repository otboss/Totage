import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';

import { FormsModule } from '@angular/forms';
import { IonContent, IonHeader, IonTitle, IonToolbar } from '@ionic/angular/standalone';

@Component({
  selector: 'app-select-location',
  templateUrl: './select-location.page.html',
  styleUrls: ['./select-location.page.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonContent, IonHeader, IonTitle, IonToolbar, FormsModule]
})
export class SelectLocationPage implements OnInit {

  constructor() { }

  ngOnInit() {
  }

}
