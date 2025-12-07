package frontend

func DefaultStylesCSS() string {
	return `:root {
  --color-primary: #2563eb;
  --color-primary-hover: #1d4ed8;
  --color-background: #0a0a0b;
  --color-surface: #141416;
  --color-surface-hover: #1a1a1d;
  --color-text: #fafafa;
  --color-text-secondary: #a1a1aa;
  --color-border: #27272a;
  --color-success: #22c55e;
  --color-error: #ef4444;
  --radius: 8px;
  --shadow: 0 1px 3px rgba(0,0,0,0.3), 0 1px 2px rgba(0,0,0,0.2);
  --shadow-lg: 0 10px 40px rgba(0,0,0,0.4);
}

* { box-sizing: border-box; margin: 0; padding: 0; }

html { scroll-behavior: smooth; }

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', sans-serif;
  background: var(--color-background);
  color: var(--color-text);
  line-height: 1.6;
  -webkit-font-smoothing: antialiased;
}

a { color: var(--color-primary); text-decoration: none; transition: color 0.2s; }
a:hover { color: var(--color-primary-hover); }

button { font-family: inherit; }

::selection { background: var(--color-primary); color: white; }
`
}

func DefaultAppRoutes() string {
	return `import { Routes } from '@angular/router';

export const routes: Routes = [
  { path: '', redirectTo: 'home', pathMatch: 'full' },
  { path: 'home', loadComponent: () => import('@features/home/home.component').then(m => m.HomeComponent) },
  { path: 'login', loadComponent: () => import('@features/auth/login/login.component').then(m => m.LoginComponent) },
  { path: 'register', loadComponent: () => import('@features/auth/register/register.component').then(m => m.RegisterComponent) },
  { path: 'dashboard', loadComponent: () => import('@features/dashboard/dashboard.component').then(m => m.DashboardComponent) },
  { path: '**', redirectTo: 'home' }
];
`
}

