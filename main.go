package main

import (
	"fmt"
	"log"

	"github.com/Arnel-rah/eco-scale/config"
	"github.com/Arnel-rah/eco-scale/docker"
	"github.com/Arnel-rah/eco-scale/scheduler"
)

func main() {
	fmt.Println("=== Eco-Scale Daemon ===")

	cfg, err := config.LoadConfig("config/policy.yaml")
	if err != nil {
		log.Fatalf("Erreur lors du chargement de la config: %v", err)
	}

	fmt.Printf("Configuration version %s chargée avec succès !\n", cfg.Version)
	fmt.Printf("Mode alerte actuel : %v\n\n", cfg.AlertMode)

	scanner, err := docker.NewDockerScanner()
	if err != nil {
		log.Fatalf("Erreur Docker: %v", err)
	}
	defer scanner.Close()

	activeContainers, err := scanner.ListActiveContainers()
	if err != nil {
		log.Fatalf("Erreur lors de la récupération des conteneurs: %v", err)
	}

	fmt.Println("=== Analyse du Système ===")
	actionsRequired := scheduler.AnalyzeSystem(cfg.Policies, activeContainers, cfg.AlertMode)

	if len(actionsRequired) == 0 {
		fmt.Println("Aucune action requise.")
		return
	}

	for _, target := range actionsRequired {
		fmt.Printf("[CIBLE TROUVÉE] Conteneur: %s (%s) | Action requise: %s\n",
			target.Name,
			target.ID,
			target.Required,
		)

		if target.Required == scheduler.ActionStop {
			fmt.Printf("-> Arrêt du conteneur %s en cours...\n", target.Name)
			err := scanner.StopContainer(target.ID)
			if err != nil {
				fmt.Printf("Erreur lors de l'arrêt de %s: %v\n", target.Name, err)
			} else {
				fmt.Printf("Conteneur %s arrêté avec succès.\n", target.Name)
			}
		} else if target.Required == scheduler.ActionScale {
			fmt.Printf("-> Bridage du conteneur %s à %d%% CPU...\n", target.Name, target.Policy.CPULimit)
			err := scanner.ScaleContainer(target.ID, target.Policy.CPULimit)
			if err != nil {
				fmt.Printf("Erreur lors du bridage de %s: %v\n", target.Name, err)
			} else {
				fmt.Printf("Conteneur %s bridé avec succès à %d%% CPU.\n", target.Name, target.Policy.CPULimit)
			}
		}
	}
}
