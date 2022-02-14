package dockertests

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"

	"github.com/ory/dockertest/v3"
)

func ClientWithDockerTest() (client *redis.Client, cleanup func(), err error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, err
	}

	resource, err := pool.Run("eqalpha/keydb", "x86_64_v6.0.18", nil)
	if err != nil {
		return nil, nil, err
	}

	if err = pool.Retry(func() error {
		client = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
		})

		return client.Ping(context.TODO()).Err()
	}); err != nil {
		return nil, nil, err
	}

	cleanup = func() {
		if errClose := client.Close(); errClose != nil {
			log.Println(errClose)
		}

		if errPurge := pool.Purge(resource); errPurge != nil {
			log.Println(errPurge)
		}
	}

	return client, cleanup, nil
}
