package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/s1ntezc0der/bazis-restapi/internal/config"
	"github.com/s1ntezc0der/bazis-restapi/internal/run"
)

const (
	testEmail    = "test@example.com"
	testPassword = "123456"
	testName     = "Test User"
)

// TestIntegration_Smoke — проверяет, что сервер запускается и отвечает на запросы
func TestIntegration_Smoke(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	// 1. Поднимаем MySQL контейнер
	mysqlContainer, err := startMySQLContainer(ctx)
	require.NoError(t, err)
	defer mysqlContainer.Terminate(ctx)

	// 2. Поднимаем Redis контейнер
	redisContainer, err := startRedisContainer(ctx)
	require.NoError(t, err)
	defer redisContainer.Terminate(ctx)

	// 3. Загружаем конфиг с динамическими портами
	cfg := config.Load()
	cfg.DB.Host = mysqlContainer.Host
	cfg.DB.Port = mysqlContainer.Port
	cfg.DB.User = "test"
	cfg.DB.Password = "test"
	cfg.DB.Name = "test"
	cfg.RedisAddr = fmt.Sprintf("%s:%s", redisContainer.Host, redisContainer.Port)
	cfg.Server.Port = "8081"

	// 4. Запускаем сервер в горутине
	go func() {
		if err := run.Run(); err != nil {
			t.Logf("Server error: %v", err)
		}
	}()

	// 5. Ждём, пока сервер запустится
	time.Sleep(2 * time.Second)

	baseURL := "http://localhost:" + cfg.Server.Port

	// 6. Тестируем регистрацию
	t.Run("Register", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    testEmail,
			"password": testPassword,
			"name":     testName,
		}
		body, _ := json.Marshal(reqBody)

		resp, err := http.Post(baseURL+"/api/v1/register", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var user map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&user)
		require.NoError(t, err)
		assert.Equal(t, testEmail, user["email"])
		assert.Equal(t, testName, user["name"])
		assert.NotEmpty(t, user["id"])
	})

	// 7. Тестируем логин
	t.Run("Login", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    testEmail,
			"password": testPassword,
		}
		body, _ := json.Marshal(reqBody)

		resp, err := http.Post(baseURL+"/api/v1/login", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.NotEmpty(t, result["token"])
		assert.NotEmpty(t, result["user"])
	})

	// 8. Тестируем создание команды (требует токен)
	token := getToken(t, baseURL)

	t.Run("CreateTeam", func(t *testing.T) {
		reqBody := map[string]string{
			"name":        "Test Team",
			"description": "Integration test team",
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", baseURL+"/api/v1/teams", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var team map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&team)
		require.NoError(t, err)
		assert.Equal(t, "Test Team", team["name"])
		assert.Equal(t, "Integration test team", team["description"])
	})

	// 9. Тестируем создание задачи
	t.Run("CreateTask", func(t *testing.T) {
		// Сначала создаём команду, чтобы получить team_id
		teamID := createTeam(t, baseURL, token)

		reqBody := map[string]interface{}{
			"title":       "Integration Task",
			"description": "Test task",
			"team_id":     teamID,
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", baseURL+"/api/v1/tasks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var task map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&task)
		require.NoError(t, err)
		assert.Equal(t, "Integration Task", task["title"])
		assert.Equal(t, "todo", task["status"])
	})
}

// startMySQLContainer — поднимает MySQL в Docker
func startMySQLContainer(ctx context.Context) (*testcontainers.DockerContainer, error) {
	req := testcontainers.ContainerRequest{
		Image: "mysql:8",
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_DATABASE":      "test",
			"MYSQL_USER":          "test",
			"MYSQL_PASSWORD":      "test",
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server").WithStartupTimeout(30 * time.Second),
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

// startRedisContainer — поднимает Redis в Docker
func startRedisContainer(ctx context.Context) (*testcontainers.DockerContainer, error) {
	req := testcontainers.ContainerRequest{
		Image: "redis:7-alpine",
		WaitingFor: wait.ForLog("Ready to accept connections").
			WithStartupTimeout(30 * time.Second),
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

// getToken — вспомогательная функция для получения JWT-токена
func getToken(t *testing.T, baseURL string) string {
	reqBody := map[string]string{
		"email":    testEmail,
		"password": testPassword,
	}
	body, _ := json.Marshal(reqBody)

	resp, err := http.Post(baseURL+"/api/v1/login", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	token, ok := result["token"].(string)
	require.True(t, ok)
	require.NotEmpty(t, token)

	return token
}

// createTeam — вспомогательная функция для создания команды
func createTeam(t *testing.T, baseURL, token string) int64 {
	reqBody := map[string]string{
		"name":        "Test Team for Tasks",
		"description": "Created for testing tasks",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", baseURL+"/api/v1/teams", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	var team map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&team)
	require.NoError(t, err)

	teamID, ok := team["id"].(float64)
	require.True(t, ok)
	return int64(teamID)
}
