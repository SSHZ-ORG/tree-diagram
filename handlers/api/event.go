package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/SSHZ-ORG/tree-diagram/apicache"
	"github.com/SSHZ-ORG/tree-diagram/models"
	"github.com/SSHZ-ORG/tree-diagram/pb"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/protobuf/proto"
)

func (t treeDiagramService) RenderEvent(ctx context.Context, req *pb.RenderEventRequest) (*pb.RenderEventResponse, error) {
	eid := req.GetId()

	if fromCache := apicache.GetRenderEvent(ctx, eid); fromCache != nil {
		r := &pb.RenderEventResponse{}
		if err := proto.Unmarshal(fromCache, r); err == nil {
			return r, nil
		} else {
			log.Errorf(ctx, "proto.Unmarshal: %+v", err)
			// Continue below
		}
	}

	res, err := models.PrepareRenderEventResponse(ctx, eid)
	if err == nil {
		if m, err := proto.Marshal(res); err == nil {
			apicache.PutRenderEvent(ctx, eid, m)
		} else {
			panic(err)
		}
	} else {
		log.Errorf(ctx, "models.PrepareRenderEventResponse: %+v", err)
	}

	return res, err
}

func queryEvents(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(r)

	err := r.ParseForm()
	if err != nil {
		log.Errorf(ctx, "r.ParseForm: %v", err)
		http.Error(w, "Malformed Query", http.StatusBadRequest)
		return
	}

	placeID := r.Form.Get("place")
	actorIDs := r.Form["actor"]

	offset := 0
	if arg := r.Form.Get("offset"); arg != "" {
		offset, err = strconv.Atoi(arg)
		if err != nil {
			http.Error(w, "Illegal arg offset", http.StatusBadRequest)
			return
		}
	}

	if arg := r.Form.Get("page"); arg != "" && offset == 0 {
		page, err := strconv.Atoi(arg)
		if err != nil {
			http.Error(w, "Illegal arg page", http.StatusBadRequest)
			return
		}
		offset = (page - 1) * queryPageSize
	}

	events, err := models.QueryEvents(ctx, placeID, actorIDs, queryPageSize, offset)
	if err != nil {
		log.Errorf(ctx, "models.QueryEvents: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fes := make([]models.FrontendEvent, 0)
	for _, e := range events {
		fes = append(fes, e.ToFrontendEvent())
	}
	writeJSON(ctx, w, fes)
}
