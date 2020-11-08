import { Component, OnInit } from '@angular/core';
import { PANEL_CREATION_SPECIALIST } from '../../constants/routes';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';

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
  ];

  constructor(private router: Router, private authService: AuthService) {}

  ngOnInit(): void {}

  moveTo(path: string): void {
    this.router.navigateByUrl(path);
  }
}
