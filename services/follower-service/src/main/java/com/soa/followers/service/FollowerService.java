package com.soa.followers.service;

import com.soa.followers.dto.CreateUserRequest;
import com.soa.followers.dto.FollowRequest;
import com.soa.followers.model.UserNode;
import com.soa.followers.repository.UserRepository;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import java.util.List;

@Service
public class FollowerService {
    private final UserRepository userRepository;

    public FollowerService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    public UserNode createUser(CreateUserRequest request) {
        if (isBlank(request.userId()) || isBlank(request.username()) || isBlank(request.role())) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "userId, username and role are required");
        }
        UserNode user = new UserNode();
        user.setUserId(request.userId().trim());
        user.setUsername(request.username().trim());
        user.setRole(request.role().trim());
        return userRepository.save(user);
    }

    public UserNode syncUser(CreateUserRequest request) {
        return createUser(request);
    }

    public List<UserNode> getAllUsers() {
        return userRepository.findAll();
    }

    public void follow(FollowRequest request) {
        validateFollowRequest(request);
        ensureExists(request.followerId());
        ensureExists(request.followingId());
        if (userRepository.isFollowing(request.followerId(), request.followingId())) {
            throw new ResponseStatusException(HttpStatus.CONFLICT, "already following user");
        }
        userRepository.follow(request.followerId(), request.followingId());
    }

    public void unfollow(FollowRequest request) {
        validateFollowRequest(request);
        if (!userRepository.isFollowing(request.followerId(), request.followingId())) {
            throw new ResponseStatusException(HttpStatus.NOT_FOUND, "follow relation not found");
        }
        userRepository.unfollow(request.followerId(), request.followingId());
    }

    public List<UserNode> getFollowing(String userId) {
        ensureExists(userId);
        return userRepository.findFollowing(userId);
    }

    public List<UserNode> getFollowers(String userId) {
        ensureExists(userId);
        return userRepository.findFollowers(userId);
    }

    public List<UserNode> getRecommendations(String userId) {
        ensureExists(userId);
        return userRepository.findRecommendations(userId);
    }

    public List<String> getFollowedAuthorIds(String userId) {
        ensureExists(userId);
        return userRepository.findFollowedAuthorIds(userId);
    }

    private void validateFollowRequest(FollowRequest request) {
        if (isBlank(request.followerId()) || isBlank(request.followingId())) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "followerId and followingId are required");
        }
        if (request.followerId().trim().equals(request.followingId().trim())) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "user cannot follow themselves");
        }
    }

    private void ensureExists(String userId) {
        if (isBlank(userId)) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "userId is required");
        }
        if (userRepository.findById(userId.trim()).isEmpty()) {
            throw new ResponseStatusException(HttpStatus.NOT_FOUND, "user not found: " + userId);
        }
    }

    private boolean isBlank(String value) {
        return value == null || value.trim().isEmpty();
    }
}
