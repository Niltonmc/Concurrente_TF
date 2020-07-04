import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Informe } from '../model/informe';
import { Grupo } from '../model/grupo';

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  private BASE_URL = "http://localhost:8000";
  private KNN_URL = `${this.BASE_URL}/clustering_covid`;
  private SAVE_UPDATE_INFORMES = `${this.BASE_URL}/classification_covid`;
  private Result = `${this.BASE_URL}/classification_covid`
  private RESULT2 = `${this.BASE_URL}/clustering_covid`;

  constructor(private http: HttpClient) { }


  classificationCovid(informe: Informe): Observable<Informe>{
    return this.http.post<Informe>(this.SAVE_UPDATE_INFORMES,informe);
  }

  clusteringCovid(grupo: Grupo):Observable<Grupo>{
    return this.http.post<Grupo>(this.KNN_URL,grupo);
  }

  result1():Observable<any>{
    return this.http.get<any>(this.Result);
  }

  result2():Observable<any>{
    return this.http.get<any>(this.RESULT2);
  }
}
