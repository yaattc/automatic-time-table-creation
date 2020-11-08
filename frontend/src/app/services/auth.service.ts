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
  private initialized = false;

  private currentPageSource$ = new BehaviorSubject<string>(window.location.href);
  currentPage$ = this.currentPageSource$.asObservable();

  public readonly isAuthenticated$ = this.isAuthenticatedSource$.asObservable();
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
    return (this.initialized ? this.isAuthenticated$ : this.checkToken()).pipe(
      map((isAuthorized) => {
        if (!isAuthorized) {
          return this.router.parseUrl(PANEL_LOGIN);
        }
        return true;
      }),
    );
  }

  public login(user: string, passwd: string): void {
    this.http
      .post<JwtResponseModel>(`${environment.apiUrl}/auth/local/login`, {
        user,
        passwd,
      } as LoginRequestModel)
      .subscribe(
        (response) => {
          this.isAuthenticatedSource$.next(true);
          this.cookieService.delete(TOKEN_COOKIE_NAME);
          this.cookieService.set(TOKEN_COOKIE_NAME, response.token);
          this.errorSource$.next(null);
          this.router.navigateByUrl(PANEL_DASHBOARD);
          this.currentPageSource$.next(`${environment.apiUrl}/${PANEL_DASHBOARD}`);
        },
        (error: HttpErrorResponse) => {
          this.errorSource$.next(error.error);
          this.isAuthenticatedSource$.next(false);
        },
      );
  }

  public logout(): void {
    this.cookieService.delete(TOKEN_COOKIE_NAME);
    this.isAuthenticatedSource$.next(false);
    this.tokenData = undefined;
    this.router.navigateByUrl(PANEL_LOGIN);
    this.currentPageSource$.next(`${environment.apiUrl}/${PANEL_LOGIN}`);
  }

  private checkToken(): Observable<boolean> {
    if (this.cookieService.check(TOKEN_COOKIE_NAME)) {
      this.isAuthenticatedSource$.next(true);
      return of(true);
    } else {
      return of(false);
    }
  }

  public setCurrentPage(path: string): void {
    this.currentPageSource$.next(path);
  }
}
