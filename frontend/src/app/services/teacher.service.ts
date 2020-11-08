import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { CreationTeacherModel } from '../model/creation-teacher-model';

@Injectable({
  providedIn: 'root',
})
export class TeacherService {
  constructor(private http: HttpClient) {}

  createTeacher(teacher: CreationTeacherModel): void {
    this.http.post<any>(`${environment.apiUrl}/api/v1/teacher`, teacher).subscribe(
      (response) => console.log(response),
      (error) => console.log(error),
    );
  }
}
