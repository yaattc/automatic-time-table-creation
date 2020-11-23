import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';
import { Group } from '../model/group';

@Injectable({
  providedIn: 'root',
})
export class GroupYearService {
  constructor(private http: HttpClient) {}

  postStudyYear(name: string): Observable<{ id: string; name: string }> {
    return this.http.post<{ id: string; name: string }>(`${environment.apiUrl}/api/v1/study_year`, {
      name,
    });
  }

  getStudyYears(): Observable<{ study_years: { id: string; name: string }[] }> {
    return this.http.get<{ study_years: { id: string; name: string }[] }>(
      `${environment.apiUrl}/api/v1/study_year`,
    );
  }

  postGroup(value: any): Observable<any> {
    const body = {
      name: value.name,
      study_year_id: value.study_year_id.id,
    };

    return this.http.post<any>(`${environment.apiUrl}/api/v1/group`, body);
  }

  getGroup(): Observable<{ groups: Group[] }> {
    return this.http.get<{ groups: Group[] }>(`${environment.apiUrl}/api/v1/group`);
  }
}
