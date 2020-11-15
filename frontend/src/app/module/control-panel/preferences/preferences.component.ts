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

  timeSlots: { name: string; value: TimeSlot }[] = [];
  selectedTimeSlots: { name: string; value: TimeSlot }[];

  teachers: { name: string; value: Staff }[] = [];
  selectedTeachers: { name: string; value: Staff }[];
  selectedTeacherForPreferences: Staff;

  rooms: { name: string }[] = [{ name: 'room #108' }, { name: 'room #109' }, { name: 'room #231' }];
  selectedRooms: { name: string }[];

  creationForm = this.formBuilder.group({
    teacher: [undefined as Staff],
    timeSlots: [[] as TimeSlot[]],
    staff: [[] as Staff[]],
    locations: [[]],
  });

  ngOnInit(): void {
    this.teacherService.getListOfTeachers().subscribe((value) => {
      this.teachers = value.teachers.map((val) => {
        return {
          name: val.degree + ' ' + val.name + ' ' + val.surname,
          value: val,
        };
      });
    });
  }

  public setSelectedTimeSlots(val: { name: string; value: TimeSlot }[]): void {
    // restore original order
    if (val !== undefined) {
      this.selectedTimeSlots = this.timeSlots
        .filter((timeSlot) => val.includes(timeSlot))
        .map((timeSlot) => {
          return {
            name:
              timeSlot.value.weekday +
              ' ' +
              timeSlot.value.start +
              ' ' +
              timeSlot.value.duration +
              ' ' +
              timeSlot.value.location,
            value: timeSlot.value,
          };
        });
    }
  }

  public setSelectedTeachers(val: { name: string; value: Staff }[]): void {
    // restore original order
    if (val !== undefined) {
      this.selectedTeachers = this.teachers.filter((teacher) => val.includes(teacher));
    }
  }

  public setSelectedRooms(val: { name: string }[]): void {
    // restore original order
    console.log(val);
    if (val !== undefined) {
      this.selectedRooms = this.rooms.filter((room) => val.includes(room));
    }
  }

  public submit(): void {}
}
