import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { PANEL, PANEL_LOGIN } from './constants/routes';

const routes: Routes = [
  {
    path: PANEL,
    loadChildren: () =>
      import('./module/control-panel/control-panel.module').then((m) => m.ControlPanelModule),
  },
  // { path: '**', redirectTo: PANEL_LOGIN },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
