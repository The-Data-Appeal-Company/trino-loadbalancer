import {NgModule} from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {MatToolbarModule} from "@angular/material/toolbar";
import {MatIconModule} from "@angular/material/icon";
import {MatButtonModule} from "@angular/material/button";
import {MatSliderModule} from "@angular/material/slider";
import {MatCardModule} from "@angular/material/card";
import {StatsCardComponent} from './stats-card/stats-card.component';
import {StatsOverviewComponentComponent} from './stats-overview-component/stats-overview-component.component';
import {MatGridListModule} from "@angular/material/grid-list";
import {HttpClientModule} from "@angular/common/http";
import {ClusterItemComponentComponent} from './cluster-item-component/cluster-item-component.component';
import {ClusterListComponentComponent} from './cluster-list-component/cluster-list-component.component';
import {MatButtonToggleModule} from "@angular/material/button-toggle";
import {MatSlideToggleModule} from "@angular/material/slide-toggle";
import {MatChipsModule} from "@angular/material/chips";
import { OverviewComponent } from './overview-component/overview.component';

@NgModule({
  declarations: [
    AppComponent,
    StatsCardComponent,
    StatsOverviewComponentComponent,
    ClusterItemComponentComponent,
    ClusterListComponentComponent,
    OverviewComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    MatToolbarModule,
    MatIconModule,
    MatButtonModule,
    MatSliderModule,
    MatCardModule,
    MatGridListModule,
    MatSlideToggleModule,
    HttpClientModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {
}
