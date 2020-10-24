import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { LoginComponent } from './login/login.component';
import { ControlPanelRoutingModule } from './control-panel-routing.module';
import { ButtonModule } from 'primeng/button';
import { InputTextModule } from 'primeng/inputtext';
import { DashboardComponent } from './dashboard/dashboard.component';

@NgModule({
  declarations: [LoginComponent, DashboardComponent],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    ControlPanelRoutingModule,
    ButtonModule,
    InputTextModule,
    FormsModule,
  ],
})
export class ControlPanelModule {}
