package frontend

func MinimalStylesCSS() string {
	return `* { box-sizing: border-box; margin: 0; padding: 0; }

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: #f8fafc;
  color: #1e293b;
}
`
}

func MinimalAppRoutes() string {
	return `import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';

export const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: '**', redirectTo: '' }
];
`
}

func MinimalHomeComponent() string {
	return `import { Component } from '@angular/core';

@Component({
  selector: 'app-home',
  standalone: true,
  template: ` + "`" + `
    <main>
      <h1>Hello, GoAstra!</h1>
      <p>Start building your app.</p>
    </main>
  ` + "`" + `,
  styles: [` + "`" + `
    main { padding: 2rem; text-align: center; }
    h1 { margin-bottom: 0.5rem; }
  ` + "`" + `]
})
export class HomeComponent {}
`
}
