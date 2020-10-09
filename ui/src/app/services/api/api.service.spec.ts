import { TestBed } from '@angular/core/testing';

import { PrestoApiService } from './prestoApiService';

describe('ApiService', () => {
  let service: PrestoApiService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(PrestoApiService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
