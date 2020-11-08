import { Component, OnInit } from '@angular/core';
import { FormBuilder, Validators } from '@angular/forms';

@Component({
  selector: 'app-course-creation',
  templateUrl: './course-creation.component.html',
  styleUrls: ['./course-creation.component.css'],
})
export class CourseCreationComponent implements OnInit {
  selectedProgram: string = null;

  programs: any[] = ['Bachelor', 'Master'];

  teachers: any[] = [{ name: 'Konyukhov' }, { name: 'Gorodetskiy' }, { name: 'Shilov' }];
  selectedTeachers: string[];

  creationForm = this.formBuilder.group({
    id: [null],
    name: [undefined, Validators.required],
    program: [this.selectedProgram, Validators.required],
    teachers: [[], Validators.required],
  });

  constructor(private formBuilder: FormBuilder) {}

  ngOnInit(): void {
    this.selectedProgram = this.programs[0];
  }

  public setSelectedTeachers(val: any[]): void {
    // restore original order
    this.selectedTeachers = this.teachers.filter((teacher) => val.includes(teacher));
  }

  submit(): void {}
}
