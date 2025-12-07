/*
 * GoAstra Frontend - Not Found Component
 *
 * 404 error page for unmatched routes.
 */
import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-not-found',
  standalone: true,
  imports: [RouterLink],
  template: `
    <div class="not-found-container">
      <div class="not-found-content">
        <h1>404</h1>
        <h2>Page Not Found</h2>
        <p>The page you're looking for doesn't exist or has been moved.</p>
        <a routerLink="/" class="btn btn-primary">Go Home</a>
      </div>
    </div>
  `,
  styles: [`
    .not-found-container {
      min-height: 100vh;
      display: flex;
      align-items: center;
      justify-content: center;
      background: var(--color-surface);
    }

    .not-found-content {
      text-align: center;
      padding: 2rem;
    }

    h1 {
      font-size: 6rem;
      color: var(--color-primary);
      margin: 0;
      line-height: 1;
    }

    h2 {
      font-size: 1.5rem;
      margin: 1rem 0;
    }

    p {
      color: var(--color-text-muted);
      margin-bottom: 2rem;
    }
  `]
})
export class NotFoundComponent {}
