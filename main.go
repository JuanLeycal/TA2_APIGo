package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var dataU = [][]string{}

var jsondata = []byte{}
var jsonKMean2 = []byte{}

type Afiliado struct {
	FECHA_CORTE            string `json:"fecha_corte"`
	REGION                 string `json:"region"`
	PROVINCIA              string `json:"provincia"`
	DISTRITO               string `json:"distrito"`
	UBIGEO                 string `json:"ubigeo"`
	COD_UNIDAD_EJECUTORA   string `json:"cod_unidad_ejecutora"`
	UNIDAD_EJECUTORA       string `json:"unidad_ejecutora"`
	AMBITO_INEI            string `json:"ambito_inei"`
	CODIGO_IPRESS          string `json:"codigo_ipress"`
	IPRESS                 string `json:"ipress"`
	VRAEM                  string `json:"vraem"`
	NACIONAL_EXTRANJERO    string `json:"nacional_extranjero"`
	PAIS_EXTRANJERO        string `json:"pais_extranjero"`
	DOCUMENTO_IDENTIDAD    string `json:"documento_identidad"`
	EDAD                   string `json:"edad"`
	SEXO                   string `json:"sexo"`
	REGIMEN_FINANCIAMIENTO string `json:"regimen_financiamiento"`
	PLAN_DE_SEGURO         string `json:"plan_de_seguro"`
	COBERTURA_FINANCIERA   string `json:"cobertura_financiera"`
	TOTAL_AFILIADOS        string `json:"total_afiliados"`
}

