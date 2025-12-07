/*
 * GoAstra CLI - Module Generator
 *
 * Generates Angular feature modules with routing,
 * components, and services.
 */
package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

/*
 * ModuleGenerator handles Angular module code generation.
 */
type ModuleGenerator struct {
	name       string
	pascalName string
	camelName  string
}

/*
 * NewModuleGenerator creates a new module generator instance.
 */
func NewModuleGenerator(name string) *ModuleGenerator {
	return &ModuleGenerator{
		name:       name,
		pascalName: toPascalCase(name),
		camelName:  toCamelCase(name),
	}
}

/*
 * GenerateModule creates the Angular module file.
 */
func (g *ModuleGenerator) GenerateModule() error {
	dir := filepath.Join("web/src/app/features", g.name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	content := fmt.Sprintf(`/*
 * %s Module
 *
 * Feature module for %s functionality.
 * Configured for lazy loading.
 */
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';

import { %sRoutingModule } from './%s-routing.module';
import { %sComponent } from './%s.component';

@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    %sRoutingModule,
    %sComponent
  ]
})
export class %sModule { }
`,
		g.pascalName, g.name,
		g.pascalName, g.name,
		g.pascalName, g.name,
		g.pascalName,
		g.pascalName,
		g.pascalName,
	)

	path := filepath.Join(dir, g.name+".module.ts")
	return os.WriteFile(path, []byte(content), 0644)
}

/*
 * GenerateRouting creates the routing module file.
 */
func (g *ModuleGenerator) GenerateRouting() error {
	dir := filepath.Join("web/src/app/features", g.name)

	content := fmt.Sprintf(`/*
 * %s Routing Module
 *
 * Route configuration for %s feature.
 */
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { %sComponent } from './%s.component';

const routes: Routes = [
  {
    path: '',
    component: %sComponent,
    children: [
      {
        path: '',
        redirectTo: 'list',
        pathMatch: 'full'
      },
      {
        path: 'list',
        loadComponent: () => import('./components/%s-list/%s-list.component')
          .then(m => m.%sListComponent)
      },
      {
        path: ':id',
        loadComponent: () => import('./components/%s-detail/%s-detail.component')
          .then(m => m.%sDetailComponent)
      },
      {
        path: 'create',
        loadComponent: () => import('./components/%s-form/%s-form.component')
          .then(m => m.%sFormComponent)
      },
      {
        path: ':id/edit',
        loadComponent: () => import('./components/%s-form/%s-form.component')
          .then(m => m.%sFormComponent)
      }
    ]
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class %sRoutingModule { }
`,
		g.pascalName, g.name,
		g.pascalName, g.name,
		g.pascalName,
		g.name, g.name,
		g.pascalName,
		g.name, g.name,
		g.pascalName,
		g.name, g.name,
		g.pascalName,
		g.name, g.name,
		g.pascalName,
		g.pascalName,
	)

	path := filepath.Join(dir, g.name+"-routing.module.ts")
	return os.WriteFile(path, []byte(content), 0644)
}

/*
 * GenerateComponent creates the main feature component.
 */
func (g *ModuleGenerator) GenerateComponent() error {
	dir := filepath.Join("web/src/app/features", g.name)

	content := fmt.Sprintf(`/*
 * %s Component
 *
 * Container component for %s feature.
 */
import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-%s',
  standalone: true,
  imports: [RouterOutlet],
  template: ` + "`" + `
    <div class="%s-container">
      <router-outlet></router-outlet>
    </div>
  ` + "`" + `,
  styles: [` + "`" + `
    .%s-container {
      padding: 1rem;
    }
  ` + "`" + `]
})
export class %sComponent { }
`,
		g.pascalName, g.name,
		g.name,
		g.name,
		g.name,
		g.pascalName,
	)

	path := filepath.Join(dir, g.name+".component.ts")
	return os.WriteFile(path, []byte(content), 0644)
}

/*
 * GenerateService creates the feature service.
 */
func (g *ModuleGenerator) GenerateService() error {
	dir := filepath.Join("web/src/app/features", g.name)

	content := fmt.Sprintf(`/*
 * %s Service
 *
 * API service for %s operations.
 */
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '@core/services/api.service';

export interface %s {
  id: number;
  createdAt: string;
  updatedAt: string;
  /* TODO: Add model fields */
}

export interface %sPaginatedResponse {
  data: %s[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface Create%sDto {
  /* TODO: Add create fields */
}

export interface Update%sDto {
  /* TODO: Add update fields */
}

@Injectable({
  providedIn: 'root'
})
export class %sService {
  private readonly basePath = '/%s';

  constructor(private api: ApiService) {}

  list(page = 1, pageSize = 10): Observable<%sPaginatedResponse> {
    return this.api.get(this.basePath, {
      page: page.toString(),
      page_size: pageSize.toString()
    });
  }

  getById(id: number): Observable<%s> {
    return this.api.get(` + "`${this.basePath}/${id}`" + `);
  }

  create(data: Create%sDto): Observable<%s> {
    return this.api.post(this.basePath, data);
  }

  update(id: number, data: Update%sDto): Observable<%s> {
    return this.api.put(` + "`${this.basePath}/${id}`" + `, data);
  }

  delete(id: number): Observable<void> {
    return this.api.delete(` + "`${this.basePath}/${id}`" + `);
  }
}
`,
		g.pascalName, g.name,
		g.pascalName,
		g.pascalName, g.pascalName,
		g.pascalName,
		g.pascalName,
		g.pascalName,
		toPlural(g.name),
		g.pascalName,
		g.pascalName,
		g.pascalName, g.pascalName,
		g.pascalName, g.pascalName,
	)

	path := filepath.Join(dir, g.name+".service.ts")
	return os.WriteFile(path, []byte(content), 0644)
}
