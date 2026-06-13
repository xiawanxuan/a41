package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"ancient-text-backend/internal/config"
	"ancient-text-backend/internal/models"

	"github.com/go-redis/redis/v8"
)

var (
	Client *redis.Client
	ctx    = context.Background()
)

const (
	annotationKeyPrefix    = "annotation:"
	annotationListKey      = "annotations:text:%s"
	annotationTTL          = 30 * time.Minute
	layoutKeyPrefix        = "layout:"
)

func Connect() {
	cfg := config.AppConfig.Redis

	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connection established")
}

func CacheAnnotation(ann *models.Annotation) error {
	key := annotationKeyPrefix + ann.ID.String()
	data, err := json.Marshal(ann)
	if err != nil {
		return err
	}
	return Client.Set(ctx, key, data, annotationTTL).Err()
}

func GetCachedAnnotation(id string) (*models.Annotation, error) {
	key := annotationKeyPrefix + id
	data, err := Client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var ann models.Annotation
	if err := json.Unmarshal(data, &ann); err != nil {
		return nil, err
	}
	return &ann, nil
}

func DeleteCachedAnnotation(id string) error {
	key := annotationKeyPrefix + id
	return Client.Del(ctx, key).Err()
}

func CacheAnnotationList(textID string, annotations []models.Annotation) error {
	key := fmt.Sprintf(annotationListKey, textID)
	data, err := json.Marshal(annotations)
	if err != nil {
		return err
	}
	return Client.Set(ctx, key, data, annotationTTL).Err()
}

func GetCachedAnnotationList(textID string) ([]models.Annotation, error) {
	key := fmt.Sprintf(annotationListKey, textID)
	data, err := Client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var annotations []models.Annotation
	if err := json.Unmarshal(data, &annotations); err != nil {
		return nil, err
	}
	return annotations, nil
}

func DeleteCachedAnnotationList(textID string) error {
	key := fmt.Sprintf(annotationListKey, textID)
	return Client.Del(ctx, key).Err()
}

func CacheLayoutConfig(layout *models.LayoutConfig) error {
	key := layoutKeyPrefix + layout.TextID.String()
	data, err := json.Marshal(layout)
	if err != nil {
		return err
	}
	return Client.Set(ctx, key, data, annotationTTL).Err()
}

func GetCachedLayoutConfig(textID string) (*models.LayoutConfig, error) {
	key := layoutKeyPrefix + textID
	data, err := Client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var layout models.LayoutConfig
	if err := json.Unmarshal(data, &layout); err != nil {
		return nil, err
	}
	return &layout, nil
}

func DeleteCachedLayoutConfig(textID string) error {
	key := layoutKeyPrefix + textID
	return Client.Del(ctx, key).Err()
}
