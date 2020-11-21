import { Component, OnInit } from '@angular/core';
import { FormBuilder } from '@angular/forms';
import { GroupYearService } from '../../../services/group-year.service';
import { MessageService } from 'primeng/api';

@Component({
  selector: 'app-group-creation',
  templateUrl: './group-creation.component.html',
  styleUrls: ['./group-creation.component.css'],
  providers: [MessageService],
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

  constructor(
    private formBuilder: FormBuilder,
    private groupYearService: GroupYearService,
    private messageService: MessageService,
  ) {}

  ngOnInit(): void {
    this.groupYearService.getStudyYears().subscribe((response) => {
      this.years = response;
    });
  }

  public setSelectedYears(val: any[]): void {
    // restore original order
    if (val !== undefined) {
      this.selectedYears = this.years.filter((year) => val.includes(year));
    }
  }

  submitYear(): void {
    this.messageService.clear();
    this.groupYearService.postStudyYear(this.year.value.name).subscribe(
      (response) =>
        this.messageService.add({
          key: 'tl',
          severity: 'success',
          summary: 'Success',
          detail: 'Study year has been added',
        }),
      (error) =>
        this.messageService.add({
          key: 'tl',
          severity: 'error',
          summary: 'Error',
          detail: 'Smth strange',
        }),
    );
    this.year.reset();
  }

  submitGroup(): void {
    this.messageService.clear();
    this.groupYearService.postGroup(this.group.value).subscribe(
      (response) =>
        this.messageService.add({
          severity: 'success',
          summary: 'Success',
          detail: 'Group has been added',
        }),
      (error) =>
        this.messageService.add({ severity: 'error', summary: 'Error', detail: 'Smth strange' }),
    );
    this.group.reset();
  }
}
