package models

// Product는 상품 정보를 담는 구조체입니다
type Product struct {
    ID          int
    Name        string
    Price       float64
    Vector      []float64
    Category    string
    Description string
}

// SearchResult는 검색 결과를 저장하는 구조체입니다
type SearchResult struct {
    Product    Product
    Similarity float64
} 