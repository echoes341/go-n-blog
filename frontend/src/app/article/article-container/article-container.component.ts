import { Component, OnInit } from '@angular/core';
import { LikeService } from '../like.service';

@Component({
  selector: 'app-article-container',
  templateUrl: './article-container.component.html',
  styleUrls: ['./article-container.component.css'],
  providers: [LikeService]
})
export class ArticleContainerComponent implements OnInit {

  constructor() { }

  ngOnInit() {
  }

}
