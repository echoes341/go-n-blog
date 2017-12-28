import { Component, OnInit, Input } from '@angular/core';

import { Article } from '../../article.model';



import { CommentService } from '../../../comments/comment.service';
import { LikeService } from '../../like.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-article-item',
  templateUrl: './article-item.component.html',
  styleUrls: ['./article-item.component.css']
})
export class ArticleItemComponent implements OnInit {

  // tslint:disable-next-line:no-input-rename
  @Input() article: Article;
  strippedText: string;
  cCount: number;
  dateFormat: string;
  likeNum: number;
  isLiked: boolean;

  constructor(private commentServ: CommentService,
    private likeServ: LikeService,
    private router: Router
  ) { }

  ngOnInit() {
    this.cCount = this.commentServ.getCountCommentByArtID(this.article.id);

    /* date formatting */
    const d = this.article.date;
    this.dateFormat = d.getHours() + ':' + d.getMinutes() + ' ';
    const day = d.getDate() < 10 ? '0' + d.getDate() : d.getDate();
    const month: string = d.getMonth() + 1 < 10 ? '0' + (d.getMonth() + 1) : (d.getMonth() + 1) + '';
    this.dateFormat += day + '-' + month + '-' + d.getFullYear();

    this.strippedText = this.article.text.slice(0, 300);

    this.likeNum = this.getLikeNum();
    this.isLiked = this.likeServ.isLiked(this.article.id /*, userid*/);

  }

  public getLikeNum(): number {
    const likes = this.likeServ.getLikeByArtID(this.article.id);
    return likes.length;
  }

  public toggleLike() {
    this.isLiked ? this.likeNum-- : this.likeNum++;
    this.isLiked = !this.isLiked;
  }

  public onGoToComments() {
    this.router.navigate(['/article', 'v', this.article.id], { fragment: 'comments'});
  }
}
