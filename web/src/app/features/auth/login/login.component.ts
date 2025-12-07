/*
 * GoAstra Frontend - Login Component
 *
 * User login form with validation.
 */
import { Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, Validators } from '@angular/forms';
import { Router, RouterLink, ActivatedRoute } from '@angular/router';
import { AuthService } from '@core/services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterLink],
  template: `
    <div class="auth-container">
      <div class="auth-card card">
        <div class="card-header">
          <h1>Sign In</h1>
          <p>Welcome back! Please enter your credentials.</p>
        </div>

        <div class="card-body">
          @if (error()) {
            <div class="alert alert-error">{{ error() }}</div>
          }

          <form [formGroup]="form" (ngSubmit)="onSubmit()">
            <div class="form-group">
              <label for="email" class="form-label">Email</label>
              <input
                type="email"
                id="email"
                formControlName="email"
                class="form-input"
                placeholder="you@example.com"
              >
              @if (form.get('email')?.touched && form.get('email')?.errors) {
                <span class="form-error">Please enter a valid email</span>
              }
            </div>

            <div class="form-group">
              <label for="password" class="form-label">Password</label>
              <input
                type="password"
                id="password"
                formControlName="password"
                class="form-input"
                placeholder="Enter your password"
              >
              @if (form.get('password')?.touched && form.get('password')?.errors) {
                <span class="form-error">Password is required</span>
              }
            </div>

            <button
              type="submit"
              class="btn btn-primary btn-full"
              [disabled]="form.invalid || loading()"
            >
              @if (loading()) {
                <span class="spinner"></span>
              } @else {
                Sign In
              }
            </button>
          </form>
        </div>

        <div class="card-footer">
          <p>Don't have an account? <a routerLink="/auth/register">Create one</a></p>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .auth-container {
      min-height: 100vh;
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 1rem;
      background: var(--color-surface);
    }

    .auth-card {
      width: 100%;
      max-width: 400px;
    }

    .card-header {
      text-align: center;
    }

    .card-header h1 {
      margin-bottom: 0.5rem;
    }

    .card-header p {
      color: var(--color-text-muted);
    }

    .card-footer {
      text-align: center;
    }

    .btn-full {
      width: 100%;
      margin-top: 1rem;
    }
  `]
})
export class LoginComponent {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);
  private router = inject(Router);
  private route = inject(ActivatedRoute);

  loading = signal(false);
  error = signal<string | null>(null);

  form = this.fb.group({
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required]]
  });

  onSubmit(): void {
    if (this.form.invalid) return;

    this.loading.set(true);
    this.error.set(null);

    const { email, password } = this.form.value;

    this.authService.login({ email: email!, password: password! }).subscribe({
      next: () => {
        const returnUrl = this.route.snapshot.queryParams['returnUrl'] || '/dashboard';
        this.router.navigateByUrl(returnUrl);
      },
      error: (err) => {
        this.error.set(err.message || 'Login failed');
        this.loading.set(false);
      }
    });
  }
}
