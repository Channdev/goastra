/*
 * GoAstra Frontend - Error Interceptor
 *
 * Global HTTP error handling.
 * Handles token refresh and redirects on authentication errors.
 */
import { HttpInterceptorFn, HttpErrorResponse } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, throwError, switchMap } from 'rxjs';
import { AuthService } from '@core/services/auth.service';

/*
 * Interceptor that handles HTTP errors globally.
 */
export const errorInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  return next(req).pipe(
    catchError((error: HttpErrorResponse) => {
      /* Handle 401 Unauthorized */
      if (error.status === 401) {
        /* Skip if this is a refresh or login request */
        if (req.url.includes('/auth/')) {
          return throwError(() => error);
        }

        /* Try to refresh the token */
        const refreshToken = authService.getRefreshToken();
        if (refreshToken) {
          return authService.refreshToken().pipe(
            switchMap(() => {
              /* Retry the original request with new token */
              const token = authService.getToken();
              const retryReq = req.clone({
                setHeaders: {
                  Authorization: `Bearer ${token}`
                }
              });
              return next(retryReq);
            }),
            catchError((refreshError) => {
              /* Refresh failed, logout user */
              authService.logout();
              return throwError(() => refreshError);
            })
          );
        }

        /* No refresh token, logout */
        authService.logout();
      }

      /* Handle 403 Forbidden */
      if (error.status === 403) {
        router.navigate(['/unauthorized']);
      }

      /* Handle 404 Not Found */
      if (error.status === 404) {
        /* Could navigate to not-found page or handle differently */
      }

      /* Handle 500 Server Error */
      if (error.status >= 500) {
        console.error('Server error:', error);
        /* Could show global error notification */
      }

      /* Format error message */
      const errorMessage = extractErrorMessage(error);

      return throwError(() => ({
        status: error.status,
        message: errorMessage,
        original: error
      }));
    })
  );
};

/*
 * Extracts user-friendly error message from response.
 */
function extractErrorMessage(error: HttpErrorResponse): string {
  if (error.error instanceof ErrorEvent) {
    /* Client-side error */
    return error.error.message;
  }

  /* Server-side error */
  if (error.error?.message) {
    return error.error.message;
  }

  if (error.error?.error) {
    return error.error.error;
  }

  /* Default messages based on status */
  switch (error.status) {
    case 400:
      return 'Invalid request. Please check your input.';
    case 401:
      return 'Please login to continue.';
    case 403:
      return 'You do not have permission to perform this action.';
    case 404:
      return 'The requested resource was not found.';
    case 422:
      return 'Validation failed. Please check your input.';
    case 500:
      return 'An unexpected error occurred. Please try again later.';
    default:
      return 'An error occurred. Please try again.';
  }
}
