package repositories_test

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github/feroddev/challengeV3/internal/core"
	"github/feroddev/challengeV3/internal/pkg/database"
	"github/feroddev/challengeV3/internal/repositories"
)

func setupTestDB() (*repositories.PostgresTelemetryRepository, func()) {
	godotenv.Load("../../.env.test")

	db, err := database.NewDatabaseConnection()
	if err != nil {
		panic(err)
	}

	err = database.MigrateDatabase(db)
	if err != nil {
		panic(err)
	}

	repo := repositories.NewPostgresTelemetryRepository(db)

	cleanup := func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}

	return repo.(*repositories.PostgresTelemetryRepository), cleanup
}

func TestSaveAndGetGyroscopeData(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION") == "true" {
		t.Skip("Pulando teste de integração")
	}

	repo, cleanup := setupTestDB()
	defer cleanup()

	gyroscopeData := core.Gyroscope{
		X:         10.5,
		Y:         -5.2,
		Z:         3.7,
		Timestamp: time.Now(),
		DeviceID:  "test-device",
	}

	err := repo.SaveGyroscopeData(gyroscopeData)
	assert.NoError(t, err)

	data, err := repo.GetGyroscopeData()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	found := false
	for _, d := range data {
		if d.DeviceID == gyroscopeData.DeviceID && d.X == gyroscopeData.X {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestSaveAndGetGPSData(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION") == "true" {
		t.Skip("Pulando teste de integração")
	}

	repo, cleanup := setupTestDB()
	defer cleanup()

	gpsData := core.GPS{
		Latitude:  -23.5505,
		Longitude: -46.6333,
		Timestamp: time.Now(),
		DeviceID:  "test-device",
	}

	err := repo.SaveGPSData(gpsData)
	assert.NoError(t, err)

	data, err := repo.GetGPSData()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	found := false
	for _, d := range data {
		if d.DeviceID == gpsData.DeviceID && d.Latitude == gpsData.Latitude {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestSaveAndGetPhotoData(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION") == "true" {
		t.Skip("Pulando teste de integração")
	}

	repo, cleanup := setupTestDB()
	defer cleanup()

	photoData := core.Photo{
		Photo:     "base64_encoded_string",
		Timestamp: time.Now(),
		DeviceID:  "test-device",
	}

	err := repo.SavePhotoData(photoData)
	assert.NoError(t, err)

	data, err := repo.GetPhotoData()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	found := false
	for _, d := range data {
		if d.DeviceID == photoData.DeviceID {
			found = true
			break
		}
	}
	assert.True(t, found)
}
