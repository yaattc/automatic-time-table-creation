import { Component, OnInit } from '@angular/core';
import { FormBuilder, Validators } from '@angular/forms';
import { TeacherService } from '../../../services/teacher.service';
import { Staff } from '../../../model/staff';
import { CourseService } from '../../../services/course.service';
import { MessageService } from 'primeng/api';

@Component({
  selector: 'app-course-creation',
  templateUrl: './course-creation.component.html',
  styleUrls: ['./course-creation.component.css'],
  providers: [MessageService],
})
export class CourseCreationComponent implements OnInit {
  selectedProgram: string = null;

  programs: any[] = ['Bachelor', 'Master'];

  teachers: any[] = [];
  selectedTeacherAssistant: Staff[];

  creationForm = this.formBuilder.group({
    name: [undefined, Validators.required],
    program: [this.selectedProgram, Validators.required],
    primary_lector: [undefined, Validators.required],
    assistant_lector: [undefined],
    teacher_assistants: [[], Validators.required],
  });

  constructor(
    private formBuilder: FormBuilder,
    private teacherService: TeacherService,
    private courseService: CourseService,
    private messageService: MessageService,
  ) {}

  ngOnInit(): void {
    this.selectedProgram = this.programs[0];
    this.teacherService.getListOfTeachers().subscribe((value) => {
      this.teachers = value.teachers.map((val) => {
        return {
          name: val.degree + ' ' + val.name + ' ' + val.surname,
          value: val,
        };
      });
    });
  }

  public setSelectedTeacherAssistant(val: any[]): void {
    // restore original order
    if (val !== null) {
      this.selectedTeacherAssistant = this.teachers.filter((teacher) => val.includes(teacher));
    }
  }

  submit(): void {
    this.courseService.postCreationCourse(this.creationForm.value).subscribe(
      (response) =>
        this.messageService.add({
          severity: 'success',
          summary: 'Success',
          detail: 'Course has been added',
        }),
      (error) =>
        this.messageService.add({ severity: 'error', summary: 'Error', detail: 'Smth strange' }),
    );
    this.creationForm.reset();
    this.selectedTeacherAssistant = [];
    this.selectedProgram = this.programs[0];
    this.creationForm.patchValue({
      program: this.selectedProgram,
    });
  }
}
