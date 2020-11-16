import { Component, OnInit } from '@angular/core';
import { CalendarOptions } from '@fullcalendar/angular';
import { FormBuilder, Validators } from '@angular/forms';

@Component({
  selector: 'app-time-table',
  templateUrl: './time-table.component.html',
  styleUrls: ['./time-table.component.css'],
})
export class TimeTableComponent implements OnInit {
  events: any[];

  calendarOptions: CalendarOptions;

  constructor() {}

  ngOnInit(): void {
    this.calendarOptions = {
      initialView: 'timeGridWeek',
      headerToolbar: {
        left: 'prev,next',
        center: 'title',
        right: 'dayGridMonth,timeGridWeek,timeGridDay',
      },
    };
  }
}
