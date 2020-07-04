package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gonum/stat"

	"encoding/csv"
	"encoding/hex"
	"encoding/json"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var columnsnames []string
var columnsNamesKNN []string
var columnsNamesKMeans []string

var columns map[string]int
var columnsKNN map[string]int
var columnsKMeans map[string]int

var fileData [][]float32
var dataframeKNN [][]float32
var dataframeKMeans [][]float32

var xDataTr [][]float32
var xDataTrainKNN [][]float32

var yDataTr [][]float32
var yDataTrainKNN [][]float32

var groupsSelected []informe

// ESTRUCTURAS PARA KNN

type Block struct {
	Index     int
	Timestamp string
	data      informe
	Hash      string
	PrevHash  string
}

var Blockchain []Block

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + fmt.Sprintf("%f", block.data.AgeGroup) + fmt.Sprintf("%f", block.data.Cancer) + fmt.Sprintf("%f", block.data.CardiovascularDisease) + fmt.Sprintf("%f", block.data.Diabetes) + fmt.Sprintf("%f", block.data.Hypertension) + fmt.Sprintf("%f", block.data.RespiratoryDisease) + fmt.Sprintf("%f", block.data.Sex) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, data informe) (Block, error) {
	var newBlock Block
	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.data = data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)
	return newBlock, nil
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

type informe struct {
	AgeGroup              float32 `json:"age_group"`
	Sex                   float32 `json:"sex"`
	CardiovascularDisease float32 `json:"cardiovascular_disease"`
	Diabetes              float32 `json:"diabetes"`
	RespiratoryDisease    float32 `json:"respiratory_disease"`
	Hypertension          float32 `json:"hypertension"`
	Cancer                float32 `json:"cancer"`
}

type knnResponse struct {
	ClassResult         int `json:"classResult"`
	DontHaveCoronavirus int `json:"dontHaveCoronavirus"`
	HaveCoronavirus     int `json:"haveCoronavirus"`
}

type res struct {
	Clase  int `json:"clase"`
	Ocurs0 int `json:"ocurs0"`
	Ocurs1 int `json:"ocurs1"`
}

// ESTRUCTURAS PARA KMEANS

type grupo struct {
	K     int `json:"clusters"`
	MaxIt int `json:"iteraciones"`
}

type kmeansResponse struct {
	Centroids []informe `json:"centroidsInforme"`
	Ncentroid []int     `json:"centroidsQuantity"`
}

var clase float32
var ocurs map[float32]int
var kmeansResultVar kmeansResponse

func readArchiveCSV(filePath string) ([]string, map[string]int, [][]float32) {

	fileArchive, err := os.Open(filePath)

	if err != nil {
		log.Fatal("No se puede leer el archivo de entrada "+filePath, err)
	}

	defer fileArchive.Close()
	csvReader := csv.NewReader(fileArchive)
	fileData, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("No se puede parsear el archivo de entrada "+filePath, err)
	}

	headers := make([]string, len(fileData[0]))
	copy(headers, fileData[0])

	columns := make(map[string]int)
	for i, header := range headers {
		columns[header] = i
	}

	fileData = fileData[1:]
	fileDataReal := make([][]float32, len(fileData))

	for i := range fileDataReal {
		fileDataReal[i] = make([]float32, len(headers))
		for j := range fileDataReal[i] {
			val, _ := strconv.ParseFloat(fileData[i][j], 32)
			fileDataReal[i][j] = float32(val)
		}
	}

	return headers, columns, fileDataReal
}

func splitColumns(headers []string, columns map[string]int, fileData [][]float32, newheaders []string) ([]string, map[string]int, [][]float32) {

	temp := make([]int, len(newheaders))
	newfileData := make([][]float32, len(fileData))

	for i, newh := range newheaders {
		temp[i] = columns[newh]
	}

	for i := range newfileData {
		newfileData[i] = make([]float32, len(temp))
		for j, t := range temp {
			newfileData[i][j] = fileData[i][t]
		}
	}

	newcolumns := make(map[string]int)

	for i, header := range newheaders {
		newcolumns[header] = i
	}

	return newheaders, newcolumns, newfileData

}

func head(fileData [][]float32, n int) {
	for i := 0; i < n; i++ {
		fmt.Println(fileData[i])
	}
}

