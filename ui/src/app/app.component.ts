import { CommonModule } from '@angular/common';
import { Component, HostListener } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { MessageService } from 'primeng/api';
import { SplitButtonModule } from 'primeng/splitbutton';
import { ToastModule } from 'primeng/toast';
import { HeaderComponent } from './components/header/header.component';
import { ScrollPanelModule } from 'primeng/scrollpanel';


@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, CommonModule, SplitButtonModule, ToastModule, HeaderComponent, ScrollPanelModule],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
  providers: [MessageService]
})
export class AppComponent {

  constructor(private messageService: MessageService) {
  }

  @HostListener("window:onbeforeunload",["$event"])
  clearLocalStorage(event: any){
      localStorage.clear();
  }
}
