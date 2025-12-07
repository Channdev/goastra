/*
 * GoAstra Frontend - Auth Guard
 *
 * Protects routes requiring authentication.
 * Redirects unauthenticated users to login.
 */
import { inject } from '@angular/core';
import { Router, CanActivateFn } from '@angular/router';
import { AuthService } from '@core/services/auth.service';

/*
 * Guard that requires user to be authenticated.
 */
export const authGuard: CanActivateFn = (route, state) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  if (authService.isAuthenticated()) {
    return true;
  }

  /* Store intended destination for redirect after login */
  router.navigate(['/auth/login'], {
    queryParams: { returnUrl: state.url }
  });

  return false;
};

/*
 * Guard that requires specific roles.
 */
export const roleGuard = (allowedRoles: string[]): CanActivateFn => {
  return (route, state) => {
    const authService = inject(AuthService);
    const router = inject(Router);

    if (!authService.isAuthenticated()) {
      router.navigate(['/auth/login'], {
        queryParams: { returnUrl: state.url }
      });
      return false;
    }

    if (authService.hasAnyRole(allowedRoles)) {
      return true;
    }

    /* Redirect to unauthorized page or home */
    router.navigate(['/unauthorized']);
    return false;
  };
};

/*
 * Guard that prevents authenticated users from accessing route.
 * Useful for login/register pages.
 */
export const guestGuard: CanActivateFn = () => {
  const authService = inject(AuthService);
  const router = inject(Router);

  if (!authService.isAuthenticated()) {
    return true;
  }

  router.navigate(['/dashboard']);
  return false;
};
