import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class GroupYearService {
  constructor(private http: HttpClient) {}

  postStudyYear(name: string): Observable<{ id: string; name: string }> {
    return this.http.post<{ id: string; name: string }>(`${environment.apiUrl}/study_year`, {
      name,
    });
  }

  getStudyYears(): Observable<[{ id: string; name: string }]> {
    return this.http.get<[{ id: string; name: string }]>(`${environment.apiUrl}/study_year`);
  }

  postGroup(value: any): Observable<any> {
    return this.http.post<any>(`${environment.apiUrl}/group`, value);
  }
}
