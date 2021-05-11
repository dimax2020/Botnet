import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {CommandsComponent} from './commands/commands.component';
import {MainComponent} from './main/main.component';
import {ReadmeComponent} from './readme/readme.component';
import {ResponseComponent} from './response/response.component';

const routes: Routes = [{ path: 'commands', component: CommandsComponent },
  { path: 'response', component: ResponseComponent},
  { path: 'readme', component: ReadmeComponent },
  { path: 'main', component: MainComponent },
  {path: '**', redirectTo: '/main', pathMatch: 'full'},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
