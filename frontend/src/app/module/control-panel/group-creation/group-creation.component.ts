import { Component, OnInit } from '@angular/core';
import { FormBuilder } from '@angular/forms';

@Component({
  selector: 'app-group-creation',
  templateUrl: './group-creation.component.html',
  styleUrls: ['./group-creation.component.css'],
})
export class GroupCreationComponent implements OnInit {
  group = this.formBuilder.group({
    name: [null],
    study_year_id: [null],
  });

  year = this.formBuilder.group({
    name: [null],
  });

  years: any[];
  selectedYears: any[];

  constructor(private formBuilder: FormBuilder) {}

  ngOnInit(): void {}

  public setSelectedYears(val: any[]): void {
    // restore original order
    if (val !== undefined) {
      this.selectedYears = this.years.filter((year) => val.includes(year));
    }
  }

  submitYear(): void {}

  submitGroup(): void {}
}
