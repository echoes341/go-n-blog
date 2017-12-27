import { Component, OnInit } from '@angular/core';
import * as Quill from 'quill';

@Component({
  selector: 'app-add-article',
  templateUrl: './add-article.component.html',
  styleUrls: ['./add-article.component.css']
})
export class AddArticleComponent implements OnInit {

  constructor() { }
  onAddArticle(content: any) {
    console.log(content);
  }

  ngOnInit() {

    const quill = new Quill('#editor', {
      theme: 'snow'
    });
  }

}
