/*
 * GoAstra CLI - Frontend GraphQL Templates
 *
 * Generates Angular templates for GraphQL with Apollo Client.
 * Includes package.json, Apollo setup, and service templates.
 */
package api

import "fmt"

// GraphQLPackageJSON returns package.json with Apollo dependencies.
func GraphQLPackageJSON(projectName string) string {
	return fmt.Sprintf(`{
  "name": "%s-web",
  "version": "1.0.0",
  "scripts": {
    "ng": "ng",
    "start": "ng serve --proxy-config proxy.conf.json",
    "build": "ng build",
    "watch": "ng build --watch --configuration development",
    "test": "ng test",
    "codegen": "graphql-codegen --config codegen.yml"
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
    "@apollo/client": "^3.8.0",
    "apollo-angular": "^6.0.0",
    "graphql": "^16.8.0",
    "rxjs": "~7.8.0",
    "tslib": "^2.6.0",
    "zone.js": "~0.14.0"
  },
  "devDependencies": {
    "@angular-devkit/build-angular": "^17.0.0",
    "@angular/cli": "^17.0.0",
    "@angular/compiler-cli": "^17.0.0",
    "@graphql-codegen/cli": "^5.0.0",
    "@graphql-codegen/typescript": "^4.0.0",
    "@graphql-codegen/typescript-operations": "^4.0.0",
    "@graphql-codegen/typescript-apollo-angular": "^4.0.0",
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

// ApolloConfigTS returns the app.config.ts with Apollo setup.
func ApolloConfigTS() string {
	return `import { ApplicationConfig, inject } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient } from '@angular/common/http';
import { provideApollo } from 'apollo-angular';
import { HttpLink } from 'apollo-angular/http';
import { InMemoryCache, ApolloLink } from '@apollo/client/core';
import { setContext } from '@apollo/client/link/context';
import { routes } from './app.routes';
import { environment } from '@env/environment';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideHttpClient(),
    provideApollo(() => {
      const httpLink = inject(HttpLink);

      // Create auth link for JWT token
      const authLink = setContext((_, { headers }) => {
        const token = localStorage.getItem('access_token');
        return {
          headers: {
            ...headers,
            authorization: token ? ` + "`Bearer ${token}`" + ` : '',
          }
        };
      });

      return {
        link: ApolloLink.from([authLink, httpLink.create({ uri: environment.graphqlUrl })]),
        cache: new InMemoryCache(),
        defaultOptions: {
          watchQuery: {
            fetchPolicy: 'cache-and-network',
          },
        },
      };
    }),
  ]
};
`
}

// CodegenYML returns the GraphQL codegen configuration.
func CodegenYML() string {
	return `overwrite: true
schema: "http://localhost:8080/graphql"
documents: "src/**/*.graphql"
generates:
  src/app/core/graphql/generated.ts:
    plugins:
      - "typescript"
      - "typescript-operations"
      - "typescript-apollo-angular"
    config:
      addExplicitOverride: true
      withHooks: false
      apolloAngularVersion: 6
`
}

// GraphQLServiceTS returns the GraphQL service template.
func GraphQLServiceTS() string {
	return `/*
 * GraphQL Service
 *
 * Angular service for GraphQL operations with Apollo Client.
 * Provides typed query and mutation helpers.
 */
import { Injectable, inject } from '@angular/core';
import { Apollo, gql } from 'apollo-angular';
import { Observable, map, catchError, throwError } from 'rxjs';
import { ApolloQueryResult } from '@apollo/client/core';

@Injectable({
  providedIn: 'root'
})
export class GraphQLService {
  private apollo = inject(Apollo);

  /**
   * Execute a GraphQL query
   */
  query<T>(query: string, variables?: Record<string, any>): Observable<T> {
    return this.apollo.query<T>({
      query: gql(query),
      variables,
    }).pipe(
      map((result: ApolloQueryResult<T>) => result.data),
      catchError(this.handleError)
    );
  }

  /**
   * Execute a GraphQL mutation
   */
  mutate<T>(mutation: string, variables?: Record<string, any>): Observable<T> {
    return this.apollo.mutate<T>({
      mutation: gql(mutation),
      variables,
    }).pipe(
      map(result => result.data as T),
      catchError(this.handleError)
    );
  }

  /**
   * Subscribe to a GraphQL query (refetches on cache update)
   */
  watch<T>(query: string, variables?: Record<string, any>): Observable<T> {
    return this.apollo.watchQuery<T>({
      query: gql(query),
      variables,
    }).valueChanges.pipe(
      map(result => result.data),
      catchError(this.handleError)
    );
  }

  private handleError(error: any): Observable<never> {
    console.error('GraphQL error:', error);
    return throwError(() => error);
  }
}
`
}

// GraphQLEnvTS returns the environment.ts with GraphQL URL.
func GraphQLEnvTS() string {
	return `export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080',
  graphqlUrl: 'http://localhost:8080/graphql',
};
`
}

// GraphQLEnvProdTS returns the environment.prod.ts with GraphQL URL.
func GraphQLEnvProdTS() string {
	return `export const environment = {
  production: true,
  apiUrl: '/api',
  graphqlUrl: '/graphql',
};
`
}
