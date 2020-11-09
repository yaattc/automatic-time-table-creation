import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { CreationTeacherModel } from '../model/creation-teacher-model';
import {Observable} from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class TeacherService {
  constructor(private http: HttpClient) {}

  createTeacher(teacher: CreationTeacherModel): Observable<any> {
    return this.http.post<any>(`${environment.apiUrl}/api/v1/teacher`, teacher);
  }
}
