package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Arnel-rah/eco-scale/config"
	"github.com/Arnel-rah/eco-scale/docker"
	"github.com/Arnel-rah/eco-scale/scheduler"
)

func main() {
	fmt.Println("=== Eco-Scale Daemon Started ===")

	scanner, err := docker.NewDockerScanner()
	if err != nil {
		log.Fatalf("Erreur Docker: %v", err)
	}
	defer scanner.Close()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-stopSignal:
			fmt.Println("\n=== Arrêt du démon Eco-Scale ===")
			return
		case <-ticker.C:
			cfg, err := config.LoadConfig("config/policy.yaml")
			if err != nil {
				fmt.Printf("Erreur rechargement config: %v\n", err)
				continue
			}

			activeContainers, err := scanner.ListActiveContainers()
			if err != nil {
				fmt.Printf("Erreur récupération conteneurs: %v\n", err)
				continue
			}

			actionsRequired := scheduler.AnalyzeSystem(cfg.Policies, activeContainers, cfg.AlertMode)
			if len(actionsRequired) == 0 {
				continue
			}

			fmt.Printf("[%s] --- Exécution de la routine de régulation ---\n", time.Now().Format("15:04:05"))
			for _, target := range actionsRequired {
				if target.Required == scheduler.ActionStop {
					fmt.Printf("-> Arrêt de %s...\n", target.Name)
					if err := scanner.StopContainer(target.ID); err != nil {
						fmt.Printf("Erreur arrêt: %v\n", err)
					} else {
						fmt.Printf("%s arrêté.\n", target.Name)
					}
				} else if target.Required == scheduler.ActionScale {
					fmt.Printf("-> Bridage de %s à %d%% CPU...\n", target.Name, target.Policy.CPULimit)
					if err := scanner.ScaleContainer(target.ID, target.Policy.CPULimit); err != nil {
						fmt.Printf("Erreur bridage: %v\n", err)
					} else {
						fmt.Printf("%s bridé.\n", target.Name)
					}
				}
			}
		}
	}
}
