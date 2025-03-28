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
	"os"
	"path/filepath"
	"time"

	"github.com/gsmith257-cyber/better-sliver-package/protobuf/clientpb"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Website - Colletions of content to serve from HTTP(S)
type Website struct {
	ID        uuid.UUID `gorm:"primaryKey;->;<-:create;type:uuid;"`
	CreatedAt time.Time `gorm:"->;<-:create;"`

	Name string `gorm:"unique;"` // Website Name

	WebContents []WebContent
}

// BeforeCreate - GORM hook
func (w *Website) BeforeCreate(tx *gorm.DB) (err error) {
	w.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	w.CreatedAt = time.Now()
	return nil
}

// ToProtobuf - Converts to protobuf object
func (w *Website) ToProtobuf(webContentDir string) *clientpb.Website {
	WebContents := map[string]*clientpb.WebContent{}
	for _, webcontent := range w.WebContents {
		contents, _ := os.ReadFile(filepath.Join(webContentDir, webcontent.Path))
		WebContents[webcontent.ID.String()] = webcontent.ToProtobuf(&contents)
	}
	return &clientpb.Website{
		ID:       w.ID.String(),
		Name:     w.Name,
		Contents: WebContents,
	}
}

// WebContent - One piece of content mapped to a path
type WebContent struct {
	ID        uuid.UUID `gorm:"primaryKey;->;<-:create;type:uuid;"`
	WebsiteID uuid.UUID `gorm:"type:uuid;"`

	Path        string `gorm:"primaryKey"`
	Size        uint64
	ContentType string
}

// BeforeCreate - GORM hook to automatically set values
func (wc *WebContent) BeforeCreate(tx *gorm.DB) (err error) {
	wc.ID, err = uuid.NewV4()
	return err
}

// ToProtobuf - Converts to protobuf object
func (wc *WebContent) ToProtobuf(content *[]byte) *clientpb.WebContent {
	return &clientpb.WebContent{
		ID:          wc.ID.String(),
		WebsiteID:   wc.WebsiteID.String(),
		Path:        wc.Path,
		Size:        uint64(wc.Size),
		ContentType: wc.ContentType,
		Content:     *content,
	}
}

func WebContentFromProtobuf(pbWebContent *clientpb.WebContent) WebContent {
	siteUUID, _ := uuid.FromString(pbWebContent.ID)
	websiteUUID, _ := uuid.FromString(pbWebContent.WebsiteID)

	return WebContent{
		ID:          siteUUID,
		WebsiteID:   websiteUUID,
		Path:        pbWebContent.Path,
		Size:        pbWebContent.Size,
		ContentType: pbWebContent.ContentType,
	}
}
