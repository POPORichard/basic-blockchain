package server

type getblocks struct {
	AddrFrom string
}

type inv struct {
	AddrFrom string
	Type     string	//区别是块还是交易
	Items    [][]byte
}

//请求数据的就够体
type getdata struct {
	AddrFrom string
	Type     string
	ID       []byte
}

//实际传输block的结构体
type block struct {
	AddrFrom string
	Block    []byte
}
//实际传输tx的结构体
type tx struct {
	AddFrom     string
	Transaction []byte
}

type storefile struct {
	FileData []byte
}

