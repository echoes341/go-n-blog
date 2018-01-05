import { Component, OnInit } from '@angular/core';
import { ArticleService } from '../../article/article.service';
import { ArticleRecap } from '../../article/article.model';

@Component({
  selector: 'app-artbydate-recap',
  templateUrl: './artbydate-recap.component.html',
  styleUrls: ['./artbydate-recap.component.css']
})
export class ArtbydateRecapComponent implements OnInit {


  list = {};
  objectKeys= Object.keys;

  constructor(private aService: ArticleService) { }

  ngOnInit() {
    const artRecap = this.aService.getArticlesRecap();
    this.buildList(artRecap);
  }

  public getMonthName(n: number): string {
    const monthName = [
      'January',   'February', 'March',    'April',
      'May',       'June',     'July',     'August',
      'September', 'October',  'November', 'December'
    ];
    return monthName[n];
  }
  /* Reorganize given articles list*/
  private buildList(aR: ArticleRecap[]) {

    aR.forEach(el => {
      if (this.list[el.year] === undefined) {
        this.list[el.year] = {};
      }
      this.list[el.year][el.month] = el.c;

    });
  }

}
