import {
  HttpErrorResponse,
  HttpEvent,
  HttpHandler,
  HttpInterceptor,
  HttpRequest,
} from '@angular/common/http';
import { Injectable } from '@angular/core';
import { CookieService } from 'ngx-cookie-service';
import { Router } from '@angular/router';
import { Observable, throwError } from 'rxjs';
import { TOKEN_COOKIE_NAME } from '../constants/cookie';
import { environment } from '../../environments/environment';
import { catchError } from 'rxjs/operators';
import { PANEL_LOGIN } from '../constants/routes';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  constructor(private cookieService: CookieService, private router: Router) {}

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    req = this.addAuthenticationToken(req);

    return next.handle(req).pipe(
      catchError((error: HttpErrorResponse) => {
        if (error && error.status === 401 && !error.url.includes('/auth')) {
          this.router.navigateByUrl(PANEL_LOGIN);
        }
        return throwError(error);
      }),
    );
  }

  private addAuthenticationToken(req: HttpRequest<any>): HttpRequest<any> {
    if (
      this.cookieService.check(TOKEN_COOKIE_NAME) ||
      req.headers.has('Authorization') ||
      !req.url.startsWith(environment.apiUrl)
    ) {
      return req;
    } else {
      return req.clone({
        headers: req.headers.set(
          'Authorization',
          `Bearer ${this.cookieService.get(TOKEN_COOKIE_NAME)}`,
        ),
      });
    }
  }
}
