import { Component, OnInit } from '@angular/core';
import { Article } from '../article.model';
import { CommentService } from '../../comments/comment.service';
import { LikeService } from '../like.service';
import { ArticleService } from '../article.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-view-article',
  templateUrl: './view-article.component.html',
  styleUrls: ['./view-article.component.css']
})
export class ViewArticleComponent implements OnInit {

  article: Article;
  dateFormat: string;
  likeNum: number;
  isLiked: boolean;

  constructor(private commentServ: CommentService,
    private likeServ: LikeService,
    private aServ: ArticleService,
    private route: ActivatedRoute ) { }

  ngOnInit() {
    const id = this.route.snapshot.params['id']; // get article ID
    console.log(this.aServ.getArticleByID(id));
    this.article = this.aServ.getArticleByID(id); // get article

    /* date formatting */
    const d = this.article.date;
    this.dateFormat = d.getHours() + ':' + d.getMinutes() + ' ';
    const day = d.getDate() < 10 ? '0' + d.getDate() : d.getDate();
    const month: string = d.getMonth() + 1 < 10 ? '0' + (d.getMonth() + 1) : (d.getMonth() + 1) + '';
    this.dateFormat += day + '-' + month + '-' + d.getFullYear();

    this.likeNum = this.getLikeNum();
    this.isLiked = this.likeServ.isLiked(id /*, userid*/);

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
