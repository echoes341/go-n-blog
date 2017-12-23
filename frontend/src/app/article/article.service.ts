import { Article, ArticleRecap } from './article.model';


export class ArticleService {
  public articles = [
    new Article(
      1,
      'Titolo',
      'echoes',
      'This is an article test, I hope it works.',
      new Date(2017, 11, 22, 17, 23)
    ),
    new Article(
      5,
      'Articolo 2',
      'echoes',
      `Lorem ipsum dolor sit amet consectetur adipisicing elit.
    Qui quasi eum eveniet perspiciatis repellat minus tempora in
    reprehenderit porro culpa! Reprehenderit dolorem esse ullam, saepe atque quia
    distinctio recusandae maxime possimus soluta!
    Aliquid exercitationem veritatis, dolores culpa ratione sint expedita. `,
      new Date(2017, 9, 21, 17, 21)
    ),
    new Article(
      2,
      'Articolo 1',
      'echoes',
      `Lorem ipsum dolor sit amet consectetur adipisicing elit.
    Qui quasi eum eveniet perspiciatis repellat minus tempora in
    reprehenderit porro culpa! Reprehenderit dolorem esse ullam, saepe atque quia
    distinctio recusandae maxime possimus soluta!
    Aliquid exercitationem veritatis, dolores culpa ratione sint expedita. `,
      new Date(2017, 9, 22, 17, 21)
    )
  ];

  public getFirstsXFromDate(x: number, d: Date): Article[] {
    let result: Article[];
    result = this.articles.filter( element => {
      return element.date < d;
    });
    this.sortArticles(result);
    return result.slice(0, x);
  }

  public getArticlesRecap(): ArticleRecap[] {
    /*
      get from API-> TODO
    */
    return [
      new ArticleRecap(2017, 11, 4),
      new ArticleRecap(2017, 9, 2),
      new ArticleRecap(2017, 0, 1),
      new ArticleRecap(2018, 11, 4),
    ];

  }


  /* sort articles by date */
  public sort() {
    this.sortArticles(this.articles);
  }
  public sortArticles(art: Article[]) {
    art.sort( function(a, b) { // order from the most recent
      const d1 = a.date;
      const d2 = b.date;
      return d1 > d2 ? -1 : d1 < d2 ? +1 : 0;
    });
  }
}
