import { Article, ArticleRecap } from './article.model';
import { Http, Response, Headers } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/map';
import { Injectable } from '@angular/core';

@Injectable()
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
      `Lorem ipsum dolor sit amet consectetur adipisicing elit.<br>
    Qui quasi eum eveniet perspiciatis repellat minus tempora in
    reprehenderit porro culpa! Reprehenderit dolorem esse ullam, saepe atque quia
    distinctio recusandae maxime possimus soluta!<br>
    Aliquid exercitationem veritatis, dolores culpa ratione sint expedita.Lorem ipsum dolor sit amet consectetur adipisicing elit.<br>
    Qui quasi eum eveniet perspiciatis repellat minus tempora in
    reprehenderit porro culpa! Reprehenderit dolorem esse ullam, saepe atque quia
    distinctio recusandae maxime possimus soluta!<br>
    Aliquid exercitationem veritatis, dolores culpa ratione sint expedita.Lorem ipsum dolor sit amet consectetur adipisicing elit.<br>
    Qui quasi eum eveniet perspiciatis repellat minus tempora in
    reprehenderit porro culpa! Reprehenderit dolorem esse ullam, saepe atque quia
    distinctio recusandae maxime possimus soluta!<br>
    Aliquid exercitationem veritatis, dolores culpa ratione sint expedita.Lorem ipsum dolor sit amet consectetur adipisicing elit.<br>
    Qui quasi eum eveniet perspiciatis repellat minus tempora in
    reprehenderit porro culpa! Reprehenderit dolorem esse ullam, saepe atque quia
    distinctio recusandae maxime possimus soluta!<br>
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

  baseUrl = 'http://localhost:8080';

  constructor(private http: Http) {}

  /* All these methods should be handled by API */
  public getArticleByID(id: number): Observable<Article> {
    /*let i = 0;
    for (i = 0; i < this.articles.length; i++) {
      // tslint:disable-next-line:triple-equals
      if (this.articles[i].id === id) {
        return this.articles[i];
      }
    }*/
    const a = this.http
      .get(`${this.baseUrl}/article/${id}`, { headers: this.getHeaders() })
      .map(mapArticle);
    return a;

    function mapArticle(response: Response): Article {
      return response.json().results.map(this.toArticle);
    }
  }

  toArticle(r: any): Article {
    const article = <Article>{
      id: r.id,
      title: r.title,
      author: r.author,
      text: r.text,
      date: r.date
    };
    console.log('Parsed article: ', article);
    return article;
  }

  private getHeaders() {
    const headers = new Headers();
    headers.append('Accept', 'application/json');
    return headers;
  }

  public getFirstsXFromDate(x: number, d: Date): Article[] {
    let result: Article[];
    result = this.articles.filter(element => {
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
      new ArticleRecap(2018, 11, 4)
    ];
  }

  lastId(): number {
    let id = 0;
    id = this.articles[0].id;
    for (let i = 1; i < this.articles.length; i++) {
      if (id < this.articles[i].id) {
        id = this.articles[i].id;
      }
    }
    return id;
  }

  addArticle(a: Article): number {
    const id = this.lastId() + 1;
    a.id = id;
    this.articles.push(a);
    return id;
  }

  editArticle(a: Article): boolean {
    let i;
    for (i = 0; i < this.articles.length; i++) {
      if (this.articles[i].id === a.id) {
        this.articles[i] = a;
        return true;
      }
    }
    return false;
  }

  /* sort articles by date */
  public sort() {
    this.sortArticles(this.articles);
  }
  public sortArticles(art: Article[]) {
    art.sort(function(a, b) {
      // order from the most recent
      const d1 = a.date;
      const d2 = b.date;
      return d1 > d2 ? -1 : d1 < d2 ? +1 : 0;
    });
  }
}
