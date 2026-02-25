package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"nexus-super-node-v3/internal/adapters/ai"
	"nexus-super-node-v3/internal/adapters/openclaw"
	"nexus-super-node-v3/internal/core/domain"
	"nexus-super-node-v3/internal/core/services/mcp"
	"nexus-super-node-v3/internal/core/services/voltagent"
	"nexus-super-node-v3/internal/ports"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spruceid/siwe-go"
)

// EchoGateway is the gateway for the application using Echo framework.
type EchoGateway struct {
	echo       *echo.Echo
	userRepo   ports.UserRepository
	mcpSvc     *mcp.MCPService
	voltSvc    *voltagent.VoltAgentService
	chatSvc    ports.ChatService
	socialSvc  ports.SocialService
	financeSvc ports.FinanceService
	agentSvc   ports.AgentService
	redpanda   ports.EventProducer
	claw       *openclaw.Client
}

// NewEchoGateway creates a new EchoGateway.
func NewEchoGateway(userRepo ports.UserRepository, mcpSvc *mcp.MCPService, voltSvc *voltagent.VoltAgentService, chatSvc ports.ChatService, socialSvc ports.SocialService, financeSvc ports.FinanceService, agentSvc ports.AgentService, rp ports.EventProducer, claw *openclaw.Client) *EchoGateway {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS Configuration
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodHead, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		ExposeHeaders:    []string{echo.HeaderContentLength},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	return &EchoGateway{
		echo:       e,
		userRepo:   userRepo,
		mcpSvc:     mcpSvc,
		voltSvc:    voltSvc,
		chatSvc:    chatSvc,
		socialSvc:  socialSvc,
		financeSvc: financeSvc,
		agentSvc:   agentSvc,
		redpanda:   rp,
		claw:       claw,
	}
}

// RegisterWebSocketRoutes registers websocket endpoints
func (g *EchoGateway) RegisterWebSocketRoutes(wsHandler *WebSocketHandler) {
	g.echo.GET("/ws/market", wsHandler.HandleConnection)
}

// Start starts the gateway.
func (g *EchoGateway) Start(ctx context.Context) error {
	g.echo.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	g.setupAuthRoutes()
	g.setupToolRoutes()
	g.setupMCPRoutes()
	g.setupVoltAgentRoutes()
	g.setupHasuraRoutes()
	g.setupChatRoutes()
	g.setupSocialRoutes()
	g.setupFinanceRoutes()
	g.setupOpenClawRoutes()
	g.setupAgentRoutes()
	g.setupWorkflowRoutes()
	g.setupBuilderRoutes()

	// Use 3000 as default load balancer/gateway port
	return g.echo.Start(":3000")
}

// In-memory store for workflows (for demonstration purposes)
var inMemoryWorkflows = []map[string]interface{}{
	{
		"id":          "workflow-1",
		"name":        "Rivet Agent Flow",
		"description": "Main logic for agent decision making",
		"status":      "Active",
		"runs":        42,
		"lastRun":     "10 mins ago",
	},
}

func (g *EchoGateway) setupHasuraRoutes() {
	g.echo.POST("/hasura/events", func(c echo.Context) error {
		// Generic endpoint for Hasura Event Triggers
		var body map[string]interface{}
		if err := c.Bind(&body); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		// In a real implementation, we would inspect the event payload
		// For now, we just marshal it and send it to Redpanda
		payload, _ := json.Marshal(body)

		err := g.redpanda.Produce(c.Request().Context(), []byte("hasura-event"), payload)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "queued"})
	})
}

func (g *EchoGateway) setupAuthRoutes() {
	auth := g.echo.Group("/auth")

	auth.GET("/nonce", func(c echo.Context) error {
		nonce := siwe.GenerateNonce()

		// In a real app, store this nonce in a session or Redis with expiration
		// For now, we set it in a cookie for the client to sign (simplified flow)
		cookie := new(http.Cookie)
		cookie.Name = "siwe-nonce"
		cookie.Value = nonce
		cookie.Expires = time.Now().Add(5 * time.Minute)
		cookie.HttpOnly = true
		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, map[string]string{"nonce": nonce})
	})

	auth.POST("/verify", func(c echo.Context) error {
		var body struct {
			Message   string `json:"message"`
			Signature string `json:"signature"`
		}

		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, "Invalid request body")
		}

		// Parse SIWE message
		message, err := siwe.ParseMessage(body.Message)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Invalid SIWE message: %s", err))
		}

		// Verify signature
		// Note: In production, verify the nonce matches what we issued
		publicKey, err := message.Verify(body.Signature, nil, nil, nil)
		if err != nil {
			return c.String(http.StatusUnauthorized, fmt.Sprintf("Invalid signature: %s", err))
		}

		// Calculate Ethereum address from public key
		address := crypto.PubkeyToAddress(*publicKey).Hex()
		fmt.Println("Verified SIWE message from:", address)

		// Check if user exists using standard context
		user, err := g.userRepo.GetUserByAddress(c.Request().Context(), address)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		if user == nil {
			// User does not exist, create a new one
			newUser := &domain.User{
				Address: address,
			}
			user, err = g.userRepo.CreateUser(c.Request().Context(), newUser)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
		}

		return c.JSON(http.StatusOK, user)
	})
}

