import { Component, OnInit } from '@angular/core';
import { Article } from '../article.model';
import { CommentService } from '../../comments/comment.service';
import { LikeService } from '../like.service';
import { ArticleService } from '../article.service';
import { ActivatedRoute } from '@angular/router';
import { Comment } from '../../comments/comment.model';

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
  comments: Comment[];

  constructor(private cmmtServ: CommentService,
    private likeServ: LikeService,
    private aServ: ArticleService,
    private route: ActivatedRoute ) { }

  ngOnInit() {
    const id = +this.route.snapshot.params['id']; // get article ID

    this.aServ.getArticle(id).subscribe(
      (data: any) => {
        this.article = data.data;
        /* date formatting */
        const d = this.article.date;
        this.dateFormat = d.getHours() + ':' + d.getMinutes() + ' ';
        const day = d.getDate() < 10 ? '0' + d.getDate() : d.getDate();
        const month: string = d.getMonth() + 1 < 10 ? '0' + (d.getMonth() + 1) : (d.getMonth() + 1) + '';
        this.dateFormat += day + '-' + month + '-' + d.getFullYear();
      },
      (error: Response) => {
        if (error.status === 404) {
          console.log('Article not found');
        }
      }
    );


    this.likeNum = this.getLikeNum();
    this.isLiked = this.likeServ.isLiked(id /*, userid*/);

    // retrieving comments from article id
    this.comments = this.cmmtServ.getCommentByArtID(this.article.id);

  }

  public getLikeNum(): number {
    const likes = this.likeServ.getLikeByArtID(this.article.id);
    return likes.length;
  }

  public toggleLike() {
    this.isLiked ? this.likeNum-- : this.likeNum++;
    this.isLiked = !this.isLiked;
  }

  public addComment(c: Comment) {
    c.idA = this.article.id;
    this.cmmtServ.addComment(c);
    this.comments = this.cmmtServ.getCommentByArtID(this.article.id);
  }

  public removeComment(id: number) {
    this.cmmtServ.removeComment(id);
    this.comments = this.cmmtServ.getCommentByArtID(this.article.id);
  }
}
