import { Routes } from '@angular/router';
import { AppRoutes } from 'src/constants/AppRoutes';
import { AuthGuard } from './user-session.service';

export const routes: Routes = [
  // {
  //   path: '',
  //   redirectTo: 'home',
  //   pathMatch: 'full',
  // },
  {
    path: AppRoutes.Home(),
    loadComponent: () => import('./home/home.page').then( m => m.HomePage),
  },
  {
    path: AppRoutes.Login(),
    loadComponent: () => import('./login/login.page').then( m => m.LoginPage)
  },
  // {
  //   path: AppRoutes.Portal(),
  //   loadComponent: () => import('./portal/portal.page').then( m => m.PortalPage)
  // },
  // {
  //   path: AppRoutes.Organizations(),
  //   loadComponent: () => import('./portal/organizations/organizations.page').then( m => m.OrganizationsPage)
  // },
  // {
  //   path: AppRoutes.Vehicles(),
  //   loadComponent: () => import('./portal/vehicles/vehicles.page').then( m => m.VehiclesPage)
  // },
  // {
  //   path: AppRoutes.VehiclesCreate(),
  //   loadComponent: () => import('./portal/vehicles/create/create.page').then( m => m.CreatePage)
  // },
  // {
  //   path: AppRoutes.OrganizationsCreate(),
  //   loadComponent: () => import('./portal/organizations/create/create.page').then( m => m.CreatePage)
  // },
  {
    path: AppRoutes.Register(),
    loadComponent: () => import('./register/register.page').then( m => m.RegisterPage)
  },
  {
    path: AppRoutes.TermsOfService(),
    loadComponent: () => import('./login/terms-of-service/terms-of-service.page').then( m => m.TermsOfServicePage)
  },
  {
    path: AppRoutes.PrivacyPolicy(),
    loadComponent: () => import('./login/privacy-policy/privacy-policy.page').then( m => m.PrivacyPolicyPage)
  },
  {
    path: AppRoutes.Organization(),
    canActivate: [AuthGuard],
    loadComponent: () => import('./home/dashboard/organization/organization.page').then( m => m.OrganizationPage)
  },
  {
    path: AppRoutes.OrganizationSettings(),
    loadComponent: () => import('./home/dashboard/organization/organization-settings/organization-settings.page').then( m => m.OrganizationSettingsPage)
  },
  {
    path: AppRoutes.VehiclesCreate(),
    canActivate: [AuthGuard],
    loadComponent: () => import('./home/dashboard/organization/add-edit-vehicle/add-edit-vehicle.page').then( m => m.AddEditVehiclePage)
  },
  {
    path: AppRoutes.VehiclesEdit(),
    canActivate: [AuthGuard],
    loadComponent: () => import('./home/dashboard/organization/add-edit-vehicle/add-edit-vehicle.page').then( m => m.AddEditVehiclePage)
  },
  {
    path: 'add-edit-account',
    loadComponent: () => import('./home/dashboard/organization/add-edit-account/add-edit-account.page').then( m => m.AddEditAccountPage)
  },
  {
    path: 'add-edit-order-payment-method',
    loadComponent: () => import('./home/dashboard/organization/add-edit-order-payment-method/add-edit-order-payment-method.page').then( m => m.AddEditOrderPaymentMethodPage)
  },
  {
    path: 'select-location',
    loadComponent: () => import('./shared/select-location/select-location.page').then( m => m.SelectLocationPage)
  },
  {
    path: AppRoutes.Dashboard(),
    canActivate: [AuthGuard],
    loadComponent: () => import('./home/dashboard/dashboard.page').then( m => m.DashboardPage),
  },
  {
    path: AppRoutes.OrganizationsCreate(),
    canActivate: [AuthGuard],
    loadComponent: () => import('./home/dashboard/organization-create/organization-create.page').then( m => m.OrganizationCreatePage)
  },
  {
    path: AppRoutes.OrganizationCreateAccount(),
    canActivate: [AuthGuard],
    loadComponent: () => import('./home/dashboard/organization/add-edit-account/add-edit-account.page').then( m => m.AddEditAccountPage),
  },
  {
    path: AppRoutes.OrganizationEditAccount(),
    canActivate: [AuthGuard],
    loadComponent: () => import('./home/dashboard/organization/add-edit-account/add-edit-account.page').then( m => m.AddEditAccountPage),
  },
  {
    path: AppRoutes.OrganizationAddEditOrderType(),
    canActivate: [AuthGuard],
    loadComponent: () => import('./home/dashboard/organization/add-edit-order-type/add-edit-order-type.page').then( m => m.AddEditOrderTypePage)
  },
  {
    path: AppRoutes.LocationSelector(),
    loadComponent: () => import('./shared/components/location-selector/location-selector.page').then( m => m.LocationSelectorPage)
  },
];
