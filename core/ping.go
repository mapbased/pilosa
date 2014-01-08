package core

import (
	"encoding/gob"
	"pilosa/db"
	"time"

	"tux21b.org/v1/gocql/uuid"
)

type PingRequest struct {
	Id     *uuid.UUID
	Source *uuid.UUID
}

type PongRequest struct {
	Id *uuid.UUID
}

func (self PongRequest) ResultId() *uuid.UUID {
	return self.Id
}
func (self PongRequest) ResultData() interface{} {
	return self.Id
}

func init() {
	gob.Register(PingRequest{})
	gob.Register(PongRequest{})
}

func (self *Service) Ping(process_id *uuid.UUID) (*time.Duration, error) {
	id := uuid.RandomUUID()
	ping := db.Message{Data: PingRequest{Id: &id, Source: self.Id}}
	start := time.Now()
	self.Transport.Send(&ping, process_id)
	_, err := self.Hold.Get(&id, 60)
	if err != nil {
		return nil, err
	}
	end := time.Now()
	dur := end.Sub(start)
	return &dur, nil
}