func (g *EchoGateway) setupToolRoutes() {
	tools := g.echo.Group("/tools")

	tools.POST("/execute", func(c echo.Context) error {
		var body struct {
			ToolBeltName string        `json:"tool_belt_name"`
			ToolName     string        `json:"tool_name"`
			Args         []interface{} `json:"args"`
		}

		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		result, err := g.mcpSvc.ExecuteTool(body.ToolBeltName, body.ToolName, body.Args...)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, result)
	})
}

func (g *EchoGateway) setupMCPRoutes() {
	mcpGroup := g.echo.Group("/mcp")

	// Register a new dynamic MCP server
	mcpGroup.POST("/register", func(c echo.Context) error {
		var config mcp.MCPServerConfig
		if err := c.Bind(&config); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		if err := g.mcpSvc.RegisterDynamicServer(config); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "registered", "id": config.ID})
	})

	// List all dynamic MCP servers
	mcpGroup.GET("/servers", func(c echo.Context) error {
		servers := g.mcpSvc.ListDynamicServers()
		return c.JSON(http.StatusOK, servers)
	})
}

func (g *EchoGateway) setupVoltAgentRoutes() {
	volt := g.echo.Group("/voltagent")

	// Manifest for VoltAgent (Tools definition)
	volt.GET("/manifest", func(c echo.Context) error {
		manifest, err := g.voltSvc.GetManifest()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, manifest)
	})

	// Execution endpoint for VoltAgent tools
	volt.POST("/execute", func(c echo.Context) error {
		var body struct {
			ToolID string                 `json:"tool_id"`
			Args   map[string]interface{} `json:"args"`
		}
		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		result, err := g.voltSvc.ExecuteTool(body.ToolID, body.Args)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, result)
	})
}

func (g *EchoGateway) setupChatRoutes() {
	api := g.echo.Group("/api")

	// --- VoltAgent Chat (Streaming) ---
	api.POST("/chat", func(c echo.Context) error {
		// Parse request
		var body ai.ChatRequest

		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		// Set headers for streaming
		c.Response().Header().Set(echo.HeaderContentType, "text/plain; charset=utf-8")
		c.Response().Header().Set("Transfer-Encoding", "chunked")
		c.Response().Header().Set("X-Content-Type-Options", "nosniff")

		// Stream via VoltAgent
		ctx := c.Request().Context()
		streamChan, errChan := g.voltSvc.StreamChat(ctx, body.Messages)

		for {
			select {
			case chunk, ok := <-streamChan:
				if !ok {
					return nil // Stream closed
				}
				if _, err := c.Response().Write([]byte(chunk)); err != nil {
					return err
				}
				c.Response().Flush()
			case err := <-errChan:
				if err != nil {
					// Log error but maybe don't break stream if it's partial?
					// For now, return error
					return err
				}
			case <-ctx.Done():
				return nil
			}
		}
	})

	// --- Unified Chat System (Matrix + LiveKit) ---
	channels := api.Group("/channels")

	// Create Channel
	channels.POST("", func(c echo.Context) error {
		var body struct {
			Name     string `json:"name"`
			Alias    string `json:"alias"`
			IsPublic bool   `json:"is_public"`
		}
		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		channelID, err := g.chatSvc.CreateChannel(c.Request().Context(), body.Name, body.Alias, body.IsPublic)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, map[string]string{"channel_id": channelID})
	})

	// Join Channel (Invite User)
	channels.POST("/:id/join", func(c echo.Context) error {
		channelID := c.Param("id")
		var body struct {
			UserID string `json:"user_id"`
		}
		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if err := g.chatSvc.JoinChannel(c.Request().Context(), channelID, body.UserID); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "joined"})
	})

	// Send Message
	channels.POST("/:id/messages", func(c echo.Context) error {
		channelID := c.Param("id")
		var body struct {
			Message string `json:"message"`
		}
		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if err := g.chatSvc.SendMessage(c.Request().Context(), channelID, body.Message); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "sent"})
	})

	// Start Call in Channel
	channels.POST("/:id/call", func(c echo.Context) error {
		channelID := c.Param("id")

		callURL, err := g.chatSvc.StartCallInChannel(c.Request().Context(), channelID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"call_url": callURL})
	})
}

