package response

import (
	"math"
	"net/http"
	"workshop-management/pkg/messages"

	"github.com/google/uuid"
)

// Success is an alias for ApiResponse for swag documentation.
type Success ApiResponse

// Error is an alias for ApiResponse for swag documentation.
type Error ApiResponse

// Pagination is an alias for PaginatedResponse for swag documentation.
type Pagination PaginatedResponse

type Errors struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type ApiResponse struct {
	Id      uuid.UUID   `json:"log_id"`
	Code    int         `json:"code,omitempty"`
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type PaginatedResponse struct {
	LogID       string      `json:"log_id"`
	Code        int         `json:"code"`
	Status      bool        `json:"status"`
	Message     string      `json:"message"`
	TotalData   int         `json:"total_data"`
	TotalPages  int         `json:"total_pages"`
	CurrentPage int         `json:"current_page"`
	NextPage    bool        `json:"next_page"`
	PrevPage    bool        `json:"prev_page"`
	Limit       int         `json:"limit"`
	Data        interface{} `json:"data,omitempty"`
	Error       interface{} `json:"error,omitempty"`
}

func Response(code int, msg string, logId uuid.UUID, data interface{}) *ApiResponse {
	res := new(ApiResponse)
	res.Id = logId
	res.Message = msg
	res.Data = data
	res.Status = code == http.StatusOK || code == http.StatusCreated

	return res
}

func PaginationResponse(code, total, page, perPage int, logId uuid.UUID, data interface{}) *PaginatedResponse {
	res := new(PaginatedResponse)

	// Count total pages
	var totalPages int
	if total > 0 && perPage > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(perPage)))
	} else if total > 0 {
		totalPages = 1
	}

	// Check for next page (hasNext)
	hasNext := page < totalPages

	message := messages.MsgSuccess
	if total == 0 || page > totalPages {
		message = messages.MsgNotFound
	}

	res.LogID = logId.String()
	res.Code = code
	res.Status = code == http.StatusOK || code == http.StatusCreated
	res.Message = message
	res.Data = data
	res.TotalData = total
	res.TotalPages = totalPages
	res.CurrentPage = page
	res.NextPage = hasNext
	res.PrevPage = page > 1
	res.Limit = perPage

	return res
}
