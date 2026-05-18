package com.soa.followers.controller;

import com.soa.followers.dto.CreateUserRequest;
import com.soa.followers.dto.FollowRequest;
import com.soa.followers.model.UserNode;
import com.soa.followers.service.FollowerService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/followers")
public class FollowerController {
    private final FollowerService followerService;

    public FollowerController(FollowerService followerService) {
        this.followerService = followerService;
    }

    @GetMapping("/health")
    public Map<String, String> health() {
        return Map.of("status", "ok");
    }

    @PostMapping("/users")
    public ResponseEntity<UserNode> createUser(@RequestBody CreateUserRequest request) {
        return ResponseEntity.status(HttpStatus.CREATED).body(followerService.createUser(request));
    }

    @PostMapping("/sync-user")
    public ResponseEntity<UserNode> syncUser(@RequestBody CreateUserRequest request) {
        return ResponseEntity.ok(followerService.syncUser(request));
    }

    @GetMapping("/users")
    public List<UserNode> getUsers() {
        return followerService.getAllUsers();
    }

    @PostMapping("/follow")
    public ResponseEntity<Map<String, String>> follow(@RequestBody FollowRequest request) {
        followerService.follow(request);
        return ResponseEntity.status(HttpStatus.CREATED).body(Map.of("message", "followed successfully"));
    }

    @DeleteMapping("/follow")
    public Map<String, String> unfollow(@RequestBody FollowRequest request) {
        followerService.unfollow(request);
        return Map.of("message", "unfollowed successfully");
    }

    @GetMapping("/{userId}/following")
    public List<UserNode> getFollowing(@PathVariable String userId) {
        return followerService.getFollowing(userId);
    }

    @GetMapping("/{userId}/followers")
    public List<UserNode> getFollowers(@PathVariable String userId) {
        return followerService.getFollowers(userId);
    }

    @GetMapping("/{userId}/recommendations")
    public List<UserNode> getRecommendations(@PathVariable String userId) {
        return followerService.getRecommendations(userId);
    }

    @GetMapping("/{userId}/followed-author-ids")
    public List<String> getFollowedAuthorIds(@PathVariable String userId) {
        return followerService.getFollowedAuthorIds(userId);
    }

    @GetMapping("/{followerId}/follows/{targetId}")
    public boolean isFollowing(@PathVariable String followerId, @PathVariable String targetId) {
        return followerService.isFollowing(followerId, targetId);
    }
}
