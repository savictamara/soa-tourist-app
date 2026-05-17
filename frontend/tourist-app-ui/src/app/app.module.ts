import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { AuthInterceptor } from './interceptors/auth.interceptor';
import { LoginComponent } from './login/login.component';
import { ProfileComponent } from './profile/profile.component';
import { RegisterComponent } from './register/register.component';
import { UsersComponent } from './users/users.component';
import { BlogCreateComponent } from './blog-create/blog-create.component';
import { TourAuthorComponent } from './tour-author/tour-author.component';
import { FollowersComponent } from './followers/followers.component';
import { PositionSimulatorComponent } from './position-simulator/position-simulator.component';
import { TouristToursComponent } from './tourist-tours/tourist-tours.component';

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    ProfileComponent,
    RegisterComponent,
    UsersComponent,
    BlogCreateComponent,
    TourAuthorComponent,
    FollowersComponent,
    PositionSimulatorComponent,
    TouristToursComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpClientModule,
    AppRoutingModule
  ],
  providers: [
    {
      provide: HTTP_INTERCEPTORS,
      useClass: AuthInterceptor,
      multi: true
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
