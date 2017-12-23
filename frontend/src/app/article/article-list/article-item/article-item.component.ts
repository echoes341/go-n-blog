import { Component, OnInit, Input } from '@angular/core';

import { Article } from '../../article.model';



import { CommentService } from '../../../comments/comment.service';
import { LikeService } from '../../like.service';

@Component({
  selector: 'app-article-item',
  templateUrl: './article-item.component.html',
  styleUrls: ['./article-item.component.css']
})
export class ArticleItemComponent implements OnInit {

  @Input() article: Article;
  cCount: number;
  dateFormat: string;
  likeNum: number;
  isLiked: boolean;

  constructor(private commentServ: CommentService, private likeServ: LikeService ) { }

  ngOnInit() {
    this.cCount = this.commentServ.getCountCommentByArtID(this.article.id);
    const d = this.article.date;
    let dF = d.getHours() + ':' + d.getMinutes() + ' ';
    const day = d.getDate() < 10 ? '0' + d.getDate() : d.getDate();
    const month: string = d.getMonth() + 1 < 10 ? '0' + (d.getMonth() + 1) : (d.getMonth() + 1) + '';
    dF += day + '-' + month + '-' + d.getFullYear();
    this.dateFormat = dF;

    this.likeNum = this.getLikeNum();
    this.isLiked = false;

  }

  public getLikeNum(): number {
    const likes = this.likeServ.getLikeByArtID(this.article.id);
    return likes.length;
  }

  public toggleLike() {
    this.isLiked ? this.likeNum-- : this.likeNum++;
    this.isLiked = !this.isLiked;
  }

}