func DefaultHomeComponent() string {
	return `import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [RouterLink],
  template: ` + "`" + `
    <div class="page">
      <header class="header">
        <nav class="nav">
          <a routerLink="/" class="logo">GoAstra</a>
          <div class="nav-links">
            <a routerLink="/login" class="nav-link">Sign In</a>
            <a routerLink="/register" class="nav-btn">Get Started</a>
          </div>
        </nav>
      </header>

      <main class="hero">
        <div class="hero-content">
          <span class="hero-badge">Full-Stack Framework</span>
          <h1>Build production-ready applications with Go and Angular</h1>
          <p>A complete development framework that combines Go's performance with Angular's powerful frontend capabilities. Type-safe, scalable, and ready for production.</p>
          <div class="hero-actions">
            <a routerLink="/register" class="btn btn-primary">Start Building</a>
            <a routerLink="/login" class="btn btn-outline">Sign In</a>
          </div>
        </div>
      </main>

      <section class="features">
        <div class="features-grid">
          <div class="feature-card">
            <div class="feature-header">
              <svg class="feature-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/></svg>
              <h3>High Performance</h3>
            </div>
            <p>Go's compiled speed combined with Angular's optimized rendering delivers exceptional application performance.</p>
          </div>
          <div class="feature-card">
            <div class="feature-header">
              <svg class="feature-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
              <h3>Type Safety</h3>
            </div>
            <p>End-to-end type safety with shared schema definitions between your backend and frontend code.</p>
          </div>
          <div class="feature-card">
            <div class="feature-header">
              <svg class="feature-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
              <h3>Production Ready</h3>
            </div>
            <p>JWT authentication, CORS handling, structured logging, and database integration included out of the box.</p>
          </div>
        </div>
      </section>

      <footer class="footer">
        <div class="footer-content">
          <span class="footer-logo">GoAstra</span>
          <div class="footer-links">
            <a href="https://github.com/channdev/goastra" target="_blank">Documentation</a>
            <a href="https://github.com/channdev/goastra" target="_blank">GitHub</a>
          </div>
        </div>
      </footer>
    </div>
  ` + "`" + `,
  styles: [` + "`" + `
    .page { min-height: 100vh; display: flex; flex-direction: column; }
    .header { position: fixed; top: 0; left: 0; right: 0; z-index: 100; background: rgba(10,10,11,0.8); backdrop-filter: blur(12px); border-bottom: 1px solid var(--color-border); }
    .nav { max-width: 1200px; margin: 0 auto; padding: 1rem 2rem; display: flex; justify-content: space-between; align-items: center; }
    .logo { font-size: 1.25rem; font-weight: 600; color: var(--color-text); letter-spacing: -0.02em; }
    .nav-links { display: flex; align-items: center; gap: 1.5rem; }
    .nav-link { color: var(--color-text-secondary); font-size: 0.875rem; font-weight: 500; }
    .nav-link:hover { color: var(--color-text); }
    .nav-btn { background: var(--color-primary); color: white; padding: 0.5rem 1rem; border-radius: var(--radius); font-size: 0.875rem; font-weight: 500; }
    .nav-btn:hover { background: var(--color-primary-hover); color: white; }

    .hero { flex: 1; display: flex; align-items: center; justify-content: center; padding: 8rem 2rem 4rem; }
    .hero-content { max-width: 720px; text-align: center; }
    .hero-badge { display: inline-block; background: var(--color-surface); border: 1px solid var(--color-border); padding: 0.375rem 0.875rem; border-radius: 50px; font-size: 0.75rem; font-weight: 500; color: var(--color-text-secondary); margin-bottom: 1.5rem; text-transform: uppercase; letter-spacing: 0.05em; }
    .hero h1 { font-size: 3rem; font-weight: 600; line-height: 1.15; letter-spacing: -0.03em; margin-bottom: 1.25rem; }
    .hero p { font-size: 1.125rem; color: var(--color-text-secondary); line-height: 1.7; margin-bottom: 2rem; }
    .hero-actions { display: flex; gap: 0.75rem; justify-content: center; }
    .btn { padding: 0.75rem 1.5rem; border-radius: var(--radius); font-size: 0.9375rem; font-weight: 500; transition: all 0.15s ease; }
    .btn-primary { background: var(--color-primary); color: white; }
    .btn-primary:hover { background: var(--color-primary-hover); color: white; }
    .btn-outline { background: transparent; color: var(--color-text); border: 1px solid var(--color-border); }
    .btn-outline:hover { background: var(--color-surface); border-color: var(--color-text-secondary); }

    .features { padding: 4rem 2rem; background: var(--color-surface); border-top: 1px solid var(--color-border); }
    .features-grid { max-width: 1200px; margin: 0 auto; display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 1.5rem; }
    .feature-card { background: var(--color-background); border: 1px solid var(--color-border); border-radius: var(--radius); padding: 1.5rem; }
    .feature-header { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 0.75rem; }
    .feature-icon { width: 20px; height: 20px; color: var(--color-primary); }
    .feature-card h3 { font-size: 0.9375rem; font-weight: 600; }
    .feature-card p { font-size: 0.875rem; color: var(--color-text-secondary); line-height: 1.6; }

    .footer { padding: 1.5rem 2rem; border-top: 1px solid var(--color-border); }
    .footer-content { max-width: 1200px; margin: 0 auto; display: flex; justify-content: space-between; align-items: center; }
    .footer-logo { font-weight: 600; color: var(--color-text-secondary); }
    .footer-links { display: flex; gap: 1.5rem; }
    .footer-links a { font-size: 0.875rem; color: var(--color-text-secondary); }
    .footer-links a:hover { color: var(--color-text); }

    @media (max-width: 640px) {
      .hero h1 { font-size: 2rem; }
      .hero-actions { flex-direction: column; }
      .footer-content { flex-direction: column; gap: 1rem; text-align: center; }
    }
  ` + "`" + `]
})
export class HomeComponent {}
`
}

