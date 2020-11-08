import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { LoginComponent } from './login/login.component';
import { ControlPanelRoutingModule } from './control-panel-routing.module';
import { ButtonModule } from 'primeng/button';
import { InputTextModule } from 'primeng/inputtext';
import { DashboardComponent } from './dashboard/dashboard.component';
import { CreationProfTaComponent } from './creation-prof-ta/creation-prof-ta.component';
import { ControlPanelComponent } from './control-panel.component';
import { RadioButtonModule } from 'primeng/radiobutton';
import { InputTextareaModule } from 'primeng/inputtextarea';
import { CourseCreationComponent } from './course-creation/course-creation.component';
import { MultiSelectModule } from 'primeng/multiselect';

@NgModule({
  declarations: [
    ControlPanelComponent,
    LoginComponent,
    DashboardComponent,
    CreationProfTaComponent,
    CourseCreationComponent,
  ],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    ControlPanelRoutingModule,
    ButtonModule,
    InputTextModule,
    FormsModule,
    RadioButtonModule,
    InputTextareaModule,
    MultiSelectModule,
  ],
})
export class ControlPanelModule {}
