import { ComponentFixture, TestBed } from '@angular/core/testing';

import { StatsOverviewComponentComponent } from './stats-overview-component.component';

describe('StatsOverviewComponentComponent', () => {
  let component: StatsOverviewComponentComponent;
  let fixture: ComponentFixture<StatsOverviewComponentComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ StatsOverviewComponentComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(StatsOverviewComponentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
