import { Component, OnInit } from '@angular/core';
import * as Quill from 'quill';
import { Article } from '../article.model';
import { ArticleService } from '../article.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-add-article',
  templateUrl: './add-article.component.html',
  styleUrls: ['./add-article.component.css']
})
export class AddArticleComponent implements OnInit {

  quill: Quill;
  constructor(
    private aServ: ArticleService,
    private route: Router
  ) { }
  onAddArticle(title: string, author: string) {
    const a = new Article(
      -1,
      title,
      author,
      // Picking HTML from editor as suggested in official quill's github.
      // Even if it's not the safest idea ever, actually.
      this.quill.root.innerHTML,
      new Date()
    );
    const id = this.aServ.addArticle(a);
    this.route.navigate(['/article', 'v', id]);
  }

  ngOnInit() {

    this.quill = new Quill('#editor', {
      theme: 'snow'
    });
  }


}
