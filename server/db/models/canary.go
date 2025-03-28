package models

/*
	Sliver Implant Framework
	Copyright (C) 2020  Bishop Fox

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"time"

	"github.com/gsmith257-cyber/better-sliver-package/protobuf/clientpb"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// DNSCanary - Colletions of content to serve from HTTP(S)
type DNSCanary struct {
	ID        uuid.UUID `gorm:"primaryKey;->;<-:create;type:uuid;"`
	CreatedAt time.Time `gorm:"->;<-:create;"`

	ImplantName   string
	Domain        string
	Triggered     bool
	FirstTrigger  time.Time
	LatestTrigger time.Time
	Count         uint32
}

// BeforeCreate - GORM hook
func (c *DNSCanary) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	c.CreatedAt = time.Now()
	return nil
}

// ToProtobuf - Converts to protobuf object
func (c *DNSCanary) ToProtobuf() *clientpb.DNSCanary {
	return &clientpb.DNSCanary{
		ImplantName:    c.ImplantName,
		Domain:         c.Domain,
		Triggered:      c.Triggered,
		FirstTriggered: c.FirstTrigger.Format(time.RFC1123),
		LatestTrigger:  c.LatestTrigger.Format(time.RFC1123),
		Count:          c.Count,
	}
}

// convert from protobuf
func DNSCanaryFromProtobuf(m *clientpb.DNSCanary) DNSCanary {
	uuid, _ := uuid.FromString(m.ID)
	firstTrigger, _ := time.Parse(time.RFC1123, m.FirstTriggered)
	latestTrigger, _ := time.Parse(time.RFC1123, m.LatestTrigger)
	return DNSCanary{
		ID:            uuid,
		ImplantName:   m.ImplantName,
		Domain:        m.Domain,
		Triggered:     m.Triggered,
		FirstTrigger:  firstTrigger,
		LatestTrigger: latestTrigger,
		Count:         m.Count,
	}
}
