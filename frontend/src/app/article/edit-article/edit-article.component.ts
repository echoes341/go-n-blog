import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import * as Quill from 'quill';

import { ArticleService } from '../article.service';
import { ActivatedRoute, Router } from '@angular/router';
import { Article } from '../article.model';
import { CanComponentDeactivate } from '../../can-deactivated-guard.service';

@Component({
  selector: 'app-edit-article',
  templateUrl: './edit-article.component.html',
  styleUrls: ['./edit-article.component.css']
})
export class EditArticleComponent implements OnInit, CanComponentDeactivate {
  article: Article;
  @ViewChild('text') text: ElementRef;
  title = '';
  isChanged = false;
  quill: Quill;
  id: number;

  constructor(private aServ: ArticleService,
    private route: ActivatedRoute,
    private router: Router
  ) {}

  ngOnInit() {
    this.id = +this.route.snapshot.params['id'];
    this.article = this.aServ.getArticleByID(this.id);
    this.title = this.article.title;
    this.text.nativeElement.innerHTML = this.article.text;

    this.quill = new Quill('#editor', {
      theme: 'snow'
    });
    this.quill.once('text-change', () => {
      this.isChanged = true;
    });
  }

  onEditArticle() {
    if ( this.title !== this.article.title || this.isChanged ) {
      this.article.title = this.title;
      this.article.text = this.quill.root.innerHTML;
      this.aServ.editArticle(this.article);
      this.router.navigate(['/article', 'v', this.article.id]);
    } else {
      alert('Nothing changed!');
    }
  }

  canDeactivate() {
    return true;
  }
}
