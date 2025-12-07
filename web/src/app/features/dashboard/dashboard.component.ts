/*
 * GoAstra Frontend - Dashboard Component
 *
 * Main dashboard for authenticated users.
 */
import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { AuthService } from '@core/services/auth.service';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, RouterLink],
  template: `
    <div class="dashboard">
      <header class="dashboard-header">
        <h1>Dashboard</h1>
        <div class="header-actions">
          <span class="user-info">{{ authService.currentUser()?.email }}</span>
          <button (click)="logout()" class="btn btn-secondary">Logout</button>
        </div>
      </header>

      <main class="dashboard-content">
        <div class="welcome-card card">
          <div class="card-body">
            <h2>Welcome, {{ authService.currentUser()?.name }}!</h2>
            <p>You are logged in as <strong>{{ authService.currentUser()?.role }}</strong></p>
          </div>
        </div>

        <div class="stats-grid">
          <div class="stat-card card">
            <div class="card-body">
              <h3>Projects</h3>
              <p class="stat-value">0</p>
            </div>
          </div>
          <div class="stat-card card">
            <div class="card-body">
              <h3>Active Tasks</h3>
              <p class="stat-value">0</p>
            </div>
          </div>
          <div class="stat-card card">
            <div class="card-body">
              <h3>Completed</h3>
              <p class="stat-value">0</p>
            </div>
          </div>
        </div>

        <div class="quick-actions card">
          <div class="card-header">
            <h3>Quick Actions</h3>
          </div>
          <div class="card-body">
            <div class="action-buttons">
              <button class="btn btn-primary">Create Project</button>
              <button class="btn btn-secondary">View Profile</button>
              <button class="btn btn-secondary">Settings</button>
            </div>
          </div>
        </div>
      </main>
    </div>
  `,
  styles: [`
    .dashboard {
      min-height: 100vh;
      background: var(--color-surface);
    }

    .dashboard-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 1rem 2rem;
      background: var(--color-background);
      border-bottom: 1px solid var(--color-border);
    }

    .header-actions {
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .user-info {
      color: var(--color-text-muted);
    }

    .dashboard-content {
      padding: 2rem;
      max-width: 1200px;
      margin: 0 auto;
    }

    .welcome-card {
      margin-bottom: 2rem;
    }

    .welcome-card h2 {
      margin-bottom: 0.5rem;
    }

    .stats-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 1rem;
      margin-bottom: 2rem;
    }

    .stat-card {
      text-align: center;
    }

    .stat-value {
      font-size: 2rem;
      font-weight: bold;
      color: var(--color-primary);
      margin-top: 0.5rem;
    }

    .action-buttons {
      display: flex;
      gap: 1rem;
      flex-wrap: wrap;
    }
  `]
})
export class DashboardComponent {
  authService = inject(AuthService);

  logout(): void {
    this.authService.logout();
  }
}
