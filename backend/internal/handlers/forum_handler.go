package handlers

import (
	"net/http"
	"strconv"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ForumHandler struct {
	svc *services.ForumService
}

func NewForumHandler(svc *services.ForumService) *ForumHandler {
	return &ForumHandler{svc: svc}
}

// ListForums GET /courses/:id/forums
func (h *ForumHandler) ListForums(c *gin.Context) {
	courseID := c.Param("id")
	forums, err := h.svc.ListForums(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, forums)
}

// CreateForum POST /courses/:id/forums
func (h *ForumHandler) CreateForum(c *gin.Context) {
	courseID := c.Param("id")
	var f models.Forum
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	f.CourseOfferingID = courseID
	
	created, err := h.svc.CreateForum(c.Request.Context(), &f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// ListTopics GET /forums/:id/topics
func (h *ForumHandler) ListTopics(c *gin.Context) {
	forumID := c.Param("id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	topics, err := h.svc.ListTopics(c.Request.Context(), forumID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, topics)
}

// CreateTopic POST /forums/:id/topics
func (h *ForumHandler) CreateTopic(c *gin.Context) {
	forumID := c.Param("id")
	var t models.Topic
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t.ForumID = forumID
	t.AuthorID = userIDFromClaims(c)

	created, err := h.svc.CreateTopic(c.Request.Context(), &t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// GetTopic GET /topics/:id
func (h *ForumHandler) GetTopic(c *gin.Context) {
	topicID := c.Param("id")
	topic, err := h.svc.GetTopic(c.Request.Context(), topicID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}
	
	// Also fetch posts
	posts, err := h.svc.ListPosts(c.Request.Context(), topicID)
	if err != nil {
		// Log error but maybe return topic anyway?
	}
	
	c.JSON(http.StatusOK, gin.H{
		"topic": topic,
		"posts": posts,
	})
}

// CreatePost POST /topics/:id/posts
func (h *ForumHandler) CreatePost(c *gin.Context) {
	topicID := c.Param("id")
	var p models.Post
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.TopicID = topicID
	p.AuthorID = userIDFromClaims(c)

	created, err := h.svc.CreatePost(c.Request.Context(), &p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}