func standardizeData(fileData [][]float32) ([][]float32, []float64, []float64) {

	newfileData := make([][]float32, len(fileData))
	for i := 0; i < len(fileData); i++ {
		newfileData[i] = make([]float32, len(fileData[i]))
	}

	mean := make([]float64, len(fileData[0]))
	std := make([]float64, len(fileData[0]))

	for i := 0; i < len(fileData[0]); i++ {
		columnsumn := make([]float64, len(fileData))
		mean[i] = float64(0)
		for j := 0; j < len(fileData); j++ {
			columnsumn[j] = float64(fileData[j][i])
			mean[i] += float64(fileData[j][i])
		}
		mean[i] = mean[i] / float64(len(fileData))
		std[i] = stat.StdDev(columnsumn, nil)
		for j := 0; j < len(fileData); j++ {
			newfileData[j][i] = (fileData[j][i] - float32(mean[i])) / float32(std[i])
		}
	}

	return newfileData, mean, std
}

func classKNN(xDataTr [][]float32, yDataTr [][]float32, xTest []float32, k int, outCh chan []float32) {
	distances := make([][]float32, k)
	for i := 0; i < len(distances); i++ {
		distances[i] = []float32{math.MaxFloat32, -1}
	}
	for i := 0; i < len(xDataTr); i++ {
		addValue := float32(0)
		for j := 0; j < len(xDataTr[i]); j++ {
			addValue += (xDataTr[i][j] - xTest[j]) * (xDataTr[i][j] - xTest[j])
		}
		addValue = float32(math.Sqrt(float64(addValue)))
		j := len(distances) - 1
		for ; j >= 0; j-- {
			if distances[j][0] <= addValue {
				temp := make([][]float32, k+1)
				for m := 0; m < j+1; m++ {
					temp[m] = make([]float32, 2)
					copy(temp[m], distances[m])
				}
				temp[j+1] = []float32{addValue, yDataTr[i][0]}
				for m := j + 1; m < k; m++ {
					temp[m+1] = make([]float32, 2)
					copy(temp[m+1], distances[m])
				}
				distances = temp[:k]
				break
			} else {
				if j == 0 {
					temp := make([][]float32, k+1)
					temp[0] = []float32{addValue, yDataTr[i][0]}
					for m := range distances {
						temp[m+1] = make([]float32, 2)
						copy(temp[m+1], distances[m])
					}
					distances = temp[:k]
				}
			}
		}
	}
	for _, dist := range distances {
		outCh <- dist
	}
}

func multiThreadKNN(xDataTr [][]float32, yDataTr [][]float32, xTest []float32, k int, chans int) (float32, map[float32]int) {
	outCh := make(chan []float32)
	xsize := len(xDataTr) / chans
	for i := 0; i < chans; i++ {
		if i < chans-1 {
			go classKNN(xDataTr[i*xsize:(i+1)*xsize], yDataTr[i*xsize:(i+1)*xsize], xTest, k, outCh)
		} else {
			go classKNN(xDataTr[i*xsize:], yDataTr[i*xsize:], xTest, k, outCh)
		}
	}

	distances := make([][]float32, k)
	for i := 0; i < len(distances); i++ {
		distances[i] = []float32{math.MaxFloat32, -1}
	}
	for i := 0; i < k*chans; i++ {
		j := len(distances) - 1
		candidate := <-outCh
		for ; j >= 0; j-- {
			if distances[j][0] <= candidate[0] {

				temp := make([][]float32, k+1)
				for m := 0; m < j+1; m++ {
					temp[m] = make([]float32, 2)
					copy(temp[m], distances[m])
				}
				temp[j+1] = []float32{candidate[0], candidate[1]}
				for m := j + 1; m < k; m++ {
					temp[m+1] = make([]float32, 2)
					copy(temp[m+1], distances[m])
				}
				distances = temp[:k]
				break
			} else {
				if j == 0 {
					temp := make([][]float32, k+1)
					temp[0] = []float32{candidate[0], candidate[1]}
					for m := range distances {
						temp[m+1] = make([]float32, 2)
						copy(temp[m+1], distances[m])
					}
					distances = temp[:k]
				}
			}
		}
	}
	close(outCh)
	clases := make(map[float32]int)
	for _, dist := range distances {
		if dist[1] == -1 {
			continue
		}
		if _, found := clases[dist[1]]; !found {
			clases[dist[1]] = 1
		} else {
			clases[dist[1]]++
		}
	}
	var res float32
	max := -1
	for key, val := range clases {
		fmt.Printf("Clase '%d': %d es la cantidad de ocurrencias que presenta\n", int(key), val)
		if val > max {
			max = val
			res = key
		}
	}

	return res, clases
}

