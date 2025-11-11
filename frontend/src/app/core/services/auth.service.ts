import { Injectable, signal, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable, tap } from 'rxjs';
import { jwtDecode } from 'jwt-decode';
import { environment } from '../../../environments/environment';
import { User, LoginRequest, RegisterRequest, AuthResponse } from '../models/user.model';

interface JwtPayload {
  user_id: number;
  email: string;
  role: string;
  exp: number;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private http = inject(HttpClient);
  private router = inject(Router);

  private readonly TOKEN_KEY = 'reservify_token';
  private readonly USER_KEY = 'reservify_user';

  // Estado reactivo con signals
  currentUser = signal<User | null>(this.getUserFromStorage());
  isAuthenticated = signal<boolean>(this.hasValidToken());
  isAdmin = signal<boolean>(this.checkIsAdmin());

  constructor() {
    this.checkTokenExpiration();
  }

  // LOGIN
  login(credentials: LoginRequest): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(`${environment.apiUrl}/auth/login`, credentials).pipe(
      tap(response => {
        if (response.success && response.data) {
          this.saveAuthData(response.data.token, response.data.user);
        }
      })
    );
  }

  // REGISTER
  register(data: RegisterRequest): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(`${environment.apiUrl}/auth/register`, data).pipe(
      tap(response => {
        if (response.success && response.data) {
          this.saveAuthData(response.data.token, response.data.user);
        }
      })
    );
  }

  // LOGOUT
  logout(): void {
    localStorage.removeItem(this.TOKEN_KEY);
    localStorage.removeItem(this.USER_KEY);
    this.currentUser.set(null);
    this.isAuthenticated.set(false);
    this.isAdmin.set(false);
    this.router.navigate(['/auth/login']);
  }

  // Obtener token
  getToken(): string | null {
    return localStorage.getItem(this.TOKEN_KEY);
  }

  // ========== MÉTODOS PRIVADOS ==========

  private saveAuthData(token: string, user: User): void {
    localStorage.setItem(this.TOKEN_KEY, token);
    localStorage.setItem(this.USER_KEY, JSON.stringify(user));
    this.currentUser.set(user);
    this.isAuthenticated.set(true);
    this.isAdmin.set(user.role === 'admin');
  }

  private getUserFromStorage(): User | null {
    const userJson = localStorage.getItem(this.USER_KEY);
    return userJson ? JSON.parse(userJson) : null;
  }

  private hasValidToken(): boolean {
    const token = this.getToken();
    if (!token) return false;

    try {
      const decoded = jwtDecode<JwtPayload>(token);
      const now = Date.now() / 1000;
      return decoded.exp > now;
    } catch {
      return false;
    }
  }

  private checkIsAdmin(): boolean {
    const user = this.getUserFromStorage();
    return user?.role === 'admin';
  }

  private checkTokenExpiration(): void {
    if (!this.hasValidToken()) {
      this.logout();
    }
  }
}