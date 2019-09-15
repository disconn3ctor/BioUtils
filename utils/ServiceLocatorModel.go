package utils

type ServiceLocatorModel struct {
	ProductService struct {
		ID      string        `json:"ID"`
		Service string        `json:"Service"`
		Tags    []interface{} `json:"Tags"`
		Meta    struct {
		} `json:"Meta"`
		Port    int    `json:"Port"`
		Address string `json:"Address"`
		Weights struct {
			Passing int `json:"Passing"`
			Warning int `json:"Warning"`
		} `json:"Weights"`
		EnableTagOverride bool `json:"EnableTagOverride"`
	} `json:"product-service"`
	Test struct {
		ID      string        `json:"ID"`
		Service string        `json:"Service"`
		Tags    []interface{} `json:"Tags"`
		Meta    struct {
		} `json:"Meta"`
		Port    int    `json:"Port"`
		Address string `json:"Address"`
		Weights struct {
			Passing int `json:"Passing"`
			Warning int `json:"Warning"`
		} `json:"Weights"`
		EnableTagOverride bool `json:"EnableTagOverride"`
	} `json:"test"`
}
