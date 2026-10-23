import { Component, HostListener, OnInit, ChangeDetectionStrategy } from '@angular/core';

import { FormsModule } from '@angular/forms';
import { NavController } from '@ionic/angular';
import { IonContent, IonHeader, IonTitle, IonToolbar, IonInput, IonButton, IonButtons, IonBackButton, IonIcon } from '@ionic/angular/standalone';
import { addIcons } from 'ionicons';
import { logoGoogle } from 'ionicons/icons';
import { AppRoutes } from 'src/constants/AppRoutes';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { Assets } from 'src/constants/Assets';
import { LoginService } from './login.service';
import { UserSessionService } from '../user-session.service';
import { create } from '@bufbuild/protobuf';
import { LoginViaOAuthCommandSchema } from '../shared/services/client/modules/user/api/user_pb';

@Component({
  selector: 'app-login',
  templateUrl: './login.page.html',
  styleUrls: ['./login.page.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonContent, IonHeader, IonTitle, IonToolbar, FormsModule, IonInput, IonButton, IonButtons, IonBackButton, IonIcon, RouterLink]
})
export class LoginPage implements OnInit {

  readonly termsofServicePage:string  = AppRoutes.TermsOfService()
  readonly privacyPolicyPage:string  = AppRoutes.PrivacyPolicy()
  readonly googleLogo: string = Assets.Images.google
  loginWithGoogle = this.loginService.loginWithGoogle

  constructor(private router: Router, private loginService: LoginService, private route: ActivatedRoute, private userSessionService: UserSessionService) {
    addIcons({ logoGoogle })
  }

  ngOnInit() {
    if (this.userSessionService.isLoggedIn()){
      this.router.navigateByUrl(AppRoutes.Home(), {replaceUrl: true})
      return
    }

    this.route.fragment.subscribe(async (fragment: string | null) => {
      if (fragment) {
        const params = new URLSearchParams(fragment);

        if(params.get("access_token") === null){
          return
        }

        const JwtTokenResponse = await this.loginService.getJwtWithGoogleAccessToken(
          create(LoginViaOAuthCommandSchema, {
            accessToken: params.get("access_token")
          })
        )

        this.userSessionService.saveAuthToken({
          jwtToken: JwtTokenResponse.jwtToken,
          expires: new Date(Date.now() + parseInt(params.get("expires_in")) * 1000),
          flow: "GOOGLE"
        })

        window.close()
      }
    })
  }

  @HostListener('window:storage', ['$event'])
  onStorageChange(event: StorageEvent) {
    if (event.key === 'auth_token') {
      this.router.navigateByUrl(AppRoutes.Home(), {replaceUrl: true})
    }
  }

  onSubmit(){
    
  }

  goToTermsOfServicePage(){
    this.router.navigateByUrl(this.termsofServicePage)
  }

  goToPrivacyPolicyPage(){
    this.router.navigateByUrl(this.privacyPolicyPage)
  }
}
