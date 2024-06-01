package router

import (
	"mongorest/common"
	"mongorest/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddBlogGroup(app *fiber.App) {
	blogGroup := app.Group("/api/blogs")

	blogGroup.Get("/", getBlogs)
	blogGroup.Get("/:id", getBlog)
	blogGroup.Post("/", createBlog)
}

func getBlogs(c *fiber.Ctx) error {
	coll := common.GetDBCollection("blogs")

	// find all
	blogs := make([]models.Blog, 0)
	cursor, err := coll.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// iterate over the cursor
	for cursor.Next(c.Context()) {
		blog := models.Blog{}
		err := cursor.Decode(&blog)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		blogs = append(blogs, blog)
	}

	return c.Status(200).JSON(fiber.Map{"data": blogs})
}

func getBlog(c *fiber.Ctx) error {
	coll := common.GetDBCollection("blogs")

	// find the blog
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "id is required",
		})
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	blog := models.Blog{}

	err = coll.FindOne(c.Context(), bson.M{"_id": objectId}).Decode(&blog)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(blog)
}

type createDTO struct {
	Title   string `json:"title" bson:"title"`
	Content string `json:"content" bson:"content"`
}

func createBlog(c *fiber.Ctx) error {

	b := new(createDTO)
	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	blog := models.Blog{
		Title:     b.Title,
		Content:   b.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	coll := common.GetDBCollection("blogs")
	res, err := coll.InsertOne(c.Context(), blog)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create a blog",
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"result": res,
	})
}