func (g *EchoGateway) setupSocialRoutes() {
	social := g.echo.Group("/api/social")

	// Create Post
	social.POST("/posts", func(c echo.Context) error {
		var body struct {
			AuthorID  string   `json:"author_id"`
			Content   string   `json:"content"`
			MediaURLs []string `json:"media_urls"`
		}
		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		post, err := g.socialSvc.CreatePost(c.Request().Context(), body.AuthorID, body.Content, body.MediaURLs)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, post)
	})

	// Get Feed
	social.GET("/feed", func(c echo.Context) error {
		var filter domain.FeedFilter
		// Bind query params manually or via struct if echo supports it well for GET
		// Simplified binding:
		// filter.Tags = c.QueryParams()["tags"]

		posts, err := g.socialSvc.GetFeed(c.Request().Context(), filter)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, posts)
	})

	// Like Post
	social.POST("/posts/:id/like", func(c echo.Context) error {
		postID := c.Param("id")
		var body struct {
			UserID string `json:"user_id"`
		}
		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if err := g.socialSvc.LikePost(c.Request().Context(), postID, body.UserID); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "liked"})
	})

	// Add Comment
	social.POST("/posts/:id/comments", func(c echo.Context) error {
		postID := c.Param("id")
		var body struct {
			UserID  string `json:"user_id"`
			Content string `json:"content"`
		}
		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		comment, err := g.socialSvc.AddComment(c.Request().Context(), postID, body.UserID, body.Content)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, comment)
	})
}

func (g *EchoGateway) setupFinanceRoutes() {
	finance := g.echo.Group("/api/finance")

	// Create Asset
	finance.POST("/assets", func(c echo.Context) error {
		var body struct {
			Name          string  `json:"name"`
			Symbol        string  `json:"symbol"`
			Type          string  `json:"type"`
			InitialSupply float64 `json:"initial_supply"`
		}
		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		asset, err := g.financeSvc.CreateAsset(c.Request().Context(), body.Name, body.Symbol, domain.AssetType(body.Type), body.InitialSupply)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, asset)
	})

	// Get Asset
	finance.GET("/assets/:id", func(c echo.Context) error {
		id := c.Param("id")
		asset, err := g.financeSvc.GetAsset(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, asset)
	})

	// Get Balance
	finance.GET("/balances/:userID/:assetID", func(c echo.Context) error {
		userID := c.Param("userID")
		assetID := c.Param("assetID")
		balance, err := g.financeSvc.GetUserBalance(c.Request().Context(), userID, assetID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, balance)
	})

	// Transfer
	finance.POST("/transfer", func(c echo.Context) error {
		var body struct {
			FromUserID string  `json:"from_user_id"`
			ToUserID   string  `json:"to_user_id"`
			AssetID    string  `json:"asset_id"`
			Amount     float64 `json:"amount"`
		}
		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if err := g.financeSvc.Transfer(c.Request().Context(), body.FromUserID, body.ToUserID, body.AssetID, body.Amount); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "success"})
	})

	// Request Loan
	finance.POST("/loans", func(c echo.Context) error {
		var body struct {
			BorrowerID   string  `json:"borrower_id"`
			CollateralID string  `json:"collateral_id"`
			Amount       float64 `json:"amount"`
		}
		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		loan, err := g.financeSvc.RequestLoan(c.Request().Context(), body.BorrowerID, body.CollateralID, body.Amount)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, loan)
	})

	// Distribute Reward
	finance.POST("/rewards", func(c echo.Context) error {
		var body struct {
			UserID  string  `json:"user_id"`
			AssetID string  `json:"asset_id"`
			Amount  float64 `json:"amount"`
			Reason  string  `json:"reason"`
		}
		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if err := g.financeSvc.DistributeReward(c.Request().Context(), body.UserID, body.AssetID, body.Amount, body.Reason); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "reward_distributed"})
	})
}

func (g *EchoGateway) setupOpenClawRoutes() {
	claw := g.echo.Group("/api/openclaw")

	// Send Message
	claw.POST("/send", func(c echo.Context) error {
		var body struct {
			To      string `json:"to"`
			Message string `json:"message"`
			Channel string `json:"channel"`
		}
		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if err := g.claw.SendMessage(body.To, body.Message, body.Channel); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "sent"})
	})
}

