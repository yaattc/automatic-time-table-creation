import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {
  CREATION_COURSE,
  CREATION_GROUP,
  CREATION_SPECIALIST,
  DASHBOARD,
  LOGIN,
  SPECIALIST_PREFERENCES,
} from '../../constants/routes';
import { LoginComponent } from './login/login.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { ControlPanelComponent } from './control-panel.component';
import { CreationProfTaComponent } from './creation-prof-ta/creation-prof-ta.component';
import { CourseCreationComponent } from './course-creation/course-creation.component';
import { AuthService } from '../../services/auth.service';
import { PreferencesComponent } from './preferences/preferences.component';
import { GroupCreationComponent } from './group-creation/group-creation.component';

const routes: Routes = [
  {
    path: '',
    component: ControlPanelComponent,
    children: [
      {
        path: CREATION_SPECIALIST,
        component: CreationProfTaComponent,
        canActivate: [AuthService],
      },
      {
        path: LOGIN,
        component: LoginComponent,
      },
      {
        path: DASHBOARD,
        component: DashboardComponent,
        canActivate: [AuthService],
      },
      {
        path: CREATION_COURSE,
        component: CourseCreationComponent,
        canActivate: [AuthService],
      },
      {
        path: CREATION_COURSE,
        component: CourseCreationComponent,
        canActivate: [AuthService],
      },
      {
        path: SPECIALIST_PREFERENCES,
        component: PreferencesComponent,
        canActivate: [AuthService],
      },
      {
        path: CREATION_GROUP,
        component: GroupCreationComponent,
        canActivate: [AuthService],
      },
      // {
      //   path: '',
      //   redirectTo: LOGIN,
      //   // canActivate: [AuthService]
      // },
    ],
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ControlPanelRoutingModule {}
