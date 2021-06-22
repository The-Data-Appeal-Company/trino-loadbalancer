import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-stats-card',
  templateUrl: './stats-card.component.html',
  styleUrls: ['./stats-card.component.scss']
})
export class StatsCardComponent implements OnInit {

  @Input() value: number = 0;

  @Input() unit: string = '';

  @Input() icon: string = '';

  constructor() {
  }

  ngOnInit(): void {
  }

}
