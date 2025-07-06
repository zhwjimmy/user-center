package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	redisclient "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	rediscontainer "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDB represents a test database instance
type TestDB struct {
	DB      *gorm.DB
	Mock    sqlmock.Sqlmock
	Cleanup func()
}

// TestRedis represents a test Redis instance
type TestRedis struct {
	Client  *redisclient.Client
	Cleanup func()
}

// SetupTestDB creates a test database using testcontainers
func SetupTestDB(t *testing.T) *TestDB {
	ctx := context.Background()

	// Start PostgreSQL container
	postgresContainer, err := postgrescontainer.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgrescontainer.WithDatabase("testdb"),
		postgrescontainer.WithUsername("testuser"),
		postgrescontainer.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(120*time.Second),
		),
	)
	require.NoError(t, err)

	// Get connection details
	host, err := postgresContainer.Host(ctx)
	require.NoError(t, err)
	port, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%s user=testuser password=testpass dbname=testdb sslmode=disable",
		host, port.Port())

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Run migrations
	sqlDB, err := db.DB()
	require.NoError(t, err)

	// Create test tables
	err = createTestTables(sqlDB)
	require.NoError(t, err)

	cleanup := func() {
		sqlDB.Close()
		postgresContainer.Terminate(ctx)
	}

	return &TestDB{
		DB:      db,
		Cleanup: cleanup,
	}
}

// SetupMockDB creates a mock database for unit tests
func SetupMockDB(t *testing.T) *TestDB {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	cleanup := func() {
		sqlDB.Close()
	}

	return &TestDB{
		DB:      db,
		Mock:    mock,
		Cleanup: cleanup,
	}
}

// SetupTestRedis creates a test Redis instance using testcontainers
func SetupTestRedis(t *testing.T) *TestRedis {
	ctx := context.Background()

	// Start Redis container
	redisContainer, err := rediscontainer.RunContainer(ctx,
		testcontainers.WithImage("redis:7-alpine"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithOccurrence(1).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err)

	// Get connection details
	host, err := redisContainer.Host(ctx)
	require.NoError(t, err)
	port, err := redisContainer.MappedPort(ctx, "6379")
	require.NoError(t, err)

	// Connect to Redis
	client := redisclient.NewClient(&redisclient.Options{
		Addr: fmt.Sprintf("%s:%s", host, port.Port()),
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = client.Ping(ctx).Result()
	require.NoError(t, err)

	cleanup := func() {
		client.Close()
		redisContainer.Terminate(ctx)
	}

	return &TestRedis{
		Client:  client,
		Cleanup: cleanup,
	}
}

// createTestTables creates the necessary tables for testing
func createTestTables(db *sql.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			first_name VARCHAR(100),
			last_name VARCHAR(100),
			phone VARCHAR(20),
			avatar_url TEXT,
			is_active BOOLEAN DEFAULT true,
			is_admin BOOLEAN DEFAULT false,
			email_verified BOOLEAN DEFAULT false,
			phone_verified BOOLEAN DEFAULT false,
			status VARCHAR(20) DEFAULT 'active',
			last_login_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		return err
	}

	// Create user_sessions table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user_sessions (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token_hash VARCHAR(255) UNIQUE NOT NULL,
			refresh_token_hash VARCHAR(255) UNIQUE NOT NULL,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		return err
	}

	// Create indexes
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`)
	if err != nil {
		return err
	}

	return nil
}

// LoadTestConfig loads test configuration
func LoadTestConfig() map[string]interface{} {
	return map[string]interface{}{
		"server": map[string]interface{}{
			"port": 8080,
			"host": "localhost",
		},
		"database": map[string]interface{}{
			"driver":   "postgres",
			"host":     "localhost",
			"port":     5432,
			"user":     "testuser",
			"password": "testpass",
			"dbname":   "testdb",
			"sslmode":  "disable",
		},
		"redis": map[string]interface{}{
			"host":     "localhost",
			"port":     6379,
			"password": "",
			"db":       0,
		},
		"jwt": map[string]interface{}{
			"secret":     "test-secret-key",
			"expiration": 24,
		},
		"log": map[string]interface{}{
			"level": "debug",
		},
	}
}

// CreateTempFile creates a temporary file for testing
func CreateTempFile(t *testing.T, content string) (string, func()) {
	tmpfile, err := os.CreateTemp("", "test-*.tmp")
	require.NoError(t, err)

	_, err = tmpfile.Write([]byte(content))
	require.NoError(t, err)

	err = tmpfile.Close()
	require.NoError(t, err)

	cleanup := func() {
		os.Remove(tmpfile.Name())
	}

	return tmpfile.Name(), cleanup
}

// WaitForCondition waits for a condition to be true
func WaitForCondition(t *testing.T, condition func() bool, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatal("Condition not met within timeout")
}

// RandomString generates a random string for testing
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// RandomEmail generates a random email for testing
func RandomEmail() string {
	return fmt.Sprintf("test-%s@example.com", RandomString(8))
}

// RandomUsername generates a random username for testing
func RandomUsername() string {
	return fmt.Sprintf("testuser_%s", RandomString(6))
}
