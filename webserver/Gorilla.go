package webserver

import (
	"encoding/json"
	"goBlockchain/blockchain"
	"goBlockchain/transactions"

	//"goBlockchain/blockchain"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Run() error {
	mux := makeMuxRouter()
	httpPort := "8080"
	log.Println("Listening on ", httpPort)
	s := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/transactions", handleGetTransactions).Methods("GET")
	muxRouter.HandleFunc("/merkle", handleGetMerkle).Methods("GET")
	muxRouter.HandleFunc("/signatures", handleGetSignatures).Methods("GET")
	//muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	blockchain := blockchain.InitBlockchain()
	blocks := blockchain.TraverseForwardBlockchain()
	bytes, err := json.MarshalIndent(blocks, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleGetMerkle(w http.ResponseWriter, r *http.Request) {
	blockchain := blockchain.InitBlockchain()
	block1 := blockchain.GetBlockById(1)

	var v = make(map[string]map[int][]string, len(block1.GetMerkleTree().PrintLevels()))

	v["BlockHash: "+block1.GetHash()] = block1.GetMerkleTree().PrintLevels()

	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleGetSignatures(w http.ResponseWriter, r *http.Request) {
	blockchain := blockchain.InitBlockchain()
	sigs := blockchain.GetAllSignatures()
	bytes, err := json.MarshalIndent(sigs, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleGetTransactions(w http.ResponseWriter, r *http.Request) {
	blockchain := blockchain.InitBlockchain()
	blocks := blockchain.TraverseForwardBlockchain()
	var txs []transactions.Transaction
	for i := range blocks {
		txs = append(txs, blocks[i].GetTransactions()...)
	}
	bytes, err := json.MarshalIndent(txs, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

//func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
//	var m Message
//
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&m); err != nil {
//		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
//		return
//	}
//	defer r.Body.Close()
//
//	newBlock, err := generateBlock(blockchain[len(blockchain)-1], m.BPM)
//	if err != nil {
//		respondWithJSON(w, r, http.StatusInternalServerError, m)
//		return
//	}
//	if isBlockValid(newBlock, blockchain[len(blockchain)-1]) {
//		newBlockchain := append(blockchain, newBlock)
//		replaceChain(newBlockchain)
//		spew.Dump(blockchain)
//	}
//
//	respondWithJSON(w, r, http.StatusCreated, newBlock)
//
//}

//func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
//	response, err := json.MarshalIndent(payload, "", "  ")
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		w.Write([]byte("HTTP 500: Internal Server Error"))
//		return
//	}
//	w.WriteHeader(code)
//	w.Write(response)
//}
