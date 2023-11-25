package controller

import (
	"log"

	"github.com/Mau005/KraynoSerer/models"
)

type ToolsController struct{}

func (tc *ToolsController) SharedLoot(content string, sc *models.StructModel) error {
	log.Println(content)

	return nil
}
