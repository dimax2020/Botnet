import { Component, OnInit } from '@angular/core';
import { DataService} from '../data.service';

@Component({
  selector: 'app-response',
  templateUrl: './response.component.html',
  styleUrls: ['./response.component.scss']
})
export class ResponseComponent implements OnInit {
  // response = [] as any;

  constructor(public dataService: DataService) { }

  ngOnInit(): void {
      // this.response = JSON.parse(this.dataService.resp);
  }

}
