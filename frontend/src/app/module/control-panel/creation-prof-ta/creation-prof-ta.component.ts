import { Component, OnInit } from '@angular/core';
import { FormBuilder, Validators } from '@angular/forms';
import { TeacherService } from '../../../services/teacher.service';

@Component({
  selector: 'app-creation-prof-ta',
  templateUrl: './creation-prof-ta.component.html',
  styleUrls: ['./creation-prof-ta.component.css'],
})
export class CreationProfTaComponent implements OnInit {
  selectedRole: string = null;

  roles: any[] = ['Professor', 'TA'];

  creationForm = this.formBuilder.group({
    id: [null],
    name: [undefined, Validators.required],
    surname: [undefined, Validators.required],
    email: [null],
    degree: [undefined, Validators.required],
    about: [undefined],
  });

  constructor(private formBuilder: FormBuilder, private teacherService: TeacherService) {}

  ngOnInit(): void {
    this.selectedRole = this.roles[0];
  }

  submit(): void {
    this.teacherService.createTeacher(this.creationForm.value);
  }
}
