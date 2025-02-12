package database

import (
    "sort"
    "search/models"
    "search/indexing"
    "search/utils"
)

// VectorDB는 벡터 데이터베이스를 관리하는 구조체입니다
type VectorDB struct {
    Products []models.Product
    lshIndex *indexing.LSHIndex
}

// NewVectorDB는 새로운 벡터 데이터베이스를 생성합니다
func NewVectorDB() *VectorDB {
    return &VectorDB{
        Products: make([]models.Product, 0),
        lshIndex: indexing.NewLSHIndex(5, 4, 2),
    }
}

// AddProduct는 상품을 데이터베이스에 추가합니다
func (db *VectorDB) AddProduct(product models.Product) {
    productIdx := len(db.Products)
    db.Products = append(db.Products, product)
    db.lshIndex.AddToTable(product.Vector, productIdx)
}

// GetProducts returns all products
func (db *VectorDB) GetProducts() []models.Product {
    return db.Products
}

// Search는 ANN(Approximate Nearest Neighbor) 검색을 수행합니다
func (db *VectorDB) Search(query []float64, topK int) []models.SearchResult {
	// LSH를 사용한 ANN 검색
	candidateSet := db.lshIndex.GetCandidates(query)
	results := make([]models.SearchResult, 0)

	// 후보 검색 (근사 이웃 검색)
	for idx := range candidateSet {
			similarity := utils.CosineSimilarity(query, db.Products[idx].Vector)
			results = append(results, models.SearchResult{
					Product:    db.Products[idx],
					Similarity: similarity,
			})
	}

	// 결과가 없으면 전체 검색 (폴백 메커니즘)
	if len(results) == 0 {
			for _, product := range db.Products {
					similarity := utils.CosineSimilarity(query, product.Vector)
					results = append(results, models.SearchResult{
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