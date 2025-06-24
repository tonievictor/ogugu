package rss

import (
	"net/http/httptest"
	"testing"
)

//	func (rc *RssController) Fetch(w http.ResponseWriter, r *http.Request) {
//		spanctx, span := tracer.Start(r.Context(), "Find RssFeedByID")
//		defer span.End()
//
//		feed, err := rc.rss.Fetch(spanctx)
//		if err != nil {
//			response.Error(w, "Resource not found", http.StatusNotFound, err.Error(), rc.log)
//			return
//		}
//
//		message := "Resources Found"
//		if len(feed) < 1 {
//			message = "No resources found"
//		}
//
//		response.Success(w, message, http.StatusOK, feed, rc.log)
//	}
func TestFetch(t *testing.T) {
	req := httptest.NewRequest("GET", "/")
}
