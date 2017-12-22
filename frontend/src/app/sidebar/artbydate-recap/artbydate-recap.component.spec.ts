import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ArtbydateRecapComponent } from './artbydate-recap.component';

describe('ArtbydateRecapComponent', () => {
  let component: ArtbydateRecapComponent;
  let fixture: ComponentFixture<ArtbydateRecapComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ArtbydateRecapComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ArtbydateRecapComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
