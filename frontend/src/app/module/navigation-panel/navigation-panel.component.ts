import { Component, OnInit } from '@angular/core';
import {
  PANEL_CREATION_COURSE, PANEL_CREATION_GROUP,
  PANEL_CREATION_SPECIALIST,
  PANEL_SPECIALIST_PREFERENCES,
} from '../../constants/routes';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { environment } from '../../../environments/environment';
import { FormBuilder } from '@angular/forms';

@Component({
  selector: 'app-navigation-panel',
  templateUrl: './navigation-panel.component.html',
  styleUrls: ['./navigation-panel.component.css'],
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
    year: [null],
    group: [null],
  });

  groups: any[];
  years: any[];

  selectedGroups: any;
  selectedYears: any;

  constructor(
    private router: Router,
    private authService: AuthService,
    private formBuilder: FormBuilder,
  ) {}

  ngOnInit(): void {
    this.authService.currentPage$.subscribe((page) => {
      this.currentPage = page;
    });
  }

  moveTo(path: string): void {
    this.router.navigateByUrl(path);
    this.authService.setCurrentPage(`${environment.apiUrl}/${path}`);
  }
}
