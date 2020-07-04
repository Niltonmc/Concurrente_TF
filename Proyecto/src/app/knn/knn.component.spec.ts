import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KnnComponent } from './knn.component';

describe('KnnComponent', () => {
  let component: KnnComponent;
  let fixture: ComponentFixture<KnnComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KnnComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KnnComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
