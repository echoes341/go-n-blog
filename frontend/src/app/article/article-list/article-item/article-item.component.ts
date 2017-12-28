import { Component, OnInit, Input } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-article-item',
  templateUrl: './article-item.component.html',
  styleUrls: ['./article-item.component.css']
})
export class ArticleItemComponent implements OnInit {

  @Input() id: number;
  @Input() title: string;
  @Input() author: string;
  @Input() text: string;
  @Input() date: string;
  @Input() commentNum: number;
  @Input() likeNum: number;
  @Input() isLiked: boolean;
  @Input() isSliced: boolean;

  constructor(private router: Router) { }

  ngOnInit() {

  }

  public toggleLike() {
    this.isLiked ? this.likeNum-- : this.likeNum++;
    this.isLiked = !this.isLiked;
  }

  public onGoToComments() {
    this.router.navigate(['/article', 'v', this.id], { fragment: 'comments'});
  }
}
