import {Component, Input, OnInit} from '@angular/core';
import {Cluster, TrinoApiService} from "../services/api.service";
import {MatSlideToggleChange} from "@angular/material/slide-toggle";

@Component({
  selector: 'app-cluster-item-component',
  templateUrl: './cluster-item-component.component.html',
  styleUrls: ['./cluster-item-component.component.scss']
})
export class ClusterItemComponentComponent implements OnInit {

  @Input() item?: Cluster;

  constructor(private readonly api: TrinoApiService) {
  }

  ngOnInit(): void {

  }

  clusterStatusChange(e: MatSlideToggleChange) {
    this.api.updateClusterStatus({
      enabled: e.checked,
      name: this.item?.name ? this.item?.name : "",
    }).subscribe(e => console.log(e));
  }

}
