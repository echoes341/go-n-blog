import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';


import { ArticleService } from '../article.service';
import { ActivatedRoute, Router } from '@angular/router';
import { Article } from '../article.model';
import { CanComponentDeactivate } from '../../can-deactivated-guard.service';
import { race } from 'q';

@Component({
  selector: 'app-edit-article',
  templateUrl: './edit-article.component.html',
  styleUrls: ['./edit-article.component.css']
})
export class EditArticleComponent implements OnInit, CanComponentDeactivate {
  article: Article;
  title = '';
  text = '';

  id: number;

  constructor(private aServ: ArticleService,
    private route: ActivatedRoute,
    private router: Router
  ) {}

  ngOnInit() {
    this.id = +this.route.snapshot.params['id'];
    this.article = this.aServ.getArticleByID(this.id);
    this.title = this.article.title;
    this.text = this.article.text;
  }

  isChanged() {
   return this.title !== this.article.title || this.text !== this.article.text;
  }
  onEditArticle() {
    if ( this.isChanged() ) {
      this.article.title = this.title;
      this.article.text = this.text;
      this.aServ.editArticle(this.article);
      this.router.navigate(['/article', 'v', this.article.id]);
    } else {
      alert('Nothing changed!');
    }
  }

  canDeactivate() {
    if (this.isChanged()) {
      return confirm('If you change the page the current content will be lost. OK?');
    } else {
      return true;
    }
  }
}
