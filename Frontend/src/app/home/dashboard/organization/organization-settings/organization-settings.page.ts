import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';

import { FormsModule } from '@angular/forms';
import { IonContent, IonHeader, IonTitle, IonToolbar } from '@ionic/angular/standalone';

@Component({
  selector: 'app-organization-settings',
  templateUrl: './organization-settings.page.html',
  styleUrls: ['./organization-settings.page.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonContent, IonHeader, IonTitle, IonToolbar, FormsModule]
})
export class OrganizationSettingsPage implements OnInit {

  constructor() { }

  ngOnInit() {
  }

}
