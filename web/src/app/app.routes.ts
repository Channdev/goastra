/*
 * GoAstra Frontend - Route Configuration
 *
 * Defines application routes with lazy loading and guards.
 */
import { Routes } from '@angular/router';
import { authGuard } from '@core/guards/auth.guard';

export const routes: Routes = [
  {
    path: '',
    redirectTo: 'home',
    pathMatch: 'full'
  },
  {
    path: 'home',
    loadComponent: () => import('@features/home/home.component')
      .then(m => m.HomeComponent)
  },
  {
    path: 'auth',
    loadChildren: () => import('@features/auth/auth.routes')
      .then(m => m.AUTH_ROUTES)
  },
  {
    path: 'dashboard',
    loadComponent: () => import('@features/dashboard/dashboard.component')
      .then(m => m.DashboardComponent),
    canActivate: [authGuard]
  },
  {
    path: '**',
    loadComponent: () => import('@features/not-found/not-found.component')
      .then(m => m.NotFoundComponent)
  }
];
