/*
 * GoAstra CLI - CRUD Generator
 *
 * Full-stack CRUD generation combining backend API
 * and frontend module with all components.
 */
package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

/*
 * CRUDGenerator handles full-stack CRUD code generation.
 */
type CRUDGenerator struct {
	name       string
	pascalName string
	camelName  string
}

/*
 * NewCRUDGenerator creates a new CRUD generator instance.
 */
func NewCRUDGenerator(name string) *CRUDGenerator {
	return &CRUDGenerator{
		name:       name,
		pascalName: toPascalCase(name),
		camelName:  toCamelCase(name),
	}
}

/*
 * GenerateModel creates the Go model definition.
 */
func (g *CRUDGenerator) GenerateModel() error {
	content := fmt.Sprintf(`/*
 * %s Model
 *
 * Database entity and DTO definitions for %s.
 */
package models

import "time"

/*
 * %s represents the database entity.
 */
type %s struct {
	ID        uint      ` + "`db:\"id\" json:\"id\"`" + `
	CreatedAt time.Time ` + "`db:\"created_at\" json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`db:\"updated_at\" json:\"updated_at\"`" + `
	/* TODO: Add model fields */
}

/*
 * %sCreateDTO defines input for creating a %s.
 */
type %sCreateDTO struct {
	/* TODO: Add create fields */
}

/*
 * %sUpdateDTO defines input for updating a %s.
 */
type %sUpdateDTO struct {
	/* TODO: Add update fields */
}

/*
 * %sResponse defines the API response structure.
 */
type %sResponse struct {
	ID        uint      ` + "`json:\"id\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
	/* TODO: Add response fields */
}

/*
 * ToResponse converts a %s to its response representation.
 */
func (m *%s) ToResponse() *%sResponse {
	return &%sResponse{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
`,
		g.pascalName, g.name,
		g.pascalName,
		g.pascalName,
		g.pascalName, g.name,
		g.pascalName,
		g.pascalName, g.name,
		g.pascalName,
		g.pascalName,
		g.pascalName,
		g.pascalName,
		g.pascalName, g.pascalName,
		g.pascalName,
	)

	path := filepath.Join("app/internal/models", g.name+".go")
	return os.WriteFile(path, []byte(content), 0644)
}

/*
 * GenerateAPI creates backend handler, service, and repository.
 */
func (g *CRUDGenerator) GenerateAPI() error {
	apiGen := NewAPIGenerator(g.name)

	if err := apiGen.GenerateHandler(); err != nil {
		return err
	}

	if err := apiGen.GenerateService(); err != nil {
		return err
	}

	if err := apiGen.GenerateRepository(); err != nil {
		return err
	}

	return apiGen.GenerateRoutes()
}

/*
 * GenerateMigration creates a database migration file.
 */
func (g *CRUDGenerator) GenerateMigration() error {
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_create_%s_table.sql", timestamp, toSnakeCase(g.name))

	content := fmt.Sprintf(`-- Migration: Create %s table
-- Generated: %s

-- +migrate Up
CREATE TABLE IF NOT EXISTS %s (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    /* TODO: Add table columns */
);

CREATE INDEX idx_%s_created_at ON %s(created_at);

-- +migrate Down
DROP TABLE IF EXISTS %s;
`,
		toSnakeCase(g.name),
		time.Now().Format(time.RFC3339),
		toSnakeCase(g.name),
		toSnakeCase(g.name), toSnakeCase(g.name),
		toSnakeCase(g.name),
	)

	path := filepath.Join("app/migrations", filename)
	return os.WriteFile(path, []byte(content), 0644)
}

/*
 * GenerateModule creates the Angular feature module.
 */
func (g *CRUDGenerator) GenerateModule() error {
	modGen := NewModuleGenerator(g.name)

	if err := modGen.GenerateModule(); err != nil {
		return err
	}

	if err := modGen.GenerateRouting(); err != nil {
		return err
	}

	if err := modGen.GenerateComponent(); err != nil {
		return err
	}

	return modGen.GenerateService()
}

/*
 * GenerateComponents creates CRUD-specific Angular components.
 */
func (g *CRUDGenerator) GenerateComponents() error {
	if err := g.generateListComponent(); err != nil {
		return err
	}

	if err := g.generateDetailComponent(); err != nil {
		return err
	}

	return g.generateFormComponent()
}

func (g *CRUDGenerator) generateListComponent() error {
	dir := filepath.Join("web/src/app/features", g.name, "components", g.name+"-list")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	content := fmt.Sprintf(`/*
 * %s List Component
 *
 * Displays paginated list of %s resources.
 */
import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { %sService, %s, %sPaginatedResponse } from '../../%s.service';

@Component({
  selector: 'app-%s-list',
  standalone: true,
  imports: [CommonModule, RouterLink],
  template: ` + "`" + `
    <div class="%s-list">
      <header class="list-header">
        <h1>%s</h1>
        <a routerLink="../create" class="btn btn-primary">Create New</a>
      </header>

      @if (loading()) {
        <div class="loading">Loading...</div>
      }

      @if (error()) {
        <div class="error">{{ error() }}</div>
      }

      <table class="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Created</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          @for (item of items(); track item.id) {
            <tr>
              <td>{{ item.id }}</td>
              <td>{{ item.createdAt | date:'short' }}</td>
              <td>
                <a [routerLink]="['..', item.id]">View</a>
                <a [routerLink]="['..', item.id, 'edit']">Edit</a>
                <button (click)="delete(item.id)">Delete</button>
              </td>
            </tr>
          }
        </tbody>
      </table>

      <div class="pagination">
        <button (click)="prevPage()" [disabled]="page() === 1">Previous</button>
        <span>Page {{ page() }} of {{ totalPages() }}</span>
        <button (click)="nextPage()" [disabled]="page() >= totalPages()">Next</button>
      </div>
    </div>
  ` + "`" + `,
  styles: [` + "`" + `
    .list-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 1rem;
    }
    .data-table {
      width: 100%%;
      border-collapse: collapse;
    }
    .data-table th,
    .data-table td {
      padding: 0.75rem;
      text-align: left;
      border-bottom: 1px solid #ddd;
    }
    .pagination {
      display: flex;
      justify-content: center;
      gap: 1rem;
      margin-top: 1rem;
    }
    .btn-primary {
      background: #007bff;
      color: white;
      padding: 0.5rem 1rem;
      text-decoration: none;
      border-radius: 4px;
    }
  ` + "`" + `]
})
export class %sListComponent implements OnInit {
  private service = inject(%sService);

  items = signal<%s[]>([]);
  loading = signal(false);
  error = signal<string | null>(null);
  page = signal(1);
  totalPages = signal(1);
  pageSize = 10;

  ngOnInit(): void {
    this.loadData();
  }

  loadData(): void {
    this.loading.set(true);
    this.error.set(null);

    this.service.list(this.page(), this.pageSize).subscribe({
      next: (response) => {
        this.items.set(response.data);
        this.totalPages.set(response.totalPages);
        this.loading.set(false);
      },
      error: (err) => {
        this.error.set(err.message);
        this.loading.set(false);
      }
    });
  }

  prevPage(): void {
    if (this.page() > 1) {
      this.page.update(p => p - 1);
      this.loadData();
    }
  }

  nextPage(): void {
    if (this.page() < this.totalPages()) {
      this.page.update(p => p + 1);
      this.loadData();
    }
  }

  delete(id: number): void {
    if (confirm('Are you sure you want to delete this item?')) {
      this.service.delete(id).subscribe({
        next: () => this.loadData(),
        error: (err) => this.error.set(err.message)
      });
    }
  }
}
`,
		g.pascalName, g.name,
		g.pascalName, g.pascalName, g.pascalName, g.name,
		g.name,
		g.name,
		g.pascalName,
		g.pascalName,
		g.pascalName,
		g.pascalName,
	)

	path := filepath.Join(dir, g.name+"-list.component.ts")
	return os.WriteFile(path, []byte(content), 0644)
}

func (g *CRUDGenerator) generateDetailComponent() error {
	dir := filepath.Join("web/src/app/features", g.name, "components", g.name+"-detail")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	content := fmt.Sprintf(`/*
 * %s Detail Component
 *
 * Displays single %s resource details.
 */
import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { %sService, %s } from '../../%s.service';

@Component({
  selector: 'app-%s-detail',
  standalone: true,
  imports: [CommonModule, RouterLink],
  template: ` + "`" + `
    <div class="%s-detail">
      @if (loading()) {
        <div class="loading">Loading...</div>
      }

      @if (error()) {
        <div class="error">{{ error() }}</div>
      }

      @if (item()) {
        <header>
          <h1>%s #{{ item()!.id }}</h1>
          <div class="actions">
            <a [routerLink]="['..', item()!.id, 'edit']" class="btn">Edit</a>
            <a routerLink=".." class="btn">Back to List</a>
          </div>
        </header>

        <div class="detail-content">
          <dl>
            <dt>ID</dt>
            <dd>{{ item()!.id }}</dd>

            <dt>Created At</dt>
            <dd>{{ item()!.createdAt | date:'medium' }}</dd>

            <dt>Updated At</dt>
            <dd>{{ item()!.updatedAt | date:'medium' }}</dd>
          </dl>
        </div>
      }
    </div>
  ` + "`" + `,
  styles: [` + "`" + `
    header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 1rem;
    }
    .actions {
      display: flex;
      gap: 0.5rem;
    }
    .btn {
      padding: 0.5rem 1rem;
      text-decoration: none;
      border: 1px solid #ddd;
      border-radius: 4px;
    }
    dl {
      display: grid;
      grid-template-columns: auto 1fr;
      gap: 0.5rem 1rem;
    }
    dt {
      font-weight: bold;
    }
  ` + "`" + `]
})
export class %sDetailComponent implements OnInit {
  private service = inject(%sService);
  private route = inject(ActivatedRoute);

  item = signal<%s | null>(null);
  loading = signal(false);
  error = signal<string | null>(null);

  ngOnInit(): void {
    const id = Number(this.route.snapshot.paramMap.get('id'));
    this.loadData(id);
  }

  loadData(id: number): void {
    this.loading.set(true);
    this.error.set(null);

    this.service.getById(id).subscribe({
      next: (data) => {
        this.item.set(data);
        this.loading.set(false);
      },
      error: (err) => {
        this.error.set(err.message);
        this.loading.set(false);
      }
    });
  }
}
`,
		g.pascalName, g.name,
		g.pascalName, g.pascalName, g.name,
		g.name,
		g.name,
		g.pascalName,
		g.pascalName,
		g.pascalName,
		g.pascalName,
	)

	path := filepath.Join(dir, g.name+"-detail.component.ts")
	return os.WriteFile(path, []byte(content), 0644)
}

func (g *CRUDGenerator) generateFormComponent() error {
	dir := filepath.Join("web/src/app/features", g.name, "components", g.name+"-form")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	content := fmt.Sprintf(`/*
 * %s Form Component
 *
 * Create and edit form for %s resources.
 */
import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { %sService, %s, Create%sDto, Update%sDto } from '../../%s.service';

@Component({
  selector: 'app-%s-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterLink],
  template: ` + "`" + `
    <div class="%s-form">
      <header>
        <h1>{{ isEdit() ? 'Edit' : 'Create' }} %s</h1>
        <a routerLink=".." class="btn">Cancel</a>
      </header>

      @if (error()) {
        <div class="error">{{ error() }}</div>
      }

      <form [formGroup]="form" (ngSubmit)="submit()">
        <!-- TODO: Add form fields -->
        <div class="form-group">
          <label for="example">Example Field</label>
          <input type="text" id="example" formControlName="example">
        </div>

        <div class="form-actions">
          <button type="submit" [disabled]="form.invalid || submitting()">
            {{ submitting() ? 'Saving...' : (isEdit() ? 'Update' : 'Create') }}
          </button>
        </div>
      </form>
    </div>
  ` + "`" + `,
  styles: [` + "`" + `
    header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 1rem;
    }
    .form-group {
      margin-bottom: 1rem;
    }
    .form-group label {
      display: block;
      margin-bottom: 0.25rem;
      font-weight: 500;
    }
    .form-group input {
      width: 100%%;
      padding: 0.5rem;
      border: 1px solid #ddd;
      border-radius: 4px;
    }
    .form-actions {
      margin-top: 1.5rem;
    }
    .btn {
      padding: 0.5rem 1rem;
      text-decoration: none;
      border: 1px solid #ddd;
      border-radius: 4px;
    }
    button[type="submit"] {
      background: #007bff;
      color: white;
      padding: 0.5rem 1.5rem;
      border: none;
      border-radius: 4px;
      cursor: pointer;
    }
    button:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }
    .error {
      color: red;
      margin-bottom: 1rem;
    }
  ` + "`" + `]
})
export class %sFormComponent implements OnInit {
  private fb = inject(FormBuilder);
  private service = inject(%sService);
  private route = inject(ActivatedRoute);
  private router = inject(Router);

  form: FormGroup = this.fb.group({
    example: ['', Validators.required]
    /* TODO: Add form controls */
  });

  isEdit = signal(false);
  submitting = signal(false);
  error = signal<string | null>(null);
  private itemId: number | null = null;

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.isEdit.set(true);
      this.itemId = Number(id);
      this.loadData(this.itemId);
    }
  }

  loadData(id: number): void {
    this.service.getById(id).subscribe({
      next: (data) => {
        this.form.patchValue(data);
      },
      error: (err) => {
        this.error.set(err.message);
      }
    });
  }

  submit(): void {
    if (this.form.invalid) return;

    this.submitting.set(true);
    this.error.set(null);

    const data = this.form.value;

    const request = this.isEdit()
      ? this.service.update(this.itemId!, data as Update%sDto)
      : this.service.create(data as Create%sDto);

    request.subscribe({
      next: (result) => {
        this.router.navigate(['..', result.id], { relativeTo: this.route });
      },
      error: (err) => {
        this.error.set(err.message);
        this.submitting.set(false);
      }
    });
  }
}
`,
		g.pascalName, g.name,
		g.pascalName, g.pascalName, g.pascalName, g.pascalName, g.name,
		g.name,
		g.name,
		g.pascalName,
		g.pascalName,
		g.pascalName,
		g.pascalName,
		g.pascalName,
	)

	path := filepath.Join(dir, g.name+"-form.component.ts")
	return os.WriteFile(path, []byte(content), 0644)
}

/*
 * UpdateRoutes modifies the main app routing to include the new module.
 */
func (g *CRUDGenerator) UpdateRoutes() error {
	/* In a real implementation, this would parse and modify app.routes.ts */
	/* For now, we output instructions */
	fmt.Printf("\nAdd the following route to web/src/app/app.routes.ts:\n")
	fmt.Printf(`
  {
    path: '%s',
    loadChildren: () => import('@features/%s/%s.module').then(m => m.%sModule),
    canActivate: [authGuard]
  }
`, toPlural(g.name), g.name, g.name, g.pascalName)

	return nil
}
