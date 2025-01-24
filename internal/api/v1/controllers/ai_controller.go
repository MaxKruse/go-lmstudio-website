package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/maxkruse/go-lmstudio-website/internal/llm_integration"
	requestdtos "github.com/maxkruse/go-lmstudio-website/internal/models/dtos/request_dtos"
)

// @Description		Endpoint to ask AI for help
// @Accept			json
// @Produce			json
// @Success			200	{object}	dtos.CompletionResult
// @Failure			400	{object}	error
// @Failure			500	{object}	error
// @Param			CompletionRequest	body	string	true "The Chat Completion Request including ChatCompletionNewParams and a Prompt"
// @Router			/ai/completion	[post]
func AiChatCompletion(e echo.Context) error {

	// step 1: get the prompt from the request body
	var request requestdtos.CompletionRequest

	if err := e.Bind(&request); err != nil {
		return e.JSON(http.StatusBadRequest, err)
	}

	// step 2: make the ai client
	aiClient := llm_integration.NewClient()

	// step 3: get the completion
	completionResult, err := aiClient.GetCompletion(e.Request().Context(), request.Prompt, &request.ParamsUsed)

	if err != nil {
		return e.JSON(http.StatusInternalServerError, err)
	}

	return e.JSON(http.StatusOK, completionResult)
}
