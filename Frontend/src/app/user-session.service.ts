import { Injectable } from '@angular/core';
import { inject } from '@angular/core';
import { CanActivateChildFn, Router } from '@angular/router';
import { AppRoutes } from 'src/constants/AppRoutes';

@Injectable({
  providedIn: 'root',
})
export class UserSessionService {
  public saveAuthToken(authToken: AuthToken){
    localStorage.setItem(ItemKeys.auth_token, JSON.stringify(authToken))
  }

  public getAuthToken(): AuthToken | null {
    return JSON.parse(localStorage.getItem(ItemKeys.auth_token))
  }

  public deleteAuthToken() {
    return localStorage.removeItem(ItemKeys.auth_token)
  }

  public isLoggedIn(): boolean {
    const authToken: AuthToken | null = JSON.parse(localStorage.getItem(ItemKeys.auth_token))
    return authToken !== null && new Date() < new Date(authToken.expires)
  }
}


enum ItemKeys {
  auth_token = "auth_token"
}

interface AuthToken {
  jwtToken: string,
  expires: Date
  flow: "GOOGLE" | "PASSWORD"
}

export const AuthGuard: CanActivateChildFn = () => {
  const authService = inject(UserSessionService);
  const router = inject(Router);
  return authService.isLoggedIn() || router.parseUrl(AppRoutes.Home());
};