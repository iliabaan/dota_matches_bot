package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type Client struct {
	Redis *redis.Client
}

func NewRedisDb(ctx context.Context, config Config) (*Client, error) {
	client := new(Client)

	client.Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	err := client.Redis.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Set(ctx context.Context, key string, value string) {
	c.Redis.Set(ctx, key, value, 0)
}

func (c *Client) GetLastMatch(ctx context.Context, playerID int) (int64, error) {
	// Generate the Redis key for the player's last match
	key := "last_matches"

	// Get the last match ID for the player from Redis
	matchIDStr, err := c.Redis.HGet(ctx, key, strconv.FormatUint(uint64(playerID), 10)).Result()
	if err == redis.Nil {
		// The player's last match is not in the hash
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	// Parse the match ID as an int64
	matchID, err := strconv.ParseInt(matchIDStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return matchID, nil
}

// SetLastMatch sets the last match ID for a player in the "last_matches" hash in Redis.
func (c *Client) SetLastMatch(ctx context.Context, playerID int, matchID int64) error {
	// Generate the Redis key for the player's last match
	key := "last_matches"

	// Set the last match ID for the player in Redis
	_, err := c.Redis.HSet(ctx, key, strconv.FormatUint(uint64(playerID), 10), matchID).Result()
	if err != nil {
		return err
	}

	return nil
}
