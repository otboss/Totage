import { Injectable } from '@angular/core';
// import { GoogleSignIn } from '@capawesome/capacitor-google-sign-in';
import { SocialLogin } from '@capgo/capacitor-social-login';
import { createClient } from '@connectrpc/connect';
import { environment } from 'src/environments/environment';
import { LoginViaOAuthCommand, UserService } from '../shared/services/client/modules/user/api/user_pb';

@Injectable({
  providedIn: 'root',
})
export class LoginService {

  async loginWithGoogle(){
    try {

      await SocialLogin.logout({
        provider: "google"
      })

      await SocialLogin.initialize({
        google: {
          webClientId: environment.googleWebAuthConfig.clientId,
          redirectUrl: environment.googleWebAuthConfig.redirectUrl
        }
      })

      const response = await SocialLogin.login({
        provider: 'google',
        options: {
          scopes: ['email', 'profile'],
          // forceRefreshToken: true // Uncomment if you specifically need a refresh token
        },
      })

      // window.close()

      // alert(JSON.stringify(response))

      // await GoogleSignIn.initialize({
      //   clientId: environment.googleWebAuthConfig.clientId,
      //   redirectUrl: environment.googleWebAuthConfig.redirectUrl,
      //   state: ""
      // })

      // Initiates native popup on iOS/Android or web overlay
      // const result = await GoogleSignIn.signIn();
      
      // console.log('Google User Data:', response);
      
      // Extract tokens to validate securely on your server
      // const idToken = result.idToken;
      // const accessToken = result.accessToken;

      // if (idToken) {
      //   // TODO: 
      //   // Send idToken to your custom backend API via HttpClient here
      //   // Example: this.http.post('https://yourdomain.com', { token: idToken })
      // }
      
      return response;
    } catch (error) {
      console.error('Google Sign-In Error:', error);
      throw error;
    }
  }

  async logoutFromGoogle() {
    try {
      await SocialLogin.logout({
        provider: "google"
      })
    } catch (error) {
      console.error('Logout error:', error);
    }
  }

  async getJwtWithGoogleAccessToken(requestBody: LoginViaOAuthCommand) {
      const client = createClient(UserService, environment.connectTransport)
      return await client.loginViaOAuthHandler(requestBody)
  }
  
}
