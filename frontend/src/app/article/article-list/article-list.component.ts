import { Component, OnInit } from '@angular/core';
import { Article } from '../article.model';
import { ArticleService } from '../article.service';

@Component({
  selector: 'app-article-list',
  templateUrl: './article-list.component.html',
  styleUrls: ['./article-list.component.css']
})
export class ArticleListComponent implements OnInit {
  articles: Article[];
  constructor(private aService: ArticleService) {}

  ngOnInit() {
    this.articles = this.aService.getFirstsXFromDate(5, new Date());
  }
}
