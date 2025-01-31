package main

import (
    "fmt"
    "math"
    "sort"
)

// Product는 상품 정보를 담는 구조체입니다
type Product struct {
    ID          int
    Name        string
    Price       float64
    Vector      []float64
    Category    string
    Description string
}

// VectorDB 구조체에 상품 정보 추가
type VectorDB struct {
    products []Product
}

func NewVectorDB() *VectorDB {
    return &VectorDB{
        products: make([]Product, 0),
    }
}

// AddProduct는 상품을 데이터베이스에 추가합니다
func (db *VectorDB) AddProduct(product Product) {
    db.products = append(db.products, product)
}

// SearchResult는 검색 결과를 저장하는 구조체입니다
type SearchResult struct {
    Product    Product
    Similarity float64
}

// 코사인 유사도 계산 함수는 그대로 유지
func cosineSimilarity(a, b []float64) float64 {
    if len(a) != len(b) {
        return 0
    }

    var dotProduct, normA, normB float64
    for i := 0; i < len(a); i++ {
        dotProduct += a[i] * b[i]
        normA += a[i] * a[i]
        normB += b[i] * b[i]
    }

    if normA == 0 || normB == 0 {
        return 0
    }

    return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// Search 함수 수정
func (db *VectorDB) Search(query []float64, topK int) []SearchResult {
    results := make([]SearchResult, 0)

    for _, product := range db.products {
        similarity := cosineSimilarity(query, product.Vector)
        results = append(results, SearchResult{
            Product:    product,
            Similarity: similarity,
        })
    }

    sort.Slice(results, func(i, j int) bool {
        return results[i].Similarity > results[j].Similarity
    })

    if len(results) > topK {
        results = results[:topK]
    }

    return results
}

func main() {
    // 벡터 데이터베이스 생성
    db := NewVectorDB()

    // 상품 데이터 추가
    // 벡터 의미: [가격대(0-1), 캐주얼스타일(0-1), 포멀스타일(0-1), 스포티스타일(0-1), 계절감(0-1)]
    db.AddProduct(Product{
        ID:          1,
        Name:        "캐주얼 티셔츠",
        Price:       29900,
        Vector:      []float64{0.3, 0.9, 0.1, 0.4, 0.7}, // 저가, 매우 캐주얼, 여름
        Category:    "의류",
        Description: "편안한 데일리 티셔츠",
    })

    db.AddProduct(Product{
        ID:          2,
        Name:        "정장 셔츠",
        Price:       89000,
        Vector:      []float64{0.7, 0.1, 0.9, 0.1, 0.5}, // 중가, 매우 포멀
        Category:    "의류",
        Description: "고급 비즈니스 셔츠",
    })

    db.AddProduct(Product{
        ID:          3,
        Name:        "운동용 반팔",
        Price:       39900,
        Vector:      []float64{0.4, 0.3, 0.1, 0.9, 0.8}, // 저가, 매우 스포티, 여름
        Category:    "의류",
        Description: "기능성 스포츠 웨어",
    })

    db.AddProduct(Product{
        ID:          4,
        Name:        "캐주얼 셔츠",
        Price:       49900,
        Vector:      []float64{0.5, 0.8, 0.3, 0.2, 0.6}, // 중가, 캐주얼
        Category:    "의류",
        Description: "데일리 캐주얼 셔츠",
    })

    // 검색 예시 1: 캐주얼한 여름 의류 찾기
    fmt.Println("\n캐주얼한 여름 의류 검색:")
    query1 := []float64{0.3, 0.8, 0.1, 0.3, 0.8} // 저가, 캐주얼, 여름 스타일
    results1 := db.Search(query1, 2)
    for _, result := range results1 {
        fmt.Printf("상품명: %s\n", result.Product.Name)
        fmt.Printf("가격: %.0f원\n", result.Product.Price)
        fmt.Printf("설명: %s\n", result.Product.Description)
        fmt.Printf("유사도: %.2f\n\n", result.Similarity)
    }

    // 검색 예시 2: 포멀한 의류 찾기
    fmt.Println("\n포멀한 의류 검색:")
    query2 := []float64{0.7, 0.1, 0.9, 0.1, 0.5} // 중가, 포멀 스타일
    results2 := db.Search(query2, 2)
    for _, result := range results2 {
        fmt.Printf("상품명: %s\n", result.Product.Name)
        fmt.Printf("가격: %.0f원\n", result.Product.Price)
        fmt.Printf("설명: %s\n", result.Product.Description)
        fmt.Printf("유사도: %.2f\n\n", result.Similarity)
    }
}