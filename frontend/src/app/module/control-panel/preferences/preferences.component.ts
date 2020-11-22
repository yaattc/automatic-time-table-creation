import { Component, OnInit } from '@angular/core';
import { FormBuilder } from '@angular/forms';
import { TeacherService } from '../../../services/teacher.service';
import { TimeSlot } from '../../../model/time-slot';
import { Staff } from '../../../model/staff';
import { PreferencesService } from '../../../services/preferences.service';
import { MessageService } from 'primeng/api';

@Component({
  selector: 'app-preferences',
  templateUrl: './preferences.component.html',
  styleUrls: ['./preferences.component.css'],
  providers: [MessageService],
})
export class PreferencesComponent implements OnInit {
  constructor(
    private formBuilder: FormBuilder,
    private teacherService: TeacherService,
    private preferencesService: PreferencesService,
    private messageService: MessageService,
  ) {}

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

    this.preferencesService.getTimeSlots().subscribe((value) => {
      this.timeSlots = value.time_slots.map((val) => {
        return {
          name: val.weekday + ' ' + val.start + ' ' + val.duration,
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
              timeSlot.value.weekday + ' ' + timeSlot.value.start + ' ' + timeSlot.value.duration,
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
    if (val !== undefined) {
      this.selectedRooms = this.rooms.filter((room) => val.includes(room));
    }
  }

  public submit(): void {
    this.messageService.clear();
    this.preferencesService.postTeacherPreferences(this.creationForm.value).subscribe(
      (response) =>
        this.messageService.add({
          severity: 'success',
          summary: 'Success',
          detail: 'Group has been added',
        }),
      (error) =>
        this.messageService.add({ severity: 'error', summary: 'Error', detail: 'Smth strange' }),
    );
    this.timeSlots = [];
    this.selectedRooms = [];
    this.selectedTeachers = [];
    this.creationForm.reset();
  }
}
