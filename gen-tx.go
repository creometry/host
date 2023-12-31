//v4: fetch variable values from cm
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/sealed-secrets/pkg/client/clientset/versioned"
)

// Ethereum node connection details
var (
	nodeURL     string
	contractAddr string
	gasLimit    *big.Int
)

const (
	contractABI = `[{"constant":false,"inputs":[{"name":"walletAddresses","type":"address[]"},{"name":"amounts","type":"uint256[]"}],"name":"distributeTokens","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
)

func main() {
	// Read Ethereum connection details from ConfigMap mounted as a file
	readEthereumConfig()

	client, err := rpc.Dial(nodeURL)
	if err != nil {
		log.Fatal(err)
	}

	// Construct contract ABI
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatal(err)
	}

	// Fetch private key from Kubernetes Secret
	privateKey, err := getPrivateKeyFromSecret("your-namespace", "your-secret-name", "private-key-field-name")
	if err != nil {
		log.Fatal(err)
	}

	// Construct private key and unlock account
	auth, err := bindPrivateKey(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Interact with the Ethereum contract
	callOpts := &ethereum.CallMsg{}
	transactOpts := &bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: gasLimit}

	// Fetch Kubernetes node metrics
	nodeMetrics, err := getKubernetesNodeMetrics("your-namespace", "wallet_ID")
	if err != nil {
		log.Fatal(err)
	}

	// Aggregate metrics based on wallet_ID
	aggregatedMetrics := aggregateMetrics(nodeMetrics)

	// Create a new Ethereum client
	contractAddr := common.HexToAddress(contractAddr)
	instance, err := NewTokenDistribution(contractAddr, client)
	if err != nil {
		log.Fatal(err)
	}

	// Call the distributeTokens function
	addresses := make([]common.Address, 0, len(aggregatedMetrics))
	amounts := make([]*big.Int, 0, len(aggregatedMetrics))

	for walletID, metric := range aggregatedMetrics {
		addresses = append(addresses, common.HexToAddress(walletID))
		amounts = append(amounts, big.NewInt(metric))
	}

	_, err = instance.DistributeTokens(transactOpts, addresses, amounts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transaction successful!")
}

func readEthereumConfig() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	configMapName := "ethereum-configmap"
	namespace := "your-namespace"

	configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	nodeURL = readConfigMapValue(configMap, "nodeURL")
	contractAddr = readConfigMapValue(configMap, "contractAddr")

	gasLimitStr := readConfigMapValue(configMap, "gasLimit")
	gasLimit, success := new(big.Int).SetString(gasLimitStr, 10)
	if !success {
		log.Fatal("Invalid gasLimit value in ConfigMap")
	}
}

func readConfigMapValue(configMap *v1.ConfigMap, key string) string {
	value, exists := configMap.Data[key]
	if !exists {
		log.Fatalf("Key %s not found in ConfigMap", key)
	}
	return value
}

// ... (rest of the script remains unchanged)


//v3
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/sealed-secrets/pkg/client/clientset/versioned"
)

// Ethereum node connection details
const (
	nodeURL        = "http://localhost:8545" // Update with your Ethereum node URL
	contractAddr   = "0xYourContractAddress"  // Replace with your deployed contract address
	contractABI    = `[{"constant":false,"inputs":[{"name":"walletAddresses","type":"address[]"},{"name":"amounts","type":"uint256[]"}],"name":"distributeTokens","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
)

func main() {
	client, err := rpc.Dial(nodeURL)
	if err != nil {
		log.Fatal(err)
	}

	// Construct contract ABI
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatal(err)
	}

	// Fetch private key from Kubernetes Secret
	privateKey, err := getPrivateKeyFromSecret("your-namespace", "your-secret-name", "private-key-field-name")
	if err != nil {
		log.Fatal(err)
	}

	// Construct private key and unlock account
	auth, err := bindPrivateKey(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Interact with the Ethereum contract
	callOpts := &ethereum.CallMsg{}
	transactOpts := &bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: big.NewInt(2000000)}

	// Fetch Kubernetes node metrics
	nodeMetrics, err := getKubernetesNodeMetrics("your-namespace", "wallet_ID")
	if err != nil {
		log.Fatal(err)
	}

	// Aggregate metrics based on wallet_ID
	aggregatedMetrics := aggregateMetrics(nodeMetrics)

	// Create a new Ethereum client
	contractAddr := common.HexToAddress(contractAddr)
	instance, err := NewTokenDistribution(contractAddr, client)
	if err != nil {
		log.Fatal(err)
	}

	// Call the distributeTokens function
	addresses := make([]common.Address, 0, len(aggregatedMetrics))
	amounts := make([]*big.Int, 0, len(aggregatedMetrics))

	for walletID, metric := range aggregatedMetrics {
		addresses = append(addresses, common.HexToAddress(walletID))
		amounts = append(amounts, big.NewInt(metric))
	}

	_, err = instance.DistributeTokens(transactOpts, addresses, amounts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transaction successful!")
}

func aggregateMetrics(nodeMetrics map[string]Metric) map[string]int64 {
	aggregatedMetrics := make(map[string]int64)

	for _, metric := range nodeMetrics {
		walletID := metric.WalletID
		amount := metric.Usage.CpuUsage + metric.Usage.MemoryUsage

		if _, exists := aggregatedMetrics[walletID]; exists {
			aggregatedMetrics[walletID] += amount
		} else {
			aggregatedMetrics[walletID] = amount
		}
	}

	return aggregatedMetrics
}

type Metric struct {
	WalletID string
	Usage    struct {
		CpuUsage    int64
		MemoryUsage int64
	}
}

func getKubernetesNodeMetrics(namespace, labelName string) (map[string]Metric, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	nodeMetrics := make(map[string]Metric)

	for _, node := range nodes.Items {
		walletID := node.GetLabels()[labelName]
		metric, err := getNodeMetrics(clientset, node.GetName())
		if err != nil {
			log.Printf("Error fetching metrics for node %s: %v", node.GetName(), err)
			continue
		}

		nodeMetrics[node.GetName()] = Metric{
			WalletID: walletID,
			Usage:    metric,
		}
	}

	return nodeMetrics, nil
}

func getNodeMetrics(clientset *kubernetes.Clientset, nodeName string) (Metric, error) {
	// Implement logic to fetch and return metrics for a Kubernetes node
	// You may use clientset to interact with metrics APIs or other tools
	// Extract and return CPU and memory usage
	return Metric{}, nil
}

// v2 
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/sealed-secrets/pkg/client/clientset/versioned"
)

// Ethereum node connection details
const (
	nodeURL        = "http://localhost:8545" // Update with your Ethereum node URL
	contractAddr   = "0xYourContractAddress"  // Replace with your deployed contract address
	contractABI    = `[{"constant":false,"inputs":[{"name":"walletAddresses","type":"address[]"},{"name":"amounts","type":"uint256[]"}],"name":"distributeTokens","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
	walletAddress1 = "0xWalletAddress1"       // Replace with the target wallet addresses
	walletAddress2 = "0xWalletAddress2"
	amount1        = 100                       // Replace with the corresponding amounts
	amount2        = 150
)

func main() {
	client, err := rpc.Dial(nodeURL)
	if err != nil {
		log.Fatal(err)
	}

	// Construct contract ABI
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatal(err)
	}

	// Fetch private key from Kubernetes Secret
	privateKey, err := getPrivateKeyFromSecret("your-namespace", "your-secret-name", "private-key-field-name")
	if err != nil {
		log.Fatal(err)
	}

	// Construct private key and unlock account
	auth, err := bindPrivateKey(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Interact with the Ethereum contract
	callOpts := &ethereum.CallMsg{}
	transactOpts := &bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: big.NewInt(2000000)}

	// Create a new Ethereum client
	contractAddr := common.HexToAddress(contractAddr)
	instance, err := NewTokenDistribution(contractAddr, client)
	if err != nil {
		log.Fatal(err)
	}

	// Call the distributeTokens function
	_, err = instance.DistributeTokens(transactOpts, []common.Address{common.HexToAddress(walletAddress1), common.HexToAddress(walletAddress2)}, []*big.Int{big.NewInt(amount1), big.NewInt(amount2)})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transaction successful!")
}

func getPrivateKeyFromSecret(namespace, secretName, keyName string) (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return "", err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	sealedSecretClient, err := versioned.NewForConfig(config)
	if err != nil {
		return "", err
	}

	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	sealedSecret, err := sealedSecretClient.SealedSecrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// Assuming the secret is base64 encoded
	encodedPrivateKey := secret.Data[keyName]
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(encodedPrivateKey)
	if err != nil {
		return "", err
	}

	return string(decodedPrivateKey), nil
}

// v1
package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

// Ethereum node connection details
const (
	nodeURL        = "http://localhost:8545" // Update with your Ethereum node URL
	contractAddr   = "0xYourContractAddress"  // Replace with your deployed contract address
	privateKey     = "0xYourPrivateKey"       // Replace with your private key
	contractABI    = `[{"constant":false,"inputs":[{"name":"walletAddresses","type":"address[]"},{"name":"amounts","type":"uint256[]"}],"name":"distributeTokens","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
	walletAddress1 = "0xWalletAddress1"       // Replace with the target wallet addresses
	walletAddress2 = "0xWalletAddress2"
	amount1        = 100                       // Replace with the corresponding amounts
	amount2        = 150
)

func main() {
	client, err := rpc.Dial(nodeURL)
	if err != nil {
		log.Fatal(err)
	}

	// Construct contract ABI
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatal(err)
	}

	// Construct private key and unlock account
	auth, err := bindPrivateKey(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Interact with the Ethereum contract
	callOpts := &ethereum.CallMsg{}
	transactOpts := &bind.TransactOpts{From: auth.From, Signer: auth.Signer, GasLimit: big.NewInt(2000000)}

	// Create a new Ethereum client
	contractAddr := common.HexToAddress(contractAddr)
	instance, err := NewTokenDistribution(contractAddr, client)
	if err != nil {
		log.Fatal(err)
	}

	// Call the distributeTokens function
	_, err = instance.DistributeTokens(transactOpts, []common.Address{common.HexToAddress(walletAddress1), common.HexToAddress(walletAddress2)}, []*big.Int{big.NewInt(amount1), big.NewInt(amount2)})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transaction successful!")
}

func bindPrivateKey(privateKey string) (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactor(strings.NewReader(privateKey))
	if err != nil {
		return nil, err
	}

	return auth, nil
}
