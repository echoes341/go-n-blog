import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { Comment } from '../comment.model';

@Component({
  selector: 'app-add-comment',
  templateUrl: './add-comment.component.html',
  styleUrls: ['./add-comment.component.css']
})
export class AddCommentComponent implements OnInit {
  @Output() addComment = new EventEmitter<Comment>();
  constructor() { }

  onAddComment(name: string, email: string, content: string) {
    const c = new Comment(-1, -1, name, email, content);
    this.addComment.emit(c);
  }
  ngOnInit() {
  }

}
