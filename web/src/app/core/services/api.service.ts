/*
 * GoAstra Frontend - API Service
 *
 * Base service for HTTP API calls with error handling.
 * Provides type-safe methods for CRUD operations.
 */
import { Injectable } from '@angular/core';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '@env/environment';

/*
 * Request options interface for HTTP calls.
 */
export interface RequestOptions {
  params?: Record<string, string>;
  headers?: Record<string, string>;
}

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private readonly baseUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  /*
   * Performs GET request to the API.
   */
  get<T>(path: string, params?: Record<string, string>): Observable<T> {
    let httpParams = new HttpParams();

    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          httpParams = httpParams.set(key, value);
        }
      });
    }

    return this.http.get<T>(`${this.baseUrl}${path}`, { params: httpParams });
  }

  /*
   * Performs POST request to the API.
   */
  post<T>(path: string, body: unknown, options?: RequestOptions): Observable<T> {
    const headers = this.buildHeaders(options?.headers);
    return this.http.post<T>(`${this.baseUrl}${path}`, body, { headers });
  }

  /*
   * Performs PUT request to the API.
   */
  put<T>(path: string, body: unknown, options?: RequestOptions): Observable<T> {
    const headers = this.buildHeaders(options?.headers);
    return this.http.put<T>(`${this.baseUrl}${path}`, body, { headers });
  }

  /*
   * Performs PATCH request to the API.
   */
  patch<T>(path: string, body: unknown, options?: RequestOptions): Observable<T> {
    const headers = this.buildHeaders(options?.headers);
    return this.http.patch<T>(`${this.baseUrl}${path}`, body, { headers });
  }

  /*
   * Performs DELETE request to the API.
   */
  delete<T>(path: string): Observable<T> {
    return this.http.delete<T>(`${this.baseUrl}${path}`);
  }

  /*
   * Uploads a file to the API.
   */
  upload<T>(path: string, file: File, fieldName = 'file'): Observable<T> {
    const formData = new FormData();
    formData.append(fieldName, file);

    return this.http.post<T>(`${this.baseUrl}${path}`, formData);
  }

  private buildHeaders(customHeaders?: Record<string, string>): HttpHeaders {
    let headers = new HttpHeaders();

    if (customHeaders) {
      Object.entries(customHeaders).forEach(([key, value]) => {
        headers = headers.set(key, value);
      });
    }

    return headers;
  }
}
