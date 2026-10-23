import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';

import { FormsModule } from '@angular/forms';
import { IonContent, IonHeader, IonTitle, IonToolbar } from '@ionic/angular/standalone';

@Component({
  selector: 'app-terms-of-service',
  templateUrl: './terms-of-service.page.html',
  styleUrls: ['./terms-of-service.page.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonContent, IonHeader, IonTitle, IonToolbar, FormsModule]
})
export class TermsOfServicePage implements OnInit {

  constructor() { }

  ngOnInit() {
  }

}
