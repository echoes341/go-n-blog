import { Component, OnInit } from '@angular/core';
import { Article } from '../article.model';
import { Like } from '../like.model';

import { ArticleService } from '../article.service';
import { LikeService } from '../like.service';
import { CommentService } from '../../comments/comment.service';

@Component({
  selector: 'app-article-list',
  templateUrl: './article-list.component.html',
  styleUrls: ['./article-list.component.css'],
})
export class ArticleListComponent implements OnInit {
  articles: Article[];
  articlesList: {
    id: number,
    title: string,
    author: string,
    date: string,
    text: string,
    cCount: number,
    lCount: number,
    isLiked: boolean,
    isSliced: boolean
  }[] = [];

  constructor(
    private aService: ArticleService,
    private likeServ: LikeService,
    private commServ: CommentService
  ) {}

  ngOnInit() {
    this.articles = this.aService.getFirstsXFromDate(5, new Date());

    // Date format
    // strip articles whose text is more than 300 chars
    // count how many likes and if is liked ([TODO] by current user)
    // count how many comments

    // should be done with a single call to API?

    this.articles.forEach( el => {
      // Date format
      const d = el.date;
      const time: string = d.getHours() + ':' + d.getMinutes();
      const day: string = d.getDate() < 10 ? '0' + d.getDate() : d.getDate() + '';
      const month: string = d.getMonth() + 1 < 10 ? '0' + (d.getMonth() + 1) : (d.getMonth() + 1) + '';
      const dateFormat = time + ' ' + day + '-' + month + '-' + d.getFullYear();

      // Strip text
      const isSliced: boolean = el.text.length > 300;
      const t: string = isSliced ? el.text.slice(0, 300) : el.text;

      // 'pack' and build data
      this.articlesList.push( {
        id: el.id,
        title: el.title,
        author: el.author,
        date: dateFormat,
        text: t,
        // count comments
        cCount: this.commServ.getCountCommentByArtID(el.id),
        // count likes
        lCount: this.likeServ.getLikeByArtID(el.id).length,
        isLiked: this.likeServ.isLiked(el.id),
        isSliced: isSliced
      });

    });



  }
}
