package main

import (
	"log"

	"sync"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pranava-mohan/library-automation-pre/naan/config"
	"github.com/pranava-mohan/library-automation-pre/naan/server/routes"
)

func UserStatusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		if user == nil {
			return fiber.ErrUnauthorized
		}
		claims := user.Claims.(jwt.MapClaims)
		userType := claims["type"].(string)

		c.Locals("user_type", userType)
		c.Locals("user_id", claims["id"].(string))

		return c.Next()
	}
}

type Room struct {
	ID        string
	Conn      *websocket.Conn
	WriteChan chan string
}

type Hub struct {
	Rooms map[string]*Room
	mu    sync.RWMutex
}

var hub = Hub{
	Rooms: make(map[string]*Room),
}

type UserCheckInReq struct {
	UserID string `json:"UserID"`
}

func main() {
	config.LoadEnv()
	app := fiber.New()
	app.Use(cors.New())

	config.ConnectDB()

	routes.InitAuth(app)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("naan running")
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		roomID := c.Params("id")
		kioskToken := c.Query("token")

		token, err := jwt.Parse(kioskToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret()), nil
		})
		if err != nil || !token.Valid {
			log.Println("Invalid kiosk token")
			c.WriteJSON(fiber.Map{"error": "Invalid kiosk token"})
			c.Close()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if claims["type"] != "kiosk" {
			log.Println("Unauthorized kiosk")
			c.WriteJSON(fiber.Map{"error": "Unauthorized kiosk"})
			c.Close()
			return
		}

		hub.mu.Lock()
		if _, exists := hub.Rooms[roomID]; exists {
			hub.mu.Unlock()
			c.WriteJSON(fiber.Map{"error": "Room already exists"})
			c.Close()
			return
		}

		room := &Room{
			ID:        roomID,
			Conn:      c,
			WriteChan: make(chan string),
		}
		hub.Rooms[roomID] = room
		hub.mu.Unlock()

		log.Printf("Host joined room: %s", roomID)

		go func() {
			for msg := range room.WriteChan {
				if err = c.WriteJSON(fiber.Map{"id": msg}); err != nil {
					log.Println("Write error:", err)
					break // Stop loop on error
				}
			}
		}()
		var (
			mt  int
			msg []byte
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				break // Connection closed
			}
			// Optional: Host can send ping/pong or commands back if needed
			log.Printf("Recv from host %s: %s type: %d", roomID, msg, mt)
		}

		hub.mu.Lock()
		delete(hub.Rooms, roomID)
		close(room.WriteChan)
		hub.mu.Unlock()
		log.Printf("Room closed: %s", roomID)
	}))

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.JWTSecret())},
	}))
	app.Use(UserStatusMiddleware())

	routes.InitUser(app)
	routes.InitShelf(app)
	routes.InitBook(app)
	routes.InitKiosk(app)

	app.Post("/check-in/:id", func(c *fiber.Ctx) error {
		roomID := c.Params("id")

		checkInData := new(UserCheckInReq)

		if err := c.BodyParser(checkInData); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		hub.mu.RLock()
		room, exists := hub.Rooms[roomID]
		hub.mu.RUnlock()

		if !exists {
			return c.Status(404).JSON(fiber.Map{"error": "Room not found or host not connected"})
		}

		room.WriteChan <- checkInData.UserID

		return c.JSON(fiber.Map{"status": "sent"})
	})

	app.Listen(":8000")
}
