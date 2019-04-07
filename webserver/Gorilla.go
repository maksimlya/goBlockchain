package webserver

import (
	"encoding/json"
	"goBlockchain/blockchain"
	"goBlockchain/imports/mux"
	"goBlockchain/transactions"
	"goBlockchain/utility"
	"goBlockchain/webserver/cors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	From   string
	To     string
	Amount int
	Tag    string
}
type Response struct {
	TxHash string
}

type AddTransaction struct {
	TxHash    string
	Signature string
}

type NewToken struct {
	Tag       string
	Voters    []string
	Signature string
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
	muxRouter.HandleFunc("/txAmount", handleGetTxAmount).Methods("GET")
	muxRouter.HandleFunc("/generateTokens", handleGenerateTokens).Methods("POST")
	muxRouter.HandleFunc("/newTransaction", handleNewTransaction).Methods("POST")
	muxRouter.HandleFunc("/addTransaction", handleAddTransaction).Methods("POST")
	muxRouter.HandleFunc("/mineBlock", handleMineBlock).Methods("POST")
	return muxRouter
}

func handleGenerateTokens(w http.ResponseWriter, r *http.Request) {
	var token NewToken

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&token); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	bc := blockchain.GetInstance()

	hash := utility.Hash(strings.Join(token.Voters, "")) // Calculates hash of all addresses that participate in poll

	autherityAssurance := utility.PostRequest(bc.GetAuthorizedTokenGenerators()[0], token.Signature) // Sends the server authorized pubKey with the signature to assure it will equal the hash we calculated before, therefore assure that token generate request came from it.

	if autherityAssurance == hash {
		for _, receiver := range token.Voters {
			tx := transactions.Tx("Generator", receiver, 1, token.Tag)
			bc.AddTransaction(tx)
		}

		tx := transactions.Tx(bc.GetAuthorizedTokenGenerators()[0], strings.Join(token.Voters, ","), 0, token.Tag)
		tx.Signature = token.Signature
		bc.AddTransaction(tx)
	}

	//fmt.Println(bc.GetAuthorizedTokenGenerators());

	//tx := transactions.Tx(token.Tag,strings.Join(token.Voters,","),0,"Control");
	//bc.AddTransaction(tx);

	//bc.AppendSignature(add.TxHash, add.Signature)

	respondWithJSON(w, r, 200, "Success")

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

func handleGetTxAmount(w http.ResponseWriter, r *http.Request) {
	bc := blockchain.GetInstance()

	txAmount := bc.GetTxAmount()

	bytes, err := json.MarshalIndent(txAmount, "", "  ")
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
	var add AddTransaction
	var res Response
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&add); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	//respondWithJSON(w, r, http.StatusCreated, m)

	bc := blockchain.GetInstance()
	res.TxHash = add.TxHash

	bc.AppendSignature(add.TxHash, add.Signature)

	respondWithJSON(w, r, 200, res)

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

func handleNewTransaction(w http.ResponseWriter, r *http.Request) {
	var m Message
	var res Response
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	//respondWithJSON(w, r, http.StatusCreated, m)

	bc := blockchain.GetInstance()

	pubKey := m.From

	tx := transactions.Tx(pubKey, m.To, m.Amount, m.Tag) // TODO - check if the address has unused votes.
	res.TxHash = tx.GetHash()

	//sign := security.Sign(tx.GetHash(), pubKey)

	bc.AddTransaction(tx)

	respondWithJSON(w, r, 200, res)

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

//func handleAddTransaction(w http.ResponseWriter, r *http.Request) {
//	var m Message
//
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&m); err != nil {
//		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
//		return
//	}
//	defer r.Body.Close()
//
//	//respondWithJSON(w, r, http.StatusCreated, m)
//
//	bc := blockchain.GetInstance()
//
//
//	pubKey := security.GenerateKey(m.From)
//
//	tx := transactions.Tx(pubKey, m.To, m.Amount, m.Tag)
//
//	sign := security.Sign(tx.GetHash(), m.From)
//
//	bc.AddTransaction(tx, sign)
//
//	respondWithJSON(w, r, 200, m)
//
//
//	//newBlock, err := generateBlock(blockchain[len(blockchain)-1], m.BPM)
//	//if err != nil {
//	//	respondWithJSON(w, r, http.StatusInternalServerError, m)
//	//	return
//	//}
//	//if isBlockValid(newBlock, blockchain[len(blockchain)-1]) {
//	//	newBlockchain := append(blockchain, newBlock)
//	//	replaceChain(newBlockchain)
//	//	spew.Dump(blockchain)
//	//}
//	//
//	//respondWithJSON(w, r, http.StatusCreated, newBlock)
//
//}

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
