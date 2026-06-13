package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AncientText struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null;index" json:"title"`
	Author      string         `gorm:"type:varchar(255)" json:"author"`
	Dynasty     string         `gorm:"type:varchar(100)" json:"dynasty"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	LayoutConfig   *LayoutConfig  `gorm:"foreignKey:TextID" json:"layout_config,omitempty"`
	Annotations    []Annotation   `gorm:"foreignKey:TextID" json:"annotations,omitempty"`
}

type LayoutConfig struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	TextID            uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"text_id"`
	Columns           int            `gorm:"default:12" json:"columns"`
	CharsPerColumn    int            `gorm:"default:20" json:"chars_per_column"`
	ColumnGap         int            `gorm:"default:20" json:"column_gap"`
	FontSize          int            `gorm:"default:18" json:"font_size"`
	FontFamily        string         `gorm:"type:varchar(255);default:'SimSun, serif'" json:"font_family"`
	LineHeight        float64        `gorm:"default:1.8" json:"line_height"`
	CharSpacing       float64        `gorm:"default:0" json:"char_spacing"`
	TextColor         string         `gorm:"type:varchar(20);default:'#1a1a1a'" json:"text_color"`
	BackgroundColor   string         `gorm:"type:varchar(20);default:'#f5f0e8'" json:"background_color"`
	ShowBorder        bool           `gorm:"default:true" json:"show_border"`
	BorderStyle       datatypes.JSON `gorm:"type:jsonb" json:"border_style"`
	Padding           datatypes.JSON `gorm:"type:jsonb" json:"padding"`
	CustomCSS         string         `gorm:"type:text" json:"custom_css"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

type LayoutPreset struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID            string         `gorm:"type:varchar(255);not null;index" json:"user_id"`
	Name              string         `gorm:"type:varchar(255);not null" json:"name"`
	Description       string         `gorm:"type:text" json:"description"`
	CharsPerColumn    int            `gorm:"default:20" json:"chars_per_column"`
	ColumnGap         int            `gorm:"default:20" json:"column_gap"`
	FontSize          int            `gorm:"default:18" json:"font_size"`
	FontFamily        string         `gorm:"type:varchar(255);default:'SimSun, serif'" json:"font_family"`
	LineHeight        float64        `gorm:"default:1.8" json:"line_height"`
	CharSpacing       float64        `gorm:"default:0" json:"char_spacing"`
	TextColor         string         `gorm:"type:varchar(20);default:'#1a1a1a'" json:"text_color"`
	BackgroundColor   string         `gorm:"type:varchar(20);default:'#f5f0e8'" json:"background_color"`
	ShowBorder        bool           `gorm:"default:true" json:"show_border"`
	IsDefault         bool           `gorm:"default:false;index" json:"is_default"`
	CustomCSS         string         `gorm:"type:text" json:"custom_css"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *LayoutPreset) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type Annotation struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	TextID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"text_id"`
	UserID    string         `gorm:"type:varchar(255);index" json:"user_id"`
	StartPos  int            `gorm:"not null" json:"start_pos"`
	EndPos    int            `gorm:"not null" json:"end_pos"`
	Type      string         `gorm:"type:varchar(50);not null" json:"type"`
	Content   string         `gorm:"type:text" json:"content"`
	Color     string         `gorm:"type:varchar(20)" json:"color"`
	Metadata  datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
	Status    string         `gorm:"type:varchar(20);default:'draft'" json:"status"`
	Version   int            `gorm:"default:1" json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type AnnotationType string

const (
	AnnotationTypePeriod    AnnotationType = "period"
	AnnotationTypeComma     AnnotationType = "comma"
	AnnotationTypeQuestion  AnnotationType = "question"
	AnnotationTypeExclaim   AnnotationType = "exclaim"
	AnnotationTypeComment   AnnotationType = "comment"
	AnnotationTypeHighlight AnnotationType = "highlight"
)

func (t *AncientText) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

func (l *LayoutConfig) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

func (a *Annotation) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
