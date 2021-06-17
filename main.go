package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var dataU = [][]string{}

var jsondata = []byte{}
var jsonKKN1 = []byte{}
var jsonKKN2 = []byte{}

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

type KKMean1 struct {
	REGION         string `json:"region"`
	PLAN_DE_SEGURO string `json:"plan_de_seguro"`
}
type KKMean2 struct {
	X                int    `json:"edad"`
	Y                int    `json:"total_afiliados"`
	GroupId          int    `json:"group_id"`
	Unidad_ejecutora string `json:"unidad_ejecutora"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome the my GO API!")
}

func distance(p KKMean2, p2 KKMean2) float64 {
	first := math.Pow(float64(p2.X-p.X), 2)
	second := math.Pow(float64(p2.Y-p.Y), 2)
	return math.Sqrt(first + second)
}

func KKMeans(w http.ResponseWriter, r *http.Request) {
	ref := mux.Vars(r)
	ref1 := ref["indice1"]
	ref2 := ref["indice2"]
	variables := []KKMean2{}
	//groupId := 0
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
			punto := KKMean2{
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
		fmt.Println(variables)
		//var distancias []int
		//KKGroups := []KKMean2{}
		groupContr := 1
		sumDist := 0
		variables[0].GroupId = groupContr
		for i := 0; i < len(variables)-1; i++ {
			dist := distance(variables[i], variables[i+1])
			sumDist += int(dist)
		}
		promedio := sumDist / len(variables)
		fmt.Println(promedio)
		for i := 0; i < len(variables)-1; i++ {
			dist := distance(variables[i], variables[i+1])
			if dist > float64(promedio) {
				groupContr++
				//KKGroups = append(KKGroups, variables[i+1])
			}
			variables[i+1].GroupId = groupContr
			//distancias = append(distancias, int(dist))
		}
		for _, each := range variables {
			fmt.Println(each.GroupId)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(variables)

		//fmt.Println(distancias)
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
	r.HandleFunc("/json/{indice1}/{indice2}", KKMeans).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", handler))

}
