import {Component, OnInit} from '@angular/core';
import Chart from 'chart.js';
import {Cluster, PrestoApiService, Stats} from '../../services/api/prestoApiService';
import {interval, Subject} from 'rxjs';
import {flatMap} from 'rxjs/operators';
import {MatSlideToggleChange} from '@angular/material/slide-toggle';
import {MatSnackBar} from '@angular/material/snack-bar';
import {MatDialog} from '@angular/material/dialog';


@Component({
  selector: 'dashboard-cmp',
  moduleId: module.id,
  templateUrl: 'dashboard.component.html',
  styleUrls: ['dashboard.component.css'],
})

export class DashboardComponent implements OnInit {


  public updatingState: boolean;

  public canvas: any;
  public ctx;
  public chartColor;
  public chartHours;

  public prestoStatistics: Stats;
  public clusterState: Cluster[];

  public runningQueries = []

  public blockedQueries = []

  public forceRefresh: Subject<any>;

  constructor(private prestoApi: PrestoApiService, private _snackBar: MatSnackBar, public dialog: MatDialog) {
  }

  ngOnInit() {

    this.forceRefresh = new Subject<any>();

    this.updatingState = false;

    this.clusterState = []

    this.prestoStatistics = {
      blocked_queries: 0,
      queued_queries: 0,
      running_queries: 0,
      total_workers: 0,
    }

    interval(2000).subscribe(
      () => this.forceRefresh.next()
    )

    this.forceRefresh.asObservable()
      .pipe(
        flatMap(_ => this.prestoApi.stats())
      )
      .subscribe(
        state => {
          this.prestoStatistics = state

          const running = state.running_queries;
          if (this.runningQueries.length > 9) {
            this.runningQueries.shift()
          }

          this.runningQueries.push({
            x: new Date(),
            y: running
          })


          const blocked = state.blocked_queries;
          if (this.blockedQueries.length > 9) {
            this.blockedQueries.shift()
          }

          this.blockedQueries.push({
            x: new Date(),
            y: blocked
          });

          this.chartHours.update()
        }
      )

    this.forceRefresh.asObservable().pipe(
      flatMap(_ => this.prestoApi.clusters())
    ).subscribe(
      res => this.clusterState = res,
      err => this.openSnackBar('error refreshing cluster topology :(', 'ok')
    );

    this.chartColor = '#FFFFFF';

    this.canvas = document.getElementById('chartHours');
    this.ctx = this.canvas.getContext('2d');

    this.chartHours = new Chart(this.ctx, {
      type: 'line',

      data: {
        // labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct'],
        datasets: [{
          borderColor: '#6bd098',
          backgroundColor: '#6bd098',
          pointRadius: 0,
          pointHoverRadius: 0,
          borderWidth: 3,
          data: this.runningQueries
        },
          {
            borderColor: '#f17e5d',
            backgroundColor: '#f17e5d',
            pointRadius: 0,
            pointHoverRadius: 0,
            borderWidth: 3,
            data: this.blockedQueries
          },
          {
            borderColor: '#fcc468',
            backgroundColor: '#fcc468',
            pointRadius: 0,
            pointHoverRadius: 0,
            borderWidth: 3,
            data: []
          }
        ]
      },
      options: {
        legend: {
          display: false
        },

        tooltips: {
          enabled: false
        },

        scales: {
          yAxes: [{

            ticks: {
              fontColor: '#9f9f9f',
              beginAtZero: false,
              maxTicksLimit: 5,
              suggestedMin: 0,
              suggestedMax: 100,
              //padding: 20
            },
            gridLines: {
              drawBorder: false,
              zeroLineColor: '#ccc',
              color: 'rgba(255,255,255,0.05)'
            }

          }],

          xAxes: [{
            type: 'time',
            time: {
              unit: 'second'
            },
            barPercentage: 1.6,
            gridLines: {
              drawBorder: false,
              color: 'rgba(255,255,255,0.1)',
              zeroLineColor: 'transparent',
              display: false,
            },
            ticks: {
              padding: 20,
              fontColor: '#9f9f9f'
            }
          }]
        },
      }
    });


  }

  openDialog(): void {
    this.openSnackBar('not implemented yet', 'ok');
  }

  launchDiscovery(): void {
    this.updatingState = true
    this.prestoApi.launchDiscover().subscribe(
      result => {
        this.openSnackBar('cluster discovery completed', 'ok')
      },
      err => {

        this.openSnackBar('cluster discovery error', 'ok')
      },
      () => {
        this.updatingState = false;
        this.forceRefresh.next();
      }
    )
  }


  openSnackBar(message: string, action: string) {
    this._snackBar.open(message, action, {
      duration: 2000,
    });
  }

  toggleClusterStatus(event: MatSlideToggleChange, cluster: Cluster) {
    this.updatingState = true;
    this.prestoApi.updateClusterStatus({
      name: cluster.name,
      enabled: event.checked,
    }).subscribe(
      result => {
      },
      err => this.openSnackBar('error', 'ok'),
      () => {
        this.updatingState = false;
        this.forceRefresh.next();
      }
    )
  }
}


