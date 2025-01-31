package main

import (
    "fmt"
    "math"
    "sort"
)

// 벡터와 ID를 포함하는 구조체
type Vector struct {
    ID     int
    Values []float64
}

// 검색 결과를 저장하는 구조체
type SearchResult struct {
    ID        int
    Similarity float64
}

// 벡터 데이터베이스 구조체
type VectorDB struct {
    vectors []Vector
}

// 새로운 벡터 데이터베이스 생성
func NewVectorDB() *VectorDB {
    return &VectorDB{
        vectors: make([]Vector, 0),
    }
}

// 벡터 추가
func (db *VectorDB) AddVector(id int, values []float64) {
    db.vectors = append(db.vectors, Vector{
        ID:     id,
        Values: values,
    })
}

// 코사인 유사도 계산
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

// 가장 유사한 벡터 검색
func (db *VectorDB) Search(query []float64, topK int) []SearchResult {
    results := make([]SearchResult, 0)

    // 모든 벡터와의 유사도 계산
    for _, vec := range db.vectors {
        similarity := cosineSimilarity(query, vec.Values)
        results = append(results, SearchResult{
            ID:         vec.ID,
            Similarity: similarity,
        })
    }

    // 유사도 기준으로 정렬
    sort.Slice(results, func(i, j int) bool {
        return results[i].Similarity > results[j].Similarity
    })

    // topK 개수만큼 반환
    if len(results) > topK {
        results = results[:topK]
    }

    return results
}

func main() {
    // 벡터 데이터베이스 생성
    db := NewVectorDB()

    // 샘플 벡터 추가
    db.AddVector(1, []float64{1.0, 0.5, 0.3})
    db.AddVector(2, []float64{0.8, 0.2, 0.9})
    db.AddVector(3, []float64{0.1, 0.9, 0.4})

    // 검색 쿼리 벡터
    query := []float64{1.0, 0.5, 0.3}

    // 상위 2개 유사한 벡터 검색
    results := db.Search(query, 2)

    // 결과 출력
    for _, result := range results {
        fmt.Printf("ID: %d, Similarity: %.4f\n", result.ID, result.Similarity)
    }
}