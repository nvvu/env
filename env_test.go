package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnv(t *testing.T) {
	cfg := struct {
		Username  string    `env:"USERNAME"`
		Password  string    `env:"PASSWORD"`
		Count     int64     `env:"COUNT"`
		Hosts     []string  `env:"HOSTS"`
		Numbers   []float64 `env:"NUMBERS"`
		ApiServer struct {
			Host string `env:"API_HOST"`
			Port int    `env:"API_PORT"`
		}
		DB *struct {
			Addr    string `env:"DB_ADDR"`
			Port    string `env:"DB_PORT"`
			MaxConn int    `env:"DB_MAX_CONN"`
		}
	}{}

	os.Setenv("USERNAME", "admin")
	os.Setenv("PASSWORD", "123")
	os.Setenv("COUNT", "1")
	os.Setenv("HOSTS", "192.168.0.1,192.168.0.2,192.168.0.3")
	os.Setenv("NUMBERS", "1,2,3.4")
	os.Setenv("API_HOST", "0.0.0.0")
	os.Setenv("API_PORT", "9001")
	os.Setenv("DB_ADDR", "localhost")
	os.Setenv("DB_PORT", "9002")

	err := OverwriteFromEnv(&cfg)
	require.NoError(t, err, "")
	require.Equal(t, "admin", cfg.Username, "")
	require.Equal(t, "123", cfg.Password, "")
	require.Equal(t, int64(1), cfg.Count, "")

	require.ElementsMatch(t, []string{"192.168.0.1", "192.168.0.2", "192.168.0.3"}, cfg.Hosts, "")
	require.ElementsMatch(t, []float64{1, 2, 3.4}, cfg.Numbers, "")

	require.Equal(t, "0.0.0.0", cfg.ApiServer.Host, "")
	require.Equal(t, 9001, cfg.ApiServer.Port, "")

	require.Equal(t, "localhost", cfg.DB.Addr, "")
	require.Equal(t, "9002", cfg.DB.Port, "")
	require.Equal(t, 0, cfg.DB.MaxConn, "")
}

func TestEnvMap(t *testing.T) {
	type Project struct {
		Name string `env:"NAME"`
		Addr string `env:"ADDR"`
	}

	type Config struct {
		Projects map[string]*Project `env_prefix:"PROJECTS_"`
	}

	cfg := Config{
		Projects: map[string]*Project{
			"project1": {
				Name: "Project 1",
				Addr: "Addr 1",
			},
			"project2": {
				Name: "Project 2",
				Addr: "Addr 2",
			},
		},
	}
	os.Setenv("PROJECTS_project1_NAME", "project_111")
	os.Setenv("PROJECTS_project2_ADDR", "addr_222")
	err := OverwriteFromEnv(&cfg)
	require.NoError(t, err, "")

	assert.Equal(t, "project_111", cfg.Projects["project1"].Name)
	assert.Equal(t, "addr_222", cfg.Projects["project2"].Addr)
}
