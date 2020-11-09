import { Component, OnInit } from '@angular/core';
import { FormBuilder, Validators } from '@angular/forms';
import { CourseService } from '../../../services/course.service';

@Component({
  selector: 'app-course-creation',
  templateUrl: './course-creation.component.html',
  styleUrls: ['./course-creation.component.css'],
})
export class CourseCreationComponent implements OnInit {
  selectedProgram: string = null;

  programs: any[] = ['Bachelor', 'Master'];

  teachers: any[] = [];
  selectedTeachers: string[];

  creationForm = this.formBuilder.group({
    id: [null],
    name: [undefined, Validators.required],
    program: [this.selectedProgram, Validators.required],
    teachers: [[], Validators.required],
  });

  constructor(private formBuilder: FormBuilder, private courseService: CourseService) {}

  ngOnInit(): void {
    this.selectedProgram = this.programs[0];
    this.courseService.getListOfTeachers().subscribe((value) => {
      this.teachers = value.teachers.map((val) => {
        return {
          name: val.degree + ' ' + val.name + ' ' + val.surname,
        };
      });
    });
  }

  public setSelectedTeachers(val: any[]): void {
    // restore original order
    if (val !== undefined) {
      this.selectedTeachers = this.teachers.filter((teacher) => val.includes(teacher));
    }
  }

  submit(): void {}
}
