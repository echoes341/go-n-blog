import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { Comment } from '../comment.model';
import { NgForm } from '@angular/forms';

@Component({
  selector: 'app-add-comment',
  templateUrl: './add-comment.component.html',
  styleUrls: ['./add-comment.component.css']
})
export class AddCommentComponent implements OnInit {
  @Output() addComment = new EventEmitter<Comment>();
  constructor() { }

  onAddComment(f: NgForm) {
    const name = f.value.name;
    const email = f.value.email;
    const content = f.value.text;
    const c = new Comment(-1, -1, name, email, content);
    this.addComment.emit(c);
    console.log(f);
  }
  ngOnInit() {
  }

}
