import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatChipsModule } from '@angular/material/chips';
import { MatPaginatorModule, PageEvent } from '@angular/material/paginator';
import { ResourceService } from '../../../core/services/resource.service';
import { Resource } from '../../../core/models/resource.model';

@Component({
  selector: 'app-resource-list',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatButtonModule,
    MatProgressSpinnerModule,
    MatChipsModule,
    MatPaginatorModule
  ],
  templateUrl: './resource-list.component.html',
  styleUrl: './resource-list.component.css'
})
export class ResourceListComponent implements OnInit {
  private resourceService = inject(ResourceService);
  private router = inject(Router);

  resources: Resource[] = [];
  loading = true;
  
  // Paginación
  totalItems = 0;
  pageSize = 9;
  currentPage = 0;

  ngOnInit(): void {
    this.loadResources();
  }

  loadResources(): void {
    this.loading = true;
    const page = this.currentPage + 1; // Backend usa página 1-indexed

    this.resourceService.getResources(page, this.pageSize).subscribe({
      next: (response) => {
        this.resources = response.data;
        this.totalItems = response.meta.total_items;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error al cargar recursos:', error);
        this.loading = false;
      }
    });
  }

  onPageChange(event: PageEvent): void {
    this.currentPage = event.pageIndex;
    this.pageSize = event.pageSize;
    this.loadResources();
  }

  viewResource(id: number): void {
    this.router.navigate(['/resources', id]);
  }

  createBooking(resourceId: number): void {
    this.router.navigate(['/bookings/create'], { queryParams: { resourceId } });
  }
}