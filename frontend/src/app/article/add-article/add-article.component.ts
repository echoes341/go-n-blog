import { Component, OnInit } from '@angular/core';
import * as Quill from 'quill';
import { Article } from '../article.model';
import { ArticleService } from '../article.service';

@Component({
  selector: 'app-add-article',
  templateUrl: './add-article.component.html',
  styleUrls: ['./add-article.component.css']
})
export class AddArticleComponent implements OnInit {

  quill: Quill;
  constructor(private aServ: ArticleService) { }
  onAddArticle(title: string, author: string) {
    const a = new Article(-1,
      title,
      author,
      this.quill.root.innerHTML,
      new Date());
      this.aServ.addArticle(a);
    console.log(this.quill.root.innerHTML);
  }

  ngOnInit() {

    this.quill = new Quill('#editor', {
      theme: 'snow'
    });
  }


}
