package word

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/bluele/mecab-golang"
)

func ReturnKeyWords(texts []string) [][]string {
	n := len(texts)
	fmt.Println(n)
	var Dec []string
	outDec := make([][]string, n)
	var s []int
	D := texts
	BoW := make([][]int, n)
	for i := 0; i < n; i++ {
		BoW[i] = make([]int, 10000)
	}
	KeyOfDj := make([][]float64, n)
	for i := range KeyOfDj {
		KeyOfDj[i] = make([]float64, 10000)
	}
	lastindex := 0
	m, err := mecab.New("-Owakati")
	if err != nil {
		panic(err)
	}
	defer m.Destroy()
	for i := 0; i < len(D); i++ {
		text := D[i]
		Dec = parseToDec(m, text, Dec, BoW[i], &lastindex)
	}

	if err != nil {
		panic(err)
	}
	CreateKOD(Dec, BoW, KeyOfDj)
	for i := range texts {
		DecOut := make([]string, 0)
		s = OutKeyWordOfD(Dec, KeyOfDj, i, 5)
		for j := 0; j < len(s); j++ {
			DecOut = append(DecOut, Dec[s[j]])
		}
		log.Println(s, DecOut)
		outDec[i] = DecOut
	}
	log.Println("outDec : ")
	log.Println(outDec)
	return outDec
}

func CreateKOD(Dec []string, BoW [][]int, KeyOfDj [][]float64) {
	for i := range BoW {
		for j := range Dec {
			KeyOfDj[i][j] = Tfidf(Dec[j], Dec, BoW, i)
		}
	}
}

func OutKeyWordOfD(Dec []string, KeyOfDj [][]float64, Dj int, N int) []int {
	RankIndex := make([]int, N)
	for i := 0; i < N; i++ {
		RankIndex[i] = 0
	}
	RankValue := make([]float64, N)
	for i := 0; i < N; i++ {
		RankValue[i] = -1
	}
	for i := range Dec {
		for j := range RankIndex {
			if RankValue[j] < KeyOfDj[Dj][i] {
				for k := len(RankIndex) - 1; k > j; k-- {
					RankValue[k] = RankValue[k-1]
					RankIndex[k] = RankIndex[k-1]
				}
				RankValue[j] = KeyOfDj[Dj][i]
				RankIndex[j] = i
				break
			} else if RankValue[j] == KeyOfDj[Dj][i] {
				for k := len(RankIndex) - 1; k-1 > j; k-- {
					RankValue[k] = RankValue[k-1]
					RankIndex[k] = RankIndex[k-1]
				}
				if j == N-1 {
					break
				}
				RankValue[j+1] = KeyOfDj[Dj][i]
				RankIndex[j+1] = i
				break
			}
		}
	}
	//fmt.Println(RankValue)
	//fmt.Println(RankIndex)
	for i := range RankIndex {
		fmt.Printf("%s ", Dec[RankIndex[i]])
	}
	fmt.Println("")
	return RankIndex
}

func parseToDec(m *mecab.MeCab, text string, Dec []string, BoWj []int, lastindex *int) []string {
	///表層形\t品詞,品詞細分類1,品詞細分類2,品詞細分類3,活用型,活用形,原形,読み,発音
	const BOSEOS = "BOS/EOS"

	tg, err := m.NewTagger()
	if err != nil {
		fmt.Println(err)
	}
	defer tg.Destroy()

	lt, err := m.NewLattice(text)
	if err != nil {
		fmt.Println(err)
	}
	defer lt.Destroy()

	node := tg.ParseToNode(lt)
	for {
		features := strings.Split(node.Feature(), ",")
		if features[0] != BOSEOS {
			if features[0] == "名詞" /*|| features[0] == "動詞" || features[0] == "形容詞" */ {
				Text := strings.Split(node.Feature(), ",")
				index := arrayContains(Dec, Text[6])
				if index != -1 {
					BoWj[index] = BoWj[index] + 1
				} else {
					NewDec := append(Dec, Text[6])
					Dec = NewDec
					BoWj[*lastindex] = BoWj[*lastindex] + 1
					*lastindex++
				}
			}
		}
		if node.Next() != nil {
			break
		}
	}
	return Dec
}

func arrayContains(arr []string, str string) int {
	for n, v := range arr {
		if v == str {
			return n
		}
	}
	return -1
}

// Tfidf returns t's TF-IDF value in ds
func Tfidf(t string, Dec []string, BoW [][]int, Dj int) float64 {
	return Tf(t, Dec, BoW, Dj) * (Idf(t, Dec, BoW) + 1)
}

// Tf returns t's TF value in d
func Tf(t string, Dec []string, BoW [][]int, Dj int) float64 {
	var count int
	n := 0
	index := arrayContains(Dec, t)
	if index != -1 {
		count = BoW[Dj][index]
	} else {
		return 0
	}
	for i := 0; i < 10000; i++ {
		n += BoW[Dj][i]
	}
	return float64(count) / float64(n)
}

// Idf returns t's IDF value in ds
func Idf(t string, Dec []string, BoW [][]int) float64 {
	D := len(BoW)
	n := 0
	index := arrayContains(Dec, t)
	if index != -1 {
		for i := range BoW {
			if BoW[i][index] != 0 {
				n++
			}
		}
	} else {
		return 0
	}
	return math.Log(float64(D) / float64(n))
}
