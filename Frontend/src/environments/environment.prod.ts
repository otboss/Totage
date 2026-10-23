import { createConnectTransport } from "@connectrpc/connect-web";

export const environment = Object.freeze({
  production: true,
  connectTransport: createConnectTransport({
    baseUrl: "<< TODO: PRODUCTION BACKEND URL >>" 
  }),
  googleWebAuthConfig: {
    clientId: "110285512027-tqjlq4chjt0tghsb3tlrp4nh6elq9g2t.apps.googleusercontent.com",
    redirectUrl: "<< TODO: PRODUCTION FRONTEND URL >>"
  }
});
