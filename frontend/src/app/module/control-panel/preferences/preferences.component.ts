import { Component, OnInit } from '@angular/core';
import { FormBuilder, Validators } from '@angular/forms';
import { TeacherService } from '../../../services/teacher.service';
import { TimeSlot } from '../../../model/time-slot';
import { Staff } from '../../../model/staff';

@Component({
  selector: 'app-preferences',
  templateUrl: './preferences.component.html',
  styleUrls: ['./preferences.component.css'],
})
export class PreferencesComponent implements OnInit {
  constructor(private formBuilder: FormBuilder, private teacherService: TeacherService) {}

  timeSlots: TimeSlot[] = [];
  selectedTimeSlots: { name: string }[];

  teachers: Staff[] = [];
  selectedTeachers: { name: string }[];

  rooms: string[] = ['room #108', 'room #109', 'room #231'];
  selectedRooms: string[];

  creationForm = this.formBuilder.group({
    time_slots: [[] as TimeSlot[]],
    staff: [[] as Staff[]],
    locations: [[]],
  });

  ngOnInit(): void {
    this.teacherService.getListOfTeachers().subscribe((value) => {
      this.teachers = value.teachers.map((val) => {
        return {
          name: val.degree + ' ' + val.name + ' ' + val.surname,
        };
      });
    });
  }

  public setSelectedTimeSlots(val: TimeSlot[]): void {
    // restore original order
    if (val !== undefined) {
      this.selectedTimeSlots = this.timeSlots
        .filter((timeSlot) => val.includes(timeSlot))
        .map((timeSlot) => {
          return {
            name:
              timeSlot.weekday +
              ' ' +
              timeSlot.start +
              ' ' +
              timeSlot.duration +
              ' ' +
              timeSlot.location,
          };
        });
    }
  }

  public setSelectedTeachers(val: Staff[]): void {
    // restore original order
    if (val !== undefined) {
      this.selectedTeachers = this.teachers.filter((teacher) => val.includes(teacher));
    }
  }

  public setSelectedRooms(val: string[]): void {
    // restore original order
    if (val !== undefined) {
      this.selectedRooms = this.rooms.filter((room) => val.includes(room));
    }
  }
}
