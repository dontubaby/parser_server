package models

type Hex struct {
	Coordinate [2]int //у каждого гекса 2 координаты Q и R
	Htype      string // тип покрытия гекса
}

type Zone struct {
	Id        string
	Hexs      []Hex  //пока берем за основу что в зоне не больше 127 гексов
	Modifiers string // бафы зоны
}

type Sector struct {
	Id    string
	Zones []Zone //число зон точно не больше числа гексов в зоне (127)
}

type AutobattlerMap struct {
	Id      string
	Sectors []Sector //карта состоит из 3 секторов
}

type Unit struct {
	Id         string
	Utype      string
	Subspecies string
	GroupID    int64
}

type UnitGroup struct {
	Id    string
	Units []Unit
}
