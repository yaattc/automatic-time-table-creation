import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {BehaviorSubject, Observable} from 'rxjs';
import { environment } from '../../environments/environment';
import { Class } from '../model/class';

@Injectable({
  providedIn: 'root',
})
export class TimeTableService {
  private timeSlotsSource$: BehaviorSubject<any[]> = new BehaviorSubject<any[]>([]);
  public timeSlots$ = this.timeSlotsSource$.asObservable();

  constructor(private http: HttpClient) {}

  postSchedule(group): void {
    const body = {
      from: '2020-06-30T22:01:53+06:00',
      till: '2020-09-30T22:01:53+06:00',
      group_id: group.group.id,
    };

    this.http
      .post<{ classes: Class[] }>(`${environment.apiUrl}/api/v1/classes`, body)
      .subscribe((response) => {
        this.timeSlotsSource$.next([...response.classes]);
      });
  }

  postCreateTimeTable(): Observable<any> {
    return this.http.post<any>(`${environment.apiUrl}/api/v1/generation`, null);
  }
}
