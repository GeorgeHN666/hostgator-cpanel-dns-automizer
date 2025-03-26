package main

import (
	"context"
	"dns-automizer/pkg/DNS"
	"dns-automizer/pkg/IP"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("./env/.env")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	if os.Getenv("ENV") == "DEV" {
		err := godotenv.Load("./env/dev.local.env")
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	} else {
		err := godotenv.Load("./env/prod.local.env")
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	}
}

func main() {

	interval, err := strconv.Atoi(os.Getenv("INTERVAL"))
	if err != nil {
		log.Fatal("--- MUST BE A VALID INTERVAL ---", err.Error())
		return
	}

	fmt.Printf("--- STARTING AUTOMATIC DNS SERVICE v%s ---\n", os.Getenv("VERSION"))
	ctx, cancel := context.WithCancel(context.Background())
	go StartDDNSService(ctx, interval)

	// Wait for SIGNALS
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	// Gracefully shutdown:
	cancel()
}

func StartDDNSService(ctx context.Context, interval int) {

	ticker := time.NewTicker(time.Duration(interval) * time.Minute)
	defer ticker.Stop()

	for {

		select {
		case <-ticker.C:
			fmt.Println("--------------------------------------------------")
			fmt.Println("--- CHECKING PUBLIC IP CHANGES ---")

			updatedAddr, oldAddr, match, err := IP.StartIPService().StartIPComprobation()
			if err != nil {
				fmt.Printf("⚠️ IP check failed (retrying in 20m): %v\n", err)
				continue // Skips to next iteration (keeps ticker alive)
			}

			if !match {
				fmt.Printf("🔄 PUBLIC IP CHANGED: %s → %s\n", oldAddr, updatedAddr)
				if err := DNS.StartDNSService().StartRecordUpdate(updatedAddr); err != nil {
					fmt.Printf("⚠️ DNS update failed (retrying in 20m): %v\n", err)
					continue
				}
			}

			fmt.Println("--------------------------------------------------")
			fmt.Println("✅ IP CHECK COMPLETE — NEXT CHECK IN 20 MINUTES")
		case <-ctx.Done(): // Triggered when ctx is cancelled
			fmt.Println("🛑 DDNS service stopped gracefully")
			return
		}

	}

}