func DefaultLoginComponent() string {
	return `import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [RouterLink, FormsModule],
  template: ` + "`" + `
    <div class="auth-page">
      <div class="auth-container">
        <a routerLink="/" class="auth-logo">GoAstra</a>
        <div class="auth-card">
          <div class="auth-header">
            <h1>Sign in</h1>
            <p>Enter your credentials to access your account</p>
          </div>
          <form (ngSubmit)="onSubmit()" class="auth-form">
            <div class="form-field">
              <label for="email">Email</label>
              <input type="email" id="email" [(ngModel)]="email" name="email" placeholder="name@company.com" required autocomplete="email">
            </div>
            <div class="form-field">
              <label for="password">Password</label>
              <input type="password" id="password" [(ngModel)]="password" name="password" placeholder="Enter your password" required autocomplete="current-password">
            </div>
            <button type="submit" class="btn-submit">Sign in</button>
          </form>
          <p class="auth-footer">Don't have an account? <a routerLink="/register">Create account</a></p>
        </div>
      </div>
    </div>
  ` + "`" + `,
  styles: [` + "`" + `
    .auth-page { min-height: 100vh; display: flex; align-items: center; justify-content: center; padding: 2rem; background: var(--color-background); }
    .auth-container { width: 100%; max-width: 400px; }
    .auth-logo { display: block; text-align: center; font-size: 1.25rem; font-weight: 600; color: var(--color-text); margin-bottom: 2rem; }
    .auth-card { background: var(--color-surface); border: 1px solid var(--color-border); border-radius: var(--radius); padding: 2rem; }
    .auth-header { margin-bottom: 1.5rem; }
    .auth-header h1 { font-size: 1.25rem; font-weight: 600; margin-bottom: 0.375rem; }
    .auth-header p { font-size: 0.875rem; color: var(--color-text-secondary); }
    .auth-form { display: flex; flex-direction: column; gap: 1rem; }
    .form-field { display: flex; flex-direction: column; gap: 0.375rem; }
    .form-field label { font-size: 0.875rem; font-weight: 500; color: var(--color-text); }
    .form-field input { padding: 0.625rem 0.875rem; background: var(--color-background); border: 1px solid var(--color-border); border-radius: var(--radius); color: var(--color-text); font-size: 0.875rem; transition: border-color 0.15s; }
    .form-field input::placeholder { color: var(--color-text-secondary); }
    .form-field input:focus { outline: none; border-color: var(--color-primary); }
    .btn-submit { margin-top: 0.5rem; padding: 0.625rem 1rem; background: var(--color-primary); color: white; border: none; border-radius: var(--radius); font-size: 0.875rem; font-weight: 500; cursor: pointer; transition: background 0.15s; }
    .btn-submit:hover { background: var(--color-primary-hover); }
    .auth-footer { text-align: center; margin-top: 1.5rem; font-size: 0.875rem; color: var(--color-text-secondary); }
    .auth-footer a { color: var(--color-primary); font-weight: 500; }
  ` + "`" + `]
})
export class LoginComponent {
  email = '';
  password = '';
  onSubmit() { console.log('Login:', this.email); }
}
`
}

func DefaultRegisterComponent() string {
	return `import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [RouterLink, FormsModule],
  template: ` + "`" + `
    <div class="auth-page">
      <div class="auth-container">
        <a routerLink="/" class="auth-logo">GoAstra</a>
        <div class="auth-card">
          <div class="auth-header">
            <h1>Create account</h1>
            <p>Get started with your free account</p>
          </div>
          <form (ngSubmit)="onSubmit()" class="auth-form">
            <div class="form-field">
              <label for="name">Full name</label>
              <input type="text" id="name" [(ngModel)]="name" name="name" placeholder="Enter your name" required autocomplete="name">
            </div>
            <div class="form-field">
              <label for="email">Email</label>
              <input type="email" id="email" [(ngModel)]="email" name="email" placeholder="name@company.com" required autocomplete="email">
            </div>
            <div class="form-field">
              <label for="password">Password</label>
              <input type="password" id="password" [(ngModel)]="password" name="password" placeholder="Create a password" required autocomplete="new-password">
            </div>
            <button type="submit" class="btn-submit">Create account</button>
          </form>
          <p class="auth-footer">Already have an account? <a routerLink="/login">Sign in</a></p>
        </div>
      </div>
    </div>
  ` + "`" + `,
  styles: [` + "`" + `
    .auth-page { min-height: 100vh; display: flex; align-items: center; justify-content: center; padding: 2rem; background: var(--color-background); }
    .auth-container { width: 100%; max-width: 400px; }
    .auth-logo { display: block; text-align: center; font-size: 1.25rem; font-weight: 600; color: var(--color-text); margin-bottom: 2rem; }
    .auth-card { background: var(--color-surface); border: 1px solid var(--color-border); border-radius: var(--radius); padding: 2rem; }
    .auth-header { margin-bottom: 1.5rem; }
    .auth-header h1 { font-size: 1.25rem; font-weight: 600; margin-bottom: 0.375rem; }
    .auth-header p { font-size: 0.875rem; color: var(--color-text-secondary); }
    .auth-form { display: flex; flex-direction: column; gap: 1rem; }
    .form-field { display: flex; flex-direction: column; gap: 0.375rem; }
    .form-field label { font-size: 0.875rem; font-weight: 500; color: var(--color-text); }
    .form-field input { padding: 0.625rem 0.875rem; background: var(--color-background); border: 1px solid var(--color-border); border-radius: var(--radius); color: var(--color-text); font-size: 0.875rem; transition: border-color 0.15s; }
    .form-field input::placeholder { color: var(--color-text-secondary); }
    .form-field input:focus { outline: none; border-color: var(--color-primary); }
    .btn-submit { margin-top: 0.5rem; padding: 0.625rem 1rem; background: var(--color-primary); color: white; border: none; border-radius: var(--radius); font-size: 0.875rem; font-weight: 500; cursor: pointer; transition: background 0.15s; }
    .btn-submit:hover { background: var(--color-primary-hover); }
    .auth-footer { text-align: center; margin-top: 1.5rem; font-size: 0.875rem; color: var(--color-text-secondary); }
    .auth-footer a { color: var(--color-primary); font-weight: 500; }
  ` + "`" + `]
})
export class RegisterComponent {
  name = '';
  email = '';
  password = '';
  onSubmit() { console.log('Register:', this.email); }
}
`
}

