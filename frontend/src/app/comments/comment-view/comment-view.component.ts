import { Component, OnInit, Input, Output, EventEmitter  } from '@angular/core';
import { Comment } from '../comment.model';
import { ElementRef } from '@angular/core/src/linker/element_ref';



@Component({
  selector: 'app-comment-view',
  templateUrl: './comment-view.component.html',
  styleUrls: ['./comment-view.component.css']
})
export class CommentViewComponent implements OnInit {
  @Input() comment: Comment;
  @Output() removeEvent = new EventEmitter<number>();
  isEdit: boolean;

  onToggleEdit() {
    this.isEdit = !this.isEdit;
  }

  onEditComment(val: string) {
    this.comment.content = val;
    this.onToggleEdit();
  }

  onRemoveComment() {
    this.removeEvent.emit(this.comment.id);
  }
  constructor() { }

  ngOnInit() {
    this.isEdit = false;
  }

}
