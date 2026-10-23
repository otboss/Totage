import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { IonGrid, IonRow, IonInput, IonButton, IonCol } from "@ionic/angular/standalone";
import { Assets } from 'src/constants/Assets';

@Component({
  selector: 'app-location',
  templateUrl: './location.component.html',
  styleUrls: ['./location.component.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonGrid, IonRow, IonInput, IonButton, IonCol]
})
export class LocationComponent  implements OnInit {

  mapImage: string = Assets.Images.map
  
  constructor() { }

  ngOnInit() {}

}
