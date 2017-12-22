import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';


import { AppComponent } from './app.component';
import { HeaderComponent } from './header/header.component';
import { ArticleListComponent } from './article/article-list/article-list.component';
import { SidebarComponent } from './sidebar/sidebar.component';
import { FooterComponent } from './footer/footer.component';
import { ArticleItemComponent } from './article/article-list/article-item/article-item.component';


@NgModule({
  declarations: [
    AppComponent,
    HeaderComponent,
    ArticleListComponent,
    SidebarComponent,
    FooterComponent,
    ArticleItemComponent
  ],
  imports: [
    BrowserModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
