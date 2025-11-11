import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Booking, CreateBookingRequest, BookingListItem } from '../models/booking.model';

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
export class BookingService {
  private http = inject(HttpClient);
  private apiUrl = `${environment.apiUrl}/bookings`;

  // Obtener mis reservas
  getMyBookings(page: number = 1, pageSize: number = 10): Observable<PaginatedResponse<BookingListItem>> {
    const params = new HttpParams()
      .set('page', page.toString())
      .set('page_size', pageSize.toString());

    return this.http.get<PaginatedResponse<BookingListItem>>(`${this.apiUrl}/my`, { params });
  }

  // Obtener próximas reservas
  getUpcomingBookings(): Observable<ApiResponse<BookingListItem[]>> {
    return this.http.get<ApiResponse<BookingListItem[]>>(`${this.apiUrl}/upcoming`);
  }

  // Obtener una reserva por ID
  getBookingById(id: number): Observable<ApiResponse<Booking>> {
    return this.http.get<ApiResponse<Booking>>(`${this.apiUrl}/${id}`);
  }

  // Crear reserva
  createBooking(booking: CreateBookingRequest): Observable<ApiResponse<Booking>> {
    return this.http.post<ApiResponse<Booking>>(this.apiUrl, booking);
  }

  // Actualizar reserva
  updateBooking(id: number, booking: any): Observable<ApiResponse<Booking>> {
    return this.http.put<ApiResponse<Booking>>(`${this.apiUrl}/${id}`, booking);
  }

  // Cancelar reserva
  cancelBooking(id: number): Observable<ApiResponse<void>> {
    return this.http.delete<ApiResponse<void>>(`${this.apiUrl}/${id}`);
  }

  // ADMIN: Obtener todas las reservas
  getAllBookings(page: number = 1, pageSize: number = 10): Observable<PaginatedResponse<BookingListItem>> {
    const params = new HttpParams()
      .set('page', page.toString())
      .set('page_size', pageSize.toString());

    return this.http.get<PaginatedResponse<BookingListItem>>(`${environment.apiUrl}/admin/bookings`, { params });
  }

  // ADMIN: Cambiar estado de reserva
  changeBookingStatus(id: number, status: string): Observable<ApiResponse<Booking>> {
    return this.http.patch<ApiResponse<Booking>>(`${environment.apiUrl}/admin/bookings/${id}/status`, { status });
  }
}