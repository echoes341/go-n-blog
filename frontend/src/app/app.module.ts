// Modules
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { QuillEditorModule } from 'ngx-quill-editor';
import { Routes, RouterModule, ActivatedRoute } from '@angular/router';

// Components
import { AppComponent } from './app.component';
import { HeaderComponent } from './header/header.component';
import { ArticleListComponent } from './article/article-list/article-list.component';
import { SidebarComponent } from './sidebar/sidebar.component';
import { FooterComponent } from './footer/footer.component';
import { ArticleItemComponent } from './article/article-list/article-item/article-item.component';
import { UserinfoComponent } from './sidebar/userinfo/userinfo.component';
import { ArtbydateRecapComponent } from './sidebar/artbydate-recap/artbydate-recap.component';
import { ViewArticleComponent } from './article/view-article/view-article.component';
import { CommentViewComponent } from './comments/comment-view/comment-view.component';
import { AddCommentComponent } from './comments/add-comment/add-comment.component';
import { ViewerComponent } from './viewer/viewer.component';

// Services
import { CommentService } from './comments/comment.service';
import { ArticleService } from './article/article.service';
import { LikeService } from './article/like.service';
import { CanDeactivateGuard } from './can-deactivated-guard.service';
import { EditArticleComponent } from './article/edit-article/edit-article.component';
import { DeleteArticleComponent } from './article/delete-article/delete-article.component';
import { AppRoutes } from './app-routes.module';
import { SignupComponent } from './user/signup/signup.component';
import { LoginComponent } from './user/login/login.component';
import { HttpModule } from '@angular/http';


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
    ViewerComponent,
    EditArticleComponent,
    DeleteArticleComponent,
    SignupComponent,
    LoginComponent,
  ],
  imports: [
    BrowserModule,
    FormsModule,
    QuillEditorModule,
    AppRoutes,
    HttpModule
  ],
  providers: [CommentService, ArticleService, LikeService, CanDeactivateGuard],
  bootstrap: [AppComponent]
})
export class AppModule { }
