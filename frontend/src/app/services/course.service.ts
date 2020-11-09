import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class CourseService {
  constructor(private http: HttpClient) {}

  public getListOfTeachers(): Observable<any> {
    return this.http.get<any>(`${environment.apiUrl}/api/v1/teacher`);
  }
}
