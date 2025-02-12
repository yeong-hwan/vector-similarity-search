package main

import (
    "fmt"
    "search/models"
    "search/database"
    "search/dimension"
)

func main() {
    // 벡터 데이터베이스 생성
    db := database.NewVectorDB()

    // 상품 데이터 추가
    // 벡터 의미: [가격대(0-1), 캐주얼스타일(0-1), 포멀스타일(0-1), 스포티스타일(0-1), 계절감(0-1)]
    db.AddProduct(models.Product{
        ID:          1,
        Name:        "캐주얼 티셔츠",
        Price:       29900,
        Vector:      []float64{0.3, 0.9, 0.1, 0.4, 0.7}, // 저가, 매우 캐주얼, 여름
        Category:    "의류",
        Description: "편안한 데일리 티셔츠",
    })

    db.AddProduct(models.Product{
        ID:          2,
        Name:        "정장 셔츠",
        Price:       89000,
        Vector:      []float64{0.7, 0.1, 0.9, 0.1, 0.5}, // 중가, 매우 포멀, 무난
        Category:    "의류",
        Description: "고급 비즈니스 셔츠",
    })

    db.AddProduct(models.Product{
        ID:          3,
        Name:        "운동용 반팔",
        Price:       39900,
        Vector:      []float64{0.4, 0.3, 0.1, 0.9, 0.8}, // 저가, 매우 스포티, 여름
        Category:    "의류",
        Description: "기능성 스포츠 웨어",
    })

    db.AddProduct(models.Product{
        ID:          4,
        Name:        "캐주얼 셔츠",
        Price:       49900,
        Vector:      []float64{0.5, 0.8, 0.3, 0.2, 0.6}, // 중가, 캐주얼, 무난
        Category:    "의류",
        Description: "데일리 캐주얼 셔츠",
    })

    // 검색 시나리오 예시들
    fmt.Println("=== 다양한 검색 시나리오 ===\n")

    // 1. 캐주얼한 여름옷 찾기
    fmt.Println("1. {캐주얼한 여름옷} 검색:")
    summerCasualQuery := []float64{0.3, 0.9, 0.1, 0.2, 0.9} // 저가, 매우 캐주얼, 여름
    results := db.Search(summerCasualQuery, 2)
    printSearchResults(results)

    // 2. 고급 정장 스타일 찾기
    fmt.Println("\n2. {고급 정장 스타일} 검색:")
    formalQuery := []float64{0.8, 0.1, 0.9, 0.1, 0.5} // 고가, 매우 포멀
    results = db.Search(formalQuery, 2)
    printSearchResults(results)

    // 3. 적당한 가격의 스포츠웨어 찾기
    fmt.Println("\n3. {가성비 스포츠웨어} 검색:")
    sportsQuery := []float64{0.4, 0.2, 0.1, 0.9, 0.7} // 중저가, 매우 스포티
    results = db.Search(sportsQuery, 2)
    printSearchResults(results)

    // PCA 예시
    vectors := make([][]float64, len(db.GetProducts()))
    for i, product := range db.GetProducts() {
        vectors[i] = product.Vector
    }
    reducedVectors := dimension.PCA(vectors, 3)
    
    fmt.Println("차원 축소된 벡터들:")
    for i, vec := range reducedVectors {
        fmt.Printf("상품 %s의 축소된 벡터: %v\n", db.GetProducts()[i].Name, vec)
    }
}

// 검색 결과 출력을 위한 헬퍼 함수
func printSearchResults(results []models.SearchResult) {
    for _, result := range results {
        fmt.Printf("상품명: %s\n", result.Product.Name)
        fmt.Printf("가격: %.0f원\n", result.Product.Price)
        fmt.Printf("카테고리: %s\n", result.Product.Category)
        fmt.Printf("설명: %s\n", result.Product.Description)
        fmt.Printf("유사도: %.2f\n\n", result.Similarity)
    }
} 