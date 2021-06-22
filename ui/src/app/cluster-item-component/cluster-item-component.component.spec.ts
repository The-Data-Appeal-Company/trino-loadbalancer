import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ClusterItemComponentComponent } from './cluster-item-component.component';

describe('ClusterItemComponentComponent', () => {
  let component: ClusterItemComponentComponent;
  let fixture: ComponentFixture<ClusterItemComponentComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ClusterItemComponentComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ClusterItemComponentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
