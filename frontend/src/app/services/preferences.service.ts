import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Staff } from '../model/staff';
import { environment } from '../../environments/environment';
import { TimeSlot } from '../model/time-slot';
@Injectable({
  providedIn: 'root',
})
export class PreferencesService {
  constructor(private http: HttpClient) {}

  postTeacherPreferences(value): Observable<Staff> {
    const body: any = {
      time_slots: [],
      staff: [],
      locations: [],
    };
    if (value.timeSlots !== undefined) {
      value.timeSlots.forEach((val) => body.time_slots.push(val.value));
    }
    if (value.staff !== undefined) {
      value.staff.forEach((val) => body.staff.push(val.value));
    }
    if (value.locations !== undefined) {
      value.locations.forEach((val) => body.locations.push(val.name));
    }

    return this.http.post<Staff>(
      `${environment.apiUrl}/api/v1/teacher/${value.teacher.value.id}/preferences`,
      body,
    );
  }

  getTimeSlots(): Observable<{ time_slots: TimeSlot[] }> {
    return this.http.get<{ time_slots: TimeSlot[] }>(`${environment.apiUrl}/api/v1/time_slot`);
  }
}
