import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class CourseService {
  constructor(private http: HttpClient) {}

  postCreationCourse(value): Observable<any> {
    const body = {
      name: value.name,
      program: value.program,
      primary_lector: value.primary_lector.value.id,
      assistant_lector: value.assistant_lector.value.id,
      teacher_assistants: [],
    };

    value.teacher_assistants.forEach((val) => {
      body.teacher_assistants.push(val.value.id);
    });

    if (body.assistant_lector === null) {
      body.assistant_lector = '';
    }

    return this.http.post<any>(`${environment.apiUrl}/api/v1/course`, body);
  }
}
