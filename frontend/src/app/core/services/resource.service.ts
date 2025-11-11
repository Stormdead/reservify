import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Resource } from '../models/resource.model';

interface ApiResponse<T> {
  success: boolean;
  message: string;
  data?: T;
}

interface PaginatedResponse<T> {
  success: boolean;
  message: string;
  data: T[];
  meta: {
    page: number;
    page_size: number;
    total_items: number;
    total_pages: number;
  };
}

@Injectable({
  providedIn: 'root'
})
export class ResourceService {
  private http = inject(HttpClient);
  private apiUrl = `${environment.apiUrl}/resources`;

  // Obtener todos los recursos (con paginación)
  getResources(page: number = 1, pageSize: number = 10, search: string = ''): Observable<PaginatedResponse<Resource>> {
    let params = new HttpParams()
      .set('page', page.toString())
      .set('page_size', pageSize.toString());

    if (search) {
      params = params.set('search', search);
    }

    return this.http.get<PaginatedResponse<Resource>>(this.apiUrl, { params });
  }

  // Obtener recursos por categoría
  getResourcesByCategory(category: string, page: number = 1, pageSize: number = 10): Observable<PaginatedResponse<Resource>> {
    const params = new HttpParams()
      .set('page', page.toString())
      .set('page_size', pageSize.toString());

    return this.http.get<PaginatedResponse<Resource>>(`${this.apiUrl}/category/${category}`, { params });
  }

  // Obtener un recurso por ID
  getResourceById(id: number): Observable<ApiResponse<Resource>> {
    return this.http.get<ApiResponse<Resource>>(`${this.apiUrl}/${id}`);
  }

  // Obtener categorías
  getCategories(): Observable<ApiResponse<{ categories: string[] }>> {
    return this.http.get<ApiResponse<{ categories: string[] }>>(`${this.apiUrl}/categories`);
  }

  // ADMIN: Crear recurso
  createResource(resource: any): Observable<ApiResponse<Resource>> {
    return this.http.post<ApiResponse<Resource>>(`${environment.apiUrl}/admin/resources`, resource);
  }

  // ADMIN: Actualizar recurso
  updateResource(id: number, resource: any): Observable<ApiResponse<Resource>> {
    return this.http.put<ApiResponse<Resource>>(`${environment.apiUrl}/admin/resources/${id}`, resource);
  }

  // ADMIN: Eliminar recurso
  deleteResource(id: number): Observable<ApiResponse<void>> {
    return this.http.delete<ApiResponse<void>>(`${environment.apiUrl}/admin/resources/${id}`);
  }
}