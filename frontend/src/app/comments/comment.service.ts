import { Comment } from './comment.model';

export class CommentService {
  public c = [
    new Comment(1, 1, 'Gianfranco', 'gg@gg.it', 'Nice!'),
    new Comment(2, 1, 'Luigino', 'sdsdgg@gg.it', 'Happy Christmas!'),
    new Comment(3, 2, 'Luise', 'luise@lovecat.com', 'I perfectly agree with you')
  ];
  constructor() { }

  public getCommentByArtID(id: number): Comment[] {
    return this.c.filter(function(singleC: Comment) {
      return singleC.idA === id;
    });
  }
  public getCountCommentByArtID(id: number): number {
    let i = 0;
    this.c.forEach(element => {
      if (element.idA === id) {
        i++;
      }
    });
    return i;
  }
}
