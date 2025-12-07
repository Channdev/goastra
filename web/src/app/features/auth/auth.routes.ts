/*
 * GoAstra Frontend - Auth Routes
 *
 * Route configuration for authentication feature.
 */
import { Routes } from '@angular/router';
import { guestGuard } from '@core/guards/auth.guard';

export const AUTH_ROUTES: Routes = [
  {
    path: '',
    redirectTo: 'login',
    pathMatch: 'full'
  },
  {
    path: 'login',
    loadComponent: () => import('./login/login.component')
      .then(m => m.LoginComponent),
    canActivate: [guestGuard]
  },
  {
    path: 'register',
    loadComponent: () => import('./register/register.component')
      .then(m => m.RegisterComponent),
    canActivate: [guestGuard]
  }
];
