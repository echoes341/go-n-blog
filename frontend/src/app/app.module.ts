// Modules
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

// Components
import { AppComponent } from './app.component';
import { HeaderComponent } from './header/header.component';
import { ArticleListComponent } from './article/article-list/article-list.component';
import { SidebarComponent } from './sidebar/sidebar.component';
import { FooterComponent } from './footer/footer.component';
import { ArticleItemComponent } from './article/article-list/article-item/article-item.component';
import { UserinfoComponent } from './sidebar/userinfo/userinfo.component';
import { ArtbydateRecapComponent } from './sidebar/artbydate-recap/artbydate-recap.component';
import { Routes, RouterModule, ActivatedRoute } from '@angular/router';
import { ViewArticleComponent } from './article/view-article/view-article.component';
import { CommentViewComponent } from './comments/comment-view/comment-view.component';
import { AddCommentComponent } from './comments/add-comment/add-comment.component';
import { AddArticleComponent } from './article/add-article/add-article.component';
import { ViewerComponent } from './viewer/viewer.component';

// Services
import { CommentService } from './comments/comment.service';
import { ArticleService } from './article/article.service';
import { LikeService } from './article/like.service';

const appRoutes: Routes = [
  {path: '', component: ArticleListComponent},
  {path: 'article/v/:id', component: ViewArticleComponent},
  {path: 'article/new', component: AddArticleComponent}
];

@NgModule({
  declarations: [
    AppComponent,
    HeaderComponent,
    ArticleListComponent,
    SidebarComponent,
    FooterComponent,
    ArticleItemComponent,
    UserinfoComponent,
    ArtbydateRecapComponent,
    ViewArticleComponent,
    CommentViewComponent,
    AddCommentComponent,
    AddArticleComponent,
    ViewerComponent
  ],
  imports: [
    BrowserModule,
    RouterModule.forRoot(appRoutes)
  ],
  providers: [CommentService, ArticleService, LikeService],
  bootstrap: [AppComponent]
})
export class AppModule { }
