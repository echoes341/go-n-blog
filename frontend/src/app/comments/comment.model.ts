export class Comment {
  constructor(
    public id: number, // comment's own id
    public idA: number, // id of the article which the comment is related to
    public name: string, // name of the author
    public email: string, // email of the autorh
    public content: string // text content of the comment
  ) { }
}
