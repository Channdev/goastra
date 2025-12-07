package frontend

func DefaultStylesCSS() string {
	return `:root {
  --color-primary: #3b82f6;
  --color-background: #0f172a;
  --color-surface: #1e293b;
  --color-text: #f8fafc;
  --color-text-muted: #94a3b8;
  --color-border: #334155;
}

* { box-sizing: border-box; margin: 0; padding: 0; }

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: var(--color-background);
  color: var(--color-text);
}

a { color: var(--color-primary); text-decoration: none; }
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
    <div class="landing">
      <nav class="navbar">
        <div class="logo">GoAstra</div>
        <div class="nav-links">
          <a routerLink="/login">Login</a>
          <a routerLink="/register" class="btn-primary">Get Started</a>
        </div>
      </nav>

      <main class="hero">
        <div class="hero-content">
          <h1>Build Full-Stack Apps<br><span class="gradient">Lightning Fast</span></h1>
          <p>GoAstra combines the power of Go backend with Angular frontend.<br>Production-ready, type-safe, and developer-friendly.</p>
          <div class="hero-buttons">
            <a routerLink="/register" class="btn btn-primary">Start Building</a>
            <a href="https://github.com/channdev/goastra" target="_blank" class="btn btn-secondary">View on GitHub</a>
          </div>
        </div>
      </main>

      <section class="features">
        <div class="feature">
          <div class="feature-icon">&#9889;</div>
          <h3>Blazing Fast</h3>
          <p>Go's performance meets Angular's reactivity for lightning-fast apps.</p>
        </div>
        <div class="feature">
          <div class="feature-icon">&#128274;</div>
          <h3>Type Safe</h3>
          <p>End-to-end type safety with shared schemas between frontend and backend.</p>
        </div>
        <div class="feature">
          <div class="feature-icon">&#128640;</div>
          <h3>Production Ready</h3>
          <p>JWT auth, CORS, logging, and database support out of the box.</p>
        </div>
      </section>

      <footer class="footer">
        <p>Built with GoAstra &middot; <a href="https://github.com/channdev/goastra">GitHub</a></p>
      </footer>
    </div>
  ` + "`" + `,
  styles: [` + "`" + `
    .landing { min-height: 100vh; display: flex; flex-direction: column; }
    .navbar { display: flex; justify-content: space-between; align-items: center; padding: 1.5rem 3rem; }
    .logo { font-size: 1.5rem; font-weight: 700; background: linear-gradient(135deg, #3b82f6, #8b5cf6); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
    .nav-links { display: flex; gap: 1.5rem; align-items: center; }
    .nav-links a { color: #94a3b8; transition: color 0.2s; }
    .nav-links a:hover { color: #f8fafc; }
    .btn-primary { background: #3b82f6 !important; color: white !important; padding: 0.5rem 1rem; border-radius: 6px; }
    .hero { flex: 1; display: flex; align-items: center; justify-content: center; text-align: center; padding: 2rem; }
    .hero h1 { font-size: 3.5rem; line-height: 1.1; margin-bottom: 1.5rem; }
    .gradient { background: linear-gradient(135deg, #3b82f6, #8b5cf6, #ec4899); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
    .hero p { color: #94a3b8; font-size: 1.25rem; margin-bottom: 2rem; line-height: 1.6; }
    .hero-buttons { display: flex; gap: 1rem; justify-content: center; }
    .btn { padding: 0.875rem 1.75rem; border-radius: 8px; font-weight: 500; transition: transform 0.2s, box-shadow 0.2s; }
    .btn:hover { transform: translateY(-2px); }
    .btn-secondary { background: #1e293b; color: #f8fafc; border: 1px solid #334155; }
    .features { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 2rem; padding: 4rem 3rem; background: #1e293b; }
    .feature { text-align: center; padding: 2rem; }
    .feature-icon { font-size: 2.5rem; margin-bottom: 1rem; }
    .feature h3 { margin-bottom: 0.5rem; }
    .feature p { color: #94a3b8; }
    .footer { padding: 2rem; text-align: center; color: #64748b; border-top: 1px solid #334155; }
    .footer a { color: #3b82f6; }
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
      <div class="auth-card">
        <h1>Welcome Back</h1>
        <p class="subtitle">Sign in to your account</p>
        <form (ngSubmit)="onSubmit()">
          <div class="form-group">
            <label>Email</label>
            <input type="email" [(ngModel)]="email" name="email" placeholder="you@example.com" required>
          </div>
          <div class="form-group">
            <label>Password</label>
            <input type="password" [(ngModel)]="password" name="password" placeholder="Enter password" required>
          </div>
          <button type="submit" class="btn-submit">Sign In</button>
        </form>
        <p class="switch">Don't have an account? <a routerLink="/register">Sign up</a></p>
      </div>
    </div>
  ` + "`" + `,
  styles: [` + "`" + `
    .auth-page { min-height: 100vh; display: flex; align-items: center; justify-content: center; padding: 2rem; }
    .auth-card { background: #1e293b; padding: 2.5rem; border-radius: 12px; width: 100%; max-width: 400px; }
    h1 { margin-bottom: 0.5rem; }
    .subtitle { color: #94a3b8; margin-bottom: 2rem; }
    .form-group { margin-bottom: 1.25rem; }
    .form-group label { display: block; margin-bottom: 0.5rem; color: #94a3b8; font-size: 0.875rem; }
    .form-group input { width: 100%; padding: 0.75rem; background: #0f172a; border: 1px solid #334155; border-radius: 6px; color: #f8fafc; font-size: 1rem; }
    .form-group input:focus { outline: none; border-color: #3b82f6; }
    .btn-submit { width: 100%; padding: 0.875rem; background: #3b82f6; color: white; border: none; border-radius: 6px; font-size: 1rem; cursor: pointer; margin-top: 0.5rem; }
    .btn-submit:hover { background: #2563eb; }
    .switch { text-align: center; margin-top: 1.5rem; color: #94a3b8; }
    .switch a { color: #3b82f6; }
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
      <div class="auth-card">
        <h1>Create Account</h1>
        <p class="subtitle">Start building with GoAstra</p>
        <form (ngSubmit)="onSubmit()">
          <div class="form-group">
            <label>Name</label>
            <input type="text" [(ngModel)]="name" name="name" placeholder="Your name" required>
          </div>
          <div class="form-group">
            <label>Email</label>
            <input type="email" [(ngModel)]="email" name="email" placeholder="you@example.com" required>
          </div>
          <div class="form-group">
            <label>Password</label>
            <input type="password" [(ngModel)]="password" name="password" placeholder="Create password" required>
          </div>
          <button type="submit" class="btn-submit">Create Account</button>
        </form>
        <p class="switch">Already have an account? <a routerLink="/login">Sign in</a></p>
      </div>
    </div>
  ` + "`" + `,
  styles: [` + "`" + `
    .auth-page { min-height: 100vh; display: flex; align-items: center; justify-content: center; padding: 2rem; }
    .auth-card { background: #1e293b; padding: 2.5rem; border-radius: 12px; width: 100%; max-width: 400px; }
    h1 { margin-bottom: 0.5rem; }
    .subtitle { color: #94a3b8; margin-bottom: 2rem; }
    .form-group { margin-bottom: 1.25rem; }
    .form-group label { display: block; margin-bottom: 0.5rem; color: #94a3b8; font-size: 0.875rem; }
    .form-group input { width: 100%; padding: 0.75rem; background: #0f172a; border: 1px solid #334155; border-radius: 6px; color: #f8fafc; font-size: 1rem; }
    .form-group input:focus { outline: none; border-color: #3b82f6; }
    .btn-submit { width: 100%; padding: 0.875rem; background: #3b82f6; color: white; border: none; border-radius: 6px; font-size: 1rem; cursor: pointer; margin-top: 0.5rem; }
    .btn-submit:hover { background: #2563eb; }
    .switch { text-align: center; margin-top: 1.5rem; color: #94a3b8; }
    .switch a { color: #3b82f6; }
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
    <div class="dashboard">
      <nav class="sidebar">
        <div class="logo">GoAstra</div>
        <div class="nav-items">
          <a routerLink="/dashboard" class="active">Dashboard</a>
          <a routerLink="/dashboard">Users</a>
          <a routerLink="/dashboard">Settings</a>
        </div>
        <a routerLink="/home" class="logout">Logout</a>
      </nav>
      <main class="content">
        <header>
          <h1>Dashboard</h1>
          <p>Welcome back! Here's an overview of your application.</p>
        </header>
        <div class="stats">
          <div class="stat-card"><h3>1,234</h3><p>Total Users</p></div>
          <div class="stat-card"><h3>567</h3><p>Active Today</p></div>
          <div class="stat-card"><h3>89%</h3><p>Uptime</p></div>
          <div class="stat-card"><h3>12ms</h3><p>Avg Response</p></div>
        </div>
      </main>
    </div>
  ` + "`" + `,
  styles: [` + "`" + `
    .dashboard { display: flex; min-height: 100vh; }
    .sidebar { width: 240px; background: #1e293b; padding: 1.5rem; display: flex; flex-direction: column; }
    .logo { font-size: 1.25rem; font-weight: 700; margin-bottom: 2rem; background: linear-gradient(135deg, #3b82f6, #8b5cf6); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
    .nav-items { flex: 1; display: flex; flex-direction: column; gap: 0.5rem; }
    .nav-items a { padding: 0.75rem 1rem; border-radius: 6px; color: #94a3b8; transition: all 0.2s; }
    .nav-items a:hover, .nav-items a.active { background: #334155; color: #f8fafc; }
    .logout { color: #94a3b8; padding: 0.75rem 1rem; }
    .content { flex: 1; padding: 2rem; }
    header { margin-bottom: 2rem; }
    header h1 { margin-bottom: 0.5rem; }
    header p { color: #94a3b8; }
    .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1.5rem; }
    .stat-card { background: #1e293b; padding: 1.5rem; border-radius: 8px; }
    .stat-card h3 { font-size: 2rem; margin-bottom: 0.25rem; }
    .stat-card p { color: #94a3b8; }
  ` + "`" + `]
})
export class DashboardComponent {}
`
}
