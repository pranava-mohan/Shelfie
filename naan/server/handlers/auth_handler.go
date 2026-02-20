package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pranava-mohan/library-automation-pre/naan/config"
	"github.com/pranava-mohan/library-automation-pre/naan/server/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type DauthUserRes struct {
	AccessToken string `json:"access_token"`
}

type DauthFinalUserRes struct {
	ID          json.Number `json:"id"`
	Email       string      `json:"email"`
	Gender      string      `json:"gender"`
	Name        string      `json:"name"`
	PhoneNumber string      `json:"phoneNumber"`
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func LoginAdmin(c *fiber.Ctx) error {
	type LoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	userCollection := config.GetAdminUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user models.AdminUser
	err := userCollection.FindOne(ctx, bson.M{"username": input.Username}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if !checkPasswordHash(input.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	claims := jwt.MapClaims{
		"id":   user.ID.Hex(),
		"exp":  time.Now().Add(time.Hour * 24 * 24).Unix(),
		"type": "admin",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.JWTSecret()))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": t})
}

// Google OAuth2 Implementation

type GoogleTokenRes struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
}

type GoogleUserRes struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

func googleStep1(code string) (string, error) {
	token_url := "https://oauth2.googleapis.com/token"
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", config.Env("GOOGLE_CLIENT_ID", ""))
	data.Set("client_secret", config.Env("GOOGLE_CLIENT_SECRET", ""))
	data.Set("redirect_uri", config.Env("GOOGLE_REDIRECT_URL", "http://localhost:8000/auth/google"))
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", token_url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var googleTokenRes GoogleTokenRes
	err = json.Unmarshal(body, &googleTokenRes)
	if err != nil {
		return "", err
	}
	return googleTokenRes.AccessToken, nil
}

func googleStep2(accessToken string) (GoogleUserRes, error) {
	user_url := "https://www.googleapis.com/oauth2/v2/userinfo"
	req, err := http.NewRequest("GET", user_url, nil)
	client := &http.Client{}
	if err != nil {
		return GoogleUserRes{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return GoogleUserRes{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GoogleUserRes{}, err
	}

	var userRes GoogleUserRes
	err = json.Unmarshal(body, &userRes)
	if err != nil {
		return GoogleUserRes{}, err
	}

	return userRes, nil
}

func LoginGoogle(c *fiber.Ctx) error {
	clientID := config.Env("GOOGLE_CLIENT_ID", "")
	redirectURI := config.Env("GOOGLE_REDIRECT_URL", "http://localhost:8000/auth/google")
	baseURL := "https://accounts.google.com/o/oauth2/v2/auth"

	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("redirect_uri", redirectURI)
	params.Add("response_type", "code")
	params.Add("scope", "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile")
	params.Add("access_type", "offline")

	finalURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	return c.Redirect(finalURL)
}

func GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")

	accessToken, err := googleStep1(code)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to exchange code for token: " + err.Error()})
	}
	userRes, err := googleStep2(accessToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get user info: " + err.Error()})
	}

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := config.GetUserCollection()
	// match by google_id OR email to link existing accounts
	filter := bson.M{
		"$or": []bson.M{
			{"google_id": userRes.ID},
			{"email": userRes.Email},
		},
	}

	err = userCollection.FindOne(ctx, filter).Decode(&user)

	update := bson.M{
		"$set": bson.M{
			"google_id": userRes.ID,
			"name":      userRes.Name,
			"email":     userRes.Email,
			"picture":   userRes.Picture,
		},
		"$setOnInsert": bson.M{
			"createdAt": time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)
	user_err := userCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&user)
	if user_err != nil {
		log.Printf("findOrInsert error: %v", user_err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	claims := jwt.MapClaims{
		"id":   user.ID.Hex(),
		"exp":  time.Now().Add(time.Hour * 24 * 24).Unix(),
		"type": "normal",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(config.JWTSecret()))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Token generation error"})
	}

	frontendURL := config.Env("FRONTEND_AUTH_URL", "http://localhost:3000")
	params := url.Values{}
	params.Add("token", t)

	finalURL := fmt.Sprintf("%s?%s", frontendURL, params.Encode())
	return c.Redirect(finalURL)
}
