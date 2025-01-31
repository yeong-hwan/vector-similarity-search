package main

import (
    "fmt"
    "math"
    "math/rand"
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

// LSH를 위한 구조체
type LSHIndex struct {
    hashTables []map[string][]int  // 해시 테이블들
    numTables  int                 // 해시 테이블 개수
    numBands   int                 // LSH 밴드 개수
    bandSize   int                 // 각 밴드의 크기
}

// LSH 인덱스 생성
func NewLSHIndex(numTables, numBands, bandSize int) *LSHIndex {
    hashTables := make([]map[string][]int, numTables)
    for i := range hashTables {
        hashTables[i] = make(map[string][]int)
    }
    return &LSHIndex{
        hashTables: hashTables,
        numTables:  numTables,
        numBands:   numBands,
        bandSize:   bandSize,
    }
}

// LSH 해시 함수
func (lsh *LSHIndex) hashVector(vector []float64, tableIdx int) string {
    rand.Seed(int64(tableIdx))
    hash := ""
    for i := 0; i < lsh.numBands; i++ {
        bandHash := 0
        for j := 0; j < lsh.bandSize && (i*lsh.bandSize+j) < len(vector); j++ {
            randVal := rand.Float64()
            if vector[i*lsh.bandSize+j] > randVal {
                bandHash = bandHash*2 + 1
            } else {
                bandHash = bandHash * 2
            }
        }
        hash += fmt.Sprintf("_%d", bandHash)
    }
    return hash
}

// PCA 차원 축소
func pca(vectors [][]float64, targetDim int) [][]float64 {
    if len(vectors) == 0 || targetDim >= len(vectors[0]) {
        return vectors
    }

    // 1. 평균 계산
    dim := len(vectors[0])
    mean := make([]float64, dim)
    for _, v := range vectors {
        for j := range v {
            mean[j] += v[j]
        }
    }
    for j := range mean {
        mean[j] /= float64(len(vectors))
    }

    // 2. 중심화
    centered := make([][]float64, len(vectors))
    for i, v := range vectors {
        centered[i] = make([]float64, dim)
        for j := range v {
            centered[i][j] = v[j] - mean[j]
        }
    }

    // 3. 공분산 행렬 계산 (간단한 버전)
    cov := make([][]float64, dim)
    for i := range cov {
        cov[i] = make([]float64, dim)
        for j := range cov[i] {
            for k := 0; k < len(vectors); k++ {
                cov[i][j] += centered[k][i] * centered[k][j]
            }
            cov[i][j] /= float64(len(vectors) - 1)
        }
    }

    // 4. 주성분 계산 (간단한 power iteration 방법)
    components := make([][]float64, targetDim)
    for i := range components {
        components[i] = make([]float64, dim)
        // 초기 벡터
        for j := range components[i] {
            components[i][j] = rand.Float64()
        }
        // Power iteration
        for iter := 0; iter < 100; iter++ {
            // 행렬-벡터 곱
            newVec := make([]float64, dim)
            for j := range newVec {
                for k := range cov[j] {
                    newVec[j] += cov[j][k] * components[i][k]
                }
            }
            // 정규화
            norm := 0.0
            for j := range newVec {
                norm += newVec[j] * newVec[j]
            }
            norm = math.Sqrt(norm)
            for j := range newVec {
                components[i][j] = newVec[j] / norm
            }
        }
    }

    // 5. 차원 축소된 데이터 계산
    reduced := make([][]float64, len(vectors))
    for i := range reduced {
        reduced[i] = make([]float64, targetDim)
        for j := 0; j < targetDim; j++ {
            for k := range vectors[i] {
                reduced[i][j] += (vectors[i][k] - mean[k]) * components[j][k]
            }
        }
    }

    return reduced
}

// VectorDB 구조체에 상품 정보 추가
type VectorDB struct {
    products []Product
    lshIndex *LSHIndex
}

func NewVectorDB() *VectorDB {
    return &VectorDB{
        products: make([]Product, 0),
        lshIndex: NewLSHIndex(5, 4, 2), // 예: 5개 테이블, 4개 밴드, 밴드당 2개 해시
    }
}

// AddProduct는 상품을 데이터베이스에 추가합니다
func (db *VectorDB) AddProduct(product Product) {
    productIdx := len(db.products)
    db.products = append(db.products, product)
    
    // LSH 인덱스에 추가
    for i := 0; i < db.lshIndex.numTables; i++ {
        hash := db.lshIndex.hashVector(product.Vector, i)
        db.lshIndex.hashTables[i][hash] = append(db.lshIndex.hashTables[i][hash], productIdx)
    }
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

// Search 함수 최적화
func (db *VectorDB) Search(query []float64, topK int) []SearchResult {
    // LSH를 사용하여 후보 상품 찾기
    candidateSet := make(map[int]bool)
    for i := 0; i < db.lshIndex.numTables; i++ {
        hash := db.lshIndex.hashVector(query, i)
        if candidates, exists := db.lshIndex.hashTables[i][hash]; exists {
            for _, idx := range candidates {
                candidateSet[idx] = true
            }
        }
    }

    // 후보 상품들에 대해서만 유사도 계산
    results := make([]SearchResult, 0)
    for idx := range candidateSet {
        similarity := cosineSimilarity(query, db.products[idx].Vector)
        results = append(results, SearchResult{
            Product:    db.products[idx],
            Similarity: similarity,
        })
    }

    // 결과가 없으면 모든 상품에 대해 검색
    if len(results) == 0 {
        for _, product := range db.products {
            similarity := cosineSimilarity(query, product.Vector)
            results = append(results, SearchResult{
                Product:    product,
                Similarity: similarity,
            })
        }
    }

    // 정렬 및 상위 K개 반환
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

    // 차원 축소 예시
    vectors := make([][]float64, len(db.products))
    for i, product := range db.products {
        vectors[i] = product.Vector
    }
    
    // 5차원을 3차원으로 축소
    reducedVectors := pca(vectors, 3)
    
    fmt.Println("\n차원 축소된 벡터들:")
    for i, vec := range reducedVectors {
        fmt.Printf("상품 %s의 축소된 벡터: %v\n", db.products[i].Name, vec)
    }
}