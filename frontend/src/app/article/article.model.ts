export class Article {
   constructor (
    public id: number,
    public title: string,
    public author: string,
    public text: string,
    public date: Date // care: js months start from 0! -> 9 is october => 10
  ) { }
}
