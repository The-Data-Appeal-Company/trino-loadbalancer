import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ClusterListComponentComponent } from './cluster-list-component.component';

describe('ClusterListComponentComponent', () => {
  let component: ClusterListComponentComponent;
  let fixture: ComponentFixture<ClusterListComponentComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ClusterListComponentComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ClusterListComponentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