func calculateNearCentroids(id int, dftemp [][]float32, centers [][]float32, GCh chan []int, idCh chan int) {

	n := len(dftemp)
	ncolumnss := len(dftemp[0])
	k := len(centers)

	dist := make([][]float32, k)
	for i := 0; i < k; i++ {
		dist[i] = make([]float32, n)
	}

	G := make([]int, n)
	for i, point := range dftemp {
		c := 0
		for cent := 0; cent < k; cent++ {
			addValue := float32(0)
			for columns := 0; columns < ncolumnss; columns++ {
				addValue += (point[columns] - centers[cent][columns]) * (point[columns] - centers[cent][columns])
			}
			dist[cent][i] = float32(math.Sqrt(float64(addValue)))
			if dist[cent][i] < dist[c][i] {
				c = cent
			}
		}
		G[i] = c
	}
	GCh <- G
	idCh <- id
}

func multiThreadTMeans(dftemp [][]float32, k int, maxIt int) ([]int, [][]float32, int) {
	n := len(dftemp)
	ncolumnss := len(dftemp[0])
	centers := make([][]float32, k)
	for i := 0; i < k; i++ {
		centers[i] = make([]float32, ncolumnss)
	}
	G := make([]int, n)
	rc := make(map[int]struct{})
	it := 0
	for len(rc) != k {
		r := rand.Intn(n - 1)
		rc[r] = struct{}{}
	}
	temp := 0
	for i := range rc {
		for j := 0; j < ncolumnss; j++ {
			centers[temp][j] = dftemp[i][j]
		}
		temp++
	}
	for repeat := false; !repeat && it < maxIt; {
		GCh := make(chan []int, 1)
		idCh := make(chan int, 1)
		tsize := n / k
		for i := 0; i < k; i++ {
			if i < k-1 {
				go calculateNearCentroids(i, dftemp[i*tsize:i*tsize+tsize], centers, GCh, idCh)
			} else {
				go calculateNearCentroids(i, dftemp[i*tsize:], centers, GCh, idCh)
			}
		}
		for i := 0; i < k; i++ {
			Gpiece := <-GCh
			id := <-idCh
			for j := 0; j < len(Gpiece); j++ {
				G[id*tsize+j] = Gpiece[j]
			}
		}
		close(GCh)
		end := make(chan bool)
		newcenters := make([][]float32, k)
		for i := 0; i < k; i++ {
			newcenters[i] = make([]float32, ncolumnss)
		}
		counters := make([]int, k)
		for i := 0; i < k; i++ {
			counters[i] = 0
		}
		for i := 0; i < k; i++ {
			go func(id int) {
				for i := id; i < n; i += k {
					for j := 0; j < ncolumnss; j++ {
						newcenters[G[i]][j] += dftemp[i][j]
					}
					counters[G[i]]++
				}
				end <- true
			}(i)
		}
		for i := 0; i < k; i++ {
			<-end
		}
		for i := 0; i < k; i++ {
			for j := 0; j < ncolumnss; j++ {
				newcenters[i][j] /= float32(counters[i])
			}
		}
		it++
		repeat = true
		for i := 0; i < k; i++ {
			for j := 0; j < ncolumnss; j++ {
				if centers[i][j] != newcenters[i][j] {
					repeat = false
					break
				}
			}

			copy(centers[i], newcenters[i])
		}
	}

	return G, centers, it
}

func classifyCovid(r http.ResponseWriter, request *http.Request) {

	var bdy informe

	err := json.NewDecoder(request.Body).Decode(&bdy)
	if err != nil {
		http.Error(r, err.Error(), http.StatusBadRequest)
		fmt.Println("Ha ocurrido un error")
		return
	}
	test := bdy

	newBlock, _ := generateBlock(Blockchain[len(Blockchain)-1], bdy)
	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
		fmt.Println("SE RECIBIERON LOS DATOS DE LA PERSONA")
		spew.Dump(newBlock)
	}

	K := 900
	Threads := 6

	xTest := []float32{test.AgeGroup, test.Sex, test.CardiovascularDisease, test.Diabetes, test.RespiratoryDisease,
		test.Hypertension, test.Cancer}

	dftemp := make([][]float32, len(xDataTrainKNN))
	for i := range xDataTrainKNN {
		dftemp[i] = make([]float32, len(xDataTrainKNN[i]))
		copy(dftemp[i], xDataTrainKNN[i])
	}
	dftemp = append(dftemp, xTest)

	dftemp, _, _ = standardizeData(dftemp)

	newdftemp := make([][]float32, len(yDataTrainKNN))
	for i := 0; i < len(newdftemp); i++ {
		newdftemp[i] = make([]float32, len(dftemp[0]))
		copy(newdftemp[i], dftemp[i])
	}
	copy(xTest, dftemp[len(dftemp)-1])

	clase, ocurs = multiThreadKNN(newdftemp, yDataTrainKNN, xTest, K, Threads)
	fmt.Println("La clase resultante es: ", clase)
}

