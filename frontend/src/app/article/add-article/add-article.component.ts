import { Component, OnInit } from '@angular/core';
import * as Quill from 'quill';
import { Article } from '../article.model';
import { ArticleService } from '../article.service';
import { Router } from '@angular/router';
import { CanComponentDeactivate } from '../../can-deactivated-guard.service';
import { Observable } from 'rxjs/Observable';

@Component({
  selector: 'app-add-article',
  templateUrl: './add-article.component.html',
  styleUrls: ['./add-article.component.css']
})
export class AddArticleComponent implements OnInit, CanComponentDeactivate {
  title: string;
  author: string;
  quill: Quill;
  constructor(private aServ: ArticleService, private route: Router) {}
  onAddArticle() {
    const a = new Article(
      -1,
      this.title,
      this.author,
      // Picking HTML from editor as suggested in official quill's github.
      // Even if it's not the safest idea ever, actually.
      this.quill.root.innerHTML,
      new Date()
    );
    const id = this.aServ.addArticle(a);
    this.route.navigate(['/article', 'v', id]);
  }

  canDeactivate(): Observable<boolean> | Promise<boolean> | boolean {

    if (
      // !!str ->Boolean(str) -> !isEmpty? false: empty | true: notEmpty
      !!this.title  || this.quill.getText().length !== 1
      ) {
      return confirm('If you change the page the current content will be lost. OK?');
    } else {
      return true;
    }
  }

  ngOnInit() {
    this.quill = new Quill('#editor', {
      theme: 'snow'
    });
  }
}