func DefaultDashboardComponent() string {
	return `import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [RouterLink],
  template: ` + "`" + `
    <div class="layout">
      <aside class="sidebar">
        <div class="sidebar-header">
          <a routerLink="/" class="sidebar-logo">GoAstra</a>
        </div>
        <nav class="sidebar-nav">
          <a routerLink="/dashboard" class="nav-item active">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>
            Dashboard
          </a>
          <a routerLink="/dashboard" class="nav-item">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
            Users
          </a>
          <a routerLink="/dashboard" class="nav-item">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
            Settings
          </a>
        </nav>
        <div class="sidebar-footer">
          <a routerLink="/" class="nav-item">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
            Sign out
          </a>
        </div>
      </aside>

      <main class="main">
        <header class="main-header">
          <div>
            <h1>Dashboard</h1>
            <p>Overview of your application metrics</p>
          </div>
        </header>

        <div class="metrics">
          <div class="metric-card">
            <div class="metric-header">
              <span class="metric-label">Total Users</span>
              <svg class="metric-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/></svg>
            </div>
            <div class="metric-value">1,234</div>
            <div class="metric-change positive">+12% from last month</div>
          </div>
          <div class="metric-card">
            <div class="metric-header">
              <span class="metric-label">Active Sessions</span>
              <svg class="metric-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 12h-4l-3 9L9 3l-3 9H2"/></svg>
            </div>
            <div class="metric-value">567</div>
            <div class="metric-change positive">+8% from last hour</div>
          </div>
          <div class="metric-card">
            <div class="metric-header">
              <span class="metric-label">Uptime</span>
              <svg class="metric-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
            </div>
            <div class="metric-value">99.9%</div>
            <div class="metric-change">Last 30 days</div>
          </div>
          <div class="metric-card">
            <div class="metric-header">
              <span class="metric-label">Avg Response</span>
              <svg class="metric-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
            </div>
            <div class="metric-value">12ms</div>
            <div class="metric-change positive">-3ms from average</div>
          </div>
        </div>
      </main>
    </div>
  ` + "`" + `,
  styles: [` + "`" + `
    .layout { display: flex; min-height: 100vh; background: var(--color-background); }

    .sidebar { width: 240px; background: var(--color-surface); border-right: 1px solid var(--color-border); display: flex; flex-direction: column; }
    .sidebar-header { padding: 1.25rem 1rem; border-bottom: 1px solid var(--color-border); }
    .sidebar-logo { font-size: 1.125rem; font-weight: 600; color: var(--color-text); }
    .sidebar-nav { flex: 1; padding: 0.75rem; display: flex; flex-direction: column; gap: 0.25rem; }
    .sidebar-footer { padding: 0.75rem; border-top: 1px solid var(--color-border); }
    .nav-item { display: flex; align-items: center; gap: 0.75rem; padding: 0.625rem 0.75rem; border-radius: var(--radius); color: var(--color-text-secondary); font-size: 0.875rem; font-weight: 500; transition: all 0.15s; }
    .nav-item:hover { background: var(--color-surface-hover); color: var(--color-text); }
    .nav-item.active { background: var(--color-primary); color: white; }
    .nav-item svg { width: 18px; height: 18px; }

    .main { flex: 1; display: flex; flex-direction: column; }
    .main-header { padding: 1.5rem 2rem; border-bottom: 1px solid var(--color-border); }
    .main-header h1 { font-size: 1.25rem; font-weight: 600; margin-bottom: 0.25rem; }
    .main-header p { font-size: 0.875rem; color: var(--color-text-secondary); }

    .metrics { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 1rem; padding: 1.5rem 2rem; }
    .metric-card { background: var(--color-surface); border: 1px solid var(--color-border); border-radius: var(--radius); padding: 1.25rem; }
    .metric-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.75rem; }
    .metric-label { font-size: 0.875rem; color: var(--color-text-secondary); }
    .metric-icon { width: 18px; height: 18px; color: var(--color-text-secondary); }
    .metric-value { font-size: 1.75rem; font-weight: 600; margin-bottom: 0.25rem; }
    .metric-change { font-size: 0.75rem; color: var(--color-text-secondary); }
    .metric-change.positive { color: var(--color-success); }

    @media (max-width: 768px) {
      .sidebar { display: none; }
      .metrics { grid-template-columns: 1fr; }
    }
  ` + "`" + `]
})
export class DashboardComponent {}
`
}
