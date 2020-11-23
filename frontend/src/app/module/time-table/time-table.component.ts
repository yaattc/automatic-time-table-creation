import { Component, OnInit } from '@angular/core';
import { CalendarOptions } from '@fullcalendar/angular';
import { TimeTableService } from '../../services/time-table.service';

@Component({
  selector: 'app-time-table',
  templateUrl: './time-table.component.html',
  styleUrls: ['./time-table.component.css'],
})
export class TimeTableComponent implements OnInit {
  events: any[];

  calendarOptions: CalendarOptions;

  constructor(private timeTableService: TimeTableService) {}

  ngOnInit(): void {
    this.calendarOptions = {
      initialView: 'timeGridWeek',
      eventColor: '#80BC00',
      headerToolbar: {
        left: 'prev,next',
        center: 'title',
        right: 'dayGridMonth,timeGridWeek,timeGridDay',
      },
    };

    this.timeTableService.timeSlots$.subscribe((val) => {
      this.calendarOptions.events = val;
    });
  }
}
