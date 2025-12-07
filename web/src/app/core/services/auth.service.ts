/*
 * GoAstra Frontend - Auth Service
 *
 * Handles authentication state, login, logout, and token management.
 * Uses signals for reactive state management.
 */
import { Injectable, signal, computed } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable, tap, BehaviorSubject } from 'rxjs';
import { environment } from '@env/environment';

/*
 * User model representing authenticated user.
 */
export interface User {
  id: number;
  email: string;
  name: string;
  role: string;
}

/*
 * Authentication response from login/register endpoints.
 */
export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
  user: User;
}

/*
 * Login request payload.
 */
export interface LoginRequest {
  email: string;
  password: string;
}

/*
 * Registration request payload.
 */
export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private readonly TOKEN_KEY = 'goastra_token';
  private readonly REFRESH_KEY = 'goastra_refresh';
  private readonly USER_KEY = 'goastra_user';

  private readonly baseUrl = `${environment.apiUrl}/auth`;

  /* Reactive state using signals */
  private userSignal = signal<User | null>(null);

  /* Public computed values */
  readonly currentUser = this.userSignal.asReadonly();
  readonly isAuthenticated = computed(() => this.userSignal() !== null);

  constructor(
    private http: HttpClient,
    private router: Router
  ) {
    this.loadStoredUser();
  }

  /*
   * Authenticates user with email and password.
   */
  login(credentials: LoginRequest): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(`${this.baseUrl}/login`, credentials)
      .pipe(tap(response => this.handleAuthResponse(response)));
  }

  /*
   * Registers a new user account.
   */
  register(data: RegisterRequest): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(`${this.baseUrl}/register`, data)
      .pipe(tap(response => this.handleAuthResponse(response)));
  }

  /*
   * Refreshes the access token using refresh token.
   */
  refreshToken(): Observable<AuthResponse> {
    const refreshToken = this.getRefreshToken();

    return this.http.post<AuthResponse>(`${this.baseUrl}/refresh`, {
      refresh_token: refreshToken
    }).pipe(tap(response => this.handleAuthResponse(response)));
  }

  /*
   * Logs out the current user.
   */
  logout(): void {
    const token = this.getToken();

    if (token) {
      /* Notify backend of logout */
      this.http.post(`${this.baseUrl}/logout`, {}).subscribe({
        error: () => { /* Ignore logout errors */ }
      });
    }

    this.clearAuth();
    this.router.navigate(['/auth/login']);
  }

  /*
   * Returns the current access token.
   */
  getToken(): string | null {
    return localStorage.getItem(this.TOKEN_KEY);
  }

  /*
   * Returns the current refresh token.
   */
  getRefreshToken(): string | null {
    return localStorage.getItem(this.REFRESH_KEY);
  }

  /*
   * Checks if user has a specific role.
   */
  hasRole(role: string): boolean {
    const user = this.userSignal();
    return user?.role === role;
  }

  /*
   * Checks if user has any of the specified roles.
   */
  hasAnyRole(roles: string[]): boolean {
    const user = this.userSignal();
    return user ? roles.includes(user.role) : false;
  }

  private handleAuthResponse(response: AuthResponse): void {
    localStorage.setItem(this.TOKEN_KEY, response.access_token);
    localStorage.setItem(this.REFRESH_KEY, response.refresh_token);
    localStorage.setItem(this.USER_KEY, JSON.stringify(response.user));

    this.userSignal.set(response.user);
  }

  private loadStoredUser(): void {
    const stored = localStorage.getItem(this.USER_KEY);
    const token = localStorage.getItem(this.TOKEN_KEY);

    if (stored && token) {
      try {
        const user = JSON.parse(stored) as User;
        this.userSignal.set(user);
      } catch {
        this.clearAuth();
      }
    }
  }

  private clearAuth(): void {
    localStorage.removeItem(this.TOKEN_KEY);
    localStorage.removeItem(this.REFRESH_KEY);
    localStorage.removeItem(this.USER_KEY);
    this.userSignal.set(null);
  }
}
