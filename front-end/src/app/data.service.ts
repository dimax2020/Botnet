import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import { Router} from '@angular/router';
import {Botinfo} from '../../model';
import {Postobj} from '../../model';

@Injectable({
  providedIn: 'root'
})
export class DataService {
  data = undefined;
  del = undefined;
  result = undefined;
  items = [] as any;
  clear = {
    id: ''};
  constructor(private http: HttpClient, private router: Router) {
  }
  getBots(): any {
    return this.data;
  }
  deleteBot(bot: Botinfo): void {
    console.log(this.data?.bots);
    const index = this.data?.bots.indexOf(bot);
    if (index > -1) {
      console.log(index);
      this.clear.id = this.data?.bots[index].id;
      console.log(this.clear);
      console.log(JSON.stringify(this.clear));
      this.http.post('/deleteBot', JSON.stringify(this.clear), { responseType: 'text' } ).subscribe(res => console.log(res));
      this.data?.bots.splice(index, 1);
    }
    console.log(this.data?.bots);
  }
  kickAll(): any {
    console.log('kickAll');
    this.http.post('/kickAll', '{"action":"clear"}', { responseType: 'text' }).subscribe(res => console.log(res));
    this.data?.bots.splice(0, this.data?.bots.length);
    return this.data?.bots;
  }
  getData(): Observable<any>{
    const data = this.http.post('/getData', '{"action":"refresh"}');
    return data;
  }
  postCommand(value: Postobj): void {
    console.log(value);
    console.log(JSON.stringify(value));
    const result = this.http.post('/postCommand', JSON.stringify(value), { responseType: 'text' }).subscribe(res => console.log(res));
    this.router.navigateByUrl('response');
  }
}
