package webserver

import (
	"encoding/json"
	"goBlockchain/blockchain"
	"goBlockchain/imports/mux"
	"goBlockchain/security"
	"goBlockchain/transactions"
	"goBlockchain/webserver/cors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Message struct {
	From   string
	To     string
	Amount int
	Tag    string
}

func Run() error {
	muxServer := makeMuxRouter()
	httpPort := "8080"
	handler := cors.Default().Handler(muxServer)
	log.Println("Listening on ", httpPort)
	s := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        handler,
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
	muxRouter.HandleFunc("/transactions/{key}", handleGetTransactions).Methods("GET")
	muxRouter.HandleFunc("/merkle", handleGetMerkle).Methods("GET")
	muxRouter.HandleFunc("/signatures", handleGetSignatures).Methods("GET")
	muxRouter.HandleFunc("/pendingTransactions", handleGetPending).Methods("GET")
	muxRouter.HandleFunc("/addTransaction", handleAddTransaction).Methods("POST")
	muxRouter.HandleFunc("/mineBlock", handleMineBlock).Methods("POST")
	return muxRouter
}

func handleGetPending(w http.ResponseWriter, r *http.Request) {
	bc := blockchain.GetInstance()
	var txs []transactions.Transaction
	txs = append(txs, bc.GetPendingTransactions()...)

	bytes, err := json.MarshalIndent(txs, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bc := blockchain.GetInstance()
	blocks := bc.TraverseForwardBlockchain()
	bytes, err := json.MarshalIndent(blocks, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleGetMerkle(w http.ResponseWriter, r *http.Request) {
	bc := blockchain.GetInstance()

	it := bc.ForwardIterator()
	var v = make(map[string]map[int][]string)
	for {
		block := it.Next()

		v["BlockHash: "+block.GetHash()] = block.GetMerkleTree().PrintLevels()

		if block.GetId() == bc.GetLastBlock().GetId() {
			break
		}

	}

	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleGetSignatures(w http.ResponseWriter, r *http.Request) {
	bc := blockchain.GetInstance()
	sigs := bc.GetAllSignatures()
	bytes, err := json.MarshalIndent(sigs, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleGetTransactions(w http.ResponseWriter, r *http.Request) {
	bc := blockchain.GetInstance()
	blockId, _ := strconv.Atoi(mux.Vars(r)["key"])
	block := bc.GetBlockById(blockId)
	var txs []transactions.Transaction
	txs = append(txs, block.GetTransactions()...)

	bytes, err := json.MarshalIndent(txs, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleAddTransaction(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	respondWithJSON(w, r, http.StatusCreated, m)

	bc := blockchain.GetInstance()

	pubKey := security.GenerateKey(m.From)

	tx := transactions.Tx(pubKey, m.To, m.Amount, m.Tag)

	sign := security.Sign(tx.GetHash(), m.From)

	bc.AddTransaction(tx, sign)

	//newBlock, err := generateBlock(blockchain[len(blockchain)-1], m.BPM)
	//if err != nil {
	//	respondWithJSON(w, r, http.StatusInternalServerError, m)
	//	return
	//}
	//if isBlockValid(newBlock, blockchain[len(blockchain)-1]) {
	//	newBlockchain := append(blockchain, newBlock)
	//	replaceChain(newBlockchain)
	//	spew.Dump(blockchain)
	//}
	//
	//respondWithJSON(w, r, http.StatusCreated, newBlock)

}

func handleMineBlock(w http.ResponseWriter, r *http.Request) {

	respondWithJSON(w, r, http.StatusCreated, "Mining new Block\n")

	bc := blockchain.GetInstance()

	bc.MineNextBlock()

}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
