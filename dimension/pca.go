package dimension

import (
	"math"
	"math/rand"
)

// PCA는 주성분 분석을 수행하여 차원을 축소합니다
func PCA(vectors [][]float64, targetDim int) [][]float64 {
	if len(vectors) == 0 || len(vectors[0]) == 0 || targetDim >= len(vectors[0]) {
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

	// 3. 공분산 행렬 계산
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

	// 4. 주성분 계산 (power iteration 방법)
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