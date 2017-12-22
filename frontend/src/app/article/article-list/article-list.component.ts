import { Component, OnInit } from '@angular/core';
import { Article } from '../shared/article.model';

@Component({
  selector: 'app-article-list',
  templateUrl: './article-list.component.html',
  styleUrls: ['./article-list.component.css']
})
export class ArticleListComponent implements OnInit {
  articles = [
    new Article(
      1,
      'Titolo',
      'echoes',
      'This is an article test, I hope it works.',
      new Date('2017-12-22')
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
      new Date('2017-12-22')
    )
  ];
  constructor() {}

  ngOnInit() {}
}
