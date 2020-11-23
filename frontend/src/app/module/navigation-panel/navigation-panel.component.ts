import { Component, OnInit } from '@angular/core';
import {
  PANEL_CREATION_COURSE,
  PANEL_CREATION_GROUP,
  PANEL_CREATION_SPECIALIST,
  PANEL_SPECIALIST_PREFERENCES,
} from '../../constants/routes';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { environment } from '../../../environments/environment';
import { FormBuilder } from '@angular/forms';
import { GroupYearService } from '../../services/group-year.service';
import { Group } from '../../model/group';
import { TimeTableService } from '../../services/time-table.service';
import { MessageService } from 'primeng/api';

@Component({
  selector: 'app-navigation-panel',
  templateUrl: './navigation-panel.component.html',
  styleUrls: ['./navigation-panel.component.css'],
  providers: [MessageService],
})
export class NavigationPanelComponent implements OnInit {
  pages = [
    {
      label: 'Create teacher',
      path: PANEL_CREATION_SPECIALIST,
    },
    {
      label: 'Create course',
      path: PANEL_CREATION_COURSE,
    },
    {
      label: 'Set up preferences',
      path: PANEL_SPECIALIST_PREFERENCES,
    },
    {
      label: 'Creation study year/group',
      path: PANEL_CREATION_GROUP,
    },
  ];

  currentPage: string;

  filterForm = this.formBuilder.group({
    group: [null],
  });

  groups: Group[];

  selectedGroups: any;

  constructor(
    private router: Router,
    private authService: AuthService,
    private formBuilder: FormBuilder,
    private groupYearService: GroupYearService,
    private timeTableService: TimeTableService,
    private messageService: MessageService,
  ) {}

  ngOnInit(): void {
    this.authService.currentPage$.subscribe((page) => {
      this.currentPage = page;
    });
    this.groupYearService.getGroup().subscribe((val) => {
      this.groups = [...val.groups];
    });
  }

  moveTo(path: string): void {
    this.router.navigateByUrl(path);
    this.authService.setCurrentPage(`${environment.apiUrl}/${path}`);
  }

  submitTimeTable(): void {
    this.timeTableService.postSchedule(this.filterForm.value);
  }

  submitCreateTimeTable(): void {
    this.timeTableService.postCreateTimeTable().subscribe(
      (response) =>
        this.messageService.add({
          severity: 'success',
          summary: 'Success',
          detail: 'Time Table has been generated',
        }),
      (error) =>
        this.messageService.add({ severity: 'error', summary: 'Error', detail: 'Smth strange' }),
    );
  }
}
