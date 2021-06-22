import {Component, OnInit} from '@angular/core';
import {Cluster, TrinoApiService} from "../services/api.service";
import {interval, Observable} from "rxjs";
import {map} from "rxjs/operators";
import {flatMap} from "rxjs/internal/operators";

@Component({
  selector: 'app-cluster-list-component',
  templateUrl: './cluster-list-component.component.html',
  styleUrls: ['./cluster-list-component.component.scss']
})
export class ClusterListComponentComponent implements OnInit {

  clusters?: Observable<Cluster[]>

  constructor(private readonly api: TrinoApiService) {
  }

  ngOnInit(): void {
    this.clusters = interval(2000).pipe(
      flatMap(_ => this.api.clusters()),
      map(cl => cl.sort((a, b) => a.name.localeCompare(b.name)))
    );
  }

}
