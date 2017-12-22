import { Component, OnInit, Input } from '@angular/core';
import { Article } from '../../article.model';
import { CommentService } from '../../../comments/comment.service';


@Component({
  selector: 'app-article-item',
  templateUrl: './article-item.component.html',
  styleUrls: ['./article-item.component.css']
})
export class ArticleItemComponent implements OnInit {

  @Input() article: Article;
  cCount: number;
  constructor(private commentServ: CommentService ) { }

  ngOnInit() {
    this.cCount = this.commentServ.getCountCommentByArtID(this.article.id);
  }

}
