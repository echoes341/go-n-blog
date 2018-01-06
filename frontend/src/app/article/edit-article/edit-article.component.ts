import { Component, OnInit, ViewChild, ElementRef, OnDestroy } from '@angular/core';

import { ArticleService } from '../article.service';
import { ActivatedRoute, Router, Params } from '@angular/router';
import { Article } from '../article.model';
import { CanComponentDeactivate } from '../../can-deactivated-guard.service';
import { race } from 'q';
import { Subscription } from 'rxjs/Subscription';

@Component({
  selector: 'app-edit-article',
  templateUrl: './edit-article.component.html',
  styleUrls: ['./edit-article.component.css']
})
export class EditArticleComponent implements OnInit, CanComponentDeactivate, OnDestroy {
  article: Article;
  title = '';
  text = '';
  author = 'echoes';
  editMode: boolean;
  isAdded = false;
  id: number;

  subscription: Subscription;


  constructor(
    private aServ: ArticleService,
    private route: ActivatedRoute,
    private router: Router
  ) {}

  ngOnInit() {
    this.subscription = this.route.params
      .subscribe((params: Params) => { // reading of id
        this.editMode = params['id'] != null;
        if (this.editMode) {
          this.id = +params['id'];
          this.article = this.aServ.getArticleByID(this.id); // id present-> edit mode
          this.title = this.article.title;
          this.text = this.article.text;
          this.author = this.article.author;
        }
    });
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  isChanged() {
    if (this.editMode) {
      return (
        this.title !== this.article.title || this.text !== this.article.text
      );
    }
    return !this.isAdded && (this.title || this.text);
  }
  onSubmit() {
    if (this.editMode) {
      if (this.isChanged()) {
        this.article.title = this.title;
        this.article.text = this.text;
        this.aServ.editArticle(this.article);
        this.router.navigate(['/', this.article.id]);
      } else {
        alert('Nothing changed!');
      }
    } else {
      // new mode
      const a = new Article(-1, this.title, this.author, this.text, new Date());
      const id = this.aServ.addArticle(a);
      this.isAdded = true;
      this.router.navigate(['/', id]);
    }
  }

  canDeactivate() {
    if (this.isChanged()) {
      return confirm(
        'If you change the page the current content will be lost. OK?'
      );
    } else {
      return true;
    }
  }
}
