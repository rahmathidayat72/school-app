package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DBUsername        string
	DBPassword        string
	DBHost            string
	DBPort            int
	DBName            string
	JWT_SECRET        string
	JWT_TIME_DURATION time.Duration
}

func InitConfig() *AppConfig {
	return readData()
}

func readData() *AppConfig {
	data := readEnv()
	if data == nil {
		err := godotenv.Load("local.env") //nonaktifkan kode baris ini sebelum membuild
		// err := godotenv.Load("prodaction.env") //diaktifkan sebelum membuild file sebelum di deploy ke server
		if err != nil {
			return nil
		}
		data = readEnv()
		if data == nil {
			return nil
		}
	}

	return data
}

func readEnv() *AppConfig {
	data := &AppConfig{} // Menggunakan inisialisasi pointer secara langsung

	// Gunakan satu variabel untuk menentukan apakah ada kesalahan pada saat membaca env
	var permit = true

	if val, found := os.LookupEnv("DBUSER"); found {
		data.DBUsername = val
	} else {
		permit = false
	}

	if val, found := os.LookupEnv("DBPASS"); found {
		data.DBPassword = val
	} else {
		permit = false
	}

	if val, found := os.LookupEnv("DBHOST"); found {
		data.DBHost = val
	} else {
		permit = false
	}

	if val, found := os.LookupEnv("DBPORT"); found {
		cnv, err := strconv.Atoi(val)
		if err != nil {
			permit = false
		}

		data.DBPort = cnv
	} else {
		permit = false
	}

	if val, found := os.LookupEnv("DBNAME"); found {
		data.DBName = val
	} else {
		permit = false
	}

	if val, found := os.LookupEnv("JWT_SECRET"); found {
		data.JWT_SECRET = val
	} else {
		permit = false
	}

	if !permit {
		return nil
	}

	return data
}
