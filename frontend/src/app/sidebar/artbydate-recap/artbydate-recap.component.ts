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
  constructor(private aService: ArticleService) { }

  ngOnInit() {
    const artRecap = this.aService.getArticlesRecap();
    this.buildList(artRecap);
  }

  private getMonthName(n: number): string {
    const monthName = [
      'Jenuary',   'February', 'March',    'April',
      'May',       'June',     'July',     'August',
      'September', 'October',  'November', 'December'
    ];
    return monthName[n];
  }
  /* Reorganize given articles list*/
  private buildList(aR: ArticleRecap[]) {

    aR.forEach(el => {
      // TODO
    });
  }

}
