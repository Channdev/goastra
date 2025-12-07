/*
 * GoAstra Frontend - Application Bootstrap
 *
 * Entry point for the Angular application.
 * Configures providers and bootstraps the root component.
 */
import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { appConfig } from './app/app.config';

bootstrapApplication(AppComponent, appConfig)
  .catch((err) => console.error(err));
