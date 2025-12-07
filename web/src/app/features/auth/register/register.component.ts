/*
 * GoAstra Frontend - Register Component
 *
 * User registration form with validation.
 */
import { Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '@core/services/auth.service';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterLink],
  template: `
    <div class="auth-container">
      <div class="auth-card card">
        <div class="card-header">
          <h1>Create Account</h1>
          <p>Join GoAstra and start building.</p>
        </div>

        <div class="card-body">
          @if (error()) {
            <div class="alert alert-error">{{ error() }}</div>
          }

          <form [formGroup]="form" (ngSubmit)="onSubmit()">
            <div class="form-group">
              <label for="name" class="form-label">Name</label>
              <input
                type="text"
                id="name"
                formControlName="name"
                class="form-input"
                placeholder="Your name"
              >
              @if (form.get('name')?.touched && form.get('name')?.errors) {
                <span class="form-error">Name is required</span>
              }
            </div>

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
                placeholder="Create a strong password"
              >
              @if (form.get('password')?.touched && form.get('password')?.errors) {
                <span class="form-error">Password must be at least 8 characters</span>
              }
            </div>

            <div class="form-group">
              <label for="confirmPassword" class="form-label">Confirm Password</label>
              <input
                type="password"
                id="confirmPassword"
                formControlName="confirmPassword"
                class="form-input"
                placeholder="Confirm your password"
              >
              @if (form.get('confirmPassword')?.touched && !passwordsMatch()) {
                <span class="form-error">Passwords do not match</span>
              }
            </div>

            <button
              type="submit"
              class="btn btn-primary btn-full"
              [disabled]="form.invalid || !passwordsMatch() || loading()"
            >
              @if (loading()) {
                <span class="spinner"></span>
              } @else {
                Create Account
              }
            </button>
          </form>
        </div>

        <div class="card-footer">
          <p>Already have an account? <a routerLink="/auth/login">Sign in</a></p>
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
export class RegisterComponent {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);
  private router = inject(Router);

  loading = signal(false);
  error = signal<string | null>(null);

  form = this.fb.group({
    name: ['', [Validators.required]],
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(8)]],
    confirmPassword: ['', [Validators.required]]
  });

  passwordsMatch(): boolean {
    const password = this.form.get('password')?.value;
    const confirmPassword = this.form.get('confirmPassword')?.value;
    return password === confirmPassword;
  }

  onSubmit(): void {
    if (this.form.invalid || !this.passwordsMatch()) return;

    this.loading.set(true);
    this.error.set(null);

    const { name, email, password } = this.form.value;

    this.authService.register({
      name: name!,
      email: email!,
      password: password!
    }).subscribe({
      next: () => {
        this.router.navigate(['/dashboard']);
      },
      error: (err) => {
        this.error.set(err.message || 'Registration failed');
        this.loading.set(false);
      }
    });
  }
}
