package com.soa.followers.repository;

import com.soa.followers.model.UserNode;
import org.springframework.data.neo4j.repository.Neo4jRepository;
import org.springframework.data.neo4j.repository.query.Query;

import java.util.List;

public interface UserRepository extends Neo4jRepository<UserNode, String> {
    @Query("MATCH (u:User {userId: $userId})-[:FOLLOWS]->(f:User) RETURN f")
    List<UserNode> findFollowing(String userId);

    @Query("MATCH (f:User)-[:FOLLOWS]->(u:User {userId: $userId}) RETURN f")
    List<UserNode> findFollowers(String userId);

    @Query("""
        MATCH (u:User {userId: $userId})-[:FOLLOWS]->(:User)-[:FOLLOWS]->(rec:User)
        WHERE rec.userId <> $userId
          AND NOT (u)-[:FOLLOWS]->(rec)
        RETURN DISTINCT rec
        """)
    List<UserNode> findRecommendations(String userId);

    @Query("MATCH (u:User {userId: $userId})-[:FOLLOWS]->(f:User) RETURN f.userId")
    List<String> findFollowedAuthorIds(String userId);

    @Query("""
        MATCH (a:User {userId: $followerId}), (b:User {userId: $followingId})
        MERGE (a)-[r:FOLLOWS]->(b)
        RETURN count(r)
        """)
    long follow(String followerId, String followingId);

    @Query("""
        MATCH (a:User {userId: $followerId})-[r:FOLLOWS]->(b:User {userId: $followingId})
        DELETE r
        RETURN count(r)
        """)
    long unfollow(String followerId, String followingId);

    @Query("MATCH (a:User {userId: $followerId})-[r:FOLLOWS]->(b:User {userId: $followingId}) RETURN count(r) > 0")
    boolean isFollowing(String followerId, String followingId);
}
