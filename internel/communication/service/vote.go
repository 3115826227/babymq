package service

import (
	"context"
	pbservices "github.com/3115826227/babymq/internel/communication/service/pbservice"
)

type VoteService struct {
}

func (voteService *VoteService) Request(context context.Context, request *pbservices.VoteRequest) (response *pbservices.VoteResponse, err error) {
	response = &pbservices.VoteResponse{
		MessageType: request.MessageType,
		Term:        request.Term,
		Version:     request.Version,
	}
	return
}
