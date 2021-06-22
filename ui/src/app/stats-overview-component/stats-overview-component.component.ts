import {Component, OnInit} from '@angular/core';
import {Cluster, Stats, TrinoApiService} from "../service/api.service";
import {combineLatest, interval, Observable} from "rxjs";
import {flatMap} from "rxjs/internal/operators";
import {map} from "rxjs/operators";

@Component({
  selector: 'app-stats-overview-component',
  templateUrl: './stats-overview-component.component.html',
  styleUrls: ['./stats-overview-component.component.scss']
})
export class StatsOverviewComponentComponent implements OnInit {
  data?: Observable<Result>;

  constructor(private readonly api: TrinoApiService) {
  }

  ngOnInit(): void {
    this.data = interval(2000).pipe(
      flatMap(_ => combineLatest([this.api.stats(), this.api.clusters()])),
      map(([s, c]) => { 
        return {stats: s, clusters: c}
      })
    );
  }

}


export interface Result {


  stats: Stats,
  clusters: Cluster[]
}
