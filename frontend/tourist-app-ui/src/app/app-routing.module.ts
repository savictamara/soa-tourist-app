import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { AdminGuard } from './guards/admin.guard';
import { AuthGuard } from './guards/auth.guard';
import { BlogUserGuard } from './guards/blog-user.guard';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { ProfileComponent } from './profile/profile.component';
import { UsersComponent } from './users/users.component';
import { BlogCreateComponent } from './blog-create/blog-create.component';
import { TourAuthorComponent } from './tour-author/tour-author.component';
import { TourAuthorGuard } from './guards/tour-author.guard';
import { FollowersComponent } from './followers/followers.component';
import { PositionSimulatorComponent } from './position-simulator/position-simulator.component';
import { TouristToursComponent } from './tourist-tours/tourist-tours.component';
import { TouristGuard } from './guards/tourist.guard';
import { ShoppingCartComponent } from './shopping-cart/shopping-cart.component';

const routes: Routes = [
  { path: '', redirectTo: 'login', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'profile', component: ProfileComponent, canActivate: [AuthGuard] },
  { path: 'blogs', component: BlogCreateComponent, canActivate: [BlogUserGuard] },
  { path: 'tours', component: TourAuthorComponent, canActivate: [TourAuthorGuard] },
  { path: 'tourist-tours', component: TouristToursComponent, canActivate: [TouristGuard] },
  { path: 'shopping-cart', component: ShoppingCartComponent, canActivate: [TouristGuard] },
  { path: 'position-simulator', component: PositionSimulatorComponent, canActivate: [TouristGuard] },
  { path: 'followers', component: FollowersComponent, canActivate: [AuthGuard] },
  { path: 'users', component: UsersComponent, canActivate: [AdminGuard] }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
