package check

import (
	"net/http"
	"techtraits.com/klaxon/router"
)

func init() {
	router.Register("/rest/internal/check/{project_id}", router.POST, nil, nil, getTick)
}

func getTick(request router.Request) (int, []byte) {

	return http.StatusOK, nil

}
