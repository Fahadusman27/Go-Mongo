package service

import (
	"fmt"
	"strconv"
	"strings"
	"tugas/domain/model"
	"tugas/domain/repository"

	"github.com/gofiber/fiber/v2"
)

func GetUsersService(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "_id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	offset := (page - 1) * limit

	sortByWhitelist := map[string]bool{"_id": true, "username": true, "email": true, "created_at": true}
	if !sortByWhitelist[sortBy] {
		sortBy = "_id"
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	users, err := repository.GetUsersRepo(search, sortBy, order, limit, offset)
	if err != nil {
		fmt.Println("GetUsersRepo error:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}

	total, err := repository.CountUsersRepo(search)
	if err != nil {
		fmt.Println("CountUsersRepo error:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count users"})
	}

	response := model.UserResponse{
		Data: users,
		MetaInfo: model.MetaInfo{
			CurrentPage: page,
			Limit:       limit,
			Total:       total,
			Pages:       (total + limit - 1) / limit,
			SortBy:      sortBy,
			Order:       order,
			Search:      search,
		},
	}
	return c.JSON(response)
}