package core

import "github.com/google/uuid"

type DataObject[T any] struct {
	_ID   string
	_From string // From id
	Label string `json:"label"`
	From  string `json:"from"`
	Data  T      `json:"data"`
	Error error  `json:"error"`
}

func (d *DataObject[T]) GetId() string {
	return d._ID
}

func (d *DataObject[T]) GetFrom() string {
	return d._From
}

func (d *DataObject[T]) SetFrom(from string) {
	d._From = from
}

func (d *DataObject[T]) SetId(id string) {
	d._ID = id
}

func (d *DataObject[T]) GetError() error {
	return d.Error
}

func (d *DataObject[T]) GetData() T {
	return d.Data
}

func (d *DataObject[T]) SetData(data T) {
	d.Data = data
}

func (d *DataObject[T]) GetLabel() string {
	return d.Label
}

func (d *DataObject[T]) SetError(err error) {
	d.Error = err
}

func (d *DataObject[T]) SetLabel(label string) {
	d.Label = label
}

func (d *DataObject[T]) IsError() bool {
	return d.Error != nil
}

type DataObjectInterface[T any] interface {
	GetFrom() string
	SetFrom(string)
	GetId() string
	SetId(string)
	GetError() error
	SetError(error)
	GetData() T
	SetData(T)
	GetLabel() string
	SetLabel(string)
	IsError() bool
}

func NewData[T any](label string, data T) DataObjectInterface[T] {
	return &DataObject[T]{
		_ID:   uuid.New().String(),
		Label: label,
		Data:  data,
	}
}

func NewErrorData[T any](err error) DataObjectInterface[T] {
	return &DataObject[T]{
		_ID:   uuid.New().String(),
		Label: "error",
		Error: err,
	}
}

func NewUserData[T any](data T) DataObjectInterface[T] {
	return NewData[T]("user", data)
}
