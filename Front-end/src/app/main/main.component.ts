import { Component, OnInit } from '@angular/core';
import { DataService } from '../data.service';

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.scss']
})
export class MainComponent implements OnInit {

  constructor(public dataService: DataService) { }

  ngOnInit(): any {
    this.dataService.getData().subscribe(res => console.log(res));
    // setInterval( () => this.dataService.getData().subscribe(res => {
    //   this.dataService.data = JSON.parse(res);
    // }), 30000);
    setInterval( () => this.dataService.getData().subscribe(res => this.dataService.data = res), 30000);

  }
  getBots(): void{
    this.dataService.getBots();
  }
  deleteBot(bot: object): any{
    this.dataService.deleteBot(bot);
  }
}
