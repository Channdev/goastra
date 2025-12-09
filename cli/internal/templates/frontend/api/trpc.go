/*
 * GoAstra CLI - Frontend tRPC Templates
 *
 * Generates Angular templates for tRPC with Connect-Web.
 * Includes package.json, Connect setup, and service templates.
 */
package api

import "fmt"

// TRPCPackageJSON returns package.json with Connect dependencies.
func TRPCPackageJSON(projectName string) string {
	return fmt.Sprintf(`{
  "name": "%s-web",
  "version": "1.0.0",
  "scripts": {
    "ng": "ng",
    "start": "ng serve --proxy-config proxy.conf.json",
    "build": "ng build",
    "watch": "ng build --watch --configuration development",
    "test": "ng test",
    "generate": "buf generate"
  },
  "private": true,
  "dependencies": {
    "@angular/animations": "^17.0.0",
    "@angular/common": "^17.0.0",
    "@angular/compiler": "^17.0.0",
    "@angular/core": "^17.0.0",
    "@angular/forms": "^17.0.0",
    "@angular/platform-browser": "^17.0.0",
    "@angular/platform-browser-dynamic": "^17.0.0",
    "@angular/router": "^17.0.0",
    "@connectrpc/connect": "^1.1.0",
    "@connectrpc/connect-web": "^1.1.0",
    "@bufbuild/protobuf": "^1.5.0",
    "rxjs": "~7.8.0",
    "tslib": "^2.6.0",
    "zone.js": "~0.14.0"
  },
  "devDependencies": {
    "@angular-devkit/build-angular": "^17.0.0",
    "@angular/cli": "^17.0.0",
    "@angular/compiler-cli": "^17.0.0",
    "@bufbuild/buf": "^1.28.0",
    "@bufbuild/protoc-gen-es": "^1.5.0",
    "@connectrpc/protoc-gen-connect-es": "^1.1.0",
    "@types/jasmine": "~5.1.0",
    "jasmine-core": "~5.1.0",
    "karma": "~6.4.0",
    "karma-chrome-launcher": "~3.2.0",
    "karma-coverage": "~2.2.0",
    "karma-jasmine": "~5.1.0",
    "karma-jasmine-html-reporter": "~2.1.0",
    "typescript": "~5.2.0"
  }
}`, projectName)
}

// TRPCAppConfigTS returns the app.config.ts for tRPC setup.
func TRPCAppConfigTS() string {
	return `import { ApplicationConfig } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient } from '@angular/common/http';
import { routes } from './app.routes';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideHttpClient(),
  ]
};
`
}

// BufGenYAMLWeb returns the buf.gen.yaml for frontend code generation.
func BufGenYAMLWeb() string {
	return `version: v1
managed:
  enabled: true
plugins:
  - plugin: buf.build/bufbuild/es
    out: src/app/core/rpc/gen
    opt: target=ts
  - plugin: buf.build/connectrpc/es
    out: src/app/core/rpc/gen
    opt: target=ts
`
}

// TRPCServiceTS returns the tRPC service template.
func TRPCServiceTS() string {
	return `/*
 * tRPC Service
 *
 * Angular service for tRPC operations with Connect-Web.
 * Provides type-safe RPC client access.
 */
import { Injectable } from '@angular/core';
import { createPromiseClient, Transport } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { from, Observable } from 'rxjs';
import { environment } from '@env/environment';

// Import generated service clients
import { HealthService } from './gen/proto/v1/service_connect';
import { AuthService } from './gen/proto/v1/service_connect';
import { UserService } from './gen/proto/v1/service_connect';

@Injectable({
  providedIn: 'root'
})
export class TRPCService {
  private transport: Transport;

  public health: ReturnType<typeof createPromiseClient<typeof HealthService>>;
  public auth: ReturnType<typeof createPromiseClient<typeof AuthService>>;
  public users: ReturnType<typeof createPromiseClient<typeof UserService>>;

  constructor() {
    this.transport = createConnectTransport({
      baseUrl: environment.apiUrl,
      interceptors: [this.authInterceptor()],
    });

    // Initialize service clients
    this.health = createPromiseClient(HealthService, this.transport);
    this.auth = createPromiseClient(AuthService, this.transport);
    this.users = createPromiseClient(UserService, this.transport);
  }

  /**
   * Convert promise to observable for Angular compatibility
   */
  toObservable<T>(promise: Promise<T>): Observable<T> {
    return from(promise);
  }

  /**
   * Auth interceptor to add JWT token to requests
   */
  private authInterceptor() {
    return (next: any) => async (req: any) => {
      const token = localStorage.getItem('access_token');
      if (token) {
        req.header.set('Authorization', ` + "`Bearer ${token}`" + `);
      }
      return next(req);
    };
  }
}
`
}

// TRPCEnvTS returns the environment.ts for tRPC.
func TRPCEnvTS() string {
	return `export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080',
};
`
}

// TRPCEnvProdTS returns the environment.prod.ts for tRPC.
func TRPCEnvProdTS() string {
	return `export const environment = {
  production: true,
  apiUrl: '',
};
`
}