type KMean struct {
	X                int    `json:"edad"`
	Y                int    `json:"total_afiliados"`
	GroupId          int    `json:"group_id"`
	Unidad_ejecutora string `json:"unidad_ejecutora"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome the my GO API!")
}

func distance(p KMean, p2 KMean) int {
	first := math.Pow(float64(p2.X-p.X), 2)
	second := math.Pow(float64(p2.Y-p.Y), 2)
	final := math.Sqrt(first + second)
	var dist int = int(final)
	return dist
}

func KMeans(w http.ResponseWriter, r *http.Request) {

	rand.Seed(time.Now().UnixNano())

	ref := mux.Vars(r)
	ref1 := ref["indice1"]
	ref2 := ref["indice2"]
	variables := []KMean{}
	// if ref1 == "region" && ref2 == "plan" {
	// }
	if ref1 == "afiliados" && ref2 == "edad" {
		afilIndex := 19
		edadIndex := 14
		for _, each := range dataU {
			x, _ := strconv.Atoi(each[edadIndex])
			y, _ := strconv.Atoi(each[afilIndex])
			// edad := each[edadIndex]
			// afiliados := each[afilIndex]
			punto := KMean{
				X:                x,
				Y:                y,
				Unidad_ejecutora: each[6],
			}
			variables = append(variables, punto)
		}
		variables = variables[1:]
		sort.SliceStable(variables, func(i, j int) bool {
			return variables[i].X < variables[j].X
		})
		//fmt.Println(variables)
		//numero de clusters
		k := 3
		Centroids := []KMean{}
		for j := 0; j < k; j++ {
			i := rand.Intn(len(variables))
			fmt.Println(i)
			variables[i].GroupId = j + 1
			Centroids = append(Centroids, variables[i])
		}
		//fmt.Println(Centroids)
		//hora de comparar

		for i := 0; i < len(variables)-1; i++ {
			aux := 10000
			for j := 0; j < len(Centroids); j++ {
				if distance(variables[i], Centroids[j]) < aux && variables[i].GroupId == 0 {
					aux = distance(variables[i], Centroids[j])

					variables[i].GroupId = Centroids[j].GroupId
				}
			}
			//fmt.Print(i, "->", variables[i].GroupId, "- ")
		}

		for a := 0; a < 14; a++ {
			newClusters := []KMean{}
			for j := 0; j < k; j++ {
				aux := 0
				newCluster := KMean{}
				for i := 0; i < len(variables)-1; i++ {
					if variables[i].GroupId == j+1 {
						newCluster.X += variables[i].X
						newCluster.Y += variables[i].Y
						aux++
					}
				}
				newCluster.X = newCluster.X / aux
				newCluster.Y = newCluster.Y / aux
				newCluster.GroupId = j + 1
				newClusters = append(newClusters, newCluster)
			}
			// fmt.Print(Centroids)
			// fmt.Print(newClusters)

			for i := 0; i < len(variables)-1; i++ {
				aux := 10000
				for j := 0; j < len(newClusters); j++ {
					if distance(variables[i], newClusters[j]) < aux {
						aux = distance(variables[i], newClusters[j])
						variables[i].GroupId = newClusters[j].GroupId
					}
				}
				fmt.Print(i, "->", variables[i].GroupId, "- ")
			}

		}
		jsondata, _ = json.Marshal(variables)

		b, _ := ioutil.ReadFile(string(jsondata))

		rawIn := json.RawMessage(string(b))
		var objmap map[string]*json.RawMessage
		err := json.Unmarshal(rawIn, &objmap)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, string(jsondata))

	}
}

func AllData(w http.ResponseWriter, r *http.Request) {

	///response := Country{Name: "Perú", Capital: "Lima"}

	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(response)
	b, _ := ioutil.ReadFile("/dataset.json")

	rawIn := json.RawMessage(string(b))
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal(rawIn, &objmap)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(jsondata))

	//json.NewEncoder(w).Encode(objmap)

}

func main() {
	// read data from CSV file

	csvFile, erro := os.Open("./dataset.csv")

	if erro != nil {
		fmt.Println(erro)
	}

	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1
	csvData, erro := reader.ReadAll()
	dataU = csvData

	if erro != nil {
		fmt.Println(erro)
		os.Exit(1)
	}

	var oneRecord Afiliado
	var allRecords []Afiliado

	fmt.Println(dataU)

	for _, each := range csvData {
		oneRecord.FECHA_CORTE = each[0]
		oneRecord.REGION = each[1]
		oneRecord.PROVINCIA = each[2]
		oneRecord.DISTRITO = each[3]
		oneRecord.UBIGEO = each[4]
		oneRecord.COD_UNIDAD_EJECUTORA = each[5]
		oneRecord.UNIDAD_EJECUTORA = each[6]
		oneRecord.AMBITO_INEI = each[7]
		oneRecord.CODIGO_IPRESS = each[8]
		oneRecord.IPRESS = each[9]
		oneRecord.VRAEM = each[10]
		oneRecord.NACIONAL_EXTRANJERO = each[11]
		oneRecord.PAIS_EXTRANJERO = each[12]
		oneRecord.DOCUMENTO_IDENTIDAD = each[13]
		oneRecord.EDAD = each[14]
		oneRecord.SEXO = each[15]
		oneRecord.REGIMEN_FINANCIAMIENTO = each[16]
		oneRecord.PLAN_DE_SEGURO = each[17]
		oneRecord.COBERTURA_FINANCIERA = each[18]
		oneRecord.TOTAL_AFILIADOS = each[19]

		/*oneRecord.Name = each[0]
		oneRecord.Age, _ = strconv.Atoi(each[1]) // need to cast integer to string
		oneRecord.Job = each[2]*/
		allRecords = append(allRecords, oneRecord)
	}

	var err error

	jsondata, err = json.Marshal(allRecords) // convert to JSON

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// sanity check
	// NOTE : You can stream the JSON data to http service as well instead of saving to file

	fmt.Println(string(jsondata))

	// now write to JSON file

	//jsonFile.Write(jsondata)
	//jsonFile.Close()

	r := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/json", AllData).Methods("GET")
	r.HandleFunc("/json/{indice1}/{indice2}", KMeans).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", handler))

}
