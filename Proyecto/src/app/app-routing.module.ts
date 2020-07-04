import { NgModule, Component } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { KnnComponent } from './knn/knn.component';


const routes: Routes = [
  {
    path:'home',
    component: HomeComponent
  },
  {
    path: 'clustering_covid',
    component: KnnComponent
  } ,
  {
    path:'',
    component: HomeComponent,
    pathMatch:'full'
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
