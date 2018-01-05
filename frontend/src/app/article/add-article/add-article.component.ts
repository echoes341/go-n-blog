import { Component, OnInit } from '@angular/core';
import { Article } from '../article.model';
import { ArticleService } from '../article.service';
import { Router } from '@angular/router';
import { CanComponentDeactivate } from '../../can-deactivated-guard.service';
import { Observable } from 'rxjs/Observable';
import { NgForm } from '@angular/forms';

@Component({
  selector: 'app-add-article',
  templateUrl: './add-article.component.html',
  styleUrls: ['./add-article.component.css']
})
export class AddArticleComponent implements OnInit, CanComponentDeactivate {
  title = '';
  author = 'echoes';
  text = '';
  isAdded = false;
  constructor(private aServ: ArticleService, private route: Router) {}
  onAddArticle() {
    const a = new Article(
      -1,
      this.title,
      this.author,
      this.text,
      new Date()
    );
    const id = this.aServ.addArticle(a);
    this.isAdded = true;
    this.route.navigate(['/', id]);
  }

  canDeactivate(): Observable<boolean> | Promise<boolean> | boolean {

    console.log(Boolean(this.text));
    if (
      /*
        Here || force boolean conversion of strings
        It works like a isNotEmpty() func:
        string  empty  -> false
      */
      !this.isAdded && (this.title  || this.text)
      ) {
      return confirm('If you change the page the current content will be lost. OK?');
    } else {
      return true;
    }
  }

  ngOnInit() {
  }
}
