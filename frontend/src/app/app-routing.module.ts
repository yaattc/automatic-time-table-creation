import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { PANEL, PANEL_LOGIN, TIMETABLE } from './constants/routes';
import { TimeTableComponent } from './module/time-table/time-table.component';

const routes: Routes = [
  {
    path: PANEL,
    loadChildren: () =>
      import('./module/control-panel/control-panel.module').then((m) => m.ControlPanelModule),
  },
  {
    path: TIMETABLE,
    component: TimeTableComponent,
  },
  { path: '**', redirectTo: TIMETABLE },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