func (g *EchoGateway) setupAgentRoutes() {
	agents := g.echo.Group("/agents")

	// Create Agent
	agents.POST("", func(c echo.Context) error {
		var agent domain.Agent
		if err := c.Bind(&agent); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		// TODO: Extract owner ID from context/auth
		// agent.OwnerID = c.Get("user_id").(string)
		if agent.OwnerID == "" {
			agent.OwnerID = "default-owner" // Fallback for now
		}

		if err := g.agentSvc.CreateAgent(c.Request().Context(), &agent); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, agent)
	})

	// List Agents
	agents.GET("", func(c echo.Context) error {
		// TODO: Extract owner ID from context/auth
		ownerID := "default-owner"

		agentsList, err := g.agentSvc.ListAgents(c.Request().Context(), ownerID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, agentsList)
	})

	// Get Agent by ID
	agents.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")
		agent, err := g.agentSvc.GetAgent(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Agent not found"})
		}
		return c.JSON(http.StatusOK, agent)
	})

	// Update Agent
	agents.PUT("/:id", func(c echo.Context) error {
		id := c.Param("id")
		var agent domain.Agent
		if err := c.Bind(&agent); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}
		agent.ID = id // Ensure ID matches URL

		if err := g.agentSvc.UpdateAgent(c.Request().Context(), &agent); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, agent)
	})

	// Delete Agent
	agents.DELETE("/:id", func(c echo.Context) error {
		id := c.Param("id")
		if err := g.agentSvc.DeleteAgent(c.Request().Context(), id); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.NoContent(http.StatusNoContent)
	})

	// Deploy Agent
	agents.POST("/:id/deploy", func(c echo.Context) error {
		id := c.Param("id")
		if err := g.agentSvc.DeployAgent(c.Request().Context(), id); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "deploying"})
	})

	// Pause Agent
	agents.POST("/:id/pause", func(c echo.Context) error {
		id := c.Param("id")
		if err := g.agentSvc.PauseAgent(c.Request().Context(), id); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "paused"})
	})
}

func (g *EchoGateway) setupWorkflowRoutes() {
	workflows := g.echo.Group("/workflows")

	workflows.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, inMemoryWorkflows)
	})

	workflows.POST("", func(c echo.Context) error {
		var body struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		newWorkflow := map[string]interface{}{
			"id":          fmt.Sprintf("workflow-%d", len(inMemoryWorkflows)+1),
			"name":        body.Name,
			"description": body.Description,
			"status":      "Active",
			"runs":        0,
			"lastRun":     "Just now",
		}

		inMemoryWorkflows = append(inMemoryWorkflows, newWorkflow)
		return c.JSON(http.StatusCreated, newWorkflow)
	})

	workflows.POST("/:id/run", func(c echo.Context) error {
		id := c.Param("id")

		// Find workflow to ensure it exists
		var wf map[string]interface{}
		for _, w := range inMemoryWorkflows {
			if w["id"] == id {
				wf = w
				break
			}
		}

		if wf == nil {
			return c.String(http.StatusNotFound, "Workflow not found")
		}

		// Create a task payload
		task := map[string]interface{}{
			"type":        "workflow_execution",
			"workflow_id": id,
			"timestamp":   time.Now(),
		}

		payload, err := json.Marshal(task)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		// Send to Redpanda to trigger Temporal/Rivet
		err = g.redpanda.Produce(c.Request().Context(), []byte("workflow-run"), payload)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "started", "workflow_id": id})
	})

	workflows.POST("/run", func(c echo.Context) error {
		var body struct {
			GraphID string                 `json:"graph_id"`
			Inputs  map[string]interface{} `json:"inputs"`
		}

		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		task := map[string]interface{}{
			"type":     "workflow",
			"graph_id": body.GraphID,
			"inputs":   body.Inputs,
		}

		payload, err := json.Marshal(task)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		err = g.redpanda.Produce(c.Request().Context(), []byte("workflow-run"), payload)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "queued", "graph_id": body.GraphID})
	})
}

// setupBuilderRoutes sets up routes for the AI Agent Builder
func (g *EchoGateway) setupBuilderRoutes() {
	builder := g.echo.Group("/builder")

	builder.POST("/generate", func(c echo.Context) error {
		var body struct {
			Prompt    string `json:"prompt"`
			Model     string `json:"model"`
			Framework string `json:"framework"`
			Theme     string `json:"theme"`
		}
		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		// Simulate AI Generation delay
		time.Sleep(2 * time.Second)

		// Mock generated code response
		response := map[string]interface{}{
			"status":  "success",
			"message": "Agent generated successfully",
			"files": map[string]string{
				"agent.ts":    fmt.Sprintf("// Generated Agent based on: %s\n// Model: %s\n\nimport { Agent } from '@nexus/core';\n\nexport class CustomAgent extends Agent {\n  async run() {\n    console.log('Executing logic for: %s');\n  }\n}", body.Prompt, body.Model, body.Prompt),
				"config.json": fmt.Sprintf("{\n  \"framework\": \"%s\",\n  \"theme\": \"%s\"\n}", body.Framework, body.Theme),
			},
		}

		return c.JSON(http.StatusOK, response)
	})
}
