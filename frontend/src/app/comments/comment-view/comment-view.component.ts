import { Component, OnInit, Input } from '@angular/core';
import { Comment } from '../comment.model';

@Component({
  selector: 'app-comment-view',
  templateUrl: './comment-view.component.html',
  styleUrls: ['./comment-view.component.css']
})
export class CommentViewComponent implements OnInit {
  @Input() comment: Comment;

  onEditComment() {

  }
  onRemoveComment() {

  }
  constructor() { }

  ngOnInit() {
  }

}
