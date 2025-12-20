package api

import (
	"log"
	"net/http"
	"strconv"
	"web3-go-demo/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Server API 服务器
type Server struct {
	router   *gin.Engine
	database *db.Database
}

// NewServer 创建新的 API 服务器
func NewServer(database *db.Database) *Server {
	router := gin.Default()

	server := &Server{
		router:   router,
		database: database,
	}

	// 注册路由
	server.setupRoutes()

	return server
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// 健康检查
	s.router.GET("/health", s.healthCheck)

	// API v1 路由组
	v1 := s.router.Group("/api/v1")
	{
		// 余额查询
		v1.GET("/balance/:address", s.getBalance)

		// 交易查询
		v1.GET("/transactions/:address", s.getTransactions)
		v1.GET("/transaction/:txhash", s.getTransaction)
	}
}

// healthCheck 健康检查
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"message": "ERC20 Tracker is running",
	})
}

// getBalance 获取地址余额
// GET /api/v1/balance/:address
func (s *Server) getBalance(c *gin.Context) {
	address := c.Param("address")

	balance, err := s.database.GetBalance(address)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Address not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"address": balance.Address,
		"balance": balance.Balance,
		"updated_at": balance.UpdatedAt,
	})
}

// getTransactions 获取地址的交易记录
// GET /api/v1/transactions/:address?limit=10&offset=0
func (s *Server) getTransactions(c *gin.Context) {
	address := c.Param("address")

	// 获取分页参数
	limit := 10
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	offset := 0
	if o := c.Query("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = val
		}
	}

	transactions, err := s.database.GetTransactionsByAddress(address, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"count": len(transactions),
		"limit": limit,
		"offset": offset,
		"transactions": transactions,
	})
}

// getTransaction 根据交易哈希获取交易详情
// GET /api/v1/transaction/:txhash
func (s *Server) getTransaction(c *gin.Context) {
	txHash := c.Param("txhash")

	transaction, err := s.database.GetTransactionByHash(txHash)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Transaction not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// Start 启动 API 服务器
func (s *Server) Start(port string) error {
	log.Printf("API 服务器启动在端口: %s", port)
	return s.router.Run(":" + port)
}
