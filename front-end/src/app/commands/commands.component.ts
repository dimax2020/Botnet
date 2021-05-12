import { Component, OnInit } from '@angular/core';
import {commands} from '../commands';
import {DataService} from '../data.service';
import { FormGroup, FormBuilder } from '@angular/forms';
import {FormControl, Validators} from '@angular/forms';

@Component({
  selector: 'app-commands',
  templateUrl: './commands.component.html',
  styleUrls: ['./commands.component.scss']
})
export class CommandsComponent implements OnInit {
  obj = {
    command: '',
    option: ''};
  myFirstReactiveForm: FormGroup;
  commands = commands;
  option = [] as any;
  constructor(private dataService: DataService, private fb: FormBuilder) { }

  ngOnInit(): void {
    this.initForm();
  }
  addCommandToJson(command: any): void {
    this.obj.command = command.command;
  }
  initForm(): void {
    this.myFirstReactiveForm = this.fb.group({
      option: ['', [
        Validators.required]]
    });
  }
  onSubmit(): void {
    const controls = this.myFirstReactiveForm.controls;
    if (this.myFirstReactiveForm.invalid) {
      Object.keys(controls)
        .forEach(controlName => controls[controlName].markAsTouched());
      return;
    }
    this.option = this.myFirstReactiveForm.value;
    this.obj.option = this.myFirstReactiveForm.value.option;
    this.dataService.postCommand(this.obj);
    console.log(this.obj);
  }
}
