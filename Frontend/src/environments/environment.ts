// This file can be replaced during build by using the `fileReplacements` array.
// `ng build` replaces `environment.ts` with `environment.prod.ts`.
// The list of file replacements can be found in `angular.json`.

import { createConnectTransport } from "@connectrpc/connect-web";

export const environment = Object.freeze({
  production: false,
  connectTransport: createConnectTransport({
    baseUrl: "http://localhost:8080" 
  }),
  googleWebAuthConfig: {
    clientId: "110285512027-tqjlq4chjt0tghsb3tlrp4nh6elq9g2t.apps.googleusercontent.com",
    redirectUrl: "http://localhost:8100/login"
  }
});

/*
 * For easier debugging in development mode, you can import the following file
 * to ignore zone related error stack frames such as `zone.run`, `zoneDelegate.invokeTask`.
 *
 * This import should be commented out in production mode because it will have a negative impact
 * on performance if an error is thrown.
 */
// import 'zone.js/plugins/zone-error';  // Included with Angular CLI.