func sendKNNResult(r http.ResponseWriter, request *http.Request) {
	response := res{Clase: int(clase), Ocurs0: ocurs[0], Ocurs1: ocurs[1]}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(r, err.Error(), http.StatusBadRequest)
		fmt.Println("Ha ocurrido un error")
		return
	}

	fmt.Fprintf(r, "%s", jsonResponse)
}

func clusteringCovid(r http.ResponseWriter, request *http.Request) {

	var bdy grupo

	err := json.NewDecoder(request.Body).Decode(&bdy)
	if err != nil {
		http.Error(r, err.Error(), http.StatusBadRequest)
		fmt.Println("Ha ocurrido un error")
		return
	}

	k := bdy.K

	dftemp := make([][]float32, len(dataframeKMeans))
	for i := range dataframeKMeans {
		dftemp[i] = make([]float32, len(dataframeKMeans[i]))
		copy(dftemp[i], dataframeKMeans[i])
	}

	dftemp, meanScalesVal, stdScalesVal := standardizeData(dftemp)
	maxIt := bdy.MaxIt

	fmt.Println("INICIA EJECUCIÓN DE KMEANS")
	head(dftemp, 1)
	G, centers, _ := multiThreadTMeans(dftemp, k, maxIt)
	head(centers, 1)

	ocurs := make(map[int]int)
	var ncentroid []int
	groupsSelected = nil

	for _, clase := range G {
		if _, found := ocurs[clase]; !found {
			ocurs[clase] = 1
		} else {
			ocurs[clase]++
		}
	}

	for _, val := range ocurs {
		ncentroid = append(ncentroid, val)
	}
	for i := 0; i < k; i++ {
		if !math.IsNaN(float64(centers[i][0])) {
			for j := 0; j < len(centers[i]); j++ {
				z := float64(centers[i][j])*stdScalesVal[j] + meanScalesVal[j]
				centers[i][j] = float32(z)
			}
			groupsSelected = append(groupsSelected, informe{AgeGroup: centers[i][0], Sex: centers[i][1], CardiovascularDisease: centers[i][2],
				Diabetes: centers[i][3], RespiratoryDisease: centers[i][4],
				Hypertension: centers[i][5], Cancer: centers[i][6]})
		}
	}
	fmt.Println("Centroides:")
	fmt.Println(groupsSelected)
	fmt.Println("Cantidad de pertenencia:")
	fmt.Println(ncentroid)
	kmeansResultVar = kmeansResponse{Centroids: groupsSelected, Ncentroid: ncentroid}

}

func kMeansResult(r http.ResponseWriter, request *http.Request) {
	jsonResponse, err := json.Marshal(kmeansResultVar)
	if err != nil {
		http.Error(r, err.Error(), http.StatusBadRequest)
		fmt.Println("Ha ocurrido un error")
		return
	}
	fmt.Fprintf(r, "%s", jsonResponse)
}

func main() {

	//Crear Blockchain
	rand.Seed(time.Now().UTC().UnixNano())
	genesisInforme := informe{3, 0, 1, 1, 0, 0, 1}
	t := time.Now()
	genesisBlock := Block{0, t.String(), genesisInforme, "3323a8545b60cc10d7d210ceeea5d0d6ce26aa6fc6365b9f52d93a0e4702973a", ""}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	//KnnFile
	columnsNamesKNN, columnsKNN, dataframeKNN = readArchiveCSV("muerteKNN.csv")
	_, _, xDataTrainKNN = splitColumns(columnsNamesKNN, columnsKNN, dataframeKNN, columnsNamesKNN[:len(columnsNamesKNN)-1])
	_, _, yDataTrainKNN = splitColumns(columnsNamesKNN, columnsKNN, dataframeKNN, []string{columnsNamesKNN[len(columnsNamesKNN)-1]})

	//KMeansFile
	columnsNamesKMeans, columnsKMeans, dataframeKMeans = readArchiveCSV("muerteKMeans.csv")

	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	router.HandleFunc("/clustering_covid", kMeansResult).Methods("GET")
	router.HandleFunc("/classification_covid", sendKNNResult).Methods("GET")
	router.HandleFunc("/classification_covid", classifyCovid).Methods("POST")
	router.HandleFunc("/clustering_covid", clusteringCovid).Methods("POST")

	fmt.Println("Now server is running on port 8000")
	http.ListenAndServe(":8000", handlers.CORS(headers, methods, origins)(router))
}
