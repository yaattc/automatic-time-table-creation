import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { LoginRequestModel } from '../model/login-request-model';
import { CookieService } from 'ngx-cookie-service';
import {
  ActivatedRouteSnapshot,
  CanActivate,
  Router,
  RouterStateSnapshot,
  UrlTree,
} from '@angular/router';
import { PANEL_DASHBOARD, PANEL_LOGIN } from '../constants/routes';
import { TOKEN_COOKIE_NAME } from '../constants/cookie';
import { JwtPayloadModel } from '../model/jwt-payload-model';
import jwt_decode from 'jwt-decode';
import { BehaviorSubject, Observable, of } from 'rxjs';
import { ErrorResponseModel } from '../model/error-response-model';
import { catchError, map } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { JwtResponseModel } from '../model/jwt-response-model';

@Injectable({
  providedIn: 'root',
})
export class AuthService implements CanActivate {
  private isAuthenticatedSource$ = new BehaviorSubject<boolean>(false);
  private errorSource$ = new BehaviorSubject<ErrorResponseModel | null>(null);
  private isLoadingSource$ = new BehaviorSubject<boolean>(false);
  private initialized = false;

  public readonly isAuthenticated$ = this.isAuthenticatedSource$.asObservable();
  public readonly error$ = this.errorSource$.asObservable();
  public readonly isLoading$ = this.isLoadingSource$.asObservable();
  public tokenData: JwtPayloadModel;

  constructor(
    private http: HttpClient,
    private cookieService: CookieService,
    private router: Router,
  ) {
    if (this.cookieService.check(TOKEN_COOKIE_NAME)) {
      this.tokenData = jwt_decode(this.cookieService.get(TOKEN_COOKIE_NAME));
    }
  }

  public canActivate(
    next: ActivatedRouteSnapshot,
    state: RouterStateSnapshot,
  ): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
    return this.initialized;
    // return (this.initialized ? this.isAuthenticated$ : this.checkToken()).pipe(
    //   map((isAuthorized) => {
    //     if (!isAuthorized) {
    //       return this.router.parseUrl(PANEL_LOGIN);
    //     }
    //     return true;
    //   }),
    // );
  }

  public login(user: string, passwd: string): void {
    // this.isLoadingSource$.next(true);
    this.http
      .post<JwtResponseModel>(`${environment.apiUrl}/auth/local/login`, {
        user,
        passwd,
      } as LoginRequestModel)
      .subscribe(
        (response) => {
          this.router.navigateByUrl(PANEL_DASHBOARD);
          // this.isAuthenticatedSource$.next(true);
          // this.cookieService.delete(TOKEN_COOKIE_NAME);
          // this.cookieService.set(TOKEN_COOKIE_NAME, response.accessToken);
          // this.tokenData = jwt_decode(response.accessToken);
          // this.isLoadingSource$.next(false);
          // this.errorSource$.next(null);
        },
        (error: HttpErrorResponse) => {
          // this.errorSource$.next(error.error);
          // this.isAuthenticatedSource$.next(false);
          // this.isLoadingSource$.next(false);
        },
      );
  }

  public logout(): void {
    // this.cookieService.delete(TOKEN_COOKIE_NAME);
    // this.isAuthenticatedSource$.next(false);
    // this.tokenData = undefined;
    this.initialized = false;
    this.router.navigateByUrl(PANEL_LOGIN);
  }

  private checkToken(): Observable<boolean> {
    if (this.cookieService.check(TOKEN_COOKIE_NAME)) {
      return this.http
        .get(`${environment.apiUrl}/auth/check`, {
          headers: {
            Authorization: `Bearer ${this.cookieService.get(TOKEN_COOKIE_NAME)}`,
          },
        })
        .pipe(
          map(() => {
            this.isAuthenticatedSource$.next(true);
            return true;
          }),
          catchError(() => {
            this.isAuthenticatedSource$.next(false);
            this.cookieService.delete(TOKEN_COOKIE_NAME);
            return of(false);
          }),
        );
    } else {
      return of(false);
    }
  }

  getInitializing(): boolean {
    return this.initialized;
  }
}
