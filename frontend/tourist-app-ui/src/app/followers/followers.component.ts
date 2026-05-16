import { Component, OnInit } from '@angular/core';
import { BlogPost } from '../models/blog-post.model';
import { FollowerUser } from '../models/follower.model';
import { PublicUser } from '../models/user.model';
import { AuthStateService } from '../services/auth-state.service';
import { BlogApiService } from '../services/blog-api.service';
import { FollowerApiService } from '../services/follower-api.service';
import { StakeholdersApiService } from '../services/stakeholders-api.service';

@Component({
  selector: 'app-followers',
  templateUrl: './followers.component.html',
  styleUrls: ['./followers.component.css']
})
export class FollowersComponent implements OnInit {
  currentUserId = '';
  currentUsername = '';
  currentRole = '';
  allProfiles: PublicUser[] = [];
  following: FollowerUser[] = [];
  recommendations: FollowerUser[] = [];
  followedAuthorIds: string[] = [];
  followedBlogs: BlogPost[] = [];
  successMessage = '';
  errorMessage = '';
  isLoading = false;

  constructor(
    private readonly authStateService: AuthStateService,
    private readonly followerApiService: FollowerApiService,
    private readonly stakeholdersApiService: StakeholdersApiService,
    private readonly blogApiService: BlogApiService
  ) {}

  ngOnInit(): void {
    const user = this.authStateService.currentUser;
    this.currentUsername = user?.username?.trim() ?? '';
    this.currentUserId = user?.id ? String(user.id) : this.currentUsername;
    this.currentRole = user?.role ?? 'Tourist';
    if (!this.currentUserId || !this.currentUsername) {
      this.errorMessage = 'Login is required.';
      return;
    }
    this.syncCurrentUser();
  }

  isFollowing(profile: PublicUser): boolean {
    return this.following.some(user => user.userId === String(profile.id));
  }

  follow(profile: PublicUser): void {
    const target: FollowerUser = {
      userId: String(profile.id),
      username: profile.username,
      role: profile.role
    };
    this.followerApiService.syncUser(target).subscribe({
      next: () => {
        this.followerApiService.follow({ followerId: this.currentUserId, followingId: target.userId }).subscribe({
          next: () => {
            this.successMessage = `Now following ${profile.username}.`;
            this.errorMessage = '';
            this.loadAllData();
          },
          error: (error) => this.errorMessage = error.error?.error ?? 'Could not follow user.'
        });
      },
      error: (error) => this.errorMessage = error.error?.error ?? 'Could not sync target user.'
    });
  }

  followRecommended(user: FollowerUser): void {
    const existing = this.allProfiles.find(profile => String(profile.id) === user.userId);
    if (existing) {
      this.follow(existing);
      return;
    }
    this.follow({
      id: Number(user.userId),
      username: user.username,
      role: user.role
    });
  }

  unfollow(profile: PublicUser | FollowerUser): void {
    const followingId = 'id' in profile ? String(profile.id) : profile.userId;
    const username = profile.username;
    this.followerApiService.unfollow({ followerId: this.currentUserId, followingId }).subscribe({
      next: () => {
        this.successMessage = `Unfollowed ${username}.`;
        this.errorMessage = '';
        this.loadAllData();
      },
      error: (error) => this.errorMessage = error.error?.error ?? 'Could not unfollow user.'
    });
  }

  private syncCurrentUser(): void {
    this.followerApiService.syncUser({
      userId: this.currentUserId,
      username: this.currentUsername,
      role: this.currentRole
    }).subscribe({
      next: () => this.loadAllData(),
      error: () => this.loadAllData()
    });
  }

  private loadAllData(): void {
    this.isLoading = true;
    this.stakeholdersApiService.getPublicUsers().subscribe({
      next: (users) => {
        this.allProfiles = (users ?? []).filter(user => String(user.id) !== this.currentUserId);
        this.syncProfilesToGraph(this.allProfiles);
        this.loadFollowerDataAndBlogs();
      },
      error: (error) => {
        this.errorMessage = error.error?.message ?? 'Could not load user profiles.';
        this.isLoading = false;
      }
    });
  }

  private syncProfilesToGraph(users: PublicUser[]): void {
    users.forEach((user) => {
      this.followerApiService.syncUser({
        userId: String(user.id),
        username: user.username,
        role: user.role
      }).subscribe({ next: () => {}, error: () => {} });
    });
  }

  private loadFollowerDataAndBlogs(): void {
    this.followerApiService.getFollowing(this.currentUserId).subscribe({
      next: (following) => {
        this.following = following ?? [];
        this.followerApiService.getRecommendations(this.currentUserId).subscribe({
          next: (recommendations) => {
            this.recommendations = recommendations ?? [];
            this.followerApiService.getFollowedAuthorIds(this.currentUserId).subscribe({
              next: (ids) => {
                this.followedAuthorIds = ids ?? [];
                this.loadFollowedBlogs();
              },
              error: (error) => {
                this.errorMessage = error.error?.error ?? 'Could not load followed author ids.';
                this.isLoading = false;
              }
            });
          },
          error: (error) => {
            this.errorMessage = error.error?.error ?? 'Could not load recommendations.';
            this.isLoading = false;
          }
        });
      },
      error: (error) => {
        this.errorMessage = error.error?.error ?? 'Could not load following list.';
        this.isLoading = false;
      }
    });
  }

  private loadFollowedBlogs(): void {
    this.blogApiService.getBlogs().subscribe({
      next: (blogs) => {
        const followedIdSet = new Set(this.followedAuthorIds.map(id => id.trim()));
        const followedUsernameSet = new Set(
          this.following
            .filter(user => followedIdSet.has(user.userId))
            .map(user => user.username.toLowerCase())
        );
        this.followedBlogs = (blogs ?? []).filter(blog => followedUsernameSet.has(blog.authorUsername.toLowerCase()));
        this.isLoading = false;
      },
      error: (error) => {
        this.errorMessage = error.error?.message ?? 'Could not load blogs.';
        this.isLoading = false;
      }
    });
  }
}
