import { Component, OnInit } from '@angular/core';
import { FormBuilder, Validators } from '@angular/forms';

@Component({
  selector: 'app-creation-prof-ta',
  templateUrl: './creation-prof-ta.component.html',
  styleUrls: ['./creation-prof-ta.component.css'],
})
export class CreationProfTaComponent implements OnInit {
  selectedRole: string = null;

  roles: any[] = ['Professor', 'TA'];

  creationForm = this.formBuilder.group({
    name: [undefined, Validators.required],
    surname: [undefined, Validators.required],
    education: [undefined, Validators.required],
    about: [undefined],
  });

  constructor(private formBuilder: FormBuilder) {}

  ngOnInit(): void {
    this.selectedRole = this.roles[0];
  }

  submit(): void {}
}
