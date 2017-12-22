import { Component, OnInit } from '@angular/core';
import { Article } from '../article.model';

@Component({
  selector: 'app-article-list',
  templateUrl: './article-list.component.html',
  styleUrls: ['./article-list.component.css']
})
export class ArticleListComponent implements OnInit {
  articles: Article[];
  constructor() {}

  ngOnInit() {
    const articles = [
      new Article(
        1,
        'Titolo',
        'echoes',
        'This is an article test, I hope it works.',
        new Date(2017, 12, 22, 17, 23)
      ),
      new Article(
        2,
        'Articolo 2',
        'echoes',
        `Lorem ipsum dolor sit amet consectetur adipisicing elit.
        Qui quasi eum eveniet perspiciatis repellat minus tempora in
        reprehenderit porro culpa! Reprehenderit dolorem esse ullam, saepe atque quia
        distinctio recusandae maxime possimus soluta!
        Aliquid exercitationem veritatis, dolores culpa ratione sint expedita. `,
        new Date(2017, 12, 22, 17, 21)
      )
    ];
    this.articles = articles;
    this.articles.sort( function(a, b) { // order from the most recent
      const d1 = a.date;
      const d2 = b.date;
      return d1 > d2 ? -1 : d1 < d2 ? +1 : 0;
    });
  }
}
