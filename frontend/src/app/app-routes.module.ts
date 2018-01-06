import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { ArticleListComponent } from './article/article-list/article-list.component';
import { ViewArticleComponent } from './article/view-article/view-article.component';
import { CanDeactivateGuard } from './can-deactivated-guard.service';
import { DeleteArticleComponent } from './article/delete-article/delete-article.component';
import { EditArticleComponent } from './article/edit-article/edit-article.component';


const appRoutes: Routes = [
  {path: '', component: ArticleListComponent},
  {path: 'new', component: EditArticleComponent, canDeactivate: [CanDeactivateGuard]},
  {path: ':id/delete', component: DeleteArticleComponent},
  {path: ':id/edit', component: EditArticleComponent, canDeactivate: [CanDeactivateGuard]},
  {path: ':id', component: ViewArticleComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(appRoutes)],
  exports: [RouterModule]
})
export class AppRoutes {}
