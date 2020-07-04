import { Component, OnInit } from '@angular/core';
import { ApiService } from '../shared/api.service';
import { Grupo } from '../model/grupo';

@Component({
  selector: 'app-knn',
  templateUrl: './knn.component.html',
  styleUrls: ['./knn.component.css']
})
export class KnnComponent implements OnInit {

  newForm: Grupo = new Grupo();
  resul:any;
  visible: boolean = false;
  centroids: any;
  quantity:any;
  genero:string [] = [];
  cancer: string [] = [];
  hypertencion: string [] = [];
  respiratorio: string [] =[];
  cardiovascular: string[]=[];
  diabetis: string [] = [];

  constructor(private apiService:ApiService) { }

  ngOnInit(): void {
  }

  createGroup(){
    this.apiService.clusteringCovid(this.newForm).subscribe(
      res =>{
        this.apiService.result2().subscribe(
          res=>{
            this.resul = res;
            this.centroids = this.resul.centroidsInforme;
            this.quantity =this.resul.centroidsQuantity;
            var n = this.centroids.length
            for(var i=0; i<n;i++){
              this.centroids[i].age_group = Math.trunc(this.centroids[i].age_group)
            }
            for(var i=0; i<n;i++){
              if(this.centroids[i].sex==0){
                this.genero.push("Masculino");
              } else{
                this.genero.push("Femenino");
              }
            }
            for(var i=0; i<n;i++){
              if(this.centroids[i].cardiovascular_disease>0.5){
                this.cardiovascular.push("Si");
              }else{
                this.cardiovascular.push("No");
              }
            }
            for(var i=0; i<n;i++){
              if(this.centroids[i].diabetes>0.5){
                this.diabetis.push("Si");
              }else{
                this.diabetis.push("No");
              }
            }
            for(var i=0; i<n;i++){
              if(this.centroids[i].respiratory_disease>0.5){
                this.respiratorio.push("Si");
              }else{
                this.respiratorio.push("No");
              }
            }
            for(var i=0; i<n;i++){
              if(this.centroids[i].hypertension>0.5){
                this.hypertencion.push("Si");
              }else{
                this.hypertencion.push("No");
              }
            }
            for(var i=0; i<n;i++){
              if(this.centroids[i].cancer>0.5){
                this.cancer.push("Si");
              }else{
                this.cancer.push("No");
              }
            }
            
            this.visible=true;
          },
          err => { alert("An error has occurred while write the Kmeans"); }
        )
      },
      err => { alert("An error has occurred while saving the Group"); }
    )
  }

}
