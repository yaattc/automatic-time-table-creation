import { Component, OnInit } from '@angular/core';
import { PANEL_CREATION_COURSE, PANEL_CREATION_SPECIALIST } from '../../constants/routes';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { environment } from '../../../environments/environment';

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
  ];

  currentPage: string;

  constructor(private router: Router, private authService: AuthService) {}

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
