package webserver

import (
	"encoding/json"
	"fmt"
	"goBlockchain/blockchain"
	"goBlockchain/imports/mux"
	"goBlockchain/transactions"
	"goBlockchain/utility"
	"goBlockchain/webserver/cors"
	"io"
	"log"
	"net/http"
	"os"
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

type ResultsHandler struct {
	PollTag string
	Choices []string
	User    string
}
type ResultsWriter struct {
	Results     map[string]int
	Voters      map[string][]string
	VoteBalance int
	VoteTarget  string
}
type Response struct {
	TxHash  string
	Message string
}

type AddTransaction struct {
	Sender    string
	Receiver  string
	Amount    int
	Tag       string
	Timestamp string
	Signature string
}

type NewToken struct {
	Tag       string
	Voters    []string
	Signature string
}

type ResAmount struct {
	Amount int
}

func Run() error {
	muxServer := makeMuxRouter()
	httpPort := os.Getenv("PORT")
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
	muxRouter.HandleFunc("/getBalance", handleGetBalance).Methods("GET")
	muxRouter.HandleFunc("/generateTokens", handleGenerateTokens).Methods("POST")
	muxRouter.HandleFunc("/addTransaction", handleAddTransaction).Methods("POST")
	muxRouter.HandleFunc("/mineBlock", handleMineBlock).Methods("POST")
	muxRouter.HandleFunc("/getResults", handleGetResults).Methods("POST")
	return muxRouter
}

func handleGetBalance(w http.ResponseWriter, r *http.Request) {
	var handler AddTransaction
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&handler); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	bc := blockchain.GetInstance()

	var results ResAmount

	results.Amount = bc.GetBalanceForAddress(handler.Sender, handler.Tag)

	fmt.Println("Balance for user " + handler.Sender + " equals " + strconv.Itoa(results.Amount) + " in poll " + handler.Tag)

	bytes, err := json.MarshalIndent(results.Amount, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleGetResults(w http.ResponseWriter, r *http.Request) {
	var handler ResultsHandler

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&handler); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	var results ResultsWriter
	results.Results = make(map[string]int, len(handler.Choices))
	results.Voters = make(map[string][]string, len(handler.Choices))

	for _, choice := range handler.Choices {
		results.Results[choice] = 0
	}

	bc := blockchain.GetInstance()

	results.VoteBalance = bc.GetBalanceForAddress(handler.User, handler.PollTag)
	results.VoteTarget = ""
	if results.VoteBalance == 0 {
		results.VoteTarget = bc.GetTargetForAddress(handler.User, handler.PollTag)
	}

	blocks := bc.TraverseBlockchain()

	for _, block := range blocks {
		txs := block.GetTransactions()

		for _, tx := range txs {
			if tx.GetTag() == handler.PollTag {
				if containedIn(tx.GetReceiver(), handler.Choices) {
					results.Results[tx.GetReceiver()]++
					results.Voters[tx.GetReceiver()] = append(results.Voters[tx.GetReceiver()], tx.GetSender())
				}
			}
		}
	}

	bytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))

}

func handleGenerateTokens(w http.ResponseWriter, r *http.Request) {
	var token NewToken

	decoder := json.NewDecoder(r.Body)
	fmt.println('RECEIVED GEN_TOKENS REQUEST ===============================================');
	fmt.println(r);
	if err := decoder.Decode(&token); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	bc := blockchain.GetInstance()

	hash := utility.Hash(strings.Join(token.Voters, ",")) // Calculates hash of all addresses that participate in poll
	fmt.println('BEFORE SENDING POST ===============================================');
	autherityAssurance := utility.PostRequest(bc.GetAuthorizedTokenGenerators()[0], token.Signature) // Sends the server authorized pubKey with the signature to assure it will equal the hash we calculated before, therefore assure that token generate request came from it.

	tx := transactions.Tx(bc.GetAuthorizedTokenGenerators()[0], strings.Join(token.Voters, ","), 0, token.Tag, time.Now().String())
	tx.Signature = token.Signature

	// Stores control transaction in separate block for later validation....
	controlId := strconv.Itoa(bc.MineControlBlock(tx))
	fmt.Println(autherityAssurance)
	fmt.Println(hash)
	if autherityAssurance == hash {
		for _, receiver := range token.Voters {
			tx := transactions.Tx("Generator", receiver, 1, token.Tag, time.Now().String())
			tx.Signature = controlId
			bc.AddTransaction(tx)
		}

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

	tx := transactions.Tx(add.Sender, add.Receiver, add.Amount, add.Tag, add.Timestamp)
	tx.Signature = add.Signature

	result := bc.AddTransaction(tx)
	res.TxHash = result[0]
	res.Message = result[1]
	//res.TxHash = tx.GetHash()

	//bc.AppendSignature(add.TxHash, add.Signature)

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

func containedIn(tag string, slice []string) bool {
	for _, value := range slice {
		if tag == value {
			return true
		}
	}
	return false
}
