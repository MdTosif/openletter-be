package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mdtosif/openletter/api"
)

func main() {
	// Load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	router := gin.Default()

	config := cors.DefaultConfig()
    config.AddAllowHeaders("Authorization")
    config.AllowCredentials = true
    config.AllowAllOrigins = false
    // I think you should whitelist a limited origins instead:
    //  config.AllowAllOrigins = []{"xxxx", "xxxx"}
    config.AllowOriginFunc = func(origin string) bool {
        return true
    }
	router.Use(cors.New(config))

	// api.GetUserName("eyJhbGciOiJSUzI1NiIsImNhdCI6ImNsX0I3ZDRQRDExMUFBQSIsImtpZCI6Imluc18yZjdDUVFFYnZZYjlLeXJDcktVWGZaYmJhNjUiLCJ0eXAiOiJKV1QifQ.eyJhenAiOiJodHRwOi8vbG9jYWxob3N0OjUxNzMiLCJleHAiOjE3MTMxOTY4ODIsImlhdCI6MTcxMzE5NjgyMiwiaXNzIjoiaHR0cHM6Ly9zbWFzaGluZy1zYWxtb24tMS5jbGVyay5hY2NvdW50cy5kZXYiLCJuYmYiOjE3MTMxOTY4MTIsInNpZCI6InNlc3NfMmY3eXhCYjI0VnU3aENJdEZUTDlRTTVpNTlZIiwic3ViIjoidXNlcl8yZjd5eEVJRXl4c1h1TmxrQXVDTHkxU2FpdGIifQ.XsUZ6QywionB6cSP4MuM7Zc1LVgy8ZJhKj4b-hn2LvR0v-K9_Ob_3pxQFCcWjk4Iuf_ebxhyx-VeZcfln8oa_kMp3w3HAKmiHs0_LnpOz602pY6Cl4eNBnfY8edjl09UTgmWLpVmfH508mJR2Z4GJDdkp53OXccctYOpM0_IZatAAN99AXnc7bmmHMJauE_NND8VdQLkNVQUtM4rrU99hgqVkLDCriFf18WyAvn_Cs-E51AgzrUVTFUXhF0SDHGycb6Q9NGecXTGKH53teTy9ap1Zy9MHrFnVbGINii1gCTs3glGDSGI0Vxm10VbwVJQzb6gzvwShaP-Mtv7xIhMSA")

	// POST endpoint to create a new user
	router.POST("/letter", api.AddLetter)

	// POST endpoint to create a new user
	router.GET("/letter", api.GetLetters)

	router.Run(":8080")
}
