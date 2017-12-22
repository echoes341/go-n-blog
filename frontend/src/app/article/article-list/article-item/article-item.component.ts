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
  dateFormat: string;
  constructor(private commentServ: CommentService ) { }

  ngOnInit() {
    this.cCount = this.commentServ.getCountCommentByArtID(this.article.id);
    const d = this.article.date;
    let dF = d.getHours() + ':' + d.getMinutes() + ' ';
    const day = d.getDate() < 10 ? '0' + d.getDate() : d.getDate();
    const month: string = d.getMonth() + 1 < 10 ? '0' + (d.getMonth() + 1) : (d.getMonth() + 1) + '';
    dF += day + '-' + month + '-' + d.getFullYear();
    this.dateFormat = dF;

  }

}
