/*
 * GoAstra CLI - Frontend REST Templates
 *
 * Generates Angular templates for REST API with HttpClient.
 * Includes API service and HTTP interceptor templates.
 */
package api

// RESTServiceTS returns the REST API service template.
func RESTServiceTS() string {
	return `/*
 * API Service
 *
 * Angular service for REST API operations with HttpClient.
 * Provides typed HTTP methods with error handling.
 */
import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { environment } from '@env/environment';

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface ApiError {
  error: string;
  code?: string;
  details?: any;
}

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private http = inject(HttpClient);
  private baseUrl = environment.apiUrl;

  /**
   * GET request
   */
  get<T>(endpoint: string, params?: Record<string, any>): Observable<T> {
    let httpParams = new HttpParams();
    if (params) {
      Object.keys(params).forEach(key => {
        if (params[key] !== undefined && params[key] !== null) {
          httpParams = httpParams.set(key, params[key].toString());
        }
      });
    }

    return this.http.get<T>(` + "`${this.baseUrl}${endpoint}`" + `, { params: httpParams })
      .pipe(catchError(this.handleError));
  }

  /**
   * POST request
   */
  post<T>(endpoint: string, body: any): Observable<T> {
    return this.http.post<T>(` + "`${this.baseUrl}${endpoint}`" + `, body)
      .pipe(catchError(this.handleError));
  }

  /**
   * PUT request
   */
  put<T>(endpoint: string, body: any): Observable<T> {
    return this.http.put<T>(` + "`${this.baseUrl}${endpoint}`" + `, body)
      .pipe(catchError(this.handleError));
  }

  /**
   * PATCH request
   */
  patch<T>(endpoint: string, body: any): Observable<T> {
    return this.http.patch<T>(` + "`${this.baseUrl}${endpoint}`" + `, body)
      .pipe(catchError(this.handleError));
  }

  /**
   * DELETE request
   */
  delete<T>(endpoint: string): Observable<T> {
    return this.http.delete<T>(` + "`${this.baseUrl}${endpoint}`" + `)
      .pipe(catchError(this.handleError));
  }

  private handleError(error: HttpErrorResponse): Observable<never> {
    let errorMessage = 'An error occurred';

    if (error.error instanceof ErrorEvent) {
      // Client-side error
      errorMessage = error.error.message;
    } else {
      // Server-side error
      errorMessage = error.error?.error || error.message;
    }

    console.error('API error:', errorMessage);
    return throwError(() => ({ error: errorMessage, status: error.status }));
  }
}
`
}

// AuthInterceptorTS returns the HTTP auth interceptor template.
func AuthInterceptorTS() string {
	return `/*
 * Auth Interceptor
 *
 * HTTP interceptor that adds JWT token to requests
 * and handles 401 unauthorized responses.
 */
import { Injectable, inject } from '@angular/core';
import {
  HttpInterceptor,
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpErrorResponse,
} from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { Router } from '@angular/router';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  private router = inject(Router);

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    // Get token from storage
    const token = localStorage.getItem('access_token');

    // Clone request and add auth header if token exists
    if (token) {
      request = request.clone({
        setHeaders: {
          Authorization: ` + "`Bearer ${token}`" + `,
        },
      });
    }

    return next.handle(request).pipe(
      catchError((error: HttpErrorResponse) => {
        if (error.status === 401) {
          // Clear token and redirect to login
          localStorage.removeItem('access_token');
          localStorage.removeItem('refresh_token');
          this.router.navigate(['/auth/login']);
        }
        return throwError(() => error);
      })
    );
  }
}
`
}

// AuthServiceTS returns the auth service template.
func AuthServiceTS() string {
	return `/*
 * Auth Service
 *
 * Handles user authentication, token management,
 * and user session state.
 */
import { Injectable, inject, signal, computed } from '@angular/core';
import { Router } from '@angular/router';
import { Observable, tap, catchError, throwError } from 'rxjs';
import { ApiService } from './api.service';

export interface User {
  id: number;
  email: string;
  name: string;
  role: string;
  active: boolean;
}

export interface AuthResponse {
  token: string;
  refresh_token?: string;
  expires_at: number;
  user: User;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private api = inject(ApiService);
  private router = inject(Router);

  // Reactive state using signals
  private currentUser = signal<User | null>(null);
  private isAuthenticated = signal(false);

  // Public computed values
  user = computed(() => this.currentUser());
  loggedIn = computed(() => this.isAuthenticated());

  constructor() {
    // Check for existing token on init
    this.checkAuth();
  }

  /**
   * Login with email and password
   */
  login(credentials: LoginRequest): Observable<AuthResponse> {
    return this.api.post<AuthResponse>('/api/v1/auth/login', credentials).pipe(
      tap(response => this.handleAuthSuccess(response)),
      catchError(error => {
        this.clearAuth();
        return throwError(() => error);
      })
    );
  }

  /**
   * Register a new user
   */
  register(data: RegisterRequest): Observable<AuthResponse> {
    return this.api.post<AuthResponse>('/api/v1/auth/register', data).pipe(
      tap(response => this.handleAuthSuccess(response)),
      catchError(error => throwError(() => error))
    );
  }

  /**
   * Logout current user
   */
  logout(): void {
    this.api.post('/api/v1/auth/logout', {}).subscribe({
      complete: () => {
        this.clearAuth();
        this.router.navigate(['/auth/login']);
      },
      error: () => {
        this.clearAuth();
        this.router.navigate(['/auth/login']);
      }
    });
  }

  /**
   * Refresh access token
   */
  refreshToken(): Observable<AuthResponse> {
    const refreshToken = localStorage.getItem('refresh_token');
    return this.api.post<AuthResponse>('/api/v1/auth/refresh', { refresh_token: refreshToken }).pipe(
      tap(response => this.handleAuthSuccess(response))
    );
  }

  private handleAuthSuccess(response: AuthResponse): void {
    localStorage.setItem('access_token', response.token);
    if (response.refresh_token) {
      localStorage.setItem('refresh_token', response.refresh_token);
    }
    this.currentUser.set(response.user);
    this.isAuthenticated.set(true);
  }

  private clearAuth(): void {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    this.currentUser.set(null);
    this.isAuthenticated.set(false);
  }

  private checkAuth(): void {
    const token = localStorage.getItem('access_token');
    if (token) {
      // TODO: Validate token and fetch user profile
      this.isAuthenticated.set(true);
    }
  }
}
`
}
