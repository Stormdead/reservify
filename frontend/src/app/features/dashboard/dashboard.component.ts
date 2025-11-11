import { Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterLink } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatToolbarModule } from '@angular/material/toolbar';
import { AuthService } from '../../core/services/auth.service';
import { BookingService } from '../../core/services/booking.service';
import { BookingListItem } from '../../core/models/booking.model';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    RouterLink,
    MatButtonModule,
    MatCardModule,
    MatIconModule,
    MatToolbarModule
  ],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.css'
})
export class DashboardComponent implements OnInit {
  authService = inject(AuthService);
  private bookingService = inject(BookingService);
  private router = inject(Router);

  upcomingBookings: BookingListItem[] = [];
  loading = true;

  ngOnInit(): void {
    this.loadUpcomingBookings();
  }

  loadUpcomingBookings(): void {
    this.bookingService.getUpcomingBookings().subscribe({
      next: (response) => {
        this.upcomingBookings = response.data || [];
        this.loading = false;
      },
      error: (error) => {
        console.error('Error al cargar próximas reservas:', error);
        this.loading = false;
      }
    });
  }

  logout(): void {
    this.authService.logout();
  }

  navigateTo(path: string): void {
    this.router.navigate([path]);
  }
}