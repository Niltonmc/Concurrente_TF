import { Component, OnInit } from '@angular/core';
import { Informe } from '../model/informe';
import { ApiService } from '../shared/api.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {

  selectedIndex: number;
  selectedSex:number;
  visible: boolean = false;
  estado:number = 0; 
  resultado:any;

  newForm: Informe = new Informe();


  ages = [
    {
      value:"1 - 19"
    },{
      value:"20 - 39"
    },
    {
      value:"40 - 59"
    },
    {
      value:"60 - 79"
    },
    {
      value:"80 - 100+"
    }
  ]

  sex = [
    {
      value: "Male"
    },
    {
      value: "Female"
    }
  ]

  conditions = [
    {
      name: "Cardiovascular disease",
      option: "No",
      value : 0,
      value2: false
    },
    {
      name: "Diabetes",
      option: "No",
      value: 0,
      value2: false
    },{
      name: "Chronic respiratory disease",
      option: "No",
      value: 0,
      value2: false
      
    },
    {
      name: "Hypertension",
      option: "No",
      value: 0,
      value2: false

    },
    {
      name: "Cancer",
      option: "No",
      value: 0,
      value2: false

    }
  ]

  constructor(private apiService: ApiService) { }

  ngOnInit(): void {
  }

  createForm(){
    this.apiService.classificationCovid(this.newForm).subscribe(

      res => {
        this.newForm.age_group = this.selectedIndex;
        if(this.selectedSex==1){
          this.newForm.sex = 1;
        }else{
          this.newForm.sex = 0;
        }
        //location.reload();
        this.apiService.result1().subscribe(
          res =>{
            this.resultado=res;
            this.estado=this.resultado.clase;
            console.log(res.clase);
            this.visible = true;
          },
          err => { alert("An error has occurred while saving the Covid Analysis"); }
  
        )
      },
      err => { alert("An error has occurred while saving the Informe"); }
    )
  }
  public setRow(_index:number){
    this.selectedIndex = _index;
    this.newForm.age_group = _index;
  }

  public setSex(_index:number){
    this.selectedSex = _index;
    if(this.selectedSex==1){
      this.newForm.sex = 1;
    }else{
      this.newForm.sex = 0;
    }
  }

  public haveConditions(_index:number){
    this.conditions[_index].value2 = !this.conditions[_index].value2;
    if(this.conditions[_index].option=="No"){
      this.conditions[_index].option="Yes"
      if(this.conditions[_index].name=="Cardiovascular disease"){
        this.newForm.cardiovascular_disease = 1;
      }
      if(this.conditions[_index].name=="Diabetes"){
        this.newForm.diabetes = 1;
      }
      if(this.conditions[_index].name=="Chronic respiratory disease"){
        this.newForm.respiratory_disease = 1;
      }
      if(this.conditions[_index].name=="Hypertension"){
        this.newForm.hypertension = 1;
      }
      if(this.conditions[_index].name=="Cancer"){
        this.newForm.cancer = 1;
      }
    }else{
      this.conditions[_index].option="No"
      if(this.conditions[_index].name=="Cardiovascular disease"){
        this.newForm.cardiovascular_disease = 0;
      }
      if(this.conditions[_index].name=="Diabetes"){
        this.newForm.diabetes = 0;
      }
      if(this.conditions[_index].name=="Chronic respiratory disease"){
        this.newForm.respiratory_disease = 0;
      }
      if(this.conditions[_index].name=="Hypertension"){
        this.newForm.hypertension = 0;
      }
      if(this.conditions[_index].name=="Cancer"){
        this.newForm.cancer = 0;
      }
    }
  }
}
