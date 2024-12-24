package user

import (
	"fmt"
	"github/Shubhpreet-Rana/projects/config"
	"github/Shubhpreet-Rana/projects/internal/logging"
	"github/Shubhpreet-Rana/projects/internal/telemetry"
	"github/Shubhpreet-Rana/projects/service/auth"
	"github/Shubhpreet-Rana/projects/types"
	"github/Shubhpreet-Rana/projects/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

// Update this method to accept a fiber.Group
func (h *Handler) RegisterRoutes(group fiber.Router) {
	// Register routes under the "/api/v1" group
	fmt.Println("Entered Home 3")
	group.Get("/", h.handleHome)
	group.Post("/login", h.handleLogin)
	group.Post("/register", h.handleRegister)
}

// handleHome godoc
// @Summary Home endpoint
// @Description Returns a welcome message
// @Tags home
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/ [get]
func (h *Handler) handleHome(c *fiber.Ctx) error {
	fmt.Println("Entered Home")
	_, span := telemetry.StartSpan(c.Context(), telemetry.Tracer, "HandleLogin")
	defer span.End()
	logging.InfoLogger.Println("Welcome to Home")
	// Return a welcome message as a JSON response
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Welcome to the Home page!"})
}

// handleLogin godoc
// @Summary Login user
// @Description Login using email and password
// @Tags user
// @Accept json
// @Produce json
// @Param login body types.LoginUserPayload true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/login [post]
func (h *Handler) handleLogin(c *fiber.Ctx) error {
	_, span := telemetry.StartSpan(c.Context(), telemetry.Tracer, "HandleLogin")
	defer span.End()

	var payload types.LoginUserPayload
	if err := c.BodyParser(&payload); err != nil {
		logging.ErrorLogger.Println("Failed to parse payload:", err)
		// Return a welcome message as a JSON response
		return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Welcome to the Home page!"})

	}

	// Validate payload
	if err := utils.Validator.Struct(payload); err != nil {
		logging.ErrorLogger.Println("Validation error:", err)
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{"error": "invalid input"})
	}

	// Check if user exists
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		logging.ErrorLogger.Println("User not found:", err)
		return c.Status(http.StatusUnauthorized).JSON(map[string]interface{}{"error": "invalid email or password"})
	}

	// Check password
	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		logging.ErrorLogger.Println("Invalid password for user:", u.Email)
		return c.Status(http.StatusUnauthorized).JSON(map[string]interface{}{"error": "invalid email or password"})
	}

	// Generate JWT
	secret := []byte(config.Env.JWTSecret)
	token, err := auth.CreateJwt(secret, u.ID)
	if err != nil {
		logging.ErrorLogger.Println("JWT creation failed:", err)
		return c.Status(http.StatusInternalServerError).JSON(map[string]interface{}{"error": "failed to generate token"})
	}

	logging.InfoLogger.Println("User logged in successfully:", u.Email)
	return c.JSON(map[string]interface{}{"token": token})
}

func (h *Handler) handleRegister(c *fiber.Ctx) error {
	logging.InfoLogger.Println("Received request for registration")

	_, span := telemetry.StartSpan(c.Context(), telemetry.Tracer, "HandleRegister")
	defer span.End()

	var payload types.RegisterUserPayload
	if err := c.BodyParser(&payload); err != nil {
		logging.ErrorLogger.Println("Failed to parse payload:", err)
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{"error": "invalid request body"})
	}

	logging.InfoLogger.Println("Parsed payload:", payload)

	if err := utils.Validator.Struct(payload); err != nil {
		logging.ErrorLogger.Println("Validation error:", err)
		return c.Status(http.StatusBadRequest).JSON(map[string]interface{}{"error": "invalid input"})
	}

	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		logging.ErrorLogger.Println("User already exists:", payload.Email)
		return c.Status(http.StatusConflict).JSON(map[string]interface{}{"error": fmt.Sprintf("user with email %s already exists", payload.Email)})
	}

	hashedPassword, err := auth.HashedPassword(payload.Password)
	if err != nil {
		logging.ErrorLogger.Println("Password hashing failed:", err)
		return c.Status(http.StatusInternalServerError).JSON(map[string]interface{}{"error": "failed to hash password"})
	}

	newUser := types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	}
	if err := h.store.Createuser(newUser); err != nil {
		logging.ErrorLogger.Println("Failed to create user:", err)
		return c.Status(http.StatusInternalServerError).JSON(map[string]interface{}{"error": "failed to create user"})
	}

	logging.InfoLogger.Println("User registered successfully:", payload.Email)
	return c.Status(http.StatusCreated).JSON(map[string]interface{}{"message": "user registered successfully"})
}
