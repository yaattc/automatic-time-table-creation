import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppComponent } from './app.component';
import { AppRoutingModule } from './app-routing.module';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';
import { NavigationPanelComponent } from './module/navigation-panel/navigation-panel.component';
import { ButtonModule } from 'primeng/button';
import { RippleModule } from 'primeng/ripple';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { AuthInterceptor } from './interceptor/auth.interceptor';
import { TimeTableComponent } from './module/time-table/time-table.component';

import { FullCalendarModule } from '@fullcalendar/angular';
import dayGridPlugin from '@fullcalendar/daygrid';
import timeGridPlugin from '@fullcalendar/timegrid';
import interactionPlugin from '@fullcalendar/interaction';
import { CalendarModule } from 'primeng/calendar';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import { DropdownModule } from 'primeng/dropdown';

FullCalendarModule.registerPlugins([dayGridPlugin, timeGridPlugin, interactionPlugin]);

@NgModule({
  declarations: [AppComponent, NavigationPanelComponent, TimeTableComponent],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    AppRoutingModule,
    HttpClientModule,
    ButtonModule,
    RippleModule,
    FullCalendarModule,
    CalendarModule,
    ReactiveFormsModule,
    DropdownModule,
    FormsModule,
  ],
  providers: [
    {
      provide: HTTP_INTERCEPTORS,
      useClass: AuthInterceptor,
      multi: true,
    },
  ],
  bootstrap: [AppComponent],
})
export class AppModule {}
