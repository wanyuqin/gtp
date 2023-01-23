package example

type User struct {
	InterfaceTest interface{}       `json:"interface_test"`
	Id            int               `json:"id"`            // id
	Name          string            `json:"name"`          // 名称
	OrderId       []string          `json:"order_id"`      // 订单ID数组
	Pet           map[string]string `json:"pet"`           // 宠物map
	Float32Test   float32           `json:"float_32_test"` // float32 测试
	Float64Test   float64           `json:"float_64_test"` // float64 测试
	BoolTest      bool              `json:"bool_test"`     // 布尔测试
	Int8Test      int8              `json:"int_8_test"`
	Int16Test     int16             `json:"int_16_test"`
	Int32Test     int32             `json:"int_32_test"`
	Int64Test     int64             `json:"int_64_test"`
	UintTest      uint              `json:"uint_test"`
	Uint8Test     uint8             `json:"uint_8_test"`
	Uint16Test    uint16            `json:"uint_16_test"`
	Uint32Test    uint32            `json:"uint_32_test"`
	Uint64Test    uint64            `json:"uint_64_test"`
	ByteTest      byte              `json:"byte_test"`
	BytesTest     []byte            `json:"bytes_test"`
}
