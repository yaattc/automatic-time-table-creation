import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppComponent } from './app.component';
import { AppRoutingModule } from './app-routing.module';
import { HttpClientModule } from '@angular/common/http';
import { NavigationPanelComponent } from './module/navigation-panel/navigation-panel.component';
import {ButtonModule} from 'primeng/button';
import {RippleModule} from 'primeng/ripple';

@NgModule({
  declarations: [AppComponent, NavigationPanelComponent],
  imports: [BrowserModule, AppRoutingModule, HttpClientModule, ButtonModule, RippleModule],
  bootstrap: [AppComponent],
})
export class AppModule {}
