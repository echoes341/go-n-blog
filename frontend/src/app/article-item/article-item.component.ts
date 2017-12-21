import { Component, OnInit } from '@angular/core';
import { Article } from './article.model'

@Component({
  selector: 'app-article-item',
  templateUrl: './article-item.component.html',
  styleUrls: ['./article-item.component.css']
})
export class ArticleItemComponent implements OnInit {
  art: Article;
  constructor() { }

  ngOnInit() {
  }

}
