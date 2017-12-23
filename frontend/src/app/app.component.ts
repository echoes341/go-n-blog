import { Component } from '@angular/core';
import { CommentService } from './comments/comment.service';
import { ArticleService } from './article/article.service';
import { LikeService } from './article/like.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
  providers: [CommentService, ArticleService]
})
export class AppComponent {
  title = 'Echoes Blog';
}
