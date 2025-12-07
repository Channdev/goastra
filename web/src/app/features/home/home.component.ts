import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [RouterLink],
  template: `
    <div class="home-container">
      <nav class="navbar">
        <div class="logo">GoAstra</div>
        <div class="nav-links">
          <a href="https://github.com/channdev/goastra" target="_blank">GitHub</a>
          <a routerLink="/auth/login" class="btn btn-outline">Sign In</a>
          <a routerLink="/auth/register" class="btn btn-primary">Get Started</a>
        </div>
      </nav>

      <header class="hero">
        <div class="hero-badge">Open Source Framework</div>
        <h1>Build Full-Stack Apps with <span class="gradient-text">GoAstra</span></h1>
        <p class="hero-subtitle">
          A production-ready framework combining Go backend power with Angular frontend elegance.
          Type-safe, fast, and developer-friendly.
        </p>
        <div class="cta-buttons">
          <a routerLink="/auth/register" class="btn btn-primary btn-lg">
            Start Building
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M5 12h14M12 5l7 7-7 7"/>
            </svg>
          </a>
          <a href="https://github.com/channdev/goastra" target="_blank" class="btn btn-outline btn-lg">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
            </svg>
            View on GitHub
          </a>
        </div>
        <div class="hero-stats">
          <div class="stat">
            <span class="stat-value">Go + Angular</span>
            <span class="stat-label">Tech Stack</span>
          </div>
          <div class="stat">
            <span class="stat-value">TypeSync</span>
            <span class="stat-label">Auto Type Generation</span>
          </div>
          <div class="stat">
            <span class="stat-value">CLI</span>
            <span class="stat-label">Code Generation</span>
          </div>
        </div>
      </header>

      <section class="features">
        <h2>Everything You Need</h2>
        <div class="features-grid">
          <div class="feature-card">
            <div class="feature-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
              </svg>
            </div>
            <h3>Go Backend</h3>
            <p>High-performance REST API with Gin framework, JWT authentication, PostgreSQL, and structured logging.</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polygon points="12 2 22 8.5 22 15.5 12 22 2 15.5 2 8.5 12 2"/>
                <line x1="12" y1="22" x2="12" y2="15.5"/>
                <polyline points="22 8.5 12 15.5 2 8.5"/>
              </svg>
            </div>
            <h3>Angular Frontend</h3>
            <p>Modern TypeScript SPA with standalone components, signals, lazy loading, and reactive forms.</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="16 18 22 12 16 6"/>
                <polyline points="8 6 2 12 8 18"/>
              </svg>
            </div>
            <h3>Code Generation</h3>
            <p>Auto-generate TypeScript interfaces from Go structs. Scaffold CRUD operations with a single command.</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                <line x1="3" y1="9" x2="21" y2="9"/>
                <line x1="9" y1="21" x2="9" y2="9"/>
              </svg>
            </div>
            <h3>CLI Tools</h3>
            <p>Powerful CLI for project scaffolding, development servers, builds, and code generation.</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
              </svg>
            </div>
            <h3>Authentication</h3>
            <p>Built-in JWT authentication with refresh tokens, role-based access control, and secure password hashing.</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <line x1="2" y1="12" x2="22" y2="12"/>
                <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
              </svg>
            </div>
            <h3>Production Ready</h3>
            <p>Environment management, database migrations, structured logging, and optimized builds included.</p>
          </div>
        </div>
      </section>

      <section class="quickstart">
        <h2>Get Started in Seconds</h2>
        <div class="code-block">
          <div class="code-header">
            <span class="dot red"></span>
            <span class="dot yellow"></span>
            <span class="dot green"></span>
            <span class="code-title">Terminal</span>
          </div>
          <pre><code><span class="prompt">$</span> go install github.com/channdev/goastra/cli&#64;latest
<span class="prompt">$</span> goastra new my-app
<span class="prompt">$</span> cd my-app
<span class="prompt">$</span> goastra dev

<span class="output">Starting GoAstra development servers...</span>
<span class="success">[Backend] Starting on port 8080...</span>
<span class="info">[Frontend] Starting on port 4200...</span></code></pre>
        </div>
      </section>

      <footer class="footer">
        <div class="footer-content">
          <div class="footer-brand">
            <span class="logo">GoAstra</span>
            <p>Built with passion by <a href="https://github.com/channdev" target="_blank">channdev</a></p>
          </div>
          <div class="footer-links">
            <a href="https://github.com/channdev/goastra" target="_blank">GitHub</a>
            <a href="https://github.com/channdev/goastra/issues" target="_blank">Issues</a>
            <a href="https://github.com/channdev/goastra#readme" target="_blank">Documentation</a>
          </div>
        </div>
        <div class="footer-bottom">
          <p>MIT License - Open Source</p>
        </div>
      </footer>
    </div>
  `,
  styles: [`
    .home-container {
      min-height: 100vh;
      background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
      color: #f8fafc;
    }

    .navbar {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 1.5rem 4rem;
      max-width: 1400px;
      margin: 0 auto;
    }

    .logo {
      font-size: 1.5rem;
      font-weight: 700;
      color: #f8fafc;
    }

    .nav-links {
      display: flex;
      align-items: center;
      gap: 2rem;
    }

    .nav-links a {
      color: #94a3b8;
      text-decoration: none;
      transition: color 0.2s;
    }

    .nav-links a:hover {
      color: #f8fafc;
    }

    .hero {
      text-align: center;
      padding: 6rem 2rem;
      max-width: 900px;
      margin: 0 auto;
    }

    .hero-badge {
      display: inline-block;
      padding: 0.5rem 1rem;
      background: rgba(59, 130, 246, 0.1);
      border: 1px solid rgba(59, 130, 246, 0.3);
      border-radius: 100px;
      color: #60a5fa;
      font-size: 0.875rem;
      margin-bottom: 1.5rem;
    }

    .hero h1 {
      font-size: 4rem;
      font-weight: 800;
      line-height: 1.1;
      margin-bottom: 1.5rem;
    }

    .gradient-text {
      background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }

    .hero-subtitle {
      font-size: 1.25rem;
      color: #94a3b8;
      max-width: 600px;
      margin: 0 auto 2.5rem;
      line-height: 1.7;
    }

    .cta-buttons {
      display: flex;
      gap: 1rem;
      justify-content: center;
      flex-wrap: wrap;
      margin-bottom: 4rem;
    }

    .btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.75rem 1.5rem;
      border-radius: 8px;
      font-weight: 600;
      text-decoration: none;
      transition: all 0.2s;
      border: none;
      cursor: pointer;
    }

    .btn-primary {
      background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
      color: white;
    }

    .btn-primary:hover {
      transform: translateY(-2px);
      box-shadow: 0 10px 40px rgba(59, 130, 246, 0.3);
    }

    .btn-outline {
      background: transparent;
      border: 1px solid #334155;
      color: #f8fafc;
    }

    .btn-outline:hover {
      background: #1e293b;
      border-color: #475569;
    }

    .btn-lg {
      padding: 1rem 2rem;
      font-size: 1.1rem;
    }

    .hero-stats {
      display: flex;
      justify-content: center;
      gap: 4rem;
      padding-top: 2rem;
      border-top: 1px solid #334155;
    }

    .stat {
      text-align: center;
    }

    .stat-value {
      display: block;
      font-size: 1.25rem;
      font-weight: 700;
      color: #f8fafc;
    }

    .stat-label {
      font-size: 0.875rem;
      color: #64748b;
    }

    .features {
      padding: 6rem 2rem;
      max-width: 1200px;
      margin: 0 auto;
    }

    .features h2 {
      text-align: center;
      font-size: 2.5rem;
      margin-bottom: 3rem;
    }

    .features-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
      gap: 1.5rem;
    }

    .feature-card {
      background: rgba(30, 41, 59, 0.5);
      border: 1px solid #334155;
      border-radius: 12px;
      padding: 2rem;
      transition: all 0.3s;
    }

    .feature-card:hover {
      border-color: #3b82f6;
      transform: translateY(-4px);
    }

    .feature-icon {
      width: 56px;
      height: 56px;
      background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      margin-bottom: 1.5rem;
    }

    .feature-icon svg {
      color: white;
    }

    .feature-card h3 {
      font-size: 1.25rem;
      margin-bottom: 0.75rem;
    }

    .feature-card p {
      color: #94a3b8;
      line-height: 1.6;
    }

    .quickstart {
      padding: 6rem 2rem;
      background: rgba(15, 23, 42, 0.5);
    }

    .quickstart h2 {
      text-align: center;
      font-size: 2.5rem;
      margin-bottom: 3rem;
    }

    .code-block {
      max-width: 700px;
      margin: 0 auto;
      background: #0f172a;
      border-radius: 12px;
      overflow: hidden;
      border: 1px solid #334155;
    }

    .code-header {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 1rem;
      background: #1e293b;
      border-bottom: 1px solid #334155;
    }

    .dot {
      width: 12px;
      height: 12px;
      border-radius: 50%;
    }

    .dot.red { background: #ef4444; }
    .dot.yellow { background: #eab308; }
    .dot.green { background: #22c55e; }

    .code-title {
      margin-left: 0.5rem;
      color: #64748b;
      font-size: 0.875rem;
    }

    .code-block pre {
      padding: 1.5rem;
      margin: 0;
      overflow-x: auto;
    }

    .code-block code {
      font-family: 'SF Mono', Monaco, monospace;
      font-size: 0.9rem;
      line-height: 1.8;
    }

    .prompt { color: #22c55e; }
    .output { color: #60a5fa; }
    .success { color: #22c55e; }
    .info { color: #3b82f6; }

    .footer {
      padding: 4rem 2rem 2rem;
      border-top: 1px solid #334155;
    }

    .footer-content {
      max-width: 1200px;
      margin: 0 auto;
      display: flex;
      justify-content: space-between;
      align-items: center;
      flex-wrap: wrap;
      gap: 2rem;
    }

    .footer-brand p {
      color: #64748b;
      margin-top: 0.5rem;
    }

    .footer-brand a {
      color: #3b82f6;
      text-decoration: none;
    }

    .footer-links {
      display: flex;
      gap: 2rem;
    }

    .footer-links a {
      color: #94a3b8;
      text-decoration: none;
      transition: color 0.2s;
    }

    .footer-links a:hover {
      color: #f8fafc;
    }

    .footer-bottom {
      text-align: center;
      padding-top: 2rem;
      margin-top: 2rem;
      border-top: 1px solid #334155;
      color: #64748b;
    }

    @media (max-width: 768px) {
      .navbar {
        padding: 1rem 1.5rem;
        flex-direction: column;
        gap: 1rem;
      }

      .hero h1 {
        font-size: 2.5rem;
      }

      .hero-stats {
        flex-direction: column;
        gap: 1.5rem;
      }

      .cta-buttons {
        flex-direction: column;
        align-items: center;
      }

      .footer-content {
        flex-direction: column;
        text-align: center;
      }
    }
  `]
})
export class HomeComponent {}
