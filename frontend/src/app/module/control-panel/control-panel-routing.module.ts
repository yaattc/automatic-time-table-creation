import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { DASHBOARD, LOGIN } from '../../constants/routes';
import { LoginComponent } from './login/login.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { AuthService } from '../../services/auth.service';

const routes: Routes = [
  {
    path: LOGIN,
    component: LoginComponent,
  },
  {
    path: DASHBOARD,
    component: DashboardComponent,
    canActivate: [AuthService],
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ControlPanelRoutingModule {}
