import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ControlPanelComponent } from './module/control-panel/control-panel.component';
import { PANEL } from './constants/routes';

const routes: Routes = [
  {
    path: PANEL,
    loadChildren: () => import('./module/control-panel/control-panel.module').then(m => m.ControlPanelModule)
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
