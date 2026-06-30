package grpc

import (
	"context"
	"encoding/json"

	"github.com/elzafadli/bookrpc/pb"
	"monitorapp/internal/application/service"
	"monitorapp/internal/domain/activity_log"
	"monitorapp/internal/pkg/validator"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ActivityLogHandler struct {
	pb.UnimplementedActivityLogServiceServer
	Service   service.ActivityLogService `inject:"activityLogService"`
	Validator validator.Validator        `inject:"validator"`
}

func (h *ActivityLogHandler) Create(ctx context.Context, req *pb.CreateActivityLogRequest) (*pb.CreateActivityLogResponse, error) {
	var requestVal interface{}
	if req.GetRequest() != "" {
		if err := json.Unmarshal([]byte(req.GetRequest()), &requestVal); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid request JSON: "+err.Error())
		}
	}

	var responseVal interface{}
	if req.GetResponse() != "" {
		if err := json.Unmarshal([]byte(req.GetResponse()), &responseVal); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid response JSON: "+err.Error())
		}
	}

	payload := &activity_log.CreateActivityLogRequest{
		ObjectName: activity_log.ObjectName(req.GetObjectName()),
		RecordID:   req.GetRecordId(),
		Action:     activity_log.Action(req.GetAction()),
		ChangedBy:  req.GetChangedBy(),
		Request:    requestVal,
		Response:   responseVal,
	}

	if err := h.Validator.Validate(ctx, payload); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := h.Service.Create(ctx, payload)
	if err != nil {
		return nil, mapActivityLogError(err)
	}

	return &pb.CreateActivityLogResponse{
		ActivityLog: convertToPbActivityLog(res),
	}, nil
}

func (h *ActivityLogHandler) Update(ctx context.Context, req *pb.UpdateActivityLogRequest) (*pb.UpdateActivityLogResponse, error) {
	var requestVal interface{}
	if req.GetRequest() != "" {
		if err := json.Unmarshal([]byte(req.GetRequest()), &requestVal); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid request JSON: "+err.Error())
		}
	}

	var responseVal interface{}
	if req.GetResponse() != "" {
		if err := json.Unmarshal([]byte(req.GetResponse()), &responseVal); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid response JSON: "+err.Error())
		}
	}

	payload := &activity_log.UpdateActivityLogRequest{
		ID:         req.GetId(),
		ObjectName: activity_log.ObjectName(req.GetObjectName()),
		RecordID:   req.GetRecordId(),
		Action:     activity_log.Action(req.GetAction()),
		ChangedBy:  req.GetChangedBy(),
		Request:    requestVal,
		Response:   responseVal,
	}

	if err := h.Validator.Validate(ctx, payload); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := h.Service.Update(ctx, payload)
	if err != nil {
		return nil, mapActivityLogError(err)
	}

	return &pb.UpdateActivityLogResponse{
		ActivityLog: convertToPbActivityLog(res),
	}, nil
}

func (h *ActivityLogHandler) Delete(ctx context.Context, req *pb.DeleteActivityLogRequest) (*pb.DeleteActivityLogResponse, error) {
	err := h.Service.Delete(ctx, req.GetId())
	if err != nil {
		return nil, mapActivityLogError(err)
	}

	return &pb.DeleteActivityLogResponse{
		Message: "Activity log successfully deleted",
	}, nil
}

func (h *ActivityLogHandler) Detail(ctx context.Context, req *pb.GetActivityLogRequest) (*pb.GetActivityLogResponse, error) {
	res, err := h.Service.Find(ctx, req.GetId())
	if err != nil {
		return nil, mapActivityLogError(err)
	}

	return &pb.GetActivityLogResponse{
		ActivityLog: convertToPbActivityLog(res),
	}, nil
}

func (h *ActivityLogHandler) List(ctx context.Context, req *pb.ListActivityLogRequest) (*pb.ListActivityLogResponse, error) {
	page := req.GetPage()
	if page < 1 {
		page = 1
	}
	limit := req.GetLimit()
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	filter := map[string]interface{}{
		"limit":       int(limit),
		"offset":      int(offset),
		"object_name": req.GetObjectName(),
		"record_id":   req.GetRecordId(),
		"action":      req.GetAction(),
		"changed_by":  req.GetChangedBy(),
	}

	logs, total, err := h.Service.List(ctx, filter)
	if err != nil {
		return nil, mapActivityLogError(err)
	}

	var pbLogs []*pb.ActivityLog
	for _, l := range logs {
		pbLogs = append(pbLogs, convertToPbActivityLog(l))
	}

	currentPage := int32(page)
	var totalPages int32 = 0
	if limit > 0 {
		totalPages = int32((total + uint64(limit) - 1) / uint64(limit))
	}

	metadata := &pb.ListActivityLogMetadata{
		Total:       total,
		Limit:       limit,
		Offset:      offset,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}

	return &pb.ListActivityLogResponse{
		ActivityLogs: pbLogs,
		Metadata:     metadata,
	}, nil
}

func convertToPbActivityLog(a *activity_log.ActivityLog) *pb.ActivityLog {
	if a == nil {
		return nil
	}

	var reqStr, resStr string
	if a.Request != nil {
		if b, err := json.Marshal(a.Request); err == nil {
			reqStr = string(b)
		}
	}
	if a.Response != nil {
		if b, err := json.Marshal(a.Response); err == nil {
			resStr = string(b)
		}
	}

	return &pb.ActivityLog{
		Id:          a.ID,
		ObjectName:  string(a.ObjectName),
		RecordId:    a.RecordID,
		Action:      string(a.Action),
		ChangedBy:   a.ChangedBy,
		Request:     reqStr,
		Response:    resStr,
		ChangeStamp: timestamppb.New(a.ChangeStamp),
	}
}

func mapActivityLogError(err error) error {
	if err == nil {
		return nil
	}
	switch err {
	case activity_log.ErrDataNotFound:
		return status.Error(codes.NotFound, err.Error())
	case activity_log.ErrDataAlreadyExists:
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
