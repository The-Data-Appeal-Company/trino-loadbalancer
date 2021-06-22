import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {environment} from '../../environments/environment';


@Injectable({
  providedIn: 'root'
})
export class TrinoApiService {

  constructor(private readonly client: HttpClient) {
  }

  private endpoint(): string {
    if (environment.production) {
      return window.location.protocol + '//' + window.location.host
    }

    return 'http://localhost:8998'
  }

  public clusters(): Observable<Cluster[]> {
    return this.client.get<Cluster[]>(this.endpoint() + '/api/clusters')
  }

  public stats(): Observable<Stats> {
    return this.client.get<Stats>(this.endpoint() + '/api/stats')
  }

  public updateClusterStatus(req: ClusterUpdateRequest): Observable<any> {
    return this.client.patch(this.endpoint() + '/api/cluster/' + req.name, req)
  }

  public launchDiscover(): Observable<any> {
    return this.client.post(this.endpoint() + '/api/cluster/discover', null)
  }
}

export interface ClusterUpdateRequest {
  name: string;
  enabled: boolean;
}

export interface Stats {
  total_workers: number;
  running_queries: number;
  blocked_queries: number;
  queued_queries: number;
}

export interface Cluster {
  name: string;
  host: string;
  available: boolean;
  enabled: boolean;
  tags: Map<string, string>;
}
