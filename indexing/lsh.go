package indexing

import (
	"fmt"
	"math/rand"
)

// LSHIndex는 LSH 인덱싱을 위한 구조체입니다
type LSHIndex struct {
	hashTables []map[string][]int
	numTables  int
	numBands   int
	bandSize   int
}

// NewLSHIndex는 새로운 LSH 인덱스를 생성합니다
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

// HashVector는 벡터를 해시값으로 변환합니다
func (lsh *LSHIndex) HashVector(vector []float64, tableIdx int) string {
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

// AddToTable은 벡터를 해시 테이블에 추가합니다
func (lsh *LSHIndex) AddToTable(vector []float64, idx int) {
	for i := 0; i < lsh.numTables; i++ {
		hash := lsh.HashVector(vector, i)
		lsh.hashTables[i][hash] = append(lsh.hashTables[i][hash], idx)
	}
}

// GetCandidates는 쿼리 벡터와 유사한 후보들을 반환합니다
func (lsh *LSHIndex) GetCandidates(query []float64) map[int]bool {
	candidateSet := make(map[int]bool)
	for i := 0; i < lsh.numTables; i++ {
		hash := lsh.HashVector(query, i)
		if candidates, exists := lsh.hashTables[i][hash]; exists {
			for _, idx := range candidates {
				candidateSet[idx] = true
			}
		}
	}
	return candidateSet
} 